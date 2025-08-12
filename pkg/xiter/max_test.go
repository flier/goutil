//go:build go1.23

package xiter_test

import (
	"cmp"
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleMax() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(Max(s))
	// Output: 5
}

func ExampleMaxByKey() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	w := MaxByKey(s, func(k string) int { return len(k) })

	fmt.Println(w)
	// Output: world
}

func ExampleMaxByKeyFunc() {
	maxLen := MaxByKeyFunc(func(k string) int { return len(k) })

	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	w := maxLen(s)

	fmt.Println(w)
	// Output: world
}

func ExampleMaxBy() {
	s := slices.Values([]string{"foo", "bar", "baz"})
	w := MaxBy(s, cmp.Compare)

	fmt.Println(w)
	// Output: foo
}

func ExampleMaxByFunc() {
	max := MaxByFunc[string](cmp.Compare)

	s := slices.Values([]string{"foo", "bar", "baz"})
	w := max(s)

	fmt.Println(w)
	// Output: foo
}
