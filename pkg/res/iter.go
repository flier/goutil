//go:build go1.23

package res

import "iter"

// Collect iterates over a sequence of Result[T] values, collecting all successful values into a slice.
//
// If any Result in the sequence is an error, Collect returns nil and the encountered error immediately.
// Otherwise, it returns the slice of unwrapped values and a nil error.
func Collect[T any](seq iter.Seq[Result[T]]) (values []T, err error) {
	for res := range seq {
		if res.IsErr() {
			return nil, res.Err
		}

		values = append(values, res.Unwrap())
	}

	return
}

// Collect2 iterates over a sequence of values and errors provided by seq.
//
// It collects all values into a slice until an error is encountered.
// If an error occurs during iteration, it returns nil and the error.
// Otherwise, it returns the collected values and a nil error.
func Collect2[T any](seq iter.Seq2[T, error]) (values []T, err error) {
	for value, err := range seq {
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return
}

// Returns an iterator over the possibly contained value.
//
// The iterator yields one value if the result is Ok, otherwise none.
func (o Result[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.IsOk() {
			yield(o.unwrap())
		}
	}
}
