//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleGroupBy() {
	fmt.Println(GroupBy(slices.Values([]byte("AAAABBBCCDAABBB"))))
	// Output:
	// map[65:[65 65 65 65 65 65] 66:[66 66 66 66 66 66] 67:[67 67] 68:[68]]
}

func ExampleGroupByKey() {
	s := slices.Values([]int{1, 2, 3, 4, 5, 6, 7})
	fmt.Println(GroupByKey(s, func(v int) bool { return v%2 == 0 }))
	// Output:
	// map[false:[1 3 5 7] true:[2 4 6]]
}
