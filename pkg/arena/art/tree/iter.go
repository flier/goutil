package tree

import (
	"github.com/flier/goutil/pkg/arena/art/node"
)

// RecursiveIter iterates over the tree using a callback function.
//
// It returns true if the iteration is interrupted by the callback function,
// otherwise it returns false.
func RecursiveIter[T any](ref node.Ref[T], cb func(key []byte, value *T) bool) bool {
	if ref.Empty() {
		return false
	}

	switch n := ref.AsNode().(type) {
	case *node.Leaf[T]:
		return cb(n.Key.Raw(), &n.Value)

	case *node.Node4[T]:
		if !n.ZeroSizedChild.Empty() {
			if RecursiveIter(n.ZeroSizedChild, cb) {
				return true
			}
		}

		for i := 0; i < n.NumChildren; i++ {
			if RecursiveIter(n.Children[i], cb) {
				return true
			}
		}

	case *node.Node16[T]:
		if !n.ZeroSizedChild.Empty() {
			if RecursiveIter(n.ZeroSizedChild, cb) {
				return true
			}
		}

		for i := 0; i < n.NumChildren; i++ {
			if RecursiveIter(n.Children[i], cb) {
				return true
			}
		}

	case *node.Node48[T]:
		if !n.ZeroSizedChild.Empty() {
			if RecursiveIter(n.ZeroSizedChild, cb) {
				return true
			}
		}

		for i := 0; i < 256; i++ {
			if idx := n.Keys[i]; idx != 0 {
				if RecursiveIter(n.Children[idx-1], cb) {
					return true
				}
			}
		}

	case *node.Node256[T]:
		if !n.ZeroSizedChild.Empty() {
			if RecursiveIter(n.ZeroSizedChild, cb) {
				return true
			}
		}

		for i := 0; i < 256; i++ {
			if !n.Children[i].Empty() {
				if RecursiveIter(n.Children[i], cb) {
					return true
				}
			}
		}
	}

	return false
}

// IterPrefix iterates over the tree with a prefix using a callback function.
//
// It returns true if the iteration is interrupted by the callback function,
// otherwise it returns false.
func IterPrefix[T any](ref node.Ref[T], prefix []byte, cb func(key []byte, value *T) bool) bool {
	var depth int

	for !ref.Empty() {
		if l := ref.AsLeaf(); l != nil {
			if l.MatchesPrefix(prefix) {
				return cb(l.Key.Raw(), &l.Value)
			}

			return false
		}

		n := ref.AsNode()

		// If the depth matches the prefix, we need to handle this node
		if depth == len(prefix) {
			if l := n.Minimum(); l != nil && l.MatchesPrefix(prefix) {
				return RecursiveIter(ref, cb)
			}

			return false
		}

		if p := n.Prefix(); p.Len() > 0 {
			prefixLen := PrefixMismatch(n, prefix, depth)

			if prefixLen > p.Len() {
				prefixLen = p.Len()
			}

			if prefixLen == 0 {
				return false
			} else if depth+prefixLen == len(prefix) {
				return RecursiveIter(n.Ref(), cb)
			}

			depth += p.Len()
		}

		child := n.FindChild(int(prefix[depth]))

		if child == nil {
			break
		}

		ref = *child
		depth++
	}

	return false
}
