//go:build go1.23

package xiter_test

import (
	"cmp"
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleMinMax() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(MinMax(s))
	// Output: (1, 5)
}

func ExampleMinMaxBy() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	fmt.Println(MinMaxBy(s, cmp.Compare))
	// Output: (bar, world)
}

func ExampleMinMaxByFunc() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	f := MinMaxByFunc[string](cmp.Compare)

	fmt.Println(f(s))
	// Output: (bar, world)
}

func ExampleMinMaxByKey() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	fmt.Println(MinMaxByKey(s, func(s string) int { return len(s) }))
	// Output: (bar, world)
}

func ExampleMinMaxByKeyFunc() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	f := MinMaxByKeyFunc(func(s string) int { return len(s) })

	fmt.Println(f(s))
	// Output: (bar, world)
}
