//go:build go1.22

package arena

import (
	"math/bits"

	"github.com/flier/goutil/pkg/xunsafe"
)

// Recycled is an arena allocator that reuses freed memory segments to improve
// performance and reduce memory fragmentation.
//
// It extends the base [Arena] allocator with a sophisticated recycling mechanism
// that maintains per-size-class free lists. When memory is released, it's
// categorized by size and stored in appropriate free lists for quick reuse.
// This approach significantly reduces the need for new memory allocations
// and helps mitigate external fragmentation.
//
// # Key Features
//
//   - Size-Class Management: Memory blocks are organized into power-of-two
//     size classes, indexed by log2 of the aligned size. This provides fast
//     O(1) lookup for recycled blocks.
//   - Zero-Copy Recycling: Released blocks are threaded into single-linked
//     lists using the first machine word of each block as the "next" pointer,
//     minimizing metadata overhead.
//   - Memory Clearing: All recycled blocks are automatically zeroed before
//     being returned, ensuring clean memory state for new allocations.
//   - Fragmentation Reduction: When the current arena chunk cannot satisfy
//     a request, trailing capacity is split into power-of-two blocks and
//     recycled, reducing external fragmentation.
//   - Alignment Handling: All allocations are aligned to the arena's
//     alignment boundary (typically 8 bytes on 64-bit systems).
//
// # Usage Pattern
//
//   - Allocate memory using Alloc() or the generic New() function
//   - Release memory using Release() or the generic Free() function
//   - Reset the allocator using Reset() when done with a batch of allocations
//
// # Performance Characteristics
//
//   - Fast Allocation: O(1) allocation when recycled blocks are available
//   - Efficient Recycling: O(1) release operations with minimal overhead
//   - Memory Efficiency: Reduces external fragmentation through size-class
//     management and trailing capacity recycling
//   - Cache Friendly: Size-class organization improves cache locality
//
// # Memory Safety
//
//   - All pointers must be released before calling Reset()
//   - Memory is automatically cleared on reuse to prevent data leakage
//   - Small allocations (< Align) are ignored by Release() to avoid
//     managing tiny fragments
//   - Zero-sized allocations are delegated to the underlying Arena
//
// # Example
//
//	arena := &Recycled{}
//
//	// Allocate some memory
//	ptr1 := arena.Alloc(64)
//	ptr2 := arena.Alloc(128)
//
//	// Release memory back to the recycler
//	arena.Release(ptr1, 64)
//	arena.Release(ptr2, 128)
//
//	// New allocations will reuse the released memory
//	newPtr1 := arena.Alloc(64)  // May reuse ptr1's memory
//	newPtr2 := arena.Alloc(128) // May reuse ptr2's memory
//
//	// Reset when done with the batch
//	arena.Reset()
type Recycled struct {
	Arena

	// free maintains per-size-class free lists for recycled memory blocks.
	// Each index represents a size class (log2 of the aligned size).
	// A nil entry means no recycled blocks are available for that size class.
	// The slice is lazily initialized when first needed.
	free []xunsafe.Addr[byte]
}

var _ Allocator = (*Recycled)(nil)

// Alloc allocates size bytes of memory, prioritizing recycled blocks from
// the appropriate size class when available.
//
// # Allocation Strategy
//
//  1. Zero-Size Handling: Zero-sized allocations are delegated to the
//     underlying Arena allocator.
//  2. Recycled Block Reuse: If a recycled block of the appropriate size
//     class is available, it's removed from the free list, cleared to zero,
//     and returned immediately.
//  3. Fragmentation Reduction: If the current arena chunk cannot satisfy
//     the request, any trailing capacity is split into power-of-two blocks
//     and recycled into the appropriate free lists.
//  4. Fallback Allocation: If no recycled blocks are available, the
//     request is delegated to the underlying Arena allocator.
//
// # Memory Clearing
//
// All recycled blocks are automatically cleared to zero using xunsafe.Clear
// before being returned. This ensures that new allocations start with a
// clean memory state, preventing data leakage from previous uses.
//
// # Size Class Optimization
//
// The allocator automatically determines the optimal size class for each
// request by rounding up to the alignment boundary and computing the
// corresponding power-of-two size class. This minimizes internal
// fragmentation while maintaining fast allocation performance.
//
// # Example
//
//	// Allocate 64 bytes - may reuse previously released memory
//	ptr1 := arena.Alloc(64)
//
//	// Allocate 128 bytes - may reuse previously released memory
//	ptr2 := arena.Alloc(128)
//
//	// Zero-sized allocation - delegated to underlying Arena
//	ptr3 := arena.Alloc(0)
//
// # Performance Notes
//
//   - O(1) allocation when recycled blocks are available
//   - Automatic memory clearing ensures clean state
//   - Trailing capacity recycling reduces external fragmentation
//   - Size class indexing provides fast block lookup
//
// Do not use this method directly, use [New] instead.
func (a *Recycled) Alloc(size int) *byte {
	// Handle zero size allocation
	if size == 0 {
		return a.Arena.Alloc(size)
	}

	if a.free != nil {
		alignedSize := alignUp(size)
		log := sizeClassIndex(alignedSize)

		if p := a.free[log].AssertValid(); p != nil {
			a.free[log] = xunsafe.Addr[byte](*xunsafe.Cast[uintptr](p))

			xunsafe.Clear(p, 1<<log)

			a.Log("reuse", "%v:%v, %d:%d", p, a.next, alignedSize, Align)

			return p
		}
	}

	if a.next != 0 && a.next.Add(size) > a.end {
		n := int(a.end - a.next)

		// Initialize free slice if needed
		a.ensureFreeList()

		for n > Align {
			log := sizeClassIndex(n)

			a.free[log] = xunsafe.AddrOf(a.next.AssertValid())

			a.next.Add(1 << log)

			n -= 1 << log
		}
	}

	return a.Arena.Alloc(size)
}

// Release returns a previously allocated memory block back to the recycler's
// free list for its corresponding size class.
//
// The provided size is automatically rounded up to the arena's alignment
// boundary (Align) before determining the appropriate size class. Blocks
// smaller than Align are ignored to avoid managing tiny fragments that
// would provide minimal benefit.
//
// The first machine word of the released block is overwritten to store
// the next pointer in the per-class single-linked list. This approach
// keeps metadata overhead minimal while maintaining efficient list traversal.
//
// # Size Class Calculation
//
// The size class is determined by:
//   - Rounding the size up to Align boundary
//   - Computing log2 of the aligned size
//   - Using this as an index into the free list array
//
// # Memory Safety
//
//   - The pointer p must point to valid memory previously allocated by
//     this Recycled allocator
//   - The size must match the actual allocated size
//   - The memory block should not be accessed after release
//
// # Example
//
//	ptr := arena.Alloc(64)
//	// ... use ptr ...
//	arena.Release(ptr, 64)  // Returns to size class 6 (64 bytes)
//
//	ptr2 := arena.Alloc(32)
//	// ... use ptr2 ...
//	arena.Release(ptr2, 32) // Returns to size class 5 (32 bytes, aligned to 40)
//
// Do not use this method directly, use [Free] instead.
func (a *Recycled) Release(p *byte, size int) {
	if p == nil || size < Align {
		return
	}

	alignedSize := alignUp(size)
	log := sizeClassIndex(alignedSize)

	// Initialize free slice if needed
	a.ensureFreeList()

	*xunsafe.Cast[*uintptr](p) = xunsafe.Cast[uintptr](a.free[log].AssertValid())

	a.free[log] = xunsafe.AddrOf(xunsafe.Cast[byte](p))

	a.Log("release", "%v:%v, %d:%d", p, a.next, alignedSize, Align)
}

// Reset clears all recycled free lists and resets the underlying Arena
// allocator to its initial state.
//
// # Behavior
//
// After calling Reset:
//   - All recycled free lists are cleared, discarding all previously
//     released memory blocks
//   - The underlying Arena is reset, potentially reusing its memory
//     blocks for future allocations
//   - Any pointers into memory managed by the arena become invalid
//     and must not be accessed
//
// # Use Cases
//
// Reset is typically called when:
//   - Starting a new batch of allocations
//   - Clearing all state between different phases of computation
//   - Preparing the allocator for reuse in a different context
//
// # Memory Safety Warning
//
// Critical: All pointers into memory managed by this Recycled allocator
// must be released or become invalid before calling Reset. Accessing memory
// after Reset will result in undefined behavior and potential crashes.
//
// # Example
//
//	// First batch of allocations
//	ptr1 := arena.Alloc(64)
//	ptr2 := arena.Alloc(128)
//	arena.Release(ptr1, 64)
//	arena.Release(ptr2, 128)
//
//	// Reset for new batch
//	arena.Reset()
//
//	// ptr1 and ptr2 are now invalid
//	// New allocations will start fresh
//	newPtr := arena.Alloc(64)
//
// # Performance Impact
//
// Reset is a relatively expensive operation that:
//   - Clears all free list entries (O(freeListCapacity))
//   - Resets the underlying Arena (complexity depends on Arena implementation)
//   - Discards all recycled memory blocks
//
// Use Reset judiciously, typically at natural boundaries in your
// allocation patterns rather than after every individual allocation.
func (a *Recycled) Reset() {
	// Clear all recycled pointers
	for i := range a.free {
		a.free[i] = 0
	}
	a.Arena.Reset()
}

// alignUp rounds the size up to the arena alignment boundary.
// This ensures all allocations are properly aligned for optimal
// performance and memory access patterns.
func alignUp(size int) int {
	size += Align - 1
	size &^= Align - 1
	return size
}

// sizeClassIndex computes the size-class index (log2) for an aligned size.
// The size must be greater than 0 and properly aligned to Align.
// This function maps allocation sizes to the appropriate free list index
// for efficient block lookup and management.
func sizeClassIndex(size int) int { // size must be > 0 and aligned
	log := bits.Len(uint(size) - 1)
	sz := 1 << log
	if sz > size {
		log--
	}

	return log
}

// freeListCapacity defines the maximum number of size classes that can be
// managed by the Recycled allocator. This limits the maximum allocation
// size to 2^(freeListCapacity-1) bytes while providing efficient
// size class indexing for common allocation patterns.
const freeListCapacity = 64

// ensureFreeList lazily initializes the free-list slice when first needed.
// This approach avoids unnecessary memory allocation for Recycled instances
// that are created but never used for recycling operations.
func (a *Recycled) ensureFreeList() {
	if a.free == nil {
		a.free = make([]xunsafe.Addr[byte], freeListCapacity)
	}
}
