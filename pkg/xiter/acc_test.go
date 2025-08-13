//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleAccumulate() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(slices.Collect(Accumulate(s)))

	// Output:
	// [1 3 6 10 15]
}

func ExampleAccumulateBy() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(slices.Collect(AccumulateBy(s, func(acc, v int) int { return acc * v })))

	// Output:
	// [1 2 6 24 120]
}

func TestAccumulate(t *testing.T) {
	Convey("Given some sequence", t, func() {
		Convey("Should accumulate integers", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 3, 6, 10, 15}

			result := slices.Collect(Accumulate(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should accumulate floats", func() {
			input := slices.Values([]float64{1.5, 2.5, 3.5})
			want := []float64{1.5, 4.0, 7.5}

			result := slices.Collect(Accumulate(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			want := []int{42}

			result := slices.Collect(Accumulate(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})

			result := slices.Collect(Accumulate(input))
			So(result, ShouldBeNil)
		})

		Convey("Should handle negative numbers", func() {
			input := slices.Values([]int{-1, -2, -3})
			want := []int{-1, -3, -6}

			result := slices.Collect(Accumulate(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle mixed positive and negative", func() {
			input := slices.Values([]int{1, -2, 3, -4})
			want := []int{1, -1, 2, -2}

			result := slices.Collect(Accumulate(input))
			So(result, ShouldResemble, want)
		})
	})
}

func TestAccumulateBy(t *testing.T) {
	Convey("AccumulateBy", t, func() {
		Convey("Should accumulate with multiplication", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			multiply := func(acc, v int) int { return acc * v }
			want := []int{1, 2, 6, 24, 120}

			result := slices.Collect(AccumulateBy(input, multiply))
			So(result, ShouldResemble, want)
		})

		Convey("Should accumulate with custom function", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			custom := func(acc, v int) int { return acc*2 + v }
			want := []int{1, 4, 11, 26}

			result := slices.Collect(AccumulateBy(input, custom))
			So(result, ShouldResemble, want)
		})

		Convey("Should accumulate with subtraction", func() {
			input := slices.Values([]int{10, 3, 2, 1})
			subtract := func(acc, v int) int { return acc - v }
			want := []int{10, 7, 5, 4}

			result := slices.Collect(AccumulateBy(input, subtract))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{100})
			multiply := func(acc, v int) int { return acc * v }
			want := []int{100}

			result := slices.Collect(AccumulateBy(input, multiply))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			multiply := func(acc, v int) int { return acc * v }

			result := slices.Collect(AccumulateBy(input, multiply))
			So(result, ShouldBeNil)
		})

		Convey("Should work with different numeric types", func() {
			input := slices.Values([]float64{1.1, 2.2, 3.3})
			add := func(acc, v float64) float64 { return acc + v }
			want := []float64{1.1, 3.3, 6.6}

			result := slices.Collect(AccumulateBy(input, add))
			So(len(result), ShouldEqual, len(want))
			for i := range result {
				So(result[i], ShouldAlmostEqual, want[i], 0.0001)
			}
		})
	})
}

func TestAccumulateByFunc(t *testing.T) {
	Convey("AccumulateByFunc", t, func() {
		Convey("Should create function that accumulates with multiplication", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			multiplyFunc := AccumulateByFunc(func(acc, v int) int { return acc * v })
			want := []int{1, 2, 6, 24}

			result := slices.Collect(multiplyFunc(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create function that accumulates with addition", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			addFunc := AccumulateByFunc(func(acc, v int) int { return acc + v })
			want := []int{1, 3, 6, 10}

			result := slices.Collect(addFunc(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			input1 := slices.Values([]int{1, 2, 3})
			input2 := slices.Values([]int{4, 5, 6})
			multiplyFunc := AccumulateByFunc(func(acc, v int) int { return acc * v })

			result1 := slices.Collect(multiplyFunc(input1))
			result2 := slices.Collect(multiplyFunc(input2))

			So(result1, ShouldResemble, []int{1, 2, 6})
			So(result2, ShouldResemble, []int{4, 20, 120})
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			customFunc := AccumulateByFunc(func(acc, v int) int { return acc + v*2 })
			want := []int{42}

			result := slices.Collect(customFunc(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			customFunc := AccumulateByFunc(func(acc, v int) int { return acc + v })

			result := slices.Collect(customFunc(input))
			So(result, ShouldBeNil)
		})
	})
}
