//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleFirst() {
	s := slices.Values([]int{1, 2, 3})
	fmt.Println(First(s)) // Some(1)

	s = slices.Values([]int{})
	fmt.Println(First(s)) // None

	// Output:
	// Some(1)
	// None
}

func ExampleFirst2() {
	s := slices.All([]string{
		"one",
		"two",
		"three",
	})

	fmt.Println(First2(s)) // Some((0, one))

	s = slices.All([]string{})
	fmt.Println(First2(s)) // None

	// Output:
	// Some((0, one))
	// None
}
