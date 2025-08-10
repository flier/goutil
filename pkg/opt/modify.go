package opt

// Inserts value into the option, then returns a pointer to it.
func (o *Option[T]) Insert(value T) *T {
	o.val = &value

	return o.val
}

// Inserts value into the option if it is None, then returns a pointer to the contained value.
func (o *Option[T]) GetOrInsert(value T) *T {
	if o.IsNone() {
		o.val = &value
	}

	return o.val
}

// Inserts the default value into the option if it is None,
// then returns a mutable reference to the contained value.
func (o *Option[T]) GetOrInsertDefault() *T {
	if o.IsNone() {
		o.val = new(T)
	}

	return o.val
}

// Inserts a value computed from f into the option if it is None,
// then returns a mutable reference to the contained value.
func (o *Option[T]) GetOrInsertWith(f func() T) *T {
	if o.IsNone() {
		v := f()

		o.val = &v
	}

	return o.val
}

// Takes the value out of the option, leaving a None in its place.
func (o *Option[T]) Take() Option[T] {
	opt := Option[T]{o.val}

	o.val = nil

	return opt
}

// Takes the value out of the option, but only if the predicate evaluates to true on a mutable reference to the value.
func (o *Option[T]) TakeIf(f func(T) bool) Option[T] {
	if o.IsSome() && f(*o.val) {
		return o.Take()
	}

	return None[T]()
}

// Replaces the actual value in the option by the value given in parameter, returning the old value if present,
// leaving a Some in its place without deinitializing either one.
func (o *Option[T]) Replace(value T) Option[T] {
	opt := Option[T]{o.val}

	o.val = &value

	return opt
}
