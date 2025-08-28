//go:build go1.21

package tree

import (
	"github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
)

func LongestCommonPrefix(l slice.Slice[byte], r slice.Slice[byte], depth int) (i int) {
	n := min(l.Len(), r.Len())
	i = depth

	for i < n && l.Load(i) == r.Load(i) {
		i++
	}

	return
}

// PrefixMismatch checks if the key has the same prefix as the partial.
//
// It returns the number of characters that match.
func PrefixMismatch[T any](n node.Node[T], key []byte, depth int) (i int) {
	partial := n.Prefix()

	for ; i < min(partial.Len(), len(key)-depth); i++ {
		if partial.Load(i) != key[depth+i] {
			return
		}
	}

	// If we have a minimum leaf, continue checking
	if l := n.Minimum(); l != nil {
		for ; i < min(l.Key.Len(), len(key))-depth; i++ {
			if l.Key.Load(depth+i) != key[depth+i] {
				return
			}
		}
	}

	return
}

// CheckPrefix checks if the key has the same prefix as the partial.
//
// It returns the number of characters that match.
func CheckPrefix(partial slice.Slice[byte], key []byte, depth int) (i int) {
	n := min(partial.Len(), len(key)-depth)

	for ; i < n; i++ {
		if partial.Load(i) != key[depth+i] {
			break
		}
	}

	return i
}
