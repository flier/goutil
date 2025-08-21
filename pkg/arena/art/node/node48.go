package node

import (
	"github.com/flier/goutil/internal/debug"
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/simd"
)

// Node48 represents a node in an adaptive radix tree that can store up to 48 children.
//
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
//
// Memory Layout:
//   - Keys array: 256 bytes (sparse mapping, most entries are 0)
//   - Children array: 48 pointers (fixed size)
//   - Base struct: prefix + child count
//   - Total overhead: moderate for large node counts
//
// Performance Characteristics:
//   - Lookup: O(1) - direct array access with sparse mapping
//   - Insertion: O(1) - direct assignment to sparse array
//   - Memory: Efficient for nodes with 17-48 children
//   - Growth: Automatic conversion to Node256 when full
//
// Sparse Array Design:
//   - Keys array uses 1-based indexing to distinguish from "no child" (0)
//   - Only 48 out of 256 possible key bytes can be used
//   - Memory overhead is fixed regardless of actual children count
//   - Excellent cache locality for the Children array
//
// Generic Type Parameter:
//   - T: The type of values stored in leaf nodes of this tree
type Node48[T any] struct {
	// Base embeds the common functionality shared by all node types.
	//
	// This includes the shared prefix and child count.
	Base

	// Keys maps key bytes to indices in the Children array.
	//
	// The array has 256 elements (one for each possible byte value).
	// A value of 0 indicates "no child" for that key byte.
	// Non-zero values are 1-based indices into the Children array.
	// This sparse representation allows O(1) key lookups.
	Keys [256]byte

	// Children stores the actual child node references.
	//
	// The array has a fixed size of 48 elements, with only the first
	// NumChildren elements containing valid references.
	// Keys[byte] maps to Children[Keys[byte]-1] for valid keys.
	Children [48]Ref[T]
}

// Ensure Node48 implements the Node interface at compile time.
//
// This compile-time check ensures that Node48 satisfies all Node interface requirements.
var _ Node[any] = (*Node48[any])(nil)

// Type returns the node type identifier for Node48.
//
// Node48 nodes always return TypeNode48 since they represent the large node type.
func (n *Node48[T]) Type() Type { return TypeNode48 }

// Full returns true if the node has reached its maximum capacity of 48 children.
//
// When this returns true, calling AddChild will trigger automatic growth to Node256.
// This method is used by the tree implementation to determine when to grow nodes.
func (n *Node48[T]) Full() bool { return n.NumChildren == 48 }

// Ref returns a reference to this Node48 instance.
//
// The reference can be used to traverse the tree and access this node
// from parent nodes without direct pointer access.
func (n *Node48[T]) Ref() Ref[T] { return NewRef[T](TypeNode48, n) }

// Minimum returns the leftmost leaf node in the subtree rooted at this node.
//
// It traverses through the first non-empty child to find the minimum key.
//
// This method scans the Keys array from index 0 to find the first non-zero
// entry, then recursively calls Minimum() on the corresponding child.
//
// Returns:
//   - The leftmost leaf node if this subtree contains any leaves
//   - nil if this subtree is empty or contains no leaf nodes
//
// Performance:
//   - Time complexity: O(1) in best case, O(256) in worst case
//   - Space complexity: O(1)
//   - SIMD acceleration: Available for finding first non-zero key
func (n *Node48[T]) Minimum() *Leaf[T] {
	if n.NumChildren == 0 {
		return nil
	}

	// Find the first non-zero key in the Keys array using SIMD optimization
	if i := simd.FindNonZeroKeyIndex(&n.Keys); i >= 0 {
		return n.Children[n.Keys[i]-1].AsNode().Minimum()
	}

	return nil
}

// Maximum returns the rightmost leaf node in the subtree rooted at this node.
//
// It traverses through the last non-empty child to find the maximum key.
//
// This method scans the Keys array from index 255 down to 0 to find the last
// non-zero entry, then recursively calls Maximum() on the corresponding child.
//
// Returns:
//   - The rightmost leaf node if this subtree contains any leaves
//   - nil if this subtree is empty or contains no leaf nodes
//
// Performance:
//   - Time complexity: O(1) in best case, O(256) in worst case
//   - Space complexity: O(1)
//   - SIMD acceleration: Available for finding last non-zero key
func (n *Node48[T]) Maximum() *Leaf[T] {
	if n.NumChildren == 0 {
		return nil
	}

	// Find the last non-zero key in the Keys array using SIMD optimization
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
// This method provides O(1) lookup performance by using direct array access
// with the sparse mapping stored in the Keys array.
//
// Parameters:
//   - b: The key byte to search for
//
// Returns:
//   - A pointer to the child reference if found
//   - nil if no child exists for the given key byte
//
// Performance:
//   - Time complexity: O(1) - direct array access
//   - Space complexity: O(1)
//   - Memory access: Two array lookups (Keys[b] then Children[index])
//
// Algorithm:
//   - Check Keys[b] for valid child index (non-zero)
//   - If valid, access Children[Keys[b]-1] (1-based indexing)
//   - Return child reference or nil if not found
func (n *Node48[T]) FindChild(b byte) *Ref[T] {
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
//
// Parameters:
//   - b: The key byte to associate with the child
//   - child: The child node to add. Must implement AsRef[T].
//
// Algorithm:
//   - Check if key already exists (Keys[b] != 0)
//   - If exists: replace existing child
//   - If new: find first available slot in Children array
//   - Map key to slot using 1-based indexing
//   - Increment NumChildren counter
//
// Performance:
//   - Time complexity: O(1) for replacement, O(48) for new child
//   - Space complexity: O(1)
//   - Memory operations: Direct assignment to sparse arrays
func (n *Node48[T]) AddChild(b byte, child AsRef[T]) {
	// Check if key already exists
	if idx := n.Keys[b]; idx != 0 {
		// Replace existing child
		n.Children[idx-1] = child.Ref()
		return
	}

	debug.Assert(!n.Full(), "node must not be full")

	// Find the first available slot in the Children array
	var i byte
	for ; i < 48; i++ {
		if n.Children[i] == 0 {
			break
		}
	}

	// Map the key byte to the slot using 1-based indexing
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
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Returns:
//   - A new Node256 containing all existing children
//   - The original Node48 is no longer valid after this call
//
// Memory Management:
//   - New Node256 is allocated in the arena
//   - All children are copied to the new node
//   - Original Node48 memory is not freed (caller's responsibility)
//
// Conversion Process:
//   - Create new Node256 with same base information
//   - Copy children from sparse array to direct array positions
//   - Preserve all key-to-child mappings
//   - Maintain child count and prefix information
//
// Performance:
//   - Time complexity: O(256) - scan all possible key positions
//   - Space complexity: O(1) (fixed array sizes)
//   - Memory allocation: One Node256 structure
func (n *Node48[T]) Grow(a arena.Allocator) Node[T] {
	newNode := arena.New(a, Node256[T]{Base: n.Base})

	// Copy children from sparse array to direct array positions
	for i := 0; i < 256; i++ {
		if n.Keys[i] != 0 {
			newNode.Children[i] = n.Children[n.Keys[i]-1]
		}
	}

	return newNode
}

// RemoveChild removes a child node from the node.
//
// This method removes the child associated with the given key byte and
// clears the corresponding entries in both the Keys and Children arrays.
//
// Parameters:
//   - b: The key byte of the child to remove
//   - child: A reference to the child being removed (for verification)
//
// Algorithm:
//   - Find the child's position in the Children array via Keys[b]
//   - Clear the key mapping (Keys[b] = 0)
//   - Clear the child reference (Children[index] = 0)
//   - Decrement the child count
//
// Performance:
//   - Time complexity: O(1) - direct array access
//   - Space complexity: O(1)
//   - Memory operations: Two array assignments
func (n *Node48[T]) RemoveChild(b byte, child *Ref[T]) {
	// Find the position of the child in the Children array
	idx := n.Keys[b]
	if idx == 0 {
		return // Key doesn't exist
	}

	// Clear the key and child
	n.Keys[b] = 0
	n.Children[idx-1] = 0
	n.NumChildren--
}

// Shrink shrinks the node to a Node16 if it has less than 12 children.
//
// This method is called when a Node48 has few children and can be optimized
// by converting to a smaller Node16. The shrinking threshold is set to 12
// to balance memory efficiency with performance.
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Returns:
//   - A new Node16 if shrinking occurs (NumChildren < 12)
//   - The original Node48 if shrinking is not beneficial (NumChildren >= 12)
//
// Shrinking Logic:
//   - If NumChildren >= 12: return self (no shrinking beneficial)
//   - If NumChildren < 12: create new Node16 and copy children
//
// Memory Management:
//   - New Node16 is allocated in the arena if shrinking occurs
//   - All children are copied to the new node
//   - Original Node48 is freed if shrinking occurs
//
// Conversion Process:
//   - Create new Node16 with same base information
//   - Scan Keys array for non-zero entries
//   - Copy valid children to new node in order
//   - Maintain sorted key order for Node16
//
// Performance:
//   - Time complexity: O(256) - scan all possible key positions
//   - Space complexity: O(1) (fixed array sizes)
//   - Memory allocation: One Node16 structure (if shrinking)
func (n *Node48[T]) Shrink(a arena.AllocatorExt) Node[T] {
	if n.NumChildren >= 12 {
		return n
	}

	// Create new Node16 with same base information
	newNode := arena.New(a, Node16[T]{Base: n.Base})

	// Copy children from sparse array to sorted array
	var child byte
	for i := 0; i < 256; i++ {
		if pos := n.Keys[i]; pos != 0 {
			newNode.Keys[child] = byte(i)
			newNode.Children[child] = n.Children[pos-1]
			child++
		}
	}

	// Free the original Node48 since we're replacing it
	arena.Free(a, n)

	return newNode
}

// Release frees all memory associated with this Node48 instance.
//
// This method frees all memory associated with this Node48 instance.
// It should be called when the node is no longer needed to prevent memory leaks.
//
// Parameters:
//   - a: The arena allocator for memory management
//
// Memory Deallocation:
//   - The prefix slice is released back to the arena
//   - The Node48 structure itself is freed
//   - All memory is properly returned to the arena allocator
func (n *Node48[T]) Release(a arena.Allocator) {
	n.Partial.Release(a)

	arena.Free(a, n)
}
