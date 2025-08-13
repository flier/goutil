//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleUniq() {
	s := Uniq(slices.Values([]int{1, 2, 3, 3, 2, 1}))

	fmt.Println(slices.Collect(s))
	// Output: [1 2 3]
}

func ExampleUniqByKey() {
	s := slices.Values([]complex128{1 + 1i, -1 + 2i, -2 + 3i, 2 + 4i, -3 + 5i})
	u := UniqByKey(s, func(c complex128) int { return int(math.Abs(real(c))) })

	fmt.Println(slices.Collect(u))
	// Output:
	// [(1+1i) (-2+3i) (-3+5i)]
}

func ExampleUniqByKeyFunc() {
	abs := UniqByKeyFunc(func(c complex128) int { return int(math.Abs(real(c))) })

	s := slices.Values([]complex128{1 + 1i, -1 + 2i, -2 + 3i, 2 + 4i, -3 + 5i})
	u := abs(s)

	fmt.Println(slices.Collect(u))
	// Output:
	// [(1+1i) (-2+3i) (-3+5i)]
}

func ExampleUniqByKey2() {
	s := slices.All([]string{"foo", "bar", "hello", "world"})
	u := UniqByKey2(s, func(i int, v string) int { return i % 2 })

	fmt.Println(maps.Collect(u))
	// Output:
	// map[0:foo 1:bar]
}

func ExampleUniqByKey2Func() {
	even := UniqByKey2Func(func(i int, v string) int { return i % 2 })

	s := slices.All([]string{"foo", "bar", "hello", "world"})
	u := even(s)

	fmt.Println(maps.Collect(u))
	// Output:
	// map[0:foo 1:bar]
}

func TestUniq(t *testing.T) {
	Convey("Uniq", t, func() {
		Convey("Should remove duplicate integers", func() {
			input := slices.Values([]int{1, 2, 3, 3, 2, 1, 4, 5, 5})
			want := []int{1, 2, 3, 4, 5}

			result := slices.Collect(Uniq(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should remove duplicate strings", func() {
			input := slices.Values([]string{"a", "b", "a", "c", "b", "d"})
			want := []string{"a", "b", "c", "d"}

			result := slices.Collect(Uniq(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle sequence with no duplicates", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 2, 3, 4, 5}

			result := slices.Collect(Uniq(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle sequence with all duplicates", func() {
			input := slices.Values([]int{1, 1, 1, 1, 1})
			want := []int{1}

			result := slices.Collect(Uniq(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})

			result := slices.Collect(Uniq(input))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			want := []int{42}

			result := slices.Collect(Uniq(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 3, 2, 1, 4, 5, 5})

			seq := Uniq(input)
			result := make([]int, 0)
			count := 0
			for v := range seq {
				result = append(result, v)
				count++
				if count >= 3 { // Early termination
					break
				}
			}

			So(result, ShouldResemble, []int{1, 2, 3})
		})

		Convey("Should work with different comparable types", func() {
			input := slices.Values([]float64{1.1, 2.2, 1.1, 3.3, 2.2, 4.4})
			want := []float64{1.1, 2.2, 3.3, 4.4}

			result := slices.Collect(Uniq(input))
			So(result, ShouldResemble, want)
		})
	})
}

func TestUniqByKey(t *testing.T) {
	Convey("UniqByKey", t, func() {
		Convey("Should remove duplicates by absolute value of real part", func() {
			input := slices.Values([]complex128{1 + 1i, -1 + 2i, -2 + 3i, 2 + 4i, -3 + 5i})
			keyFunc := func(c complex128) int { return int(math.Abs(real(c))) }
			want := []complex128{1 + 1i, -2 + 3i, -3 + 5i}

			result := slices.Collect(UniqByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should remove duplicates by string length", func() {
			input := slices.Values([]string{"a", "bb", "ccc", "a", "dd", "eee"})
			keyFunc := func(s string) int { return len(s) }
			want := []string{"a", "bb", "ccc"}

			result := slices.Collect(UniqByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should remove duplicates by first character", func() {
			input := slices.Values([]string{"apple", "banana", "cherry"})
			keyFunc := func(s string) byte { return s[0] }

			result := slices.Collect(UniqByKey(input, keyFunc))
			// All strings start with different characters, so all should be kept
			So(result, ShouldResemble, []string{"apple", "banana", "cherry"})
		})

		Convey("Should handle sequence with no duplicate keys", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			keyFunc := func(n int) int { return n * n }
			want := []int{1, 2, 3, 4, 5}

			result := slices.Collect(UniqByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle sequence with all duplicate keys", func() {
			input := slices.Values([]int{1, -1, 2, -2, 3, -3})
			keyFunc := func(n int) int { return n * n }
			want := []int{1, 2, 3}

			result := slices.Collect(UniqByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			keyFunc := func(n int) int { return n % 2 }

			result := slices.Collect(UniqByKey(input, keyFunc))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			keyFunc := func(n int) int { return n % 10 }

			result := slices.Collect(UniqByKey(input, keyFunc))
			So(result, ShouldResemble, []int{42})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, -1, 2, -2, 3, -3})
			keyFunc := func(n int) int { return n * n }

			seq := UniqByKey(input, keyFunc)
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
	})
}

func TestUniqByKeyFunc(t *testing.T) {
	Convey("UniqByKeyFunc", t, func() {
		Convey("Should create function that removes duplicates by absolute value", func() {
			absKey := UniqByKeyFunc(func(n int) int { return int(math.Abs(float64(n))) })

			input := slices.Values([]int{1, -1, 2, -2, 3, -3})
			want := []int{1, 2, 3}

			result := slices.Collect(absKey(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			lengthKey := UniqByKeyFunc(func(s string) int { return len(s) })

			input1 := slices.Values([]string{"a", "bb", "ccc", "a", "dd"})
			input2 := slices.Values([]string{"x", "yy", "zzz", "w", "vv"})

			result1 := slices.Collect(lengthKey(input1))
			result2 := slices.Collect(lengthKey(input2))

			So(result1, ShouldResemble, []string{"a", "bb", "ccc"})
			So(result2, ShouldResemble, []string{"x", "yy", "zzz"})
		})

		Convey("Should work with different types", func() {
			modKey := UniqByKeyFunc(func(n float64) int { return int(n) % 3 })

			input := slices.Values([]float64{1.1, 2.2, 4.4})

			result := slices.Collect(modKey(input))
			// 1.1 % 3 = 1, 2.2 % 3 = 2, 4.4 % 3 = 1 (duplicate)
			So(len(result), ShouldEqual, 2)
			So(result[0], ShouldEqual, 1.1)
			So(result[1], ShouldEqual, 2.2)
		})
	})
}

func TestUniqByKey2(t *testing.T) {
	Convey("UniqByKey2", t, func() {
		Convey("Should remove duplicates by index modulo 2", func() {
			input := slices.All([]string{"foo", "bar", "hello", "world"})
			keyFunc := func(i int, v string) int { return i % 2 }
			want := map[int]string{0: "foo", 1: "bar"}

			result := maps.Collect(UniqByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should remove duplicates by value length", func() {
			input := maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3})
			keyFunc := func(k string, v int) int { return len(k) }
			want := map[string]int{"a": 1, "bb": 2, "ccc": 3}

			result := maps.Collect(UniqByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should remove duplicates by sum of key and value", func() {
			input := maps.All(map[int]int{1: 1, 2: 2, 3: 3})
			keyFunc := func(k, v int) int { return k + v }

			result := maps.Collect(UniqByKey2(input, keyFunc))
			// 1+1=2, 2+2=4, 3+3=6 - all different, so all should be kept
			So(len(result), ShouldEqual, 3)
			So(result[1], ShouldEqual, 1)
			So(result[2], ShouldEqual, 2)
			So(result[3], ShouldEqual, 3)
		})

		Convey("Should handle map with no duplicate keys", func() {
			input := maps.All(map[string]int{"a": 1, "b": 2, "c": 3})
			keyFunc := func(k string, v int) int { return v * v }
			want := map[string]int{"a": 1, "b": 2, "c": 3}

			result := maps.Collect(UniqByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle map with all duplicate keys", func() {
			input := maps.All(map[int]int{1: 1, 2: 2, 3: 3})
			keyFunc := func(k, v int) int { return 42 }

			result := maps.Collect(UniqByKey2(input, keyFunc))
			// All key-value pairs have the same key (42), so only the first one should remain
			So(len(result), ShouldEqual, 1)
			// The first key-value pair should be kept
			So(result[1], ShouldEqual, 1)
		})

		Convey("Should handle empty map", func() {
			input := maps.All(map[string]int{})
			keyFunc := func(k string, v int) int { return len(k) }

			result := maps.Collect(UniqByKey2(input, keyFunc))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single key-value pair", func() {
			input := maps.All(map[string]int{"single": 42})
			keyFunc := func(k string, v int) int { return len(k) }

			result := maps.Collect(UniqByKey2(input, keyFunc))
			So(result, ShouldResemble, map[string]int{"single": 42})
		})

		Convey("Should handle early termination", func() {
			input := maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3, "d": 4, "eee": 5})
			keyFunc := func(k string, v int) int { return len(k) }

			seq := UniqByKey2(input, keyFunc)
			result := make(map[string]int)
			count := 0
			for k, v := range seq {
				result[k] = v
				count++
				if count >= 2 { // Early termination
					break
				}
			}

			So(len(result), ShouldEqual, 2)
		})
	})
}

func TestUniqByKey2Func(t *testing.T) {
	Convey("UniqByKey2Func", t, func() {
		Convey("Should create function that removes duplicates by index modulo 2", func() {
			evenOdd := UniqByKey2Func(func(i int, v string) int { return i % 2 })

			input := slices.All([]string{"foo", "bar", "hello", "world"})
			want := map[int]string{0: "foo", 1: "bar"}

			result := maps.Collect(evenOdd(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should create reusable function", func() {
			lengthKey := UniqByKey2Func(func(k string, v int) int { return len(k) })

			input1 := maps.All(map[string]int{"a": 1, "bb": 2, "ccc": 3})
			input2 := maps.All(map[string]int{"x": 10, "yy": 20, "zzz": 30})

			result1 := maps.Collect(lengthKey(input1))
			result2 := maps.Collect(lengthKey(input2))

			So(result1, ShouldResemble, map[string]int{"a": 1, "bb": 2, "ccc": 3})
			So(result2, ShouldResemble, map[string]int{"x": 10, "yy": 20, "zzz": 30})
		})

		Convey("Should work with different types", func() {
			sumKey := UniqByKey2Func(func(k int, v float64) int { return k + int(v) })

			input := maps.All(map[int]float64{1: 1.5, 2: 2.5, 3: 3.5, 4: 4.5})
			want := map[int]float64{1: 1.5, 2: 2.5, 3: 3.5, 4: 4.5}

			result := maps.Collect(sumKey(input))
			So(result, ShouldResemble, want)
		})
	})
}
