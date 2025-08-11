// A finite heterogeneous sequence, (T0, T1, ..).
package tuple

import (
	"errors"
	"fmt"
)

var ErrOutOfRange = errors.New("out of range")

// Tuple is a generic tuple.
type Tuple interface {
	fmt.Stringer

	// Len returns the number of elements in the tuple.
	Len() int

	// Get returns the element at the given index.
	Get(i int) any

	// Put sets the element at the given index to the given value and returns the new tuple and the old value.
	Put(i int, v any) (new Tuple, old any)

	// Del removes the element at the given index and returns the new tuple.
	Del(i int) Tuple
}

var (
	_ Tuple = Tuple0{}
	_ Tuple = Tuple1[int]{}
	_ Tuple = Tuple2[int, int]{}
	_ Tuple = Tuple3[int, int, int]{}
	_ Tuple = Tuple4[int, int, int, int]{}
	_ Tuple = Tuple5[int, int, int, int, int]{}
	_ Tuple = Tuple6[int, int, int, int, int, int]{}
	_ Tuple = Tuple7[int, int, int, int, int, int, int]{}
)

// Tuple0 is a tuple with 0 elements.
type Tuple0 struct{}

func New0() Tuple0                             { return Tuple0{} }
func (t Tuple0) String() string                { return "()" }
func (t Tuple0) Len() int                      { return 0 }
func (t Tuple0) Get(i int) any                 { panic(indexOutOfRangeError(i, t)) }
func (t Tuple0) Put(i int, v any) (Tuple, any) { panic(indexOutOfRangeError(i, t)) }
func (t Tuple0) Del(i int) Tuple               { panic(indexOutOfRangeError(i, t)) }

func indexOutOfRangeError(i int, t Tuple) error {
	return fmt.Errorf("index %d with length %d, %w", i, t.Len(), ErrOutOfRange)
}

type Tuple1[T0 any] struct {
	V0 T0
}

func New1[T0 any](v0 T0) Tuple1[T0] { return Tuple1[T0]{v0} }

func (t Tuple1[T0]) Unpack() T0         { return t.V0 }
func (t Tuple1[T0]) Head() (T0, Tuple0) { return t.V0, Tuple0{} }
func (t Tuple1[T0]) Tail() (Tuple0, T0) { return Tuple0{}, t.V0 }
func (t Tuple1[T0]) String() string     { return fmt.Sprintf("(%v)", t.V0) }
func (t Tuple1[T0]) Len() int           { return 1 }
func (t Tuple1[T0]) Get(i int) any {
	if i == 0 {
		return t.V0
	}

	panic(indexOutOfRangeError(i, t))
}

func (t Tuple1[T0]) Put(i int, v any) (new Tuple, old any) {
	if i == 0 {
		return t.Put0(v.(T0))
	}

	panic(indexOutOfRangeError(i, t))
}

func (t Tuple1[T0]) Put0(v T0) (new Tuple1[T0], old T0) {
	return Tuple1[T0]{v}, t.V0
}

func (t Tuple1[T0]) Del(i int) Tuple {
	if i == 0 {
		return t.Del0()
	}

	panic(indexOutOfRangeError(i, t))
}

func (t Tuple1[T0]) Del0() Tuple0 { return Tuple0{} }

type Tuple2[T0, T1 any] struct {
	V0 T0
	V1 T1
}

func New2[T0, T1 any](v0 T0, v1 T1) Tuple2[T0, T1] {
	return Tuple2[T0, T1]{v0, v1}
}

func (t Tuple2[T0, T1]) Unpack() (T0, T1)       { return t.V0, t.V1 }
func (t Tuple2[T0, T1]) Head() (T0, Tuple1[T1]) { return t.V0, Tuple1[T1]{t.V1} }
func (t Tuple2[T0, T1]) Tail() (Tuple1[T0], T1) { return Tuple1[T0]{t.V0}, t.V1 }
func (t Tuple2[T0, T1]) String() string         { return fmt.Sprintf("(%v, %v)", t.V0, t.V1) }

func (t Tuple2[T0, T1]) Len() int { return 2 }

func (t Tuple2[T0, T1]) Get(i int) any {
	switch i {
	case 0:
		return t.V0
	case 1:
		return t.V1
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple2[T0, T1]) Put(i int, v any) (new Tuple, old any) {
	switch i {
	case 0:
		return t.Put0(v.(T0))
	case 1:
		return t.Put1(v.(T1))
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple2[T0, T1]) Put0(v T0) (new Tuple2[T0, T1], old T0) {
	return Tuple2[T0, T1]{v, t.V1}, t.V0
}
func (t Tuple2[T0, T1]) Put1(v T1) (new Tuple2[T0, T1], old T1) {
	return Tuple2[T0, T1]{t.V0, v}, t.V1
}

func (t Tuple2[T0, T1]) Del(i int) Tuple {
	switch i {
	case 0:
		return t.Del0()
	case 1:
		return t.Del1()
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple2[T0, T1]) Del0() Tuple1[T1] { return Tuple1[T1]{t.V1} }
func (t Tuple2[T0, T1]) Del1() Tuple1[T0] { return Tuple1[T0]{t.V0} }

type Tuple3[T0, T1, T2 any] struct {
	V0 T0
	V1 T1
	V2 T2
}

func New3[T0, T1, T2 any](v0 T0, v1 T1, v2 T2) Tuple3[T0, T1, T2] {
	return Tuple3[T0, T1, T2]{v0, v1, v2}
}

func (t Tuple3[T0, T1, T2]) Unpack() (T0, T1, T2)       { return t.V0, t.V1, t.V2 }
func (t Tuple3[T0, T1, T2]) Head() (T0, Tuple2[T1, T2]) { return t.V0, Tuple2[T1, T2]{t.V1, t.V2} }
func (t Tuple3[T0, T1, T2]) Tail() (Tuple2[T0, T1], T2) { return Tuple2[T0, T1]{t.V0, t.V1}, t.V2 }
func (t Tuple3[T0, T1, T2]) String() string             { return fmt.Sprintf("(%v, %v, %v)", t.V0, t.V1, t.V2) }

func (t Tuple3[T0, T1, T2]) Len() int { return 3 }

func (t Tuple3[T0, T1, T2]) Get(i int) any {
	switch i {
	case 0:
		return t.V0
	case 1:
		return t.V1
	case 2:
		return t.V2
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple3[T0, T1, T2]) Put(i int, v any) (new Tuple, old any) {
	switch i {
	case 0:
		return t.Put0(v.(T0))
	case 1:
		return t.Put1(v.(T1))
	case 2:
		return t.Put2(v.(T2))
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple3[T0, T1, T2]) Put0(v T0) (new Tuple3[T0, T1, T2], old T0) {
	return Tuple3[T0, T1, T2]{v, t.V1, t.V2}, t.V0
}
func (t Tuple3[T0, T1, T2]) Put1(v T1) (new Tuple3[T0, T1, T2], old T1) {
	return Tuple3[T0, T1, T2]{t.V0, v, t.V2}, t.V1
}
func (t Tuple3[T0, T1, T2]) Put2(v T2) (new Tuple3[T0, T1, T2], old T2) {
	return Tuple3[T0, T1, T2]{t.V0, t.V1, v}, t.V2
}

func (t Tuple3[T0, T1, T2]) Del(i int) Tuple {
	switch i {
	case 0:
		return t.Del0()
	case 1:
		return t.Del1()
	case 2:
		return t.Del2()
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple3[T0, T1, T2]) Del0() Tuple2[T1, T2] { return Tuple2[T1, T2]{t.V1, t.V2} }
func (t Tuple3[T0, T1, T2]) Del1() Tuple2[T0, T2] { return Tuple2[T0, T2]{t.V0, t.V2} }
func (t Tuple3[T0, T1, T2]) Del2() Tuple2[T0, T1] { return Tuple2[T0, T1]{t.V0, t.V1} }

type Tuple4[T0, T1, T2, T3 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
}

func New4[T0, T1, T2, T3 any](v0 T0, v1 T1, v2 T2, v3 T3) Tuple4[T0, T1, T2, T3] {
	return Tuple4[T0, T1, T2, T3]{v0, v1, v2, v3}
}

func (t Tuple4[T0, T1, T2, T3]) Unpack() (T0, T1, T2, T3) { return t.V0, t.V1, t.V2, t.V3 }

func (t Tuple4[T0, T1, T2, T3]) Head() (T0, Tuple3[T1, T2, T3]) {
	return t.V0, Tuple3[T1, T2, T3]{t.V1, t.V2, t.V3}
}
func (t Tuple4[T0, T1, T2, T3]) Tail() (Tuple3[T0, T1, T2], T3) {
	return Tuple3[T0, T1, T2]{t.V0, t.V1, t.V2}, t.V3
}

func (t Tuple4[T0, T1, T2, T3]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v)", t.V0, t.V1, t.V2, t.V3)
}

func (t Tuple4[T0, T1, T2, T3]) Len() int { return 4 }

func (t Tuple4[T0, T1, T2, T3]) Get(i int) any {
	switch i {
	case 0:
		return t.V0
	case 1:
		return t.V1
	case 2:
		return t.V2
	case 3:
		return t.V3
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple4[T0, T1, T2, T3]) Put(i int, v any) (new Tuple, old any) {
	switch i {
	case 0:
		return t.Put0(v.(T0))
	case 1:
		return t.Put1(v.(T1))
	case 2:
		return t.Put2(v.(T2))
	case 3:
		return t.Put3(v.(T3))
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple4[T0, T1, T2, T3]) Put0(v T0) (new Tuple4[T0, T1, T2, T3], old T0) {
	return Tuple4[T0, T1, T2, T3]{v, t.V1, t.V2, t.V3}, t.V0
}
func (t Tuple4[T0, T1, T2, T3]) Put1(v T1) (new Tuple4[T0, T1, T2, T3], old T1) {
	return Tuple4[T0, T1, T2, T3]{t.V0, v, t.V2, t.V3}, t.V1
}
func (t Tuple4[T0, T1, T2, T3]) Put2(v T2) (new Tuple4[T0, T1, T2, T3], old T2) {
	return Tuple4[T0, T1, T2, T3]{t.V0, t.V1, v, t.V3}, t.V2
}
func (t Tuple4[T0, T1, T2, T3]) Put3(v T3) (new Tuple4[T0, T1, T2, T3], old T3) {
	return Tuple4[T0, T1, T2, T3]{t.V0, t.V1, t.V2, v}, t.V3
}

func (t Tuple4[T0, T1, T2, T3]) Del(i int) Tuple {
	switch i {
	case 0:
		return t.Del0()
	case 1:
		return t.Del1()
	case 2:
		return t.Del2()
	case 3:
		return t.Del3()
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple4[T0, T1, T2, T3]) Del0() Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3]{t.V1, t.V2, t.V3}
}
func (t Tuple4[T0, T1, T2, T3]) Del1() Tuple3[T0, T2, T3] {
	return Tuple3[T0, T2, T3]{t.V0, t.V2, t.V3}
}
func (t Tuple4[T0, T1, T2, T3]) Del2() Tuple3[T0, T1, T3] {
	return Tuple3[T0, T1, T3]{t.V0, t.V1, t.V3}
}
func (t Tuple4[T0, T1, T2, T3]) Del3() Tuple3[T0, T1, T2] {
	return Tuple3[T0, T1, T2]{t.V0, t.V1, t.V2}
}

type Tuple5[T0, T1, T2, T3, T4 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
}

func New5[T0, T1, T2, T3, T4 any](v0 T0, v1 T1, v2 T2, v3 T3, v4 T4) Tuple5[T0, T1, T2, T3, T4] {
	return Tuple5[T0, T1, T2, T3, T4]{v0, v1, v2, v3, v4}
}

func (t Tuple5[T0, T1, T2, T3, T4]) Unpack() (T0, T1, T2, T3, T4) {
	return t.V0, t.V1, t.V2, t.V3, t.V4
}

func (t Tuple5[T0, T1, T2, T3, T4]) Head() (T0, Tuple4[T1, T2, T3, T4]) {
	return t.V0, Tuple4[T1, T2, T3, T4]{t.V1, t.V2, t.V3, t.V4}
}
func (t Tuple5[T0, T1, T2, T3, T4]) Tail() (Tuple4[T0, T1, T2, T3], T4) {
	return Tuple4[T0, T1, T2, T3]{t.V0, t.V1, t.V2, t.V3}, t.V4
}

func (t Tuple5[T0, T1, T2, T3, T4]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v, %v)", t.V0, t.V1, t.V2, t.V3, t.V4)
}

func (t Tuple5[T0, T1, T2, T3, T4]) Len() int { return 5 }

func (t Tuple5[T0, T1, T2, T3, T4]) Get(i int) any {
	switch i {
	case 0:
		return t.V0
	case 1:
		return t.V1
	case 2:
		return t.V2
	case 3:
		return t.V3
	case 4:
		return t.V4
	default:
		panic(indexOutOfRangeError(i, t))
	}
}
func (t Tuple5[T0, T1, T2, T3, T4]) Put(i int, v any) (new Tuple, old any) {
	switch i {
	case 0:
		return t.Put0(v.(T0))
	case 1:
		return t.Put1(v.(T1))
	case 2:
		return t.Put2(v.(T2))
	case 3:
		return t.Put3(v.(T3))
	case 4:
		return t.Put4(v.(T4))
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple5[T0, T1, T2, T3, T4]) Put0(v T0) (new Tuple5[T0, T1, T2, T3, T4], old T0) {
	return Tuple5[T0, T1, T2, T3, T4]{v, t.V1, t.V2, t.V3, t.V4}, t.V0
}
func (t Tuple5[T0, T1, T2, T3, T4]) Put1(v T1) (new Tuple5[T0, T1, T2, T3, T4], old T1) {
	return Tuple5[T0, T1, T2, T3, T4]{t.V0, v, t.V2, t.V3, t.V4}, t.V1
}
func (t Tuple5[T0, T1, T2, T3, T4]) Put2(v T2) (new Tuple5[T0, T1, T2, T3, T4], old T2) {
	return Tuple5[T0, T1, T2, T3, T4]{t.V0, t.V1, v, t.V3, t.V4}, t.V2
}
func (t Tuple5[T0, T1, T2, T3, T4]) Put3(v T3) (new Tuple5[T0, T1, T2, T3, T4], old T3) {
	return Tuple5[T0, T1, T2, T3, T4]{t.V0, t.V1, t.V2, v, t.V4}, t.V3
}
func (t Tuple5[T0, T1, T2, T3, T4]) Put4(v T4) (new Tuple5[T0, T1, T2, T3, T4], old T4) {
	return Tuple5[T0, T1, T2, T3, T4]{t.V0, t.V1, t.V2, t.V3, v}, t.V4
}

func (t Tuple5[T0, T1, T2, T3, T4]) Del(i int) Tuple {
	switch i {
	case 0:
		return t.Del0()
	case 1:
		return t.Del1()
	case 2:
		return t.Del2()
	case 3:
		return t.Del3()
	case 4:
		return t.Del4()
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple5[T0, T1, T2, T3, T4]) Del0() Tuple4[T1, T2, T3, T4] {
	return Tuple4[T1, T2, T3, T4]{t.V1, t.V2, t.V3, t.V4}
}
func (t Tuple5[T0, T1, T2, T3, T4]) Del1() Tuple4[T0, T2, T3, T4] {
	return Tuple4[T0, T2, T3, T4]{t.V0, t.V2, t.V3, t.V4}
}
func (t Tuple5[T0, T1, T2, T3, T4]) Del2() Tuple4[T0, T1, T3, T4] {
	return Tuple4[T0, T1, T3, T4]{t.V0, t.V1, t.V3, t.V4}
}
func (t Tuple5[T0, T1, T2, T3, T4]) Del3() Tuple4[T0, T1, T2, T4] {
	return Tuple4[T0, T1, T2, T4]{t.V0, t.V1, t.V2, t.V4}
}
func (t Tuple5[T0, T1, T2, T3, T4]) Del4() Tuple4[T0, T1, T2, T3] {
	return Tuple4[T0, T1, T2, T3]{t.V0, t.V1, t.V2, t.V3}
}

type Tuple6[T0, T1, T2, T3, T4, T5 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
}

func New6[T0, T1, T2, T3, T4, T5 any](v0 T0, v1 T1, v2 T2, v3 T3, v4 T4, v5 T5) Tuple6[T0, T1, T2, T3, T4, T5] {
	return Tuple6[T0, T1, T2, T3, T4, T5]{v0, v1, v2, v3, v4, v5}
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Unpack() (T0, T1, T2, T3, T4, T5) {
	return t.V0, t.V1, t.V2, t.V3, t.V4, t.V5
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Head() (T0, Tuple5[T1, T2, T3, T4, T5]) {
	return t.V0, Tuple5[T1, T2, T3, T4, T5]{t.V1, t.V2, t.V3, t.V4, t.V5}
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Tail() (Tuple5[T0, T1, T2, T3, T4], T5) {
	return Tuple5[T0, T1, T2, T3, T4]{t.V0, t.V1, t.V2, t.V3, t.V4}, t.V5
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v, %v, %v)", t.V0, t.V1, t.V2, t.V3, t.V4, t.V5)
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Len() int { return 6 }

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Get(i int) any {
	switch i {
	case 0:
		return t.V0
	case 1:
		return t.V1
	case 2:
		return t.V2
	case 3:
		return t.V3
	case 4:
		return t.V4
	case 5:
		return t.V5
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Put(i int, v any) (new Tuple, old any) {
	switch i {
	case 0:
		return t.Put0(v.(T0))
	case 1:
		return t.Put1(v.(T1))
	case 2:
		return t.Put2(v.(T2))
	case 3:
		return t.Put3(v.(T3))
	case 4:
		return t.Put4(v.(T4))
	case 5:
		return t.Put5(v.(T5))
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Put0(v T0) (new Tuple6[T0, T1, T2, T3, T4, T5], old T0) {
	return Tuple6[T0, T1, T2, T3, T4, T5]{v, t.V1, t.V2, t.V3, t.V4, t.V5}, t.V0
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Put1(v T1) (new Tuple6[T0, T1, T2, T3, T4, T5], old T1) {
	return Tuple6[T0, T1, T2, T3, T4, T5]{t.V0, v, t.V2, t.V3, t.V4, t.V5}, t.V1
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Put2(v T2) (new Tuple6[T0, T1, T2, T3, T4, T5], old T2) {
	return Tuple6[T0, T1, T2, T3, T4, T5]{t.V0, t.V1, v, t.V3, t.V4, t.V5}, t.V2
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Put3(v T3) (new Tuple6[T0, T1, T2, T3, T4, T5], old T3) {
	return Tuple6[T0, T1, T2, T3, T4, T5]{t.V0, t.V1, t.V2, v, t.V4, t.V5}, t.V3
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Put4(v T4) (new Tuple6[T0, T1, T2, T3, T4, T5], old T4) {
	return Tuple6[T0, T1, T2, T3, T4, T5]{t.V0, t.V1, t.V2, t.V3, v, t.V5}, t.V4
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Put5(v T5) (new Tuple6[T0, T1, T2, T3, T4, T5], old T5) {
	return Tuple6[T0, T1, T2, T3, T4, T5]{t.V0, t.V1, t.V2, t.V3, t.V4, v}, t.V5
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Del(i int) Tuple {
	switch i {
	case 0:
		return t.Del0()
	case 1:
		return t.Del1()
	case 2:
		return t.Del2()
	case 3:
		return t.Del3()
	case 4:
		return t.Del4()
	case 5:
		return t.Del5()
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple6[T0, T1, T2, T3, T4, T5]) Del0() Tuple5[T1, T2, T3, T4, T5] {
	return Tuple5[T1, T2, T3, T4, T5]{t.V1, t.V2, t.V3, t.V4, t.V5}
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Del1() Tuple5[T0, T2, T3, T4, T5] {
	return Tuple5[T0, T2, T3, T4, T5]{t.V0, t.V2, t.V3, t.V4, t.V5}
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Del2() Tuple5[T0, T1, T3, T4, T5] {
	return Tuple5[T0, T1, T3, T4, T5]{t.V0, t.V1, t.V3, t.V4, t.V5}
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Del3() Tuple5[T0, T1, T2, T4, T5] {
	return Tuple5[T0, T1, T2, T4, T5]{t.V0, t.V1, t.V2, t.V4, t.V5}
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Del4() Tuple5[T0, T1, T2, T3, T5] {
	return Tuple5[T0, T1, T2, T3, T5]{t.V0, t.V1, t.V2, t.V3, t.V5}
}
func (t Tuple6[T0, T1, T2, T3, T4, T5]) Del5() Tuple5[T0, T1, T2, T3, T4] {
	return Tuple5[T0, T1, T2, T3, T4]{t.V0, t.V1, t.V2, t.V3, t.V4}
}

type Tuple7[T0, T1, T2, T3, T4, T5, T6 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
	V4 T4
	V5 T5
	V6 T6
}

func New7[T0, T1, T2, T3, T4, T5, T6 any](v0 T0, v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6) Tuple7[T0, T1, T2, T3, T4, T5, T6] {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{v0, v1, v2, v3, v4, v5, v6}
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Unpack() (T0, T1, T2, T3, T4, T5, T6) {
	return t.V0, t.V1, t.V2, t.V3, t.V4, t.V5, t.V6
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Head() (T0, Tuple6[T1, T2, T3, T4, T5, T6]) {
	return t.V0, Tuple6[T1, T2, T3, T4, T5, T6]{t.V1, t.V2, t.V3, t.V4, t.V5, t.V6}
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Tail() (Tuple6[T0, T1, T2, T3, T4, T5], T6) {
	return Tuple6[T0, T1, T2, T3, T4, T5]{t.V0, t.V1, t.V2, t.V3, t.V4, t.V5}, t.V6
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) String() string {
	return fmt.Sprintf("(%v, %v, %v, %v, %v, %v, %v)", t.V0, t.V1, t.V2, t.V3, t.V4, t.V5, t.V6)
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Len() int { return 7 }

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Get(i int) any {
	switch i {
	case 0:
		return t.V0
	case 1:
		return t.V1
	case 2:
		return t.V2
	case 3:
		return t.V3
	case 4:
		return t.V4
	case 5:
		return t.V5
	case 6:
		return t.V6
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put(i int, v any) (new Tuple, old any) {
	switch i {
	case 0:
		return t.Put0(v.(T0))
	case 1:
		return t.Put1(v.(T1))
	case 2:
		return t.Put2(v.(T2))
	case 3:
		return t.Put3(v.(T3))
	case 4:
		return t.Put4(v.(T4))
	case 5:
		return t.Put5(v.(T5))
	case 6:
		return t.Put6(v.(T6))
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put0(v T0) (new Tuple7[T0, T1, T2, T3, T4, T5, T6], old T0) {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{v, t.V1, t.V2, t.V3, t.V4, t.V5, t.V6}, t.V0
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put1(v T1) (new Tuple7[T0, T1, T2, T3, T4, T5, T6], old T1) {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{t.V0, v, t.V2, t.V3, t.V4, t.V5, t.V6}, t.V1
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put2(v T2) (new Tuple7[T0, T1, T2, T3, T4, T5, T6], old T2) {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{t.V0, t.V1, v, t.V3, t.V4, t.V5, t.V6}, t.V2
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put3(v T3) (new Tuple7[T0, T1, T2, T3, T4, T5, T6], old T3) {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{t.V0, t.V1, t.V2, v, t.V4, t.V5, t.V6}, t.V3
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put4(v T4) (new Tuple7[T0, T1, T2, T3, T4, T5, T6], old T4) {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{t.V0, t.V1, t.V2, t.V3, v, t.V5, t.V6}, t.V4
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put5(v T5) (new Tuple7[T0, T1, T2, T3, T4, T5, T6], old T5) {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{t.V0, t.V1, t.V2, t.V3, t.V4, v, t.V6}, t.V5
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Put6(v T6) (new Tuple7[T0, T1, T2, T3, T4, T5, T6], old T6) {
	return Tuple7[T0, T1, T2, T3, T4, T5, T6]{t.V0, t.V1, t.V2, t.V3, t.V4, t.V5, v}, t.V6
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del(i int) Tuple {
	switch i {
	case 0:
		return t.Del0()
	case 1:
		return t.Del1()
	case 2:
		return t.Del2()
	case 3:
		return t.Del3()
	case 4:
		return t.Del4()
	case 5:
		return t.Del5()
	case 6:
		return t.Del6()
	default:
		panic(indexOutOfRangeError(i, t))
	}
}

func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del0() Tuple6[T1, T2, T3, T4, T5, T6] {
	return Tuple6[T1, T2, T3, T4, T5, T6]{t.V1, t.V2, t.V3, t.V4, t.V5, t.V6}
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del1() Tuple6[T0, T2, T3, T4, T5, T6] {
	return Tuple6[T0, T2, T3, T4, T5, T6]{t.V0, t.V2, t.V3, t.V4, t.V5, t.V6}
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del2() Tuple6[T0, T1, T3, T4, T5, T6] {
	return Tuple6[T0, T1, T3, T4, T5, T6]{t.V0, t.V1, t.V3, t.V4, t.V5, t.V6}
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del3() Tuple6[T0, T1, T2, T4, T5, T6] {
	return Tuple6[T0, T1, T2, T4, T5, T6]{t.V0, t.V1, t.V2, t.V4, t.V5, t.V6}
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del4() Tuple6[T0, T1, T2, T3, T5, T6] {
	return Tuple6[T0, T1, T2, T3, T5, T6]{t.V0, t.V1, t.V2, t.V3, t.V5, t.V6}
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del5() Tuple6[T0, T1, T2, T3, T4, T6] {
	return Tuple6[T0, T1, T2, T3, T4, T6]{t.V0, t.V1, t.V2, t.V3, t.V4, t.V6}
}
func (t Tuple7[T0, T1, T2, T3, T4, T5, T6]) Del6() Tuple6[T0, T1, T2, T3, T4, T5] {
	return Tuple6[T0, T1, T2, T3, T4, T5]{t.V0, t.V1, t.V2, t.V3, t.V4, t.V5}
}
