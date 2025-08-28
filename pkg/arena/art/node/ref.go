package node

import (
	"unsafe"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/xunsafe"
)

// AsRef is the interface for the node reference.
//
// This interface provides a way to obtain a reference to a node without
// knowing its specific type. All node types implement this interface.
//
// The AsRef interface is part of the reference system that allows nodes
// to be treated uniformly while maintaining type safety and efficient
// memory representation.
//
// Generic Type Parameter:
//   - T: The type of values stored in leaf nodes of this tree
type AsRef[T any] interface {
	// Ref returns the reference to the node.
	// This method provides access to the underlying Ref[T] value
	// that can be used for tree traversal and manipulation.
	Ref() Ref[T]
}

// Ref represents a type-safe reference to a node in the ART tree.
//
// It combines a pointer to the node with type information in a single
// uintptr value, providing efficient memory representation and type safety.
//
// The Ref type uses bit manipulation to store both the node pointer and
// type information in a single value, allowing for compact storage and
// fast type checking without additional memory overhead.
//
// Memory Layout:
//   - Lower bits: Node type identifier (Type)
//   - Upper bits: Node pointer address
//   - Total size: One uintptr (8 bytes on 64-bit systems)
//
// Type Safety:
//   - Compile-time type checking for value types
//   - Runtime type validation for node types
//   - Safe conversion between different node types
//
// Generic Type Parameter:
//   - T: The type of values stored in leaf nodes of this tree
type Ref[T any] uintptr

// NewRef creates a new reference to a node of the specified type.
//
// This function combines the node pointer with type information into
// a single Ref value that can be used throughout the tree.
//
// The function uses bit manipulation to encode both the node type and
// pointer address in a single uintptr value, ensuring efficient storage
// and fast type checking.
//
// Parameters:
//   - t: The type identifier for the node
//   - p: A pointer to the node instance
//
// Returns:
//   - A Ref[T] value containing both type and pointer information
//
// Bit Layout:
//   - Lower bits (nodeTypeMask): Type identifier
//   - Upper bits (nodePtrMask): Node pointer address
//
// Example:
//
//	node4 := arena.New(a, Node4[any]{})
//	ref := NewRef[any](TypeNode4, node4)
func NewRef[T, N any](t Type, p *N) Ref[T] {
	addr := xunsafe.AddrOf(p)

	return Ref[T]((uintptr(addr) & nodePtrMask) | (uintptr(t) & nodeTypeMask))
}

// Constants for bit manipulation in Ref values.
const (
	// nodePtrMask is used to extract the pointer address from a Ref value.
	// It masks out the type information, leaving only the pointer bits.
	nodePtrMask = ^nodeTypeMask

	// nodeTypeMask is used to extract the type information from a Ref value.
	// It masks out the pointer address, leaving only the type bits.
	// The size is determined by the arena alignment requirements.
	nodeTypeMask = uintptr(arena.Align - 1)
)

// Ref returns the reference to the node.
//
// This method satisfies the AsRef interface and provides access to the
// underlying Ref value. For Ref types, this simply returns the receiver.
func (r Ref[T]) Ref() Ref[T] { return r }

// Type returns the node type identifier for this reference.
//
// This method extracts the type information from the lower bits of the Ref value.
// The returned type can be used for type checking and determining node capabilities.
func (r Ref[T]) Type() Type { return Type(uintptr(r) & nodeTypeMask) }

// Empty returns true if this reference is empty (zero value).
//
// An empty reference indicates that no node is associated with this reference.
// This is commonly used to check for uninitialized or invalid references.
func (r Ref[T]) Empty() bool { return r == 0 }

// IsLeaf returns true if this reference points to a leaf node.
//
// Leaf nodes are terminal nodes that store key-value pairs and cannot have children.
func (r Ref[T]) IsLeaf() bool { return r.Type() == TypeLeaf }

// IsNode4 returns true if this reference points to a Node4.
//
// Node4 is the smallest node type, storing up to 4 children.
func (r Ref[T]) IsNode4() bool { return r.Type() == TypeNode4 }

// IsNode16 returns true if this reference points to a Node16.
//
// Node16 is a medium node type, storing up to 16 children.
func (r Ref[T]) IsNode16() bool { return r.Type() == TypeNode16 }

// IsNode48 returns true if this reference points to a Node48.
//
// Node48 is a large node type, storing up to 48 children.
func (r Ref[T]) IsNode48() bool { return r.Type() == TypeNode48 }

// IsNode256 returns true if this reference points to a Node256.
//
// Node256 is the largest node type, storing up to 256 children.
func (r Ref[T]) IsNode256() bool { return r.Type() == TypeNode256 }

// IsNode returns true if this reference points to any internal node type.
//
// Internal nodes are Node4, Node16, Node48, or Node256 - any node that can have children.
func (r Ref[T]) IsNode() bool { return r.IsNode4() || r.IsNode16() || r.IsNode48() || r.IsNode256() }

// AsLeaf returns a pointer to the Leaf node if this reference points to a leaf.
//
// This method provides type-safe access to leaf nodes for operations that
// specifically require leaf functionality.
//
// Returns:
//   - A pointer to the Leaf[T] if this reference points to a leaf node
//   - nil if this reference points to an internal node or is empty
//
// Type Safety:
//   - Returns *Leaf[T] only for leaf nodes
//   - Returns nil for internal nodes, preventing type errors
func (r Ref[T]) AsLeaf() *Leaf[T] {
	if r.IsLeaf() {
		return (*Leaf[T])(r.ptr())
	}

	return nil
}

// AsNode4 returns a pointer to the Node4 if this reference points to a Node4.
//
// This method provides type-safe access to Node4 nodes for operations that
// specifically require Node4 functionality.
//
// Returns:
//   - A pointer to the Node4[T] if this reference points to a Node4
//   - nil if this reference points to a different node type or is empty
//
// Type Safety:
//   - Returns *Node4[T] only for Node4 nodes
//   - Returns nil for other node types, preventing type errors
func (r Ref[T]) AsNode4() *Node4[T] {
	if r.IsNode4() {
		return (*Node4[T])(r.ptr())
	}

	return nil
}

// AsNode16 returns a pointer to the Node16 if this reference points to a Node16.
//
// This method provides type-safe access to Node16 nodes for operations that
// specifically require Node16 functionality.
//
// Returns:
//   - A pointer to the Node16[T] if this reference points to a Node16
//   - nil if this reference points to a different node type or is empty
//
// Type Safety:
//   - Returns *Node16[T] only for Node16 nodes
//   - Returns nil for other node types, preventing type errors
func (r Ref[T]) AsNode16() *Node16[T] {
	if r.IsNode16() {
		return (*Node16[T])(r.ptr())
	}

	return nil
}

// AsNode48 returns a pointer to the Node48 if this reference points to a Node48.
//
// This method provides type-safe access to Node48 nodes for operations that
// specifically require Node48 functionality.
//
// Returns:
//   - A pointer to the Node48[T] if this reference points to a Node48
//   - nil if this reference points to a different node type or is empty
//
// Type Safety:
//   - Returns *Node48[T] only for Node48 nodes
//   - Returns nil for other node types, preventing type errors
func (r Ref[T]) AsNode48() *Node48[T] {
	if r.IsNode48() {
		return (*Node48[T])(r.ptr())
	}

	return nil
}

// AsNode256 returns a pointer to the Node256 if this reference points to a Node256.
//
// This method provides type-safe access to Node256 nodes for operations that
// specifically require Node256 functionality.
//
// Returns:
//   - A pointer to the Node256[T] if this reference points to a Node256
//   - nil if this reference points to a different node type or is empty
//
// Type Safety:
//   - Returns *Node256[T] only for Node256 nodes
//   - Returns nil for other node types, preventing type errors
func (r Ref[T]) AsNode256() *Node256[T] {
	if r.IsNode256() {
		return (*Node256[T])(r.ptr())
	}

	return nil
}

// AsNode returns the node as a Node[T] interface if this reference is valid.
//
// This method provides a unified way to access any node type through the
// common Node interface, allowing for polymorphic operations.
//
// Returns:
//   - The node as a Node[T] interface if this reference is valid
//   - nil if this reference is empty
//
// Type Safety:
//   - Returns the appropriate concrete node type wrapped in the Node interface
//   - Panics if the reference contains an invalid node type (programming error)
//   - Provides access to all Node interface methods
//
// Panics:
//   - If the reference contains an invalid node type identifier
//
//go:nosplit
func (r Ref[T]) AsNode() Node[T] {
	if r == 0 {
		return nil
	}

	p := r.ptr()
	if p == nil {
		return nil
	}

	switch r.Type() {
	case TypeLeaf:
		return (*Leaf[T])(p)
	case TypeNode4:
		return (*Node4[T])(p)
	case TypeNode16:
		return (*Node16[T])(p)
	case TypeNode48:
		return (*Node48[T])(p)
	case TypeNode256:
		return (*Node256[T])(p)
	default:
		panic("invalid node type")
	}
}

// Replace updates this reference with a new node reference.
//
// This method is commonly used during tree restructuring operations
// to replace one node with another while maintaining the reference.
//
// Parameters:
//   - new: The new node reference to assign, or nil to clear the reference
//
// Returns:
//   - The previous node that was referenced (may be nil)
//
// Usage:
//   - Used during node growth/shrinking operations
//   - Used during tree restructuring and balancing
//   - Provides atomic reference replacement
//
// Example:
//
//	oldNode := ref.Replace(newNode)
//	if oldNode != nil {
//	    oldNode.Release(arena)
//	}
func (r *Ref[T]) Replace(new AsRef[T]) (current Node[T]) {
	current = r.AsNode()

	if new != nil {
		*r = new.Ref()
	} else {
		*r = 0
	}

	return
}

// ptr extracts the raw pointer from this reference.
//
// This method removes the type information and returns only the
// pointer address, which can be cast to the appropriate node type.
//
// Returns:
//   - An unsafe.Pointer to the node instance
//
// Safety:
//   - The returned pointer should only be used with the correct node type
//   - Type checking should be performed before using this method
//   - This method is primarily used internally by the type-safe accessors
func (r Ref[T]) ptr() unsafe.Pointer {
	return unsafe.Pointer(xunsafe.Addr[byte](uintptr(r) & nodePtrMask).AssertValid())
}
