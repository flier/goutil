//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleIsSorted() {
	fmt.Println(IsSorted(slices.Values([]int{1, 2, 3})))
	fmt.Println(IsSorted(slices.Values([]int{1, 3, 2})))

	// Output:
	// true
	// false
}

func ExampleIsSortedBy() {
	type User struct {
		Name string
		Age  int
	}

	s := slices.Values([]User{{"tom", 8}, {"joe", 12}})

	fmt.Println(IsSortedBy(s, func(x, y User) bool { return x.Age <= y.Age }))

	// Output:
	// true
}

func ExampleIsSortedByFunc() {
	type User struct {
		Name string
		Age  int
	}

	sortedByAge := IsSortedByFunc(func(x, y User) bool { return x.Age <= y.Age })

	s := slices.Values([]User{{"tom", 8}, {"joe", 12}})

	fmt.Println(sortedByAge(s))

	// Output:
	// true
}

func ExampleIsSortedByKey() {
	type User struct {
		Name string
		Age  int
	}

	s := slices.Values([]User{{"joe", 12}, {"tom", 8}})

	fmt.Println(IsSortedByKey(s, func(u User) string { return u.Name }))
	fmt.Println(IsSortedByKey(s, func(u User) int { return u.Age }))

	// Output:
	// true
	// false
}

func ExampleIsSortedByKeyFunc() {
	type User struct {
		Name string
		Age  int
	}

	sortedByAge := IsSortedByKeyFunc(func(u User) int { return u.Age })

	s := slices.Values([]User{{"joe", 12}, {"tom", 8}})

	fmt.Println(sortedByAge(s))

	// Output:
	// false
}
