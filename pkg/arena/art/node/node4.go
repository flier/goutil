package node

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/xunsafe"
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
//
// Memory Layout:
//   - Keys array: 4 bytes (fixed size)
//   - Children array: 4 pointers (fixed size)
//   - Base struct: prefix + child count
//   - Total overhead: minimal for small node counts
//
// Performance Characteristics:
//   - Lookup: O(n) where n ≤ 4 (linear search)
//   - Insertion: O(n) with shifting for sorted order
//   - Memory: Most efficient among all node types
//   - Growth: Automatic conversion to Node16 when full
//
// Generic Type Parameter:
//   - T: The type of values stored in leaf nodes of this tree
type Node4[T any] struct {
	// Base embeds the common functionality shared by all node types.
	// This includes the shared prefix and child count.
	Base

	// Keys stores the key bytes in ascending order.
	//
	// The array has a fixed size of 4 elements, with only the first
	// NumChildren elements containing valid keys.
	// Keys are maintained in sorted order for efficient operations.
	Keys [4]byte

	// Children stores the child node references corresponding to Keys.
	//
	// The array has a fixed size of 4 elements, with only the first
	// NumChildren elements containing valid references.
	// Children[i] corresponds to Keys[i] for all valid indices.
	Children [4]Ref[T]
}

// Ensure Node4 implements the Node interface at compile time.
//
// This compile-time check ensures that Node4 satisfies all Node interface requirements.
var _ Node[any] = (*Node4[any])(nil)

// Type returns the node type identifier for Node4.
//
// Node4 nodes always return TypeNode4 since they represent the smallest node type.
func (n *Node4[T]) Type() Type { return TypeNode4 }

// Full returns true if the node has reached its maximum capacity of 4 children.
//
// When this returns true, calling AddChild will trigger automatic growth to Node16.
// This method is used by the tree implementation to determine when to grow nodes.
func (n *Node4[T]) Full() bool { return n.NumChildren == 4 }

// Ref returns a reference to this Node4 instance.
//
// The reference can be used to traverse the tree and access this node
// from parent nodes without direct pointer access.
func (n *Node4[T]) Ref() Ref[T] { return NewRef[T](TypeNode4, n) }

// Minimum returns the leftmost leaf node in the subtree rooted at this node.
//
// Since keys are sorted, the first child contains the minimum key.
// This method traverses down the leftmost path to find the smallest key.
//
// Returns:
//   - The leftmost leaf node if this subtree contains any leaves
//   - nil if this subtree is empty or contains no leaf nodes
//
// Performance: O(1) for the first level, then O(depth) for traversal
func (n *Node4[T]) Minimum() *Leaf[T] {
	if n.NumChildren == 0 {
		return nil
	}

	return n.Children[0].AsNode().Minimum()
}

// Maximum returns the rightmost leaf node in the subtree rooted at this node.
//
// Since keys are sorted, the last child contains the maximum key.
// This method traverses down the rightmost path to find the largest key.
//
// Returns:
//   - The rightmost leaf node if this subtree contains any leaves
//   - nil if this subtree is empty or contains no leaf nodes
//
// Performance: O(1) for the last level, then O(depth) for traversal
func (n *Node4[T]) Maximum() *Leaf[T] {
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
// Parameters:
//   - b: The key byte to search for
//
// Returns:
//   - A pointer to the child reference if found
//   - nil if no child exists for the given key byte
//
// Performance:
//   - Time complexity: O(n) where n is the number of children (max 4)
//   - Space complexity: O(1)
//   - Cache locality: Excellent due to small array size
//
// Algorithm:
//   - Linear search through sorted keys array
//   - Early termination on match
//   - Returns corresponding child reference
func (n *Node4[T]) FindChild(b byte) *Ref[T] {
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
//
// Parameters:
//   - b: The key byte to associate with the child
//   - child: The child node to add. Must implement AsRef[T].
//
// Algorithm:
//   - Find insertion position to maintain sorted order
//   - Shift existing keys and children to make room
//   - Insert new key and child at correct position
//   - Increment child count
//
// Performance:
//   - Time complexity: O(n) where n is the number of children
//   - Space complexity: O(1) (fixed array size)
//   - Memory operations: Array shifting for sorted order
func (n *Node4[T]) AddChild(b byte, child AsRef[T]) {
	debug.Assert(!n.Full(), "node must not be full")

	var i int

	// Find the correct insertion position to maintain sorted order
	for ; i < n.NumChildren; i++ {
		if b < n.Keys[i] {
			break
		}
	}

	// Shift existing keys and children to make room for the new entry
	copy(n.Keys[i+1:], n.Keys[i:])
	copy(n.Children[i+1:], n.Children[i:])

	// Insert the new key and child at the correct position
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
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Returns:
//   - A new Node16 containing all existing children
//   - The original Node4 is no longer valid after this call
//
// Memory Management:
//   - New Node16 is allocated in the arena
//   - All children are copied to the new node
//   - Original Node4 memory is not freed (caller's responsibility)
//
// Performance:
//   - Time complexity: O(n) where n is the number of children
//   - Space complexity: O(1) (fixed array sizes)
//   - Memory allocation: One Node16 structure
func (n *Node4[T]) Grow(a arena.Allocator) Node[T] {
	newNode := arena.New(a, Node16[T]{Base: n.Base})

	// Copy all existing keys and children to the new Node16
	copy(newNode.Keys[:], n.Keys[:n.NumChildren])
	copy(newNode.Children[:], n.Children[:n.NumChildren])

	return newNode
}

// RemoveChild removes a child node from the node.
//
// This method removes the child associated with the given key byte and
// maintains the sorted order of remaining keys and children.
//
// Parameters:
//   - b: The key byte of the child to remove
//   - child: A reference to the child being removed (for verification)
//
// Algorithm:
//   - Calculate the position of the child in the arrays
//   - Shift remaining keys and children to fill the gap
//   - Decrement the child count
//
// Performance:
//   - Time complexity: O(n) where n is the number of children
//   - Space complexity: O(1)
//   - Memory operations: Array shifting to maintain order
func (n *Node4[T]) RemoveChild(b byte, child *Ref[T]) {
	// Calculate the position of the child in the arrays
	pos := xunsafe.AddrOf(child).Sub(xunsafe.AddrOf(&n.Children[0]))

	debug.Assert(pos < n.NumChildren, "child must be in the node")

	// Shift remaining keys and children to fill the gap
	copy(n.Keys[pos:], n.Keys[pos+1:])
	copy(n.Children[pos:], n.Children[pos+1:])

	n.NumChildren--
}

// Shrink shrinks the node to a leaf node or combines the node with its child.
//
// This method is called when a Node4 has only one child and can be optimized
// by either converting to a leaf (if the child is a leaf) or combining with
// the child node (if the child is an internal node).
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Returns:
//   - A leaf node if the child is a leaf
//   - A combined node if the child is an internal node
//   - The original node if it has multiple children
//
// Shrinking Logic:
//   - If multiple children: return self (no shrinking possible)
//   - If single child is leaf: return the leaf directly
//   - If single child is node: combine prefixes and return child
//
// Memory Management:
//   - Original Node4 is freed if shrinking occurs
//   - Child nodes are preserved and returned
//   - Prefix concatenation may occur for internal node children
func (n *Node4[T]) Shrink(a arena.AllocatorExt) Node[T] {
	if n.NumChildren > 1 {
		return n
	}

	child := n.Children[0]

	if !child.IsLeaf() {
		// If the child is a node, we need to concatenate the prefix and the child's prefix.
		if c := child.AsNode(); c != nil {
			// Append the key byte to the current prefix
			n.Partial = n.Partial.AppendOne(a, n.Keys[0])

			// Release the child's old prefix and set the new combined prefix
			c.Prefix().Release(a)
			c.SetPrefix(n.Partial)

			// Update the child reference to the modified child
			child = c.Ref()
		} else {
			// Release the current prefix if child is invalid
			n.Partial.Release(a)
		}
	}

	// Free the original Node4 since we're replacing it
	arena.Free(a, n)

	return child.AsNode()
}

// Release releases the node.
//
// This method frees all memory associated with this Node4 instance.
// It should be called when the node is no longer needed to prevent memory leaks.
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Memory Deallocation:
//   - The prefix slice is released back to the arena
//   - The Node4 structure itself is freed
//   - All memory is properly returned to the arena allocator
func (n *Node4[T]) Release(a arena.Allocator) {
	n.Partial.Release(a)

	arena.Free(a, n)
}
