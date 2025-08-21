package art

import "github.com/flier/goutil/pkg/arena/art/node"

type Tree[T any] struct {
	root node.Ref[T]
}
