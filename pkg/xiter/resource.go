//go:build go1.23

package xiter

import (
	"bufio"
	"io"
	"iter"
)

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

// Lines returns an iterator of lines from the given [io.ReadCloser].
func Lines(r io.ReadCloser) iter.Seq[string] {
	return Resource(
		func() (*bufio.Scanner, error) { return bufio.NewScanner(r), nil },
		func(s *bufio.Scanner) (string, error) { return s.Text(), s.Err() },
		func(s *bufio.Scanner) { _ = r.Close() },
	)
}
