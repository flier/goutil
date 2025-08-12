//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"math"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleUniq() {
	s := Uniq(slices.Values([]int{1, 2, 3, 3, 2, 1}))

	fmt.Println(slices.Collect(s))
	// Output: [1 2 3]
}

func ExampleUniqByKey() {
	s := slices.Values([]complex128{1 + 1i, -1 + 2i, -2 + 3i, 2 + 4i, -3 + 5i})
	u := UniqByKey(s, func(c complex128) int { return int(math.Abs(real(c))) })

	fmt.Println(slices.Collect(u))
	// Output:
	// [(1+1i) (-2+3i) (-3+5i)]
}

func ExampleUniqByKeyFunc() {
	abs := UniqByKeyFunc(func(c complex128) int { return int(math.Abs(real(c))) })

	s := slices.Values([]complex128{1 + 1i, -1 + 2i, -2 + 3i, 2 + 4i, -3 + 5i})
	u := abs(s)

	fmt.Println(slices.Collect(u))
	// Output:
	// [(1+1i) (-2+3i) (-3+5i)]
}

func ExampleUniqByKey2() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})
	u := UniqByKey2(s, func(i int, v string) int { return i % 2 })

	fmt.Println(maps.Collect(u))
	// Output:
	// map[0:foo 1:bar]
}

func ExampleUniqByKey2Func() {
	even := UniqByKey2Func(func(i int, v string) int { return i % 2 })

	s := slices.All([]string{"foo", "bar", "hello", "world"})
	u := even(s)

	fmt.Println(maps.Collect(u))
	// Output:
	// map[0:foo 1:bar]
}
