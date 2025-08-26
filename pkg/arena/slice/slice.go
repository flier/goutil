//go:build go1.20

package slice

import (
	"fmt"
	"unsafe"

	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/opt"
	"github.com/flier/goutil/pkg/xunsafe"
	"github.com/flier/goutil/pkg/xunsafe/layout"
)

// Slice is a slice that points into an arena.
//
// Unlike an ordinary slice, it does not contain pointers; in order to work
// correctly, it must be kept alive no longer than its owning arena.
type Slice[T any] struct {
	ptr      *T
	len, cap uint32
}

// Static assert that the size of Slice[T] is 16 bytes.
var _ [16]byte = [unsafe.Sizeof(Slice[byte]{})]byte{}

// FromBytes allocates a slice for the given bytes.
func FromBytes(a arena.Allocator, b []byte) Slice[byte] {
	return Of(a, b...)
}

// FromString allocates a slice for the given string.
func FromString(a arena.Allocator, s string) Slice[byte] {
	return Of(a, []byte(s)...)
}

// FromParts assembles a slice from its raw components.
func FromParts[T any](ptr *T, len, cap uint32) Slice[T] {
	return Slice[T]{ptr, len, cap}
}

// Wrap creates a Slice[T] from an existing Go slice without copying or allocating memory.
//
// This function wraps an existing Go slice into the arena slice type, allowing it to be used
// with arena-related functions. The returned slice shares the same underlying memory as the
// input slice, so any modifications to the input slice will be reflected in the returned slice.
//
// Parameters:
//   - s: The Go slice to wrap.
//
// Returns:
//   - A Slice[T] that wraps the input slice without copying data.
//
// Example:
//
//	goSlice := []int{1, 2, 3, 4, 5}
//	arenaSlice := slice.Wrap(goSlice)
//	// arenaSlice now wraps goSlice without copying
func Wrap[T any](s []T) Slice[T] {
	if len(s) == 0 {
		return Slice[T]{}
	}

	return Slice[T]{xunsafe.Cast[T](unsafe.SliceData(s)), uint32(len(s)), uint32(cap(s))}
}

// Of allocates a slice for the given values.
func Of[T any](a arena.Allocator, values ...T) Slice[T] {
	s := Make[T](a, len(values))
	copy(s.Raw(), values)
	return s
}

// Clone clones a slice.
func Clone[T any](a arena.Allocator, s Slice[T]) Slice[T] {
	return Of(a, s.Raw()...)
}

// Make allocates a slice of the given length.
func Make[T any](a arena.Allocator, n int) Slice[T] {
	cap := sliceLayout[T](n)
	p := xunsafe.Cast[T](a.Alloc(cap))

	size := layout.Size[T]()
	s := FromParts(p, uint32(n), uint32(cap/size))
	return s
}

// Release releases the slice.
func (s Slice[T]) Release(a arena.Allocator) {
	a.Release(xunsafe.Cast[byte](s.ptr), s.Cap()*layout.Size[T]())
}

// Equal returns true if a and b are equal.
//
//go:nosplit
func Equal[T comparable](a, b Slice[T]) bool {
	if a.Ptr() == nil && b.Ptr() == nil {
		return true
	}

	if a.Ptr() == nil || b.Ptr() == nil {
		return false
	}

	if a.Len() != b.Len() {
		return false
	}

	if a.Ptr() == b.Ptr() {
		return true
	}

	for i := 0; i < a.Len(); i++ {
		if a.unsafeLoad(i) != b.unsafeLoad(i) {
			return false
		}
	}

	return true
}

// EqualTo returns true if a and b are equal.
//
//go:nosplit
func EqualTo[T comparable](a Slice[T], b []T) bool {
	if a.Len() != len(b) {
		return false
	}

	for i := 0; i < a.Len(); i++ {
		if a.unsafeLoad(i) != b[i] {
			return false
		}
	}

	return true
}

// HasPrefix checks if a has the given prefix.
//
//go:nosplit
func HasPrefix[T comparable](a Slice[T], b []T) bool {
	if a.Len() < len(b) {
		return false
	}

	for i := 0; i < len(b); i++ {
		if a.unsafeLoad(i) != b[i] {
			return false
		}
	}

	return true
}

// Addr converts this slice into an address slice.
//
// See the caveats of [xunsafe.AddrOf].
func (s Slice[T]) Addr() Addr[T] {
	return Addr[T]{xunsafe.AddrOf(s.ptr), s.len, s.cap}
}

// Ptr returns this slice's pointer value.
func (s Slice[T]) Ptr() *T { return xunsafe.Cast[T](s.ptr) }

// Empty returns true if this slice is empty.
func (s Slice[_]) Empty() bool { return s.len == 0 }

// Len returns this slice's length.
func (s Slice[_]) Len() int { return int(s.len) }

// SetLen directly sets the length of s.
func (s Slice[T]) SetLen(n int) Slice[T] {
	if debug.Enabled && n > int(s.cap) {
		panic(fmt.Errorf("runtime error: SetLen(%v) with Cap() = %v", n, s.cap))
	}

	debug.Log(nil, "set len", "%v->%d", s.Addr(), n)
	s.len = uint32(n)
	return s
}

// Cap returns this slice's capacity.
func (s Slice[_]) Cap() int { return int(s.cap) }

// Get returns the pointer to the given index.
func (s Slice[T]) Get(n int) *T {
	if debug.Enabled {
		return &s.Raw()[n]
	}

	return s.unsafeGet(n)
}

// CheckedGet returns the pointer to the given index, returning None if the index is out of bounds.
func (s Slice[T]) CheckedGet(n int) opt.Option[*T] {
	if n < 0 || n >= s.Len() {
		return opt.None[*T]()
	}

	return opt.Some(s.unsafeGet(n))
}

func (s Slice[T]) unsafeGet(n int) *T { return xunsafe.Add(s.Ptr(), n) }

// Load loads a value at the given index.
func (s Slice[T]) Load(n int) T {
	if debug.Enabled {
		return s.Raw()[n]
	}

	return s.unsafeLoad(n)
}

// CheckedLoad loads a value at the given index, returning None if the index is out of bounds.
func (s Slice[T]) CheckedLoad(n int) opt.Option[T] {
	if n < 0 || n >= s.Len() {
		return opt.None[T]()
	}

	return opt.Some(s.unsafeLoad(n))
}

// unsafeLoad loads a value at the given index.
//
// This function is used to avoid the overhead of the debug.Enabled check.
//
// It is only used in the unsafe code path, the caller must ensure that the index is in bounds.
//
//go:nosplit
func (s Slice[T]) unsafeLoad(n int) T {
	return xunsafe.Load(s.Ptr(), n)
}

// Store stores a value at the given index.
func (s Slice[T]) Store(n int, v T) {
	if debug.Enabled {
		s.Raw()[n] = v
	}

	xunsafe.Store(s.Ptr(), n, v)
}

// Raw returns the underlying slice for this slice.
//
// The return value of this function must never escape outside of this module.
func (s Slice[T]) Raw() []T {
	if s.ptr == nil || s.len == 0 {
		return nil
	}

	return unsafe.Slice(s.Ptr(), s.cap)[:s.len]
}

// Rest returns the portion of s between the length and the capacity.
//
// The return value of this function must never escape outside of this module.
func (s Slice[T]) Rest() []T {
	return unsafe.Slice(xunsafe.Add(s.Ptr(), s.len), s.cap-s.len)
}

// Slice returns a slice of s between the given start and end indices.
//
// Parameters:
//
//	start:
//		Zero-based index at which to start extraction.
//
//		Negative index counts back from the end of the slice — if -slice.len <= start < 0, start + slice.len is used.
//		If start < -slice.len or start is omitted, 0 is used.
//		If start >= slice.len, an empty slice is returned.
//
//	end:
//		Zero-based index at which to end extraction. slice() extracts up to but not including end.
//
//		Negative index counts back from the end of the slice — if -slice.len <= end < 0, end + slice.len is used.
//		If end < -slice.len, 0 is used.
//		If end >= slice.len or end is omitted or undefined, slice.len is used, causing all elements until the end to be extracted.
//		If end implies a position before or at the position that start implies, an empty slice is returned.
func (s Slice[T]) Slice(start, end int) Slice[T] {
	// Early return for empty slice
	if s.len == 0 {
		return Slice[T]{}
	}

	if start < 0 {
		if start >= -int(s.len) {
			start += int(s.len)
		} else {
			start = 0
		}
	} else if start >= int(s.len) {
		return Slice[T]{}
	}

	if end < 0 {
		if end >= -int(s.len) {
			end += int(s.len)
		} else {
			end = 0
		}
	} else if end >= int(s.len) {
		end = int(s.len)
	}

	// Return empty slice if indices are invalid
	if start >= end {
		return Slice[T]{}
	}

	// Calculate new capacity more accurately
	cap := s.cap - uint32(start)
	// Ensure capacity doesn't go below the new length
	if cap < uint32(end-start) {
		cap = uint32(end - start)
	}

	return Slice[T]{
		ptr: xunsafe.Add(s.ptr, start),
		len: uint32(end - start),
		cap: cap,
	}
}

// SplitAt splits a slice at the given index, returning two new slices that share
// the same underlying memory as the original slice.
//
// The method creates two slices:
//   - Left slice (l): contains elements from index 0 to n (exclusive)
//   - Right slice (r): contains elements from index n to the end
//
// Parameters:
//   - n: The index at which to split the slice. Can be negative or beyond the slice length.
//     Negative indices are interpreted as counting from the end (-1 = last element).
//     Indices beyond the slice length are clamped to the slice length.
//
// Returns:
//   - l: Left slice containing elements [0:n). Empty if n <= 0.
//   - r: Right slice containing elements [n:len). Empty if n >= len.
//
// Behavior:
//   - If n < 0: n is converted to a positive index by adding the slice length.
//     If the result is still negative, n is clamped to 0.
//   - If n >= len: n is clamped to len, making the right slice empty.
//   - Both returned slices have capacity equal to their length for memory efficiency.
//   - The slices share the same underlying memory, so modifications to the original
//     slice will affect both parts.
//   - If the original slice is empty, both returned slices are empty.
//
// Examples:
//
//	s := slice.Of(a, 1, 2, 3, 4, 5)
//	left, right := s.SplitAt(2)
//	// left: [1, 2], right: [3, 4, 5]
//
//	left, right := s.SplitAt(-2)
//	// left: [1, 2, 3], right: [4, 5]
//
//	left, right := s.SplitAt(0)
//	// left: [], right: [1, 2, 3, 4, 5]
//
//	left, right := s.SplitAt(5)
//	// left: [1, 2, 3, 4, 5], right: []
//
// Memory characteristics:
//   - O(1) time complexity - no data copying
//   - Both slices reference the same arena memory
//   - Capacity is set to length for both slices to prevent accidental growth
func (s Slice[T]) SplitAt(n int) (l Slice[T], r Slice[T]) {
	if s.len == 0 {
		return
	}

	if n < 0 {
		if n >= -int(s.len) {
			n += int(s.len)
		} else {
			n = 0
		}
	} else if n >= int(s.len) {
		n = int(s.len)
	}

	// Left slice: from start to n, capacity equals length
	l = Slice[T]{s.ptr, uint32(n), uint32(n)}

	// Right slice: from n to end, capacity equals length
	r = Slice[T]{xunsafe.Add(s.ptr, n), s.len - uint32(n), s.cap - uint32(n)}

	return
}

// Clone clones a slice.
func (s Slice[T]) Clone(a arena.Allocator) Slice[T] {
	return Clone(a, s)
}

// Prepend prepends the given elements to a slice, reallocating on the given
// arena if necessary.
func (s Slice[T]) Prepend(a arena.AllocatorExt, elems ...T) Slice[T] {
	var z T
	a.Log("prepend", "%p[%d:%d], %T x %d", s.ptr, s.len, s.cap, z, len(elems))

	if s.Cap()-s.Len() < len(elems) {
		s = s.Grow(a, len(elems))
	}

	buf := unsafe.Slice(s.Ptr(), s.cap)

	copy(buf[len(elems):], buf[:s.len])
	copy(buf[:len(elems)], elems)

	s.len += uint32(len(elems))

	return s
}

// Append appends the given elements to a slice, reallocating on the given
// arena if necessary.
func (s Slice[T]) Append(a arena.AllocatorExt, elems ...T) Slice[T] {
	var z T
	a.Log("append", "%p[%d:%d], %T x %d", s.ptr, s.len, s.cap, z, len(elems))

	if s.Cap()-s.Len() < len(elems) {
		s = s.Grow(a, len(elems))
	}

	copy(s.Rest(), elems)
	s.len += uint32(len(elems))

	return s
}

// AppendOne is an optimized version of append for one element.
//
//go:nosplit
func (s Slice[T]) AppendOne(a arena.AllocatorExt, elem T) Slice[T] {
	a.Log("append", "%p[%d:%d], %T x 1", s.ptr, s.len, s.cap, elem)

	if s.Len() == s.Cap() {
		s = s.Grow(a, 1)
	}

	xunsafe.Store(s.Ptr(), s.len, elem)
	s.len += 1
	return s
}

// Grow extends the capacity of this slice by n bytes.
func (s Slice[T]) Grow(a arena.AllocatorExt, n int) Slice[T] {
	var z T
	size := layout.Size[T]()
	a.Log("grow", "%p[%d:%d], %d x %T", s.ptr, s.len, s.cap, n, z)

	if s.ptr == nil {
		cap := sliceLayout[T](n)
		s.ptr = xunsafe.Cast[T](a.Alloc(cap))
		s.cap = uint32(cap) / uint32(size)
		return s
	}

	oldSize := sliceLayout[T](s.Cap())
	newSize := sliceLayout[T](s.Cap() + n)

	// Originally, this was arena.Realloc. It is inlined in-place for speed.
	p := xunsafe.Cast[byte](s.ptr)
	for i := 0; i < 1; i++ {
		// This Just Works regardless of whether the allocation is growing or
		// shrinking. If it's shrinking, delta will be negative, and a.left
		// is never negative, so this will add back the spare capacity.
		i := a.Next().Add(-oldSize)
		j := i.Add(newSize)
		if xunsafe.AddrOf(p) == i && j <= a.End() {
			a.Advance(newSize)
			a.Log("fast realloc", "%p, %d->%d:%d", p, oldSize, newSize, arena.Align)
			break
		}

		if newSize < oldSize {
			a.Log("realloc", "%p, %d->%d:%d", p, oldSize, newSize, arena.Align)
			break
		}

		q := a.Alloc(newSize)
		a.Log("realloc", "%p->%p, %d->%d:%d", p, q, oldSize, newSize, arena.Align)
		if oldSize > 0 {
			xunsafe.Copy(q, p, oldSize)
		}

		p = q
	}

	s.ptr = xunsafe.Cast[T](p)
	s.cap = uint32(newSize) / uint32(size)
	return s
}

// Format implements [fmt.Formatter].
func (s Slice[T]) Format(state fmt.State, v rune) {
	if s.Ptr() == nil && (s.Len() != 0 || s.Cap() != 0) {
		_, _ = fmt.Fprintf(state, "%v", s.Addr())
		return
	}

	_, _ = fmt.Fprintf(state, fmt.FormatString(state, v), s.Raw())
}

// Addr is like [Slice], but its pointer is replaced with an address, so
// loading/storing values of this type issues no write barriers.
type Addr[T any] struct {
	Ptr      xunsafe.Addr[T]
	Len, Cap uint32
}

// AssertValid converts this address slice into a true [Slice].
//
// See the caveats of [xunsafe.Addr.AssertValid].
func (s Addr[T]) AssertValid() Slice[T] {
	return Slice[T]{s.Ptr.ClearSignBit().AssertValid(), s.Len, s.Cap}
}

// Untyped converts this address slice into a true [Slice].
//
// See the caveats of [xunsafe.Addr.AssertValid].
func (s Addr[T]) Untyped() Untyped {
	return Untyped{
		Ptr: xunsafe.Addr[byte](s.Ptr),
		Len: s.Len,
		Cap: s.Cap,
	}
}

// String implements [fmt.Stringer].
func (s Addr[T]) String() string {
	return s.Untyped().String()
}

// Untyped is an [Addr] that has forgotten what type it is.
type Untyped struct {
	Ptr      xunsafe.Addr[byte]
	Len, Cap uint32
}

// OffArena creates a new off-arena slice.
//
// When cast to a concrete type, this will clear.
func OffArena[T any](ptr *T, len int) Untyped {
	return Untyped{
		Ptr: ^xunsafe.Addr[byte](xunsafe.AddrOf(ptr)),
		Len: uint32(len),
		Cap: uint32(len),
	}
}

// CastUntyped gives a type to a [Untyped], asserting it as valid in the
// process.
func CastUntyped[To any](s Untyped) Slice[To] {
	return Slice[To]{
		ptr: xunsafe.Addr[To](s.Ptr.ClearSignBit()).AssertValid(),
		len: s.Len,
		cap: s.Cap,
	}
}

// OffArena returns if this is off-arena memory, i.e., as created with
// [OffArena].
func (s Untyped) OffArena() bool {
	return s.Ptr.SignBit()
}

// String implements [fmt.Stringer].
func (s Untyped) String() string {
	return fmt.Sprintf("%v[%d:%d]", s.Ptr, s.Len, s.Cap)
}

func sliceLayout[T any](n int) (size int) {
	layout := layout.Of[T]()

	if debug.Enabled && layout.Align > arena.Align {
		// This doesn't seem to inline correctly if we don't use debug.Enabled.
		panic("over-aligned object")
	}

	return arena.SuggestSize(layout.Size * n)
}
