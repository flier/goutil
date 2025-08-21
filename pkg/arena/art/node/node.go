// Package node implements the core node types for an Adaptive Radix Tree (ART).
//
// ART trees are a space-efficient data structure for storing and retrieving
// string keys. This package provides four different node types that automatically
// adapt their storage strategy based on the number of children:
//
//   - Node4: Stores up to 4 children using sorted arrays (most memory efficient)
//   - Node16: Stores up to 16 children using sorted arrays (balanced)
//   - Node48: Stores up to 48 children using sparse arrays (efficient for medium nodes)
//   - Node256: Stores up to 256 children using direct arrays (fastest lookup)
//
// The tree automatically grows and shrinks nodes as children are added or removed,
// optimizing both memory usage and lookup performance. All nodes use arena allocation
// for efficient memory management and support prefix compression to reduce memory overhead.
//
// Key Features:
//   - Automatic node type adaptation based on child count
//   - Prefix compression for memory efficiency
//   - Arena-based memory allocation
//   - Generic value types with compile-time type safety
//   - Efficient ordered traversal support
package node

import (
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

// Type represents the type identifier for different node implementations in the ART tree.
//
// Each node type has different characteristics and capacity limits, allowing the tree
// to adapt its structure based on the number of children at each node.
//
// The type system follows a progression from memory-efficient (Node4) to
// performance-optimized (Node256) implementations, with automatic conversion
// between types as the tree structure evolves.
type Type int

const (
	// TypeUnknown represents an unknown or invalid node type.
	//
	// This is typically used for uninitialized references or error conditions.
	TypeUnknown Type = iota

	// TypeLeaf represents a leaf node containing key-value pairs.
	//
	// Leaf nodes are terminal nodes that cannot have children and store
	// the actual data associated with complete keys.
	TypeLeaf

	// TypeNode4 represents a small node storing up to 4 children.
	//
	// This is the most memory-efficient node type, ideal for sparse
	// trees or nodes near the leaves with few children.
	TypeNode4

	// TypeNode16 represents a medium node storing up to 16 children.
	//
	// This type provides a good balance between memory usage and
	// lookup performance for moderately populated nodes.
	TypeNode16

	// TypeNode48 represents a large node storing up to 48 children.
	//
	// This type uses sparse array storage for efficient memory usage
	// while maintaining good lookup performance.
	TypeNode48

	// TypeNode256 represents the largest node type storing up to 256 children.
	//
	// This type provides the fastest possible lookup performance using
	// direct array indexing, but uses the most memory.
	TypeNode256
)

// Node is the core interface for all node types in the Adaptive Radix Tree (ART).
//
// It provides a unified interface for operations like finding children, adding/removing
// children, and managing node growth/shrinking based on capacity requirements.
//
// The interface is designed to support the adaptive nature of ART trees, where nodes
// can dynamically change their type (e.g., from Node4 to Node16) as the number of
// children grows or shrinks, optimizing both memory usage and performance.
//
// All node implementations must satisfy this interface, ensuring consistent behavior
// across different node types while allowing for type-specific optimizations.
//
// Generic Type Parameter:
//   - T: The type of values stored in leaf nodes of this tree
type Node[T any] interface {
	// AsRef[T] embeds the reference interface, allowing nodes to be
	// treated as references for tree traversal and manipulation.
	AsRef[T]

	// Type returns the specific type identifier for this node implementation.
	//
	// This is useful for type assertions and determining the node's capabilities.
	// The returned type will be one of the Type constants defined above.
	Type() Type

	// Full returns true if the node has reached its maximum capacity and cannot
	// accommodate more children without growing to a larger node type.
	// When this returns true, calling AddChild will trigger automatic growth.
	Full() bool

	// Prefix returns the shared prefix bytes that all keys in this subtree have in common.
	//
	// This prefix compression is a key optimization in ART trees, reducing memory usage
	// by storing common key prefixes only once at each node level.
	// The prefix is stored as a slice.Slice[byte] for efficient memory management.
	Prefix() slice.Slice[byte]

	// SetPrefix updates the shared prefix for this node and its subtree.
	//
	// This is typically called during tree restructuring operations like node splitting
	// or merging to maintain the correct prefix relationships.
	// The prefix parameter should be a valid slice.Slice[byte] instance.
	SetPrefix(prefix slice.Slice[byte])

	// Minimum returns the leftmost leaf node in the subtree rooted at this node.
	//
	// This is useful for ordered traversal operations and finding the smallest key.
	// Returns nil if the subtree is empty or contains no leaf nodes.
	Minimum() *Leaf[T]

	// Maximum returns the rightmost leaf node in the subtree rooted at this node.
	//
	// This is useful for ordered traversal operations and finding the largest key.
	// Returns nil if the subtree is empty or contains no leaf nodes.
	Maximum() *Leaf[T]

	// FindChild locates and returns a reference to the child node associated with
	// the given key byte. Returns nil if no child exists for that byte value.
	//
	// The key byte represents the next character in the key being searched.
	// This method is the core of tree traversal and search operations.
	FindChild(b byte) *Ref[T]

	// AddChild adds a new child node to this node, associating it with the given key byte.
	// If the node becomes full, it may need to grow to a larger node type.
	//
	// The child parameter must implement AsRef[T] to provide a reference interface.
	// If a child with the same key already exists, it will be replaced.
	AddChild(b byte, child AsRef[T])

	// RemoveChild removes the child node associated with the given key byte.
	//
	// The child parameter is used to verify the removal operation and may trigger
	// node shrinking if the remaining children are few enough.
	// This method maintains the tree's structural integrity during deletions.
	RemoveChild(b byte, child *Ref[T])

	// Grow converts this node to a larger node type when it reaches capacity.
	//
	// The new node type will have more storage space for children while maintaining
	// all existing child relationships. This is called automatically by AddChild.
	// The arena allocator is used for memory allocation during the conversion.
	Grow(a arena.Allocator) Node[T]

	// Shrink converts this node to a smaller node type when it has few children.
	//
	// This optimization reduces memory usage for sparsely populated nodes.
	// The shrinking threshold varies by node type to balance memory and performance.
	// The arena allocator is used for memory management during the conversion.
	Shrink(a arena.AllocatorExt) Node[T]

	// Release frees all memory associated with this node and its children.
	//
	// This should be called when the node is no longer needed to prevent memory leaks.
	// The arena allocator is used to properly deallocate all allocated memory.
	// After calling Release, the node should not be used again.
	Release(a arena.Allocator)
}

// Base provides common functionality shared by all node implementations.
//
// It contains the shared prefix and child count that every node type needs,
// reducing code duplication and ensuring consistent behavior across implementations.
//
// The Base struct is embedded in all concrete node types, providing a foundation
// for common operations while allowing each node type to implement its own
// storage strategy for children.
type Base struct {
	// Partial stores the shared prefix bytes for this subtree.
	//
	// This prefix compression reduces memory usage by storing common
	// key prefixes only once at each node level.
	Partial slice.Slice[byte]

	// NumChildren tracks the current number of children in this node.
	//
	// This count is used to determine when to grow or shrink the node
	// and for various tree operations that need to know the node's state.
	NumChildren int
}

// Prefix returns the shared prefix bytes for this node.
//
// This method satisfies the Node interface requirement for prefix access.
// The returned slice represents the common prefix shared by all keys in this subtree.
func (n *Base) Prefix() slice.Slice[byte] { return n.Partial }

// SetPrefix updates the shared prefix for this node.
//
// This method satisfies the Node interface requirement for prefix modification.
// The prefix parameter should be a valid slice.Slice[byte] instance.
// This method is typically called during tree restructuring operations.
func (n *Base) SetPrefix(prefix slice.Slice[byte]) { n.Partial = prefix }
