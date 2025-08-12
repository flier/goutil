//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	"github.com/flier/goutil/pkg/tuple"
	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleIterate() {
	s := Iterate(1, func(n int) int { return n * 2 })

	fmt.Println(slices.Collect(Take(s, 5)))
	// Output: [1 2 4 8 16]
}

func ExampleIterate2() {
	s := Iterate2[int](tuple.New2(1, 1), func(k, v int) (int, int) { return k + 1, v * 2 })

	fmt.Println(maps.Collect(Take2(s, 5)))
	// Output: map[1:1 2:2 3:4 4:8 5:16]
}
