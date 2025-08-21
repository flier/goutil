package tree

import (
	"github.com/flier/goutil/pkg/arena/art/node"
)

// Search searches for a key in the ART tree.
//
// It returns the value pointer if the key is found, otherwise it returns nil.
func Search[T any](ref node.Ref[T], key []byte) *T {
	var depth int

	for !ref.Empty() {
		// If the current node is a leaf, we need to check if the key matches
		if l := ref.AsLeaf(); l != nil {
			if l.Matches(key) {
				return &l.Value
			}

			return nil
		}

		curr := ref.AsNode()

		// Check if the key has the same prefix as the current node
		if partial := curr.Prefix(); partial.Len() > 0 {
			if prefixMatch := CheckPrefix(partial, key, depth); prefixMatch != partial.Len() {
				return nil
			}

			depth += partial.Len()
		}

		var b byte
		if depth < len(key) {
			b = key[depth]
		}

		// Recursively search
		child := curr.FindChild(b)
		if child == nil {
			break
		}

		ref = *child
		depth++
	}

	return nil
}
