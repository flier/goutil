//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

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

func ExampleChunkExact() {
	seq := slices.Values([]int{1, 2, 3, 4, 5})
	c := ChunkExact(seq, 2)

	fmt.Println(slices.Collect(c))
	// Output:
	// [[1 2] [3 4]]
}

func ExampleChunkExactFunc() {
	twain := ChunkExactFunc[int](2)

	seq := slices.Values([]int{1, 2, 3, 4, 5})
	c := twain(seq)

	fmt.Println(slices.Collect(c))
	// Output:
	// [[1 2] [3 4]]
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

func TestChunks(t *testing.T) {
	Convey("Chunks", t, func() {
		Convey("Should split sequence into chunks of specified size", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8})
			want := [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}

			result := slices.Collect(Chunks(input, 2))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle last chunk with fewer elements", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := [][]int{{1, 2}, {3, 4}, {5}}

			result := slices.Collect(Chunks(input, 2))
			So(result, ShouldResemble, want)
		})

		Convey("Should return empty when chunk size is 0", func() {
			input := slices.Values([]int{1, 2, 3})

			result := slices.Collect(Chunks(input, 0))
			So(result, ShouldBeEmpty)
		})

		Convey("Should return empty when chunk size is negative", func() {
			input := slices.Values([]int{1, 2, 3})

			result := slices.Collect(Chunks(input, -1))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})

			result := slices.Collect(Chunks(input, 3))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})

			result := slices.Collect(Chunks(input, 3))
			So(result, ShouldResemble, [][]int{{42}})
		})

		Convey("Should handle chunk size larger than sequence", func() {
			input := slices.Values([]int{1, 2, 3})

			result := slices.Collect(Chunks(input, 5))
			So(result, ShouldResemble, [][]int{{1, 2, 3}})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := [][]int{{1, 2}}

			var chunks [][]int

			for chunk := range Chunks(input, 2) {
				if len(chunks) == 1 {
					break
				}

				chunks = append(chunks, chunk)
			}

			So(chunks, ShouldResemble, want)
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "b", "c", "d", "e"})
			want := [][]string{{"a", "b", "c"}, {"d", "e"}}

			result := slices.Collect(Chunks(input, 3))
			So(result, ShouldResemble, want)
		})
	})
}

func TestChunksFunc(t *testing.T) {
	Convey("ChunksFunc", t, func() {
		Convey("Should create function that splits sequence into chunks", func() {
			chunks := ChunksFunc[int](3)

			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7})
			want := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}

			result := slices.Collect(chunks(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			chunks := ChunksFunc[string](2)

			input1 := slices.Values([]string{"a", "b", "c", "d"})
			input2 := slices.Values([]string{"x", "y", "z"})

			result1 := slices.Collect(chunks(input1))
			result2 := slices.Collect(chunks(input2))

			So(result1, ShouldResemble, [][]string{{"a", "b"}, {"c", "d"}})
			So(result2, ShouldResemble, [][]string{{"x", "y"}, {"z"}})
		})

		Convey("Should work with different types", func() {
			chunks := ChunksFunc[float64](2)

			input := slices.Values([]float64{1.1, 2.2, 3.3, 4.4, 5.5})
			want := [][]float64{{1.1, 2.2}, {3.3, 4.4}, {5.5}}

			result := slices.Collect(chunks(input))
			So(result, ShouldResemble, want)
		})
	})
}

func TestChunkBy(t *testing.T) {
	Convey("ChunkBy", t, func() {
		Convey("Should group consecutive elements by predicate", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			predicate := func(n int) bool { return n%2 == 1 }
			want := [][]int{{1}, {2, 2}, {3}, {4, 4, 6}, {7, 7}}

			result := slices.Collect(ChunkBy(input, predicate))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			predicate := func(n int) bool { return n%2 == 1 }

			result := slices.Collect(ChunkBy(input, predicate))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			predicate := func(n int) bool { return n%2 == 1 }

			result := slices.Collect(ChunkBy(input, predicate))
			So(result, ShouldResemble, [][]int{{42}})
		})

		Convey("Should handle all elements with same predicate result", func() {
			input := slices.Values([]int{1, 3, 5, 7})
			predicate := func(n int) bool { return n%2 == 1 }

			result := slices.Collect(ChunkBy(input, predicate))
			So(result, ShouldResemble, [][]int{{1, 3, 5, 7}})
		})

		Convey("Should handle alternating predicate results", func() {
			input := slices.Values([]int{1, 2, 2, 3, 3, 3})
			predicate := func(n int) bool { return n%2 == 1 }
			want := [][]int{{1}, {2, 2}, {3, 3, 3}}

			result := slices.Collect(ChunkBy(input, predicate))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 2, 3, 3, 3})
			predicate := func(n int) bool { return n%2 == 1 }
			want := [][]int{{1}, {2, 2}}

			var chunks [][]int

			for chunk := range ChunkBy(input, predicate) {
				if len(chunks) == 2 {
					break
				}

				chunks = append(chunks, chunk)
			}

			So(chunks, ShouldResemble, want)
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "bb", "bb", "c", "dd", "dd", "eee"})
			predicate := func(s string) bool { return len(s)%2 == 1 }

			result := slices.Collect(ChunkBy(input, predicate))
			So(result, ShouldResemble, [][]string{{"a"}, {"bb", "bb"}, {"c"}, {"dd", "dd"}, {"eee"}})
		})
	})
}

func TestChunkByFunc(t *testing.T) {
	Convey("ChunkByFunc", t, func() {
		Convey("Should create function that groups elements by predicate", func() {
			groupByOdd := ChunkByFunc(func(n int) bool { return n%2 == 1 })

			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			want := [][]int{{1}, {2, 2}, {3}, {4, 4, 6}, {7, 7}}

			result := slices.Collect(groupByOdd(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			groupByLength := ChunkByFunc(func(s string) bool { return len(s)%2 == 1 })

			input1 := slices.Values([]string{"a", "bb", "bb", "c"})
			input2 := slices.Values([]string{"x", "yy", "yy", "z"})

			result1 := slices.Collect(groupByLength(input1))
			result2 := slices.Collect(groupByLength(input2))

			So(result1, ShouldResemble, [][]string{{"a"}, {"bb", "bb"}, {"c"}})
			So(result2, ShouldResemble, [][]string{{"x"}, {"yy", "yy"}, {"z"}})
		})

		Convey("Should work with different types", func() {
			groupByPositive := ChunkByFunc(func(n float64) bool { return n > 0.0 })

			input := slices.Values([]float64{1.1, -2.2, -2.2, 3.3, -4.4, -4.4, 5.5})
			want := [][]float64{{1.1}, {-2.2, -2.2}, {3.3}, {-4.4, -4.4}, {5.5}}

			result := slices.Collect(groupByPositive(input))
			So(result, ShouldResemble, want)
		})
	})
}

func TestChunkByKey(t *testing.T) {
	Convey("ChunkByKey", t, func() {
		Convey("Should group consecutive elements by key", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			keyFunc := func(n int) bool { return n%2 == 1 }
			want := [][]int{{1}, {2, 2}, {3}, {4, 4, 6}, {7, 7}}

			result := slices.Collect(ChunkByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			keyFunc := func(n int) bool { return n%2 == 1 }

			result := slices.Collect(ChunkByKey(input, keyFunc))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			keyFunc := func(n int) bool { return n%2 == 1 }

			result := slices.Collect(ChunkByKey(input, keyFunc))
			So(result, ShouldResemble, [][]int{{42}})
		})

		Convey("Should handle all elements with same key", func() {
			input := slices.Values([]int{1, 3, 5, 7})
			keyFunc := func(n int) bool { return n%2 == 1 }

			result := slices.Collect(ChunkByKey(input, keyFunc))
			So(result, ShouldResemble, [][]int{{1, 3, 5, 7}})
		})

		Convey("Should handle alternating keys", func() {
			input := slices.Values([]int{1, 2, 2, 3, 3, 3})
			keyFunc := func(n int) int { return n % 2 }

			result := slices.Collect(ChunkByKey(input, keyFunc))
			So(result, ShouldResemble, [][]int{{1}, {2, 2}, {3, 3, 3}})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 2, 3, 3, 3})
			predicate := func(n int) int { return n % 2 }
			want := [][]int{{1}}

			var chunks [][]int

			for chunk := range ChunkByKey(input, predicate) {
				if len(chunks) == 1 {
					break
				}

				chunks = append(chunks, chunk)
			}

			So(chunks, ShouldResemble, want)
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "bb", "bb", "c", "dd", "dd", "eee"})
			keyFunc := func(s string) int { return len(s) }

			result := slices.Collect(ChunkByKey(input, keyFunc))
			So(result, ShouldResemble, [][]string{{"a"}, {"bb", "bb"}, {"c"}, {"dd", "dd"}, {"eee"}})
		})

		Convey("Should work with complex keys", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			keyFunc := func(n int) string {
				if n%2 == 0 {
					return "even"
				}
				return "odd"
			}

			result := slices.Collect(ChunkByKey(input, keyFunc))
			So(result, ShouldResemble, [][]int{{1}, {2, 2}, {3}, {4, 4, 6}, {7, 7}})
		})
	})
}

func TestChunkByKeyFunc(t *testing.T) {
	Convey("ChunkByKeyFunc", t, func() {
		Convey("Should create function that groups elements by key", func() {
			groupByOdd := ChunkByKeyFunc(func(n int) bool { return n%2 == 1 })

			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			want := [][]int{{1}, {2, 2}, {3}, {4, 4, 6}, {7, 7}}

			result := slices.Collect(groupByOdd(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			groupByLength := ChunkByKeyFunc(func(s string) int { return len(s) })

			input1 := slices.Values([]string{"a", "bb", "bb", "c"})
			input2 := slices.Values([]string{"x", "yy", "yy", "z"})

			result1 := slices.Collect(groupByLength(input1))
			result2 := slices.Collect(groupByLength(input2))

			So(result1, ShouldResemble, [][]string{{"a"}, {"bb", "bb"}, {"c"}})
			So(result2, ShouldResemble, [][]string{{"x"}, {"yy", "yy"}, {"z"}})
		})

		Convey("Should work with different types", func() {
			groupBySign := ChunkByKeyFunc(func(n float64) string {
				if n > 0 {
					return "positive"
				} else if n < 0 {
					return "negative"
				}
				return "zero"
			})

			input := slices.Values([]float64{1.1, -2.2, -2.2, 0.0, 3.3, -4.4, -4.4, 5.5})
			want := [][]float64{{1.1}, {-2.2, -2.2}, {0.0}, {3.3}, {-4.4, -4.4}, {5.5}}

			result := slices.Collect(groupBySign(input))
			So(result, ShouldResemble, want)
		})
	})
}

func TestChunkExact(t *testing.T) {
	Convey("ChunkExact", t, func() {
		Convey("Should split sequence into exact chunks of specified size", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8})
			want := [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}

			result := slices.Collect(ChunkExact(input, 2))
			So(result, ShouldResemble, want)
		})

		Convey("Should discard incomplete final chunk", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := [][]int{{1, 2}, {3, 4}}

			result := slices.Collect(ChunkExact(input, 2))
			So(result, ShouldResemble, want)
		})

		Convey("Should return empty when chunk size is 0", func() {
			input := slices.Values([]int{1, 2, 3})

			result := slices.Collect(ChunkExact(input, 0))
			So(result, ShouldBeEmpty)
		})

		Convey("Should return empty when chunk size is negative", func() {
			input := slices.Values([]int{1, 2, 3})

			result := slices.Collect(ChunkExact(input, -1))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})

			result := slices.Collect(ChunkExact(input, 3))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle sequence shorter than chunk size", func() {
			input := slices.Values([]int{1, 2})

			result := slices.Collect(ChunkExact(input, 3))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle sequence exactly matching chunk size", func() {
			input := slices.Values([]int{1, 2, 3})

			result := slices.Collect(ChunkExact(input, 3))
			So(result, ShouldResemble, [][]int{{1, 2, 3}})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5, 6})
			want := [][]int{{1, 2, 3}}

			var chunks [][]int

			for chunk := range ChunkExact(input, 3) {
				if len(chunks) == 1 {
					break
				}

				chunks = append(chunks, chunk)
			}

			So(chunks, ShouldResemble, want)
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "b", "c", "d", "e", "f"})
			want := [][]string{{"a", "b", "c"}, {"d", "e", "f"}}

			result := slices.Collect(ChunkExact(input, 3))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with custom struct types", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 35},
				{"David", 40}, {"Eve", 45}, {"Frank", 50},
			})
			want := [][]Person{
				{{"Alice", 25}, {"Bob", 30}, {"Charlie", 35}},
				{{"David", 40}, {"Eve", 45}, {"Frank", 50}},
			}

			result := slices.Collect(ChunkExact(input, 3))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle large chunk sizes", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			want := [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}

			result := slices.Collect(ChunkExact(input, 10))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle chunk size of 1", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			want := [][]int{{1}, {2}, {3}, {4}}

			result := slices.Collect(ChunkExact(input, 1))
			So(result, ShouldResemble, want)
		})
	})
}

func TestChunkExactFunc(t *testing.T) {
	Convey("ChunkExactFunc", t, func() {
		Convey("Should create function that splits sequence into exact chunks", func() {
			chunks := ChunkExactFunc[int](3)

			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7})
			want := [][]int{{1, 2, 3}, {4, 5, 6}}

			result := slices.Collect(chunks(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			chunks := ChunkExactFunc[string](2)

			input1 := slices.Values([]string{"a", "b", "c", "d"})
			input2 := slices.Values([]string{"x", "y", "z"})

			result1 := slices.Collect(chunks(input1))
			result2 := slices.Collect(chunks(input2))

			So(result1, ShouldResemble, [][]string{{"a", "b"}, {"c", "d"}})
			So(result2, ShouldResemble, [][]string{{"x", "y"}})
		})

		Convey("Should work with different types", func() {
			chunks := ChunkExactFunc[float64](2)

			input := slices.Values([]float64{1.1, 2.2, 3.3, 4.4, 5.5})
			want := [][]float64{{1.1, 2.2}, {3.3, 4.4}}

			result := slices.Collect(chunks(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle edge cases consistently", func() {
			chunks := ChunkExactFunc[int](2)

			// Empty sequence
			emptyInput := slices.Values([]int{})
			emptyResult := slices.Collect(chunks(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element (too short for chunk size 2)
			singleInput := slices.Values([]int{42})
			singleResult := slices.Collect(chunks(singleInput))
			So(singleResult, ShouldBeEmpty)

			// Exact multiple of chunk size
			exactInput := slices.Values([]int{1, 2, 3, 4})
			exactResult := slices.Collect(chunks(exactInput))
			So(exactResult, ShouldResemble, [][]int{{1, 2}, {3, 4}})
		})

		Convey("Should preserve function behavior across multiple calls", func() {
			chunks := ChunkExactFunc[int](3)

			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

			// First call
			result1 := slices.Collect(chunks(input))
			So(result1, ShouldResemble, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})

			// Second call with same function
			result2 := slices.Collect(chunks(input))
			So(result2, ShouldResemble, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})

			// Results should be identical
			So(result1, ShouldResemble, result2)
		})
	})
}
