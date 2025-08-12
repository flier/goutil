//go:build go1.23

package xiter_test

import (
	"cmp"
	"fmt"
	"slices"

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
