//go:build go1.23

package opt

import "iter"

// Returns an iterator over the possibly contained value.
func (o Option[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.IsSome() {
			yield(o.unwrap())
		}
	}
}
