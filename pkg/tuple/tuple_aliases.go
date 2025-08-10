//go:build go1.24

package tuple

type (
	T2[T0, T1 any]                     = Tuple2[T0, T1]
	T3[T0, T1, T2 any]                 = Tuple3[T0, T1, T2]
	T4[T0, T1, T2, T3 any]             = Tuple4[T0, T1, T2, T3]
	T5[T0, T1, T2, T3, T4 any]         = Tuple5[T0, T1, T2, T3, T4]
	T6[T0, T1, T2, T3, T4, T5 any]     = Tuple6[T0, T1, T2, T3, T4, T5]
	T7[T0, T1, T2, T3, T4, T5, T6 any] = Tuple7[T0, T1, T2, T3, T4, T5, T6]
)
