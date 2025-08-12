//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleEmpty() {
	fmt.Println(slices.Collect(Empty[int]()))
	// Output: []
}

func ExampleEmpty2() {
	fmt.Println(maps.Collect(Empty2[string, string]()))
	// Output: map[]
}
