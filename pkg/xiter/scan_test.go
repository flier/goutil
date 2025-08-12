//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleScan() {
	s := slices.Values([]int{1, 2, 3, 4})

	state := 1

	l := Scan(s, &state, func(ctx *int, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), *ctx <= 6
	})

	fmt.Println(slices.Collect(l))
	// Output: [-1 -2 -6]
}

func ExampleScanFunc() {
	state := 1
	product := ScanFunc(&state, func(ctx *int, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), *ctx <= 6
	})

	s := slices.Values([]int{1, 2, 3, 4})
	l := product(s)

	fmt.Println(slices.Collect(l))
	// Output: [-1 -2 -6]
}

func ExampleScan2() {
	s := slices.All([]int{1, 2, 3, 4})

	state := 1

	l := Scan2(s, &state, func(ctx *int, i, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), i < 3
	})

	fmt.Println(maps.Collect(l))
	// Output: map[0:-1 1:-2 2:-6]
}

func ExampleScan2Func() {
	state := 1
	product := Scan2Func(&state, func(ctx *int, i, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), i < 3
	})

	s := slices.All([]int{1, 2, 3, 4})
	l := product(s)

	fmt.Println(maps.Collect(l))
	// Output: map[0:-1 1:-2 2:-6]
}
