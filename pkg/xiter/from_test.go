//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleFromIndex() {
	s := FromIndex(1)

	fmt.Println(slices.Collect(Take(s, 3)))
	// Output: [1 2 3]
}

func ExampleFromIndexBy() {
	s := FromIndexBy(0, func(n int) int { return n * 10 })

	fmt.Println(slices.Collect(Take(s, 3)))
	// Output: [0 10 20]
}

func ExampleFromFunc() {
	var count int
	s := FromFunc(func() (int, bool) {
		count += 1

		return count, count < 6
	})

	fmt.Println(slices.Collect(s))
	// Output:
	// [1 2 3 4 5]
}

func ExampleFromFunc2() {
	var count int
	s := FromFunc2(func() (int, int, bool) {
		count += 1

		return count, count * count, count < 6
	})

	fmt.Println(maps.Collect(s))
	// Output:
	// map[1:1 2:4 3:9 4:16 5:25]
}

func ExampleFromChan() {
	c := make(chan int)

	go func() {
		defer close(c)

		for i := 0; i < 6; i++ {
			c <- i
		}
	}()

	s := FromChan(c)
	fmt.Println(slices.Collect(s))
	// Output:
	// [0 1 2 3 4 5]
}
