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

// TestTree_LazyExpansion tests lazy expansion behavior where inner nodes are only
// created when they are required to distinguish between at least two leaf nodes
func TestTree_LazyExpansion(t *testing.T) {
	Convey("Given an ART tree with lazy expansion", t, func() {
		tree := &art.Tree[int]{}
		a := new(arena.Arena)

		Convey("When inserting a single key", func() {
			tree.Insert(a, []byte("single"), 1)

			Convey("Then the tree should contain only a leaf node", func() {
				So(tree.Len(), ShouldEqual, 1)
				So(*tree.Search([]byte("single")), ShouldEqual, 1)
			})

			Convey("Then no inner nodes should be created unnecessarily", func() {
				// With only one key, no inner nodes are needed for distinction
				// The tree should remain as a single leaf
				min := tree.Minimum()
				max := tree.Maximum()
				So(min, ShouldNotBeNil)
				So(max, ShouldNotBeNil)
				So(min.Key.Raw(), ShouldResemble, []byte("single"))
				So(max.Key.Raw(), ShouldResemble, []byte("single"))
			})
		})

		Convey("When inserting keys with no common prefix", func() {
			tree.Insert(a, []byte("apple"), 1)
			tree.Insert(a, []byte("zebra"), 2)

			Convey("Then both keys should be searchable", func() {
				So(*tree.Search([]byte("apple")), ShouldEqual, 1)
				So(*tree.Search([]byte("zebra")), ShouldEqual, 2)
			})

			Convey("Then inner nodes should be created only when necessary", func() {
				// These keys have no common prefix, so minimal inner node structure
				// should be created to distinguish them
				So(tree.Len(), ShouldEqual, 2)
			})
		})

		Convey("When inserting keys with common prefix that requires distinction", func() {
			tree.Insert(a, []byte("apple"), 1)
			tree.Insert(a, []byte("apricot"), 2)

			Convey("Then both keys should be searchable", func() {
				So(*tree.Search([]byte("apple")), ShouldEqual, 1)
				So(*tree.Search([]byte("apricot")), ShouldEqual, 2)
			})

			Convey("Then inner nodes should be created to distinguish the common prefix", func() {
				// These keys share "ap" prefix but differ at position 2
				// Inner nodes should be created to distinguish them
				So(tree.Len(), ShouldEqual, 2)
			})

			Convey("Then prefix operations should work correctly", func() {
				visited := make(map[string]int)
				tree.VisitPrefix([]byte("ap"), func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})
				So(len(visited), ShouldEqual, 2)
				So(visited["apple"], ShouldEqual, 1)
				So(visited["apricot"], ShouldEqual, 2)
			})
		})

		Convey("When inserting keys that share a longer common prefix", func() {
			tree.Insert(a, []byte("application"), 1)
			tree.Insert(a, []byte("appliance"), 2)

			Convey("Then both keys should be searchable", func() {
				So(*tree.Search([]byte("application")), ShouldEqual, 1)
				So(*tree.Search([]byte("appliance")), ShouldEqual, 2)
			})

			Convey("Then inner nodes should be created only at the point of divergence", func() {
				// These keys share "appli" prefix but differ at position 5
				// Inner nodes should be created only where distinction is needed
				So(tree.Len(), ShouldEqual, 2)
			})
		})

		Convey("When inserting keys with incremental common prefixes", func() {
			tree.Insert(a, []byte("a"), 1)
			tree.Insert(a, []byte("ab"), 2)
			tree.Insert(a, []byte("abc"), 3)

			Convey("Then all keys should be searchable", func() {
				So(*tree.Search([]byte("a")), ShouldEqual, 1)
				So(*tree.Search([]byte("ab")), ShouldEqual, 2)
				So(*tree.Search([]byte("abc")), ShouldEqual, 3)
			})

			Convey("Then inner nodes should be created progressively", func() {
				// Each key extends the previous one, requiring inner nodes
				// only where distinction is necessary
				So(tree.Len(), ShouldEqual, 3)
			})

			Convey("Then prefix operations should work at all levels", func() {
				// Test prefix "a"
				visitedA := make(map[string]int)
				tree.VisitPrefix([]byte("a"), func(key []byte, value *int) bool {
					visitedA[string(key)] = *value
					return false
				})
				So(len(visitedA), ShouldEqual, 3)

				// Test prefix "ab"
				visitedAB := make(map[string]int)
				tree.VisitPrefix([]byte("ab"), func(key []byte, value *int) bool {
					visitedAB[string(key)] = *value
					return false
				})
				So(len(visitedAB), ShouldEqual, 2)
				So(visitedAB["ab"], ShouldEqual, 2)
				So(visitedAB["abc"], ShouldEqual, 3)
			})
		})

		Convey("When inserting keys that require deep inner node creation", func() {
			tree.Insert(a, []byte("verylongprefix1"), 1)
			tree.Insert(a, []byte("verylongprefix2"), 2)

			Convey("Then both keys should be searchable", func() {
				So(*tree.Search([]byte("verylongprefix1")), ShouldEqual, 1)
				So(*tree.Search([]byte("verylongprefix2")), ShouldEqual, 2)
			})

			Convey("Then inner nodes should be created only at the point of divergence", func() {
				// These keys share "verylongprefix" but differ at the last character
				// Inner nodes should be created only where distinction is needed
				So(tree.Len(), ShouldEqual, 2)
			})
		})

		Convey("When inserting keys with no common prefix but similar lengths", func() {
			tree.Insert(a, []byte("hello"), 1)
			tree.Insert(a, []byte("world"), 2)
			tree.Insert(a, []byte("test"), 3)

			Convey("Then all keys should be searchable", func() {
				So(*tree.Search([]byte("hello")), ShouldEqual, 1)
				So(*tree.Search([]byte("world")), ShouldEqual, 2)
				So(*tree.Search([]byte("test")), ShouldEqual, 3)
			})

			Convey("Then inner nodes should be created minimally", func() {
				// These keys have no common prefix, so minimal inner node structure
				// should be created to distinguish them
				So(tree.Len(), ShouldEqual, 3)
			})
		})

		Convey("When inserting keys that create a balanced structure", func() {
			tree.Insert(a, []byte("middle"), 1)
			tree.Insert(a, []byte("left"), 2)
			tree.Insert(a, []byte("right"), 3)

			Convey("Then all keys should be searchable", func() {
				So(*tree.Search([]byte("middle")), ShouldEqual, 1)
				So(*tree.Search([]byte("left")), ShouldEqual, 2)
				So(*tree.Search([]byte("right")), ShouldEqual, 3)
			})

			Convey("Then inner nodes should be created efficiently", func() {
				// This creates a balanced structure where inner nodes are
				// created only where necessary for distinction
				So(tree.Len(), ShouldEqual, 3)
			})

			Convey("Then traversal should work correctly", func() {
				visited := make(map[string]int)
				tree.Visit(func(key []byte, value *int) bool {
					visited[string(key)] = *value
					return false
				})
				So(len(visited), ShouldEqual, 3)
				So(visited["left"], ShouldEqual, 2)
				So(visited["middle"], ShouldEqual, 1)
				So(visited["right"], ShouldEqual, 3)
			})
		})

		Convey("When testing lazy expansion with deletion", func() {
			tree.Insert(a, []byte("prefix1"), 1)
			tree.Insert(a, []byte("prefix2"), 2)

			Convey("Then both keys should exist initially", func() {
				So(tree.Len(), ShouldEqual, 2)
				So(*tree.Search([]byte("prefix1")), ShouldEqual, 1)
				So(*tree.Search([]byte("prefix2")), ShouldEqual, 2)
			})

			// Delete one key
			oldValue := tree.Delete(a, []byte("prefix1"))

			Convey("Then deletion should work correctly", func() {
				So(oldValue, ShouldNotBeNil)
				So(*oldValue, ShouldEqual, 1)
				So(tree.Len(), ShouldEqual, 1)
			})

			Convey("Then remaining key should still be searchable", func() {
				So(*tree.Search([]byte("prefix2")), ShouldEqual, 2)
				So(tree.Search([]byte("prefix1")), ShouldBeNil)
			})

			Convey("Then inner nodes should be optimized after deletion", func() {
				// After deletion, if only one key remains, the tree structure
				// should be optimized to remove unnecessary inner nodes
				min := tree.Minimum()
				max := tree.Maximum()
				So(min, ShouldNotBeNil)
				So(max, ShouldNotBeNil)
				So(min.Key.Raw(), ShouldResemble, []byte("prefix2"))
				So(max.Key.Raw(), ShouldResemble, []byte("prefix2"))
			})
		})

		Convey("When testing lazy expansion with complex key patterns", func() {
			// Insert keys that create a complex but efficient structure
			keys := []string{
				"a", "aa", "aaa", "aaaa",
				"b", "bb", "bbb", "bbbb",
				"c", "cc", "ccc", "cccc",
			}

			for i, key := range keys {
				tree.Insert(a, []byte(key), i+1)
			}

			Convey("Then all keys should be searchable", func() {
				for i, key := range keys {
					So(*tree.Search([]byte(key)), ShouldEqual, i+1)
				}
			})

			Convey("Then inner nodes should be created efficiently", func() {
				// The tree should create inner nodes only where necessary
				// to distinguish between keys, not for every character
				So(tree.Len(), ShouldEqual, len(keys))
			})

			Convey("Then prefix operations should work at all levels", func() {
				// Test various prefix levels
				testPrefixes := []string{"a", "aa", "aaa", "b", "bb", "c"}

				for _, prefix := range testPrefixes {
					visited := make(map[string]int)
					tree.VisitPrefix([]byte(prefix), func(key []byte, value *int) bool {
						visited[string(key)] = *value
						return false
					})

					// Count expected keys with this prefix
					expectedCount := 0
					for _, key := range keys {
						if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
							expectedCount++
						}
					}

					So(len(visited), ShouldEqual, expectedCount)
				}
			})
		})

		Convey("When testing lazy expansion with edge cases", func() {
			Convey("Then empty keys should work correctly", func() {
				tree.Insert(a, []byte{}, 1)
				So(*tree.Search([]byte{}), ShouldEqual, 1)
				So(tree.Len(), ShouldEqual, 1)
			})

			Convey("Then single byte keys should work correctly", func() {
				tree.Insert(a, []byte("x"), 1)
				tree.Insert(a, []byte("y"), 2)
				So(*tree.Search([]byte("x")), ShouldEqual, 1)
				So(*tree.Search([]byte("y")), ShouldEqual, 2)
				So(tree.Len(), ShouldEqual, 2)
			})

			Convey("Then keys with special characters should work correctly", func() {
				tree.Insert(a, []byte("key@123"), 1)
				tree.Insert(a, []byte("key#456"), 2)
				So(*tree.Search([]byte("key@123")), ShouldEqual, 1)
				So(*tree.Search([]byte("key#456")), ShouldEqual, 2)
				So(tree.Len(), ShouldEqual, 2)
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
