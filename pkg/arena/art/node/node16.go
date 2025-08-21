package node

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/simd"
	"github.com/flier/goutil/pkg/xunsafe"
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
//
// Memory Layout:
//   - Keys array: 16 bytes (fixed size)
//   - Children array: 16 pointers (fixed size)
//   - Base struct: prefix + child count
//   - Total overhead: moderate for medium node counts
//
// Performance Characteristics:
//   - Lookup: O(n) where n ≤ 16 (linear search with SIMD optimization)
//   - Insertion: O(n) with shifting for sorted order
//   - Memory: Good balance between efficiency and performance
//   - Growth: Automatic conversion to Node48 when full
//
// SIMD Optimization:
//   - Uses AVX2 instructions for key search operations on AMD64
//   - Falls back to scalar implementation on other architectures
//   - Provides significant performance improvement for key lookups
//
// Generic Type Parameter:
//   - T: The type of values stored in leaf nodes of this tree
type Node16[T any] struct {
	// Base embeds the common functionality shared by all node types.
	//
	// This includes the shared prefix and child count.
	Base

	// Keys stores the key bytes in ascending order.
	//
	// The array has a fixed size of 16 elements, with only the first
	// NumChildren elements containing valid keys.
	// Keys are maintained in sorted order for efficient operations.
	Keys [16]byte

	// Children stores the child node references corresponding to Keys.
	//
	// The array has a fixed size of 16 elements, with only the first
	// NumChildren elements containing valid references.
	// Children[i] corresponds to Keys[i] for all valid indices.
	Children [16]Ref[T]
}

// Ensure Node16 implements the Node interface at compile time.
//
// This compile-time check ensures that Node16 satisfies all Node interface requirements.
var _ Node[any] = (*Node16[any])(nil)

// Type returns the node type identifier for Node16.
//
// Node16 nodes always return TypeNode16 since they represent the medium node type.
func (n *Node16[T]) Type() Type { return TypeNode16 }

// Full returns true if the node has reached its maximum capacity of 16 children.
//
// When this returns true, calling AddChild will trigger automatic growth to Node48.
// This method is used by the tree implementation to determine when to grow nodes.
func (n *Node16[T]) Full() bool { return n.NumChildren == 16 }

// Ref returns a reference to this Node16 instance.
//
// The reference can be used to traverse the tree and access this node
// from parent nodes without direct pointer access.
func (n *Node16[T]) Ref() Ref[T] { return NewRef[T](TypeNode16, n) }

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
func (n *Node16[T]) Minimum() *Leaf[T] {
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
func (n *Node16[T]) Maximum() *Leaf[T] {
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
// The search is optimized using SIMD instructions on AMD64 architectures,
// providing significant performance improvement for key lookups.
//
// Parameters:
//   - b: The key byte to search for
//
// Returns:
//   - A pointer to the child reference if found
//   - nil if no child exists for the given key byte
//
// Performance:
//   - Time complexity: O(n) where n is the number of children (max 16)
//   - Space complexity: O(1)
//   - Cache locality: Good due to moderate array size
//   - SIMD acceleration: Available on AMD64 for improved performance
//
// Algorithm:
//   - SIMD-optimized search through sorted keys array
//   - Early termination on match
//   - Returns corresponding child reference
func (n *Node16[T]) FindChild(b byte) *Ref[T] {
	if i := simd.FindKeyIndex(&n.Keys, n.NumChildren, b); i >= 0 {
		return &n.Children[i]
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
// The key finding operations are optimized using SIMD instructions for better
// performance on modern processors.
//
// Precondition: The node must not be full (n.NumChildren < 16)
// Postcondition: Keys remain sorted, children are properly aligned, NumChildren is incremented (unless replacing)
//
// Parameters:
//   - b: The key byte to associate with the child
//   - child: The child node to add. Must implement AsRef[T].
//
// Algorithm:
//   - Use SIMD-optimized search to find insertion position
//   - Shift existing keys and children to make room
//   - Insert new key and child at correct position
//   - Increment child count
//
// Performance:
//   - Time complexity: O(n) where n is the number of children
//   - Space complexity: O(1) (fixed array size)
//   - Memory operations: Array shifting for sorted order
//   - SIMD acceleration: For finding insertion position
func (n *Node16[T]) AddChild(b byte, child AsRef[T]) {
	debug.Assert(!n.Full(), "node must not be full")

	// Check if key already exists using SIMD-optimized search
	i := simd.FindInsertPosition(&n.Keys, n.NumChildren, b)
	if i >= 0 {
		// Key exists, shift elements to make room for insertion
		copy(n.Keys[i+1:], n.Keys[i:])
		copy(n.Children[i+1:], n.Children[i:])
	} else {
		// Key doesn't exist, insert at the end
		i = n.NumChildren
	}

	// Insert the new key and child at the correct position
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
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Returns:
//   - A new Node48 containing all existing children
//   - The original Node16 is no longer valid after this call
//
// Memory Management:
//   - New Node48 is allocated in the arena
//   - All children are copied to the new node
//   - Original Node16 memory is not freed (caller's responsibility)
//
// Conversion Process:
//   - Create new Node48 with same base information
//   - Copy all children to the new node
//   - Map keys to sparse array positions using 1-based indexing
//   - Preserve all child relationships
//
// Performance:
//   - Time complexity: O(n) where n is the number of children
//   - Space complexity: O(1) (fixed array sizes)
//   - Memory allocation: One Node48 structure
func (n *Node16[T]) Grow(a arena.Allocator) Node[T] {
	newNode := arena.New(a, Node48[T]{Base: n.Base})

	// Copy all existing children to the new Node48
	copy(newNode.Children[:], n.Children[:n.NumChildren])

	// Map each key to its position in the sparse array using 1-based indexing
	for i := 0; i < n.NumChildren; i++ {
		newNode.Keys[n.Keys[i]] = byte(i + 1)
	}

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
func (n *Node16[T]) RemoveChild(b byte, child *Ref[T]) {
	// Calculate the position of the child in the arrays
	pos := xunsafe.AddrOf(child).Sub(xunsafe.AddrOf(&n.Children[0]))

	debug.Assert(pos < n.NumChildren, "child must be in the node")

	// Shift remaining keys and children to fill the gap
	copy(n.Keys[pos:], n.Keys[pos+1:])
	copy(n.Children[pos:], n.Children[pos+1:])

	n.NumChildren--
}

// Shrink shrinks the node to a Node4 if it has less than 3 children.
//
// This method is called when a Node16 has few children and can be optimized
// by converting to a smaller Node4. The shrinking threshold is set to 3
// to balance memory efficiency with performance.
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Returns:
//   - A new Node4 if shrinking occurs (NumChildren < 3)
//   - The original Node16 if shrinking is not beneficial (NumChildren >= 3)
//
// Shrinking Logic:
//   - If NumChildren >= 3: return self (no shrinking beneficial)
//   - If NumChildren < 3: create new Node4 and copy children
//
// Memory Management:
//   - New Node4 is allocated in the arena if shrinking occurs
//   - All children are copied to the new node
//   - Original Node16 is freed if shrinking occurs
//
// Performance:
//   - Time complexity: O(n) where n is the number of children
//   - Space complexity: O(1) (fixed array sizes)
//   - Memory allocation: One Node4 structure (if shrinking)
func (n *Node16[T]) Shrink(a arena.AllocatorExt) Node[T] {
	if n.NumChildren >= 3 {
		return n
	}

	// Create new Node4 with same base information
	newNode := arena.New(a, Node4[T]{Base: n.Base})

	// Copy all existing keys and children to the new Node4
	copy(newNode.Keys[:], n.Keys[:n.NumChildren])
	copy(newNode.Children[:], n.Children[:n.NumChildren])

	// Free the original Node16 since we're replacing it
	arena.Free(a, n)

	return newNode
}

// Release frees all memory associated with this Node16 instance.
//
// This method frees all memory associated with this Node16 instance.
// It should be called when the node is no longer needed to prevent memory leaks.
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Memory Deallocation:
//   - The prefix slice is released back to the arena
//   - The Node16 structure itself is freed
//   - All memory is properly returned to the arena allocator
func (n *Node16[T]) Release(a arena.Allocator) {
	n.Partial.Release(a)

	arena.Free(a, n)
}
