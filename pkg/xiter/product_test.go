//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleProduct() {
	fmt.Println(Product(slices.Values([]int{1, 2, 3, 4, 5})))
	fmt.Println(Product(Empty[int]()))
	// Output:
	// 120
	// 1
}
