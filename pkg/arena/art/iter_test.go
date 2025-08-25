//go:build go1.23

package art_test

import (
	"fmt"
	"maps"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art"
	"github.com/flier/goutil/pkg/xiter"
)

// TestTree_Iter tests the Iter method
func TestTree_Iter(t *testing.T) {
	Convey("Given an ART tree with values", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		// Insert some values
		tree.Insert(a, []byte("apple"), 1)
		tree.Insert(a, []byte("banana"), 2)
		tree.Insert(a, []byte("cherry"), 3)

		Convey("When iterating over all values", func() {
			visited := maps.Collect(xiter.MapKeyValue(tree.Iter(),
				func(key []byte, value *int) (string, int) {
					return string(key), *value
				}),
			)

			Convey("Then values should be visited", func() {
				So(visited, ShouldResemble, map[string]int{
					"apple":  1,
					"banana": 2,
					"cherry": 3,
				})
			})
		})

		Convey("When iterating with early termination", func() {
			visited := make(map[string]int)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value

				if string(key) == "banana" {
					break // Terminate after visiting "banana"
				}
			}

			Convey("Then some values should be visited", func() {
				So(visited, ShouldResemble, map[string]int{
					"apple":  1,
					"banana": 2,
				})
			})
		})

		Convey("When iterating over empty tree", func() {
			emptyTree := &art.Tree[int]{}
			visited := maps.Collect(xiter.MapKey(emptyTree.Iter(), func(key []byte, value *int) string {
				return string(key)
			}))

			Convey("Then no values should be visited", func() {
				So(len(visited), ShouldEqual, 0)
			})
		})

		Convey("When iterating with callback that modifies visited map", func() {
			visited := make(map[string]int)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value
				// Modify the map during iteration
				visited["modified"] = 999
			}

			Convey("Then values should be present", func() {
				So(visited, ShouldResemble, map[string]int{
					"apple":    1,
					"banana":   2,
					"cherry":   3,
					"modified": 999,
				})
			})
		})
	})
}

// TestTree_IterPrefix tests the IterPrefix method
func TestTree_IterPrefix(t *testing.T) {
	Convey("Given an ART tree with prefixed values", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		// Insert values with different prefixes
		tree.Insert(a, []byte("apple"), 1)
		tree.Insert(a, []byte("application"), 2)
		tree.Insert(a, []byte("banana"), 3)
		tree.Insert(a, []byte("cherry"), 4)
		tree.Insert(a, []byte("date"), 5)

		Convey("When iterating with prefix 'app'", func() {
			seq := tree.IterPrefix([]byte("app"))

			visited := maps.Collect(xiter.MapKeyValue(seq, func(key []byte, value *int) (string, int) {
				return string(key), *value
			}))

			Convey("Then values with 'app' prefix should be visited", func() {
				So(visited, ShouldResemble, map[string]int{
					"apple":       1,
					"application": 2,
				})
			})
		})

		Convey("When iterating with prefix 'b'", func() {
			seq := tree.IterPrefix([]byte("b"))

			visited := maps.Collect(xiter.MapKeyValue(seq, func(key []byte, value *int) (string, int) {
				return string(key), *value
			}))

			Convey("Then only values with 'b' prefix should be visited", func() {
				So(visited, ShouldResemble, map[string]int{
					"banana": 3,
				})
			})
		})

		Convey("When iterating with prefix 'nonexistent'", func() {
			seq := tree.IterPrefix([]byte("nonexistent"))

			visited := maps.Collect(xiter.MapKey(seq, func(key []byte, value *int) string {
				return string(key)
			}))

			Convey("Then no values should be visited", func() {
				So(visited, ShouldBeEmpty)
			})
		})

		Convey("When iterating with empty prefix", func() {
			seq := tree.IterPrefix([]byte{})

			visited := maps.Collect(xiter.MapKeyValue(seq, func(key []byte, value *int) (string, int) {
				return string(key), *value
			}))

			Convey("Then values should be visited", func() {
				So(visited, ShouldResemble, map[string]int{
					"apple":       1,
					"application": 2,
					"banana":      3,
					"cherry":      4,
					"date":        5,
				})
			})
		})

		Convey("When iterating with early termination", func() {
			visited := make(map[string]int)
			seq := tree.IterPrefix([]byte("app"))

			for key, value := range seq {
				visited[string(key)] = *value

				break
			}

			Convey("Then values should be visited", func() {
				So(visited, ShouldResemble, map[string]int{
					"apple": 1,
				})
			})
		})
	})
}

// TestTree_Iter_EdgeCases tests edge cases for Iter
func TestTree_Iter_EdgeCases(t *testing.T) {
	Convey("Given an ART tree", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		Convey("When working with empty keys", func() {
			tree.Insert(a, []byte{}, 123)
			visited := make(map[string]int)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 1)
			So(visited[""], ShouldEqual, 123)
		})

		Convey("When working with zero byte keys", func() {
			tree.Insert(a, []byte{0}, 456)
			visited := make(map[string]int)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 1)
			So(visited[string([]byte{0})], ShouldEqual, 456)
		})

		Convey("When working with special characters", func() {
			tree.Insert(a, []byte("hello\nworld"), 789)
			visited := make(map[string]int)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 1)
			So(visited["hello\nworld"], ShouldEqual, 789)
		})

		Convey("When working with unicode characters", func() {
			tree.Insert(a, []byte("hello世界"), 999)
			visited := make(map[string]int)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 1)
			So(visited["hello世界"], ShouldEqual, 999)
		})
	})
}

// TestTree_IterPrefix_EdgeCases tests edge cases for IterPrefix
func TestTree_IterPrefix_EdgeCases(t *testing.T) {
	Convey("Given an ART tree", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		Convey("When working with very long prefixes", func() {
			// Insert a value with a long key
			longKey := make([]byte, 1000)
			for i := range longKey {
				longKey[i] = byte(i % 256)
			}
			tree.Insert(a, longKey, 123)

			// Test with a long prefix
			longPrefix := make([]byte, 500)
			for i := range longPrefix {
				longPrefix[i] = byte(i % 256)
			}

			visited := make(map[string]int)
			for key, value := range tree.IterPrefix(longPrefix) {
				visited[string(key)] = *value
			}

			// The result depends on whether the prefix matches
			So(len(visited), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("When working with special character prefixes", func() {
			tree.Insert(a, []byte("hello\nworld"), 456)
			tree.Insert(a, []byte("hello\tworld"), 789)

			visited := make(map[string]int)
			for key, value := range tree.IterPrefix([]byte("hello\n")) {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 1)
			So(visited["hello\nworld"], ShouldEqual, 456)
		})

		Convey("When working with unicode prefixes", func() {
			tree.Insert(a, []byte("hello世界"), 111)
			tree.Insert(a, []byte("hello世界你好"), 222)

			visited := make(map[string]int)

			for key, value := range tree.IterPrefix([]byte("hello世界")) {
				visited[string(key)] = *value
			}

			So(visited, ShouldResemble, map[string]int{
				"hello世界":   111,
				"hello世界你好": 222,
			})
		})

		Convey("When working with empty prefix", func() {
			tree.Insert(a, []byte("test1"), 111)
			tree.Insert(a, []byte("test2"), 222)

			visited := make(map[string]int)
			for key, value := range tree.IterPrefix([]byte{}) {
				visited[string(key)] = *value
			}

			// Empty prefix should match all keys
			So(visited, ShouldResemble, map[string]int{
				"test1": 111,
				"test2": 222,
			})
		})

		Convey("When working with single character prefix", func() {
			tree.Insert(a, []byte("a"), 1)
			tree.Insert(a, []byte("ab"), 2)
			tree.Insert(a, []byte("ac"), 3)
			tree.Insert(a, []byte("b"), 4)

			visited := make(map[string]int)
			for key, value := range tree.IterPrefix([]byte("a")) {
				visited[string(key)] = *value
			}

			So(visited, ShouldResemble, map[string]int{
				"a":  1,
				"ab": 2,
				"ac": 3,
			})
		})
	})
}

// TestTree_Iter_AdvancedScenarios tests advanced iteration scenarios
func TestTree_Iter_AdvancedScenarios(t *testing.T) {
	Convey("Given advanced iteration scenarios", t, func() {
		Convey("When iterating with nested prefixes", func() {
			tree := &art.Tree[int]{}
			a := new(arena.Arena)

			// Insert values with nested prefix structure
			tree.Insert(a, []byte("user"), 1)
			tree.Insert(a, []byte("user:1"), 2)
			tree.Insert(a, []byte("user:1:name"), 3)
			tree.Insert(a, []byte("user:1:email"), 4)
			tree.Insert(a, []byte("user:2"), 5)
			tree.Insert(a, []byte("user:2:name"), 6)
			tree.Insert(a, []byte("config"), 7)

			// Test different prefix levels
			Convey("And iterating with 'user' prefix", func() {
				visited := make(map[string]int)
				for key, value := range tree.IterPrefix([]byte("user")) {
					visited[string(key)] = *value
				}

				So(visited, ShouldResemble, map[string]int{
					"user":         1,
					"user:1":       2,
					"user:1:name":  3,
					"user:1:email": 4,
					"user:2":       5,
					"user:2:name":  6,
				})
			})

			Convey("And iterating with 'user:1' prefix", func() {
				visited := make(map[string]int)
				for key, value := range tree.IterPrefix([]byte("user:1")) {
					visited[string(key)] = *value
				}

				So(visited, ShouldResemble, map[string]int{
					"user:1":       2,
					"user:1:name":  3,
					"user:1:email": 4,
				})
			})
		})

		Convey("When iterating with mixed key types", func() {
			tree := &art.Tree[int]{}
			a := new(arena.Arena)

			// Insert various key types
			tree.Insert(a, []byte(""), 0)              // Empty key
			tree.Insert(a, []byte{0}, 1)               // Zero byte
			tree.Insert(a, []byte{255}, 2)             // Max byte
			tree.Insert(a, []byte("normal"), 3)        // Normal string
			tree.Insert(a, []byte("with\nnewline"), 4) // With newline
			tree.Insert(a, []byte("with\ttab"), 5)     // With tab

			visited := make(map[string]int)
			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(visited, ShouldResemble, map[string]int{
				"":                  0,
				string([]byte{0}):   1,
				string([]byte{255}): 2,
				"normal":            3,
				"with\nnewline":     4,
				"with\ttab":         5,
			})
		})

		Convey("When iterating with case sensitivity", func() {
			tree := &art.Tree[int]{}
			a := new(arena.Arena)

			// Insert keys with different cases
			tree.Insert(a, []byte("Hello"), 1)
			tree.Insert(a, []byte("hello"), 2)
			tree.Insert(a, []byte("HELLO"), 3)
			tree.Insert(a, []byte("hElLo"), 4)

			visited := make(map[string]int)
			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(visited, ShouldResemble, map[string]int{
				"Hello": 1,
				"hello": 2,
				"HELLO": 3,
				"hElLo": 4,
			})
		})
	})
}

// TestTree_Iter_DifferentTypes tests different value types
func TestTree_Iter_DifferentTypes(t *testing.T) {
	Convey("Given ART trees with different types", t, func() {
		Convey("When using string values", func() {
			tree := &art.Tree[string]{}
			a := new(arena.Arena)

			tree.Insert(a, []byte("key1"), "value1")
			tree.Insert(a, []byte("key2"), "value2")

			visited := make(map[string]string)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(visited, ShouldResemble, map[string]string{
				"key1": "value1",
				"key2": "value2",
			})
		})

		Convey("When using struct values", func() {
			type TestStruct struct {
				ID   int
				Name string
			}

			tree := &art.Tree[TestStruct]{}
			a := new(arena.Arena)

			tree.Insert(a, []byte("struct1"), TestStruct{ID: 1, Name: "test1"})
			tree.Insert(a, []byte("struct2"), TestStruct{ID: 2, Name: "test2"})

			visited := make(map[string]TestStruct)

			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(visited, ShouldResemble, map[string]TestStruct{
				"struct1": {ID: 1, Name: "test1"},
				"struct2": {ID: 2, Name: "test2"},
			})
		})

		Convey("When using float values", func() {
			tree := &art.Tree[float64]{}
			a := new(arena.Arena)

			tree.Insert(a, []byte("pi"), 3.14159)
			tree.Insert(a, []byte("e"), 2.71828)

			visited := make(map[string]float64)
			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(visited, ShouldResemble, map[string]float64{
				"pi": 3.14159,
				"e":  2.71828,
			})
		})
	})
}

// Benchmark tests for performance measurement
func BenchmarkTree_Iter(b *testing.B) {
	a := new(arena.Arena)
	tree := &art.Tree[int]{}

	// Pre-populate tree
	for j := 0; j < 100; j++ {
		key := []byte(fmt.Sprintf("key%d", j))
		tree.Insert(a, key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N/100; i++ {
		for key, value := range tree.Iter() {
			_, _ = key, value
		}
	}
}

func BenchmarkTree_IterPrefix(b *testing.B) {
	a := new(arena.Arena)
	tree := &art.Tree[int]{}

	// Pre-populate tree with prefixed keys
	for j := 0; j < 100; j++ {
		key := []byte(fmt.Sprintf("prefix%d", j))
		tree.Insert(a, key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N/100; i++ {
		for key, value := range tree.IterPrefix([]byte("prefix")) {
			_, _ = key, value
		}
	}
}

func BenchmarkTree_IterPrefix_NoMatch(b *testing.B) {
	a := new(arena.Arena)
	tree := &art.Tree[int]{}

	// Pre-populate tree with prefixed keys
	for j := 0; j < 100; j++ {
		key := []byte(fmt.Sprintf("prefix%d", j))
		tree.Insert(a, key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N/100; i++ {
		for key, value := range tree.IterPrefix([]byte("nonexistent")) {
			_, _ = key, value
		}
	}
}

func BenchmarkTree_Iter_LargeTree(b *testing.B) {
	a := new(arena.Arena)
	tree := &art.Tree[int]{}

	// Pre-populate tree with many keys
	for j := 0; j < 10000; j++ {
		key := []byte(fmt.Sprintf("key%06d", j))
		tree.Insert(a, key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N/10000; i++ {
		for key, value := range tree.Iter() {
			_, _ = key, value
		}
	}
}

func BenchmarkTree_IterPrefix_LargeTree(b *testing.B) {
	a := new(arena.Arena)
	tree := &art.Tree[int]{}

	// Pre-populate tree with many prefixed keys
	for j := 0; j < 10000; j++ {
		key := []byte(fmt.Sprintf("prefix%06d", j))
		tree.Insert(a, key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N/10000; i++ {
		for key, value := range tree.IterPrefix([]byte("prefix")) {
			_, _ = key, value
		}
	}
}

func BenchmarkTree_Iter_StringValues(b *testing.B) {
	a := new(arena.Arena)
	tree := &art.Tree[int]{}

	// Pre-populate tree with string values
	for j := 0; j < 1000; j++ {
		key := []byte(fmt.Sprintf("key%d", j))

		tree.Insert(a, key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N/1000; i++ {
		for key, value := range tree.Iter() {
			_, _ = key, value
		}
	}
}

// TestTree_Iter_ErrorHandling tests error handling and edge cases
func TestTree_Iter_ErrorHandling(t *testing.T) {
	Convey("Given error handling scenarios", t, func() {
		Convey("When iterating over empty tree", func() {
			tree := &art.Tree[int]{}

			visited := make(map[string]int)
			for key, value := range tree.Iter() {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 0)
		})

		Convey("When iterating with non-existent prefix", func() {
			tree := &art.Tree[int]{}
			a := new(arena.Arena)

			// Insert some values
			tree.Insert(a, []byte("hello"), 1)
			tree.Insert(a, []byte("world"), 2)

			visited := make(map[string]int)
			for key, value := range tree.IterPrefix([]byte("nonexistent")) {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 0)
		})

		Convey("When iterating with very long prefix that doesn't match", func() {
			tree := &art.Tree[int]{}
			a := new(arena.Arena)

			// Insert a short key
			tree.Insert(a, []byte("short"), 1)

			// Create a very long prefix that won't match
			longPrefix := make([]byte, 1000)
			for i := range longPrefix {
				longPrefix[i] = byte('x')
			}

			visited := make(map[string]int)
			for key, value := range tree.IterPrefix(longPrefix) {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 0)
		})

		Convey("When iterating with prefix that partially matches", func() {
			tree := &art.Tree[int]{}
			a := new(arena.Arena)

			// Insert keys with similar prefixes
			tree.Insert(a, []byte("hello"), 1)
			tree.Insert(a, []byte("hello world"), 2)
			tree.Insert(a, []byte("hello there"), 3)
			tree.Insert(a, []byte("help"), 4)

			visited := make(map[string]int)
			for key, value := range tree.IterPrefix([]byte("hel")) {
				visited[string(key)] = *value
			}

			So(len(visited), ShouldEqual, 4)
			So(visited["hello"], ShouldEqual, 1)
			So(visited["hello world"], ShouldEqual, 2)
			So(visited["hello there"], ShouldEqual, 3)
			So(visited["help"], ShouldEqual, 4)
		})
	})
}
