//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	"github.com/flier/goutil/pkg/xiter"
)

func ExampleOnce() {
	s := xiter.Once(123)
	fmt.Println(slices.Collect(s))
	// Output: [123]
}

func ExampleOnce2() {
	s := xiter.Once2("hello", "word")
	fmt.Println(maps.Collect(s))
	// Output: map[hello:word]
}

func ExampleOnceWith() {
	s := xiter.OnceWith(func() int { return 123 })
	fmt.Println(slices.Collect(s))
	// Output: [123]
}

func ExampleOnceWith2() {
	s := xiter.OnceWith2(func() (string, string) { return "hello", "word" })
	fmt.Println(maps.Collect(s))
	// Output: map[hello:word]
}
