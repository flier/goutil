//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleEnumerate() {
	s := slices.Values([]string{"foo", "bar"})
	e := Enumerate(s)

	fmt.Println(maps.Collect(e))

	// Output: map[0:foo 1:bar]
}
