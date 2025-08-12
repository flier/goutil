//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"
	"strings"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleIntersperse() {
	s := slices.Values([]string{"foo", "bar", "baz"})
	i := Intersperse(s, ",")

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}

func ExampleIntersperseFunc() {
	sep := IntersperseFunc(",")

	s := slices.Values([]string{"foo", "bar", "baz"})
	i := sep(s)

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}

func ExampleIntersperseWith() {
	s := slices.Values([]string{"foo", "bar", "baz"})
	i := IntersperseWith(s, func() string { return "," })

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}

func ExampleIntersperseWithFunc() {
	sep := IntersperseWithFunc(func() string { return "," })

	s := slices.Values([]string{"foo", "bar", "baz"})
	i := sep(s)

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}
