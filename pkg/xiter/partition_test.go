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

func ExamplePartition() {
	s := slices.Values([]int{1, 2, 3, 4, 5})
	odd, even := Partition(s, func(n int) bool { return n%2 != 0 }).Unpack()

	fmt.Println(slices.Collect(odd), slices.Collect(even))
	// Output:
	// [1 3 5] [2 4]
}

func ExamplePartitionFunc() {
	byEven := PartitionFunc(func(n int) bool { return n%2 != 0 })

	s := slices.Values([]int{1, 2, 3, 4, 5})
	odd, even := byEven(s).Unpack()

	fmt.Println(slices.Collect(odd), slices.Collect(even))
	// Output:
	// [1 3 5] [2 4]
}

func ExamplePartition2() {
	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	short, long := Partition2(s, func(k, v string) bool { return len(k) < 4 }).Unpack()

	fmt.Println(maps.Collect(short), maps.Collect(long))
	// Output:
	// map[foo:bar] map[hello:world]
}

func ExamplePartition2Func() {
	byLen := Partition2Func(func(k, v string) bool { return len(k) < 4 })

	s := maps.All(map[string]string{"foo": "bar", "hello": "world"})
	short, long := byLen(s).Unpack()

	fmt.Println(maps.Collect(short), maps.Collect(long))
	// Output:
	// map[foo:bar] map[hello:world]
}

func TestPartition(t *testing.T) {
	Convey("Partition", t, func() {
		Convey("Should partition integers by even/odd", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8})
			predicate := func(n int) bool { return n%2 != 0 }

			odd, even := Partition(input, predicate).Unpack()

			oddResult := slices.Collect(odd)
			evenResult := slices.Collect(even)

			So(oddResult, ShouldResemble, []int{1, 3, 5, 7})
			So(evenResult, ShouldResemble, []int{2, 4, 6, 8})
		})

		Convey("Should partition strings by length", func() {
			input := slices.Values([]string{"a", "bb", "ccc", "dddd", "eeeee"})
			predicate := func(s string) bool { return len(s) <= 2 }

			short, long := Partition(input, predicate).Unpack()

			shortResult := slices.Collect(short)
			longResult := slices.Collect(long)

			So(shortResult, ShouldResemble, []string{"a", "bb"})
			So(longResult, ShouldResemble, []string{"ccc", "dddd", "eeeee"})
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			predicate := func(n int) bool { return n > 0 }

			positive, negative := Partition(input, predicate).Unpack()

			positiveResult := slices.Collect(positive)
			negativeResult := slices.Collect(negative)

			So(positiveResult, ShouldBeEmpty)
			So(negativeResult, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			predicate := func(n int) bool { return n > 50 }

			large, small := Partition(input, predicate).Unpack()

			largeResult := slices.Collect(large)
			smallResult := slices.Collect(small)

			So(largeResult, ShouldBeEmpty)
			So(smallResult, ShouldResemble, []int{42})
		})

		Convey("Should handle all elements matching predicate", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			predicate := func(n int) bool { return n > 0 }

			positive, negative := Partition(input, predicate).Unpack()

			positiveResult := slices.Collect(positive)
			negativeResult := slices.Collect(negative)

			So(positiveResult, ShouldResemble, []int{1, 2, 3, 4, 5})
			So(negativeResult, ShouldBeEmpty)
		})

		Convey("Should handle no elements matching predicate", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			predicate := func(n int) bool { return n > 10 }

			large, small := Partition(input, predicate).Unpack()

			largeResult := slices.Collect(large)
			smallResult := slices.Collect(small)

			So(largeResult, ShouldBeEmpty)
			So(smallResult, ShouldResemble, []int{1, 2, 3, 4, 5})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8})
			predicate := func(n int) bool { return n%2 != 0 }

			odd, even := Partition(input, predicate).Unpack()

			// Early termination on odd sequence
			oddResult := make([]int, 0)
			count := 0
			for v := range odd {
				oddResult = append(oddResult, v)
				count++
				if count >= 2 {
					break
				}
			}

			// Collect even sequence
			evenResult := slices.Collect(even)

			So(oddResult, ShouldResemble, []int{1, 3})
			So(evenResult, ShouldResemble, []int{2, 4, 6, 8})
		})
	})
}

func TestPartitionFunc(t *testing.T) {
	Convey("PartitionFunc", t, func() {
		Convey("Should create function that partitions by even/odd", func() {
			byEven := PartitionFunc(func(n int) bool { return n%2 != 0 })

			input := slices.Values([]int{1, 2, 3, 4, 5})
			odd, even := byEven(input).Unpack()

			oddResult := slices.Collect(odd)
			evenResult := slices.Collect(even)

			So(oddResult, ShouldResemble, []int{1, 3, 5})
			So(evenResult, ShouldResemble, []int{2, 4})
		})

		Convey("Should create reusable function", func() {
			byLength := PartitionFunc(func(s string) bool { return len(s) <= 3 })

			input1 := slices.Values([]string{"a", "bb", "ccc", "dddd"})
			input2 := slices.Values([]string{"x", "yy", "zzz", "wwww"})

			short1, long1 := byLength(input1).Unpack()
			short2, long2 := byLength(input2).Unpack()

			So(slices.Collect(short1), ShouldResemble, []string{"a", "bb", "ccc"})
			So(slices.Collect(long1), ShouldResemble, []string{"dddd"})
			So(slices.Collect(short2), ShouldResemble, []string{"x", "yy", "zzz"})
			So(slices.Collect(long2), ShouldResemble, []string{"wwww"})
		})

		Convey("Should work with different types", func() {
			byPositive := PartitionFunc(func(n float64) bool { return n > 0.0 })

			input := slices.Values([]float64{-1.5, 2.7, -3.2, 4.1, -5.8})
			positive, negative := byPositive(input).Unpack()

			positiveResult := slices.Collect(positive)
			negativeResult := slices.Collect(negative)

			So(positiveResult, ShouldResemble, []float64{2.7, 4.1})
			So(negativeResult, ShouldResemble, []float64{-1.5, -3.2, -5.8})
		})
	})
}

func TestPartition2(t *testing.T) {
	Convey("Partition2", t, func() {
		Convey("Should partition key-value pairs by key length", func() {
			input := maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3, "dddd": 4})
			predicate := func(k string, v int) bool { return len(k) <= 2 }

			short, long := Partition2(input, predicate).Unpack()

			shortResult := maps.Collect(short)
			longResult := maps.Collect(long)

			So(shortResult, ShouldResemble, map[string]int{"a": 1, "bb": 2})
			So(longResult, ShouldResemble, map[string]int{"ccc": 3, "dddd": 4})
		})

		Convey("Should partition by value criteria", func() {
			input := maps.All(map[string]int{"small": 5, "medium": 15, "large": 25, "huge": 35})
			predicate := func(k string, v int) bool { return v <= 20 }

			small, large := Partition2(input, predicate).Unpack()

			smallResult := maps.Collect(small)
			largeResult := maps.Collect(large)

			So(smallResult, ShouldResemble, map[string]int{"small": 5, "medium": 15})
			So(largeResult, ShouldResemble, map[string]int{"large": 25, "huge": 35})
		})

		Convey("Should handle empty map", func() {
			input := maps.All(map[string]int{})
			predicate := func(k string, v int) bool { return len(k) > 0 }

			nonEmpty, empty := Partition2(input, predicate).Unpack()

			nonEmptyResult := maps.Collect(nonEmpty)
			emptyResult := maps.Collect(empty)

			So(nonEmptyResult, ShouldBeEmpty)
			So(emptyResult, ShouldBeEmpty)
		})

		Convey("Should handle single key-value pair", func() {
			input := maps.All(map[string]int{"single": 42})
			predicate := func(k string, v int) bool { return v > 50 }

			large, small := Partition2(input, predicate).Unpack()

			largeResult := maps.Collect(large)
			smallResult := maps.Collect(small)

			So(largeResult, ShouldBeEmpty)
			So(smallResult, ShouldResemble, map[string]int{"single": 42})
		})

		Convey("Should handle all pairs matching predicate", func() {
			input := maps.All(map[string]int{"a": 1, "b": 2, "c": 3})
			predicate := func(k string, v int) bool { return len(k) == 1 }

			short, long := Partition2(input, predicate).Unpack()

			shortResult := maps.Collect(short)
			longResult := maps.Collect(long)

			So(shortResult, ShouldResemble, map[string]int{"a": 1, "b": 2, "c": 3})
			So(longResult, ShouldBeEmpty)
		})

		Convey("Should handle no pairs matching predicate", func() {
			input := maps.All(map[string]int{"aa": 1, "bb": 2, "cc": 3})
			predicate := func(k string, v int) bool { return len(k) == 1 }

			short, long := Partition2(input, predicate).Unpack()

			shortResult := maps.Collect(short)
			longResult := maps.Collect(long)

			So(shortResult, ShouldBeEmpty)
			So(longResult, ShouldResemble, map[string]int{"aa": 1, "bb": 2, "cc": 3})
		})

		Convey("Should handle early termination", func() {
			input := maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3, "dddd": 4})
			predicate := func(k string, v int) bool { return len(k) <= 2 }

			short, long := Partition2(input, predicate).Unpack()

			// Early termination on short sequence
			shortResult := make(map[string]int)
			count := 0
			for k, v := range short {
				shortResult[k] = v
				count++
				if count >= 1 {
					break
				}
			}

			// Collect long sequence
			longResult := maps.Collect(long)

			So(len(shortResult), ShouldEqual, 1)
			So(longResult, ShouldResemble, map[string]int{"ccc": 3, "dddd": 4})
		})
	})
}

func TestPartition2Func(t *testing.T) {
	Convey("Partition2Func", t, func() {
		Convey("Should create function that partitions by key length", func() {
			byKeyLength := Partition2Func(func(k string, v int) bool { return len(k) <= 2 })

			input := maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3, "dddd": 4})
			short, long := byKeyLength(input).Unpack()

			shortResult := maps.Collect(short)
			longResult := maps.Collect(long)

			So(shortResult, ShouldResemble, map[string]int{"a": 1, "bb": 2})
			So(longResult, ShouldResemble, map[string]int{"ccc": 3, "dddd": 4})
		})

		Convey("Should create reusable function", func() {
			byValueRange := Partition2Func(func(k string, v int) bool { return v <= 20 })

			input1 := maps.All(map[string]int{"small": 5, "medium": 15, "large": 25})
			input2 := maps.All(map[string]int{"tiny": 1, "huge": 50, "normal": 10})

			small1, large1 := byValueRange(input1).Unpack()
			small2, large2 := byValueRange(input2).Unpack()

			So(maps.Collect(small1), ShouldResemble, map[string]int{"small": 5, "medium": 15})
			So(maps.Collect(large1), ShouldResemble, map[string]int{"large": 25})
			So(maps.Collect(small2), ShouldResemble, map[string]int{"tiny": 1, "normal": 10})
			So(maps.Collect(large2), ShouldResemble, map[string]int{"huge": 50})
		})

		Convey("Should work with different types", func() {
			byKeyPrefix := Partition2Func(func(k string, v float64) bool { return k[0] == 'a' })

			input := maps.All(map[string]float64{"apple": 1.5, "banana": 2.3, "apricot": 3.1, "cherry": 4.2})
			aWords, otherWords := byKeyPrefix(input).Unpack()

			aWordsResult := maps.Collect(aWords)
			otherWordsResult := maps.Collect(otherWords)

			So(aWordsResult, ShouldResemble, map[string]float64{"apple": 1.5, "apricot": 3.1})
			So(otherWordsResult, ShouldResemble, map[string]float64{"banana": 2.3, "cherry": 4.2})
		})
	})
}
