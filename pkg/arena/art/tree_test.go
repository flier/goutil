package art_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art"
)

var (
	kHello  = []byte("hello")
	kKey    = []byte("key")
	kPrefix = []byte("prefix")
)

// TestTree_BasicOperations tests basic tree operations
func TestTree_BasicOperations(t *testing.T) {
	Convey("Given a new ART tree", t, func() {
		a := new(arena.Arena)
		tree := &art.Tree[int]{}

		Convey("When the tree is empty", func() {
			Convey("Then Len should return 0", func() {
				So(tree.Len(), ShouldEqual, 0)
			})

			Convey("Then Search should return nil", func() {
				result := tree.Search(kKey)
				So(result, ShouldBeNil)
			})

			Convey("Then Minimum should return nil", func() {
				result := tree.Minimum()
				So(result, ShouldBeNil)
			})

			Convey("Then Maximum should return nil", func() {
				result := tree.Maximum()
				So(result, ShouldBeNil)
			})

			Convey("Then Visit should not call callback", func() {
				visited := make(map[string]int)
				result := tree.Visit(func(key []byte, value *int) bool {
					visited[string(key)] = *value

					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 0)
			})

			Convey("Then VisitPrefix should not call callback", func() {
				visited := make(map[string]int)
				result := tree.VisitPrefix(kPrefix, func(key []byte, value *int) bool {
					visited[string(key)] = *value

					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 0)
			})
		})

		Convey("When inserting a single value", func() {
			oldValue := tree.Insert(a, kHello, 123)

			Convey("Then Len should return 1", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			Convey("Then Insert should return nil (no old value)", func() {
				So(oldValue, ShouldBeNil)
			})

			Convey("Then Search should find the value", func() {
				result := tree.Search(kHello)
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 123)
			})

			Convey("Then Search with non-existent key should return nil", func() {
				result := tree.Search([]byte("world"))
				So(result, ShouldBeNil)
			})

			Convey("Then Minimum should return the leaf", func() {
				min := tree.Minimum()
				So(min, ShouldNotBeNil)
				So(min.Key.Raw(), ShouldResemble, kHello)
				So(min.Value, ShouldEqual, 123)
			})

			Convey("Then Maximum should return the leaf", func() {
				max := tree.Maximum()
				So(max, ShouldNotBeNil)
				So(max.Key.Raw(), ShouldResemble, kHello)
				So(max.Value, ShouldEqual, 123)
			})

			Convey("Then Visit should call callback", func() {
				visited := make(map[string]int)
				result := tree.Visit(func(key []byte, value *int) bool {
					visited[string(key)] = *value

					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 1)
				So(visited["hello"], ShouldEqual, 123)
			})

			Convey("Then VisitPrefix with matching prefix should call callback", func() {
				visited := make(map[string]int)
				result := tree.VisitPrefix([]byte("hel"), func(key []byte, value *int) bool {
					visited[string(key)] = *value

					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 1)
				So(visited["hello"], ShouldEqual, 123)
			})

			Convey("Then VisitPrefix with non-matching prefix should not call callback", func() {
				visited := make(map[string]int)
				result := tree.VisitPrefix([]byte("wor"), func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 0)
			})

			Convey("Then VisitPrefix with empty prefix should call callback", func() {
				visited := make(map[string]int)
				result := tree.VisitPrefix([]byte{}, func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 1)
				So(visited["hello"], ShouldEqual, 123)
			})
		})
	})
}

// TestTree_InsertOperations tests insert operations
func TestTree_InsertOperations(t *testing.T) {
	Convey("Given an ART tree", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		Convey("When inserting multiple values", func() {
			So(tree.Insert(a, []byte("apple"), 1), ShouldBeNil)
			So(tree.Insert(a, []byte("banana"), 2), ShouldBeNil)
			So(tree.Insert(a, []byte("cherry"), 3), ShouldBeNil)

			Convey("Then Len should return 3", func() {
				So(tree.Len(), ShouldEqual, 3)
			})

			Convey("Then all values should be searchable", func() {
				So(*tree.Search([]byte("apple")), ShouldEqual, 1)
				So(*tree.Search([]byte("banana")), ShouldEqual, 2)
				So(*tree.Search([]byte("cherry")), ShouldEqual, 3)
			})

			Convey("Then Minimum should return the first alphabetically", func() {
				min := tree.Minimum()
				So(min, ShouldNotBeNil)
				So(min.Key.Raw(), ShouldResemble, []byte("apple"))
				So(min.Value, ShouldEqual, 1)
			})

			Convey("Then Maximum should return the last alphabetically", func() {
				max := tree.Maximum()
				So(max, ShouldNotBeNil)
				So(max.Key.Raw(), ShouldResemble, []byte("cherry"))
				So(max.Value, ShouldEqual, 3)
			})

			Convey("Then Visit should visit all values", func() {
				visited := make(map[string]int)
				result := tree.Visit(func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 3)
				So(visited, ShouldResemble, map[string]int{
					"apple":  1,
					"banana": 2,
					"cherry": 3,
				})
			})

			Convey("Then VisitPrefix with 'a' should visit apple", func() {
				visited := make(map[string]int)
				result := tree.VisitPrefix([]byte("a"), func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 1)
				So(visited, ShouldResemble, map[string]int{
					"apple": 1,
				})
			})

			Convey("Then VisitPrefix with 'b' should visit banana", func() {
				visited := make(map[string]int)
				result := tree.VisitPrefix([]byte("b"), func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 1)
				So(visited, ShouldResemble, map[string]int{
					"banana": 2,
				})
			})
		})

		Convey("When inserting with replace", func() {
			// Insert initial value
			tree.Insert(a, kKey, 100)

			Convey("Then Len should be 1 after first insert", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			// Replace with new value
			oldValue := tree.Insert(a, kKey, 200)

			Convey("Then Len should remain 1 after replace", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			Convey("Then Insert should return old value", func() {
				So(oldValue, ShouldNotBeNil)
				So(*oldValue, ShouldEqual, 100)
			})

			Convey("Then Search should return new value", func() {
				result := tree.Search(kKey)
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 200)
			})
		})

		Convey("When inserting with InsertNoReplace", func() {
			// Insert initial value
			tree.Insert(a, kKey, 100)

			Convey("Then Len should be 1 after first insert", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			// Try to insert without replace
			oldValue := tree.InsertNoReplace(a, kKey, 200)

			Convey("Then Len should remain 1 after InsertNoReplace", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			Convey("Then InsertNoReplace should return old value", func() {
				So(oldValue, ShouldNotBeNil)
				So(*oldValue, ShouldEqual, 100)
			})

			Convey("Then Search should return old value", func() {
				result := tree.Search(kKey)
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 100)
			})
		})
	})
}

// TestTree_DeleteOperations tests delete operations
func TestTree_DeleteOperations(t *testing.T) {
	Convey("Given an ART tree with values", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		// Insert some values
		tree.Insert(a, []byte("apple"), 1)
		tree.Insert(a, []byte("banana"), 2)
		tree.Insert(a, []byte("cherry"), 3)

		Convey("When deleting an existing key", func() {
			Convey("Then Len should be 3 before deletion", func() {
				So(tree.Len(), ShouldEqual, 3)
			})

			oldValue := tree.Delete(a, []byte("banana"))

			Convey("Then Len should be 2 after deletion", func() {
				So(tree.Len(), ShouldEqual, 2)
			})

			Convey("Then Delete should return the old value", func() {
				So(oldValue, ShouldNotBeNil)
				So(*oldValue, ShouldEqual, 2)
			})

			Convey("Then Search should not find the deleted key", func() {
				result := tree.Search([]byte("banana"))
				So(result, ShouldBeNil)
			})

			Convey("Then other keys should still be searchable", func() {
				So(*tree.Search([]byte("apple")), ShouldEqual, 1)
				So(*tree.Search([]byte("cherry")), ShouldEqual, 3)
			})

			Convey("Then Visit should not visit deleted key", func() {
				visited := make(map[string]int)
				tree.Visit(func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(len(visited), ShouldEqual, 2)
				So(visited, ShouldResemble, map[string]int{
					"apple":  1,
					"cherry": 3,
				})
			})
		})

		Convey("When deleting a non-existent key", func() {
			Convey("Then Len should be 3 before deletion attempt", func() {
				So(tree.Len(), ShouldEqual, 3)
			})

			oldValue := tree.Delete(a, []byte("nonexistent"))

			Convey("Then Len should remain 3 after deletion attempt", func() {
				So(tree.Len(), ShouldEqual, 3)
			})

			Convey("Then Delete should return nil", func() {
				So(oldValue, ShouldBeNil)
			})

			Convey("Then all existing keys should remain", func() {
				So(*tree.Search([]byte("apple")), ShouldEqual, 1)
				So(*tree.Search([]byte("banana")), ShouldEqual, 2)
				So(*tree.Search([]byte("cherry")), ShouldEqual, 3)
			})
		})

		Convey("When deleting all keys", func() {
			Convey("Then Len should be 3 before any deletions", func() {
				So(tree.Len(), ShouldEqual, 3)
			})

			// Delete all keys
			oldValue1 := tree.Delete(a, []byte("apple"))
			Convey("Then Len should be 2 after first deletion", func() {
				So(tree.Len(), ShouldEqual, 2)
			})

			oldValue2 := tree.Delete(a, []byte("banana"))
			Convey("Then Len should be 1 after second deletion", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			oldValue3 := tree.Delete(a, []byte("cherry"))
			Convey("Then Len should be 0 after third deletion", func() {
				So(tree.Len(), ShouldEqual, 0)
			})

			Convey("Then all deletes should return values", func() {
				So(oldValue1, ShouldNotBeNil)
				So(*oldValue1, ShouldEqual, 1)
				So(oldValue2, ShouldNotBeNil)
				So(*oldValue2, ShouldEqual, 2)
				So(oldValue3, ShouldNotBeNil)
				So(*oldValue3, ShouldEqual, 3)
			})

			Convey("Then tree should be empty", func() {
				So(tree.Search([]byte("apple")), ShouldBeNil)
				So(tree.Search([]byte("banana")), ShouldBeNil)
				So(tree.Search([]byte("cherry")), ShouldBeNil)
				So(tree.Minimum(), ShouldBeNil)
				So(tree.Maximum(), ShouldBeNil)
			})

			Convey("Then Visit should not call callback", func() {
				visited := make(map[string]int)
				result := tree.Visit(func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(result, ShouldBeFalse)
				So(len(visited), ShouldEqual, 0)
			})
		})
	})
}

// TestTree_VisitOperations tests visit operations
func TestTree_VisitOperations(t *testing.T) {
	Convey("Given an ART tree with values", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		// Insert values
		tree.Insert(a, []byte("a"), 1)
		tree.Insert(a, []byte("b"), 2)
		tree.Insert(a, []byte("c"), 3)
		tree.Insert(a, []byte("d"), 4)

		Convey("When visiting with early termination", func() {
			visited := make(map[string]int)
			result := tree.Visit(func(key []byte, value *int) bool {
				visited[string(key)] = *value
				return string(key) == "b" // Terminate after visiting "b"
			})

			Convey("Then Visit should return true", func() {
				So(result, ShouldBeTrue)
			})

			Convey("Then only some values should be visited", func() {
				So(len(visited), ShouldEqual, 2)
				So(visited["a"], ShouldEqual, 1)
				So(visited["b"], ShouldEqual, 2)
				So(visited["c"], ShouldEqual, 0) // Not visited
				So(visited["d"], ShouldEqual, 0) // Not visited
			})
		})

		Convey("When visiting with prefix and early termination", func() {
			visited := make(map[string]int)
			result := tree.VisitPrefix([]byte("a"), func(key []byte, value *int) bool {
				visited[string(key)] = *value
				return true // Terminate immediately
			})

			Convey("Then VisitPrefix should return true", func() {
				So(result, ShouldBeTrue)
			})

			Convey("Then only one value should be visited", func() {
				So(len(visited), ShouldEqual, 1)
				So(visited["a"], ShouldEqual, 1)
			})
		})

		Convey("When visiting with callback that modifies visited map", func() {
			visited := make(map[string]int)
			result := tree.Visit(func(key []byte, value *int) bool {
				visited[string(key)] = *value
				// Modify the map during iteration
				visited["modified"] = 999
				return false
			})

			So(result, ShouldBeFalse)
			So(len(visited), ShouldEqual, 5) // 4 original + 1 modified
			So(visited["modified"], ShouldEqual, 999)
		})
	})
}

// TestTree_EdgeCases tests edge cases
func TestTree_EdgeCases(t *testing.T) {
	Convey("Given an ART tree", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		Convey("When working with empty keys", func() {
			oldValue := tree.Insert(a, []byte{}, 123)
			So(oldValue, ShouldBeNil)

			Convey("Then Search should find the value", func() {
				result := tree.Search([]byte{})
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 123)
			})

			Convey("Then Visit should visit the empty key", func() {
				visited := make(map[string]int)
				tree.Visit(func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})

				So(len(visited), ShouldEqual, 1)
				So(visited[""], ShouldEqual, 123)
			})
		})

		Convey("When working with zero byte keys", func() {
			oldValue := tree.Insert(a, []byte{0}, 456)
			So(oldValue, ShouldBeNil)

			Convey("Then Search should find the value", func() {
				result := tree.Search([]byte{0})
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 456)
			})
		})

		Convey("When working with very long keys", func() {
			longKey := make([]byte, 1000)
			for i := range longKey {
				longKey[i] = byte(i % 256)
			}

			oldValue := tree.Insert(a, longKey, 789)
			So(oldValue, ShouldBeNil)

			Convey("Then Search should find the value", func() {
				result := tree.Search(longKey)
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 789)
			})
		})

		Convey("When working with special characters", func() {
			specialKey := []byte("hello\nworld\t")
			oldValue := tree.Insert(a, specialKey, 999)
			So(oldValue, ShouldBeNil)

			Convey("Then Search should find the value", func() {
				result := tree.Search(specialKey)
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 999)
			})
		})

		Convey("When working with unicode characters", func() {
			unicodeKey := []byte("hello世界")
			oldValue := tree.Insert(a, unicodeKey, 888)
			So(oldValue, ShouldBeNil)

			Convey("Then Search should find the value", func() {
				result := tree.Search(unicodeKey)
				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 888)
			})
		})
	})
}

// TestTree_DifferentTypes tests different value types
func TestTree_DifferentTypes(t *testing.T) {
	Convey("Given ART trees with different types", t, func() {
		Convey("When using string values", func() {
			tree := &art.Tree[string]{}
			a := new(arena.Arena)

			tree.Insert(a, []byte("key1"), "value1")
			tree.Insert(a, []byte("key2"), "value2")

			So(*tree.Search([]byte("key1")), ShouldEqual, "value1")
			So(*tree.Search([]byte("key2")), ShouldEqual, "value2")
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

			result1 := tree.Search([]byte("struct1"))
			So(result1, ShouldNotBeNil)
			So(result1.ID, ShouldEqual, 1)
			So(result1.Name, ShouldEqual, "test1")

			result2 := tree.Search([]byte("struct2"))
			So(result2, ShouldNotBeNil)
			So(result2.ID, ShouldEqual, 2)
			So(result2.Name, ShouldEqual, "test2")
		})

		Convey("When using float values", func() {
			tree := &art.Tree[float64]{}
			a := new(arena.Arena)

			tree.Insert(a, []byte("pi"), 3.14159)
			tree.Insert(a, []byte("e"), 2.71828)

			So(*tree.Search([]byte("pi")), ShouldEqual, 3.14159)
			So(*tree.Search([]byte("e")), ShouldEqual, 2.71828)
		})
	})
}

// TestTree_LenOperations tests the Len method comprehensively
func TestTree_LenOperations(t *testing.T) {
	Convey("Given an ART tree", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		Convey("When the tree is newly created", func() {
			Convey("Then Len should return 0", func() {
				So(tree.Len(), ShouldEqual, 0)
			})
		})

		Convey("When inserting values incrementally", func() {
			Convey("Then Len should be 0 initially", func() {
				So(tree.Len(), ShouldEqual, 0)
			})

			tree.Insert(a, []byte("key1"), 1)
			Convey("Then Len should be 1 after first insert", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			tree.Insert(a, []byte("key2"), 2)
			Convey("Then Len should be 2 after second insert", func() {
				So(tree.Len(), ShouldEqual, 2)
			})

			tree.Insert(a, []byte("key3"), 3)
			Convey("Then Len should be 3 after third insert", func() {
				So(tree.Len(), ShouldEqual, 3)
			})
		})

		Convey("When inserting with duplicate keys", func() {
			tree.Insert(a, []byte("duplicate"), 1)
			Convey("Then Len should be 1 after first insert", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			tree.Insert(a, []byte("duplicate"), 2)
			Convey("Then Len should remain 1 after duplicate insert", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			tree.Insert(a, []byte("duplicate"), 3)
			Convey("Then Len should remain 1 after another duplicate insert", func() {
				So(tree.Len(), ShouldEqual, 1)
			})
		})

		Convey("When inserting with InsertNoReplace", func() {
			tree.Insert(a, []byte("no-replace"), 1)
			Convey("Then Len should be 1 after first insert", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			tree.InsertNoReplace(a, []byte("no-replace"), 2)
			Convey("Then Len should remain 1 after InsertNoReplace", func() {
				So(tree.Len(), ShouldEqual, 1)
			})
		})

		Convey("When deleting values incrementally", func() {
			// Insert some values first
			tree.Insert(a, []byte("del1"), 1)
			tree.Insert(a, []byte("del2"), 2)
			tree.Insert(a, []byte("del3"), 3)

			Convey("Then Len should be 3 before any deletions", func() {
				So(tree.Len(), ShouldEqual, 3)
			})

			tree.Delete(a, []byte("del1"))
			Convey("Then Len should be 2 after first deletion", func() {
				So(tree.Len(), ShouldEqual, 2)
			})

			tree.Delete(a, []byte("del2"))
			Convey("Then Len should be 1 after second deletion", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			tree.Delete(a, []byte("del3"))
			Convey("Then Len should be 0 after third deletion", func() {
				So(tree.Len(), ShouldEqual, 0)
			})
		})

		Convey("When deleting non-existent keys", func() {
			tree.Insert(a, []byte("exists"), 1)
			Convey("Then Len should be 1 before deletion attempt", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			tree.Delete(a, []byte("non-existent"))
			Convey("Then Len should remain 1 after deletion attempt", func() {
				So(tree.Len(), ShouldEqual, 1)
			})
		})

		Convey("When mixing insert and delete operations", func() {
			Convey("Then Len should be 0 initially", func() {
				So(tree.Len(), ShouldEqual, 0)
			})

			tree.Insert(a, []byte("mix1"), 1)
			tree.Insert(a, []byte("mix2"), 2)
			Convey("Then Len should be 2 after two inserts", func() {
				So(tree.Len(), ShouldEqual, 2)
			})

			tree.Delete(a, []byte("mix1"))
			Convey("Then Len should be 1 after one deletion", func() {
				So(tree.Len(), ShouldEqual, 1)
			})

			tree.Insert(a, []byte("mix3"), 3)
			Convey("Then Len should be 2 after another insert", func() {
				So(tree.Len(), ShouldEqual, 2)
			})

			tree.Delete(a, []byte("mix2"))
			tree.Delete(a, []byte("mix3"))
			Convey("Then Len should be 0 after deleting remaining keys", func() {
				So(tree.Len(), ShouldEqual, 0)
			})
		})

		Convey("When working with different value types", func() {
			stringTree := &art.Tree[string]{}
			floatTree := &art.Tree[float64]{}

			Convey("Then string tree Len should be 0 initially", func() {
				So(stringTree.Len(), ShouldEqual, 0)
			})

			stringTree.Insert(a, []byte("str1"), "hello")
			Convey("Then string tree Len should be 1 after insert", func() {
				So(stringTree.Len(), ShouldEqual, 1)
			})

			Convey("Then float tree Len should be 0 initially", func() {
				So(floatTree.Len(), ShouldEqual, 0)
			})

			floatTree.Insert(a, []byte("float1"), 3.14)
			Convey("Then float tree Len should be 1 after insert", func() {
				So(floatTree.Len(), ShouldEqual, 1)
			})
		})

		Convey("When performing rapid insert/delete cycles", func() {
			for i := 0; i < 10; i++ {
				key := []byte(fmt.Sprintf("cycle%d", i))
				tree.Insert(a, key, i)
				So(tree.Len(), ShouldEqual, i+1)
			}

			Convey("Then Len should be 10 after all inserts", func() {
				So(tree.Len(), ShouldEqual, 10)
			})

			for i := 9; i >= 0; i-- {
				key := []byte(fmt.Sprintf("cycle%d", i))
				tree.Delete(a, key)
				So(tree.Len(), ShouldEqual, i)
			}

			Convey("Then Len should be 0 after all deletions", func() {
				So(tree.Len(), ShouldEqual, 0)
			})
		})
	})
}

// Benchmark tests for performance measurement
func BenchmarkTree_Insert(b *testing.B) {
	tree := &art.Tree[int]{}
	a := new(arena.Arena)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		tree.Insert(a, key, i)
	}
}

func BenchmarkTree_Search(b *testing.B) {
	a := new(arena.Arena)
	tree := arena.New(a, art.Tree[int]{})

	// Pre-populate tree
	for i := 0; i < 1000; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		tree.Insert(a, key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i%1000))

		_ = tree.Search(key)
	}
}

func BenchmarkTree_Visit(b *testing.B) {
	tree := &art.Tree[int]{}
	a := new(arena.Arena)

	// Pre-populate tree
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		tree.Insert(a, key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N/100; i++ {
		tree.Visit(func(key []byte, value *int) bool {
			_, _ = key, value

			return false
		})
	}
}

func BenchmarkTree_VisitPrefix(b *testing.B) {
	tree := &art.Tree[int]{}
	a := new(arena.Arena)

	// Pre-populate tree with prefixed keys
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("prefix%d", i))
		tree.Insert(a, key, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N/100; i++ {
		tree.VisitPrefix(kPrefix, func(key []byte, value *int) bool {
			_, _ = key, value

			return false
		})
	}
}
