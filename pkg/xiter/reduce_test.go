//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleReduce() {
	s := slices.Values([]int{1, 2, 3})
	r := Reduce(s, func(x int, y int) int { return x + y })

	fmt.Println(r)
	// Output: 6
}

func ExampleReduceFunc() {
	sum := ReduceFunc(func(x int, y int) int { return x + y })

	s := slices.Values([]int{1, 2, 3})
	r := sum(s)

	fmt.Println(r)
	// Output: 6
}
