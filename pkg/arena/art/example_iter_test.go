//go:build go1.23

package art_test

import (
	"fmt"
	"strings"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art"
)

// ExampleTree_go123Iterators demonstrates using Go 1.23+ iterators for efficient iteration.
// This example requires Go 1.23 or later to compile and run.
func ExampleTree_go123Iterators() {
	a := new(arena.Arena)

	tree := &art.Tree[int]{}

	// Insert some values
	tree.Insert(a, []byte("a"), 1)
	tree.Insert(a, []byte("b"), 2)
	tree.Insert(a, []byte("c"), 3)

	// Iterate over all key-value pairs using Go 1.23+ iterators
	fmt.Println("All key-value pairs:")

	for key, value := range tree.All() {
		fmt.Printf("  %s -> %d\n", string(key), *value)
	}

	// Iterate over keys with a specific prefix
	fmt.Println("Keys starting with 'a':")

	for key, value := range tree.AllPrefix([]byte("a")) {
		fmt.Printf("  %s -> %d\n", string(key), *value)
	}

	// Output:
	// All key-value pairs:
	//   a -> 1
	//   b -> 2
	//   c -> 3
	// Keys starting with 'a':
	//   a -> 1
}

// ExampleTree_earlyTermination demonstrates early termination during iteration.
func ExampleTree_earlyTermination() {
	a := new(arena.Arena)

	tree := &art.Tree[string]{}

	// Insert many values
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		tree.Insert(a, []byte(key), value)
	}

	var found string

	// Find the first key that contains "50"
	for key, value := range tree.All() {
		if strings.Contains(string(key), "50") {
			found = *value

			break
		}
	}

	fmt.Printf("Found value containing '50': %s\n", found)

	// Output:
	// Found value containing '50': value50
}
