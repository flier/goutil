//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleDedup() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	d := Dedup(s)

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 6 7]
}

func ExampleDedupBy() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	d := DedupBy(s, func(x, y int) bool { return x == y })

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 6 7]
}

func ExampleDedupByFunc() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	dedup := DedupByFunc(func(x, y int) bool { return x == y })
	d := dedup(s)

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 6 7]
}

func ExampleDedupByKey() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	d := DedupByKey(s, func(n int) int { return n % 2 })

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 7]
}

func ExampleDedupByKeyFunc() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	dedup := DedupByKeyFunc(func(n int) int { return n % 2 })
	d := dedup(s)

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 7]
}

func ExampleDedupByKey2() {
	s := Zip(
		slices.Values([]int{1, 2, 2, 3, 4}),
		slices.Values([]string{"a", "b", "c", "d", "e"}))

	d := DedupByKey2(s, func(k int, v string) int { return k })

	fmt.Println(maps.Collect(d))

	// Output:
	// map[1:a 2:b 3:d 4:e]
}

func ExampleDedupByKey2Func() {
	s := Zip(
		slices.Values([]int{1, 2, 2, 3, 4}),
		slices.Values([]string{"a", "b", "c", "d", "e"}))

	dedup := DedupByKey2Func(func(k int, v string) int { return k })
	d := dedup(s)

	fmt.Println(maps.Collect(d))

	// Output:
	// map[1:a 2:b 3:d 4:e]
}
