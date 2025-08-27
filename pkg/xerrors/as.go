package xerrors

import "errors"

// AsA is a helper function to check if an error is of a specific type.
//
// AsA returns the error as the target type T if possible.
//
// This is a generic wrapper around [errors.As] for convenience.
func AsA[T error](err error) (_ T, ok bool) {
	var e T

	if ok := errors.As(err, &e); ok {
		return e, true
	}

	var zero T

	return zero, false
}
