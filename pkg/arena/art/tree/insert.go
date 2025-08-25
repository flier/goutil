package tree

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
)

func RecursiveInsert[T any](a arena.Allocator, ref *node.Ref[T], leaf *node.Leaf[T], depth int, replace bool) *T {
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

// InsertToLeaf inserts a leaf into a leaf node.
//
//   - If the leaf matches the key, we need to return the old value, or replace the value if replace is true.
//   - If the leaf does not match the key, we need to split the leaf into a node4.
//
// Do not use this method directly, use [RecursiveInsert] instead.
func InsertToLeaf[T any](a arena.Allocator, ref *node.Ref[T], leaf *node.Leaf[T], depth int, replace bool) *T {
	debug.Assert(ref.IsLeaf(), "current node must be a leaf")

	curr := ref.AsLeaf()

	// If the leaf matches the key, we need to return the old value
	if slice.Equal(curr.Key, leaf.Key) {
		old := curr.Value

		if replace {
			curr.Value = leaf.Value
		}

		return &old
	}

	// If the leaf does not match the key, we need to split the leaf into a node4
	newNode := arena.New(a, node.Node4[T]{})

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

// InsertToNode inserts a leaf into a node.
//
//   - If the node has a prefix, we need to check if the key has the same prefix.
//   - If the node does not have a prefix, we need to insert the leaf to the node.
//
// Returns the old value if the key matches the existing key, or nil if the key is inserted.
//
// Do not use this method directly, use [RecursiveInsert] instead.
func InsertToNode[T any](a arena.Allocator, ref *node.Ref[T], leaf *node.Leaf[T], depth int, replace bool) *T {
	debug.Assert(ref.IsNode(), "current node must be a node")

	// If the ref is a node, we need to split the node into a node4
	n := ref.AsNode()

	// If the node has a prefix, we need to check if the key has the same prefix
	if partial := n.Prefix(); !partial.Empty() {
		if diff := PrefixMismatch(n, leaf.Key, depth); diff >= partial.Len() {
			depth += partial.Len()
		} else {
			// If the key has the same prefix, we need to add the prefix to the new node
			newNode := arena.New(a, node.Node4[T]{})
			newNode.Partial = partial.Slice(0, diff).Clone(a)

			// Add the current node to the new node
			newNode.AddChild(n.Prefix().CheckedLoad(diff).UnwrapOrDefault(), n)
			n.SetPrefix(partial.Slice(diff+1, partial.Len()))

			// Add the leaf to the new node
			newNode.AddChild(leaf.Key.CheckedLoad(depth+diff).UnwrapOrDefault(), leaf)

			ref.Replace(newNode)

			return nil
		}
	}

	key := leaf.Key.CheckedLoad(depth).UnwrapOrDefault()

	// If the child is found, we need to recurse
	if child := n.FindChild(key); child != nil && !child.Empty() {
		return RecursiveInsert(a, child, leaf, depth+1, replace)
	}

	AddChild(a, ref, key, leaf)

	return nil
}

func AddChild[T any](a arena.Allocator, ref *node.Ref[T], key byte, leaf *node.Leaf[T]) {
	debug.Assert(ref.IsNode(), "current node must be a node")

	curr := ref.AsNode()

	// If the child is not found, we need to insert a new leaf
	if curr.Full() {
		newNode := curr.Grow(a)
		newNode.AddChild(key, leaf)

		ref.Replace(newNode)

		if newNode != curr {
			curr.Release(a)
		}
	} else {
		curr.AddChild(key, leaf)
	}
}
