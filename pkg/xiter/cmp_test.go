//go:build go1.23

package xiter_test

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleCompare() {
	fmt.Println(Compare(slices.Values([]int{1}), slices.Values([]int{1})))    // 0
	fmt.Println(Compare(slices.Values([]int{1}), slices.Values([]int{1, 2}))) // -1
	fmt.Println(Compare(slices.Values([]int{1, 2}), slices.Values([]int{1}))) // 1

	// Output:
	// 0
	// -1
	// 1
}

func ExampleCompareBy() {
	type User struct {
		Name string
		Age  int
	}

	byAge := func(l, r User) int { return cmp.Compare(l.Age, r.Age) }

	empty := slices.Values([]User{})
	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})

	fmt.Println(CompareBy(empty, cathy, byAge)) // -1
	fmt.Println(CompareBy(joe, joe, byAge))     // 0
	fmt.Println(CompareBy(joe, cathy, byAge))   // 1
	fmt.Println(CompareBy(joe, empty, byAge))   // 1

	// Output:
	// -1
	// 0
	// 1
	// 1
}

func ExampleCompareByFunc() {
	type User struct {
		Name string
		Age  int
	}

	compareByAge := CompareByFunc(func(l, r User) int { return cmp.Compare(l.Age, r.Age) })

	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})

	fmt.Println(compareByAge(joe, cathy))

	// Output:
	// 1
}

func ExampleCompareByKey() {
	type User struct {
		Name string
		Age  int
	}

	empty := slices.Values([]User{})
	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})
	byAge := func(u User) int { return u.Age }

	fmt.Println(CompareByKey(empty, cathy, byAge)) // -1
	fmt.Println(CompareByKey(joe, joe, byAge))     // 0
	fmt.Println(CompareByKey(joe, cathy, byAge))   // 1
	fmt.Println(CompareByKey(joe, empty, byAge))   // 1

	// Output:
	// -1
	// 0
	// 1
	// 1
}

func ExampleCompareByKeyFunc() {
	type User struct {
		Name string
		Age  int
	}

	compareByAge := CompareByKeyFunc(func(u User) int { return u.Age })

	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})

	fmt.Println(compareByAge(joe, cathy))

	// Output:
	// 1
}

func ExampleEqual() {
	fmt.Println(Equal(slices.Values([]int{1}), slices.Values([]int{1})))
	fmt.Println(Equal(slices.Values([]int{1}), slices.Values([]int{1, 2})))

	// Output:
	// true
	// false
}

func ExampleEqualBy() {
	type User struct {
		Name string
		Age  int
	}

	sameAge := func(l, r User) bool { return l.Age == r.Age }

	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})

	fmt.Println(EqualBy(joe, joe, sameAge))
	fmt.Println(EqualBy(joe, cathy, sameAge))

	// Output:
	// true
	// false
}

func ExampleEqualByFunc() {
	type User struct {
		Name string
		Age  int
	}

	sameAge := EqualByFunc(func(l, r User) bool { return l.Age == r.Age })

	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})

	fmt.Println(sameAge(joe, joe))
	fmt.Println(sameAge(joe, cathy))

	// Output:
	// true
	// false
}

func ExampleEqualByKey() {
	type User struct {
		Name string
		Age  int
	}

	userAge := func(u User) int { return u.Age }

	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})

	fmt.Println(EqualByKey(joe, joe, userAge))
	fmt.Println(EqualByKey(joe, cathy, userAge))

	// Output:
	// true
	// false
}

func ExampleEqualByKeyFunc() {
	type User struct {
		Name string
		Age  int
	}

	sameAge := EqualByKeyFunc(func(u User) int { return u.Age })

	joe := slices.Values([]User{{"joe", 12}})
	cathy := slices.Values([]User{{"cathy", 8}})

	fmt.Println(sameAge(joe, joe))
	fmt.Println(sameAge(joe, cathy))

	// Output:
	// true
	// false
}
func ExampleNotEqual() {
	fmt.Println(NotEqual(slices.Values([]int{1}), slices.Values([]int{1})))
	fmt.Println(NotEqual(slices.Values([]int{1}), slices.Values([]int{1, 2})))

	// Output:
	// false
	// true
}

func ExampleLessThan() {
	fmt.Println(LessThan(slices.Values([]int{1}), slices.Values([]int{1})))
	fmt.Println(LessThan(slices.Values([]int{1}), slices.Values([]int{1, 2})))
	fmt.Println(LessThan(slices.Values([]int{1, 2}), slices.Values([]int{1})))

	// Output:
	// false
	// true
	// false
}

func ExampleLessOrEqual() {
	fmt.Println(LessOrEqual(slices.Values([]int{1}), slices.Values([]int{1})))
	fmt.Println(LessOrEqual(slices.Values([]int{1}), slices.Values([]int{1, 2})))
	fmt.Println(LessOrEqual(slices.Values([]int{1, 2}), slices.Values([]int{1})))

	// Output:
	// true
	// true
	// false
}

func ExampleGreaterThan() {
	fmt.Println(GreaterThan(slices.Values([]int{1}), slices.Values([]int{1})))
	fmt.Println(GreaterThan(slices.Values([]int{1}), slices.Values([]int{1, 2})))
	fmt.Println(GreaterThan(slices.Values([]int{1, 2}), slices.Values([]int{1})))

	// Output:
	// false
	// false
	// true
}

func ExampleGreaterOrEqual() {
	fmt.Println(GreaterOrEqual(slices.Values([]int{1}), slices.Values([]int{1})))
	fmt.Println(GreaterOrEqual(slices.Values([]int{1}), slices.Values([]int{1, 2})))
	fmt.Println(GreaterOrEqual(slices.Values([]int{1, 2}), slices.Values([]int{1})))

	// Output:
	// true
	// false
	// true
}

func TestCompare(t *testing.T) {
	Convey("Compare", t, func() {
		Convey("Should return 0 for equal sequences", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 3})

			result := Compare(seq1, seq2)
			So(result, ShouldEqual, 0)
		})

		Convey("Should return -1 when first sequence is shorter", func() {
			seq1 := slices.Values([]int{1, 2})
			seq2 := slices.Values([]int{1, 2, 3})

			result := Compare(seq1, seq2)
			So(result, ShouldEqual, -1)
		})

		Convey("Should return 1 when second sequence is shorter", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2})

			result := Compare(seq1, seq2)
			So(result, ShouldEqual, 1)
		})

		Convey("Should return -1 when first sequence is lexicographically less", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 4})

			result := Compare(seq1, seq2)
			So(result, ShouldEqual, -1)
		})

		Convey("Should return 1 when first sequence is lexicographically greater", func() {
			seq1 := slices.Values([]int{1, 2, 4})
			seq2 := slices.Values([]int{1, 2, 3})

			result := Compare(seq1, seq2)
			So(result, ShouldEqual, 1)
		})

		Convey("Should handle empty sequences", func() {
			empty := slices.Values([]int{})
			nonEmpty := slices.Values([]int{1, 2, 3})

			result1 := Compare(empty, empty)
			result2 := Compare(empty, nonEmpty)
			result3 := Compare(nonEmpty, empty)

			So(result1, ShouldEqual, 0)
			So(result2, ShouldEqual, -1)
			So(result3, ShouldEqual, 1)
		})

		Convey("Should work with different types", func() {
			seq1 := slices.Values([]string{"a", "b", "c"})
			seq2 := slices.Values([]string{"a", "b", "d"})

			result := Compare(seq1, seq2)
			So(result, ShouldEqual, -1)
		})
	})
}

func TestCompareBy(t *testing.T) {
	Convey("CompareBy", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		compareByAge := func(l, r Person) int {
			return cmp.Compare(l.Age, r.Age)
		}

		Convey("Should compare using custom comparison function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 20}, {"David", 35}})

			result := CompareBy(seq1, seq2, compareByAge)
			So(result, ShouldEqual, 1) // Alice(25) > Charlie(20)
		})

		Convey("Should return 0 for equal sequences", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 25}, {"David", 30}})

			result := CompareBy(seq1, seq2, compareByAge)
			So(result, ShouldEqual, 0)
		})

		Convey("Should handle empty sequences", func() {
			empty := slices.Values([]Person{})
			nonEmpty := slices.Values([]Person{{"Alice", 25}})

			result1 := CompareBy(empty, empty, compareByAge)
			result2 := CompareBy(empty, nonEmpty, compareByAge)
			result3 := CompareBy(nonEmpty, empty, compareByAge)

			So(result1, ShouldEqual, 0)
			So(result2, ShouldEqual, -1)
			So(result3, ShouldEqual, 1)
		})
	})
}

func TestCompareByFunc(t *testing.T) {
	Convey("CompareByFunc", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		compareByAge := CompareByFunc(func(l, r Person) int {
			return cmp.Compare(l.Age, r.Age)
		})

		Convey("Should create reusable comparison function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 20}, {"David", 35}})

			result1 := compareByAge(seq1, seq2)
			result2 := compareByAge(seq2, seq1)

			So(result1, ShouldEqual, 1)
			So(result2, ShouldEqual, -1)
		})

		Convey("Should work with different inputs", func() {
			seq1 := slices.Values([]Person{{"Alice", 30}})
			seq2 := slices.Values([]Person{{"Bob", 25}})

			result := compareByAge(seq1, seq2)
			So(result, ShouldEqual, 1)
		})
	})
}

func TestCompareByKey(t *testing.T) {
	Convey("CompareByKey", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		extractAge := func(p Person) int { return p.Age }

		Convey("Should compare using key extraction function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 20}, {"David", 35}})

			result := CompareByKey(seq1, seq2, extractAge)
			So(result, ShouldEqual, 1) // Alice(25) > Charlie(20)
		})

		Convey("Should return 0 for equal sequences", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 25}, {"David", 30}})

			result := CompareByKey(seq1, seq2, extractAge)
			So(result, ShouldEqual, 0)
		})

		Convey("Should handle empty sequences", func() {
			empty := slices.Values([]Person{})
			nonEmpty := slices.Values([]Person{{"Alice", 25}})

			result1 := CompareByKey(empty, empty, extractAge)
			result2 := CompareByKey(empty, nonEmpty, extractAge)
			result3 := CompareByKey(nonEmpty, empty, extractAge)

			So(result1, ShouldEqual, 0)
			So(result2, ShouldEqual, -1)
			So(result3, ShouldEqual, 1)
		})

		Convey("Should work with different key types", func() {
			extractName := func(p Person) string { return p.Name }

			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 20}, {"David", 35}})

			result := CompareByKey(seq1, seq2, extractName)
			So(result, ShouldEqual, -1) // "Alice" < "Charlie"
		})
	})
}

func TestEqual(t *testing.T) {
	Convey("Equal", t, func() {
		Convey("Should return true for equal sequences", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 3})

			result := Equal(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return false for different sequences", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 4})

			result := Equal(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return false for sequences of different lengths", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2})

			result := Equal(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should handle empty sequences", func() {
			empty := slices.Values([]int{})
			nonEmpty := slices.Values([]int{1, 2, 3})

			result1 := Equal(empty, empty)
			result2 := Equal(empty, nonEmpty)

			So(result1, ShouldBeTrue)
			So(result2, ShouldBeFalse)
		})
	})
}

func TestLessThan(t *testing.T) {
	Convey("LessThan", t, func() {
		Convey("Should return true when first sequence is less", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 4})

			result := LessThan(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return false when sequences are equal", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 3})

			result := LessThan(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return false when first sequence is greater", func() {
			seq1 := slices.Values([]int{1, 2, 4})
			seq2 := slices.Values([]int{1, 2, 3})

			result := LessThan(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return true when first sequence is shorter", func() {
			seq1 := slices.Values([]int{1, 2})
			seq2 := slices.Values([]int{1, 2, 3})

			result := LessThan(seq1, seq2)
			So(result, ShouldBeTrue)
		})
	})
}

func TestGreaterThan(t *testing.T) {
	Convey("GreaterThan", t, func() {
		Convey("Should return true when first sequence is greater", func() {
			seq1 := slices.Values([]int{1, 2, 4})
			seq2 := slices.Values([]int{1, 2, 3})

			result := GreaterThan(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return false when sequences are equal", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 3})

			result := GreaterThan(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return false when first sequence is less", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 4})

			result := GreaterThan(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return true when second sequence is shorter", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2})

			result := GreaterThan(seq1, seq2)
			So(result, ShouldBeTrue)
		})
	})
}

func TestNotEqual(t *testing.T) {
	Convey("NotEqual", t, func() {
		Convey("Should return false for equal sequences", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 3})

			result := NotEqual(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return true for different sequences", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 4})

			result := NotEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return true for sequences of different lengths", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2})

			result := NotEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})
	})
}

func TestLessOrEqual(t *testing.T) {
	Convey("LessOrEqual", t, func() {
		Convey("Should return true when first sequence is less", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 4})

			result := LessOrEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return true when sequences are equal", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 3})

			result := LessOrEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return false when first sequence is greater", func() {
			seq1 := slices.Values([]int{1, 2, 4})
			seq2 := slices.Values([]int{1, 2, 3})

			result := LessOrEqual(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return true when first sequence is shorter", func() {
			seq1 := slices.Values([]int{1, 2})
			seq2 := slices.Values([]int{1, 2, 3})

			result := LessOrEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})
	})
}

func TestGreaterOrEqual(t *testing.T) {
	Convey("GreaterOrEqual", t, func() {
		Convey("Should return true when first sequence is greater", func() {
			seq1 := slices.Values([]int{1, 2, 4})
			seq2 := slices.Values([]int{1, 2, 3})

			result := GreaterOrEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return true when sequences are equal", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 3})

			result := GreaterOrEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})

		Convey("Should return false when first sequence is less", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2, 4})

			result := GreaterOrEqual(seq1, seq2)
			So(result, ShouldBeFalse)
		})

		Convey("Should return true when second sequence is shorter", func() {
			seq1 := slices.Values([]int{1, 2, 3})
			seq2 := slices.Values([]int{1, 2})

			result := GreaterOrEqual(seq1, seq2)
			So(result, ShouldBeTrue)
		})
	})
}

func TestEqualBy(t *testing.T) {
	Convey("EqualBy", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		sameAge := func(l, r Person) bool { return l.Age == r.Age }

		Convey("Should compare using custom equality function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 25}, {"David", 30}})

			result := EqualBy(seq1, seq2, sameAge)
			So(result, ShouldBeTrue)
		})

		Convey("Should return false for different sequences", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 20}, {"David", 35}})

			result := EqualBy(seq1, seq2, sameAge)
			So(result, ShouldBeFalse)
		})

		Convey("Should handle empty sequences", func() {
			empty := slices.Values([]Person{})
			nonEmpty := slices.Values([]Person{{"Alice", 25}})

			result1 := EqualBy(empty, empty, sameAge)
			result2 := EqualBy(empty, nonEmpty, sameAge)

			So(result1, ShouldBeTrue)
			So(result2, ShouldBeFalse)
		})
	})
}

func TestEqualByFunc(t *testing.T) {
	Convey("EqualByFunc", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		sameAge := EqualByFunc(func(l, r Person) bool { return l.Age == r.Age })

		Convey("Should create reusable equality function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 25}, {"David", 30}})

			result1 := sameAge(seq1, seq2)
			result2 := sameAge(seq2, seq1)

			So(result1, ShouldBeTrue)
			So(result2, ShouldBeTrue)
		})

		Convey("Should work with different inputs", func() {
			seq1 := slices.Values([]Person{{"Alice", 30}})
			seq2 := slices.Values([]Person{{"Bob", 25}})

			result := sameAge(seq1, seq2)
			So(result, ShouldBeFalse)
		})
	})
}

func TestEqualByKey(t *testing.T) {
	Convey("EqualByKey", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		extractAge := func(p Person) int { return p.Age }

		Convey("Should compare using key extraction function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 25}, {"David", 30}})

			result := EqualByKey(seq1, seq2, extractAge)
			So(result, ShouldBeTrue)
		})

		Convey("Should return false for different sequences", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 20}, {"David", 35}})

			result := EqualByKey(seq1, seq2, extractAge)
			So(result, ShouldBeFalse)
		})

		Convey("Should handle empty sequences", func() {
			empty := slices.Values([]Person{})
			nonEmpty := slices.Values([]Person{{"Alice", 25}})

			result1 := EqualByKey(empty, empty, extractAge)
			result2 := EqualByKey(empty, nonEmpty, extractAge)

			So(result1, ShouldBeTrue)
			So(result2, ShouldBeFalse)
		})
	})
}

func TestEqualByKeyFunc(t *testing.T) {
	Convey("EqualByKeyFunc", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		sameAge := EqualByKeyFunc(func(p Person) int { return p.Age })

		Convey("Should create reusable key-based equality function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 25}, {"David", 30}})

			result1 := sameAge(seq1, seq2)
			result2 := sameAge(seq2, seq1)

			So(result1, ShouldBeTrue)
			So(result2, ShouldBeTrue)
		})

		Convey("Should work with different inputs", func() {
			seq1 := slices.Values([]Person{{"Alice", 30}})
			seq2 := slices.Values([]Person{{"Bob", 25}})

			result := sameAge(seq1, seq2)
			So(result, ShouldBeFalse)
		})
	})
}

func TestCompareByKeyFunc(t *testing.T) {
	Convey("CompareByKeyFunc", t, func() {
		type Person struct {
			Name string
			Age  int
		}

		compareByAge := CompareByKeyFunc(func(p Person) int { return p.Age })

		Convey("Should create reusable key-based comparison function", func() {
			seq1 := slices.Values([]Person{{"Alice", 25}, {"Bob", 30}})
			seq2 := slices.Values([]Person{{"Charlie", 20}, {"David", 35}})

			result1 := compareByAge(seq1, seq2)
			result2 := compareByAge(seq2, seq1)

			So(result1, ShouldEqual, 1)
			So(result2, ShouldEqual, -1)
		})

		Convey("Should work with different inputs", func() {
			seq1 := slices.Values([]Person{{"Alice", 30}})
			seq2 := slices.Values([]Person{{"Bob", 25}})

			result := compareByAge(seq1, seq2)
			So(result, ShouldEqual, 1)
		})
	})
}
