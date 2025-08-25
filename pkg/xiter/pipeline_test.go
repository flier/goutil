//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExamplePipeline() {
	r := Pipeline(slices.Values([]int{1, 2, 3, 4, 5}),
		MapFunc(func(x int) int { return x * 2 }),
		FilterFunc(func(x int) bool { return x > 5 }))

	fmt.Println(slices.Collect(r))

	// Output:
	// [6 8 10]
}

func ExamplePipelineFunc() {
	p := PipelineFunc(
		MapFunc(func(x int) int { return x * 2 }),
		FilterFunc(func(x int) bool { return x > 5 }))

	s := slices.Values([]int{1, 2, 3, 4, 5})
	r := p(s)

	fmt.Println(slices.Collect(r))

	// Output:
	// [6 8 10]
}

func ExamplePipeline2() {
	r := Pipeline2(slices.All([]int{1, 2, 3, 4, 5}),
		Filter2Func(func(i, x int) bool { return i%2 == 0 }),
		MapValueFunc(func(i, x int) int { return x * 2 }))

	fmt.Println(maps.Collect(r))

	// Output:
	// map[0:2 2:6 4:10]
}

func ExamplePipeline2Func() {
	p := Pipeline2Func(
		Filter2Func(func(i, x int) bool { return i%2 == 0 }),
		MapValueFunc(func(i, x int) int { return x * 2 }))

	s := slices.All([]int{1, 2, 3, 4, 5})
	r := p(s)

	fmt.Println(maps.Collect(r))

	// Output:
	// map[0:2 2:6 4:10]
}
