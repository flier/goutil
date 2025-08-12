//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleFind() {
	s := slices.Values([]int{1, 2, 3})

	fmt.Println(Find(s, func(n int) bool { return n%2 == 0 })) // Some(2)
	fmt.Println(Find(s, func(n int) bool { return n > 7 }))    // None

	// Output:
	// Some(2)
	// None
}

func ExampleFindFunc() {
	s := slices.Values([]int{1, 2, 3})

	even := FindFunc(func(n int) bool { return n%2 == 0 })
	fmt.Println(even(s)) // Some(2)

	greatThan7 := FindFunc(func(n int) bool { return n > 7 })
	fmt.Println(greatThan7(s)) // None

	// Output:
	// Some(2)
	// None
}

func ExampleFind2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	fmt.Println(Find2(s, func(k, v string) bool { return k == "foo" })) // Some({foo bar})
	fmt.Println(Find2(s, func(k, v string) bool { return k == "baz" })) // None

	// Output:
	// Some((foo, bar))
	// None
}

func ExampleFind2Func() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	foo := Find2Func(func(k, v string) bool { return k == "foo" })
	fmt.Println(foo(s)) // Some({foo bar})

	baz := Find2Func(func(k, v string) bool { return k == "baz" })
	fmt.Println(baz(s)) // None

	// Output:
	// Some((foo, bar))
	// None
}

func ExampleFindMap() {
	s := slices.Values([]int{1, 2, 3})

	fmt.Println(FindMap(s, func(n int) (int, bool) { return n * n, n%2 == 0 })) // Some(4)
	fmt.Println(FindMap(s, func(n int) (int, bool) { return n * n, n > 7 }))    // None
	// Output:
	// Some(4)
	// None
}

func ExampleFindMapFunc() {
	s := slices.Values([]int{1, 2, 3})

	squareEven := FindMapFunc(func(n int) (int, bool) { return n * n, n%2 == 0 })
	fmt.Println(squareEven(s)) // Some(4)

	squareGt7 := FindMapFunc(func(n int) (int, bool) { return n * n, n > 7 })
	fmt.Println(squareGt7(s)) // None
	// Output:
	// Some(4)
	// None
}

func ExampleFindMap2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	fmt.Println(FindMap2(s, func(k, v string) (int, bool) { return len(v), k == "foo" })) // Some({foo 3})
	fmt.Println(FindMap2(s, func(k, v string) (int, bool) { return len(v), k == "baz" })) // None
	// Output:
	// Some((foo, 3))
	// None
}

func ExampleFindMap2Func() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})

	fooLen := FindMap2Func(func(k, v string) (int, bool) { return len(v), k == "foo" })
	fmt.Println(fooLen(s)) // Some({foo 3})

	bazLen := FindMap2Func(func(k, v string) (int, bool) { return len(v), k == "baz" })
	fmt.Println(bazLen(s)) // None
	// Output:
	// Some((foo, 3))
	// None
}
