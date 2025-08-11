//go:build go1.22

// Package arena provides a simple memory arena allocator for Go, inspired by the article
// [Cheating the Reaper in Go].
//
// The Arena type allows efficient allocation of memory for short-lived objects,
// reducing GC pressure by allocating memory in large chunks and freeing all allocations
// at once when the arena is reset or finalized.
//
// Usage:
//
//	a := new(arena.Arena)
//	p := arena.Alloc[MyStruct](a)
//
// The package uses unsafe operations and reflection to manage memory chunks,
// and sync.Pool to reuse memory efficiently.
//
// Note: This package is experimental and relies on unsafe features. Use with caution.
//
// Types:
//   - [Arena]: Represents a memory arena for fast allocation.
//   - Alloc[T]: Allocates an object of type T in the arena.
//
// Functions:
//   - [Arena.Alloc]: Allocates a chunk of memory with the specified size and alignment.
//   - [Arena.Reset]: Resets the arena, releasing all allocations.
//
// [Cheating the Reaper in Go]: https://mcyoung.xyz/2025/04/21/go-arenas/
package arena

import (
	"math/bits"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

type Allocator interface {
	Alloc(size, align uintptr) unsafe.Pointer
}

// New allocates memory for a value of type T from the provided Arena,
// initializes it with the given value v, and returns a pointer to the allocated value.
//
// The memory is allocated with the size and alignment required for type T.
func New[T any](a *Arena, v T) (p *T) {
	p = Alloc[T](a)
	*p = v

	return
}

// Alloc allocates memory for a value of type T from the provided Arena and returns a pointer to it.
//
// The allocated memory is properly sized and aligned for type T.
//
// Note: The returned pointer points to uninitialized memory.
func Alloc[T any](a *Arena) (p *T) {
	p = (*T)(a.Alloc(unsafe.Sizeof(*p), unsafe.Alignof(*p)))

	return
}

// Realloc reallocates memory for an object of type T to an object of type R within the given Arena.
//
// It takes a pointer to the original object and returns a pointer to the newly allocated object of type R.
// The function uses the Arena's Realloc method to handle the memory reallocation, adjusting the size and alignment as needed.
func Realloc[R, T any](a *Arena, p *T) (r *R) {
	return (*R)(a.Realloc(unsafe.Pointer(p), unsafe.Sizeof(*p), unsafe.Sizeof(*r), unsafe.Alignof(*r)))
}

type Arena struct {
	next      uintptr
	left, cap uintptr
	chunks    []unsafe.Pointer
}

func (a *Arena) Empty() bool { return a.next == 0 }
func (a *Arena) Reset()      { a.next, a.left, a.cap = 0, 0, 0 }

const (
	maxAlign uintptr = 8 // Depends on target, this is for 64-bit.
	minWords uintptr = 8
)

func (a *Arena) Alloc(size, align uintptr) unsafe.Pointer {
	// First, round the size up to the alignment of every object in the arena.
	mask := maxAlign - 1
	size = (size + mask) &^ mask

	// Then, replace the size with the size in pointer-sized words. This does not
	// result in any loss of size, since size is now a multiple of the uintptr
	// size.
	words := size / maxAlign

	// Next, check if we have enough space left for this chunk. If there isn't,
	// we need to grow.
	if a.left < words {
		// Pick whichever is largest: the minimum allocation size, twice the last
		// allocation, or the next power of two after words.
		a.cap = max(minWords, a.cap*2, nextPow2(words))
		p := a.allocChunk(a.cap)
		a.next = uintptr(p)
		a.left = a.cap
	}

	// Allocate the chunk by incrementing the pointer.
	p := a.next
	a.next += size
	a.left -= words

	return unsafe.Pointer(p) //nolint:govet
}

func (a *Arena) Realloc(ptr unsafe.Pointer, oldSize, newSize, align uintptr) unsafe.Pointer {
	// First, round the size up to the alignment of every object in the arena.
	mask := maxAlign - 1
	oldSize = (oldSize + mask) &^ mask
	newSize = (newSize + mask) &^ mask

	if newSize <= oldSize {
		return ptr
	}

	// Check if this is the most recent allocation. If it is, we can grow in-place.
	if a.next-oldSize == uintptr(ptr) {
		// Check if we have enough space available for the
		// requisite extra space.
		need := (newSize - oldSize) / maxAlign
		if a.left >= need {
			// Grow in-place.
			a.left -= need
			return ptr
		}
	}

	// Can't grow in place, allocate new memory and copy to it.
	new := a.Alloc(newSize, align)
	copy(
		unsafe.Slice((*byte)(new), newSize),
		unsafe.Slice((*byte)(ptr), oldSize),
	)

	return new
}

var pools [64]sync.Pool

func init() {
	for i := range pools {
		pools[i].New = func() any {
			return reflect.New(reflect.StructOf([]reflect.StructField{
				{
					Name: "X0",
					Type: reflect.ArrayOf(1<<i, reflect.TypeFor[uintptr]()),
				},
				{Name: "X1", Type: reflect.TypeFor[unsafe.Pointer]()},
			})).UnsafePointer()
		}
	}
}

func (a *Arena) allocChunk(words uintptr) unsafe.Pointer {
	log := bits.TrailingZeros(uint(words))
	chunk := pools[log].Get().(unsafe.Pointer)

	if len(a.chunks) > log {
		setChunkEnd(chunk, words, a.chunks[log])

		a.chunks[log] = chunk
	} else {
		setChunkEnd(chunk, words, unsafe.Pointer(a))

		// If this is the first chunk allocated, set a finalizer.
		if a.chunks == nil {
			runtime.SetFinalizer(a, (*Arena).finalize)
		}

		// Place the returned chunk at the offset in a.chunks that
		// corresponds to its log, so we can identify its size easily
		// in the loop above.
		a.chunks = append(a.chunks, make([]unsafe.Pointer, log+1-len(a.chunks))...)
		a.chunks[log] = chunk
	}

	return chunk
}

func (a *Arena) finalize() {
	for log, chunk := range a.chunks {
		if chunk == nil {
			continue
		}

		words := uintptr(1) << log

		setChunkEnd(chunk, words, nil) // Make sure that we don't leak the arena.

		pools[log].Put(chunk)
	}
}

func nextPow2(n uintptr) uintptr {
	return uintptr(1) << bits.Len(uint(n))
}

func setChunkEnd(chunk unsafe.Pointer, words uintptr, p unsafe.Pointer) {
	end := unsafe.Add(chunk, words*unsafe.Sizeof(uintptr(0)))
	*(*unsafe.Pointer)(end) = p
}
