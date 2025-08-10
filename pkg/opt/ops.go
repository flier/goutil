package opt

import (
	"github.com/flier/goutil/pkg/res"
	"github.com/flier/goutil/pkg/tuple"
)

// Maps an Option<T> by applying a function to a contained value (if Some) or returns None (if None).
func (o Option[T]) Map(f func(T) T) Option[T] { return Map(o, f) }

// Maps an Option<T> to Option<U> by applying a function to a contained value (if Some) or returns None (if None).
func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if o.IsNone() {
		return None[U]()
	}

	return Some(f(o.unwrap()))
}

// Returns the provided default result (if none), or applies a function to the contained value (if any).
func (o Option[T]) MapOr(def T, f func(T) T) T { return MapOr(o, def, f) }

// Returns the provided default result (if none), or applies a function to the contained value (if any).
func MapOr[T, U any](o Option[T], def U, f func(T) U) U {
	if o.IsNone() {
		return def
	}

	return f(o.unwrap())
}

// Computes a default function result (if none), or applies a different function to the contained value (if any).
func (o Option[T]) MapOrElse(def func() T, f func(T) T) T { return MapOrElse(o, def, f) }

// Computes a default function result (if none), or applies a different function to the contained value (if any).
func MapOrElse[T, U any](o Option[T], def func() U, f func(T) U) U {
	if o.IsNone() {
		return def()
	}

	return f(o.unwrap())
}

// Calls a function with a reference to the contained value if Some.
func (o Option[T]) Inspect(f func(T)) Option[T] {
	if o.IsSome() {
		f(o.unwrap())
	}

	return o
}

// Transforms the Option[T] into a Result[T], mapping Some(v) to Ok(v) and None to Err(err).
func (o Option[T]) OkOr(err error) res.Result[T] {
	if o.IsSome() {
		return res.Ok(o.unwrap())
	}

	return res.Err[T](err)
}

// Converts from Result[T] to Option[T].
func Ok[T any](r res.Result[T]) Option[T] {
	if r.IsOk() {
		return Some(r.Unwrap())
	}

	return None[T]()
}

// Converts from Result<T, E> to Option<E>.
func Err[T any](r res.Result[T]) Option[error] {
	if r.IsErr() {
		return Some(r.UnwrapErr())
	}

	return None[error]()
}

// Transforms the Option[T] into a Result[T], mapping Some(v) to Ok(v) and None to Err(err()).
func (o Option[T]) OkOrElse(err func() error) res.Result[T] {
	if o.IsSome() {
		return res.Ok(o.unwrap())
	}

	return res.Err[T](err())
}

// Returns None if the option is None, otherwise returns optb.
func (o Option[T]) And(optb Option[T]) Option[T] { return And(o, optb) }

// Returns None if the option is None, otherwise returns optb.
func And[T, U any](opta Option[T], optb Option[U]) Option[U] {
	if opta.IsNone() {
		return None[U]()
	}

	return optb
}

// Returns None if the option is None, otherwise calls f with the wrapped value and returns the result.
func (o Option[T]) AndThen(f func(T) Option[T]) Option[T] { return AndThen(o, f) }

// Returns None if the option is None, otherwise calls f with the wrapped value and returns the result.
func AndThen[T, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.IsNone() {
		return None[U]()
	}

	return f(o.unwrap())
}

// Returns None if the option is None, otherwise calls f with the wrapped value and returns:
//
//   - Some(t) if f returns true (where t is the wrapped value), and
//   - None if f returns false.
func (o Option[T]) Filter(f func(T) bool) Option[T] {
	if o.IsSome() && f(o.unwrap()) {
		return o
	}

	return None[T]()
}

// Returns the option if it contains a value, otherwise returns optb.
func (o Option[T]) Or(optb Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}

	return optb
}

// Returns the option if it contains a value, otherwise calls f and returns the result.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}

	return f()
}

// Returns Some if exactly one of self, optb is Some, otherwise returns None.
func (o Option[T]) Xor(optb Option[T]) Option[T] {
	if o.IsSome() && optb.IsNone() {
		return o
	}

	if o.IsNone() && optb.IsSome() {
		return optb
	}

	return None[T]()
}

// Converts from Option[Option[T]] to Option[T].
func Flatten[T any](o Option[Option[T]]) Option[T] {
	if o.IsNone() {
		return None[T]()
	}

	return o.unwrap()
}

// Zips x with y Option.
//
// If x is Some(s) and y is Some(o), this method returns Some((s, o)). Otherwise, None is returned.
func Zip[T, U any](x Option[T], y Option[U]) Option[tuple.Tuple2[T, U]] {
	if x.IsSome() && y.IsSome() {
		return Some(tuple.New2(x.unwrap(), y.unwrap()))
	}

	return None[tuple.Tuple2[T, U]]()
}

// Zips x and y Option with function f.
//
// If x is Some(s) and y is Some(o), this method returns Some(f(s, o)). Otherwise, None is returned.
func ZipWith[T, U, R any](x Option[T], y Option[U], f func(T, U) R) Option[R] {
	if x.IsSome() && y.IsSome() {
		return Some(f(x.unwrap(), y.unwrap()))
	}

	return None[R]()
}

// Unzips an option containing a tuple of two options.
//
// If x is Some((a, b)) this method returns (Some(a), Some(b)). Otherwise, (None, None) is returned.
func Unzip[T, U any](x Option[tuple.Tuple2[T, U]]) (Option[T], Option[U]) {
	if x.IsSome() {
		v0, v1 := x.Unwrap().Unpack()

		return Some(v0), Some(v1)
	}

	return None[T](), None[U]()
}
