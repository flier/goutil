// A finite heterogeneous sequence, (T0, T1, ..).
package tuple

import "fmt"

type Tuple2[T0, T1 any] struct {
	V0 T0
	V1 T1
}

func New2[T0, T1 any](v0 T0, v1 T1) Tuple2[T0, T1] {
	return Tuple2[T0, T1]{v0, v1}
}

func (t Tuple2[T0, T1]) Unpack() (T0, T1) { return t.V0, t.V1 }
func (t Tuple2[T0, T1]) String() string   { return fmt.Sprintf("(%v, %v)", t.V0, t.V1) }

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
