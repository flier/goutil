//go:build go1.22

// Package arena provides a low-level, relatively unsafe arena allocation
// abstraction for high-performance memory management.
//
// Arena allocation is a memory management technique where memory is allocated
// from large, pre-allocated blocks rather than individual system allocations.
// This approach can significantly improve performance by reducing the overhead
// of frequent small allocations and improving cache locality.
//
// # Key Concepts
//
// Arena: A large block of pre-allocated memory from which smaller allocations
// are made. All memory in an arena is freed together when the arena is reset.
//
// Recycled Allocator: An enhanced arena allocator that maintains free lists
// for previously released memory blocks, enabling memory reuse and reducing
// fragmentation.
//
// Memory Safety: Arena-allocated memory must not be accessed after the arena
// is reset. The package provides mechanisms to ensure proper memory management.
//
// # Design
//
// See [Cheating the Reaper in Go] for detailed design information.
//
// Arenas are designed to only return pointers to data with pointer-free shape.
// However, we would like to store pointers in this data, so that the arena can
// point to itself (and to no other memory).
//
// This means that to store such data, the pointers must either live for the
// same lifetime as the [Arena] value (such as by storing them alongside it) or
// must point back into the Arena.
//
// We ensure this by making it so that holding a pointer onto any memory
// allocated by an [Arena] will keep all memory reachable from it alive.
// We achieve this by having the shape of each chunk allocated for the arena
// contain a pointer to the arena as a header; each chunk thus must have the
// shape:
//
//	type chunk struct {
//		memory [N]uint64
//		arena  *Arena
//	}
//
// By holding a pointer into chunk.memory anywhere reachable by a GC root (such
// as in a local variable) the GC will mark the allocation for the whole chunk
// as live, and therefore mark the [*Arena] field as live. Tracing through
// chunk.arena.chunks will mark all the other chunks as alive.
//
// Memory not directly allocated by an arena can be tied to it using
// [Arena.KeepAlive]. Using this operation is very slow, since this is the one
// part of the arena that is not re-used when calling [Arena.Reset].
//
// # Usage Patterns
//
// Basic Arena Usage
//
//	arena := &Arena{}
//
//	// Allocate memory
//	ptr := arena.Alloc(1024)
//
//	// Use the memory...
//
//	// Reset when done (frees all memory)
//	arena.Reset()
//
// Recycled Allocator Usage
//
//	recycled := &Recycled{}
//
//	// Allocate memory
//	ptr1 := recycled.Alloc(64)
//	ptr2 := recycled.Alloc(128)
//
//	// Release memory for reuse
//	recycled.Release(ptr1, 64)
//	recycled.Release(ptr2, 128)
//
//	// New allocations may reuse released memory
//	newPtr1 := recycled.Alloc(64)  // May reuse ptr1's memory
//	newPtr2 := recycled.Alloc(128) // May reuse ptr2's memory
//
//	// Reset when done with batch
//	recycled.Reset()
//
// Generic Allocation with New/Free
//
//	type MyStruct struct {
//		ID   int64
//		Name string
//	}
//
//	arena := &Recycled{}
//
//	// Allocate and initialize
//	ptr := New(arena, MyStruct{ID: 1, Name: "test"})
//
//	// Release automatically determines size
//	Free(arena, ptr)
//
// # Performance Characteristics
//
//   - Fast Allocation: O(1) allocation from pre-allocated blocks
//   - Memory Efficiency: Reduces fragmentation and improves cache locality
//   - Batch Operations: Efficient for scenarios with many small allocations
//   - Recycling Benefits: Recycled allocator provides additional performance
//     improvements through memory reuse
//
// # Memory Safety Considerations
//
//   - Lifetime Management: All pointers into arena memory become invalid
//     after calling [Arena.Reset]()
//   - No Individual Free: Individual memory blocks cannot be freed; only
//     the entire arena can be reset
//   - Pointer Restrictions: Arena-allocated data should not contain pointers
//     to memory outside the arena
//   - Reset Timing: Call [Arena.Reset]() only when all arena-allocated memory is
//     no longer needed
//
// # When to Use
//
// Arena allocation is most beneficial for:
//
//   - High-performance applications with many small allocations
//   - Batch processing where memory can be freed together
//   - Scenarios where memory fragmentation is a concern
//   - Performance-critical code that can tolerate the memory safety restrictions
//
// # Alternatives
//
// For simpler use cases, consider:
//
//   - Standard Go allocation for general-purpose memory management
//   - Object pools for frequently allocated/freed objects of the same type
//   - Custom allocators for domain-specific memory patterns
//
// [Cheating the Reaper in Go]: https://mcyoung.xyz/2025/04/21/go-arenas/
package arena

import (
	"unsafe"

	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/xunsafe"
	"github.com/flier/goutil/pkg/xunsafe/layout"
)

// Allocator is the interface that wraps the basic memory allocation and release
// operations. It provides a unified abstraction for different types of memory
// allocators, enabling polymorphic usage across the codebase.
//
// The Allocator interface is implemented by both [Arena] and [Recycled] types,
// allowing code to work with either allocator without modification. This design
// enables easy switching between different allocation strategies based on
// performance requirements and memory usage patterns.
//
// # Core Operations
//
//   - Alloc: Allocates a block of memory with the specified size
//   - Release: Returns previously allocated memory to the allocator
//
// # Implementation Differences
//
// Different allocator implementations provide varying levels of functionality:
//
//   - [Arena]: Basic arena allocation with no memory recycling
//   - [Recycled]: Enhanced arena allocation with automatic memory recycling
//
// # Usage Patterns
//
// The Allocator interface enables several common usage patterns:
//
//   - Generic Functions: Functions like [New] and [Free] work with any Allocator implementation
//   - Polymorphic Allocation: Code can accept any Allocator and work with it transparently
//   - Strategy Selection: Applications can choose the appropriate allocator based on their needs
//
// # Example Usage
//
//	func ProcessData(a arena.Allocator) {
//		// Allocate memory using the interface
//		data := a.Alloc(1024)
//
//		// Use the memory...
//
//		// Release when done
//		a.Release(data, 1024)
//	}
//
//	// Can be called with different allocator types
//	arena := &arena.Arena{}
//	recycled := &arena.Recycled{}
//
//	ProcessData(arena)     // Uses basic arena allocation
//	ProcessData(recycled)  // Uses recycled allocation with memory reuse
//
// # Memory Safety
//
// All implementations of Allocator must ensure:
//
//   - Allocated memory is valid until explicitly released
//   - Released memory is not accessed by the application
//   - Memory alignment requirements are satisfied
//   - Proper cleanup on allocator reset or destruction
//
// # Performance Characteristics
//
// The choice of Allocator implementation affects performance:
//
//   - [Arena]: Fast allocation, no recycling overhead, higher memory usage
//   - [Recycled]: Slightly slower allocation, significant memory reuse,
//     lower overall memory consumption
//
// # When to Use Each Implementation
//
// Use [Arena] when:
//
//   - Memory recycling is not needed
//   - Maximum allocation performance is required
//   - Simple memory management is preferred
//   - Memory usage patterns are predictable
//
// Use [Recycled] when:
//
//   - Frequent allocation/deallocation cycles occur
//   - Memory fragmentation is a concern
//   - Long-running applications need memory efficiency
//   - Complex allocation patterns exist
//
// # Thread Safety
//
// The Allocator interface does not guarantee thread safety. If multiple
// goroutines access the same allocator concurrently, external synchronization
// must be provided by the caller.
//
// # Error Handling
//
// Allocator methods do not return errors. Allocation failures (e.g., out of
// memory) typically result in panics. Applications should ensure adequate
// memory is available or handle panics appropriately.
type Allocator interface {
	// Alloc allocates size bytes of memory and returns a pointer to the
	// allocated block.
	//
	// The returned pointer is guaranteed to be valid until explicitly
	// released or until the allocator is reset. The memory is aligned
	// according to the allocator's alignment requirements.
	//
	// Args:
	//   size: The number of bytes to allocate. Must be non-negative.
	//
	// Returns:
	//   A pointer to the allocated memory block. The memory contents
	//   are undefined and should be initialized before use.
	//
	// Panics:
	//   May panic if size is negative or if allocation fails due to
	//   insufficient memory or other system constraints.
	//
	// Example:
	//   ptr := allocator.Alloc(64)
	//   // ptr points to 64 bytes of uninitialized memory
	//   // Use ptr, then release when done
	Alloc(size int) *byte

	// Release returns previously allocated memory back to the allocator
	// for potential reuse or cleanup.
	//
	// After calling Release, the memory block should not be accessed
	// by the application. The behavior of Release varies by implementation:
	//
	//   - [Arena.Release]: No-op (memory is freed on Reset)
	//   - [Recycled.Release]: Memory is added to free lists for reuse
	//
	// Args:
	//   p: Pointer to the memory block to release. Must be a valid
	//      pointer previously returned by Alloc.
	//   size: The size of the memory block in bytes. Must match the
	//         size used when allocating the memory.
	//
	// Behavior:
	//   - The memory block becomes invalid and should not be accessed
	//   - For Recycled allocators, the memory may be reused in future
	//     allocations of the same size
	//   - For Arena allocators, the memory is freed when Reset is called
	//
	// Example:
	//   ptr := allocator.Alloc(64)
	//   // ... use ptr ...
	//   allocator.Release(ptr, 64)  // Release the memory
	//   // ptr is now invalid
	Release(p *byte, size int)
}

type AllocatorExt interface {
	Allocator

	// Next returns the next available address in the arena.
	Next() xunsafe.Addr[byte]

	// End returns the end of the arena.
	End() xunsafe.Addr[byte]

	// Cap returns the current capacity of the arena.
	Cap() int

	// Advance advances the next available address in the arena by n bytes.
	Advance(n int)

	// Log logs a message to the arena.
	Log(op, format string, args ...any)
}

// Arena is an Arena for holding values of any type which does not contain
// pointers.
//
// A zero Arena is empty and ready to use.
type Arena struct {
	_ xunsafe.NoCopy

	// Exported to allow for open-coding of Alloc() in some hot callsites,
	// because Go won't inline it >_>
	next, end xunsafe.Addr[byte]
	cap       int // Always a power of 2.

	// Blocks of memory allocated by this arena. Indexed by their size log 2.
	blocks []*byte

	// Data to keep around for the GC to mark whenever it marks an arena.
	// Holding any pointer to the arena will keep anything here alive, too.
	keep []unsafe.Pointer
}

var _ Allocator = (*Arena)(nil)

// Align is the alignment of all objects on the arena.
const Align = int(unsafe.Sizeof(uintptr(0)))

// New allocates a new value of type T on an arena.
func New[T any](a Allocator, value T) *T {
	layout := layout.Of[T]()
	if layout.Align > Align {
		panic("over-aligned object")
	}

	p := xunsafe.Cast[T](a.Alloc(layout.Size))
	*p = value
	return p
}

// Free releases a value of type T previously allocated from the given allocator
// back to its free list for recycling.
//
// This is a convenience function that automatically determines the size of type T
// using layout metadata and calls the allocator's Release method. It's particularly
// useful with Recycled allocators, where it enables automatic memory recycling
// without manual size tracking.
//
// # Type Safety
//
// The function automatically derives the correct size from the type T using
// [layout.Of][T]().Size, ensuring that the exact allocated size is used when
// releasing memory. This prevents size mismatches that could lead to memory
// corruption or inefficient recycling.
//
// # Usage with Recycled Allocators
//
// When used with a [Recycled] allocator, Free automatically:
//   - Determines the appropriate size class for type T
//   - Adds the memory block to the correct free list
//   - Enables future allocations of the same size to reuse this memory
//
// # Usage with Regular Arena Allocators
//
// When used with a regular [Arena] allocator, [Free] calls [Arena.Release] which is a no-op,
// effectively doing nothing. This allows the same code to work with both allocator
// types without modification.
//
// # Example
//
//	type MyStruct struct {
//		ID   int64
//		Name string
//	}
//
//	arena := &Recycled{}
//
//	// Allocate using New
//	ptr := New(arena, MyStruct{ID: 1, Name: "test"})
//
//	// Release using Free - automatically determines size
//	Free(arena, ptr)
//
//	// Future allocations of MyStruct may reuse this memory
//	newPtr := New(arena, MyStruct{ID: 2, Name: "reused"})
//
// # Memory Safety
//
//   - The pointer p must point to valid memory previously allocated by the
//     given allocator
//   - The memory should not be accessed after calling Free
//   - For Recycled allocators, the memory will be automatically cleared
//     when reused to prevent data leakage
//
// # Performance
//
// Free is a lightweight wrapper that:
//   - Performs a single layout lookup to determine size
//   - Calls the allocator's Release method
//   - Has minimal overhead compared to manual size tracking
//
// For performance-critical code that allocates many objects of the same type,
// consider manually tracking sizes to avoid repeated layout lookups.
func Free[T any](a Allocator, p *T) {
	size := layout.Of[T]().Size

	a.Release(xunsafe.Cast[byte](p), size)
}

// KeepAlive ensures that v is not swept by the GC until all pointers into the
// arena go away.
func (a *Arena) KeepAlive(v any) {
	a.keep = append(a.keep, unsafe.Pointer(xunsafe.AnyData(v)))
}

// Alloc allocates memory with the given size.
//
// All memory is pointer-aligned. The memory may be uninitialized.
//
// Do not use this method directly, use [New] instead.
func (a *Arena) Alloc(size int) *byte {
	// Align size to a pointer boundary.
	alignedSize := alignUp(size)

	if a.next.Add(alignedSize) <= a.end {
		// Duplicating this block ensures that Go schedules this branch
		// correctly. This block is the "hot" side of the branch.
		p := a.next.AssertValid()
		a.next = a.next.Add(alignedSize)
		a.Log("alloc", "%v:%v, %d:%d", p, a.next, alignedSize, Align)
		return p
	}

	a.Grow(alignedSize)
	p := a.next.AssertValid()
	a.next = a.next.Add(alignedSize)
	a.Log("alloc", "%v:%v, %d:%d", p, a.next, alignedSize, Align)
	return p
}

// Release is a no-op for Arena.
//
// Do not use this method directly, use [Free] instead.
func (a *Arena) Release(p *byte, size int) {}

// Reserve ensures that at least size bytes can be allocated without calling
// [Arena.Grow].
func (a *Arena) Reserve(size int) {
	if a.next.Add(size) > a.end {
		a.Grow(size)
	}
}

// Reset resets this arena to an "empty" state, allowing all memory allocated by
// it to be re-used.
//
// Although this can be used to amortize trips into Go's allocator, doing so
// trades off safety: any memory allocated by the arena must not be referenced
// after a call to Reset.
func (a *Arena) Reset() {
	if len(a.blocks) == 0 {
		return
	}

	// Discard all but the largest block, which we clear. This means that as
	// an arena is re-used, we will eventually wind up learning the size of the
	// largest block we need to allocate, and use only that one, meaning that
	// "average" calls should never have to call Grow().
	end := len(a.blocks) - 1
	clear(a.blocks[:end])
	xunsafe.Clear(a.blocks[end], 1<<end)

	// Set up next/end/cap to point to the largest block.
	a.next = xunsafe.AddrOf(a.blocks[end])
	a.end = a.next.Add(1 << end)
	a.cap = 1 << end

	// Order doesn't matter here: nothing in a.blocks can point into a.keep,
	// because the only GC-visible pointers in a.blocks are pointers back to
	// a, the arena header.
	//
	// We set this to nil because clearing this will walk us right into an
	// unavoidable bulk write barrier. By writing nil, we only pay for a fast
	// single-pointer write barrier, and make cleaning up the handful of bytes
	// this throws out the GC's problem.
	//
	// In profiling, it turns out that doing clear(a.keep) is several times
	// more expensive than the noscan clear that happens below.
	a.keep = nil
}

// Grow allocates fresh memory onto next of at least the given size.
//
//go:nosplit
func (a *Arena) Grow(size int) {
	xunsafe.Escape(a)
	p, n := a.allocChunk(max(size, a.cap*2))
	// No need to KeepAlive(p) this pointer, since allocChunk sticks it in the
	// dedicated memory block array.

	a.next = xunsafe.AddrOf(p)
	a.end = a.next.Add(n)
	a.cap = n
	a.Log("grow", "%v:%v:%d\n", a.next, a.end, a.cap)
}

func (a *Arena) Next() xunsafe.Addr[byte] { return a.next }
func (a *Arena) End() xunsafe.Addr[byte]  { return a.end }
func (a *Arena) Cap() int                 { return a.cap }
func (a *Arena) Advance(n int)            { a.next.Add(n) }

func (a *Arena) Log(op, format string, args ...any) {
	debug.Log([]any{"%p %v:%v", a, a.next, a.end}, op, format, args...)
}
