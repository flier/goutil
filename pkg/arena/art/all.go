//go:build go1.23

package art

import (
	"iter"

	"github.com/flier/goutil/pkg/arena/art/tree"
)

// All iterates over the tree.
//
// It returns a sequence of key-value pairs.
func (t *Tree[T]) All() iter.Seq2[[]byte, *T] {
	return func(yield func([]byte, *T) bool) {
		tree.RecursiveIter(t.root, func(key []byte, value *T) bool {
			return !yield(key, value)
		})
	}
}

// AllPrefix iterates over the tree with a prefix.
//
// It returns a sequence of key-value pairs.
func (t *Tree[T]) AllPrefix(prefix []byte) iter.Seq2[[]byte, *T] {
	return func(yield func([]byte, *T) bool) {
		tree.IterPrefix(t.root, prefix, func(key []byte, value *T) bool {
			return !yield(key, value)
		})
	}
}
