//go:build go1.22

package arena

import (
	"math/bits"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/xunsafe"
	"github.com/flier/goutil/pkg/xunsafe/layout"
)

//go:generate ./make_shapes.sh shapes.go 49

func suggestSizeLog(bytes int) uint {
	// Snap to the next power of two.
	return max(4, uint(bits.Len(uint(bytes)-1)))
}

// SuggestSize suggests an allocation size by rounding up to a power of 2.
func SuggestSize(bytes int) int {
	// Snap to the next power of two.
	n := 1 << suggestSizeLog(bytes)
	if bytes == 0 {
		return n
	}
	return n
}

func (a *Arena) allocChunk(size int) (*byte, int) {
	log := suggestSizeLog(size)
	n := 1 << log
	if int(log) < len(a.blocks) {
		if a.blocks[log] == nil {
			a.blocks[log] = allocTraceable(n, unsafe.Pointer(a))
		}
		return a.blocks[log], n
	}

	p := allocTraceable(n, unsafe.Pointer(a))
	if a.blocks == nil {
		a.blocks = make([]*byte, 64)
		if debug.Enabled {
			addr := xunsafe.AddrOf(a)
			runtime.SetFinalizer(unsafe.SliceData(a.blocks), func(**byte) {
				debug.Log(nil, "arena collected", "addr: %v", addr)
			})
		}
	}
	a.blocks = a.blocks[:log+1]
	debug.Log(nil, "saving block", "a.blocks[%d] = %p -> %p", log, a.blocks[log], p)
	a.blocks[log] = p

	return p, n
}

// allocTraceable allocates size bytes of garbage-collected memory and returns
// a pointer to them.
//
// This function will also store ptr in the same allocation in such a way that
// as long as any pointer into the allocated memory is live, ptr will be marked
// as live by the garbage collector.
func allocTraceable(size int, ptr unsafe.Pointer) *byte {
	// This needs to be done with reflection, because we need a weirdly-shaped
	// allocation: a bunch of bytes followed by a pointer.
	//
	// To avoid the overhead of hammering reflection, we cache the shape for
	// each power of two size. For non-powers of two, we hammer reflection
	// every time, because that path is not used by the arena implementation.
	var shape reflect.Type
	size = layout.RoundUp(size, layout.Align[*byte]())

	if isPow2(size) {
		// Power-of-two shapes avoid needing to take a trip through reflection.
		// Calling reflect.New() on one of these types will immediately go to
		// runtime.mallocgc(), as if by new().
		shape = shapes[bits.TrailingZeros(uint(size))]
	} else {
		shape = reflect.StructOf([]reflect.StructField{
			{Name: "Data", Type: reflect.ArrayOf(size, reflect.TypeFor[byte]())},
			{Name: "Arena", Type: reflect.TypeFor[*Arena]()},
		})
	}

	p := (*byte)(reflect.New(shape).UnsafePointer())
	xunsafe.ByteStore(p, size, ptr) // Store the tracee pointer at the end.

	return p
}

func isPow2(n int) bool {
	return n&(n-1) == 0
}
