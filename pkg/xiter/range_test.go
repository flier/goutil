//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleRange() {
	s := Range(1, 5)
	fmt.Println(slices.Collect(s))

	mul2 := MapFunc(func(n int) int { return n * 2 })

	m := mul2(s)
	fmt.Println(slices.Collect(m))

	// Output:
	// [1 2 3 4]
	// [2 4 6 8]
}

func ExampleRange_map() {
	mul2 := MapFunc(func(n int) int { return n * 2 })
	add1 := MapFunc(func(n int) int { return n + 1 })

	s := Range(1, 4)
	r := add1(mul2(s))

	fmt.Println(slices.Collect(r))

	// Output:
	// [3 5 7]
}

func ExampleRangeFrom() {
	s := RangeFrom(5)
	fmt.Println(slices.Collect(Take(s, 5))) // [5 6 7 8 9]

	s = RangeFrom(0)
	fmt.Println(slices.Collect(Take(s, 5))) // [0 1 2 3 4]

	s = RangeFrom(-3)
	fmt.Println(slices.Collect(Take(s, 5))) // [-3 -2 -1 0 1]

	// Output:
	// [5 6 7 8 9]
	// [0 1 2 3 4]
	// [-3 -2 -1 0 1]
}

func ExampleRangeTo() {
	s := RangeTo(5)
	fmt.Println(slices.Collect(s)) // [0 1 2 3 4]

	s = RangeTo(0)
	fmt.Println(slices.Collect(s)) // []

	s = RangeTo(-3)
	fmt.Println(slices.Collect(s)) // []

	// Output:
	// [0 1 2 3 4]
	// []
	// []
}
