//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

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

func TestTake(t *testing.T) {
	Convey("Take", t, func() {
		Convey("Should take first n elements", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 2, 3}

			result := slices.Collect(Take(input, 3))
			So(result, ShouldResemble, want)
		})

		Convey("Should take all elements when n is larger than sequence", func() {
			input := slices.Values([]int{1, 2, 3})
			want := []int{1, 2, 3}

			result := slices.Collect(Take(input, 5))
			So(result, ShouldResemble, want)
		})

		Convey("Should return empty when n is 0", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})

			result := slices.Collect(Take(input, 0))
			So(result, ShouldBeEmpty)
		})

		Convey("Should return empty when n is negative", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})

			result := slices.Collect(Take(input, -1))
			// For negative n, the function will return empty since i >= n is always true
			// for any positive i, causing the loop to exit immediately
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})

			result := slices.Collect(Take(input, 3))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})

			result := slices.Collect(Take(input, 1))
			So(result, ShouldResemble, []int{42})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})

			seq := Take(input, 3)
			result := make([]int, 0)
			count := 0
			for v := range seq {
				result = append(result, v)
				count++
				if count >= 2 { // Early termination
					break
				}
			}

			So(result, ShouldResemble, []int{1, 2})
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "bb", "ccc", "dddd"})
			want := []string{"a", "bb"}

			result := slices.Collect(Take(input, 2))
			So(result, ShouldResemble, want)
		})
	})
}

func TestTakeFunc(t *testing.T) {
	Convey("TakeFunc", t, func() {
		Convey("Should create function that takes first n elements", func() {
			take2 := TakeFunc[int](2)

			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 2}

			result := slices.Collect(take2(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			take3 := TakeFunc[string](3)

			input1 := slices.Values([]string{"a", "b", "c", "d"})
			input2 := slices.Values([]string{"x", "y", "z", "w"})

			result1 := slices.Collect(take3(input1))
			result2 := slices.Collect(take3(input2))

			So(result1, ShouldResemble, []string{"a", "b", "c"})
			So(result2, ShouldResemble, []string{"x", "y", "z"})
		})

		Convey("Should work with different types", func() {
			take1 := TakeFunc[float64](1)

			input := slices.Values([]float64{1.1, 2.2, 3.3})
			want := []float64{1.1}

			result := slices.Collect(take1(input))
			So(result, ShouldResemble, want)
		})
	})
}

func TestTake2(t *testing.T) {
	Convey("Take2", t, func() {
		Convey("Should take first n key-value pairs", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})
			want := map[int]string{0: "foo", 1: "bar", 2: "hello"}

			result := maps.Collect(Take2(input, 3))
			So(result, ShouldResemble, want)
		})

		Convey("Should take all key-value pairs when n is larger than sequence", func() {
			input := slices.All([]string{"foo", "bar"})
			want := map[int]string{0: "foo", 1: "bar"}

			result := maps.Collect(Take2(input, 5))
			So(result, ShouldResemble, want)
		})

		Convey("Should return empty when n is 0", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})

			result := maps.Collect(Take2(input, 0))
			So(result, ShouldBeEmpty)
		})

		Convey("Should return empty when n is negative", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})

			result := maps.Collect(Take2(input, -1))
			// For negative n, the function will return empty since i >= n is always true
			// for any positive i, causing the loop to exit immediately
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.All([]string{})

			result := maps.Collect(Take2(input, 3))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single key-value pair", func() {
			input := slices.All([]string{"single"})

			result := maps.Collect(Take2(input, 1))
			So(result, ShouldResemble, map[int]string{0: "single"})
		})

		Convey("Should handle early termination", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})

			seq := Take2(input, 3)
			result := make(map[int]string)
			count := 0
			for k, v := range seq {
				result[k] = v
				count++
				if count >= 2 { // Early termination
					break
				}
			}

			So(len(result), ShouldEqual, 2)
			So(result[0], ShouldEqual, "foo")
			So(result[1], ShouldEqual, "bar")
		})

		Convey("Should work with different types", func() {
			input := slices.All([]string{"a", "bb", "ccc"})

			result := maps.Collect(Take2(input, 2))
			// Since slice iteration order is guaranteed, we can check the exact result
			So(result, ShouldResemble, map[int]string{0: "a", 1: "bb"})
		})
	})
}

func TestTake2Func(t *testing.T) {
	Convey("Take2Func", t, func() {
		Convey("Should create function that takes first n key-value pairs", func() {
			take2 := Take2Func[int, string](2)

			input := slices.All([]string{"foo", "bar", "hello", "world"})
			want := map[int]string{0: "foo", 1: "bar"}

			result := maps.Collect(take2(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			take1 := Take2Func[int, string](1)

			input1 := slices.All([]string{"a", "b", "c"})
			input2 := slices.All([]string{"x", "y", "z"})

			result1 := maps.Collect(take1(input1))
			result2 := maps.Collect(take1(input2))

			So(len(result1), ShouldEqual, 1)
			So(len(result2), ShouldEqual, 1)
		})

		Convey("Should work with different types", func() {
			take3 := Take2Func[int, float64](3)

			input := slices.All([]float64{1.1, 2.2, 3.3, 4.4})

			result := maps.Collect(take3(input))
			// Since slice iteration order is guaranteed, we can check the exact result
			So(result, ShouldResemble, map[int]float64{0: 1.1, 1: 2.2, 2: 3.3})
		})
	})
}

func TestTakeWhile(t *testing.T) {
	Convey("TakeWhile", t, func() {
		Convey("Should take elements while predicate is true", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			predicate := func(n int) bool { return n <= 3 }
			want := []int{1, 2, 3}

			result := slices.Collect(TakeWhile(input, predicate))
			So(result, ShouldResemble, want)
		})

		Convey("Should take no elements when predicate is false from start", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			predicate := func(n int) bool { return n > 10 }

			result := slices.Collect(TakeWhile(input, predicate))
			So(result, ShouldBeEmpty)
		})

		Convey("Should take all elements when predicate is always true", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			predicate := func(n int) bool { return n > 0 }
			want := []int{1, 2, 3, 4, 5}

			result := slices.Collect(TakeWhile(input, predicate))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			predicate := func(n int) bool { return n > 0 }

			result := slices.Collect(TakeWhile(input, predicate))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			predicate := func(n int) bool { return n > 40 }

			result := slices.Collect(TakeWhile(input, predicate))
			So(result, ShouldResemble, []int{42})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			predicate := func(n int) bool { return n <= 3 }

			seq := TakeWhile(input, predicate)
			result := make([]int, 0)
			count := 0
			for v := range seq {
				result = append(result, v)
				count++
				if count >= 2 { // Early termination
					break
				}
			}

			So(result, ShouldResemble, []int{1, 2})
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "bb", "ccc", "dddd"})
			predicate := func(s string) bool { return len(s) <= 2 }
			want := []string{"a", "bb"}

			result := slices.Collect(TakeWhile(input, predicate))
			So(result, ShouldResemble, want)
		})
	})
}

func TestTakeWhileFunc(t *testing.T) {
	Convey("TakeWhileFunc", t, func() {
		Convey("Should create function that takes elements while predicate is true", func() {
			takeLe3 := TakeWhileFunc(func(n int) bool { return n <= 3 })

			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 2, 3}

			result := slices.Collect(takeLe3(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			takeShort := TakeWhileFunc(func(s string) bool { return len(s) <= 3 })

			input1 := slices.Values([]string{"a", "bb", "ccc", "dddd"})
			input2 := slices.Values([]string{"x", "yy", "zzz", "wwww"})

			result1 := slices.Collect(takeShort(input1))
			result2 := slices.Collect(takeShort(input2))

			So(result1, ShouldResemble, []string{"a", "bb", "ccc"})
			So(result2, ShouldResemble, []string{"x", "yy", "zzz"})
		})

		Convey("Should work with different types", func() {
			takePositive := TakeWhileFunc(func(n float64) bool { return n > 0.0 })

			input := slices.Values([]float64{1.1, 2.2, -3.3, 4.4})
			want := []float64{1.1, 2.2}

			result := slices.Collect(takePositive(input))
			So(result, ShouldResemble, want)
		})
	})
}

func TestTakeWhile2(t *testing.T) {
	Convey("TakeWhile2", t, func() {
		Convey("Should take key-value pairs while predicate is true", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})
			predicate := func(i int, v string) bool { return len(v) <= 3 }
			want := map[int]string{0: "foo", 1: "bar"}

			result := maps.Collect(TakeWhile2(input, predicate))
			So(result, ShouldResemble, want)
		})

		Convey("Should take no key-value pairs when predicate is false from start", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})
			predicate := func(i int, v string) bool { return len(v) > 10 }

			result := maps.Collect(TakeWhile2(input, predicate))
			So(result, ShouldBeEmpty)
		})

		Convey("Should take all key-value pairs when predicate is always true", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})
			predicate := func(i int, v string) bool { return len(v) > 0 }
			want := map[int]string{0: "foo", 1: "bar", 2: "hello", 3: "world"}

			result := maps.Collect(TakeWhile2(input, predicate))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.All([]string{})
			predicate := func(i int, v string) bool { return len(v) > 0 }

			result := maps.Collect(TakeWhile2(input, predicate))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single key-value pair", func() {
			input := slices.All([]string{"single"})
			predicate := func(i int, v string) bool { return len(v) > 5 }

			result := maps.Collect(TakeWhile2(input, predicate))
			So(result, ShouldResemble, map[int]string{0: "single"})
		})

		Convey("Should handle early termination", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})
			predicate := func(i int, v string) bool { return len(v) <= 3 }

			seq := TakeWhile2(input, predicate)
			result := make(map[int]string)
			count := 0
			for k, v := range seq {
				result[k] = v
				count++
				if count >= 1 { // Early termination
					break
				}
			}

			So(len(result), ShouldEqual, 1)
			So(result[0], ShouldEqual, "foo")
		})

		Convey("Should work with different types", func() {
			input := slices.All([]string{"a", "bb", "ccc", "dddd"})
			predicate := func(i int, v string) bool { return len(v) <= 2 }

			result := maps.Collect(TakeWhile2(input, predicate))
			// Since slice iteration order is guaranteed, we can check the exact result
			So(len(result), ShouldEqual, 2)
			So(result[0], ShouldEqual, "a")
			So(result[1], ShouldEqual, "bb")
		})
	})
}

func TestTakeWhile2Func(t *testing.T) {
	Convey("TakeWhile2Func", t, func() {
		Convey("Should create function that takes key-value pairs while predicate is true", func() {
			takeShort := TakeWhile2Func(func(i int, v string) bool { return len(v) <= 3 })

			input := slices.All([]string{"foo", "bar", "hello", "world"})
			want := map[int]string{0: "foo", 1: "bar"}

			result := maps.Collect(takeShort(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			takeEvenIndex := TakeWhile2Func(func(i int, v string) bool { return i%2 == 0 })

			input1 := slices.All([]string{"a", "b", "c", "d"})
			input2 := slices.All([]string{"x", "y", "z", "w"})

			result1 := maps.Collect(takeEvenIndex(input1))
			result2 := maps.Collect(takeEvenIndex(input2))

			// TakeWhile2 stops at the first element that doesn't satisfy the predicate
			// So it will only take the first element (index 0)
			So(result1, ShouldResemble, map[int]string{0: "a"})
			So(result2, ShouldResemble, map[int]string{0: "x"})
		})

		Convey("Should work with different types", func() {
			takeSmallValue := TakeWhile2Func(func(i int, v string) bool { return len(v) <= 10 })

			input := slices.All([]string{"small", "tiny"})

			result := maps.Collect(takeSmallValue(input))
			// Since slice iteration order is guaranteed, we can check the exact result
			So(len(result), ShouldEqual, 2)
			So(result[0], ShouldEqual, "small")
			So(result[1], ShouldEqual, "tiny")
		})
	})
}
