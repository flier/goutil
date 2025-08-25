//go:build go1.23

package xiter_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleMap() {
	s := slices.Values([]int{1, 2, 3})
	m := Map(s, func(n int) int { return n * n })

	fmt.Println(slices.Collect(m))
	// Output: [1 4 9]
}

func ExampleMapFunc() {
	square := MapFunc(func(n int) int { return n * n })

	s := slices.Values([]int{1, 2, 3})
	m := square(s)

	fmt.Println(slices.Collect(m))
	// Output: [1 4 9]
}

func ExampleMapValue() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})
	m := MapValue(s, func(n int, v string) int { return len(v) })

	fmt.Println(maps.Collect(m))
	// Output: map[0:3 1:3 2:5 3:5]
}

func ExampleMapValueFunc() {
	lengthOfValue := MapValueFunc(func(n int, v string) int { return len(v) })

	s := slices.All([]string{"foo", "bar", "hello", "world"})
	m := lengthOfValue(s)

	fmt.Println(maps.Collect(m))
	// Output: map[0:3 1:3 2:5 3:5]
}

func ExampleFlatMapFunc() {
	square := FlatMapFunc(func(n int) iter.Seq[int] {
		return slices.Values([]int{n, n * n})
	})

	s := slices.Values([]int{1, 2, 3})
	m := square(s)

	fmt.Println(slices.Collect(m))
	// Output: [1 1 2 4 3 9]
}

func ExampleFlatMap2Func() {
	keyAndLenghtOfValue := FlatMap2Func(func(k string, v string) iter.Seq2[string, int] {
		return Once2(k, len(v))
	})

	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	m := keyAndLenghtOfValue(s)

	fmt.Println(maps.Collect(m))
	// Output: map[foo:3 hello:5]
}

func ExampleMapWhile() {
	s := slices.Values([]int{1, 2, 3})
	m := MapWhile(s, func(n int) (int, bool) { return n * n, n < 3 })

	fmt.Println(slices.Collect(m))
	// Output: [1 4]
}

func ExampleMapWhileFunc() {
	square := MapWhileFunc(func(n int) (int, bool) { return n * n, n < 3 })

	s := slices.Values([]int{1, 2, 3})
	m := square(s)

	fmt.Println(slices.Collect(m))
	// Output: [1 4]
}

func ExampleMapWhile2() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})
	m := MapWhile2(s, func(n int, v string) (int, bool) { return len(v), len(v) <= 3 })

	fmt.Println(maps.Collect(m))
	// Output: map[0:3 1:3]
}

func ExampleMapWhile2Func() {
	lengthOfValue := MapWhile2Func(func(n int, v string) (int, bool) { return len(v), len(v) <= 3 })

	s := slices.All([]string{"foo", "bar", "hello", "world"})
	m := lengthOfValue(s)

	fmt.Println(maps.Collect(m))
	// Output: map[0:3 1:3]
}
