//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleTake() {
	s := slices.Values([]int{1, 2, 3})

	fmt.Println(slices.Collect(Take(s, 0)))
	fmt.Println(slices.Collect(Take(s, 2))) // [1 2]
	fmt.Println(slices.Collect(Take(s, 5))) // [1 2 3]

	// Output:
	// []
	// [1 2]
	// [1 2 3]
}

func ExampleTakeFunc() {
	s := slices.Values([]int{1, 2, 3})

	take0 := TakeFunc[int](0)
	fmt.Println(slices.Collect(take0(s))) // []

	take2 := TakeFunc[int](2)
	fmt.Println(slices.Collect(take2(s))) // [1 2]

	take5 := TakeFunc[int](5)
	fmt.Println(slices.Collect(take5(s))) // [1 2 3]

	// Output:
	// []
	// [1 2]
	// [1 2 3]
}

func ExampleTake2() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})

	fmt.Println(maps.Collect(Take2(s, 0))) // map[]
	fmt.Println(maps.Collect(Take2(s, 1))) // map[0:foo]
	fmt.Println(maps.Collect(Take2(s, 5))) // map[0:foo 1:bar 2:hello 3:world]
	// Output:
	// map[]
	// map[0:foo]
	// map[0:foo 1:bar 2:hello 3:world]
}

func ExampleTake2Func() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})

	take0 := Take2Func[int, string](0)
	fmt.Println(maps.Collect(take0(s))) // map[]

	take1 := Take2Func[int, string](1)
	fmt.Println(maps.Collect(take1(s))) // map[0:foo]

	take5 := Take2Func[int, string](5)
	fmt.Println(maps.Collect(take5(s))) // map[0:foo 1:bar 2:hello 3:world]
	// Output:
	// map[]
	// map[0:foo]
	// map[0:foo 1:bar 2:hello 3:world]
}

func ExampleTakeWhile() {
	s := slices.Values([]int{1, 2, 3})
	t := TakeWhile(s, func(n int) bool { return n <= 2 })

	fmt.Println(slices.Collect(t))

	// Output:
	// [1 2]
}

func ExampleTakeWhileFunc() {
	takeLe2 := TakeWhileFunc(func(n int) bool { return n <= 2 })

	s := slices.Values([]int{1, 2, 3})
	t := takeLe2(s)

	fmt.Println(slices.Collect(t))

	// Output:
	// [1 2]
}

func ExampleTakeWhile2() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})
	t := TakeWhile2(s, func(i int, w string) bool { return len(w) <= 3 })

	fmt.Println(maps.Collect(t))

	// Output:
	// map[0:foo 1:bar]
}

func ExampleTakeWhile2Func() {
	takeLenGe3 := TakeWhile2Func(func(i int, w string) bool { return len(w) <= 3 })

	s := slices.All([]string{"foo", "bar", "hello", "world"})
	t := takeLenGe3(s)

	fmt.Println(maps.Collect(t))

	// Output:
	// map[0:foo 1:bar]
}
