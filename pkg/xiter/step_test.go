//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleStepBy() {
	s := slices.Values([]int{0, 1, 2, 3, 4, 5})
	l := StepBy(s, 2)

	fmt.Println(slices.Collect(l))
	// Output: [0 2 4]
}

func ExampleStepByFunc() {
	stepBy2 := StepByFunc[int](2)

	s := slices.Values([]int{0, 1, 2, 3, 4, 5})
	l := stepBy2(s)

	fmt.Println(slices.Collect(l))
	// Output: [0 2 4]
}

func ExampleStepBy2() {
	s := slices.All([]int{0, 1, 2, 3, 4, 5})
	l := StepBy2(s, 2)

	fmt.Println(maps.Collect(l))
	// Output: map[0:0 2:2 4:4]
}

func ExampleStepBy2Func() {
	stepBy2 := StepBy2Func[int, int](2)

	s := slices.All([]int{0, 1, 2, 3, 4, 5})
	l := stepBy2(s)

	fmt.Println(maps.Collect(l))
	// Output: map[0:0 2:2 4:4]
}
