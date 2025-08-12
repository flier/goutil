//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleAny() {
	s := slices.Values([]int{1, 2, 3})

	isPowerOfTwo := func(n int) bool { return n&(n-1) == 0 }
	fmt.Println(Any(s, isPowerOfTwo))

	isZero := func(n int) bool { return n == 0 }
	fmt.Println(Any(s, isZero))

	// Output:
	// true
	// false
}

func ExampleAnyFunc() {
	s := slices.Values([]int{1, 2, 3})

	isPowerOfTwo := AnyFunc(func(n int) bool { return n&(n-1) == 0 })
	fmt.Println(isPowerOfTwo(s))

	isZero := AnyFunc(func(n int) bool { return n == 0 })
	fmt.Println(isZero(s))

	// Output:
	// true
	// false
}

func ExampleAny2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	isShortKey := func(k, v string) bool { return len(k) < 4 }
	fmt.Println(Any2(s, isShortKey))

	isLongValue := func(k, v string) bool { return len(v) > 10 }
	fmt.Println(Any2(s, isLongValue))

	// Output:
	// true
	// false
}

func ExampleAny2Func() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	isShortKey := Any2Func(func(k, v string) bool { return len(k) < 4 })
	fmt.Println(isShortKey(s))

	isLongValue := Any2Func(func(k, v string) bool { return len(v) > 10 })
	fmt.Println(isLongValue(s))

	// Output:
	// true
	// false
}
