//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleNext() {
	s := slices.Values([]int{1, 2, 3, 4, 5})
	fmt.Println(Next(s)) // Some(1)

	take5 := SkipFunc[int](5)
	fmt.Println(Next(take5(s))) // None

	// Output:
	// Some(1)
	// None
}

func ExampleNext2() {
	fmt.Println(Next2(slices.All([]int{1, 2, 3})))
	fmt.Println(Next2(Empty2[int, int]()))

	// Output:
	// Some((0, 1))
	// None
}
