//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleKeys() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	k := Keys(s)

	fmt.Println(slices.Sorted(k))
	// Output:
	// [foo hello]
}

func ExampleValues() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	k := Values(s)

	fmt.Println(slices.Sorted(k))
	// Output:
	// [bar world]
}

func ExampleSwap() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	k := Swap(s)

	fmt.Println(maps.Collect(k))
	// Output:
	// map[bar:foo world:hello]
}
