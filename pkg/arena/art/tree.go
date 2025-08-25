package art

import (
	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/art/tree"
)

// Tree represents an Adaptive Radix Tree.
//
// It is a generic type that can store any type of value.
type Tree[T any] struct {
	root node.Ref[T]
}

// Search searches for a value in the tree.
//
// It returns the value if found, otherwise nil.
func (t *Tree[T]) Search(key []byte) *T {
	return tree.Search(t.root, key)
}

// Minimum returns the minimum leaf in the tree.
//
// It returns nil if the tree is empty.
func (t *Tree[T]) Minimum() *node.Leaf[T] {
	if t.root.Empty() {
		return nil
	}

	return t.root.AsNode().Minimum()
}

// Maximum returns the maximum leaf in the tree.
//
// It returns nil if the tree is empty.
func (t *Tree[T]) Maximum() *node.Leaf[T] {
	if t.root.Empty() {
		return nil
	}

	return t.root.AsNode().Maximum()
}

// Insert inserts a new value into the tree.
//
// It returns the old value if the key matches the existing key, or nil if the key is inserted.
func (t *Tree[T]) Insert(a arena.Allocator, key []byte, value T) *T {
	return tree.RecursiveInsert(a, &t.root, node.NewLeaf(a, key, value), 0, true)
}

// InsertNoReplace inserts a new value into the tree without replacing the existing value.
//
// It returns the old value if the key matches the existing key, or nil if the key is inserted.
func (t *Tree[T]) InsertNoReplace(a arena.Allocator, key []byte, value T) *T {
	return tree.RecursiveInsert(a, &t.root, node.NewLeaf(a, key, value), 0, false)
}

// Delete deletes a value from the tree.
//
// It returns the old value if the key matches the existing key, or nil if the key is not found.
func (t *Tree[T]) Delete(a arena.AllocatorExt, key []byte) *T {
	l := tree.RecursiveDelete(a, &t.root, key, 0)
	if l == nil {
		return nil
	}

	old := l.Value

	arena.Free(a, l)

	return &old
}

// Visit visits the tree.
//
// It returns true if the iteration is interrupted by the callback function,
// otherwise it returns false.
func (t *Tree[T]) Visit(cb func(key []byte, value *T) bool) bool {
	return tree.RecursiveIter(t.root, cb)
}

// VisitPrefix visits the tree with a prefix.
//
// It returns true if the iteration is interrupted by the callback function,
// otherwise it returns false.
func (t *Tree[T]) VisitPrefix(prefix []byte, cb func(key []byte, value *T) bool) bool {
	return tree.IterPrefix(t.root, prefix, cb)
}
