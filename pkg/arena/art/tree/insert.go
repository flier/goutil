//go:build go1.21

package tree

import (
	"unsafe"

	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
)

func RecursiveInsert(a *arena.Arena, ref *node.Ref, leaf *node.Leaf, depth int, replace bool) unsafe.Pointer {
	// If the ref is empty, we need to inject a leaf
	if ref.Empty() {
		ref.Replace(leaf)

		return nil
	}

	// If the ref is a leaf, we need to replace it with a node4, or update the existing leaf
	if ref.IsLeaf() {
		return InsertToLeaf(a, ref, leaf, depth, replace)
	}

	// If the ref is a node, we need to insert the leaf to the node
	return InsertToNode(a, ref, leaf, depth, replace)
}

func InsertToLeaf(a *arena.Arena, ref *node.Ref, leaf *node.Leaf, depth int, replace bool) unsafe.Pointer {
	curr := ref.AsLeaf()

	debug.Assert(curr != nil, "current node must be a leaf")

	// If the leaf matches the key, we need to return the old value
	if slice.Equal(curr.Key, leaf.Key) {
		old := curr.Value

		if replace {
			curr.Value = leaf.Value
		}

		return old
	}

	// If the leaf does not match the key, we need to split the leaf into a node4
	newNode := arena.New(a, node.Node4{})

	// If the key and the current key have a common prefix, we need to add it to the node4
	if i := LongestCommonPrefix(leaf.Key, curr.Key, depth); i > depth {
		newNode.Partial = leaf.Key.Slice(depth, i)

		depth = i
	}

	// Add the leafs to the new node4
	newNode.AddChild(leaf.Key.CheckedLoad(depth).UnwrapOrDefault(), leaf)
	newNode.AddChild(curr.Key.CheckedLoad(depth).UnwrapOrDefault(), ref)

	ref.Replace(newNode)

	return nil
}

func InsertToNode(a *arena.Arena, ref *node.Ref, leaf *node.Leaf, depth int, replace bool) unsafe.Pointer {
	// If the ref is a node, we need to split the node into a node4
	curr := ref.AsNode()

	debug.Assert(curr != nil, "current node must be a node")

	if partial := curr.Prefix(); !partial.Empty() {
		if diff := PrefixMismatch(curr, partial.Raw(), depth); diff >= partial.Len() {
			depth += partial.Len()
		} else {
			newNode := arena.New(a, node.Node4{})
			newNode.Partial = partial.Slice(0, diff)

			newNode.AddChild(curr.Prefix().CheckedLoad(diff).UnwrapOrDefault(), curr)
			curr.Prefix().SetLen(curr.Prefix().Len() - diff + 1)
		}
	}

	b := leaf.Key.CheckedLoad(depth).UnwrapOrDefault()

	// If the child is found, we need to recurse
	if child := curr.FindChild(b); child != nil {
		return RecursiveInsert(a, child, leaf, depth+1, replace)
	}

	// If the child is not found, we need to insert a new leaf
	AddChild(a, b, ref, leaf)

	return nil
}

func AddChild(a *arena.Arena, b byte, curr *node.Ref, child node.AsRef) {
	switch n := curr.AsNode().(type) {
	case *node.Node4:
		if n.NumChildren < 4 {
			n.AddChild(b, child)
		} else {
			newNode := n.Grow(a)
			newNode.AddChild(b, child)

			curr.Replace(newNode)
		}

	case *node.Node16:
		if n.NumChildren < 16 {
			n.AddChild(b, child)
		} else {
			newNode := n.Grow(a)
			newNode.AddChild(b, child)

			curr.Replace(newNode)
		}

	case *node.Node48:
		if n.NumChildren < 48 {
			n.AddChild(b, child)
		} else {
			newNode := n.Grow(a)
			newNode.AddChild(b, child)

			curr.Replace(newNode)
		}
	case *node.Node256:
		n.AddChild(b, child)

	default:
		panic("invalid node type")
	}
}

func LongestCommonPrefix(l slice.Slice[byte], r slice.Slice[byte], depth int) (i int) {
	n := min(l.Len(), r.Len())
	i = depth

	for i < n && l.Load(i) == r.Load(i) {
		i++
	}

	return
}

func PrefixMismatch(n node.Node, partial []byte, depth int) (i int) {
	key := n.Prefix().Raw()

	for ; i < min(len(key)-depth, len(partial)); i++ {
		if key[depth+i] != partial[i] {
			return i
		}
	}

	l := n.Minimum()

	for ; i < min(l.Key.Len(), len(key))-depth; i++ {
		if l.Key.Load(depth+i) != key[depth+i] {
			return i
		}
	}

	return
}
