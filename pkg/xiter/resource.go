//go:build go1.23

package xiter

import "iter"

// Resource emits an iterator of values for the given resource.
func Resource[S, T any](start func() (S, error), next func(S) (T, error), stop func(S)) iter.Seq[T] {
	return func(yield func(T) bool) {
		s, err := start()
		if err != nil {
			return
		}

		defer stop(s)

		for {
			v, err := next(s)
			if err != nil {
				break
			}

			if !yield(v) {
				break
			}
		}
	}
}
