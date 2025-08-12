//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExamplePartition() {
	s := slices.Values([]int{1, 2, 3, 4, 5})
	odd, even := Partition(s, func(n int) bool { return n%2 != 0 }).Unpack()

	fmt.Println(slices.Collect(odd), slices.Collect(even))
	// Output:
	// [1 3 5] [2 4]
}

func ExamplePartitionFunc() {
	byEven := PartitionFunc(func(n int) bool { return n%2 != 0 })

	s := slices.Values([]int{1, 2, 3, 4, 5})
	odd, even := byEven(s).Unpack()

	fmt.Println(slices.Collect(odd), slices.Collect(even))
	// Output:
	// [1 3 5] [2 4]
}

func ExamplePartition2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	short, long := Partition2(s, func(k, v string) bool { return len(k) < 4 }).Unpack()

	fmt.Println(maps.Collect(short), maps.Collect(long))
	// Output:
	// map[foo:bar] map[hello:world]
}

func ExamplePartition2Func() {
	byLen := Partition2Func(func(k, v string) bool { return len(k) < 4 })

	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	short, long := byLen(s).Unpack()

	fmt.Println(maps.Collect(short), maps.Collect(long))
	// Output:
	// map[foo:bar] map[hello:world]
}
