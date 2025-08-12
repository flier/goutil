//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleAll() {
	s := slices.Values([]int{1, 2, 3})

	lessThan10 := func(n int) bool { return n < 10 }
	fmt.Println(All(s, lessThan10))

	isEven := func(n int) bool { return n%2 == 0 }
	fmt.Println(All(s, isEven))

	// Output:
	// true
	// false
}

func ExampleAllFunc() {
	s := slices.Values([]int{1, 2, 3})

	lessThan10 := AllFunc(func(n int) bool { return n < 10 })
	fmt.Println(lessThan10(s)) // true

	isEven := AllFunc(func(n int) bool { return n%2 == 0 })
	fmt.Println(isEven(s)) // false

	// Output:
	// true
	// false
}

func ExampleAll2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	isShortKey := func(k, v string) bool { return len(k) < 10 }
	fmt.Println(All2(s, isShortKey))

	isLongValue := func(k, v string) bool { return len(v) > 3 }
	fmt.Println(All2(s, isLongValue))

	// Output:
	// true
	// false
}

func ExampleAll2Func() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	isShortKey := All2Func(func(k, v string) bool { return len(k) < 10 })
	fmt.Println(isShortKey(s)) // true

	isLongValue := All2Func(func(k, v string) bool { return len(v) > 3 })
	fmt.Println(isLongValue(s)) // false

	// Output:
	// true
	// false
}
