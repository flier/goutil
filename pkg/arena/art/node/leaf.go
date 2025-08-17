package node

import (
	"unsafe"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

type Leaf struct {
	Key   slice.Slice[byte]
	Value unsafe.Pointer
}

var _ Node = (*Leaf)(nil)

func (l *Leaf) Type() Type                   { return TypeLeaf }
func (l *Leaf) Full() bool                   { return true }
func (l *Leaf) Ref() Ref                     { return NewRef(TypeLeaf, l) }
func (l *Leaf) Prefix() *slice.Slice[byte]   { return &l.Key }
func (l *Leaf) Minimum() *Leaf               { return l }
func (l *Leaf) Maximum() *Leaf               { return l }
func (l *Leaf) AddChild(b byte, child AsRef) { panic("leaf cannot have children") }
func (l *Leaf) FindChild(b byte) *Ref        { panic("leaf cannot have children") }
func (l *Leaf) Grow(a *arena.Arena) Node     { panic("leaf cannot have children") }
