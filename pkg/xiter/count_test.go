//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleCount() {
	s := slices.Values([]int{1, 2, 3})
	n := Count(s)

	fmt.Println(n)

	// Output: 3
}

func ExampleCount2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	n := Count2(s)

	fmt.Println(n)

	// Output: 2
}
