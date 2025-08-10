package untrust

import (
	"bytes"

	"github.com/flier/goutil/pkg/opt"
)

// Input is a wrapper around []byte that helps in writing panic-free code.
type Input []byte

func (i Input) GoString() string { return "Input" }

// Returns true if the input is empty and false otherwise.
func (i Input) Empty() bool { return len(i) == 0 }

// Returns the length of the Input.
func (i Input) Len() int { return len(i) }

// Access the input as a slice so it can be processed by functions
// that are not written using the Input/Reader framework.
func (i Input) AsSliceLessSafe() []byte { return i }

// Clone returns a copy of the `Input`.
//
// The elements are copied using assignment, so this is a shallow clone.
func (i Input) Clone() Input { return bytes.Clone(i) }

// Calls read with the given input as a [Reader], ensuring that read consumed the entire input.
//
// If read does not consume the entire input, incomplete error is returned.
func ReadAll[T any](input Input, incomplete error, read func(r *Reader) (T, error)) (res T, err error) {
	r := NewReader(input)

	res, err = read(r)
	if err != nil {
		return
	}

	if !r.AtEnd() {
		err = incomplete
	}

	return
}

// Calls read with the given input as a [Reader], ensuring that read consumed the entire input.
//
// When input is None, read will be called with None.
func ReadAllOptional[T any](input opt.Option[Input], incomplete error, read func(r opt.Option[*Reader]) (T, error)) (res T, err error) {
	if input.IsSome() {
		r := NewReader(input.Unwrap())

		res, err = read(opt.Some(r))
		if err != nil {
			return
		}

		if !r.AtEnd() {
			err = incomplete
		}
	} else {
		res, err = read(opt.None[*Reader]())
	}

	return
}
