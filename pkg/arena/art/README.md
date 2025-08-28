# Adaptive Radix Tree (ART) Package

The `goutil/pkg/arena/art` package provides an efficient, memory-optimized implementation of Adaptive Radix Trees (ART) using arena allocation for superior performance and memory management.

## Overview

Adaptive Radix Trees are a space-efficient data structure for storing and retrieving string keys. This implementation automatically adapts its node storage strategy based on the number of children, optimizing both memory usage and lookup performance.

## Key Features

- **Adaptive Node Types**: Automatically switches between Node4, Node16, Node48, and Node256 based on child count
- **Arena Allocation**: Uses efficient memory management with the `goutil/pkg/arena` package
- **Generic Types**: Supports any Go type as values with compile-time type safety
- **Prefix Compression**: Reduces memory overhead by storing common key prefixes only once
- **Lazy Expansion**: Inner nodes are only created when necessary to distinguish between leaf nodes
- **Path Compression**: Removes inner nodes with only a single child
- **Go 1.23+ Support**: Leverages modern Go features like `iter.Seq2` for efficient iteration

## Architecture

### Node Types

The tree automatically adapts its node structure:

- **Node4**: Stores up to 4 children using sorted arrays (most memory efficient)
- **Node16**: Stores up to 16 children using sorted arrays (balanced)
- **Node48**: Stores up to 48 children using sparse arrays (efficient for medium nodes)
- **Node256**: Stores up to 256 children using direct arrays (fastest lookup)

### Memory Management

All memory is managed through arena allocation, providing:
- Efficient bulk memory allocation
- Automatic cleanup when the arena is released
- Better cache locality
- Reduced memory fragmentation

## Usage

### Basic Operations

```go
package main

import (
    "fmt"
    "github.com/flier/goutil/pkg/arena"
    "github.com/flier/goutil/pkg/arena/art"
)

func main() {
    // Create a new arena for memory management
    a := new(arena.Arena)

    // Create a tree that stores string values
    tree := &art.Tree[string]{}

    // Insert key-value pairs
    tree.Insert(a, []byte("hello"), "world")
    tree.Insert(a, []byte("foo"), "bar")
    tree.Insert(a, []byte("test"), "value")

    // Search for values
    if value := tree.Search([]byte("hello")); value != nil {
        fmt.Printf("Found: %s\n", *value) // Output: Found: world
    }

    // Get tree statistics
    fmt.Printf("Tree size: %d\n", tree.Len()) // Output: Tree size: 3

    // Find minimum and maximum keys
    if min := tree.Minimum(); min != nil {
        fmt.Printf("Min key: %s\n", string(min.Key.Raw()))
    }
    if max := tree.Maximum(); max != nil {
        fmt.Printf("Max key: %s\n", string(max.Key.Raw()))
    }
}
```

### Iteration

#### Using Callback Functions

```go
// Visit all key-value pairs
tree.Visit(func(key []byte, value *string) bool {
    fmt.Printf("Key: %s, Value: %s\n", string(key), *value)
    return false // Continue iteration
})

// Visit keys with a specific prefix
tree.VisitPrefix([]byte("h"), func(key []byte, value *string) bool {
    fmt.Printf("Prefix match: %s -> %s\n", string(key), *value)
    return false // Continue iteration
})
```

#### Using Go 1.23+ Iterators

```go
// Iterate over all key-value pairs
for key, value := range tree.All() {
    fmt.Printf("Key: %s, Value: %s\n", string(key), *value)
}

// Iterate over keys with a specific prefix
for key, value := range tree.AllPrefix([]byte("h")) {
    fmt.Printf("Prefix match: %s -> %s\n", string(key), *value)
}
```

### Advanced Operations

```go
// Insert without replacing existing values
if oldValue := tree.InsertNoReplace(a, []byte("hello"), "new world"); oldValue != nil {
    fmt.Printf("Key already exists with value: %s\n", *oldValue)
}

// Delete values
if deletedValue := tree.Delete(a, []byte("foo")); deletedValue != nil {
    fmt.Printf("Deleted: %s\n", *deletedValue)
}

// Check if tree is empty
if tree.Len() == 0 {
    fmt.Println("Tree is empty")
}
```

### Working with Different Value Types

```go
// Integer values
intTree := &art.Tree[int]{}
intTree.Insert(a, []byte("count"), 42)
intTree.Insert(a, []byte("max"), 100)

// Struct values
type User struct {
    ID   int
    Name string
}

userTree := &art.Tree[User]{}
userTree.Insert(a, []byte("user:1"), User{ID: 1, Name: "Alice"})
userTree.Insert(a, []byte("user:2"), User{ID: 2, Name: "Bob"})

// Pointer values
ptrTree := &art.Tree[*User]{}
user1 := &User{ID: 1, Name: "Alice"}
ptrTree.Insert(a, []byte("user:1"), user1)
```

## Performance Characteristics

### Time Complexity

- **Search**: O(k) where k is the key length
- **Insert**: O(k) where k is the key length
- **Delete**: O(k) where k is the key length
- **Prefix Search**: O(k + m) where k is prefix length and m is number of matching keys

### Space Complexity

- **Node4**: 4 bytes + 4 pointers + overhead
- **Node16**: 16 bytes + 16 pointers + overhead
- **Node48**: 256 bytes + 48 pointers + overhead
- **Node256**: 256 pointers + overhead

### Memory Optimization Features

- **Prefix Compression**: Common key prefixes are stored only once
- **Lazy Expansion**: Inner nodes are created only when necessary
- **Path Compression**: Single-child nodes are eliminated
- **Arena Allocation**: Efficient bulk memory management

## Best Practices

### Memory Management

1. **Use a single arena per tree**: This provides the best memory efficiency
2. **Batch operations**: Perform multiple operations within the same arena context

```go
a := new(arena.Arena)

tree := &art.Tree[string]{}

// Perform all operations
for i := 0; i < 1000; i++ {
    key := fmt.Sprintf("key%d", i)
    tree.Insert(a, []byte(key), fmt.Sprintf("value%d", i))
}
```

### Key Design

1. **Use meaningful prefixes**: Keys with common prefixes benefit from compression
2. **Avoid very long keys**: While supported, very long keys increase memory usage
3. **Consider key ordering**: Sequential keys can improve cache locality

### Performance Tuning

1. **Monitor node type distribution**: Use tree inspection tools to understand node distribution
2. **Balance memory vs. speed**: Node4 is most memory-efficient, Node256 is fastest
3. **Profile your use case**: Different access patterns may benefit from different optimizations

## Examples

### Database Index

```go
type Record struct {
    ID   int
    Data string
}

// Create an index on the Data field
index := &art.Tree[Record]{}
a := new(arena.Arena)

// Index records
records := []Record{
    {ID: 1, Data: "apple"},
    {ID: 2, Data: "banana"},
    {ID: 3, Data: "cherry"},
}

for _, record := range records {
    index.Insert(a, []byte(record.Data), record)
}

// Search for records with prefix "a"
index.VisitPrefix([]byte("a"), func(key []byte, record *Record) bool {
    fmt.Printf("Found: %+v\n", *record)
    return false
})
```

### URL Router

```go
type Handler func(w http.ResponseWriter, r *http.Request)

router := &art.Tree[Handler]{}
a := new(arena.Arena)

// Register routes
router.Insert(a, []byte("/api/users"), handleUsers)
router.Insert(a, []byte("/api/users/"), handleUserDetails)
router.Insert(a, []byte("/api/posts"), handlePosts)

// Find handler for request
func findHandler(path []byte) Handler {
    if handler := router.Search(path); handler != nil {
        return *handler
    }

    // Try prefix matching for wildcard routes
    var found Handler
    router.VisitPrefix(path, func(key []byte, h *Handler) bool {
        found = *h
        return true // Stop at first match
    })

    return found
}
```

### Configuration Store

```go
type Config struct {
    Value     string
    Timestamp time.Time
}

configStore := &art.Tree[Config]{}
a := new(arena.Arena)

// Store hierarchical configuration
configStore.Insert(a, []byte("database.host"), Config{Value: "localhost", Timestamp: time.Now()})
configStore.Insert(a, []byte("database.port"), Config{Value: "5432", Timestamp: time.Now()})
configStore.Insert(a, []byte("database.name"), Config{Value: "myapp", Timestamp: time.Now()})

// Retrieve all database configuration
configStore.VisitPrefix([]byte("database."), func(key []byte, config *Config) bool {
    fmt.Printf("%s: %s\n", string(key), config.Value)
    return false
})
```

## Benchmarks

The package includes comprehensive benchmarks for all operations:

```bash
# Run all benchmarks
go test ./pkg/arena/art -bench=.

# Run specific benchmarks
go test ./pkg/arena/art -bench=BenchmarkTree_Insert
go test ./pkg/arena/art -bench=BenchmarkTree_Search
go test ./pkg/arena/art -bench=BenchmarkTree_Visit
```

## Contributing

When contributing to this package:

1. **Maintain performance**: All changes should preserve or improve performance
2. **Add tests**: New features must include comprehensive tests
3. **Update documentation**: Keep this README and code comments up to date
4. **Follow Go conventions**: Use standard Go formatting and naming conventions

## License

This package is part of the `goutil` project and follows the same license terms.
