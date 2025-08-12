//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleNth() {
	s := slices.Values([]int{1, 2, 3, 4, 5, 6})

	fmt.Println(Nth(s, 0))  // Some(1)
	fmt.Println(Nth(s, 4))  // Some(5)
	fmt.Println(Nth(s, 10)) // None

	// Output:
	// Some(1)
	// Some(5)
	// None
}

func ExampleNthFunc() {
	s := slices.Values([]int{1, 2, 3, 4, 5, 6})

	n0 := NthFunc[int](0)
	fmt.Println(n0(s)) // Some(1)

	n4 := NthFunc[int](4)
	fmt.Println(n4(s)) // Some(5)

	n10 := NthFunc[int](10)
	fmt.Println(n10(s)) // None

	// Output:
	// Some(1)
	// Some(5)
	// None
}

func ExampleNth2() {
	s := slices.All([]int{1, 2, 3, 4, 5, 6})

	fmt.Println(Nth2(s, 0))  // Some((0, 1))
	fmt.Println(Nth2(s, 4))  // Some((4, 5))
	fmt.Println(Nth2(s, 10)) // None

	// Output:
	// Some((0, 1))
	// Some((4, 5))
	// None
}

func ExampleNth2Func() {
	s := slices.All([]int{1, 2, 3, 4, 5, 6})

	n0 := Nth2Func[int, int](0)
	fmt.Println(n0(s)) // Some((0, 1))

	n4 := Nth2Func[int, int](4)
	fmt.Println(n4(s)) // Some((4, 5))

	n10 := Nth2Func[int, int](10)
	fmt.Println(n10(s)) // None

	// Output:
	// Some((0, 1))
	// Some((4, 5))
	// None
}
