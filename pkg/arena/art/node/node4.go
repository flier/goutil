package node

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
)

// Node4 represents the smallest node type in an adaptive radix tree, capable of
// storing up to 4 children. It is the entry point for most tree operations and
// provides the most memory-efficient storage for nodes with few children.
//
// Node4 uses a simple linear array representation where:
// - Keys are stored in ascending order for efficient binary search
// - Children are stored in the same order as their corresponding keys
// - Both arrays have a fixed size of 4 elements
//
// This design prioritizes memory efficiency over lookup performance for small
// node counts, making it ideal for sparse trees or tree nodes near the leaves.
type Node4 struct {
	Base
	Keys     [4]byte // Array of key bytes, sorted in ascending order
	Children [4]Ref  // Array of child node references, corresponding to Keys
}

var _ Node = (*Node4)(nil)

// Type returns the node type identifier for Node4.
func (n *Node4) Type() Type { return TypeNode4 }

// Full returns true if the node has reached its maximum capacity of 4 children.
func (n *Node4) Full() bool { return n.NumChildren == 4 }

// Ref returns a reference to this Node4 instance.
func (n *Node4) Ref() Ref { return NewRef(TypeNode4, n) }

// Minimum returns the leftmost leaf node in the subtree rooted at this node.
// Since keys are sorted, the first child contains the minimum key.
func (n *Node4) Minimum() *Leaf {
	if n.NumChildren == 0 {
		return nil
	}
	return n.Children[0].AsNode().Minimum()
}

// Maximum returns the rightmost leaf node in the subtree rooted at this node.
// Since keys are sorted, the last child contains the maximum key.
func (n *Node4) Maximum() *Leaf {
	if n.NumChildren == 0 {
		return nil
	}
	return n.Children[n.NumChildren-1].AsNode().Maximum()
}

// FindChild returns the child node for the given key byte.
//
// The method performs a linear search through the sorted keys array to find
// a matching key. While not optimal for larger node types, this approach
// is efficient for Node4 due to its small size and provides good cache locality.
//
// Returns:
//   - A pointer to the child reference if found
//   - nil if no child exists for the given key byte
//
// Time complexity: O(n) where n is the number of children (max 4)
func (n *Node4) FindChild(b byte) *Ref {
	for i := 0; i < n.NumChildren; i++ {
		if n.Keys[i] == b {
			return &n.Children[i]
		}
	}

	return nil
}

// AddChild adds a child node to the node while maintaining key ordering.
//
// The method inserts the new key in sorted order by shifting existing keys
// and children to make room. This preserves the sorted invariant required
// for efficient Minimum/Maximum operations and binary search capabilities.
// If a key already exists, it replaces the existing child.
//
// Precondition: The node must not be full (n.NumChildren < 4)
// Postcondition: Keys remain sorted, children are properly aligned, NumChildren is incremented (unless replacing)
func (n *Node4) AddChild(b byte, child AsRef) {
	// Check if key already exists
	for i := 0; i < n.NumChildren; i++ {
		if n.Keys[i] == b {
			// Replace existing child
			n.Children[i] = child.Ref()
			return
		}
	}

	debug.Assert(!n.Full(), "node must not be full")

	i := 0

	for ; i < n.NumChildren; i++ {
		if b < n.Keys[i] {
			break
		}
	}

	copy(n.Keys[i+1:], n.Keys[i:])
	copy(n.Children[i+1:], n.Children[i:])

	n.Keys[i] = b
	n.Children[i] = child.Ref()
	n.NumChildren++
}

// Grow converts this Node4 to a Node16 when it reaches capacity.
//
// This method is called when a Node4 becomes full and needs to accommodate
// more children. It creates a new Node16 and transfers all existing children
// while preserving the sorted key order.
//
// The conversion maintains the existing key-value mappings and allows the
// tree to continue growing efficiently.
func (n *Node4) Grow(a *arena.Arena) Node {
	newNode := arena.New(a, Node16{Base: n.Base})

	copy(newNode.Keys[:], n.Keys[:n.NumChildren])
	copy(newNode.Children[:], n.Children[:n.NumChildren])

	return newNode
}
