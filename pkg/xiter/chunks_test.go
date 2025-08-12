//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleChunks() {
	s := slices.Values([]int{1, 2, 3, 4, 5})
	c := Chunks(s, 2)

	fmt.Println(slices.Collect(c))
	// Output:
	// [[1 2] [3 4] [5]]
}

func ExampleChunksFunc() {
	twain := ChunksFunc[int](2)

	s := slices.Values([]int{1, 2, 3, 4, 5})
	c := twain(s)

	fmt.Println(slices.Collect(c))
	// Output:
	// [[1 2] [3 4] [5]]
}

func ExampleChunkBy() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	groupByOdd := func(n int) bool { return (n % 2) == 1 }

	fmt.Println(slices.Collect(ChunkBy(s, groupByOdd)))
	// Output:
	// [[1] [2 2] [3] [4 4 6] [7 7]]
}

func ExampleChunkByFunc() {
	groupByOdd := ChunkByFunc(func(n int) bool { return (n % 2) == 1 })

	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	g := groupByOdd(s)

	fmt.Println(slices.Collect(g))
	// Output:
	// [[1] [2 2] [3] [4 4 6] [7 7]]
}

func ExampleChunkByKey() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	groupByOdd := func(n int) bool { return (n % 2) == 1 }

	fmt.Println(slices.Collect(ChunkByKey(s, groupByOdd)))
	// Output:
	// [[1] [2 2] [3] [4 4 6] [7 7]]
}

func ExampleChunkByKeyFunc() {
	groupByOdd := ChunkByKeyFunc(func(n int) bool { return (n % 2) == 1 })

	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	g := groupByOdd(s)

	fmt.Println(slices.Collect(g))
	// Output:
	// [[1] [2 2] [3] [4 4 6] [7 7]]
}
