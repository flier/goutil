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
//
// Memory Layout:
//   - Children array: 256 pointers (fixed size, one for each byte value)
//   - Base struct: prefix + child count
//   - Total overhead: highest among all node types
//   - Memory usage: constant regardless of actual children count
//
// Performance Characteristics:
//   - Lookup: O(1) - direct array indexing
//   - Insertion: O(1) - direct assignment
//   - Memory: Highest usage but predictable overhead
//   - Growth: No growth possible (maximum size)
//   - Shrinking: Automatic conversion to Node48 when sparse
//
// Use Cases:
//   - Nodes with many children across diverse byte values
//   - Root nodes of large trees
//   - Performance-critical lookup operations
//   - When memory efficiency is less important than speed
//
// Generic Type Parameter:
//   - T: The type of values stored in leaf nodes of this tree
type Node256[T any] struct {
	// Base embeds the common functionality shared by all node types.
	//
	// This includes the shared prefix and child count.
	Base[T]

	// Children stores child node references in a direct array mapping.
	//
	// The array has 256 elements, one for each possible byte value (0-255).
	// A zero value (Ref[T](0)) indicates "no child" for that byte value.
	// Non-zero values are valid child references.
	// This direct mapping provides O(1) lookup performance.
	Children [256]Ref[T]
}

// Ensure Node256 implements the Node interface at compile time.
//
// This compile-time check ensures that Node256 satisfies all Node interface requirements.
var _ Node[any] = (*Node256[any])(nil)

// Type returns the node type identifier for Node256.
//
// Node256 nodes always return TypeNode256 since they represent the largest node type.
func (n *Node256[T]) Type() Type { return TypeNode256 }

// Full returns true if the node has reached its maximum capacity of 256 children.
//
// Note that Node256 can theoretically store 256 children, but this limit is
// rarely reached in practice due to the sparse nature of most key distributions.
// This method is included for interface compatibility but always returns false
// since Node256 can never be truly "full" in the same sense as smaller nodes.
func (n *Node256[T]) Full() bool { return n.NumChildren == 256 }

// Ref returns a reference to this Node256 instance.
//
// The reference can be used to traverse the tree and access this node
// from parent nodes without direct pointer access.
func (n *Node256[T]) Ref() Ref[T] { return NewRef[T](TypeNode256, n) }

// Minimum returns the leftmost leaf node in the subtree rooted at this node.
//
// The method scans the children array from index 0 to find the first non-empty
// child, then recursively calls Minimum() on that child. This approach works
// because Node256 stores children in a direct byte-to-index mapping.
//
// Returns:
//   - The leftmost leaf node if this subtree contains any leaves
//   - nil if this subtree is empty or contains no leaf nodes
//
// Performance:
//   - Time complexity: O(1) in the best case (first child exists), O(256) in the worst case
//   - Space complexity: O(1)
//   - Memory access: Sequential scan of children array
//
// Algorithm:
//   - Scan children array from index 0 to 255
//   - Find first non-zero child reference
//   - Recursively call Minimum() on that child
//   - Return result or nil if no children found
func (n *Node256[T]) Minimum() *Leaf[T] {
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
// Returns:
//   - The rightmost leaf node if this subtree contains any leaves
//   - nil if this subtree is empty or contains no leaf nodes
//
// Performance:
//   - Time complexity: O(1) in the best case (last child exists), O(256) in the worst case
//   - Space complexity: O(1)
//   - Memory access: Reverse sequential scan of children array
//
// Algorithm:
//   - Scan children array from index 255 down to 0
//   - Find last non-zero child reference
//   - Recursively call Maximum() on that child
//   - Return result or nil if no children found
func (n *Node256[T]) Maximum() *Leaf[T] {
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
// This method is the primary reason for using Node256 - it provides constant-time
// lookup performance regardless of how many children the node actually contains.
//
// Parameters:
//   - b: The key byte to search for (0-255)
//
// Returns:
//   - A pointer to the child reference if found
//   - A pointer to a zero Ref if no child exists for the given key byte
//
// Performance:
//   - Time complexity: O(1) - direct array access
//   - Space complexity: O(1)
//   - Memory access: Single array lookup
//   - Cache locality: Excellent for frequently accessed keys
//
// Algorithm:
//   - Direct array access: Children[b]
//   - Check if reference is non-zero
//   - Return pointer to reference or nil
func (n *Node256[T]) FindChild(b int) *Ref[T] {
	if b < 0 {
		if n.ZeroSizedChild.Empty() {
			return nil
		}

		return &n.ZeroSizedChild
	}

	k := byte(b)

	if !n.Children[k].Empty() {
		return &n.Children[k]
	}

	return nil
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
//
// Parameters:
//   - b: The key byte to associate with the child (0-255)
//   - child: The child node to add. Must implement AsRef[T].
//
// Algorithm:
//   - Check if position already has a child (Children[b] != 0)
//   - If no existing child: increment NumChildren counter
//   - Assign child reference to Children[b]
//   - No shifting or reordering required
//
// Performance:
//   - Time complexity: O(1) - direct array assignment
//   - Space complexity: O(1)
//   - Memory operations: Single array assignment
//   - No shifting or reordering overhead
func (n *Node256[T]) AddChild(b int, child AsRef[T]) {
	if b < 0 {
		n.ZeroSizedChild = child.Ref()

		return
	}

	k := byte(b)

	if n.Children[k] == 0 {
		n.NumChildren++
	}

	n.Children[k] = child.Ref()
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
//
// Parameters:
//   - a: The arena allocator (unused, included for interface compatibility)
//
// Returns:
//   - The node itself (no growth possible)
//
// Performance:
//   - Time complexity: O(1) - no operation performed
//   - Space complexity: O(1) - no memory allocation
func (n *Node256[T]) Grow(arena.Allocator) Node[T] {
	return n
}

// RemoveChild removes a child node from the node.
//
// This method removes the child associated with the given key byte by
// setting the corresponding array position to zero and decrementing
// the child count.
//
// Parameters:
//   - b: The key byte of the child to remove (0-255)
//   - child: A reference to the child being removed (for verification)
//
// Algorithm:
//   - Set Children[b] = 0 (clear the child reference)
//   - Decrement NumChildren counter
//   - No shifting or reordering required
//
// Performance:
//   - Time complexity: O(1) - direct array assignment
//   - Space complexity: O(1)
//   - Memory operations: Single array assignment
func (n *Node256[T]) RemoveChild(b int, child *Ref[T]) {
	if b < 0 {
		if &n.ZeroSizedChild == child {
			n.ZeroSizedChild = 0
		}

		return
	}

	k := byte(b)

	n.Children[k] = 0
	n.NumChildren--
}

// Shrink shrinks the node to a Node48 if it has less than 37 children.
//
// This method is called when a Node256 has few children and can be optimized
// by converting to a smaller Node48. The shrinking threshold is set to 37
// to balance memory efficiency with performance.
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Returns:
//   - A new Node48 if shrinking occurs (NumChildren < 37)
//   - The original Node256 if shrinking is not beneficial (NumChildren >= 37)
//
// Shrinking Logic:
//   - If NumChildren >= 37: return self (no shrinking beneficial)
//   - If NumChildren < 37: create new Node48 and copy children
//
// Memory Management:
//   - New Node48 is allocated in the arena if shrinking occurs
//   - All children are copied to the new node
//   - Original Node256 is freed if shrinking occurs
//
// Conversion Process:
//   - Create new Node48 with same base information
//   - Scan Children array for non-zero entries
//   - Copy valid children to sparse array positions
//   - Map keys to sparse array using 1-based indexing
//
// Performance:
//   - Time complexity: O(256) - scan all possible positions
//   - Space complexity: O(1) (fixed array sizes)
//   - Memory allocation: One Node48 structure (if shrinking)
func (n *Node256[T]) Shrink(a arena.AllocatorExt) Node[T] {
	if n.NumChildren >= 37 {
		return n
	}

	// Create new Node48 with same base information
	newNode := arena.New(a, Node48[T]{Base: n.Base})

	// Copy children from direct array to sparse array
	var pos byte
	for i := 0; i < 256; i++ {
		if n.Children[i] != 0 {
			newNode.Children[pos] = n.Children[i]
			newNode.Keys[i] = pos + 1
			pos++
		}
	}

	// Free the original Node256 since we're replacing it
	arena.Free(a, n)

	return newNode
}

// Release frees all memory associated with this Node256 instance.
//
// This method frees all memory associated with this Node256 instance.
// It should be called when the node is no longer needed to prevent memory leaks.
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Memory Deallocation:
//   - The prefix slice is released back to the arena
//   - The Node256 structure itself is freed
//   - All memory is properly returned to the arena allocator
func (n *Node256[T]) Release(a arena.Allocator) {
	n.Partial.Release(a)

	arena.Free(a, n)
}
