# Project 2.1: CSR Graph Implementation

## Overview
Implement a Compressed Sparse Row (CSR) graph representation using Go 1.23 iterators for cache-efficient graph traversal and queries.

**Duration:** 15-18 hours
**Difficulty:** Medium-Hard

## Learning Objectives
- Understand CSR graph representation
- Master Go 1.23 iterators (`iter.Seq`)
- Optimize for cache locality
- Implement iterator composition
- Profile memory layout and performance
- Compare to adjacency list representation

## Concepts Covered
- Compressed Sparse Row format
- Go 1.23 `iter.Seq` and `iter.Seq2`
- Cache-aware data structures
- Iterator patterns and composition
- Memory layout optimization
- Graph traversal algorithms
- Performance profiling with perf

## Requirements

### Core Functionality

#### 1. CSR Graph Structure
```
Adjacency stored in two arrays:
┌─────────────────────────────┐
│  Offsets Array              │
│  [0, 3, 5, 8, ...]          │  offset[i] = start of node i's edges
│                             │  offset[i+1] = end of node i's edges
├─────────────────────────────┤
│  Edges Array                │
│  [1, 2, 3, 0, 4, 2, 5, 6]   │  flattened edge lists
│                             │
└─────────────────────────────┘
```

#### 2. Graph Operations
- **AddNode()** - Add node to graph
- **AddEdge(src, dst)** - Add directed edge
- **Neighbors(node)** - Iterator over neighbors
- **Degree(node)** - Get node degree
- **Has2Hop(src, dst)** - Check 2-hop connectivity
- **Build()** - Finalize graph (create CSR arrays)

#### 3. Iterator Support (Go 1.23)
- Neighbor iteration using `iter.Seq[NodeID]`
- Edge iteration using `iter.Seq2[NodeID, NodeID]`
- Early exit with break
- Iterator composition (filter, map, take)

#### 4. Memory Efficiency
- Compact representation
- No duplicate edge storage
- Minimal metadata overhead
- Cache-friendly access patterns

## Getting Started

```bash
# Initialize module
cd csr-graph
go mod init csrgraph

# Run tests
go test -v
go test -race
go test -cover

# Run benchmarks
go test -bench=. -benchmem

# Memory profiling
go test -bench=BenchmarkIteration -memprofile=mem.prof
go tool pprof mem.prof
```

## API Design

```go
package csrgraph

import "iter"

type NodeID uint32

// CSRGraph represents a directed graph in CSR format
type CSRGraph struct {
	nodeCount uint32
	edgeCount uint32
	offsets   []uint32  // nodeCount + 1 elements
	edges     []NodeID  // edgeCount elements
}

// Builder for constructing CSR graph
type GraphBuilder struct {
	adjList map[NodeID][]NodeID
}

// Create new graph builder
func NewBuilder() *GraphBuilder

// Add node to graph
func (b *GraphBuilder) AddNode(node NodeID)

// Add directed edge
func (b *GraphBuilder) AddEdge(src, dst NodeID)

// Build CSR graph from adjacency list
func (b *GraphBuilder) Build() *CSRGraph

// CSR Graph methods

// Number of nodes
func (g *CSRGraph) NodeCount() uint32

// Number of edges
func (g *CSRGraph) EdgeCount() uint32

// Degree of a node
func (g *CSRGraph) Degree(node NodeID) uint32

// Iterate over neighbors using Go 1.23 iterators
func (g *CSRGraph) Neighbors(node NodeID) iter.Seq[NodeID]

// Iterate over all edges
func (g *CSRGraph) Edges() iter.Seq2[NodeID, NodeID]

// Check if path exists within 2 hops
func (g *CSRGraph) Has2Hop(src, dst NodeID) bool

// Get all 2-hop neighbors
func (g *CSRGraph) TwoHopNeighbors(node NodeID) iter.Seq[NodeID]
```

## Test Cases

### Correctness Tests
- **TestBuild** - Build CSR from adjacency list
- **TestNeighbors** - Iterate over neighbors
- **TestDegree** - Compute node degrees
- **TestEdges** - Iterate over all edges
- **TestEarlyExit** - Break from iteration
- **Test2Hop** - 2-hop connectivity
- **TestEmptyGraph** - Handle empty graph
- **TestSingleNode** - Graph with one node
- **TestIsolatedNodes** - Nodes with no edges

### Iterator Tests
- **TestIteratorComposition** - Chain filter/map
- **TestIteratorBreak** - Early exit works
- **TestMultipleIterators** - Concurrent iterations
- **TestIteratorReuse** - Reuse iterator multiple times

### Performance Tests
- **TestLargeGraph** - 1M nodes, 10M edges
- **TestMemoryUsage** - Compare to adjacency list
- **TestCacheLocality** - Measure cache misses

## Benchmarks

```go
BenchmarkBuild              - Graph construction
BenchmarkNeighborIteration  - Iterate neighbors
BenchmarkDegree             - Degree computation
Benchmark2HopQuery          - 2-hop queries
BenchmarkVsAdjList          - Compare to adjacency list
BenchmarkCacheMisses        - Cache performance (with perf)
```

## Implementation Hints

### CSR Construction
```go
func (b *GraphBuilder) Build() *CSRGraph {
	// Find max node ID
	maxNode := NodeID(0)
	for node := range b.adjList {
		if node > maxNode {
			maxNode = node
		}
	}

	nodeCount := maxNode + 1
	offsets := make([]uint32, nodeCount+1)
	edgeCount := uint32(0)

	// Compute offsets
	for node := NodeID(0); node < nodeCount; node++ {
		offsets[node] = edgeCount
		if neighbors, ok := b.adjList[node]; ok {
			edgeCount += uint32(len(neighbors))
		}
	}
	offsets[nodeCount] = edgeCount

	// Fill edges array
	edges := make([]NodeID, edgeCount)
	idx := 0
	for node := NodeID(0); node < nodeCount; node++ {
		if neighbors, ok := b.adjList[node]; ok {
			copy(edges[idx:], neighbors)
			idx += len(neighbors)
		}
	}

	return &CSRGraph{
		nodeCount: nodeCount,
		edgeCount: edgeCount,
		offsets:   offsets,
		edges:     edges,
	}
}
```

### Neighbor Iterator (Go 1.23)
```go
func (g *CSRGraph) Neighbors(node NodeID) iter.Seq[NodeID] {
	return func(yield func(NodeID) bool) {
		if node >= g.nodeCount {
			return
		}

		start := g.offsets[node]
		end := g.offsets[node+1]

		for i := start; i < end; i++ {
			if !yield(g.edges[i]) {
				return  // early exit
			}
		}
	}
}
```

### 2-Hop Query
```go
func (g *CSRGraph) Has2Hop(src, dst NodeID) bool {
	// Check direct edge
	for neighbor := range g.Neighbors(src) {
		if neighbor == dst {
			return true
		}
	}

	// Check 2-hop paths
	for intermediate := range g.Neighbors(src) {
		for neighbor := range g.Neighbors(intermediate) {
			if neighbor == dst {
				return true
			}
		}
	}

	return false
}
```

### Iterator Composition
```go
// Filter iterator
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if pred(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Map iterator
func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Take first N
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for v := range seq {
			if count >= n {
				return
			}
			if !yield(v) {
				return
			}
			count++
		}
	}
}
```

## Performance Goals

- Build time: < 1s for 1M nodes, 10M edges
- Neighbor iteration: < 5ns per edge
- Degree computation: < 50ns
- 2-hop query: < 1µs (avg degree 10)
- Memory: < 8 bytes per edge (vs 24 for adjacency list)
- Cache miss rate: < 5% for sequential access

## Stretch Goals

### 1. Weighted Edges
- Add edge weights array
- Weighted neighbor iterator
- Shortest path queries

### 2. Bidirectional Graph
- Store both forward and reverse edges
- Reverse neighbor iteration
- Bidirectional BFS

### 3. Node Properties
- Store node attributes efficiently
- Property filtering during iteration
- Columnar property storage

### 4. Graph Compression
- Delta encoding for sorted edges
- Bit packing for small node IDs
- Measure compression ratio

## Common Pitfalls

1. **Off-by-One Errors**
   - Offsets array has nodeCount + 1 elements
   - Check bounds carefully
   - Handle empty neighbor lists

2. **Iterator Misuse**
   - Don't modify graph during iteration
   - Handle early exit correctly
   - Test with break statements

3. **Memory Layout**
   - Ensure arrays are contiguous
   - Avoid pointer chasing
   - Profile cache behavior

4. **Large Graphs**
   - Watch for integer overflow
   - Use appropriate types (uint32 vs uint64)
   - Test with realistic sizes

## Debugging Tips

```bash
# Visualize small graph
go run tools/visualize.go graph.txt

# Profile cache misses (Linux)
perf stat -e cache-misses go test -bench=BenchmarkIteration

# Memory layout
go build -gcflags='-m' graph.go

# Benchmark comparison
benchstat old.txt new.txt
```

## Validation Checklist

Your implementation should:
- [ ] Pass all unit tests
- [ ] Support Go 1.23 iterators
- [ ] Handle early exit with break
- [ ] Work with graphs of 1M+ nodes
- [ ] Use < 8 bytes per edge
- [ ] Have cache-friendly access
- [ ] Beat adjacency list in iteration
- [ ] Support 2-hop queries efficiently

## Learning Outcomes

After completing this project, you will understand:
- CSR graph representation
- Cache-aware data structure design
- Go 1.23 iterator patterns
- Iterator composition techniques
- Memory layout optimization
- Graph algorithm implementation
- Performance profiling for graphs

## Time Estimate
- Core implementation: 8-10 hours
- Testing and optimization: 3-4 hours
- Iterator composition: 2-3 hours
- Stretch goals: 5-7 hours (optional)

## Next Steps
After completing this project, move on to **Project 2.2: Columnar Property Store** which builds on memory-efficient storage techniques.
