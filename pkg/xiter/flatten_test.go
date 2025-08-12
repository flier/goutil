//go:build go1.23

package xiter_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleFlatten() {
	s1 := slices.Values([]int{1, 2, 3})
	s2 := slices.Values([]int{4, 5})
	s := slices.Values([]iter.Seq[int]{s1, s2})
	f := Flatten(s)

	fmt.Println(slices.Collect(f))
	// Output: [1 2 3 4 5]
}

func ExampleFlatten2() {
	s1 := slices.All([]int{1, 2, 3})
	s2 := slices.All([]int{4, 5})
	s := slices.Values([]iter.Seq2[int, int]{s1, s2})
	f := Flatten2(s)

	fmt.Println(maps.Collect(f))
	// Output: map[0:4 1:5 2:3]
}
