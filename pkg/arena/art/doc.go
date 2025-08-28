// Package art provides an efficient, memory-optimized implementation of Adaptive Radix Trees (ART).
//
// Adaptive Radix Trees are a space-efficient data structure for storing and retrieving string keys.
// This implementation automatically adapts its node storage strategy based on the number of children,
// optimizing both memory usage and lookup performance.
//
// # Overview
//
// ART trees are designed to provide excellent performance characteristics for string key operations
// while maintaining reasonable memory usage. The tree automatically switches between different node
// types (Node4, Node16, Node48, Node256) based on the number of children at each node, ensuring
// optimal performance for various data distributions.
//
// # Key Features
//
//   - **Adaptive Node Types**: Automatically switches between Node4, Node16, Node48, and Node256
//     based on child count for optimal memory usage and performance
//   - **Arena Allocation**: Uses efficient memory management with the `goutil/pkg/arena` package
//     for better cache locality and reduced allocation overhead
//   - **Generic Types**: Supports any Go type as values with compile-time type safety using
//     Go's generics (`[T any]`)
//   - **Prefix Compression**: Reduces memory overhead by sharing common prefixes between nodes
//   - **Lazy Expansion**: Inner nodes are only created when required to distinguish between
//     at least two leaf nodes
//   - **Path Compression**: Removes inner nodes with only a single child to reduce memory usage
//   - **Go 1.23+ Support**: Optional support for Go 1.23+ iterators for more ergonomic iteration
//
// # Performance Characteristics
//
//   - **Lookup**: O(k) where k is the key length, with excellent cache locality
//   - **Insertion**: O(k) with automatic node type adaptation
//   - **Deletion**: O(k) with automatic node type shrinking
//   - **Memory**: Adaptive to data distribution, typically 2-4x more space-efficient than hash tables
//   - **Cache Performance**: Optimized for CPU cache utilization with arena allocation
//
// # Node Types
//
// The tree automatically adapts its node storage strategy:
//
//   - **Node4**: For nodes with 1-4 children, using sorted arrays for optimal small-node performance
//   - **Node16**: For nodes with 5-16 children, using SIMD-optimized search for medium-sized nodes
//   - **Node48**: For nodes with 17-48 children, using indirect indexing for space efficiency
//   - **Node256**: For nodes with 49-256 children, using direct array indexing for large nodes
//
// # Memory Management
//
// The package uses arena allocation for efficient memory management:
//
//   - **Arena Allocation**: Memory is allocated from large, pre-allocated blocks
//   - **Batch Operations**: All memory in an arena is freed together when the arena is reset
//   - **Cache Locality**: Related data is stored together, improving CPU cache performance
//   - **Reduced Fragmentation**: Eliminates memory fragmentation common with individual allocations
//
// # Usage Patterns
//
// ## Basic Operations
//
//	arena := new(arena.Arena)
//	defer arena.Reset()
//
//	tree := &art.Tree[string]{}
//	tree.Insert(arena, []byte("key"), "value")
//
//	if value := tree.Search([]byte("key")); value != nil {
//	    fmt.Printf("Found: %s\n", *value)
//	}
//
// ## Iteration
//
//	// Iterate over all key-value pairs
//	tree.Visit(func(key []byte, value *string) bool {
//	    fmt.Printf("%s -> %s\n", string(key), *value)
//	    return false // Continue iteration
//	})
//
//	// Iterate over keys with a specific prefix
//	tree.VisitPrefix([]byte("user:"), func(key []byte, value *string) bool {
//	    fmt.Printf("%s -> %s\n", string(key), *value)
//	    return false
//	})
//
// ## Go 1.23+ Iterators (Optional)
//
//	// For Go 1.23+ users, more ergonomic iteration is available
//	for key, value := range tree.All() {
//	    fmt.Printf("%s -> %s\n", string(key), *value)
//	}
//
//	for key, value := range tree.AllPrefix([]byte("user:")) {
//	    fmt.Printf("%s -> %s\n", string(key), *value)
//	}
//
// # When to Use
//
// ART trees are most beneficial for:
//
//   - **String Key Storage**: Applications requiring efficient string key lookups
//   - **Prefix Operations**: Scenarios needing prefix-based searches and iterations
//   - **Memory-Constrained Environments**: Where memory efficiency is important
//   - **High-Performance Applications**: That can benefit from arena allocation
//   - **Sorted Data Requirements**: When lexicographic ordering is needed
//
// # Alternatives
//
// Consider alternatives based on your specific needs:
//
//   - **Hash Tables**: For simple key-value storage without ordering requirements
//   - **B-Trees**: For disk-based storage or when memory is not a constraint
//   - **Tries**: For simpler prefix operations without adaptive optimization
//   - **Standard Go Maps**: For general-purpose key-value storage
//
// # Thread Safety
//
// The Tree type is not thread-safe. If multiple goroutines access the same tree
// concurrently, external synchronization must be provided by the caller.
//
// # Memory Safety
//
//   - All memory allocated through the arena must not be accessed after calling `arena.Reset()`
//   - The tree and all its nodes become invalid after the arena is reset
//   - Use defer statements to ensure proper cleanup: `defer arena.Reset()`
//
// # Examples
//
// See the `example_test.go` file for comprehensive usage examples covering:
//
//   - Basic tree operations (insert, search, delete)
//   - Prefix-based operations and iteration
//   - Finding minimum and maximum keys
//   - Working with different value types
//   - Early termination during iteration
//   - Go 1.23+ iterator usage
//
// # References
//
//   - [The Adaptive Radix Tree: ARTful Indexing for Main-Memory Databases](https://db.in.tum.de/~leis/papers/ART.pdf)
//   - [Go Arena Allocation](https://go.dev/blog/arena)
//   - [Go Generics](https://go.dev/doc/tutorial/generics)
package art
