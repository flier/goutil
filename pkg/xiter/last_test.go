//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleLast() {
	fmt.Println(Last(slices.Values([]int{1, 2, 3})))
	fmt.Println(Last(Empty[int]()))

	// Output:
	// Some(3)
	// None
}

func ExampleLast2() {
	fmt.Println(Last2(slices.All([]int{1, 2, 3})))
	fmt.Println(Last2(Empty2[int, int]()))

	// Output:
	// Some((2, 3))
	// None
}
