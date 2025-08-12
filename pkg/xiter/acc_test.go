//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleAccumulate() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(slices.Collect(Accumulate(s)))

	// Output:
	// [1 3 6 10 15]
}

func ExampleAccumulateBy() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(slices.Collect(AccumulateBy(s, func(acc, v int) int { return acc * v })))

	// Output:
	// [1 2 6 24 120]
}
