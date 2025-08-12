//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleRepeat() {
	s := Repeat(123)

	fmt.Println(slices.Collect(Take(s, 3)))

	// Output: [123 123 123]
}

func ExampleRepeat2() {
	s := Repeat2("foo", "bar")

	fmt.Println(slices.Collect(Pairs(Take2(s, 3))))

	// Output: [(foo, bar) (foo, bar) (foo, bar)]
}

func ExampleRepeatN() {
	s := RepeatN(123, 3)

	fmt.Println(slices.Collect(s))

	// Output: [123 123 123]
}

func ExampleRepeatN2() {
	s := RepeatN2("foo", "bar", 3)

	fmt.Println(slices.Collect(Pairs(s)))

	// Output: [(foo, bar) (foo, bar) (foo, bar)]
}

func ExampleRepeatWith() {
	n := 0
	s := RepeatWith(func() int {
		n += 1

		return n
	})

	fmt.Println(slices.Collect(Take(s, 3)))

	// Output: [1 2 3]
}

func ExampleRepeatWith2() {
	n := 0
	s := RepeatWith2(func() (int, int) {
		n += 1

		return n, n * n
	})

	fmt.Println(maps.Collect(Take2(s, 3)))

	// Output: map[1:1 2:4 3:9]
}
