// Optional values.
//
// Type Option represents an optional value: every Option is either Some and contains a value, or None, and does not.
package opt

import "fmt"

// The Option type.
type Option[T any] struct {
	Value *T
}

// Some value of type T.
func Some[T any](value T) Option[T] { return Option[T]{&value} }

// No value.
func None[T any]() Option[T] { return Option[T]{nil} }

// Wrap a optional value of type T.
func Wrap[T any](value *T) Option[T] { return Option[T]{value} }

func (o Option[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("Some(%v)", o.unwrap())
	}

	return "None"
}

// Returns true if the option is a Some value.
func (o Option[T]) IsSome() bool { return o.Value != nil }

// Returns true if the option is a Some and the value inside of it matches a predicate.
func (o Option[T]) IsSomeAnd(f func(T) bool) bool { return o.IsSome() && f(o.unwrap()) }

// Returns true if the option is a None value.
func (o Option[T]) IsNone() bool { return o.Value == nil }

// Returns true if the option is a None or the value inside of it matches a predicate.
func (o Option[T]) IsNoneOr(f func(T) bool) bool { return o.IsNone() || f(o.unwrap()) }

// Returns the contained Some value, orp panics if the value is a None with a custom panic message provided by msg.
func (o Option[T]) Expect(msg string) T {
	if o.IsNone() {
		panic(msg)
	}

	return o.unwrap()
}

// Returns the contained Some value.
func (o Option[T]) Unwrap() T {
	return o.Expect("called `Option.Unwrap()` on a `None` value")
}

// Returns the contained Some value or a provided default.
func (o Option[T]) UnwrapOr(def T) T {
	if o.Value == nil {
		return def
	}

	return o.unwrap()
}

// Returns the contained Some value or computes it from a closure.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.Value == nil {
		return f()
	}

	return o.unwrap()
}

// Returns the contained Some value or a default.
func (o Option[T]) UnwrapOrDefault() (v T) {
	if o.Value != nil {
		v = o.unwrap()
	}

	return
}

func (o Option[T]) unwrap() T { return *o.Value }
