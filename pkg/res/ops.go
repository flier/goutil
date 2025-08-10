package res

// Maps a Result[T] to Result[U] by applying a function to a contained Ok value, leaving an Err value untouched.
func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}

	return Ok(f(r.unwrap()))
}

// Maps a Result[T] by applying a function to a contained Ok value, leaving an Err value untouched.
func (r Result[T]) Map(f func(T) T) Result[T] { return Map(r, f) }

// Returns the provided default (if Err), or applies a function to the contained value (if Ok).
func MapOr[T any, U any](r Result[T], def U, f func(T) U) U {
	if r.IsErr() {
		return def
	}

	return f(r.unwrap())
}

// Returns the provided default (if Err), or applies a function to the contained value (if Ok).
func (r Result[T]) MapOr(def T, f func(T) T) T { return MapOr(r, def, f) }

// Maps a Result[T] to U by applying fallback function default to a contained Err value,
// or function f to a contained Ok value.
func MapOrElse[T any, U any](r Result[T], def func() U, f func(T) U) U {
	if r.IsErr() {
		return def()
	}

	return f(r.unwrap())
}

// Maps a Result[T] by applying fallback function default to a contained Err value,
// or function f to a contained Ok value.
func (r Result[T]) MapOrElse(def func() T, f func(T) T) T { return MapOrElse(r, def, f) }

// Maps a Result[T] by applying a function to a contained Err value, leaving an Ok value untouched.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.IsErr() {
		return Err[T](f(r.err))
	}

	return Ok(r.unwrap())
}

// Calls a function with a reference to the contained value if Ok.
func (r Result[T]) Inspect(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.unwrap())
	}

	return r
}

// Calls a function with a reference to the contained value if Err.
func (r Result[T]) InspectErr(f func(error)) Result[T] {
	if r.IsErr() {
		f(r.err)
	}

	return r
}

// Returns res2 if the res1 is Ok, otherwise returns the Err value of res1.
func And[T, U any](res1 Result[T], res2 Result[U]) Result[U] {
	if res1.IsErr() {
		return Err[U](res1.err)
	}

	return res2
}

// Returns res if the r is Ok, otherwise returns the Err value of r.
func (r Result[T]) And(res Result[T]) Result[T] { return And(r, res) }

// Calls op if the res is Ok, otherwise returns the Err value of self.
func AndThen[T, U any](res Result[T], op func(T) Result[U]) Result[U] {
	if res.IsErr() {
		return Err[U](res.err)
	}

	return op(res.unwrap())
}

// Calls op if the r is Ok, otherwise returns the Err value of self.
func (r Result[T]) AndThen(op func(T) Result[T]) Result[T] { return AndThen(r, op) }

// Returns res if the res is Err, otherwise returns the Ok value of self.
func (r Result[T]) Or(res Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}

	return res
}

// Calls op if the result is Err, otherwise returns the Ok value of self.
func (r Result[T]) OrElse(op func(error) Result[T]) Result[T] {
	if r.IsOk() {
		return Ok(r.unwrap())
	}

	return op(r.err)
}

// Converts from Result[Result[T]] to Result[T].
func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.IsErr() {
		return Err[T](r.err)
	}

	return r.unwrap()
}
