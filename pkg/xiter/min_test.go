//go:build go1.23

package xiter_test

import (
	"cmp"
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleMin() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(Min(s))
	// Output: 1
}

func ExampleMinByKey() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	w := MinByKey(s, func(k string) int { return len(k) })

	fmt.Println(w)
	// Output: bar
}

func ExampleMinByKeyFunc() {
	minLen := MinByKeyFunc(func(k string) int { return len(k) })

	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	w := minLen(s)

	fmt.Println(w)
	// Output: bar
}

func ExampleMinBy() {
	s := slices.Values([]string{"foo", "bar", "baz"})
	w := MinBy(s, cmp.Compare)

	fmt.Println(w)
	// Output: bar
}

func ExampleMinByFunc() {
	min := MinByFunc[string](cmp.Compare)

	s := slices.Values([]string{"foo", "bar", "baz"})
	w := min(s)

	fmt.Println(w)
	// Output: bar
}
