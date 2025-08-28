package art_test

import (
	"fmt"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art"
)

// ExampleTree_basic demonstrates basic tree operations including insert, search, and iteration.
func ExampleTree_basic() {
	// Create a new arena for memory management
	a := new(arena.Arena)

	// Create a tree that stores string values
	tree := &art.Tree[string]{}

	// Insert some key-value pairs
	tree.Insert(a, []byte("apple"), "red fruit")
	tree.Insert(a, []byte("banana"), "yellow fruit")
	tree.Insert(a, []byte("cherry"), "red berry")

	// Search for a value
	if value := tree.Search([]byte("apple")); value != nil {
		fmt.Printf("Found: %s\n", *value)
	}

	// Get the tree size
	fmt.Printf("Tree size: %d\n", tree.Len())

	// Iterate over all key-value pairs using callbacks
	tree.Visit(func(key []byte, value *string) bool {
		fmt.Printf("Key: %s, Value: %s\n", string(key), *value)

		return false // Continue iteration
	})

	// Output:
	// Found: red fruit
	// Tree size: 3
	// Key: apple, Value: red fruit
	// Key: banana, Value: yellow fruit
	// Key: cherry, Value: red berry
}

// ExampleTree_prefix demonstrates prefix-based operations and iteration.
func ExampleTree_prefix() {
	a := new(arena.Arena)

	tree := &art.Tree[string]{}

	// Insert hierarchical keys
	tree.Insert(a, []byte("user:1"), "Alice")
	tree.Insert(a, []byte("user:2"), "Bob")
	tree.Insert(a, []byte("user:1:name"), "Alice Smith")
	tree.Insert(a, []byte("user:1:email"), "alice@example.com")
	tree.Insert(a, []byte("config:database"), "postgres")
	tree.Insert(a, []byte("config:cache"), "redis")

	// Find all user-related keys
	fmt.Println("User keys:")
	tree.VisitPrefix([]byte("user:"), func(key []byte, value *string) bool {
		fmt.Printf("  %s -> %s\n", string(key), *value)

		return false // Continue iteration
	})

	// Find all configuration keys
	fmt.Println("Config keys:")
	tree.VisitPrefix([]byte("config:"), func(key []byte, value *string) bool {
		fmt.Printf("  %s -> %s\n", string(key), *value)

		return false // Continue iteration
	})

	// Output:
	// User keys:
	//   user:1 -> Alice
	//   user:1:email -> alice@example.com
	//   user:1:name -> Alice Smith
	//   user:2 -> Bob
	// Config keys:
	//   config:cache -> redis
	//   config:database -> postgres
}

// ExampleTree_minMax demonstrates finding minimum and maximum keys in the tree.
func ExampleTree_minMax() {
	a := new(arena.Arena)

	tree := &art.Tree[int]{}

	// Insert some numeric values
	tree.Insert(a, []byte("zebra"), 100)
	tree.Insert(a, []byte("ant"), 1)
	tree.Insert(a, []byte("cat"), 50)
	tree.Insert(a, []byte("dog"), 75)

	// Find the minimum key (lexicographically smallest)
	if min := tree.Minimum(); min != nil {
		fmt.Printf("Minimum key: %s (value: %d)\n", string(min.Key.Raw()), min.Value)
	}

	// Find the maximum key (lexicographically largest)
	if max := tree.Maximum(); max != nil {
		fmt.Printf("Maximum key: %s (value: %d)\n", string(max.Key.Raw()), max.Value)
	}

	// Output:
	// Minimum key: ant (value: 1)
	// Maximum key: zebra (value: 100)
}

// ExampleTree_differentTypes demonstrates using the tree with different value types.
func ExampleTree_differentTypes() {
	a := new(arena.Arena)

	// Tree with integer values
	intTree := &art.Tree[int]{}
	intTree.Insert(a, []byte("count"), 42)
	intTree.Insert(a, []byte("max"), 100)

	// Tree with struct values
	type User struct {
		ID   int
		Name string
	}

	userTree := &art.Tree[User]{}
	userTree.Insert(a, []byte("user:1"), User{ID: 1, Name: "Alice"})
	userTree.Insert(a, []byte("user:2"), User{ID: 2, Name: "Bob"})

	// Search and display results
	if count := intTree.Search([]byte("count")); count != nil {
		fmt.Printf("Count: %d\n", *count)
	}

	if user := userTree.Search([]byte("user:1")); user != nil {
		fmt.Printf("User: %+v\n", *user)
	}

	// Output:
	// Count: 42
	// User: {ID:1 Name:Alice}
}

// ExampleTree_insertNoReplace demonstrates inserting without replacing existing values.
func ExampleTree_insertNoReplace() {
	a := new(arena.Arena)

	tree := &art.Tree[string]{}

	// Insert a value
	tree.Insert(a, []byte("config"), "original")

	// Try to insert without replacing
	if existing := tree.InsertNoReplace(a, []byte("config"), "new"); existing != nil {
		fmt.Printf("Key already exists, keeping: %s\n", *existing)
	} else {
		fmt.Println("New value inserted")
	}

	// Try to insert a new key
	if existing := tree.InsertNoReplace(a, []byte("newkey"), "value"); existing != nil {
		fmt.Printf("Key already exists, keeping: %s\n", *existing)
	} else {
		fmt.Println("New value inserted")
	}

	// Output:
	// Key already exists, keeping: original
	// New value inserted
}

// ExampleTree_delete demonstrates deleting values from the tree.
func ExampleTree_delete() {
	a := new(arena.Arena)

	tree := &art.Tree[string]{}

	// Insert some values
	tree.Insert(a, []byte("apple"), "red")
	tree.Insert(a, []byte("banana"), "yellow")
	tree.Insert(a, []byte("cherry"), "red")

	fmt.Printf("Before deletion: %d items\n", tree.Len())

	// Delete a value
	if deleted := tree.Delete(a, []byte("banana")); deleted != nil {
		fmt.Printf("Deleted: %s\n", *deleted)
	} else {
		fmt.Println("Key not found for deletion")
	}

	fmt.Printf("After deletion: %d items\n", tree.Len())

	// Try to delete a non-existent key
	if deleted := tree.Delete(a, []byte("nonexistent")); deleted != nil {
		fmt.Printf("Deleted: %s\n", *deleted)
	} else {
		fmt.Println("Key not found for deletion")
	}

	// Output:
	// Before deletion: 3 items
	// Deleted: yellow
	// After deletion: 2 items
	// Key not found for deletion
}
