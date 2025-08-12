//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleFold() {
	s := slices.Values([]int{1, 2, 3})
	f := Fold(s, 0, func(acc int, n int) int { return acc + n })

	fmt.Println(f)
	// Output: 6
}

func ExampleFoldFunc() {
	sum := FoldFunc(0, func(acc int, n int) int { return acc + n })

	s := slices.Values([]int{1, 2, 3})
	f := sum(s)

	fmt.Println(f)
	// Output: 6
}

func ExampleFold2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	f := Fold2(s, 0, func(sz int, k, v string) int {
		return sz + len(k) + len(v)
	})

	fmt.Println(f)
	// Output: 16
}

func ExampleFold2Func() {
	sizeOf := Fold2Func(0, func(sz int, k, v string) int {
		return sz + len(k) + len(v)
	})

	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	f := sizeOf(s)

	fmt.Println(f)
	// Output: 16
}
