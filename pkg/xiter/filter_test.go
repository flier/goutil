//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleFilter() {
	s := slices.Values([]int{1, 2, 3, 4, 5})
	f := Filter(s, func(n int) bool { return n%2 == 0 })

	fmt.Println(slices.Collect(f))

	// Output: [2 4]
}

func ExampleFilterFunc() {
	isEven := FilterFunc(func(n int) bool { return n%2 == 0 })

	s := slices.Values([]int{1, 2, 3, 4, 5})
	f := isEven(s)

	fmt.Println(slices.Collect(f))

	// Output: [2 4]
}

func ExampleFilter2() {
	s := slices.All([]int{1, 2, 3, 4, 5})
	f := Filter2(s, func(i, n int) bool { return i%2 == 0 })

	fmt.Println(maps.Collect(f))

	// Output: map[0:1 2:3 4:5]
}

func ExampleFilter2Func() {
	isEven := Filter2Func(func(i, n int) bool { return i%2 == 0 })

	s := slices.All([]int{1, 2, 3, 4, 5})
	f := isEven(s)

	fmt.Println(maps.Collect(f))

	// Output: map[0:1 2:3 4:5]
}

func ExampleFilterMap() {
	s := slices.Values([]int{1, 2, 3, 4, 5})
	f := FilterMap(s, func(n int) (int, bool) { return n * n, n%2 == 0 })

	fmt.Println(slices.Collect(f))

	// Output: [4 16]
}

func ExampleFilterMapFunc() {
	squareEven := FilterMapFunc(func(n int) (int, bool) { return n * n, n%2 == 0 })

	s := slices.Values([]int{1, 2, 3, 4, 5})
	f := squareEven(s)

	fmt.Println(slices.Collect(f))

	// Output: [4 16]
}

func ExampleFilterMap2() {
	s := slices.All([]int{1, 2, 3, 4, 5})
	f := FilterMap2(s, func(i, n int) (int, bool) { return n * n, i%2 == 0 })

	fmt.Println(maps.Collect(f))

	// Output: map[0:1 2:9 4:25]
}

func ExampleFilterMap2Func() {
	squareEven := FilterMap2Func(func(i, n int) (int, bool) { return n * n, i%2 == 0 })

	s := slices.All([]int{1, 2, 3, 4, 5})
	f := squareEven(s)

	fmt.Println(maps.Collect(f))

	// Output: map[0:1 2:9 4:25]
}
