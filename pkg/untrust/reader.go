package untrust

import (
	"io"
	"math"
)

// The error type used to indicate the end of the input was reached before the operation could be completed.
var ErrEndOfInput = io.ErrUnexpectedEOF

// A read-only, forward-only cursor into the data in an Input.
type Reader struct {
	b []byte
	i int
}

// Construct a new Reader for the given input.
//
// Use [ReadAll] or [ReadAllOptional] instead of [NewReader] whenever possible.
func NewReader(i Input) *Reader { return &Reader{b: i, i: 0} }

func (r *Reader) GoString() string { return "Reader" }

// Returns a copy of the Reader.
func (r *Reader) Clone() *Reader { return &Reader{b: r.b, i: r.i} }

// Returns true if the reader is at the end of the input, and false otherwise.
func (r *Reader) AtEnd() bool { return r.i == len(r.b) }

// Returns true if there is at least one more byte in the input and that byte is equal to b, and false otherwise.
func (r *Reader) Peek(b byte) bool { return len(r.b) > r.i && r.b[r.i] == b }

// Reads the next input byte.
func (r *Reader) ReadByte() (byte, error) {
	if len(r.b) <= r.i {
		return 0, ErrEndOfInput
	}

	b := r.b[r.i]
	r.i++

	return b, nil
}

// Skips n bytes of the input, returning the skipped input as an Input.
func (r *Reader) ReadBytes(n int) (Input, error) {
	if n < 0 || r.i > math.MaxInt-n {
		return nil, ErrEndOfInput
	}

	i := r.i + n

	if len(r.b) < i {
		r.i = len(r.b)

		return nil, ErrEndOfInput
	}

	b := r.b[r.i:i]
	r.i = i

	return Input(b), nil
}

// Skips the reader to the end of the input, returning the skipped input as an `Input`.
func (r *Reader) ReadBytesToEnd() (Input, error) {
	return r.ReadBytes(len(r.b) - r.i)
}

// Skips n bytes of the input.
func (r *Reader) Skip(n int) error {
	_, err := r.ReadBytes(n)

	return err
}

// Skips the reader to the end of the input.
func (r *Reader) SkipToEnd() error {
	_, err := r.ReadBytesToEnd()

	return err
}

// Calls read with the given input as a [Reader].
//
// On success, returns a pair (bytes_read, r)
// where bytes_read is what read consumed and r is read's return value.
func ReadPartial[T any](r *Reader, read func(*Reader) (T, error)) (Input, T, error) {
	start := r.i
	res, err := read(r)
	b := r.b[start:r.i]

	return Input(b), res, err
}
