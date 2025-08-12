//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleZip() {
	s3 := Zip(slices.Values([]int{1, 2, 3}), slices.Values([]string{"foo", "bar", "hello"}))

	fmt.Println(maps.Collect(s3))

	// Output: map[1:foo 2:bar 3:hello]
}

func ExampleZipWith() {
	s1 := slices.Values([]int{1, 2, 3})
	s2 := slices.Values([]int{4, 5, 6})
	s3 := ZipWith(s1, s2, func(x, y int) int { return x + y })

	fmt.Println(slices.Collect(s3))

	// Output: [5 7 9]
}

func ExampleZipWithFunc() {
	zipAndAdd := ZipWithFunc(func(x, y int) int { return x + y })

	s1 := slices.Values([]int{1, 2, 3})
	s2 := slices.Values([]int{4, 5, 6})
	s3 := zipAndAdd(s1, s2)

	fmt.Println(slices.Collect(s3))

	// Output: [5 7 9]
}

func ExampleUnzip() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	k, v := Unzip(s).Unpack()

	keys := slices.Collect(k)
	values := slices.Collect(v)

	sort.Strings(keys)
	sort.Strings(values)

	fmt.Println(keys)
	fmt.Println(values)

	// Output:
	// [foo hello]
	// [bar world]
}
