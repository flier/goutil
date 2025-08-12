//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleForEach() {
	s := slices.Values([]int{1, 2, 3})

	ForEach(s, func(n int) { fmt.Println(n) })

	// Output:
	// 1
	// 2
	// 3
}

func ExampleForEachFunc() {
	printAll := ForEachFunc(func(n int) { fmt.Println(n) })

	s := slices.Values([]int{1, 2, 3})
	printAll(s)

	// Output:
	// 1
	// 2
	// 3
}

func ExampleForEach2() {
	s := slices.All([]int{1, 2, 3})
	printAll := func(i, n int) { fmt.Println(i, n) }

	ForEach2(s, printAll)

	// Output:
	// 0 1
	// 1 2
	// 2 3
}

func ExampleForEach2Func() {
	printAll := ForEach2Func(func(i, n int) { fmt.Println(i, n) })

	s := slices.All([]int{1, 2, 3})
	printAll(s)

	// Output:
	// 0 1
	// 1 2
	// 2 3
}
