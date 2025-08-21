package node

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

// Leaf represents a leaf node in the Adaptive Radix Tree (ART).
//
// Leaf nodes are the terminal nodes that store the actual key-value pairs.
// Unlike internal nodes, leaves cannot have children and represent the end
// of a key path in the tree.
//
// The Leaf type is generic over the value type T, allowing the tree to store
// any type of values while maintaining type safety at compile time.
//
// Leaf nodes are immutable in terms of their structure - they cannot gain or
// lose children, but their key and value can be modified through the SetPrefix
// method and direct field access respectively.
//
// Memory Management:
//   - Leaf nodes use arena allocation for efficient memory management
//   - The key is stored as a slice.Slice[byte] for memory efficiency
//   - All memory is properly managed through the arena allocator
//
// Generic Type Parameter:
//   - T: The type of values stored in this leaf node
type Leaf[T any] struct {
	// Key stores the complete key bytes for this leaf.
	//
	// This represents the full key path from the root to this leaf.
	// The key is stored as a slice.Slice[byte] for efficient memory management
	// and to support prefix operations.
	Key slice.Slice[byte]

	// Value stores the actual data associated with this key.
	//
	// The type T can be any Go type, providing flexibility for different use cases.
	// Common types include strings, integers, pointers, or custom structs.
	Value T
}

// Ensure Leaf implements the Node interface at compile time.
//
// This compile-time check ensures that Leaf satisfies all Node interface requirements.
var _ Node[any] = (*Leaf[any])(nil)

// NewLeaf creates a new leaf node with the given key and value.
//
// The key is converted to a slice.Slice[byte] for efficient memory management
// and prefix operations. The arena allocator is used to allocate memory for
// both the leaf structure and the key slice.
//
// Parameters:
//   - a: The arena allocator for memory management. Must not be nil.
//   - key: The byte slice representing the key. Can be nil or empty.
//   - value: The value to associate with the key. Can be any type T.
//
// Returns:
//   - A pointer to the newly created leaf node.
//
// Memory Allocation:
//   - The leaf structure is allocated in the arena
//   - The key is converted to a slice.Slice[byte] and allocated in the arena
//   - All memory is managed by the arena allocator
//
// Example:
//
//	a := &arena.Arena{}
//	leaf := NewLeaf(a, []byte("hello"), "world")
func NewLeaf[T any](a arena.Allocator, key []byte, value T) *Leaf[T] {
	debug.Assert(a != nil, "arena must not be nil")
	debug.Assert(len(key) > 0, "key must not be nil or empty")

	return arena.New(a, Leaf[T]{slice.FromBytes(a, key), value})
}

// Type returns the node type identifier for Leaf nodes.
//
// Leaf nodes always return TypeLeaf since they represent terminal nodes.
func (l *Leaf[T]) Type() Type { return TypeLeaf }

// Full always returns true for leaf nodes since they cannot have children.
//
// Leaf nodes are always considered "full" because they cannot accommodate
// additional children - they represent the end of a key path.
func (l *Leaf[T]) Full() bool { return true }

// Ref returns a reference to this leaf node.
//
// The reference can be used to traverse the tree and access the leaf
// from parent nodes without direct pointer access.
func (l *Leaf[T]) Ref() Ref[T] { return NewRef[T](TypeLeaf, l) }

// Prefix returns the complete key as the prefix for this leaf.
//
// Since leaves represent complete keys, their prefix is the entire key.
// The returned prefix represents the full key path from root to this leaf.
func (l *Leaf[T]) Prefix() slice.Slice[byte] { return l.Key }

// SetPrefix updates the key for this leaf node.
//
// This method satisfies the Node interface requirement and allows
// the leaf's key to be modified during tree restructuring operations.
// The prefix parameter becomes the new key for this leaf.
//
// Note: This method is typically called during tree restructuring operations
// like node splitting or merging, not during normal usage.
func (l *Leaf[T]) SetPrefix(prefix slice.Slice[byte]) { l.Key = prefix }

// Minimum returns the leaf itself since it has no children.
//
// Since leaves are terminal nodes, they are both the minimum and maximum
// of their own subtree.
func (l *Leaf[T]) Minimum() *Leaf[T] { return l }

// Maximum returns the leaf itself since it has no children.
//
// Since leaves are terminal nodes, they are both the minimum and maximum
// of their own subtree.
func (l *Leaf[T]) Maximum() *Leaf[T] { return l }

// FindChild panics since leaf nodes cannot have children.
//
// If this method is called, it indicates a programming error in the tree
// implementation or usage.
func (l *Leaf[T]) FindChild(b byte) *Ref[T] { panic("leaf cannot have children") }

// AddChild panics since leaf nodes cannot have children.
//
// If this method is called, it indicates a programming error in the tree
// implementation or usage.
func (l *Leaf[T]) AddChild(b byte, child AsRef[T]) { panic("leaf cannot have children") }

// RemoveChild panics since leaf nodes cannot have children.
//
// If this method is called, it indicates a programming error in the tree
// implementation or usage.
func (l *Leaf[T]) RemoveChild(b byte, child *Ref[T]) { panic("leaf cannot have children") }

// Grow panics since leaf nodes cannot have children.
//
// If this method is called, it indicates a programming error in the tree
// implementation or usage.
func (l *Leaf[T]) Grow(a arena.Allocator) Node[T] { panic("leaf cannot have children") }

// Shrink panics since leaf nodes cannot have children.
//
// If this method is called, it indicates a programming error in the tree
// implementation or usage.
func (l *Leaf[T]) Shrink(a arena.AllocatorExt) Node[T] { panic("leaf cannot have children") }

// Release frees all memory associated with this leaf node.
//
// This includes the key slice and the leaf structure itself.
// The arena allocator is used to properly deallocate all allocated memory.
// After calling this method, the leaf should not be used again.
//
// Parameters:
//   - a: The arena allocator for memory management. Must not be nil.
//
// Memory Deallocation:
//   - The key slice is released back to the arena
//   - The leaf structure itself is freed
//   - All memory is properly returned to the arena allocator
func (l *Leaf[T]) Release(a arena.Allocator) {
	l.Key.Release(a)

	arena.Free(a, l)
}

// Matches checks if this leaf's key matches the given key.
//
// The comparison is done using slice.EqualTo for efficient byte-by-byte
// comparison. This method is useful for verifying key matches during
// search and delete operations.
//
// Parameters:
//   - key: The byte slice to compare against this leaf's key.
//     Can be nil or empty.
//
// Returns:
//   - true if the keys match exactly (same length and content)
//   - false if the keys differ in any way
//
// Performance:
//   - Uses slice.EqualTo for efficient comparison
//   - Early termination on length mismatch
//   - Optimized for byte-by-byte comparison
//
// Example:
//
//	leaf := NewLeaf(a, []byte("hello"), "world")
//	if leaf.Matches([]byte("hello")) {
//	    // Keys match
//	}
func (l *Leaf[T]) Matches(key []byte) bool {
	return slice.EqualTo(l.Key, key)
}
