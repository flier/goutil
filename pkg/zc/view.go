// Package zc provides helpers for working with zero-copy ranges.
package zc

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/xunsafe"
)

// View is a representation of a []byte as a slice relative to some larger byte
// array, such as the source of a parsed message.
//
// This is a packed representation of a value with the layout
//
//	struct {
//	  offset, length uint32
//	}
//
// The zero value faithfully represents an empty slice.
type View uint64

// New creates a new View over the given source buffer with the given start
// and length.
func New(src *byte, start *byte, len int) View {
	offset := xunsafe.Sub(start, src)
	return Raw(offset, len)
}

// Raw is like New, but it only takes the offset and length.
func Raw(offset, len int) View {
	debug.Assert(offset <= math.MaxUint32 && len <= math.MaxUint32,
		"offset too large for zc: [%d:%d]", offset, len)
	return View(offset) | View(len)<<32
}

// Start returns the start offset of this slice within its source.
func (r View) Start() int { return int(uint32(r)) }

// End returns the end offset of this slice within its source.
func (r View) End() int { return r.Start() + r.Len() }

// Len returns the length of this View.
func (r View) Len() int { return int(r >> 32) }

// Bytes converts this View into a byte slice, given its source.
//
// NOTE: Go refuses to inline this function sometimes. View.String does not
// appear to have this problem.
func (r View) Bytes(src *byte) []byte {
	if r.Len() == 0 {
		return nil
	}
	return unsafe.Slice(xunsafe.Add(src, r.Start()), r.Len())
}

// Format implements [fmt.Formatter].
func (r View) Format(s fmt.State, verb rune) {
	debug.Fprintf("[%d:%d]", r.Start(), r.End()).Format(s, verb)
}

// ExtractFrom is a helper for creating extraction funcs. It exists to work
// around an inliner limitation.
type ExtractFrom struct {
	Src *byte
}

// ExtractBytes returns a func that calls [View.Bytes].
//
// This exists to work around inlining failure in certain places in the parser.
func (e ExtractFrom) Bytes(raw uint64) []byte {
	r := View(raw)
	p := (*byte)(unsafe.Add(unsafe.Pointer(e.Src), r.Start()))
	return unsafe.Slice(p, r.Len())
}
