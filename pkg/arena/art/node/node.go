package node

import (
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

type Type int

const (
	TypeUnknown Type = iota
	TypeLeaf
	TypeNode4
	TypeNode16
	TypeNode48
	TypeNode256
)

type Node interface {
	Type() Type

	Full() bool

	Ref() Ref

	Prefix() *slice.Slice[byte]

	Minimum() *Leaf

	Maximum() *Leaf

	FindChild(b byte) *Ref

	AddChild(b byte, child AsRef)

	Grow(a *arena.Arena) Node
}

type Base struct {
	Partial     slice.Slice[byte]
	NumChildren int
}

func (n *Base) Prefix() *slice.Slice[byte] { return &n.Partial }
