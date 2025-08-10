// Package either provides a generic [Either] type implementation that represents
// a value of one of two possible types (a disjoint union).
//
// [Either] is commonly used for handling success/failure cases where the Left
// variant typically represents failure or error cases, and the Right variant
// represents success cases. However, it can be used for any scenario requiring
// a choice between two types.
//
// The package provides:
//   - Type-safe construction via [Left] and [Right] functions
//   - Pattern matching via [Either.HasLeft] and [Either.HasRight]
//   - Value extraction via [Either.UnwrapLeft], [Either.UnwrapRight], and various [Either.LeftOr]/[Either.LeftOrElse]/[Either.RightOr]/[Either.RightOrElse] methods
//   - Transformation functions like [MapLeft], [MapRight], and [MapEither]
//   - Monadic operations via [LeftAndThen] and [RightAndThen]
//   - Utility functions like [Either.Flip] and [Reduce]
//
// Example usage:
//
//	result := either.Right[string, int](42)
//	mapped := either.MapRight(result, func(i int) string {
//	    return fmt.Sprintf("value: %d", i)
//	})
package either

import "fmt"

// Either with variants Left and Right is a general purpose sum type with two cases.
type Either[L, R any] struct {
	left  *L // A value of type L.
	right *R // A value of type R.
}

func Empty[L, R any]() Either[L, R] {
	return Either[L, R]{}
}

// Left creates a new Either value with the given left value.
func Left[L, R any](left L) Either[L, R] {
	return Either[L, R]{left: &left}
}

// Right creates a new Either value with the given right value.
func Right[L, R any](right R) Either[L, R] {
	return Either[L, R]{right: &right}
}

func (e Either[L, R]) String() string {
	if e.left != nil {
		return fmt.Sprintf("Left(%v)", *e.left)
	}

	if e.right != nil {
		return fmt.Sprintf("Right(%v)", *e.right)
	}

	return "Empty"
}

func (e Either[L, R]) GoString() string {
	if e.left != nil {
		return fmt.Sprintf("Either { Left: %v }", *e.left)
	}

	if e.right != nil {
		return fmt.Sprintf("Either { Right: %v }", *e.right)
	}

	return "Either {}"
}

// HasLeft returns true if the value is the Left variant.
func (e Either[L, R]) HasLeft() bool { return e.left != nil }

// HasRight returns true if the value is the Right variant.
func (e Either[L, R]) HasRight() bool { return e.right != nil }

// Flip converts Either[L, R] to Either[R, L].
func (e Either[L, R]) Flip() Either[R, L] { return Either[R, L]{e.right, e.left} }

// LeftOr returns left value or given value
func (e Either[L, R]) LeftOr(other L) L {
	if e.left != nil {
		return *e.left
	}

	return other
}

// LeftOrEmpty returns left or a empty value
func (e Either[L, R]) LeftOrEmpty() (l L) {
	if e.left != nil {
		l = *e.left
	}

	return
}

// LeftOrElse returns left value or computes it from a function f
func (e Either[L, R]) LeftOrElse(f func() L) L {
	if e.left != nil {
		return *e.left
	}

	return f()
}

// RightOr returns right value or given value
func (e Either[L, R]) RightOr(other R) R {
	if e.right != nil {
		return *e.right
	}

	return other
}

// RightOrEmpty returns right or a empty value
func (e Either[L, R]) RightOrEmpty() (r R) {
	if e.right != nil {
		r = *e.right
	}

	return
}

// RightOrElse returns right value or computes it from a function f
func (e Either[L, R]) RightOrElse(f func() R) R {
	if e.right != nil {
		return *e.right
	}

	return f()
}

// UnwrapLeft returns the left value or panic
func (e Either[L, R]) UnwrapLeft() L {
	if e.left == nil {
		unwrapFail("called Either.UnwrapLeft on a Right value: %v", e.RightOrEmpty())
	}

	return *e.left
}

// UnwrapRight returns the right value or panic
func (e Either[L, R]) UnwrapRight() R {
	if e.right == nil {
		unwrapFail("called Either.UnwrapRight on a Left value: %v", e.LeftOrEmpty())
	}

	return *e.right
}

// ExpectLeft returns the left value or panic with message
func (e Either[L, R]) ExpectLeft(msg string) L {
	if e.left == nil {
		unwrapFail("%s: %v", msg, e.RightOrEmpty())
	}

	return *e.left
}

// ExpectRight returns the right value or panic with message
func (e Either[L, R]) ExpectRight(msg string) R {
	if e.right == nil {
		unwrapFail("%s: %v", msg, e.LeftOrEmpty())
	}

	return *e.right
}

func unwrapFail(format string, a ...any) { panic(fmt.Sprintf(format, a...)) }

// MapLeft applies the function f on the value in the Left variant if it is present rewrapping the result in Left.
func MapLeft[L, R, M any](e Either[L, R], f func(L) M) Either[M, R] {
	if e.left == nil {
		return Either[M, R]{nil, e.right}
	}

	m := f(*e.left)

	return Either[M, R]{&m, e.right}
}

// MapRight applies the function f on the value in the Right variant if it is present rewrapping the result in Right.
func MapRight[L, R, M any](e Either[L, R], f func(R) M) Either[L, M] {
	if e.right == nil {
		return Either[L, M]{left: e.left, right: nil}
	}

	m := f(*e.right)

	return Either[L, M]{e.left, &m}
}

// MapEither applies the functions f and g to the Left and Right variants respectively.
func MapEither[L, R, M, S any](e Either[L, R], f func(L) M, g func(R) S) Either[M, S] {
	var m Either[M, S]

	if e.left != nil {
		l := f(*e.left)
		m.left = &l
	}

	if e.right != nil {
		r := g(*e.right)
		m.right = &r
	}

	return m
}

// Reduce applies one of two functions depending on contents, unifying their result.
func Reduce[L, R, T any](e Either[L, R], f func(L) T, g func(R) T) (t T) {
	if e.left != nil {
		t = f(*e.left)

	}

	if e.right != nil {
		t = g(*e.right)
	}

	return
}

// LeftAndThen applies the function f on the value in the Left variant if it is present.
func LeftAndThen[L, R, S any](e Either[L, R], f func(L) Either[S, R]) Either[S, R] {
	if e.left == nil {
		return Either[S, R]{nil, e.right}
	}

	return f(*e.left)
}

// RightAndThen applies the function f on the value in the Right variant if it is present.
func RightAndThen[L, R, S any](e Either[L, R], f func(R) Either[L, S]) Either[L, S] {
	if e.right == nil {
		return Either[L, S]{left: e.left, right: nil}
	}

	return f(*e.right)
}
