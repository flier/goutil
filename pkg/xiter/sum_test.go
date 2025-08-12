//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"
	"time"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleSum() {
	fmt.Println(Sum(slices.Values([]int{1, 2, 3})))
	fmt.Println(Sum(slices.Values([]time.Duration{time.Second * 15, time.Minute})))

	// Output:
	// 6
	// 1m15s
}

func ExampleSumBy() {

	s := slices.Values([]string{"foo", "bar", "baz"})
	n := SumBy(s, func(s string) int { return len(s) })

	fmt.Println(n)

	// Output:
	// 9
}

func ExampleSumByFunc() {
	sumLen := SumByFunc(func(s string) int { return len(s) })

	s := slices.Values([]string{"foo", "bar", "baz"})
	n := sumLen(s)

	fmt.Println(n)

	// Output:
	// 9
}
