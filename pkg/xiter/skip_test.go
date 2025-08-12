//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleSkip() {
	s := slices.Values([]int{1, 2, 3})
	l := Skip(s, 1)

	fmt.Println(slices.Collect(l))

	// Output: [2 3]
}

func ExampleSkipFunc() {
	skip1 := SkipFunc[int](1)

	s := slices.Values([]int{1, 2, 3})
	l := skip1(s)

	fmt.Println(slices.Collect(l))

	// Output: [2 3]
}

func ExampleSkip2() {
	s := slices.All([]int{1, 2, 3})
	l := Skip2(s, 1)

	fmt.Println(maps.Collect(l))

	// Output: map[1:2 2:3]
}

func ExampleSkip2Func() {
	skip1 := Skip2Func[int, int](1)

	s := slices.All([]int{1, 2, 3})
	l := skip1(s)

	fmt.Println(maps.Collect(l))

	// Output: map[1:2 2:3]
}

func ExampleSkipWhile() {
	s := slices.Values([]int{1, 2, 3})
	l := SkipWhile(s, func(n int) bool { return n < 2 })

	fmt.Println(slices.Collect(l))

	// Output: [2 3]
}

func ExampleSkipWhileFunc() {
	lt2 := SkipWhileFunc(func(n int) bool { return n < 2 })

	s := slices.Values([]int{1, 2, 3})
	l := lt2(s)

	fmt.Println(slices.Collect(l))

	// Output: [2 3]
}

func ExampleSkipWhile2() {
	s := slices.All([]int{1, 2, 3})
	l := SkipWhile2(s, func(i, n int) bool { return n < 2 })

	fmt.Println(maps.Collect(l))

	// Output: map[1:2 2:3]
}
func ExampleSkipWhile2Func() {
	lt2 := SkipWhile2Func(func(i, n int) bool { return n < 2 })

	s := slices.All([]int{1, 2, 3})
	l := lt2(s)

	fmt.Println(maps.Collect(l))

	// Output: map[1:2 2:3]
}
