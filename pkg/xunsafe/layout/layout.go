//go:build go1.21

// Package layout includes helpers for working with type layouts.
//
// It is separate from xunsafe, because nothing in this package is actually
// unsafe.
package layout

import (
	"unsafe"

	"github.com/flier/goutil/internal/debug"
)

// Int is any integer type.
type Int interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

// Signed is an interface that represents any signed integer type in Go, including signed integers of various bit widths.
type Signed interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// Unsigned is an interface that represents any unsigned integer type in Go, including unsigned integers of various bit widths.
type Unsigned interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~uintptr
}

// Size returns T's size in bytes.
func Size[T any]() int {
	var z T

	return int(unsafe.Sizeof(z))
}

// Size returns T's size in bits.
func Bits[T any]() int {
	return Size[T]() * 8
}

// Size returns T's alignment in bytes.
func Align[T any]() int {
	var z T
	return int(unsafe.Alignof(z))
}

// Layout is the layout of some type.
type Layout struct {
	Size, Align int
}

// Of returns the size and alignment of a given type.
func Of[T any]() Layout {
	return Layout{Size[T](), Align[T]()}
}

// Max returns a layout whose size and alignment are both as large as the
// largest among l and that.
func (l Layout) Max(that Layout) Layout {
	return Layout{max(l.Size, that.Size), max(l.Align, that.Align)}
}

// RoundDown rounds v down to a power of two.
func RoundDown[T Int](v, align T) T {
	debug.Assert(v >= 0, "v must be greater than 0")
	debug.Assert(align > 0, "align must be greater than 0")

	if align <= 0 {
		return v
	}

	return v &^ (align - 1)
}

// RoundDown rounds v up to a power of two.
func RoundUp[T Int](v, align T) T {
	debug.Assert(v >= 0, "v must be greater than 0")
	debug.Assert(align > 0, "align must be greater than 0")

	if align <= 0 {
		return v
	}

	return (v + align - 1) &^ (align - 1)
}

// Padding returns [RoundUp](v, align) - v.
func Padding[T Int](v, align T) T {
	debug.Assert(v >= 0, "v must be greater than 0")
	debug.Assert(align > 0, "align must be greater than 0")

	if align <= 0 {
		return 0
	}

	return (align - v) & (align - 1)
}

// PadSlice appends zeros to buf until its length is a multiple of align.
func PadSlice(buf []byte, align int) []byte {
	debug.Assert(align > 0, "align must be greater than 0")

	return append(buf, make([]byte, Padding(len(buf), align))...)
}
