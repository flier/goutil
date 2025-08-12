//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExamplePosition() {
	s := slices.Values([]int{1, 2, 3})

	fmt.Println(Position(s, func(n int) bool { return n%2 == 0 })) // 1
	fmt.Println(Position(s, func(n int) bool { return n > 7 }))    // -1

	// Output:
	// 1
	// -1
}

func ExamplePositionFunc() {
	s := slices.Values([]int{1, 2, 3})

	even := PositionFunc(func(n int) bool { return n%2 == 0 })
	fmt.Println(even(s)) // 1

	gt7 := PositionFunc(func(n int) bool { return n > 7 })
	fmt.Println(gt7(s)) // -1

	// Output:
	// 1
	// -1
}

func ExamplePosition2() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})

	fmt.Println(Position2(s, func(n int, v string) bool { return v == "bar" })) // 1
	fmt.Println(Position2(s, func(n int, v string) bool { return v == "baz" })) // -1

	// Output:
	// 1
	// -1
}

func ExamplePosition2Func() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})

	bar := Position2Func(func(n int, v string) bool { return v == "bar" })
	fmt.Println(bar(s)) // 1

	baz := Position2Func(func(n int, v string) bool { return v == "baz" })
	fmt.Println(baz(s)) // -1

	// Output:
	// 1
	// -1
}
