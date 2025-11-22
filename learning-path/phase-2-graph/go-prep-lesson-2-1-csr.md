# Phase 2 Lesson 2.1: Go Prep - Compressed Sparse Row (CSR)

**Prerequisites:** Phase 1 complete (Storage Layer)
**Time:** 4-5 hours Go prep + 25-30 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 2.1

## Overview

Compressed Sparse Row (CSR) is the foundation of graph storage. Before implementing it, master these Go concepts:
- **Go 1.23:** `iter.Seq` and range-over-func iterators ⭐⭐⭐ **CRITICAL!**
- Custom iterator implementation with yield functions
- Iterator composition and pipelines
- Memory layout optimization for cache locality
- **Go 1.25:** Experimental GreenTeaGC testing

**This lesson focuses heavily on Go 1.23 iterators - they're transformative for graph traversal!**

## Go Concepts for This Lesson

### 1. Go 1.23 Range-Over-Func: The Basics

**Go 1.23's killer feature for graph databases!**

```go
package main

import (
    "fmt"
    "iter"
)

// Traditional approach: Return a slice (allocates!)
func GetNumbersSlice() []int {
    return []int{1, 2, 3, 4, 5}  // Heap allocation
}

// Go 1.23: Return an iterator (zero allocations!)
func GetNumbersIter() iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := 1; i <= 5; i++ {
            if !yield(i) {
                return  // Early exit
            }
        }
    }
}

func main() {
    // Old way
    for _, n := range GetNumbersSlice() {
        fmt.Println(n)
    }

    // New way - same syntax, better performance!
    for n := range GetNumbersIter() {
        fmt.Println(n)
    }

    // Early exit (stops iteration)
    for n := range GetNumbersIter() {
        fmt.Println(n)
        if n == 3 {
            break  // yield returns false
        }
    }
}
```

**Key insight:** The `yield` function returns `false` when the loop breaks, allowing cleanup!

### 2. Iterator Pattern for Graphs

**This is the pattern you'll use constantly!**

```go
package main

import (
    "fmt"
    "iter"
)

type NodeID uint32

type Graph struct {
    offsets []uint64  // offsets[i] = start index for node i's edges
    targets []NodeID  // targets[j] = destination of edge j
}

// Neighbors returns an iterator over a node's neighbors
func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        start := g.offsets[node]
        end := g.offsets[node+1]

        for i := start; i < end; i++ {
            if !yield(g.targets[i]) {
                return  // Early exit
            }
        }
    }
}

func main() {
    // Build a simple graph:
    //   0 -> 1, 2
    //   1 -> 2, 3
    //   2 -> 3
    g := Graph{
        offsets: []uint64{0, 2, 4, 5},
        targets: []NodeID{1, 2, 2, 3, 3},
    }

    // Iterate neighbors of node 0
    fmt.Print("Node 0 neighbors: ")
    for neighbor := range g.Neighbors(0) {
        fmt.Printf("%d ", neighbor)
    }
    fmt.Println()

    // Iterate neighbors of node 1
    fmt.Print("Node 1 neighbors: ")
    for neighbor := range g.Neighbors(1) {
        fmt.Printf("%d ", neighbor)
    }
    fmt.Println()
}
```

**Output:**
```
Node 0 neighbors: 1 2
Node 1 neighbors: 2 3
```

### 3. Iterator Composition (The Power!)

**Compose iterators to build complex queries!**

```go
package main

import (
    "fmt"
    "iter"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
}

func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        start := g.offsets[node]
        end := g.offsets[node+1]
        for i := start; i < end; i++ {
            if !yield(g.targets[i]) {
                return
            }
        }
    }
}

// Two-hop neighbors (friends of friends)
func (g *Graph) TwoHop(start NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        // For each neighbor
        for n1 := range g.Neighbors(start) {
            // For each neighbor's neighbor
            for n2 := range g.Neighbors(n1) {
                if !yield(n2) {
                    return
                }
            }
        }
    }
}

// Filter iterator (only nodes matching predicate)
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for item := range seq {
            if predicate(item) {
                if !yield(item) {
                    return
                }
            }
        }
    }
}

// Map iterator (transform values)
func Map[T, U any](seq iter.Seq[T], transform func(T) U) iter.Seq[U] {
    return func(yield func(U) bool) {
        for item := range seq {
            if !yield(transform(item)) {
                return
            }
        }
    }
}

func main() {
    g := Graph{
        offsets: []uint64{0, 2, 4, 5},
        targets: []NodeID{1, 2, 2, 3, 3},
    }

    // Two-hop neighbors of node 0
    fmt.Print("Two-hop from 0: ")
    for node := range g.TwoHop(0) {
        fmt.Printf("%d ", node)
    }
    fmt.Println()

    // Filter: only neighbors > 1
    fmt.Print("Node 1 neighbors > 1: ")
    filtered := Filter(g.Neighbors(1), func(n NodeID) bool {
        return n > 1
    })
    for node := range filtered {
        fmt.Printf("%d ", node)
    }
    fmt.Println()

    // Map: multiply node IDs by 10
    fmt.Print("Node 0 neighbors * 10: ")
    mapped := Map(g.Neighbors(0), func(n NodeID) NodeID {
        return n * 10
    })
    for node := range mapped {
        fmt.Printf("%d ", node)
    }
    fmt.Println()
}
```

**Output:**
```
Two-hop from 0: 2 3 3
Node 1 neighbors > 1: 2 3
Node 0 neighbors * 10: 10 20
```

**This is HUGE for graph queries!** No intermediate allocations, early exit support, composable operations.

### 4. iter.Seq2 for Key-Value Pairs

**Use `iter.Seq2` for edges with properties!**

```go
package main

import (
    "fmt"
    "iter"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
    weights []float64  // Edge weights
}

// NeighborsWithWeight returns (target, weight) pairs
func (g *Graph) NeighborsWithWeight(node NodeID) iter.Seq2[NodeID, float64] {
    return func(yield func(NodeID, float64) bool) {
        start := g.offsets[node]
        end := g.offsets[node+1]

        for i := start; i < end; i++ {
            if !yield(g.targets[i], g.weights[i]) {
                return
            }
        }
    }
}

func main() {
    g := Graph{
        offsets: []uint64{0, 2, 4},
        targets: []NodeID{1, 2, 0, 2},
        weights: []float64{0.5, 1.0, 0.3, 0.8},
    }

    // Iterate edges with weights
    fmt.Println("Node 0 edges:")
    for target, weight := range g.NeighborsWithWeight(0) {
        fmt.Printf("  -> %d (weight: %.1f)\n", target, weight)
    }

    fmt.Println("Node 1 edges:")
    for target, weight := range g.NeighborsWithWeight(1) {
        fmt.Printf("  -> %d (weight: %.1f)\n", target, weight)
    }
}
```

**Output:**
```
Node 0 edges:
  -> 1 (weight: 0.5)
  -> 2 (weight: 1.0)
Node 1 edges:
  -> 0 (weight: 0.3)
  -> 2 (weight: 0.8)
```

### 5. Memory Layout for Cache Locality

**CSR is cache-friendly because edges are contiguous!**

```go
package main

import (
    "fmt"
    "unsafe"
)

type NodeID uint32

// Bad: Pointer-based graph (poor cache locality)
type AdjacencyList struct {
    nodes []Node
}

type Node struct {
    id    NodeID
    edges []*Edge  // Pointers! Cache misses!
}

type Edge struct {
    target NodeID
}

// Good: CSR format (excellent cache locality)
type CSRGraph struct {
    offsets []uint64  // Small array
    targets []NodeID  // Contiguous! Cache-friendly!
}

func main() {
    // Adjacency list: nodes scattered in memory
    adjList := AdjacencyList{
        nodes: []Node{
            {id: 0, edges: []*Edge{{target: 1}, {target: 2}}},
            {id: 1, edges: []*Edge{{target: 2}}},
        },
    }

    // CSR: edges packed tightly
    csr := CSRGraph{
        offsets: []uint64{0, 2, 3},
        targets: []NodeID{1, 2, 2},
    }

    fmt.Println("Adjacency list size:")
    fmt.Printf("  Node size: %d bytes\n", unsafe.Sizeof(Node{}))
    fmt.Printf("  Edge pointer: %d bytes\n", unsafe.Sizeof(&Edge{}))

    fmt.Println("\nCSR size:")
    fmt.Printf("  Offset size: %d bytes\n", unsafe.Sizeof(uint64(0)))
    fmt.Printf("  Target size: %d bytes\n", unsafe.Sizeof(NodeID(0)))

    fmt.Println("\nCSR stores edges contiguously - great for CPU cache!")
}
```

**Key insight:** CSR stores all edges in one array → CPU prefetcher can load multiple edges per cache line!

### 6. Performance: Iterator vs Slice Return

**Benchmark the difference!**

```go
package main

import (
    "iter"
    "testing"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
}

// Old way: Return slice (allocation!)
func (g *Graph) NeighborsSlice(node NodeID) []NodeID {
    start := g.offsets[node]
    end := g.offsets[node+1]

    neighbors := make([]NodeID, end-start)  // Heap allocation
    copy(neighbors, g.targets[start:end])
    return neighbors
}

// New way: Return iterator (zero allocation!)
func (g *Graph) NeighborsIter(node NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        start := g.offsets[node]
        end := g.offsets[node+1]
        for i := start; i < end; i++ {
            if !yield(g.targets[i]) {
                return
            }
        }
    }
}

func makeGraph() Graph {
    // Graph with 1000 nodes, avg 10 edges each
    offsets := make([]uint64, 1001)
    targets := make([]NodeID, 10000)

    for i := uint64(0); i < 1000; i++ {
        offsets[i] = i * 10
        for j := uint64(0); j < 10; j++ {
            targets[i*10+j] = NodeID((i + j) % 1000)
        }
    }
    offsets[1000] = 10000

    return Graph{offsets: offsets, targets: targets}
}

func BenchmarkNeighborsSlice(b *testing.B) {
    g := makeGraph()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, neighbor := range g.NeighborsSlice(NodeID(i % 1000)) {
            _ = neighbor
        }
    }
}

func BenchmarkNeighborsIter(b *testing.B) {
    g := makeGraph()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for neighbor := range g.NeighborsIter(NodeID(i % 1000)) {
            _ = neighbor
        }
    }
}
```

**Expected results:**
- Slice: ~150ns + allocation
- Iterator: ~50ns, zero allocation (3x faster!)

### 7. Go 1.25: Experimental GreenTeaGC

**Test your graph algorithms with the new GC!**

```go
// Set environment variable to try GreenTeaGC
// GODEBUG=gcpercent=100,gctrace=1 go run main.go
//
// Or in Go 1.25:
// GOEXPERIMENT=greentea go test

package main

import (
    "fmt"
    "runtime"
)

func main() {
    // Force GC and print stats
    var m runtime.MemStats

    runtime.ReadMemStats(&m)
    fmt.Printf("Alloc: %d MB\n", m.Alloc/1024/1024)

    // Build large graph
    offsets := make([]uint64, 1000000)
    targets := make([]uint32, 10000000)

    runtime.ReadMemStats(&m)
    fmt.Printf("After graph: %d MB\n", m.Alloc/1024/1024)

    // Force GC
    runtime.GC()

    runtime.ReadMemStats(&m)
    fmt.Printf("After GC: %d MB\n", m.Alloc/1024/1024)

    _ = offsets
    _ = targets
}
```

**Note:** GreenTeaGC is experimental in Go 1.25. Test but don't rely on it yet!

## Pre-Implementation Exercises

### Exercise 1: Basic CSR Iterator

```go
package main

import (
    "fmt"
    "iter"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
}

func NewGraph(edges [][]NodeID) Graph {
    // TODO: Convert edge list to CSR format
    // edges[i] = list of neighbors for node i
    return Graph{}
}

func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    // TODO: Implement iterator
    return nil
}

func (g *Graph) NumNodes() int {
    // TODO: Return number of nodes
    return 0
}

func (g *Graph) NumEdges() int {
    // TODO: Return number of edges
    return 0
}

func main() {
    // Build graph: 0->1,2  1->2  2->0
    edges := [][]NodeID{
        {1, 2},
        {2},
        {0},
    }

    g := NewGraph(edges)

    fmt.Printf("Nodes: %d, Edges: %d\n", g.NumNodes(), g.NumEdges())

    for i := 0; i < g.NumNodes(); i++ {
        fmt.Printf("Node %d neighbors: ", i)
        for neighbor := range g.Neighbors(NodeID(i)) {
            fmt.Printf("%d ", neighbor)
        }
        fmt.Println()
    }
}
```

### Exercise 2: Iterator Composition

```go
package main

import (
    "iter"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
}

func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    // TODO: (from Exercise 1)
    return nil
}

// TODO: Implement TwoHop - returns friends-of-friends
func (g *Graph) TwoHop(start NodeID) iter.Seq[NodeID] {
    return nil
}

// TODO: Implement Filter
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
    return nil
}

// TODO: Implement Collect (iterator -> slice)
func Collect[T any](seq iter.Seq[T]) []T {
    return nil
}

// TODO: Implement Unique (deduplicate)
func Unique[T comparable](seq iter.Seq[T]) iter.Seq[T] {
    return nil
}

func main() {
    // TODO: Test your implementations
}
```

### Exercise 3: Weighted Graph with iter.Seq2

```go
package main

import (
    "iter"
)

type NodeID uint32

type WeightedGraph struct {
    offsets []uint64
    targets []NodeID
    weights []float64
}

func (g *WeightedGraph) Neighbors(node NodeID) iter.Seq2[NodeID, float64] {
    // TODO: Return (target, weight) pairs
    return nil
}

// TODO: Find neighbors with weight > threshold
func (g *WeightedGraph) FilterByWeight(node NodeID, minWeight float64) iter.Seq[NodeID] {
    return nil
}

func main() {
    // TODO: Build weighted graph and test
}
```

### Exercise 4: Benchmark Iterator Performance

```go
package main

import (
    "iter"
    "testing"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
}

func makeTestGraph(numNodes, avgDegree int) Graph {
    // TODO: Generate random graph
    return Graph{}
}

func (g *Graph) NeighborsSlice(node NodeID) []NodeID {
    // TODO: Return slice (allocates)
    return nil
}

func (g *Graph) NeighborsIter(node NodeID) iter.Seq[NodeID] {
    // TODO: Return iterator (zero alloc)
    return nil
}

func BenchmarkSliceIteration(b *testing.B) {
    g := makeTestGraph(10000, 10)
    // TODO: Benchmark slice approach
}

func BenchmarkIteratorIteration(b *testing.B) {
    g := makeTestGraph(10000, 10)
    // TODO: Benchmark iterator approach
}
```

Run with:
```bash
go test -bench=. -benchmem
```

Compare allocations!

### Exercise 5: BFS with Iterators

```go
package main

import (
    "iter"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
}

func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    // TODO: (from Exercise 1)
    return nil
}

// BFS returns nodes in breadth-first order
func (g *Graph) BFS(start NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        // TODO: Implement BFS using iterator
        // Hint: Use a queue and visited map
    }
}

func main() {
    // TODO: Build graph and test BFS
}
```

## Performance Benchmarks

### Benchmark 1: CSR vs Adjacency List

```go
func BenchmarkCSRTraversal(b *testing.B) {
    // CSR format
    g := makeCSRGraph(10000, 10)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        count := 0
        for neighbor := range g.Neighbors(0) {
            count++
        }
    }
}

func BenchmarkAdjacencyListTraversal(b *testing.B) {
    // Pointer-based
    g := makeAdjacencyList(10000, 10)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        count := 0
        for _, neighbor := range g.nodes[0].edges {
            count++
        }
    }
}
```

**Expected: CSR 2-3x faster (cache locality!)**

### Benchmark 2: Early Exit Performance

```go
func BenchmarkEarlyExitSlice(b *testing.B) {
    g := makeGraph(10000, 100)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        neighbors := g.NeighborsSlice(0)
        for _, n := range neighbors {
            if n > 10 {
                break  // But we already allocated the whole slice!
            }
        }
    }
}

func BenchmarkEarlyExitIter(b *testing.B) {
    g := makeGraph(10000, 100)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for n := range g.NeighborsIter(0) {
            if n > 10 {
                break  // Only iterated what we needed!
            }
        }
    }
}
```

**Expected: Iterator much faster when exiting early!**

## Common Gotchas to Avoid

### Gotcha 1: Off-By-One in Offsets Array

```go
// WRONG: offsets array too small
offsets := make([]uint64, numNodes)  // Missing final offset!

// RIGHT: Need numNodes + 1
offsets := make([]uint64, numNodes+1)
offsets[numNodes] = uint64(len(targets))
```

### Gotcha 2: Forgetting Early Exit Check

```go
// WRONG: Doesn't respect break!
func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        start := g.offsets[node]
        end := g.offsets[node+1]
        for i := start; i < end; i++ {
            yield(g.targets[i])  // Ignores return value!
        }
    }
}

// RIGHT: Check yield return value
func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        start := g.offsets[node]
        end := g.offsets[node+1]
        for i := start; i < end; i++ {
            if !yield(g.targets[i]) {
                return  // Respect early exit!
            }
        }
    }
}
```

### Gotcha 3: Mutating During Iteration

```go
// WRONG: Modifying graph during iteration!
for neighbor := range g.Neighbors(0) {
    g.AddEdge(neighbor, 5)  // RACE or panic!
}

// RIGHT: Collect first, then modify
toAdd := Collect(g.Neighbors(0))
for _, neighbor := range toAdd {
    g.AddEdge(neighbor, 5)
}
```

### Gotcha 4: Not Pre-Sizing Slices

```go
// WRONG: Growing slice repeatedly
var nodes []NodeID
for n := range g.Neighbors(0) {
    nodes = append(nodes, n)  // Multiple reallocations!
}

// RIGHT: Pre-size if you know capacity
degree := g.offsets[1] - g.offsets[0]
nodes := make([]NodeID, 0, degree)
for n := range g.Neighbors(0) {
    nodes = append(nodes, n)  // No reallocation!
}
```

## Checklist Before Starting Lesson 2.1

- [ ] I understand Go 1.23 `iter.Seq` and range-over-func
- [ ] I can implement custom iterators with yield functions
- [ ] I understand CSR memory layout
- [ ] I can compose iterators (filter, map, two-hop)
- [ ] I know how to use `iter.Seq2` for key-value pairs
- [ ] I understand cache locality benefits of CSR
- [ ] I've benchmarked iterator vs slice return
- [ ] I always check yield return value for early exit
- [ ] I understand when to use iterators vs slices
- [ ] I know how to profile cache misses with pprof

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 2.1

You'll implement:
- CSR graph storage format
- Iterator-based graph traversal (1-hop, 2-hop, n-hop)
- BFS and DFS using iterators
- Filter and map operations on graphs
- Benchmarks comparing CSR vs adjacency list
- Cache locality measurements

**Time estimate:** 25-30 hours for full implementation

**This is the most important Go 1.23 lesson - iterators are the future!**

Good luck! 🚀
