//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
	. "github.com/flier/goutil/pkg/xiter/inspect"
)

func ExampleInspect() {
	s := RangeTo(5)

	square := MapFunc(func(n int) int { return n * n })
	dump := InspectFunc[int](Label("square"))
	s = dump(square(s))

	add2 := MapFunc(func(n int) int { return n + 2 })
	dump = InspectFunc[int](Label("add2"))
	s = dump(add2(s))

	fmt.Println(slices.Collect(s))

	s = RangeTo(20)
	dump = InspectFunc[int](Width(20), Limit(15))
	n := Sum(dump(s))
	fmt.Println(n)

	// Output:
	// square: [0 1 4 9 16]
	// add2: [2 3 6 11 18]
	// [2 3 6 11 18]
	// [0 1 2 3 4 5 6 7 8 9
	//  10 11 12 13 14 ...]
	// 190
}

func ExampleInspectFunc() {
	s := RangeTo(5)

	square := MapFunc(func(n int) int { return n * n })
	dump1 := InspectFunc[int](Label("square"))

	add2 := MapFunc(func(n int) int { return n + 2 })
	dump2 := InspectFunc[int](Label("add2"))

	s = dump2(add2(dump1(square(s))))

	fmt.Println(slices.Collect(s))

	s = RangeTo(20)
	dump3 := InspectFunc[int](Width(20), Limit(15))
	n := Sum(dump3(s))
	fmt.Println(n)

	// Output:
	// square: [0 1 4 9 16]
	// add2: [2 3 6 11 18]
	// [2 3 6 11 18]
	// [0 1 2 3 4 5 6 7 8 9
	//  10 11 12 13 14 ...]
	// 190
}

func ExampleInspect2() {
	s := slices.All([]string{"foo", "bar", "hello"})

	lengthOf := Map2Func(func(n int, k string) int { return len(k) })

	fmt.Println(maps.Collect(lengthOf(Inspect2(s, Label("len")))))

	// Output:
	// len: [0:foo 1:bar 2:hello]
	// map[0:3 1:3 2:5]
}

func ExampleInspect2Func() {
	s := slices.All([]string{"foo", "bar", "hello"})

	lengthOf := Map2Func(func(n int, k string) int { return len(k) })
	dump := Inspect2Func[int, string](Label("len"))

	fmt.Println(maps.Collect(lengthOf(dump(s))))

	// Output:
	// len: [0:foo 1:bar 2:hello]
	// map[0:3 1:3 2:5]
}
