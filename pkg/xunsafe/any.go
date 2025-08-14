//go:build go1.23

package xunsafe

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/flier/goutil/internal/xsync"
)

var isDirectMap xsync.Map[reflect.Type, bool]

// iface is the internal representation an a Go interface value.
type iface struct {
	itab uintptr
	data *byte
}

// AnyData extracts the pointer value from an any.
func AnyData(v any) *byte {
	return Cast[iface](NoEscape(&v)).data
}

// AnyType extracts the opaque type from an any.
func AnyType(v any) uintptr {
	return Cast[iface](NoEscape(&v)).itab
}

// AnyBytes extracts a slice pointing to the variable-length data of an any.
func AnyBytes(v any) []byte {
	if v == nil {
		return nil
	}

	p := AnyData(v)
	if !IsDirectAny(v) {
		return unsafe.Slice(p, reflect.TypeOf(v).Size())
	}

	p2 := p // Work around https://github.com/golang/go/issues/74364
	return Bytes(&p2)
}

// MakeAny builds an any out of the given data.
func MakeAny(typ uintptr, data *byte) any {
	raw := iface{typ, data}
	return BitCast[any](raw)
}

// IsDirectAny returns whether or not this any is a direct interface.
//
// This is much slower than [IsDirect], because we can't use the same trick
// we use in [IsDirect] without making a heap allocation in some cases.
func IsDirectAny(v any) bool {
	t := reflect.TypeOf(v)
again:
	switch t.Kind() {
	case reflect.Pointer, reflect.UnsafePointer, reflect.Func,
		reflect.Map, reflect.Chan:
		return true

	case reflect.Array:
		if t.Len() != 1 {
			return false
		}
		t = t.Elem()
		goto again

	case reflect.Struct:
		if t.NumField() == 1 {
			t = t.Field(0).Type
			goto again
		}

		direct, _ := isDirectMap.LoadOrStore(t, func() bool {
			z := reflect.Zero(t).Interface()
			p := AnyData(z)
			return p == nil // See InlinedAny below.
		})
		return direct

	default:
		return false
	}
}

// IsDirect returns whether converting T into an interface requires calling
// the allocator.
//
// This is true for any type which is one of the inlined primitives
// (pointers, interfaces, channels, maps) or one-field structs/arrays whose
// sole element is inlined. Note that this does *not* include all pointer-shaped
// structs, such as struct{struct{}; int}.
func IsDirect[T any]() bool {
	var x T
	p := AnyData(any(x))

	// x's bit pattern is all zero, no matter what T is. If T is indirect,
	// the pointer extracted by AnyData will be nil. Otherwise, it will be a
	// stack pointer of some kind.
	return p == nil
}

// AssertInlinedAny is a helper for testing that T does not allocate when
// converted to an interface.
func AssertInlinedAny[T any](t testing.TB) {
	t.Helper()
	var z T
	assert.True(t, IsDirect[T](), "expected %T to be pointer-shaped", z)
}
