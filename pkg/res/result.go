// Error handling with the Result type.
//
// Result[T] is the type used for returning and propagating errors.
// It is an enum with the variants, Ok(T), representing success and containing a value,
// and Err(E), representing error and containing an error value.
package res

import (
	"fmt"
)

// Result is a type that represents either success (Ok) or failure (Err).
type Result[T any] struct {
	Value *T
	Err   error
}

// Contains the success value.
func Ok[T any](value T) Result[T] { return Result[T]{&value, nil} }

// Contains the error value.
func Err[T any](err error) Result[T] { return Result[T]{nil, err} }

// Wrap a value and error
func Wrap[T any](value T, err error) Result[T] {
	if err != nil {
		return Result[T]{nil, err}
	}

	return Result[T]{&value, nil}
}

func (r Result[T]) String() string {
	if r.IsOk() {
		return fmt.Sprintf("Ok(%v)", r.unwrap())
	}

	return fmt.Sprintf("Err(%v)", r.Err)
}

// Returns true if the result is Ok.
func (r Result[T]) IsOk() bool { return r.Value != nil }

// Returns true if the result is Ok and the value inside of it matches a predicate.
func (r Result[T]) IsOkAnd(f func(T) bool) bool { return r.IsOk() && f(r.unwrap()) }

// Returns true if the result is Err.
func (r Result[T]) IsErr() bool { return r.Err != nil }

// Returns true if the result is Err and the value inside of it matches a predicate.
func (r Result[T]) IsErrAnd(f func(error) bool) bool { return r.IsErr() && f(r.Err) }

// Returns the contained Ok value, or panics if the value is an Err,
// with a panic message including the passed message, and the content of the Err..
func (r Result[T]) Expect(msg string) T {
	if r.IsErr() {
		unwrapFail("%s: %s", msg, r.Err)
	}

	return r.unwrap()
}

// Returns the contained Err value, or panics if the value is an Ok,
// with a panic message including the passed message, and the content of the Ok.
func (r Result[T]) ExpectErr(msg string) error {
	if r.IsOk() {
		unwrapFail("%s: %v", msg, r.unwrap())
	}

	return r.Err
}

// Returns the contained Ok value, or panics if the value is an Err.
func (r Result[T]) Unwrap() T {
	return r.Expect("called `Result.Unwrap()` on an `Err` value")
}

// Returns the contained Ok value or a provided default value.
func (r Result[T]) UnwrapOr(def T) T {
	if r.IsOk() {
		return r.unwrap()
	}

	return def
}

// Returns the contained Ok value or a default.
func (r Result[T]) UnwrapOrDefault() (v T) {
	if r.IsOk() {
		v = r.unwrap()
	}

	return
}

// Returns the contained Ok value or computes it from a closure.
func (r Result[T]) UnwrapOrElse(f func() T) T {
	if r.IsOk() {
		return r.unwrap()
	}

	return f()
}

// Returns the contained Err value, or panics if the value is an Ok.
func (r Result[T]) UnwrapErr() error {
	return r.ExpectErr("called `Result.UnwrapErr()` on an `Ok` value")
}

func (r Result[T]) unwrap() T { return *r.Value }

func unwrapFail(format string, a ...any) { panic(fmt.Sprintf(format, a...)) }
