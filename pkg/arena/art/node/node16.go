package node

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
)

// Node16 represents a medium-sized node in an adaptive radix tree, capable of
// storing up to 16 children. It serves as an intermediate node type that balances
// memory efficiency with lookup performance for nodes that have outgrown Node4.
//
// Node16 uses a sorted array representation similar to Node4 but with increased
// capacity. The design maintains the sorted key order for efficient traversal
// operations while providing more storage space for growing subtrees.
//
// This node type is ideal for nodes that have more than 4 children but fewer
// than 17, offering a good compromise between memory usage and performance.
type Node16 struct {
	Base
	Keys     [16]byte // Array of key bytes, sorted in ascending order
	Children [16]Ref  // Array of child node references, corresponding to Keys
}

var _ Node = (*Node16)(nil)

// Type returns the node type identifier for Node16.
func (n *Node16) Type() Type { return TypeNode16 }

// Full returns true if the node has reached its maximum capacity of 16 children.
func (n *Node16) Full() bool { return n.NumChildren == 16 }

// Ref returns a reference to this Node16 instance.
func (n *Node16) Ref() Ref { return NewRef(TypeNode16, n) }

// Minimum returns the leftmost leaf node in the subtree rooted at this node.
// Since keys are sorted, the first child contains the minimum key.
func (n *Node16) Minimum() *Leaf {
	if n.NumChildren == 0 {
		return nil
	}
	return n.Children[0].AsNode().Minimum()
}

// Maximum returns the rightmost leaf node in the subtree rooted at this node.
// Since keys are sorted, the last child contains the maximum key.
func (n *Node16) Maximum() *Leaf {
	if n.NumChildren == 0 {
		return nil
	}
	return n.Children[n.NumChildren-1].AsNode().Maximum()
}

// FindChild returns the child node for the given key byte.
//
// The method performs a linear search through the sorted keys array to find
// a matching key. While linear search has O(n) complexity, it remains efficient
// for Node16 due to its moderate size and provides good cache locality.
//
// For larger node types, more sophisticated search algorithms (like binary search
// or SIMD operations) would be more appropriate, but Node16's size makes linear
// search a practical choice.
//
// Returns:
//   - A pointer to the child reference if found
//   - nil if no child exists for the given key byte
//
// Time complexity: O(n) where n is the number of children (max 16)
func (n *Node16) FindChild(b byte) *Ref {
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
// for efficient Minimum/Maximum operations and maintains the tree's structural
// properties. If a key already exists, it replaces the existing child.
//
// Precondition: The node must not be full (n.NumChildren < 16)
// Postcondition: Keys remain sorted, children are properly aligned, NumChildren is incremented (unless replacing)
func (n *Node16) AddChild(b byte, child AsRef) {
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

// Grow converts this Node16 to a Node48 when it reaches capacity.
//
// This method is called when a Node16 becomes full and needs to accommodate
// more children. It creates a new Node48 and transfers all existing children
// while converting from the sorted array representation to the sparse array
// representation used by Node48.
//
// The conversion involves mapping each key to its position in the new sparse
// array structure, ensuring that all existing key-value relationships are
// preserved in the larger node.
func (n *Node16) Grow(a *arena.Arena) Node {
	newNode := arena.New(a, Node48{Base: n.Base})

	copy(newNode.Children[:], n.Children[:n.NumChildren])

	for i := 0; i < n.NumChildren; i++ {
		newNode.Keys[n.Keys[i]] = byte(i + 1)
	}

	return newNode
}
