package node

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/simd"
)

// Package node implements the core node types for an adaptive radix tree.
// Node48 is one of the intermediate node types that provides efficient
// storage and lookup for nodes with 17-48 children.

// Node48 represents a node in an adaptive radix tree that can store up to 48 children.
// It is an intermediate node type that provides a balance between memory usage and
// lookup performance for nodes with more than 16 children but fewer than 256.
//
// Node48 uses a sparse array representation where:
// - Keys[byte] stores the index into the Children array (1-based indexing)
// - Children stores the actual child node references
// - A key byte maps to a child through Keys[byte] -> Children[Keys[byte]-1]
//
// This design allows for efficient lookups while maintaining reasonable memory usage
// compared to Node256, which would use 256 pointers regardless of actual children count.
type Node48 struct {
	Base
	Keys     [256]byte // Maps key bytes to indices in Children array (0 = no child)
	Children [48]Ref   // Array of child node references
}

var _ Node = (*Node48)(nil)

// Type returns the node type identifier for Node48.
func (n *Node48) Type() Type { return TypeNode48 }

// Full returns true if the node has reached its maximum capacity of 48 children.
func (n *Node48) Full() bool { return n.NumChildren == 48 }

// Ref returns a reference to this Node48 instance.
func (n *Node48) Ref() Ref { return NewRef(TypeNode48, n) }

// Minimum returns the leftmost leaf node in the subtree rooted at this node.
// It traverses through the first non-empty child to find the minimum key.
func (n *Node48) Minimum() *Leaf {
	if n.NumChildren == 0 {
		return nil
	}
	// Find the first non-zero key in the Keys array
	if i := simd.FindNonZeroKeyIndex(&n.Keys); i >= 0 {
		return n.Children[n.Keys[i]-1].AsNode().Minimum()
	}

	return nil
}

// Maximum returns the rightmost leaf node in the subtree rooted at this node.
// It traverses through the last non-empty child to find the maximum key.
func (n *Node48) Maximum() *Leaf {
	if n.NumChildren == 0 {
		return nil
	}
	// Find the last non-zero key in the Keys array
	if i := simd.FindLastNonZeroKeyIndex(&n.Keys); i >= 0 {
		return n.Children[n.Keys[i]-1].AsNode().Maximum()
	}

	return nil
}

// FindChild returns the child node for the given key byte.
//
// The method uses the sparse array representation to efficiently locate
// the child node. It first checks if the key byte maps to a valid child
// index (non-zero), then returns the corresponding child reference.
//
// Returns:
//   - A pointer to the child reference if found
//   - nil if no child exists for the given key byte
//
// Time complexity: O(1) - direct array access
func (n *Node48) FindChild(b byte) *Ref {
	if idx := n.Keys[b]; idx != 0 {
		return &n.Children[idx-1]
	}

	return nil
}

// AddChild adds a child node to the node.
//
// The method finds the first available slot in the Children array and maps
// the key byte to that slot using 1-based indexing. This ensures that
// a key byte of 0 can be distinguished from "no child" (which is also 0).
// If a key already exists, it replaces the existing child.
//
// Precondition: The node must not be full (n.NumChildren < 48)
// Postcondition: The child is added and NumChildren is incremented (unless replacing)
func (n *Node48) AddChild(b byte, child AsRef) {
	// Check if key already exists
	if idx := n.Keys[b]; idx != 0 {
		// Replace existing child
		n.Children[idx-1] = child.Ref()
		return
	}

	debug.Assert(!n.Full(), "node must not be full")

	var i byte

	for ; i < 48; i++ {
		if n.Children[i] == 0 {
			break
		}
	}

	n.Keys[b] = i + 1
	n.Children[i] = child.Ref()
	n.NumChildren++
}

// Grow converts this Node48 to a Node256 when it reaches capacity.
//
// This method is called when a Node48 becomes full and needs to accommodate
// more children. It creates a new Node256 and transfers all existing children
// to their corresponding positions in the larger node structure.
//
// The conversion preserves the existing key-value mappings while allowing
// the tree to continue growing efficiently.
func (n *Node48) Grow(a *arena.Arena) Node {
	newNode := arena.New(a, Node256{Base: n.Base})

	for i := 0; i < 256; i++ {
		if n.Keys[i] != 0 {
			newNode.Children[i] = n.Children[n.Keys[i]-1]
		}
	}

	return newNode
}
