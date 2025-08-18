package node

import (
	"github.com/flier/goutil/pkg/arena"
)

// Node256 represents the largest node type in an adaptive radix tree, capable of
// storing up to 256 children (one for each possible byte value). It is the final
// destination for nodes that have grown beyond the capacity of smaller node types.
//
// Node256 uses a direct array representation where each possible byte value (0-255)
// directly maps to a child reference. This design provides O(1) lookup performance
// but uses the most memory among all node types, as it allocates space for all
// 256 possible children regardless of how many are actually used.
//
// This node type is ideal for dense nodes that have many children across a wide
// range of byte values, such as nodes near the root of the tree or nodes in
// highly populated subtrees.
type Node256 struct {
	Base
	Children [256]Ref // Direct array mapping from byte values to child references
}

var _ Node = (*Node256)(nil)

// Type returns the node type identifier for Node256.
func (n *Node256) Type() Type { return TypeNode256 }

// Full returns true if the node has reached its maximum capacity of 256 children.
// Note that Node256 can theoretically store 256 children, but this limit is
// rarely reached in practice due to the sparse nature of most key distributions.
func (n *Node256) Full() bool { return n.NumChildren == 256 }

// Ref returns a reference to this Node256 instance.
func (n *Node256) Ref() Ref { return NewRef(TypeNode256, n) }

// Minimum returns the leftmost leaf node in the subtree rooted at this node.
//
// The method scans the children array from index 0 to find the first non-empty
// child, then recursively calls Minimum() on that child. This approach works
// because Node256 stores children in a direct byte-to-index mapping.
//
// Time complexity: O(1) in the best case (first child exists), O(256) in the worst case
func (n *Node256) Minimum() *Leaf {
	for i := 0; i < 256; i++ {
		if n.Children[i] != 0 {
			return n.Children[i].AsNode().Minimum()
		}
	}

	return nil
}

// Maximum returns the rightmost leaf node in the subtree rooted at this node.
//
// The method scans the children array from index 255 down to 0 to find the last
// non-empty child, then recursively calls Maximum() on that child. This approach
// works because Node256 stores children in a direct byte-to-index mapping.
//
// Time complexity: O(1) in the best case (last child exists), O(256) in the worst case
func (n *Node256) Maximum() *Leaf {
	for i := 255; i >= 0; i-- {
		if n.Children[i] != 0 {
			return n.Children[i].AsNode().Maximum()
		}
	}

	return nil
}

// FindChild returns the child node for the given key byte.
//
// Node256 provides the fastest possible lookup performance by using direct array
// indexing. The byte value directly maps to the corresponding child reference
// without any search or computation required.
//
// Returns:
//   - A pointer to the child reference if found
//   - A pointer to a zero Ref if no child exists for the given key byte
//
// Time complexity: O(1) - direct array access
func (n *Node256) FindChild(b byte) *Ref {
	return &n.Children[b]
}

// AddChild adds a child node to the node.
//
// The method directly assigns the child reference to the array position
// corresponding to the key byte. If this is a new child (replacing a zero
// reference), the NumChildren counter is incremented.
//
// Unlike smaller node types, Node256 does not need to maintain sorted order
// or shift elements, making insertion extremely fast.
//
// If a child with the same key already exists, it is replaced without
// affecting the NumChildren count.
func (n *Node256) AddChild(b byte, child AsRef) {
	if n.Children[b] == 0 {
		n.NumChildren++
	}

	n.Children[b] = child.Ref()
}

// Grow is a no-op for Node256 as it is the largest node type.
//
// Since Node256 can accommodate all possible byte values (0-255), it never
// needs to grow to a larger node type. This method simply returns the node
// itself, maintaining the Node interface contract.
//
// In practice, if a Node256 becomes too sparse (many unused slots), the tree
// implementation might consider shrinking it to a smaller node type for
// memory efficiency, but this would be handled at a higher level.
func (n *Node256) Grow(a *arena.Arena) Node {
	return n
}
