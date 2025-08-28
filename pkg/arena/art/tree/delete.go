package tree

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/node"
)

// RecursiveDelete finds and returns a leaf node that matches the given key.
func RecursiveDelete[T any](a arena.AllocatorExt, ref *node.Ref[T], key []byte, depth int) *node.Leaf[T] {
	if ref.Empty() {
		return nil
	}

	// If the ref is a leaf, we need to delete the leaf if it matches the key.
	if l := ref.AsLeaf(); l != nil {
		if l.Matches(key) {
			ref.Replace(nil)

			return l
		}

		return nil
	}

	n := ref.AsNode()

	// If the current node has a partial prefix, we need to check if the key matches the prefix.
	if partial := n.Prefix(); partial.Len() > 0 {
		if n := CheckPrefix(partial, key, depth); n != partial.Len() {
			return nil
		}

		depth += partial.Len()
	}

	// Check if depth exceeds key length
	if depth > len(key) {
		return nil
	}

	b := -1

	if depth < len(key) {
		b = int(key[depth])
	}

	// Find the child node
	child := n.FindChild(b)
	if child == nil {
		// If the child is not found, return nil
		return nil
	}

	// If the child is a leaf, check if it matches the key
	if l := child.AsLeaf(); l != nil {
		if l.Matches(key) {
			RemoveChild(a, ref, b, child)

			return l
		}

		return nil
	}

	// Recursively search in the child node
	return RecursiveDelete(a, child, key, depth+1)
}

// RemoveChild removes a child node from the current node.
func RemoveChild[T any](a arena.AllocatorExt, ref *node.Ref[T], key int, child *node.Ref[T]) {
	debug.Assert(ref.IsNode(), "ref must be a node")

	curr := ref.AsNode()
	curr.RemoveChild(key, child)

	if n := curr.Shrink(a); n != curr {
		ref.Replace(n)
	}
}
