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

func ExampleScan() {
	s := slices.Values([]int{1, 2, 3, 4})

	state := 1

	l := Scan(s, &state, func(ctx *int, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), *ctx <= 6
	})

	fmt.Println(slices.Collect(l))
	// Output: [-1 -2 -6]
}

func ExampleScanFunc() {
	state := 1
	product := ScanFunc(&state, func(ctx *int, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), *ctx <= 6
	})

	s := slices.Values([]int{1, 2, 3, 4})
	l := product(s)

	fmt.Println(slices.Collect(l))
	// Output: [-1 -2 -6]
}

func ExampleScan2() {
	s := slices.All([]int{1, 2, 3, 4})

	state := 1

	l := Scan2(s, &state, func(ctx *int, i, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), i < 3
	})

	fmt.Println(maps.Collect(l))
	// Output: map[0:-1 1:-2 2:-6]
}

func ExampleScan2Func() {
	state := 1
	product := Scan2Func(&state, func(ctx *int, i, n int) (int, bool) {
		*ctx *= n

		return -(*ctx), i < 3
	})

	s := slices.All([]int{1, 2, 3, 4})
	l := product(s)

	fmt.Println(maps.Collect(l))
	// Output: map[0:-1 1:-2 2:-6]
}

func TestScan(t *testing.T) {
	Convey("Scan", t, func() {
		Convey("Should scan with accumulating state", func() {
			input := slices.Values([]int{1, 2, 3, 4})
			state := 1

			result := slices.Collect(Scan(input, &state, func(ctx *int, n int) (int, bool) {
				*ctx *= n
				return -*ctx, *ctx <= 6
			}))

			So(result, ShouldResemble, []int{-1, -2, -6})
			So(state, ShouldEqual, 24) // Final state should be 1*1*2*3*4=24
		})

		Convey("Should handle empty sequence", func() {
			input := slices.Values([]int{})
			state := 42

			result := slices.Collect(Scan(input, &state, func(ctx *int, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			}))

			So(result, ShouldBeEmpty)
			So(state, ShouldEqual, 42) // State should remain unchanged
		})

		Convey("Should handle single element", func() {
			input := slices.Values([]int{5})
			state := 10

			result := slices.Collect(Scan(input, &state, func(ctx *int, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			}))

			So(result, ShouldResemble, []int{15})
			So(state, ShouldEqual, 15)
		})

		Convey("Should skip elements when function returns false", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			state := 0

			result := slices.Collect(Scan(input, &state, func(ctx *int, n int) (int, bool) {
				*ctx += n
				return *ctx, n%2 == 0 // Only yield even numbers
			}))

			So(result, ShouldResemble, []int{3, 10}) // 0+1+2=3 (even), 3+3=6 (odd, skipped), 6+4=10 (even), 10+5=15 (odd, skipped)
			So(state, ShouldEqual, 15)               // Final state: 0+1+2+3+4+5=15
		})

		Convey("Should handle early termination", func() {
			input := slices.Values([]int{1, 2, 3, 4, 5})
			state := 0

			seq := Scan(input, &state, func(ctx *int, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			})

			result := make([]int, 0)
			count := 0
			for v := range seq {
				result = append(result, v)
				count++
				if count >= 3 { // Early termination
					break
				}
			}

			So(result, ShouldResemble, []int{1, 3, 6})
			So(state, ShouldEqual, 6) // State should be updated even with early termination
		})

		Convey("Should work with different types", func() {
			input := slices.Values([]string{"a", "bb", "ccc"})
			state := ""

			result := slices.Collect(Scan(input, &state, func(ctx *string, s string) (int, bool) {
				*ctx += s
				return len(*ctx), true
			}))

			So(result, ShouldResemble, []int{1, 3, 6}) // "a"=1, "abb"=3, "abbccc"=6
			So(state, ShouldEqual, "abbccc")
		})

		Convey("Should handle nil context", func() {
			input := slices.Values([]int{1, 2, 3})
			var state *int = nil

			result := slices.Collect(Scan(input, state, func(ctx *int, n int) (int, bool) {
				if ctx == nil {
					return n * 2, true
				}
				return n, true
			}))

			So(result, ShouldResemble, []int{2, 4, 6})
		})
	})
}

func TestScanFunc(t *testing.T) {
	Convey("ScanFunc", t, func() {
		Convey("Should create function that scans with accumulating state", func() {
			state := 1
			product := ScanFunc(&state, func(ctx *int, n int) (int, bool) {
				*ctx *= n
				return -*ctx, *ctx <= 6
			})

			input := slices.Values([]int{1, 2, 3, 4})
			result := slices.Collect(product(input))

			So(result, ShouldResemble, []int{-1, -2, -6})
			So(state, ShouldEqual, 24) // Final state should be 1*1*2*3*4=24
		})

		Convey("Should create reusable function", func() {
			state := 0
			sum := ScanFunc(&state, func(ctx *int, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			})

			input1 := slices.Values([]int{1, 2, 3})
			input2 := slices.Values([]int{4, 5, 6})

			result1 := slices.Collect(sum(input1))
			result2 := slices.Collect(sum(input2))

			So(result1, ShouldResemble, []int{1, 3, 6})
			So(result2, ShouldResemble, []int{10, 15, 21}) // State accumulates: 6+4=10, 10+5=15, 15+6=21
		})

		Convey("Should work with different types", func() {
			state := ""
			length := ScanFunc(&state, func(ctx *string, s string) (int, bool) {
				*ctx += s
				return len(*ctx), true
			})

			input := slices.Values([]string{"a", "bb", "ccc"})
			result := slices.Collect(length(input))

			So(result, ShouldResemble, []int{1, 3, 6})
			So(state, ShouldEqual, "abbccc")
		})
	})
}

func TestScan2(t *testing.T) {
	Convey("Scan2", t, func() {
		Convey("Should scan key-value pairs with accumulating state", func() {
			input := slices.All([]int{1, 2, 3, 4})
			state := 1

			result := maps.Collect(Scan2(input, &state, func(ctx *int, i, n int) (int, bool) {
				*ctx *= n
				return -*ctx, i < 3
			}))

			So(result, ShouldResemble, map[int]int{0: -1, 1: -2, 2: -6})
			So(state, ShouldEqual, 24) // Final state should be 1*1*2*3*4=24
		})

		Convey("Should handle empty sequence", func() {
			input := slices.All([]int{})
			state := 42

			result := maps.Collect(Scan2(input, &state, func(ctx *int, i, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			}))

			So(result, ShouldBeEmpty)
			So(state, ShouldEqual, 42) // State should remain unchanged
		})

		Convey("Should handle single key-value pair", func() {
			input := slices.All([]int{5})
			state := 10

			result := maps.Collect(Scan2(input, &state, func(ctx *int, i, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			}))

			So(result, ShouldResemble, map[int]int{0: 15})
			So(state, ShouldEqual, 15)
		})

		Convey("Should skip key-value pairs when function returns false", func() {
			input := slices.All([]int{1, 2, 3, 4, 5})
			state := 0

			result := maps.Collect(Scan2(input, &state, func(ctx *int, i, n int) (int, bool) {
				*ctx += n
				return *ctx, i%2 == 0 // Only yield even indices
			}))

			So(result, ShouldResemble, map[int]int{0: 1, 2: 6, 4: 15}) // 0:1, 1:3(skipped), 2:6, 3:10(skipped), 4:15
			So(state, ShouldEqual, 15)                                 // Final state: 0+1+2+3+4+5=15
		})

		Convey("Should handle early termination", func() {
			input := slices.All([]int{1, 2, 3, 4, 5})
			state := 0

			seq := Scan2(input, &state, func(ctx *int, i, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			})

			result := make(map[int]int)
			count := 0
			for k, v := range seq {
				result[k] = v
				count++
				if count >= 3 { // Early termination
					break
				}
			}

			So(len(result), ShouldEqual, 3)
			So(result[0], ShouldEqual, 1)
			So(result[1], ShouldEqual, 3)
			So(result[2], ShouldEqual, 6)
			So(state, ShouldEqual, 6) // State should be updated even with early termination
		})

		Convey("Should work with different types", func() {
			input := slices.All([]string{"a", "bb", "ccc"})
			state := ""

			result := maps.Collect(Scan2(input, &state, func(ctx *string, i int, s string) (int, bool) {
				*ctx += s
				return len(*ctx), true
			}))

			So(result, ShouldResemble, map[int]int{0: 1, 1: 3, 2: 6}) // "a"=1, "abb"=3, "abbccc"=6
			So(state, ShouldEqual, "abbccc")
		})

		Convey("Should handle nil context", func() {
			input := slices.All([]int{1, 2, 3})
			var state *int = nil

			result := maps.Collect(Scan2(input, state, func(ctx *int, i, n int) (int, bool) {
				if ctx == nil {
					return n * 2, true
				}
				return n, true
			}))

			So(result, ShouldResemble, map[int]int{0: 2, 1: 4, 2: 6})
		})
	})
}

func TestScan2Func(t *testing.T) {
	Convey("Scan2Func", t, func() {
		Convey("Should create function that scans key-value pairs with accumulating state", func() {
			state := 1
			product := Scan2Func(&state, func(ctx *int, i, n int) (int, bool) {
				*ctx *= n
				return -*ctx, i < 3
			})

			input := slices.All([]int{1, 2, 3, 4})
			result := maps.Collect(product(input))

			So(result, ShouldResemble, map[int]int{0: -1, 1: -2, 2: -6})
			So(state, ShouldEqual, 24) // Final state should be 1*1*2*3*4=24
		})

		Convey("Should create reusable function", func() {
			state := 0
			sum := Scan2Func(&state, func(ctx *int, i, n int) (int, bool) {
				*ctx += n
				return *ctx, true
			})

			input1 := slices.All([]int{1, 2, 3})
			input2 := slices.All([]int{4, 5, 6})

			result1 := maps.Collect(sum(input1))
			result2 := maps.Collect(sum(input2))

			So(result1, ShouldResemble, map[int]int{0: 1, 1: 3, 2: 6})
			So(result2, ShouldResemble, map[int]int{0: 10, 1: 15, 2: 21}) // State accumulates: 6+4=10, 10+5=15, 15+6=21
		})

		Convey("Should work with different types", func() {
			state := ""
			length := Scan2Func(&state, func(ctx *string, i int, s string) (int, bool) {
				*ctx += s
				return len(*ctx), true
			})

			input := slices.All([]string{"a", "bb", "ccc"})
			result := maps.Collect(length(input))

			So(result, ShouldResemble, map[int]int{0: 1, 1: 3, 2: 6})
			So(state, ShouldEqual, "abbccc")
		})
	})
}
