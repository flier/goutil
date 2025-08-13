//go:build go1.23

package xiter_test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	. "github.com/flier/goutil/pkg/xiter"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleDedup() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	d := Dedup(s)

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 6 7]
}

func ExampleDedupBy() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	d := DedupBy(s, func(x, y int) bool { return x == y })

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 6 7]
}

func ExampleDedupByFunc() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	dedup := DedupByFunc(func(x, y int) bool { return x == y })
	d := dedup(s)

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 6 7]
}

func ExampleDedupByKey() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	d := DedupByKey(s, func(n int) int { return n % 2 })

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 7]
}

func ExampleDedupByKeyFunc() {
	s := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
	dedup := DedupByKeyFunc(func(n int) int { return n % 2 })
	d := dedup(s)

	fmt.Println(slices.Collect(d))
	// Output: [1 2 3 4 7]
}

func ExampleDedupByKey2() {
	s := Zip(
		slices.Values([]int{1, 2, 2, 3, 4}),
		slices.Values([]string{"a", "b", "c", "d", "e"}))

	d := DedupByKey2(s, func(k int, v string) int { return k })

	fmt.Println(maps.Collect(d))

	// Output:
	// map[1:a 2:b 3:d 4:e]
}

func ExampleDedupByKey2Func() {
	s := Zip(
		slices.Values([]int{1, 2, 2, 3, 4}),
		slices.Values([]string{"a", "b", "c", "d", "e"}))

	dedup := DedupByKey2Func(func(k int, v string) int { return k })
	d := dedup(s)

	fmt.Println(maps.Collect(d))

	// Output:
	// map[1:a 2:b 3:d 4:e]
}

func TestDedup(t *testing.T) {
	Convey("Dedup", t, func() {
		Convey("Should remove consecutive duplicates", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			want := []int{1, 2, 3, 4, 6, 7}

			result := slices.Collect(Dedup(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle no duplicates", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			want := []int{1, 2, 3, 4, 5}

			result := slices.Collect(Dedup(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle all duplicates", func() {
			input := slices.Values([]int{1, 1, 1, 1})
			want := []int{1}

			result := slices.Collect(Dedup(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})

			result := slices.Collect(Dedup(input))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})

			result := slices.Collect(Dedup(input))
			So(result, ShouldResemble, []int{42})
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "a", "b", "c", "c", "d"})
			want := []string{"a", "b", "c", "d"}

			result := slices.Collect(Dedup(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			want := []int{1, 2, 3}

			var result []int

			for v := range Dedup(input) {
				if len(result) == 3 {
					break
				}
				result = append(result, v)
			}

			So(result, ShouldResemble, want)
		})
	})
}

func TestDedupBy(t *testing.T) {
	Convey("DedupBy", t, func() {
		Convey("Should remove consecutive duplicates using custom function", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			equal := func(x, y int) bool { return x == y }
			want := []int{1, 2, 3, 4, 6, 7}

			result := slices.Collect(DedupBy(input, equal))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with custom equality logic", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 25}, {"Charlie", 30},
				{"David", 30}, {"Eve", 35},
			})
			sameAge := func(l, r Person) bool { return l.Age == r.Age }
			want := []Person{{"Alice", 25}, {"Charlie", 30}, {"Eve", 35}}

			result := slices.Collect(DedupBy(input, sameAge))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle no duplicates", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			equal := func(x, y int) bool { return x == y }
			want := []int{1, 2, 3, 4, 5}

			result := slices.Collect(DedupBy(input, equal))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			equal := func(x, y int) bool { return x == y }

			result := slices.Collect(DedupBy(input, equal))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			equal := func(x, y int) bool { return x == y }

			result := slices.Collect(DedupBy(input, equal))
			So(result, ShouldResemble, []int{42})
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			equal := func(x, y int) bool { return x == y }
			want := []int{1, 2, 3}

			var result []int

			for v := range DedupBy(input, equal) {
				if len(result) == 3 {
					break
				}
				result = append(result, v)
			}

			So(result, ShouldResemble, want)
		})
	})
}

func TestDedupByFunc(t *testing.T) {
	Convey("DedupByFunc", t, func() {
		Convey("Should create reusable deduplication function", func() {
			dedup := DedupByFunc(func(x, y int) bool { return x == y })

			input1 := slices.Values([]int{1, 2, 2, 3, 4, 4})
			input2 := slices.Values([]int{5, 5, 6, 7, 7, 8})

			result1 := slices.Collect(dedup(input1))
			result2 := slices.Collect(dedup(input2))

			So(result1, ShouldResemble, []int{1, 2, 3, 4})
			So(result2, ShouldResemble, []int{5, 6, 7, 8})
		})

		Convey("Should work with different types", func() {
			type Person struct {
				Name string
				Age  int
			}

			dedup := DedupByFunc(func(l, r Person) bool { return l.Age == r.Age })

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 25}, {"Charlie", 30},
				{"David", 30}, {"Eve", 35},
			})
			want := []Person{{"Alice", 25}, {"Charlie", 30}, {"Eve", 35}}

			result := slices.Collect(dedup(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle edge cases consistently", func() {
			dedup := DedupByFunc(func(x, y int) bool { return x == y })

			// Empty sequence
			emptyInput := slices.Values([]int{})
			emptyResult := slices.Collect(dedup(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element
			singleInput := slices.Values([]int{42})
			singleResult := slices.Collect(dedup(singleInput))
			So(singleResult, ShouldResemble, []int{42})

			// No duplicates
			noDupInput := slices.Values([]int{1, 2, 3})
			noDupResult := slices.Collect(dedup(noDupInput))
			So(noDupResult, ShouldResemble, []int{1, 2, 3})
		})
	})
}

func TestDedupByKey(t *testing.T) {
	Convey("DedupByKey", t, func() {
		Convey("Should remove consecutive duplicates based on key", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			keyFunc := func(n int) int { return n % 2 }
			want := []int{1, 2, 3, 4, 7}

			result := slices.Collect(DedupByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with custom key extraction", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 25}, {"Charlie", 30},
				{"David", 30}, {"Eve", 35},
			})
			keyFunc := func(p Person) int { return p.Age }
			want := []Person{{"Alice", 25}, {"Charlie", 30}, {"Eve", 35}}

			result := slices.Collect(DedupByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle no duplicates", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			keyFunc := func(n int) int { return n % 2 }
			want := []int{1, 2, 3, 4, 5}

			result := slices.Collect(DedupByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			keyFunc := func(n int) int { return n % 2 }

			result := slices.Collect(DedupByKey(input, keyFunc))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{42})
			keyFunc := func(n int) int { return n % 2 }

			result := slices.Collect(DedupByKey(input, keyFunc))
			So(result, ShouldResemble, []int{42})
		})

		Convey("Should work with different key types", func() {
			type Person struct {
				Name string
				Age  int
			}

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 25}, {"Charlie", 30},
				{"David", 30}, {"Eve", 35},
			})
			keyFunc := func(p Person) string { return p.Name[:1] } // First letter of name
			want := []Person{{"Alice", 25}, {"Bob", 25}, {"Charlie", 30}, {"David", 30}, {"Eve", 35}}

			result := slices.Collect(DedupByKey(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 2, 3, 4, 4, 6, 7, 7})
			keyFunc := func(n int) int { return n % 2 }
			want := []int{1, 2, 3}

			var result []int

			for v := range DedupByKey(input, keyFunc) {
				if len(result) == 3 {
					break
				}
				result = append(result, v)
			}

			So(result, ShouldResemble, want)
		})
	})
}

func TestDedupByKeyFunc(t *testing.T) {
	Convey("DedupByKeyFunc", t, func() {
		Convey("Should create reusable key-based deduplication function", func() {
			dedup := DedupByKeyFunc(func(n int) int { return n % 2 })

			input1 := slices.Values([]int{1, 2, 2, 3, 4, 4})
			input2 := slices.Values([]int{5, 5, 6, 7, 7, 8})

			result1 := slices.Collect(dedup(input1))
			result2 := slices.Collect(dedup(input2))

			So(result1, ShouldResemble, []int{1, 2, 3, 4})
			So(result2, ShouldResemble, []int{5, 6, 7, 8})
		})

		Convey("Should work with different types", func() {
			type Person struct {
				Name string
				Age  int
			}

			dedup := DedupByKeyFunc(func(p Person) int { return p.Age })

			input := slices.Values([]Person{
				{"Alice", 25}, {"Bob", 25}, {"Charlie", 30},
				{"David", 30}, {"Eve", 35},
			})
			want := []Person{{"Alice", 25}, {"Charlie", 30}, {"Eve", 35}}

			result := slices.Collect(dedup(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle edge cases consistently", func() {
			dedup := DedupByKeyFunc(func(n int) int { return n % 2 })

			// Empty sequence
			emptyInput := slices.Values([]int{})
			emptyResult := slices.Collect(dedup(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element
			singleInput := slices.Values([]int{42})
			singleResult := slices.Collect(dedup(singleInput))
			So(singleResult, ShouldResemble, []int{42})

			// No duplicates
			noDupInput := slices.Values([]int{1, 2, 3})
			noDupResult := slices.Collect(dedup(noDupInput))
			So(noDupResult, ShouldResemble, []int{1, 2, 3})
		})
	})
}

func TestDedupByKey2(t *testing.T) {
	Convey("DedupByKey2", t, func() {
		Convey("Should remove consecutive duplicates from key-value pairs based on key", func() {
			keys := slices.Values([]int{1, 2, 2, 3, 4})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)

			keyFunc := func(k int, v string) int { return k }
			want := map[int]string{1: "a", 2: "b", 3: "d", 4: "e"}

			result := maps.Collect(DedupByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with custom key extraction", func() {
			keys := slices.Values([]int{1, 2, 3, 4, 5})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)

			keyFunc := func(k int, v string) int { return k % 2 }
			want := map[int]string{1: "a", 2: "b", 3: "c", 4: "d", 5: "e"}

			result := maps.Collect(DedupByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should remove consecutive duplicates with same key", func() {
			keys := slices.Values([]int{1, 1, 2, 2, 3, 3})
			values := slices.Values([]string{"a", "b", "c", "d", "e", "f"})
			input := Zip(keys, values)

			keyFunc := func(k int, v string) int { return k }
			// 应该去除连续的重复键，保留第一个
			want := map[int]string{1: "a", 2: "c", 3: "e"}

			result := maps.Collect(DedupByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle no duplicates", func() {
			keys := slices.Values([]int{1, 2, 3, 4, 5})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)

			keyFunc := func(k int, v string) int { return k }
			want := map[int]string{1: "a", 2: "b", 3: "c", 4: "d", 5: "e"}

			result := maps.Collect(DedupByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle empty sequence", func() {
			input := Zip(slices.Values([]int{}), slices.Values([]string{}))
			keyFunc := func(k int, v string) int { return k }

			result := maps.Collect(DedupByKey2(input, keyFunc))
			So(result, ShouldBeEmpty)
		})

		Convey("Should handle single element", func() {
			keys := slices.Values([]int{42})
			values := slices.Values([]string{"answer"})
			input := Zip(keys, values)

			keyFunc := func(k int, v string) int { return k }
			want := map[int]string{42: "answer"}

			result := maps.Collect(DedupByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should work with different key types", func() {
			keys := slices.Values([]int{1, 2, 3, 4, 5})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)

			keyFunc := func(k int, v string) string { return v } // Use value as key
			want := map[int]string{1: "a", 2: "b", 3: "c", 4: "d", 5: "e"}

			result := maps.Collect(DedupByKey2(input, keyFunc))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle early termination", func() {
			keys := slices.Values([]int{1, 2, 3, 4, 5})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)

			keyFunc := func(k int, v string) string { return v } // Use value as key
			want := []int{1, 2, 3}

			var result []int

			for v := range DedupByKey2(input, keyFunc) {
				if len(result) == 3 {
					break
				}
				result = append(result, v)
			}

			So(result, ShouldResemble, want)
		})
	})
}

func TestDedupByKey2Func(t *testing.T) {
	Convey("DedupByKey2Func", t, func() {
		Convey("Should create reusable key-value deduplication function", func() {
			dedup := DedupByKey2Func(func(k int, v string) int { return k })

			keys1 := slices.Values([]int{1, 2, 2, 3, 4})
			values1 := slices.Values([]string{"a", "b", "c", "d", "e"})
			input1 := Zip(keys1, values1)

			keys2 := slices.Values([]int{5, 5, 6, 7, 7})
			values2 := slices.Values([]string{"f", "g", "h", "i", "j"})
			input2 := Zip(keys2, values2)

			result1 := maps.Collect(dedup(input1))
			result2 := maps.Collect(dedup(input2))

			So(result1, ShouldResemble, map[int]string{1: "a", 2: "b", 3: "d", 4: "e"})
			So(result2, ShouldResemble, map[int]string{5: "f", 6: "h", 7: "i"})
		})

		Convey("Should work with different types", func() {
			dedup := DedupByKey2Func(func(k int, v string) int { return k % 2 })

			keys := slices.Values([]int{1, 2, 3, 4, 5})
			values := slices.Values([]string{"a", "b", "c", "d", "e"})
			input := Zip(keys, values)
			want := map[int]string{1: "a", 2: "b", 3: "c", 4: "d", 5: "e"}

			result := maps.Collect(dedup(input))
			So(result, ShouldResemble, want)
		})

		Convey("Should handle edge cases consistently", func() {
			dedup := DedupByKey2Func(func(k int, v string) int { return k })

			// Empty sequence
			emptyInput := Zip(slices.Values([]int{}), slices.Values([]string{}))
			emptyResult := maps.Collect(dedup(emptyInput))
			So(emptyResult, ShouldBeEmpty)

			// Single element
			singleKeys := slices.Values([]int{42})
			singleValues := slices.Values([]string{"answer"})
			singleInput := Zip(singleKeys, singleValues)
			singleResult := maps.Collect(dedup(singleInput))
			So(singleResult, ShouldResemble, map[int]string{42: "answer"})

			// No duplicates
			noDupKeys := slices.Values([]int{1, 2, 3})
			noDupValues := slices.Values([]string{"a", "b", "c"})
			noDupInput := Zip(noDupKeys, noDupValues)
			noDupResult := maps.Collect(dedup(noDupInput))
			So(noDupResult, ShouldResemble, map[int]string{1: "a", 2: "b", 3: "c"})
		})
	})
}
