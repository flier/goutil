//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleCycle() {
	s := slices.Values([]int{1, 2, 3})
	c := Cycle(s)
	take6 := TakeFunc[int](6)

	fmt.Println(slices.Collect(take6(c)))
	// Output: [1 2 3 1 2 3]
}

func ExampleCycle2() {
	s := slices.All([]int{1, 2, 3})
	c := Cycle2(s)
	take6 := Take2Func[int, int](6)

	fmt.Println(slices.Collect(Pairs(take6(c))))
	// Output: [(0, 1) (1, 2) (2, 3) (0, 1) (1, 2) (2, 3)]
}
