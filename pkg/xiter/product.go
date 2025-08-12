//go:build go1.23

package xiter

import "iter"

// Produce iterates over the entire iterator, multiplying all the elements
//
// An empty iterator returns the one value of the type.
func Product[T Number](x iter.Seq[T]) (p T) {
	p = 1

	for v := range x {
		p *= v
	}

	return
}
