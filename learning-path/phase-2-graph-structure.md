# Phase 2: Graph Structure

**Prerequisites:** Complete Phase 1 (Storage Layer)

**Time Estimate:** 3 weeks
**Go Version:** 1.23+

In this phase, you'll learn how to efficiently represent and store graph data. You'll implement compressed sparse row (CSR) format, columnar storage for properties, and parallel query execution.

---

## Table of Contents

- [Lesson 2.1: Compressed Sparse Row (CSR)](#lesson-21-compressed-sparse-row-week-4)
- [Lesson 2.2: Columnar Node Properties](#lesson-22-columnar-node-properties-week-5)
- [Lesson 2.3: NodeGroups and Morsel-Driven Parallelism](#lesson-23-nodegroups-and-morsel-driven-parallelism-week-6)

---

## Lesson 2.1: Compressed Sparse Row (Week 4)

**Your Mission:** Store 1M nodes with 10M edges using <80MB RAM.

**The Naive Approach:**
```go
type Graph struct {
    nodes map[NodeID]*Node
}

type Node struct {
    ID    NodeID
    Edges []NodeID  // Uh oh... 10 million []NodeID slices
}
```
Calculate: How much memory does this use? Why is it terrible?

### Investigation Questions

1. **CSR Structure Deep Dive:**
   ```
   Nodes: A B C D
   Edges: A→B, A→C, B→C, C→D, C→B

   Offsets: [0, 2, 3, 5, 5]
   Targets: [B, C, C, D, B]
   ```
   - How do you find A's neighbors?
   - How do you add a new edge to C?
   - What if you need to delete an edge?

2. **Cache Locality Experiment:**
   ```go
   // Approach 1: Adjacency list (random access)
   for _, node := range nodes {
       for _, edge := range node.edges {
           visit(edge)
       }
   }

   // Approach 2: CSR (sequential access)
   for i := 0; i < len(offsets)-1; i++ {
       for j := offsets[i]; j < offsets[i+1]; j++ {
           visit(targets[j])
       }
   }
   ```
   Run `perf stat -e cache-misses` on both. What's the difference?

3. **Bidirectional Traversal:**
   - CSR gives you A→B efficiently
   - How do you answer "who points TO A?" (B→A)?
   - Do you build a second CSR in reverse?
   - How much extra memory does that cost?

### 🆕 Go 1.23: Iterator-Based Traversal

The modern approach uses iterators:

```go
import "iter"

type Graph struct {
    offsets []uint64
    targets []NodeID
}

// 🆕 Return an iterator instead of a slice
func (g *Graph) Neighbors(n NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        start := g.offsets[n]
        end := g.offsets[n+1]
        for i := start; i < end; i++ {
            if !yield(g.targets[i]) {
                return  // Early exit if consumer stops
            }
        }
    }
}

// Clean, composable traversals:
func (g *Graph) TwoHopNeighbors(start NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        for n1 := range g.Neighbors(start) {
            for n2 := range g.Neighbors(n1) {
                if !yield(n2) {
                    return
                }
            }
        }
    }
}

// Usage:
for neighbor := range graph.Neighbors(startNode) {
    fmt.Println(neighbor)
}
```

**Investigation:** Compare iterator vs slice return:
```go
func BenchmarkTraversalMethods(b *testing.B) {
    g := buildGraph(1_000_000, 10_000_000)

    b.Run("Iterator", func(b *testing.B) {
        for b.Loop() {
            count := 0
            for _ = range g.Neighbors(42) {
                count++
            }
        }
    })

    b.Run("Slice", func(b *testing.B) {
        for b.Loop() {
            neighbors := g.GetNeighborsSlice(42)
            count := len(neighbors)
            _ = count
        }
    })
}
```

Which allocates less memory? Use `go test -benchmem`.

### Your Implementation

- [ ] Build CSR from edge list in O(E) time
- [ ] Query neighbors in O(degree(n)) time
- [ ] Measure cache misses (use `pprof` or `perf`)
- [ ] Support both forward and backward traversal
- [ ] 🆕 Implement iterator-based API

### Benchmark (Go 1.24+)

```go
// Graph: 1M nodes, 10M edges (average degree 10)
func BenchmarkTraversal(b *testing.B) {
    g := buildGraph(1_000_000, 10_000_000)

    for b.Loop() {  // 🆕 Go 1.24 syntax
        // 2-hop traversal
        count := 0
        for n1 := range g.Neighbors(startNode) {
            for n2 := range g.Neighbors(n1) {
                count++
            }
        }
    }
}
```
Target: <100ns per edge access

### 🆕 Go 1.25: Experimental GC Performance

Test with the new garbage collector:
```bash
GOEXPERIMENT=greenteagc go test -bench=BenchmarkTraversal
```

The new GC improves cache locality during marking. For pointer-heavy workloads (graphs!), expect 10-35% improvement.

### What You Should Discover

- Sequential memory access is 10-100x faster
- CSR trades insert performance for query performance
- Cache line prefetching is automatic (and magical)
- 🆕 Iterators provide zero-cost abstractions
- 🆕 Experimental GC significantly helps graph traversals

---

## Lesson 2.2: Columnar Node Properties (Week 5)

**Your Mission:** Store properties for 1M nodes where:
- 50% have just `id` and `name`
- 30% add `age` and `email`
- 20% add `address`, `phone`, `salary`

**The Wrong Way:**
```go
type Node struct {
    ID         uint64
    Properties map[string]interface{}  // 🚨 SLOW
}
```
Why is this slow? Measure it.

### Investigation Questions

1. **Column vs Row Storage:**
   ```
   Row-oriented:
   [id=1, name="Alice", age=30] [id=2, name="Bob", age=25]

   Column-oriented:
   ids:   [1, 2]
   names: ["Alice", "Bob"]
   ages:  [30, 25]
   ```
   - Read all names: which layout is faster?
   - Read one full node: which layout is faster?
   - What if names are NULL for 70% of nodes?

2. **Compression Techniques:**
   Implement these and measure compression ratio + decompression speed:

   **Bit-packing:** Ages range 0-100 → store in 7 bits not 64
   ```go
   func packAges(ages []uint8) []byte {
       // Your code here
   }
   ```

   **Dictionary Encoding:** Country names
   ```go
   dict := []string{"USA", "Canada", "Mexico"}
   encoded := []uint8{0, 0, 1, 0, 2} // indices
   ```

   **Run-Length Encoding:** Gender column
   ```go
   values: [M, M, M, M, F, F, M, M, M]
   rle:    [(M, 4), (F, 2), (M, 3)]
   ```

3. **NULL Handling:**
   ```go
   // 1M nodes, only 10% have 'phone' property
   type Column struct {
       Values  []string
       Nulls   []bool   // 1M bools = 125KB
   }
   ```
   - Can you do better than a bool array?
   - Research: Validity bitmaps
   - How does Arrow format handle this?

### 🆕 Go 1.23: String Interning with `unique` Package

For low-cardinality columns, use interning:

```go
import "unique"

type StringColumn struct {
    values []unique.Handle[string]  // Interned strings
}

func (c *StringColumn) Set(idx int, value string) {
    c.values[idx] = unique.Make(value)
}

// Comparison is O(1) pointer comparison!
func (c *StringColumn) Equals(i, j int) bool {
    return c.values[i] == c.values[j]
}

// Count distinct values efficiently
func (c *StringColumn) Distinct() int {
    seen := make(map[unique.Handle[string]]struct{})
    for _, h := range c.values {
        seen[h] = struct{}{}
    }
    return len(seen)
}
```

**Memory Savings Example:**
```
1M nodes with countries:
- Without interning: 1M strings × ~15 bytes = 15MB
- With interning: 200 unique strings = ~3KB + 1M handles × 8 bytes = 8MB
- Savings: 47% reduction
```

**Investigation:** When does interning hurt performance?
- Measure Set() latency with/without interning
- Measure Equals() throughput
- Find the cardinality threshold where interning helps

### Challenge

- [ ] Store 1M heterogeneous nodes in <50MB
- [ ] Query: "Find all ages > 50" without decompressing all columns
- [ ] Add new property type without rewriting all data
- [ ] Support NULL values efficiently
- [ ] 🆕 Use `unique.Make()` for low-cardinality columns

### Benchmark Query (Go 1.24+)

```go
// SELECT name WHERE age > 30 AND country = "USA"
func BenchmarkColumnScan(b *testing.B) {
    store := buildColumnStore(1_000_000)

    for b.Loop() {
        results := store.Query(func(n Node) bool {
            return n.Age > 30 && n.Country == "USA"
        })
    }
}
```
Target: Scan 1M rows in <10ms

### What You Should Discover

- Column stores are 10x better for analytical queries
- Compression works better on columns (same type together)
- Sparse columns with mostly NULLs compress to almost nothing
- 🆕 `unique.Handle` dramatically reduces memory for repeated values
- 🆕 Go 1.24: Swiss Tables make dictionary encoding faster

---

## Lesson 2.3: NodeGroups and Morsel-Driven Parallelism (Week 6)

**Your Mission:** Scan 1M nodes using all your CPU cores.

**Single-Threaded Baseline:**
```go
for nodeID := 0; nodeID < 1_000_000; nodeID++ {
    if predicate(nodeID) {
        results = append(results, nodeID)
    }
}
// Takes 100ms on my machine
```

### Investigation Questions

1. **Work Partitioning:**
   Kuzu uses "NodeGroups" of 131,072 nodes (why this number?).
   ```go
   const NodeGroupSize = 131072
   numGroups := (totalNodes + NodeGroupSize - 1) / NodeGroupSize
   ```
   - Start 8 goroutines, each processes full groups
   - What if you have 10 groups but 8 cores?
   - What if you have 3 groups but 8 cores?
   - Measure: overhead of goroutine creation

2. **Work Stealing:**
   ```go
   // Naive: Static assignment
   goroutine1 := groups[0:3]
   goroutine2 := groups[3:6]
   // Problem: What if group 0 is huge?

   // Better: Work queue
   queue := make(chan *NodeGroup, numGroups)
   for i := 0; i < numWorkers; i++ {
       go worker(queue)
   }
   ```
   - Implement both. Which is faster?
   - What's the sweet spot for queue buffer size?
   - When does channel contention matter?

3. **Result Aggregation:**
   ```go
   // Each goroutine finds matches, now what?

   // Approach A: Shared slice + mutex
   var results []NodeID
   var mu sync.Mutex

   // Approach B: Local slices + merge
   resultChan := make(chan []NodeID)

   // Approach C: Lock-free atomic counter
   resultArray := make([]NodeID, totalNodes)
   var count atomic.Uint64
   ```
   Implement all three. Which scales best to 32 cores?

### 🆕 Go 1.25: Container-Aware GOMAXPROCS

```go
func NewExecutor() *Executor {
    // 🆕 Go 1.25: Automatically detects container CPU limits
    // No manual tuning needed in Docker/Kubernetes!
    workers := runtime.GOMAXPROCS(0)

    return &Executor{
        workers: workers,
        queue:   make(chan *NodeGroup, workers*2),
    }
}
```

### 🆕 Go 1.25: Testing Parallel Execution

Use `testing/synctest` for deterministic concurrent tests:

```go
import "testing/synctest"

func TestWorkStealingFairness(t *testing.T) {
    synctest.Run(func() {
        queue := make(chan int, 100)
        results := make([]atomic.Int32, 8)

        // Spawn 8 workers
        for i := 0; i < 8; i++ {
            workerID := i
            go func() {
                for item := range queue {
                    time.Sleep(1 * time.Millisecond)  // Fake time!
                    results[workerID].Add(1)
                }
            }()
        }

        // Send 100 work items
        for i := 0; i < 100; i++ {
            queue <- i
        }
        close(queue)

        synctest.Wait()  // Wait for all goroutines to finish

        // Verify fair distribution
        for i := 0; i < 8; i++ {
            count := results[i].Load()
            if count < 10 || count > 15 {
                t.Errorf("Worker %d: got %d items, want 10-15", i, count)
            }
        }
    })
}
```

This test runs in fake time and is deterministic!

### Your Implementation

- [ ] Achieve near-linear speedup with cores (8 cores → 7.5x faster)
- [ ] Handle unbalanced workloads (some groups scan faster)
- [ ] Measure goroutine overhead (<5% of total time)
- [ ] Zero allocations in hot path (use `go test -benchmem`)
- [ ] 🆕 Use `synctest` for concurrent correctness testing

### Scaling Test (Go 1.24+)

```bash
# Run with different GOMAXPROCS
for cores in 1 2 4 8 16; do
    GOMAXPROCS=$cores go test -bench=ParallelScan -cpu=$cores
done

# Expected output:
GOMAXPROCS=1   100ms
GOMAXPROCS=2    52ms  (1.9x speedup)
GOMAXPROCS=4    27ms  (3.7x speedup)
GOMAXPROCS=8    14ms  (7.1x speedup)
```

### Benchmark (Go 1.24+)

```go
func BenchmarkParallelScan(b *testing.B) {
    store := buildColumnStore(1_000_000)

    for cores := 1; cores <= 16; cores *= 2 {
        b.Run(fmt.Sprintf("cores=%d", cores), func(b *testing.B) {
            runtime.GOMAXPROCS(cores)

            for b.Loop() {
                _ = store.ParallelScan(func(n Node) bool {
                    return n.Age > 30
                })
            }
        })
    }
}
```

### What You Should Discover

- Goroutines are cheap but not free
- Work stealing beats static partitioning
- False sharing kills parallel performance
- CPU cache lines are 64 bytes (measure with `-cpu=1` vs `-cpu=8`)
- 🆕 `synctest` makes concurrency tests reproducible

---

## Phase 2 Complete!

**What you've built:**
- ✅ Compressed sparse row (CSR) graph representation
- ✅ Columnar storage with compression
- ✅ Parallel query execution engine

**Key Insights:**
- Memory layout affects performance by orders of magnitude
- Compression is almost always worth it
- Parallelism requires careful design to scale

**Next:** Phase 3 - Query Engine, where you'll parse Cypher queries, optimize query plans, and implement join algorithms.

---

## Resources

### Papers to Read
- "Worst-Case Optimal Join Algorithms" - Understanding graph joins
- "MonetDB/X100: Hyper-Pipelining Query Execution" - Vectorized execution
- "Morsel-Driven Parallelism" - HyPer's approach

### Reference Implementations
- Kuzu's storage format
- DuckDB's columnar storage
- GraphBLAS sparse matrix format

### Go-Specific
- `iter` package (Go 1.23+)
- `unique` package (Go 1.23+)
- `testing/synctest` (Go 1.25+)
