//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
	. "github.com/flier/goutil/pkg/xiter/inspect"
)

func ExampleInspect() {
	s := RangeTo(5)

	square := MapFunc(func(n int) int { return n * n })
	dump := InspectFunc[int](Label("square"))
	s = dump(square(s))

	add2 := MapFunc(func(n int) int { return n + 2 })
	dump = InspectFunc[int](Label("add2"))
	s = dump(add2(s))

	fmt.Println(slices.Collect(s))

	s = RangeTo(20)
	dump = InspectFunc[int](Width(20), Limit(15))
	n := Sum(dump(s))
	fmt.Println(n)

	// Output:
	// square: [0 1 4 9 16]
	// add2: [2 3 6 11 18]
	// [2 3 6 11 18]
	// [0 1 2 3 4 5 6 7 8 9
	//  10 11 12 13 14 ...]
	// 190
}

func ExampleInspectFunc() {
	s := RangeTo(5)

	square := MapFunc(func(n int) int { return n * n })
	dump1 := InspectFunc[int](Label("square"))

	add2 := MapFunc(func(n int) int { return n + 2 })
	dump2 := InspectFunc[int](Label("add2"))

	s = dump2(add2(dump1(square(s))))

	fmt.Println(slices.Collect(s))

	s = RangeTo(20)
	dump3 := InspectFunc[int](Width(20), Limit(15))
	n := Sum(dump3(s))
	fmt.Println(n)

	// Output:
	// square: [0 1 4 9 16]
	// add2: [2 3 6 11 18]
	// [2 3 6 11 18]
	// [0 1 2 3 4 5 6 7 8 9
	//  10 11 12 13 14 ...]
	// 190
}

func ExampleInspect2() {
	s := slices.All([]string{"foo", "bar", "hello"})

	lengthOf := Map2Func(func(n int, k string) int { return len(k) })

	fmt.Println(maps.Collect(lengthOf(Inspect2(s, Label("len")))))

	// Output:
	// len: [0:foo 1:bar 2:hello]
	// map[0:3 1:3 2:5]
}

func ExampleInspect2Func() {
	s := slices.All([]string{"foo", "bar", "hello"})

	lengthOf := Map2Func(func(n int, k string) int { return len(k) })
	dump := Inspect2Func[int, string](Label("len"))

	fmt.Println(maps.Collect(lengthOf(dump(s))))

	// Output:
	// len: [0:foo 1:bar 2:hello]
	// map[0:3 1:3 2:5]
}

func TestInspect(t *testing.T) {
	Convey("Inspect", t, func() {
		Convey("Should inspect elements without modifying sequence", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 2, 3, 4, 5}
			buf := new(strings.Builder)

			result := slices.Collect(Inspect(input, Output(buf)))
			So(result, ShouldResemble, want)
			So(buf.String(), ShouldResemble, "[1 2 3 4 5]\n")
		})

		Convey("Should work with empty sequence", func() {
			input := slices.Values([]int{})
			buf := new(strings.Builder)

			result := slices.Collect(Inspect(input, Output(buf)))
			So(result, ShouldBeEmpty)
			So(buf.String(), ShouldEqual, "[]\n")
		})

		Convey("Should work with single element", func() {
			input := slices.Values([]int{42})
			buf := new(strings.Builder)

			result := slices.Collect(Inspect(input, Output(buf)))
			So(result, ShouldResemble, []int{42})
			So(buf.String(), ShouldResemble, "[42]\n")
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "b", "c"})

			buf := new(strings.Builder)

			result := slices.Collect(Inspect(input, Output(buf)))
			So(result, ShouldResemble, []string{"a", "b", "c"})
			So(buf.String(), ShouldResemble, "[a b c]\n")
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 2, 3}

			var result []int
			count := 0

			buf := new(strings.Builder)

			for v := range Inspect(input, Output(buf)) {
				if count == 3 {
					break
				}
				result = append(result, v)
				count++
			}

			So(result, ShouldResemble, want)
			So(buf.String(), ShouldResemble, "[1 2 3]\n")
		})

		Convey("Should work with custom options", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})

			buf := new(strings.Builder)

			result := slices.Collect(Inspect(input, Label("test"), Output(buf)))
			So(result, ShouldResemble, []int{1, 2, 3, 4, 5})
			So(buf.String(), ShouldResemble, "test: [1 2 3 4 5]\n")
		})

		Convey("Should preserve sequence order", func() {
			input := slices.Values([]int{5, 4, 3, 2, 1})

			buf := new(strings.Builder)

			result := slices.Collect(Inspect(input, Output(buf)))
			So(result, ShouldResemble, []int{5, 4, 3, 2, 1})
			So(buf.String(), ShouldResemble, "[5 4 3 2 1]\n")
		})
	})
}

func TestInspectFunc(t *testing.T) {
	Convey("InspectFunc", t, func() {
		Convey("Should create reusable inspection function", func() {
			buf := new(strings.Builder)
			inspect := InspectFunc[int](Label("test"), Output(buf))

			input1 := slices.Values([]int{1, 2, 3})
			input2 := slices.Values([]int{4, 5, 6})

			result1 := slices.Collect(inspect(input1))
			result2 := slices.Collect(inspect(input2))

			So(result1, ShouldResemble, []int{1, 2, 3})
			So(result2, ShouldResemble, []int{4, 5, 6})
			So(buf.String(), ShouldResemble, "test: [1 2 3]\ntest: [4 5 6]\n")
		})

		Convey("Should work with different types", func() {
			type Person struct {
				Name string
				Age  int
			}

			buf := new(strings.Builder)
			inspect := InspectFunc[Person](Label("person"), Output(buf))

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 35},
			})
			want := []Person{{"Alice", 25}, {"Bob", 30}, {"Charlie", 35}}

			result := slices.Collect(inspect(input))
			So(result, ShouldResemble, want)
			So(buf.String(), ShouldResemble, "person: [{Alice 25} {Bob 30} {Charlie 35}]\n")
		})

		Convey("Should work with multiple options", func() {
			buf := new(strings.Builder)
			inspect := InspectFunc[int](Label("test"), Width(10), Limit(3), Output(buf))

			input := slices.Values([]int{1, 2, 3, 4, 5})

			result := slices.Collect(inspect(input))
			So(result, ShouldResemble, []int{1, 2, 3, 4, 5})
			So(buf.String(), ShouldResemble, "test: [1 2 3 ...]\n")
		})

		Convey("Should handle edge cases consistently", func() {
			buf := new(strings.Builder)
			inspect := InspectFunc[int](Label("test"), Output(buf))

			// Empty sequence
			emptyInput := slices.Values([]int{})
			emptyResult := slices.Collect(inspect(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element
			singleInput := slices.Values([]int{42})
			singleResult := slices.Collect(inspect(singleInput))
			So(singleResult, ShouldResemble, []int{42})

			// No options
			noOptsInspect := InspectFunc[int]()
			noOptsResult := slices.Collect(noOptsInspect(singleInput))
			So(noOptsResult, ShouldResemble, []int{42})
			So(buf.String(), ShouldResemble, "test: []\ntest: [42]\n")
		})
	})
}

func TestInspect2(t *testing.T) {
	Convey("Inspect2", t, func() {
		Convey("Should inspect key-value pairs without modifying sequence", func() {
			keys := slices.Values([]int{1, 2, 3})
			values := slices.Values([]string{"a", "b", "c"})
			input := Zip(keys, values)

			want := map[int]string{1: "a", 2: "b", 3: "c"}

			buf := new(strings.Builder)
			result := maps.Collect(Inspect2(input, Output(buf)))
			So(result, ShouldResemble, want)
			So(buf.String(), ShouldResemble, "[1:a 2:b 3:c]\n")
		})

		Convey("Should work with empty sequence", func() {
			input := Zip(slices.Values([]int{}), slices.Values([]string{}))

			buf := new(strings.Builder)
			result := maps.Collect(Inspect2(input, Output(buf)))
			So(result, ShouldBeEmpty)
			So(buf.String(), ShouldEqual, "[]\n")
		})

		Convey("Should work with single element", func() {
			keys := slices.Values([]int{42})
			values := slices.Values([]string{"answer"})
			input := Zip(keys, values)

			want := map[int]string{42: "answer"}

			buf := new(strings.Builder)
			result := maps.Collect(Inspect2(input, Output(buf)))
			So(result, ShouldResemble, want)
			So(buf.String(), ShouldResemble, "[42:answer]\n")
		})

		Convey("Should work with different types", func() {
			keys := slices.Values([]string{"a", "b", "c"})
			values := slices.Values([]int{1, 2, 3})
			input := Zip(keys, values)

			want := map[string]int{"a": 1, "b": 2, "c": 3}

			buf := new(strings.Builder)
			result := maps.Collect(Inspect2(input, Output(buf)))
			So(result, ShouldResemble, want)
			So(buf.String(), ShouldResemble, "[a:1 b:2 c:3]\n")
		})

		Convey("Should handle early termination", func() {
			keys := slices.Values([]int{1, 2, 3, 4, 5})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)

			var result map[int]string
			count := 0

			buf := new(strings.Builder)
			for k, v := range Inspect2(input, Output(buf)) {
				if count == 3 {
					break
				}
				if result == nil {
					result = make(map[int]string)
				}
				result[k] = v
				count++
			}

			expected := map[int]string{1: "a", 2: "b", 3: "c"}
			So(result, ShouldResemble, expected)
			So(buf.String(), ShouldResemble, "[1:a 2:b 3:c]\n")
		})

		Convey("Should work with custom options", func() {
			keys := slices.Values([]int{1, 2, 3})
			values := slices.Values([]string{"a", "b", "c"})
			input := Zip(keys, values)

			buf := new(strings.Builder)
			result := maps.Collect(Inspect2(input, Label("test"), Output(buf)))
			So(result, ShouldResemble, map[int]string{1: "a", 2: "b", 3: "c"})
			So(buf.String(), ShouldResemble, "test: [1:a 2:b 3:c]\n")
		})

		Convey("Should preserve sequence order", func() {
			keys := slices.Values([]int{3, 2, 1})
			values := slices.Values([]string{"c", "b", "a"})
			input := Zip(keys, values)

			buf := new(strings.Builder)
			result := maps.Collect(Inspect2(input, Output(buf)))
			So(result, ShouldResemble, map[int]string{3: "c", 2: "b", 1: "a"})
			So(buf.String(), ShouldResemble, "[3:c 2:b 1:a]\n")
		})
	})
}

func TestInspect2Func(t *testing.T) {
	Convey("Inspect2Func", t, func() {
		Convey("Should create reusable key-value inspection function", func() {
			buf := new(strings.Builder)
			inspect := Inspect2Func[int, string](Label("test"), Output(buf))

			keys1 := slices.Values([]int{1, 2, 3})
			values1 := slices.Values([]string{"a", "b", "c"})
			input1 := Zip(keys1, values1)

			keys2 := slices.Values([]int{4, 5, 6})
			values2 := slices.Values([]string{"d", "e", "f"})
			input2 := Zip(keys2, values2)

			result1 := maps.Collect(inspect(input1))
			result2 := maps.Collect(inspect(input2))

			So(result1, ShouldResemble, map[int]string{1: "a", 2: "b", 3: "c"})
			So(result2, ShouldResemble, map[int]string{4: "d", 5: "e", 6: "f"})
			So(buf.String(), ShouldResemble, "test: [1:a 2:b 3:c]\ntest: [4:d 5:e 6:f]\n")
		})

		Convey("Should work with different types", func() {
			type Person struct {
				Name string
				Age  int
			}

			buf := new(strings.Builder)
			inspect := Inspect2Func[string, Person](Label("person"), Output(buf))

			keys := slices.Values([]string{"alice", "bob", "charlie"})
			values := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 30}, {"Charlie", 35},
			})
			input := Zip(keys, values)

			want := map[string]Person{
				"alice":   {"Alice", 25},
				"bob":     {"Bob", 30},
				"charlie": {"Charlie", 35},
			}

			result := maps.Collect(inspect(input))
			So(result, ShouldResemble, want)
			So(buf.String(), ShouldResemble, "person: [alice:{Alice 25} bob:{Bob 30} charlie:{Charlie 35}]\n")
		})

		Convey("Should work with multiple options", func() {
			buf := new(strings.Builder)
			inspect := Inspect2Func[int, string](Label("test"), Width(10), Limit(3), Output(buf))

			keys := slices.Values([]int{1, 2, 3, 4, 5})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)

			result := maps.Collect(inspect(input))
			So(result, ShouldResemble, map[int]string{1: "a", 2: "b", 3: "c", 4: "d", 5: "e"})
			So(buf.String(), ShouldResemble, "test: [1:a 2:b\n 3:c ...]\n")
		})

		Convey("Should handle edge cases consistently", func() {
			buf := new(strings.Builder)
			inspect := Inspect2Func[int, string](Label("test"), Output(buf))

			// Empty sequence
			emptyInput := Zip(slices.Values([]int{}), slices.Values([]string{}))
			emptyResult := maps.Collect(inspect(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element
			singleKeys := slices.Values([]int{42})
			singleValues := slices.Values([]string{"answer"})
			singleInput := Zip(singleKeys, singleValues)
			singleResult := maps.Collect(inspect(singleInput))
			So(singleResult, ShouldResemble, map[int]string{42: "answer"})

			// No options
			noOptsInspect := Inspect2Func[int, string](Output(buf))
			keys := slices.Values([]int{1, 2, 3})
			values := slices.Values([]string{"a", "b", "c"})
			input := Zip(keys, values)
			noOptsResult := maps.Collect(noOptsInspect(input))
			So(noOptsResult, ShouldResemble, map[int]string{1: "a", 2: "b", 3: "c"})
			So(buf.String(), ShouldResemble, "test: []\ntest: [42:answer]\n[1:a 2:b 3:c]\n")
		})
	})
}
