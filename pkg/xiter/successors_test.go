//go:build go1.23

package xiter_test

import (
	"fmt"
	"math"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleSuccessors() {
	s := Successors(1, func(n int) (int, bool) {
		n *= 10

		return n, n < math.MaxUint16
	})

	fmt.Println(slices.Collect(s))

	// Output: [1 10 100 1000 10000]
}
