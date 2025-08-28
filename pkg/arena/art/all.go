//go:build go1.23

package art

import (
	"iter"

	"github.com/flier/goutil/pkg/arena/art/tree"
)

// All iterates over all key-value pairs in the tree using Go 1.23+ iterators.
//
// This method provides a modern, ergonomic way to iterate over the entire tree
// using Go's built-in range syntax. The iteration follows lexicographic key
// ordering, making it suitable for sorted traversal and range operations.
//
// The returned iterator is lazy and only processes elements as they are requested,
// making it memory-efficient for large trees. The iteration can be interrupted
// early by using break statements or return statements within the range loop.
//
// Returns:
//   - An iter.Seq2[[]byte, *T] that can be used with Go's range syntax
//
// Performance: O(n) where n is the number of key-value pairs in the tree
//
// Example:
//
//	// Iterate over all key-value pairs
//	for key, value := range tree.All() {
//	    fmt.Printf("Key: %s, Value: %v\n", string(key), *value)
//	}
//
//	// Early termination
//	for key, value := range tree.All() {
//	    if string(key) == "stop" {
//	        break
//	    }
//	    fmt.Printf("Processing: %s -> %v\n", string(key), *value)
//	}
//
//	// Collect into a map
//	visited := make(map[string]int)
//	for key, value := range tree.All() {
//	    visited[string(key)] = *value
//	}
//
// Note: This method requires Go 1.23 or later due to the use of iter.Seq2.
// For compatibility with earlier Go versions, use the Visit method instead.
func (t *Tree[T]) All() iter.Seq2[[]byte, *T] {
	return func(yield func([]byte, *T) bool) {
		tree.RecursiveIter(t.root, func(key []byte, value *T) bool {
			return !yield(key, value)
		})
	}
}

// AllPrefix iterates over key-value pairs with a specific prefix using Go 1.23+ iterators.
//
// This method provides efficient prefix-based iteration using Go's modern iterator
// syntax. It's particularly useful for implementing hierarchical data structures,
// configuration systems, or any scenario where you need to process keys that share
// a common prefix.
//
// The iteration maintains lexicographic ordering within the prefix matches and
// can be interrupted early using break statements or return statements within
// the range loop.
//
// Parameters:
//   - prefix: The byte slice representing the prefix to match against. Keys that
//     start with this prefix will be included in the iteration.
//
// Returns:
//   - An iter.Seq2[[]byte, *T] that can be used with Go's range syntax
//
// Performance: O(k + m) where k is the prefix length and m is the number of matching keys
//
// Example:
//
//	// Find all user-related keys
//	for key, user := range tree.AllPrefix([]byte("user:")) {
//	    fmt.Printf("User key: %s, data: %+v\n", string(key), *user)
//	}
//
//	// Process configuration keys
//	for key, config := range tree.AllPrefix([]byte("config.")) {
//	    if strings.Contains(string(key), "database") {
//	        fmt.Printf("Database config: %s -> %v\n", string(key), *config)
//	    }
//	}
//
//	// Early termination for performance
//	count := 0
//	for key, value := range tree.AllPrefix([]byte("temp.")) {
//	    if count >= 100 {
//	        break // Stop after processing 100 items
//	    }
//	    processTemporaryData(key, value)
//	    count++
//	}
//
//	// Collect prefix matches into a slice
//	var matches []string
//	for key, value := range tree.AllPrefix([]byte("api.")) {
//	    matches = append(matches, fmt.Sprintf("%s:%v", string(key), *value))
//	}
//
// Note: This method requires Go 1.23 or later due to the use of iter.Seq2.
// For compatibility with earlier Go versions, use the VisitPrefix method instead.
func (t *Tree[T]) AllPrefix(prefix []byte) iter.Seq2[[]byte, *T] {
	return func(yield func([]byte, *T) bool) {
		tree.IterPrefix(t.root, prefix, func(key []byte, value *T) bool {
			return !yield(key, value)
		})
	}
}
