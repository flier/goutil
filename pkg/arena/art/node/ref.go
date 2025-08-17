package node

import (
	"unsafe"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/xunsafe"
)

type AsRef interface {
	Ref() Ref
}

type Ref uintptr

func NewRef[T any](t Type, p *T) Ref {
	addr := xunsafe.AddrOf(p)

	return Ref((uintptr(addr) & nodePtrMask) | (uintptr(t) & nodeTypeMask))
}

const (
	nodePtrMask  = ^nodeTypeMask
	nodeTypeMask = uintptr(arena.Align - 1)
)

func (r Ref) Ref() Ref        { return r }
func (r Ref) Type() Type      { return Type(uintptr(r) & nodeTypeMask) }
func (r Ref) Empty() bool     { return r == 0 }
func (r Ref) IsLeaf() bool    { return r.Type() == TypeLeaf }
func (r Ref) IsNode4() bool   { return r.Type() == TypeNode4 }
func (r Ref) IsNode16() bool  { return r.Type() == TypeNode16 }
func (r Ref) IsNode48() bool  { return r.Type() == TypeNode48 }
func (r Ref) IsNode256() bool { return r.Type() == TypeNode256 }
func (r Ref) IsNode() bool    { return r.IsNode4() || r.IsNode16() || r.IsNode48() || r.IsNode256() }

func (r Ref) AsLeaf() *Leaf {
	if r.IsLeaf() {
		return (*Leaf)(r.ptr())
	}

	return nil
}

func (r Ref) AsNode4() *Node4 {
	if r.IsNode4() {
		return (*Node4)(r.ptr())
	}

	return nil
}

func (r Ref) AsNode16() *Node16 {
	if r.IsNode16() {
		return (*Node16)(r.ptr())
	}

	return nil
}

func (r Ref) AsNode48() *Node48 {
	if r.IsNode48() {
		return (*Node48)(r.ptr())
	}

	return nil
}

func (r Ref) AsNode256() *Node256 {
	if r.IsNode256() {
		return (*Node256)(r.ptr())
	}

	return nil
}

func (r Ref) AsNode() Node {
	if r == 0 {
		return nil
	}

	switch r.Type() {
	case TypeLeaf:
		return (*Leaf)(r.ptr())
	case TypeNode4:
		return (*Node4)(r.ptr())
	case TypeNode16:
		return (*Node16)(r.ptr())
	case TypeNode48:
		return (*Node48)(r.ptr())
	case TypeNode256:
		return (*Node256)(r.ptr())
	default:
		panic("invalid node type")
	}
}

func (r *Ref) Replace(new Node) {
	*r = new.Ref()
}

func (r Ref) ptr() unsafe.Pointer {
	return unsafe.Pointer(uintptr(r) & nodePtrMask)
}
