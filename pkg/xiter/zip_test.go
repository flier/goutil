//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleZip() {
	s3 := Zip(slices.Values([]int{1, 2, 3}), slices.Values([]string{"foo", "bar", "hello"}))

	fmt.Println(maps.Collect(s3))

	// Output: map[1:foo 2:bar 3:hello]
}

func ExampleZipWith() {
	s1 := slices.Values([]int{1, 2, 3})
	s2 := slices.Values([]int{4, 5, 6})
	s3 := ZipWith(s1, s2, func(x, y int) int { return x + y })

	fmt.Println(slices.Collect(s3))

	// Output: [5 7 9]
}

func ExampleZipWithFunc() {
	zipAndAdd := ZipWithFunc(func(x, y int) int { return x + y })

	s1 := slices.Values([]int{1, 2, 3})
	s2 := slices.Values([]int{4, 5, 6})
	s3 := zipAndAdd(s1, s2)

	fmt.Println(slices.Collect(s3))

	// Output: [5 7 9]
}

func ExampleUnzip() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	k, v := Unzip(s).Unpack()

	keys := slices.Collect(k)
	values := slices.Collect(v)

	sort.Strings(keys)
	sort.Strings(values)

	fmt.Println(keys)
	fmt.Println(values)

	// Output:
	// [foo hello]
	// [bar world]
}

func TestZip(t *testing.T) {
	Convey("Zip", t, func() {
		Convey("Should zip two sequences of equal length", func() {
			keys := slices.Values([]int{1, 2, 3})
			values := slices.Values([]string{"foo", "bar", "hello"})
			want := map[int]string{1: "foo", 2: "bar", 3: "hello"}

			result := maps.Collect(Zip(keys, values))
			So(result, ShouldResemble, want)
		})

		Convey("Should zip sequences with different types", func() {
			keys := slices.Values([]string{"a", "b", "c"})
			values := slices.Values([]int{1, 2, 3})
			want := map[string]int{"a": 1, "b": 2, "c": 3}

			result := maps.Collect(Zip(keys, values))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element sequences", func() {
			keys := slices.Values([]int{42})
			values := slices.Values([]string{"answer"})
			want := map[int]string{42: "answer"}

			result := maps.Collect(Zip(keys, values))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequences", func() {
			keys := slices.Values([]int{})
			values := slices.Values([]string{})
			want := map[int]string{}

			result := maps.Collect(Zip(keys, values))
			So(result, ShouldResemble, want)
		})

		Convey("Should stop at shorter sequence length", func() {
			keys := slices.Values([]int{1, 2, 3, 4})
			values := slices.Values([]string{"foo", "bar"})
			want := map[int]string{1: "foo", 2: "bar"}

			result := maps.Collect(Zip(keys, values))
			So(result, ShouldResemble, want)
		})

		Convey("Should stop at shorter sequence length (keys shorter)", func() {
			keys := slices.Values([]int{1, 2})
			values := slices.Values([]string{"foo", "bar", "hello", "world"})
			want := map[int]string{1: "foo", 2: "bar"}

			result := maps.Collect(Zip(keys, values))
			So(result, ShouldResemble, want)
		})
	})
}

func TestZipWith(t *testing.T) {
	Convey("ZipWith", t, func() {
		Convey("Should zip and apply function to elements", func() {
			s1 := slices.Values([]int{1, 2, 3})
			s2 := slices.Values([]int{4, 5, 6})
			add := func(x, y int) int { return x + y }
			want := []int{5, 7, 9}

			result := slices.Collect(ZipWith(s1, s2, add))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with different types", func() {
			keys := slices.Values([]string{"a", "b", "c"})
			values := slices.Values([]int{1, 2, 3})
			combine := func(k string, v int) string { return fmt.Sprintf("%s%d", k, v) }
			want := []string{"a1", "b2", "c3"}

			result := slices.Collect(ZipWith(keys, values, combine))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element sequences", func() {
			s1 := slices.Values([]int{10})
			s2 := slices.Values([]int{20})
			multiply := func(x, y int) int { return x * y }
			want := []int{200}

			result := slices.Collect(ZipWith(s1, s2, multiply))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequences", func() {
			s1 := slices.Values([]int{})
			s2 := slices.Values([]int{})
			add := func(x, y int) int { return x + y }

			result := slices.Collect(ZipWith(s1, s2, add))
			So(result, ShouldBeNil)
		})

		Convey("Should stop at shorter sequence length", func() {
			s1 := slices.Values([]int{1, 2, 3, 4})
			s2 := slices.Values([]int{10, 20})
			add := func(x, y int) int { return x + y }
			want := []int{11, 22}

			result := slices.Collect(ZipWith(s1, s2, add))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with complex functions", func() {
			s1 := slices.Values([]int{1, 2, 3})
			s2 := slices.Values([]int{2, 3, 4})
			complex := func(x, y int) int { return x*x + y*y }
			want := []int{5, 13, 25}

			result := slices.Collect(ZipWith(s1, s2, complex))
			So(result, ShouldResemble, want)
		})
	})
}

func TestZipWithFunc(t *testing.T) {
	Convey("ZipWithFunc", t, func() {
		Convey("Should create function that zips and applies transformation", func() {
			s1 := slices.Values([]int{1, 2, 3})
			s2 := slices.Values([]int{4, 5, 6})
			zipAndAdd := ZipWithFunc(func(x, y int) int { return x + y })
			want := []int{5, 7, 9}

			result := slices.Collect(zipAndAdd(s1, s2))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			s1 := slices.Values([]int{1, 2})
			s2 := slices.Values([]int{3, 4})
			s3 := slices.Values([]int{5, 6})
			zipAndMultiply := ZipWithFunc(func(x, y int) int { return x * y })

			result1 := slices.Collect(zipAndMultiply(s1, s2))
			result2 := slices.Collect(zipAndMultiply(s2, s3))

			So(result1, ShouldResemble, []int{3, 8})
			So(result2, ShouldResemble, []int{15, 24})
		})

		Convey("Should work with different types", func() {
			keys := slices.Values([]string{"a", "b"})
			values := slices.Values([]int{1, 2})
			zipAndFormat := ZipWithFunc(func(k string, v int) string { return fmt.Sprintf("%s%d", k, v) })
			want := []string{"a1", "b2"}

			result := slices.Collect(zipAndFormat(keys, values))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle single element sequences", func() {
			s1 := slices.Values([]int{10})
			s2 := slices.Values([]int{20})
			zipAndSubtract := ZipWithFunc(func(x, y int) int { return x - y })
			want := []int{-10}

			result := slices.Collect(zipAndSubtract(s1, s2))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequences", func() {
			s1 := slices.Values([]int{})
			s2 := slices.Values([]int{})
			zipAndAdd := ZipWithFunc(func(x, y int) int { return x + y })

			result := slices.Collect(zipAndAdd(s1, s2))
			So(result, ShouldBeNil)
		})
	})
}

func TestUnzip(t *testing.T) {
	Convey("Unzip", t, func() {
		Convey("Should unzip key-value pairs into separate sequences", func() {
			input := maps.All(map[string]int{"foo": 1, "bar": 2, "hello": 3})
			keys, values := Unzip(input).Unpack()

			keysResult := slices.Collect(keys)
			valuesResult := slices.Collect(values)

			// Sort for consistent comparison since map iteration order is not guaranteed
			sort.Strings(keysResult)
			sort.Ints(valuesResult)

			So(keysResult, ShouldResemble, []string{"bar", "foo", "hello"})
			// The values order depends on the order keys are consumed
			// Since we sort keys, we need to check that all expected values are present
			So(len(valuesResult), ShouldEqual, 3)
			So(slices.Contains(valuesResult, 1), ShouldBeTrue)
			So(slices.Contains(valuesResult, 2), ShouldBeTrue)
			So(slices.Contains(valuesResult, 3), ShouldBeTrue)
		})

		Convey("Should handle single key-value pair", func() {
			input := maps.All(map[string]int{"answer": 42})
			keys, values := Unzip(input).Unpack()

			keysResult := slices.Collect(keys)
			valuesResult := slices.Collect(values)

			So(keysResult, ShouldResemble, []string{"answer"})
			So(valuesResult, ShouldResemble, []int{42})
		})

		Convey("Should handle empty map", func() {
			input := maps.All(map[string]int{})
			keys, values := Unzip(input).Unpack()

			keysResult := slices.Collect(keys)
			valuesResult := slices.Collect(values)

			So(keysResult, ShouldBeNil)
			So(valuesResult, ShouldBeNil)
		})

		Convey("Should work with different types", func() {
			input := maps.All(map[int]string{1: "one", 2: "two", 3: "three"})
			keys, values := Unzip(input).Unpack()

			keysResult := slices.Collect(keys)
			valuesResult := slices.Collect(values)

			// Sort for consistent comparison
			sort.Ints(keysResult)
			sort.Strings(valuesResult)

			So(keysResult, ShouldResemble, []int{1, 2, 3})
			So(valuesResult, ShouldResemble, []string{"one", "three", "two"})
		})

		Convey("Should handle the shared state behavior correctly", func() {
			input := maps.All(map[string]int{"a": 1, "b": 2})
			keys, values := Unzip(input).Unpack()

			// First iteration - keys iterator will consume keys and fill values
			keysResult1 := slices.Collect(keys)
			sort.Strings(keysResult1)
			So(keysResult1, ShouldResemble, []string{"a", "b"})

			// Second iteration - values iterator will consume values and fill keys
			valuesResult := slices.Collect(values)
			sort.Ints(valuesResult)
			So(valuesResult, ShouldResemble, []int{1, 2})

			// The keys iterator cannot be reused after values iterator has run
			// This is the expected behavior due to shared state
		})
	})
}
