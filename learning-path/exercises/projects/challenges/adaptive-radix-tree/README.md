# Challenge 1: Adaptive Radix Tree (ART)

## Overview
Implement an Adaptive Radix Tree index with multiple node types, lazy expansion, prefix compression, and iterator support.

**Duration:** 30-40 hours
**Difficulty:** Very Hard

## Background
ART is a space-efficient, cache-friendly trie variant used in high-performance databases. It adapts node sizes based on fan-out.

## Node Types

### Node4 (1-4 children)
```
┌─────────────────────┐
│ Type: Node4         │
│ Keys: [a, b, d, f]  │ (sorted array)
│ Ptrs: [*, *, *, *]  │ (4 child pointers)
└─────────────────────┘
```

### Node16 (5-16 children)
```
┌─────────────────────────────┐
│ Type: Node16                │
│ Keys: [16 bytes]           │
│ Ptrs: [16 pointers]        │
└─────────────────────────────┘
```

### Node48 (17-48 children)
```
┌─────────────────────────────┐
│ Type: Node48                │
│ Index: [256 bytes]         │ (maps key -> child index)
│ Ptrs: [48 pointers]        │
└─────────────────────────────┘
```

### Node256 (49-256 children)
```
┌─────────────────────────────┐
│ Type: Node256               │
│ Ptrs: [256 pointers]       │ (direct lookup)
└─────────────────────────────┘
```

## Core Operations

### Insert
- Start at root
- Follow path, creating nodes as needed
- Grow nodes when full (Node4 -> Node16 -> Node48 -> Node256)
- Use prefix compression

### Lookup
- Follow path using appropriate node type
- Handle prefix compression
- Return value or nil

### Delete
- Remove entry
- Shrink nodes when sparse
- Merge nodes when possible

### Iterator
- Depth-first traversal
- Support range queries
- Lexicographic order

## Optimizations

### 1. Prefix Compression
Store common prefixes in nodes to save space
```
Without compression:
  "testing" -> t -> e -> s -> t -> i -> n -> g

With compression:
  "test" (prefix) -> i -> n -> g
```

### 2. Path Compression
Compress single-child paths
```
a -> b -> c -> d (leaf)
Becomes:
a -> "bcd" (leaf)
```

### 3. SIMD Lookups (Optional)
Use SIMD for searching keys in Node16/Node48

## API Design

```go
type ART struct {
	root *Node
	size int
}

type Node interface {
	Insert(key []byte, value interface{}) (Node, bool)
	Delete(key []byte) (Node, bool)
	Search(key []byte) interface{}
}

type Node4 struct {
	prefix []byte
	keys   [4]byte
	children [4]Node
	numChildren int
}

// Create new ART
func New() *ART

// Insert key-value pair
func (art *ART) Insert(key []byte, value interface{})

// Search for key
func (art *ART) Search(key []byte) (interface{}, bool)

// Delete key
func (art *ART) Delete(key []byte) bool

// Iterator over range
func (art *ART) Range(start, end []byte) iter.Seq2[[]byte, interface{}]

// Size returns number of keys
func (art *ART) Size() int
```

## Test Cases

### Correctness
- Insert and retrieve 1M keys
- Delete keys and verify
- Range queries
- Prefix searches
- Node type transitions

### Performance
- Compare to Go map
- Compare to B-tree
- Measure memory usage
- Cache performance

## Benchmarks

```go
BenchmarkInsert        - Insertion speed
BenchmarkLookup        - Lookup speed
BenchmarkIterate       - Iteration speed
BenchmarkMemory        - Memory per key
BenchmarkVsMap         - Compare to map
BenchmarkVsBTree       - Compare to B-tree
```

## Performance Goals

- Insert: >2M ops/sec
- Lookup: >5M ops/sec
- Memory: <32 bytes per key (avg)
- Range query: >1M keys/sec
- Better than map for string keys

## Implementation Hints

### Node Growth
```go
func (n *Node4) Insert(key byte, child Node) Node {
	if n.numChildren < 4 {
		// Add to Node4
		n.keys[n.numChildren] = key
		n.children[n.numChildren] = child
		n.numChildren++
		return n
	}

	// Grow to Node16
	node16 := &Node16{}
	copy(node16.keys[:], n.keys[:])
	copy(node16.children[:], n.children[:])
	node16.keys[4] = key
	node16.children[4] = child
	node16.numChildren = 5
	return node16
}
```

### Prefix Compression
```go
func longestCommonPrefix(a, b []byte) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}
```

## Stretch Goals

### 1. SIMD Optimizations
Use SIMD for Node16 key search

### 2. Concurrent ART
Optimistic lock coupling for concurrency

### 3. Persistent ART
Make ART durable with WAL

### 4. Bulk Loading
Optimize for sorted insert

## Learning Outcomes
- Advanced tree data structures
- Memory-efficient design
- Cache optimization
- Iterator implementation
- Benchmark-driven development

## References
- "The Adaptive Radix Tree: ARTful Indexing for Main-Memory Databases" (Leis et al.)
- DuckDB ART implementation
- PostgreSQL VACUUM

## Time Estimate
- Core implementation: 20-25 hours
- Optimizations: 5-8 hours
- Testing and benchmarking: 5-7 hours
