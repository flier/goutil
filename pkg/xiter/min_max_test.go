//go:build go1.23

package xiter_test

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/flier/goutil/pkg/tuple"
	. "github.com/flier/goutil/pkg/xiter"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleMinMax() {
	s := slices.Values([]int{1, 2, 3, 4, 5})

	fmt.Println(MinMax(s))
	// Output: (1, 5)
}

func ExampleMinMaxBy() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	fmt.Println(MinMaxBy(s, cmp.Compare))
	// Output: (bar, world)
}

func ExampleMinMaxByFunc() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	f := MinMaxByFunc[string](cmp.Compare)

	fmt.Println(f(s))
	// Output: (bar, world)
}

func ExampleMinMaxByKey() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	fmt.Println(MinMaxByKey(s, func(s string) int { return len(s) }))
	// Output: (bar, world)
}

func ExampleMinMaxByKeyFunc() {
	s := slices.Values([]string{"foo", "bar", "hello", "world"})
	f := MinMaxByKeyFunc(func(s string) int { return len(s) })

	fmt.Println(f(s))
	// Output: (bar, world)
}

func TestMinMax(t *testing.T) {
	Convey("MinMax", t, func() {
		Convey("Should find min and max of integer sequence", func() {
			input := slices.Values([]int{3, 1, 4, 1, 5, 9, 2, 6})
			want := tuple.New2(1, 9)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should find min and max of string sequence", func() {
			input := slices.Values([]string{"foo", "bar", "hello", "world", "a"})
			want := tuple.New2("a", "world")

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			want := tuple.New2(0, 0)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			want := tuple.New2(42, 42)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle two elements", func() {
			input := slices.Values([]int{10, 5})
			want := tuple.New2(5, 10)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle duplicate min/max values", func() {
			input := slices.Values([]int{5, 1, 5, 1, 5})
			want := tuple.New2(1, 5)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with float64", func() {
			input := slices.Values([]float64{3.14, 2.71, 1.41, 2.23})
			want := tuple.New2(1.41, 3.14)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with negative numbers", func() {
			input := slices.Values([]int{-5, -10, -1, -100, -50})
			want := tuple.New2(-100, -1)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle all same values", func() {
			input := slices.Values([]int{7, 7, 7, 7, 7})
			want := tuple.New2(7, 7)

			result := MinMax(input)
			So(result, ShouldResemble, want)
		})
	})
}

func TestMinMaxBy(t *testing.T) {
	Convey("MinMaxBy", t, func() {
		Convey("Should find min and max using custom comparison", func() {
			input := slices.Values([]int{3, 1, 4, 1, 5, 9, 2, 6})
			compareFunc := func(a, b int) int { return a - b }
			want := tuple.New2(1, 9)

			result := MinMaxBy(input, compareFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should find min and max of strings using cmp.Compare", func() {
			input := slices.Values([]string{"foo", "bar", "hello", "world", "a"})
			want := tuple.New2("a", "world")

			result := MinMaxBy(input, cmp.Compare[string])
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			compareFunc := func(a, b int) int { return a - b }
			want := tuple.New2(0, 0)

			result := MinMaxBy(input, compareFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			compareFunc := func(a, b int) int { return a - b }
			want := tuple.New2(42, 42)

			result := MinMaxBy(input, compareFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle reverse comparison", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			reverseCompare := func(a, b int) int { return b - a }
			want := tuple.New2(5, 1) // Reverse order

			result := MinMaxBy(input, reverseCompare)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with custom struct comparison", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 20}, {"David", 35},
			})
			// Use MinMaxByKey instead since Person doesn't implement cmp.Ordered
			keyFunc := func(p Person) int { return p.Age }
			want := tuple.New2(Person{"Charlie", 20}, Person{"David", 35})

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle equal values in comparison", func() {
			input := slices.Values([]int{5, 5, 5, 5})
			compareFunc := func(a, b int) int { return a - b }
			want := tuple.New2(5, 5)

			result := MinMaxBy(input, compareFunc)
			So(result, ShouldResemble, want)
		})
	})
}

func TestMinMaxByFunc(t *testing.T) {
	Convey("MinMaxByFunc", t, func() {
		Convey("Should create reusable minmax function", func() {
			minmaxFunc := MinMaxByFunc[int](func(a, b int) int { return a - b })

			input1 := slices.Values([]int{3, 1, 4, 1, 5})
			input2 := slices.Values([]int{10, 5, 15, 2})

			result1 := minmaxFunc(input1)
			result2 := minmaxFunc(input2)

			So(result1, ShouldResemble, tuple.New2(1, 5))
			So(result2, ShouldResemble, tuple.New2(2, 15))
		})

		Convey("Should work with string comparison", func() {
			minmaxFunc := MinMaxByFunc[string](cmp.Compare[string])

			input := slices.Values([]string{"foo", "bar", "hello", "world"})
			want := tuple.New2("bar", "world")

			result := minmaxFunc(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle edge cases consistently", func() {
			minmaxFunc := MinMaxByFunc[int](func(a, b int) int { return a - b })

			// Empty sequence
			emptyInput := slices.Values([]int{})
			emptyResult := minmaxFunc(emptyInput)
			So(emptyResult, ShouldResemble, tuple.New2(0, 0))

			// Single element
			singleInput := slices.Values([]int{42})
			singleResult := minmaxFunc(singleInput)
			So(singleResult, ShouldResemble, tuple.New2(42, 42))

			// Two elements
			twoInput := slices.Values([]int{10, 5})
			twoResult := minmaxFunc(twoInput)
			So(twoResult, ShouldResemble, tuple.New2(5, 10))
		})

		Convey("Should work with custom comparison logic", func() {
			// Compare by absolute value
			absCompare := func(a, b int) int {
				absA := a
				if absA < 0 {
					absA = -absA
				}
				absB := b
				if absB < 0 {
					absB = -absB
				}
				return absA - absB
			}

			minmaxFunc := MinMaxByFunc[int](absCompare)

			input := slices.Values([]int{-5, 3, -10, 7, -1})
			want := tuple.New2(-1, -10) // -1 has smallest abs, -10 has largest abs

			result := minmaxFunc(input)
			So(result, ShouldResemble, want)
		})
	})
}

func TestMinMaxByKey(t *testing.T) {
	Convey("MinMaxByKey", t, func() {
		Convey("Should find min and max by key function", func() {
			input := slices.Values([]string{"foo", "bar", "hello", "world", "a"})
			keyFunc := func(s string) int { return len(s) }
			// "a" has length 1 (min), "world" has length 5 (max)
			want := tuple.New2("a", "world")

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			keyFunc := func(i int) int { return i * 2 }
			want := tuple.New2(0, 0)

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			keyFunc := func(i int) int { return i * 2 }
			want := tuple.New2(42, 42)

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with struct key extraction", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 20}, {"David", 35},
			})
			keyFunc := func(p Person) int { return p.Age }
			want := tuple.New2(Person{"Charlie", 20}, Person{"David", 35})

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with different key types", func() {
			input := slices.Values([]string{"apple", "banana", "cherry", "date"})
			keyFunc := func(s string) float64 { return float64(len(s)) * 0.5 }
			// "date" has length 4, key = 2.0 (min), "cherry" has length 6, key = 3.0 (max)
			// Since we iterate through apple(2.5), banana(3.0), cherry(3.0), date(2.0)
			// Min: date(2.0), Max: cherry(3.0) - cherry comes after banana with same key
			want := tuple.New2("date", "cherry")

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle duplicate key values", func() {
			input := slices.Values([]string{"a", "b", "c", "d"})
			keyFunc := func(s string) int { return 1 } // All have same key
			// When all keys are equal, should return first and last elements
			// Since we iterate through a,b,c,d, we get d for both min and max
			// (last element with equal key value)
			want := tuple.New2("d", "d")

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should work with negative key values", func() {
			input := slices.Values([]int{-5, -10, -1, -100, -50})
			keyFunc := func(i int) int { return -i } // Negate the key
			// -100 has key 100 (max), -1 has key 1 (min)
			want := tuple.New2(-1, -100)

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle complex key calculation", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			keyFunc := func(i int) int { return i*i + i } // i^2 + i
			want := tuple.New2(1, 5)                      // 1 has key 2, 5 has key 30

			result := MinMaxByKey(input, keyFunc)
			So(result, ShouldResemble, want)
		})
	})
}

func TestMinMaxByKeyFunc(t *testing.T) {
	Convey("MinMaxByKeyFunc", t, func() {
		Convey("Should create reusable key-based minmax function", func() {
			minmaxFunc := MinMaxByKeyFunc(func(s string) int { return len(s) })

			input1 := slices.Values([]string{"foo", "bar", "hello"})
			input2 := slices.Values([]string{"a", "world", "test", "x"})

			result1 := minmaxFunc(input1)
			result2 := minmaxFunc(input2)

			So(result1, ShouldResemble, tuple.New2("bar", "hello"))
			So(result2, ShouldResemble, tuple.New2("x", "world"))
		})

		Convey("Should work with different types", func() {
			minmaxFunc := MinMaxByKeyFunc(func(i int) float64 { return float64(i) * 0.5 })

			input := slices.Values([]int{2, 4, 6, 8, 10})
			want := tuple.New2(2, 10)

			result := minmaxFunc(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should handle edge cases consistently", func() {
			minmaxFunc := MinMaxByKeyFunc(func(s string) int { return len(s) })

			// Empty sequence
			emptyInput := slices.Values([]string{})
			emptyResult := minmaxFunc(emptyInput)
			So(emptyResult, ShouldResemble, tuple.New2("", ""))

			// Single element
			singleInput := slices.Values([]string{"hello"})
			singleResult := minmaxFunc(singleInput)
			So(singleResult, ShouldResemble, tuple.New2("hello", "hello"))

			// Two elements
			twoInput := slices.Values([]string{"a", "hello"})
			twoResult := minmaxFunc(twoInput)
			So(twoResult, ShouldResemble, tuple.New2("a", "hello"))
		})

		Convey("Should work with complex key functions", func() {
			// Key function that returns a combination of multiple properties
			type Item struct {
				Name  string
				Price float64
				Score int
			}

			input := slices.Values([]Item{
				{"A", 10.0, 5}, {"B", 20.0, 3}, {"C", 15.0, 4}, {"D", 25.0, 2},
			})

			// Complex key: price * score
			keyFunc := func(item Item) float64 { return item.Price * float64(item.Score) }
			minmaxFunc := MinMaxByKeyFunc(keyFunc)

			// A: 10*5=50, B: 20*3=60, C: 15*4=60, D: 25*2=50
			// Min: A or D (both have key 50), Max: B or C (both have key 60)
			// Since we iterate through A,B,C,D, we get D for min and C for max
			want := tuple.New2(Item{"D", 25.0, 2}, Item{"C", 15.0, 4})

			result := minmaxFunc(input)
			So(result, ShouldResemble, want)
		})

		Convey("Should preserve original element types", func() {
			// Test that the function returns the original elements, not the key values
			input := slices.Values([]string{"apple", "banana", "cherry"})
			keyFunc := func(s string) int { return len(s) }
			minmaxFunc := MinMaxByKeyFunc(keyFunc)

			result := minmaxFunc(input)
			// Should return strings, not ints
			So(result, ShouldResemble, tuple.New2("apple", "cherry"))
		})
	})
}
