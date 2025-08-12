//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	"github.com/flier/goutil/pkg/tuple"
	. "github.com/flier/goutil/pkg/xiter"
)

func ExamplePairs() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})
	p := Pairs(s)

	ForEach(p, func(p tuple.Tuple2[int, string]) { fmt.Println(p) })

	// unordered output:
	// (0, foo)
	// (1, bar)
	// (2, hello)
	// (3, world)
}

func ExampleUnpairs() {
	s := slices.Values([]tuple.Tuple2[string, string]{
		tuple.New2("foo", "bar"),
		tuple.New2("hello", "world"),
	})
	p := Unpairs(s)

	fmt.Println(maps.Collect(p))

	// Output: map[foo:bar hello:world]
}
