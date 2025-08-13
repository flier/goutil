//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleReduce() {
	s := slices.Values([]int{1, 2, 3})
	r := Reduce(s, func(x int, y int) int { return x + y })

	fmt.Println(r)
	// Output: 6
}

func ExampleReduceFunc() {
	sum := ReduceFunc(func(x int, y int) int { return x + y })

	s := slices.Values([]int{1, 2, 3})
	r := sum(s)

	fmt.Println(r)
	// Output: 6
}

func TestReduce(t *testing.T) {
	Convey("Reduce", t, func() {
		Convey("Should reduce integer sequence with addition", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			reduceFunc := func(a, b int) int { return a + b }
			want := 15 // 1 + 2 + 3 + 4 + 5

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should reduce integer sequence with multiplication", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			reduceFunc := func(a, b int) int { return a * b }
			want := 120 // 1 * 2 * 3 * 4 * 5

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should reduce string sequence with concatenation", func() {
			input := slices.Values([]string{"hello", " ", "world", "!"})
			reduceFunc := func(a, b string) string { return a + b }
			want := "hello world!"

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			reduceFunc := func(a, b int) int { return a + b }
			want := 0 // Zero value for int

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			reduceFunc := func(a, b int) int { return a + b }
			want := 42

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should handle two elements", func() {
			input := slices.Values([]int{10, 5})
			reduceFunc := func(a, b int) int { return a - b }
			want := 5 // 10 - 5

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should work with custom struct types", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 35},
			})
			reduceFunc := func(a, b Person) Person {
				return Person{
					Name: a.Name + " & " + b.Name,
					Age:  a.Age + b.Age,
				}
			}
			want := Person{Name: "Alice & Bob & Charlie", Age: 90}

			result := Reduce(input, reduceFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with float64", func() {
			input := slices.Values([]float64{1.5, 2.5, 3.5, 4.5})
			reduceFunc := func(a, b float64) float64 { return a + b }
			want := 12.0

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should work with boolean operations", func() {
			input := slices.Values([]bool{true, false, true, true})
			reduceFunc := func(a, b bool) bool { return a && b }
			want := false // true && false && true && true = false

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should work with boolean OR operations", func() {
			input := slices.Values([]bool{false, false, true, false})
			reduceFunc := func(a, b bool) bool { return a || b }
			want := true // false || false || true || false = true

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should work with max operation", func() {
			input := slices.Values([]int{3, 1, 4, 1, 5, 9, 2, 6})
			reduceFunc := func(a, b int) int { return max(a, b) }
			want := 9

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should work with min operation", func() {
			input := slices.Values([]int{3, 1, 4, 1, 5, 9, 2, 6})
			reduceFunc := func(a, b int) int { return min(a, b) }
			want := 1

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should preserve order of operations", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			reduceFunc := func(a, b int) int { return a - b }
			want := -8 // ((1 - 2) - 3) - 4 = -1 - 3 - 4 = -4 - 4 = -8

			result := Reduce(input, reduceFunc)
			So(result, ShouldEqual, want)
		})

		Convey("Should work with slice concatenation", func() {
			input := slices.Values([][]int{{1, 2}, {3, 4}, {5, 6}})
			reduceFunc := func(a, b []int) []int { return append(a, b...) }
			want := []int{1, 2, 3, 4, 5, 6}

			result := Reduce(input, reduceFunc)
			So(result, ShouldResemble, want)
		})
	})
}

func TestReduceFunc(t *testing.T) {
	Convey("ReduceFunc", t, func() {
		Convey("Should create reusable reduce function", func() {
			sum := ReduceFunc(func(a, b int) int { return a + b })

			input1 := slices.Values([]int{1, 2, 3})
			input2 := slices.Values([]int{4, 5, 6})

			result1 := sum(input1)
			result2 := sum(input2)

			So(result1, ShouldEqual, 6)
			So(result2, ShouldEqual, 15)
		})

		Convey("Should work with string concatenation", func() {
			concat := ReduceFunc(func(a, b string) string { return a + b })

			input := slices.Values([]string{"hello", " ", "world"})
			want := "hello world"

			result := concat(input)
			So(result, ShouldEqual, want)
		})

		Convey("Should handle edge cases consistently", func() {
			sum := ReduceFunc(func(a, b int) int { return a + b })

			// Empty sequence
			emptyInput := slices.Values([]int{})
			emptyResult := sum(emptyInput)
			So(emptyResult, ShouldEqual, 0)

			// Single element
			singleInput := slices.Values([]int{42})
			singleResult := sum(singleInput)
			So(singleResult, ShouldEqual, 42)

			// Two elements
			twoInput := slices.Values([]int{10, 5})
			twoResult := sum(twoInput)
			So(twoResult, ShouldEqual, 15)
		})

		Convey("Should work with custom comparison logic", func() {
			// Custom reduce function that finds the longest string
			longest := ReduceFunc(func(a, b string) string {
				if len(a) > len(b) {
					return a
				}
				return b
			})

			input := slices.Values([]string{"a", "hello", "world", "test"})
			want := "world"

			result := longest(input)
			So(result, ShouldEqual, want)
		})

		Convey("Should work with complex data types", func() {
			type Point struct {
				X, Y int
			}

			// Reduce to find the point with maximum distance from origin
			maxDistance := ReduceFunc(func(a, b Point) Point {
				distA := a.X*a.X + a.Y*a.Y
				distB := b.X*b.X + b.Y*b.Y
				if distA > distB {
					return a
				}
				return b
			})

			input := slices.Values([]Point{
				{1, 1}, {3, 4}, {0, 5}, {2, 2},
			})
			// Distances: {1,1}=2, {3,4}=25, {0,5}=25, {2,2}=8
			// Max distance is 25, and {0,5} comes after {3,4} with same distance
			want := Point{0, 5}

			result := maxDistance(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with different reduction strategies", func() {
			// Test different reduction strategies with the same function factory
			input := slices.Values([]int{1, 2, 3, 4, 5})

			// Sum
			sum := ReduceFunc(func(a, b int) int { return a + b })
			sumResult := sum(input)
			So(sumResult, ShouldEqual, 15)

			// Product
			product := ReduceFunc(func(a, b int) int { return a * b })
			productResult := product(input)
			So(productResult, ShouldEqual, 120)

			// Max
			max := ReduceFunc(func(a, b int) int {
				if a > b {
					return a
				}
				return b
			})
			maxResult := max(input)
			So(maxResult, ShouldEqual, 5)
		})

		Convey("Should handle nil and empty sequences consistently", func() {
			sum := ReduceFunc(func(a, b int) int { return a + b })

			// Empty sequence
			emptyResult := sum(slices.Values([]int{}))
			So(emptyResult, ShouldEqual, 0)

			// Single element
			singleResult := sum(slices.Values([]int{42}))
			So(singleResult, ShouldEqual, 42)

			// Multiple elements
			multiResult := sum(slices.Values([]int{1, 2, 3}))
			So(multiResult, ShouldEqual, 6)
		})
	})
}
