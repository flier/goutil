//go:build go1.22

package arena

import (
	"math/bits"

	"github.com/flier/goutil/pkg/xunsafe"
	"github.com/flier/goutil/pkg/xunsafe/layout"
)

// Recycled is an arena allocator that reuses freed segments.
//
// It embeds an underlying Arena to satisfy new allocations and maintains
// per-size-class free lists to return previously released memory quickly.
// Size classes are indexed by log2 of the rounded-up request size to the
// arena alignment boundary.
//
// Implementation details:
//   - Released blocks are threaded into a single-linked list, using the first
//     machine word of the block as the "next" pointer. This keeps metadata
//     overhead minimal.
//   - On allocation, if a block is available in the matching size class, it is
//     popped from the list and cleared (zeroed) before being returned.
//   - If the current arena chunk does not have enough space for the request,
//     any trailing capacity is broken into powers-of-two sized blocks and fed
//     back into the free lists, reducing external fragmentation.
//   - Zero-sized allocations are delegated to the embedded Arena.
//   - Releasing a block with a size smaller than Align is ignored to avoid
//     managing tiny fragments.
type Recycled struct {
	Arena

	free []xunsafe.Addr[byte]
}

// Free releases a value of type T previously allocated from the given Recycled
// arena back to its size-class free list.
// The size of T is derived from layout metadata.
func Free[T any](a *Recycled, p *T) {
	size := layout.Of[T]().Size

	a.Release(xunsafe.Cast[byte](p), size)
}

// Release returns a previously allocated block back to the recycler's free list
// for its size class. The provided size is rounded up to Align before selecting
// a class. Blocks smaller than Align are ignored.
// The first machine word of the block is overwritten to store the next pointer
// in the per-class single-linked list.
func (a *Recycled) Release(p *byte, size int) {
	if size < Align {
		return
	}

	log := sizeClassIndex(alignUp(size))

	// Initialize free slice if needed
	a.ensureFreeList()

	*xunsafe.Cast[*uintptr](p) = xunsafe.Cast[uintptr](a.free[log].AssertValid())

	a.free[log] = xunsafe.AddrOf(xunsafe.Cast[byte](p))
}

// Alloc returns size bytes, first attempting to reuse a recycled block
// from the appropriate size class. Recycled blocks are zeroed before being
// returned. If no recycled block is available, the request is delegated to
// the embedded Arena. When the active chunk cannot satisfy the request,
// any trailing capacity is split into power-of-two blocks and recycled.
// A size of zero is delegated to the embedded Arena.
func (a *Recycled) Alloc(size int) *byte {
	// Handle zero size allocation
	if size == 0 {
		return a.Arena.Alloc(size)
	}

	if a.free != nil {
		log := sizeClassIndex(alignUp(size))

		if p := a.free[log].AssertValid(); p != nil {
			a.free[log] = xunsafe.Addr[byte](*xunsafe.Cast[uintptr](p))

			xunsafe.Clear(p, 1<<log)

			return p
		}
	}

	if a.Next != 0 && a.Next.Add(size) > a.End {
		n := int(a.End - a.Next)

		// Initialize free slice if needed
		a.ensureFreeList()

		for n > Align {
			log := sizeClassIndex(n)

			a.free[log] = xunsafe.AddrOf(a.Next.AssertValid())

			a.Next.Add(1 << log)

			n -= 1 << log
		}
	}

	return a.Arena.Alloc(size)
}

// Reset clears all recycled free lists and resets the embedded Arena.
// After Reset, released blocks are no longer tracked by the recycler, but the
// underlying arena memory may still be reused by future allocations.
// Any pointers into memory managed by the arena must not be used after Reset.
func (a *Recycled) Reset() {
	// Clear all recycled pointers
	for i := range a.free {
		a.free[i] = 0
	}
	a.Arena.Reset()
}

// alignUp rounds the size up to the arena alignment boundary.
func alignUp(size int) int {
	size += Align - 1
	size &^= Align - 1
	return size
}

// sizeClassIndex computes the size-class index (log2) for an aligned size.
func sizeClassIndex(size int) int { // size must be > 0 and aligned
	log := bits.Len(uint(size) - 1)
	sz := 1 << log
	if sz > size {
		log--
	}

	return log
}

const freeListCapacity = 64

// ensureFreeList lazily initializes the free-list slice.
func (a *Recycled) ensureFreeList() {
	if a.free == nil {
		a.free = make([]xunsafe.Addr[byte], freeListCapacity)
	}
}
