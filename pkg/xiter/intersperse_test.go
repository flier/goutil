//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	. "github.com/flier/goutil/pkg/xiter"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleIntersperse() {
	s := slices.Values([]string{"foo", "bar", "baz"})
	i := Intersperse(s, ",")

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}

func ExampleIntersperseFunc() {
	sep := IntersperseFunc(",")

	s := slices.Values([]string{"foo", "bar", "baz"})
	i := sep(s)

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}

func ExampleIntersperseWith() {
	s := slices.Values([]string{"foo", "bar", "baz"})
	i := IntersperseWith(s, func() string { return "," })

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}

func ExampleIntersperseWithFunc() {
	sep := IntersperseWithFunc(func() string { return "," })

	s := slices.Values([]string{"foo", "bar", "baz"})
	i := sep(s)

	fmt.Println(strings.Join(slices.Collect(i), ""))
	// Output: foo,bar,baz
}

func TestIntersperse(t *testing.T) {
	Convey("Intersperse", t, func() {
		Convey("Should insert separator between elements", func() {
			input := slices.Values([]string{"foo", "bar", "baz"})
			separator := ","
			want := []string{"foo", ",", "bar", ",", "baz"}

			result := slices.Collect(Intersperse(input, separator))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]string{})
			separator := ","

			result := slices.Collect(Intersperse(input, separator))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]string{"foo"})
			separator := ","
			want := []string{"foo"}

			result := slices.Collect(Intersperse(input, separator))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle two elements", func() {
			input := slices.Values([]string{"foo", "bar"})
			separator := ","
			want := []string{"foo", ",", "bar"}

			result := slices.Collect(Intersperse(input, separator))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			separator := 0
			want := []int{1, 0, 2, 0, 3, 0, 4}

			result := slices.Collect(Intersperse(input, separator))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with custom struct types", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 35},
			})
			separator := Person{"---", 0}
			want := []Person{
				{"Alice", 25}, {"---", 0}, {"Bob", 30}, {"---", 0}, {"Charlie", 35},
			}

			result := slices.Collect(Intersperse(input, separator))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination at value", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			separator := 0
			want := []int{1, 0, 2, 0, 3}

			var result []int
			count := 0

			for v := range Intersperse(input, separator) {
				if count == 5 {
					break
				}
				result = append(result, v)
				count++
			}

			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination at separator", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			separator := 0
			want := []int{1, 0, 2, 0}

			var result []int
			count := 0

			for v := range Intersperse(input, separator) {
				if count == 4 {
					break
				}
				result = append(result, v)
				count++
			}

			So(result, ShouldResemble, want)
		})

		Convey("Should preserve sequence order", func() {
			input := slices.Values([]int{5, 4, 3, 2, 1})
			separator := 0
			want := []int{5, 0, 4, 0, 3, 0, 2, 0, 1}

			result := slices.Collect(Intersperse(input, separator))
			So(result, ShouldResemble, want)
		})
	})
}

func TestIntersperseFunc(t *testing.T) {
	Convey("IntersperseFunc", t, func() {
		Convey("Should create reusable intersperse function", func() {
			separator := IntersperseFunc[string](",")

			input1 := slices.Values([]string{"foo", "bar"})
			input2 := slices.Values([]string{"hello", "world", "test"})

			result1 := slices.Collect(separator(input1))
			result2 := slices.Collect(separator(input2))

			So(result1, ShouldResemble, []string{"foo", ",", "bar"})
			So(result2, ShouldResemble, []string{"hello", ",", "world", ",", "test"})
		})

		Convey("Should work with different types", func() {
			separator := IntersperseFunc[int](0)

			input := slices.Values([]int{1, 2, 3, 4})
			want := []int{1, 0, 2, 0, 3, 0, 4}

			result := slices.Collect(separator(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle edge cases consistently", func() {
			separator := IntersperseFunc[string](",")

			// Empty sequence
			emptyInput := slices.Values([]string{})
			emptyResult := slices.Collect(separator(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element
			singleInput := slices.Values([]string{"foo"})
			singleResult := slices.Collect(separator(singleInput))
			So(singleResult, ShouldResemble, []string{"foo"})

			// Two elements
			twoInput := slices.Values([]string{"foo", "bar"})
			twoResult := slices.Collect(separator(twoInput))
			So(twoResult, ShouldResemble, []string{"foo", ",", "bar"})
		})
	})
}

func TestIntersperseWith(t *testing.T) {
	Convey("IntersperseWith", t, func() {
		Convey("Should insert dynamically generated separator between elements", func() {
			input := slices.Values([]string{"foo", "bar", "baz"})
			separatorFunc := func() string { return "," }
			want := []string{"foo", ",", "bar", ",", "baz"}

			result := slices.Collect(IntersperseWith(input, separatorFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]string{})
			separatorFunc := func() string { return "," }

			result := slices.Collect(IntersperseWith(input, separatorFunc))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]string{"foo"})
			separatorFunc := func() string { return "," }
			want := []string{"foo"}

			result := slices.Collect(IntersperseWith(input, separatorFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			separatorFunc := func() int { return 0 }
			want := []int{1, 0, 2, 0, 3, 0, 4}

			result := slices.Collect(IntersperseWith(input, separatorFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should call separator function for each gap", func() {
			input := slices.Values([]int{1, 2, 3})
			counter := 0
			separatorFunc := func() int {
				counter++
				return counter * 10
			}
			want := []int{1, 10, 2, 20, 3}

			result := slices.Collect(IntersperseWith(input, separatorFunc))
			So(result, ShouldResemble, want)
			So(counter, ShouldEqual, 2) // Should be called twice for 3 elements
		})

		Convey("Should work with custom struct types", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 35},
			})
			separatorFunc := func() Person { return Person{"---", 0} }
			want := []Person{
				{"Alice", 25}, {"---", 0}, {"Bob", 30}, {"---", 0}, {"Charlie", 35},
			}

			result := slices.Collect(IntersperseWith(input, separatorFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			separatorFunc := func() int { return 0 }
			want := []int{1, 0, 2, 0, 3}

			var result []int
			count := 0

			for v := range IntersperseWith(input, separatorFunc) {
				if count == 5 {
					break
				}
				result = append(result, v)
				count++
			}

			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination at value", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			separatorFunc := func() int { return 0 }
			want := []int{1, 0, 2, 0, 3}

			var result []int
			count := 0

			for v := range IntersperseWith(input, separatorFunc) {
				if count == 5 {
					break
				}
				result = append(result, v)
				count++
			}

			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination at separator", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			separatorFunc := func() int { return 0 }
			want := []int{1, 0, 2, 0}

			var result []int
			count := 0

			for v := range IntersperseWith(input, separatorFunc) {
				if count == 4 {
					break
				}
				result = append(result, v)
				count++
			}

			So(result, ShouldResemble, want)
		})
	})
}

func TestIntersperseWithFunc(t *testing.T) {
	Convey("IntersperseWithFunc", t, func() {
		Convey("Should create reusable dynamic intersperse function", func() {
			separator := IntersperseWithFunc(func() string { return "," })

			input1 := slices.Values([]string{"foo", "bar"})
			input2 := slices.Values([]string{"hello", "world", "test"})

			result1 := slices.Collect(separator(input1))
			result2 := slices.Collect(separator(input2))

			So(result1, ShouldResemble, []string{"foo", ",", "bar"})
			So(result2, ShouldResemble, []string{"hello", ",", "world", ",", "test"})
		})

		Convey("Should work with different types", func() {
			separator := IntersperseWithFunc(func() int { return 0 })

			input := slices.Values([]int{1, 2, 3, 4})
			want := []int{1, 0, 2, 0, 3, 0, 4}

			result := slices.Collect(separator(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should call separator function for each gap", func() {
			counter := 0
			separator := IntersperseWithFunc(func() int {
				counter++
				return counter * 10
			})

			input := slices.Values([]int{1, 2, 3})
			want := []int{1, 10, 2, 20, 3}

			result := slices.Collect(separator(input))
			So(result, ShouldResemble, want)
			So(counter, ShouldEqual, 2) // Should be called twice for 3 elements
		})

		Convey("Should handle edge cases consistently", func() {
			separator := IntersperseWithFunc(func() string { return "," })

			// Empty sequence
			emptyInput := slices.Values([]string{})
			emptyResult := slices.Collect(separator(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element
			singleInput := slices.Values([]string{"foo"})
			singleResult := slices.Collect(separator(singleInput))
			So(singleResult, ShouldResemble, []string{"foo"})

			// Two elements
			twoInput := slices.Values([]string{"foo", "bar"})
			twoResult := slices.Collect(separator(twoInput))
			So(twoResult, ShouldResemble, []string{"foo", ",", "bar"})
		})

		Convey("Should work with complex separator logic", func() {
			separator := IntersperseWithFunc(func() string {
				return "---"
			})

			input := slices.Values([]string{"a", "b", "c", "d"})
			want := []string{"a", "---", "b", "---", "c", "---", "d"}

			result := slices.Collect(separator(input))
			So(result, ShouldResemble, want)
		})
	})
}
