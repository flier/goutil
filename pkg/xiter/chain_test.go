//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	"github.com/flier/goutil/pkg/xiter"
)

func ExampleChain() {
	s := xiter.Chain(
		slices.Values([]int{1, 2, 3}),
		slices.Values([]int{4, 5, 6}),
		slices.Values([]int{7, 8, 9}))

	fmt.Println(slices.Collect(s))

	// Output: [1 2 3 4 5 6 7 8 9]
}

func ExampleChain2() {
	s := xiter.Chain2(
		maps.All(map[string]string{"foo": "bar"}),
		maps.All(map[string]string{"hello": "world"}))

	fmt.Println(maps.Collect(s))

	// Output: map[foo:bar hello:world]
}
