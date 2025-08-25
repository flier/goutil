//go:build go1.23

package xiter

import "iter"

// Mapper maps a sequence of values of type T to a sequence of values of type O.
type Mapper[T, O any] interface {
	Map(iter.Seq[T]) iter.Seq[O]
}

// KeyMapper maps a sequence of key-value pairs of type [K, V] to a sequence of key-value pairs of type [O, V].
type KeyMapper[K, V, O any] interface {
	MapKey(iter.Seq2[K, V]) iter.Seq2[O, V]
}

// ValueMapper maps a sequence of key-value pairs of type [K, V] to a sequence of key-value pairs of type [K, O].
type ValueMapper[K, V, O any] interface {
	MapValue(iter.Seq2[K, V]) iter.Seq2[K, O]
}

// KeyValueMapper maps a sequence of key-value pairs of type [K, V] to a sequence of key-value pairs of type [O, P].
type KeyValueMapper[K, V, O, P any] interface {
	MapKeyValue(iter.Seq2[K, V]) iter.Seq2[O, P]
}

// Reducer reduces a sequence of values of type T to a single value of type B.
type Reducer[T, B any] interface {
	Reduce(iter.Seq[T]) B
}

// Reducer2 reduces a sequence of key-value pairs of type [K, V] to a single value of type B.
type Reducer2[K, V, B any] interface {
	Reduce2(iter.Seq2[K, V]) B
}

// Comparer compares two sequence of values of type T and returns an integer.
type Comparer[T any] interface {
	Compare(iter.Seq[T], iter.Seq[T]) int
}

// MappingFunc is a functor that maps a sequence of values of type T to a sequence of values of type O.
type MappingFunc[T, O any] func(iter.Seq[T]) iter.Seq[O]

var _ Mapper[int, int] = MappingFunc[int, int](nil)

func (f MappingFunc[T, O]) Map(s iter.Seq[T]) iter.Seq[O] { return f(s) }

// MappingKeyValueFunc is a functor that maps a sequence of key-value pairs of type [K, V] to a sequence of key-value pairs of type [O, P].
type MappingKeyValueFunc[K, V, O, P any] func(iter.Seq2[K, V]) iter.Seq2[O, P]

var _ KeyValueMapper[string, int, bool, rune] = MappingKeyValueFunc[string, int, bool, rune](nil)

func (f MappingKeyValueFunc[K, V, O, P]) MapKeyValue(s iter.Seq2[K, V]) iter.Seq2[O, P] { return f(s) }

// MappingKeyFunc is a functor that maps a sequence of key-value pairs of type [K, V] to a sequence of key-value pairs of type [O, V].
type MappingKeyFunc[K, V, O any] func(iter.Seq2[K, V]) iter.Seq2[O, V]

var _ KeyMapper[string, int, bool] = MappingKeyFunc[string, int, bool](nil)

func (f MappingKeyFunc[K, V, O]) MapKey(s iter.Seq2[K, V]) iter.Seq2[O, V] { return f(s) }

// MappingValueFunc is a functor that maps a sequence of key-value pairs of type [K, V] to a sequence of key-value pairs of type [K, O].
type MappingValueFunc[K, V, O any] func(iter.Seq2[K, V]) iter.Seq2[K, O]

var _ ValueMapper[int, int, int] = MappingValueFunc[int, int, int](nil)

func (f MappingValueFunc[K, V, O]) MapValue(s iter.Seq2[K, V]) iter.Seq2[K, O] { return f(s) }

// ReductionFunc is a functor that reduces a sequence of values of type T to a single value of type B.
type ReductionFunc[T, B any] func(iter.Seq[T]) B

var _ Reducer[int, int] = ReductionFunc[int, int](nil)

func (f ReductionFunc[T, B]) Reduce(s iter.Seq[T]) B { return f(s) }

// Reduction2Func is a functor that reduces a sequence of key-value pairs of type [K, V] to a single value of type B.
type Reduction2Func[K, V, B any] func(iter.Seq2[K, V]) B

var _ Reducer2[int, int, int] = Reduction2Func[int, int, int](nil)

func (f Reduction2Func[K, V, B]) Reduce2(s iter.Seq2[K, V]) B { return f(s) }

// CompareFunc is a functor that compare two sequence of values of type T and returns
type CompareFunc[T any] func(iter.Seq[T], iter.Seq[T]) int

var _ Comparer[int] = CompareFunc[int](nil)

func (f CompareFunc[T]) Compare(x iter.Seq[T], y iter.Seq[T]) int { return f(x, y) }

func bind2[T1, T2, O any](f func(T1, T2) O, arg T2) func(T1) O {
	return func(x T1) O {
		return f(x, arg)
	}
}

func bind2rest[T, A, O any](f func(T, ...A) O, args []A) func(v T) O {
	return func(v T) O {
		return f(v, args...)
	}
}

func bind23[T1, T2, T3, O any](f func(T1, T2, T3) O, v1 T2, v2 T3) func(T1) O {
	return func(v0 T1) O {
		return f(v0, v1, v2)
	}
}

func bind3[T1, T2, T3, O any](f func(T1, T2, T3) O, arg T3) func(T1, T2) O {
	return func(v0 T1, v1 T2) O {
		return f(v0, v1, arg)
	}
}
