# Graph Database Learning Path in Go (Go 1.23+ Edition)

A comprehensive, hands-on curriculum for learning modern Go by building an embeddable graph database inspired by Kuzu DB.

**Target Go Version:** 1.23+ (with optional 1.24 and 1.25 features)  
**Last Updated:** November 2025

**Philosophy:** Each lesson follows the pattern: Goal → Investigation → Implementation → Verification → Reflection. You discover through experimentation and measurement, not by reading solutions.

---

## Table of Contents

- [Getting Started](#getting-started)
- [Phase 1: Storage Layer](#phase-1-storage-layer)
- [Phase 2: Graph Structure](#phase-2-graph-structure)
- [Phase 3: Query Engine](#phase-3-query-engine)
- [Phase 4: Transactions](#phase-4-transactions)
- [Bonus Challenges](#bonus-challenges)
- [Progress Tracking](#your-progress-tracking)
- [Resources](#recommended-resources)

---

## Getting Started

### Prerequisites

**Install Go 1.25 (Latest):**
```bash
# Download from https://go.dev/dl/
# Or use version manager
go install golang.org/dl/go1.25@latest
go1.25 download
```

**Verify installation:**
```bash
go version  # Should show go1.25 or later
```

**Optional: Enable Experimental Features**
```bash
# Add to your shell profile
export GOEXPERIMENT=greenteagc  # Experimental GC for better graph performance
```

### Project Setup

```bash
mkdir kuzu-go
cd kuzu-go
go mod init github.com/yourusername/kuzu-go

# Create initial structure
mkdir -p {storage,graph,query,index,transaction}
touch storage/page.go
```

### What's New in Go 1.23-1.25?

This curriculum uses modern Go features:
- **Go 1.23:** Range-over-func iterators, `unique` package for memory efficiency, improved timers
- **Go 1.24:** Swiss Tables maps (30-35% faster hash joins), weak pointers, `testing.B.Loop`
- **Go 1.25:** Stable `testing/synctest` for concurrent testing, experimental new GC

These features are integrated throughout the lessons. We'll call out when features require specific versions.

---

## PHASE 1: STORAGE LAYER

### Lesson 1.1: The Page Abstraction (Week 1)

**Your Mission:** Store 1 million integers on disk and read them back.

**Constraints:**
- File size must be exactly 4MB (1M integers × 4 bytes)
- Read them back in under 100ms
- Prove data survives program crash

#### Investigation Questions

1. **Why 4KB pages?** 
   - Run: `getconf PAGESIZE` on your system
   - Open any file and check: `stat -f %k yourfile`
   - What happens if you use 3KB pages? 8KB?

2. **Buffered vs Direct I/O:**
   ```go
   // Try both approaches
   file, _ := os.Create("data.db")
   file.Write(data) // buffered
   
   // vs
   syscall.Open(name, syscall.O_DIRECT, 0644) // direct
   ```
   - Which is faster for sequential writes?
   - Which is faster for random reads?
   - When does `file.Sync()` actually matter?

3. **Memory Mapping Mystery:**
   ```go
   import "syscall"
   
   data, _ := syscall.Mmap(fd, 0, size, 
       syscall.PROT_READ|syscall.PROT_WRITE, 
       syscall.MAP_SHARED)
   ```
   - Write to `data[100]`. When does it hit disk?
   - Call `syscall.Msync()`. Measure the latency.
   - Read 1000 random locations. Compare to file reads.

#### Your Implementation Must

- [ ] Store exactly N pages with zero wasted space
- [ ] Read any page in O(1) time using page ID
- [ ] Survive `kill -9` without data loss (after sync)
- [ ] Print stats: total pages, bytes per page, storage efficiency

#### Verification

```bash
# Your program should work like this:
$ go run storage.go write 1000000
Wrote 1,000,000 integers to 244 pages in 23ms

$ go run storage.go read 42
Page 42 contains: [172032, 172033, ..., 172163]

$ kill -9 <pid>  # crash it mid-write
$ go run storage.go verify
✓ All 244 pages intact
```

#### What You Should Discover

- OS page cache is doing a lot of work for you
- `mmap` trades simplicity for control
- Alignment matters more than you think

---

### Lesson 1.2: The Buffer Pool (Week 2)

**Your Mission:** You can only keep 10 pages in memory, but need to access 1000 pages randomly.

**The Problem:**
```go
// This will OOM with 1M pages
cache := make(map[uint32]*Page)
for i := 0; i < 1000000; i++ {
    cache[i] = readPageFromDisk(i)
}
```

#### Investigation Questions

1. **Eviction Policies Deep Dive:**
   - Implement LRU: How do you get O(1) get AND O(1) evict?
   - Implement Clock (Second Chance): Why does SQLite use this?
   - Implement LRU-K: When does K=2 beat K=1?
   
   **Measure:** Hit rate on these access patterns:
   - Sequential (1,2,3,4...)
   - Zipfian (page 1 accessed 50%, page 2 25%, page 3 12.5%...)
   - Random uniform

2. **Dirty Page Tracking:**
   ```go
   type Page struct {
       ID    uint32
       Data  [4096]byte
       Dirty bool
       Pins  atomic.Int32  // Thread-safe counter
   }
   ```
   - What happens if you evict a dirty page?
   - What if a goroutine is reading while you evict?
   - How do you flush dirty pages in the background?

3. **Concurrency Control:**
   ```go
   // 100 goroutines all want the same page
   for i := 0; i < 100; i++ {
       go func() {
           page := pool.Get(pageID)
           // use page
       }()
   }
   ```
   - Do you lock the entire pool or per-page?
   - What's the cost of `sync.RWMutex` vs channels?
   - Use `go test -race` and make it fail, then fix it

#### 🆕 Go 1.24: Weak Reference Cache (Optional Advanced Challenge)

```go
import "weak"

type WeakBufferPool struct {
    pages map[uint32]weak.Pointer[*Page]
    mu    sync.RWMutex
}

func (bp *WeakBufferPool) Get(id uint32) *Page {
    bp.mu.RLock()
    if wp, ok := bp.pages[id]; ok {
        if page := wp.Value(); page != nil {
            bp.mu.RUnlock()
            return page  // Still in memory!
        }
    }
    bp.mu.RUnlock()
    
    // Load from disk and store weak reference
    page := loadFromDisk(id)
    bp.mu.Lock()
    bp.pages[id] = weak.Make(page)
    bp.mu.Unlock()
    return page
}
```

**Investigation:** Compare LRU vs Weak cache under memory pressure. Which has better hit rates?

#### Challenge Requirements

- [ ] Support 1000 concurrent goroutines accessing pages
- [ ] Never exceed memory budget (configurable)
- [ ] Achieve >90% hit rate on Zipfian workload
- [ ] Background flusher writes dirty pages every 1s
- [ ] Zero data races (`go test -race` passes)

#### Verification Benchmark (Go 1.24+)

```go
func BenchmarkBufferPool(b *testing.B) {
    pool := NewBufferPool(100) // 100 pages in memory
    
    for b.Loop() {  // 🆕 Go 1.24: Cleaner than for i := 0; i < b.N
        pageID := zipfian.Uint64() % 10000
        page := pool.Get(pageID)
        // simulate work
        pool.Unpin(pageID)
    }
}
```
Target: >500K ops/sec on your machine

**🆕 Go 1.24 Bonus:** The underlying map uses Swiss Tables—measure if it's faster!

#### What You Should Discover

- LRU needs a doubly-linked list + hashmap
- Pin counts prevent use-after-free bugs
- Write-back caching is complex but essential
- 🆕 Go 1.24: Swiss Tables make large maps 30% faster
- 🆕 Go 1.24: Weak pointers enable memory-aware caching

---

### Lesson 1.3: Write-Ahead Log (Week 3)

**Your Mission:** Crash your program randomly during writes. Always recover correctly.

**The Scenario:**
```go
// This code must be crash-safe:
db.CreateNode(id: 1, name: "Alice")
db.CreateEdge(from: 1, to: 2)
// CRASH HERE - what's on disk?
```

#### Investigation Questions

1. **What Goes in the Log?**
   ```go
   type LogRecord struct {
       LSN    uint64  // Log Sequence Number
       TxnID  uint64
       Type   RecordType // BEGIN, INSERT, COMMIT
       // ??? what else?
   }
   ```
   - Do you log the old value (undo) or new value (redo)?
   - How do you handle multi-page operations?
   - Where do you store uncommitted data?

2. **Recovery Algorithm:**
   You wake up after a crash. The log contains:
   ```
   [LSN=1] BEGIN txn=10
   [LSN=2] INSERT page=5 offset=100 data="Alice"
   [LSN=3] INSERT page=5 offset=200 data="Bob"
   [LSN=4] BEGIN txn=11
   [LSN=5] INSERT page=6 offset=50 data="Carol"
   [LSN=6] COMMIT txn=10
   ```
   - Which records do you replay?
   - What if page 5 was already written to disk?
   - What about txn=11 (never committed)?

3. **Idempotency Challenge:**
   ```go
   // This must work if called twice:
   func (wal *WAL) Recover() error {
       // Read log
       // Apply changes
       // ???
   }
   ```
   - What if recovery crashes halfway through?
   - How do you know what's already been applied?
   - Research: What's a checkpoint?

#### 🆕 Go 1.23: Timer Improvements

Background flusher is now simpler:
```go
func (wal *WAL) StartBackgroundSync() {
    ticker := time.NewTicker(1 * time.Second)
    
    go func() {
        for range ticker.C {
            wal.Sync()
        }
    }()
    
    // 🆕 Go 1.23: No need to call ticker.Stop()
    // It will be GC'd automatically when WAL is GC'd
}
```

#### Your Implementation

- [ ] Log every write operation before modifying data pages
- [ ] `fsync()` the log on commit (measure this latency)
- [ ] Implement recovery that's idempotent
- [ ] Add chaos testing: inject crashes randomly

#### Chaos Test

```go
func TestCrashRecovery(t *testing.T) {
    for trial := 0; trial < 100; trial++ {
        db := NewDB("test.db")
        
        // Perform 1000 random operations
        done := make(chan struct{})
        go func() {
            time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
            os.Exit(137) // simulate SIGKILL
        }()
        
        // Write operations...
        close(done)
        
        // After "crash", verify all committed data intact
        db2 := NewDB("test.db")
        db2.Recover()
        // Verify...
    }
}
```

#### What You Should Discover

- `fsync()` is SLOW (5-10ms on SSDs)
- Group commit amortizes fsync cost
- ARIES algorithm has three phases: Analysis, Redo, Undo

---

## PHASE 2: GRAPH STRUCTURE

### Lesson 2.1: Compressed Sparse Row (Week 4)

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

#### Investigation Questions

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

#### 🆕 Go 1.23: Iterator-Based Traversal

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

#### Your Implementation

- [ ] Build CSR from edge list in O(E) time
- [ ] Query neighbors in O(degree(n)) time
- [ ] Measure cache misses (use `pprof` or `perf`)
- [ ] Support both forward and backward traversal
- [ ] 🆕 Implement iterator-based API

#### Benchmark (Go 1.24+)

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

#### 🆕 Go 1.25: Experimental GC Performance

Test with the new garbage collector:
```bash
GOEXPERIMENT=greenteagc go test -bench=BenchmarkTraversal
```

The new GC improves cache locality during marking. For pointer-heavy workloads (graphs!), expect 10-35% improvement.

#### What You Should Discover

- Sequential memory access is 10-100x faster
- CSR trades insert performance for query performance
- Cache line prefetching is automatic (and magical)
- 🆕 Iterators provide zero-cost abstractions
- 🆕 Experimental GC significantly helps graph traversals

---

### Lesson 2.2: Columnar Node Properties (Week 5)

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

#### Investigation Questions

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

#### 🆕 Go 1.23: String Interning with `unique` Package

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

#### Challenge

- [ ] Store 1M heterogeneous nodes in <50MB
- [ ] Query: "Find all ages > 50" without decompressing all columns
- [ ] Add new property type without rewriting all data
- [ ] Support NULL values efficiently
- [ ] 🆕 Use `unique.Make()` for low-cardinality columns

#### Benchmark Query (Go 1.24+)

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

#### What You Should Discover

- Column stores are 10x better for analytical queries
- Compression works better on columns (same type together)
- Sparse columns with mostly NULLs compress to almost nothing
- 🆕 `unique.Handle` dramatically reduces memory for repeated values
- 🆕 Go 1.24: Swiss Tables make dictionary encoding faster

---

### Lesson 2.3: NodeGroups and Morsel-Driven Parallelism (Week 6)

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

#### Investigation Questions

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

#### 🆕 Go 1.25: Container-Aware GOMAXPROCS

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

#### 🆕 Go 1.25: Testing Parallel Execution

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

#### Your Implementation

- [ ] Achieve near-linear speedup with cores (8 cores → 7.5x faster)
- [ ] Handle unbalanced workloads (some groups scan faster)
- [ ] Measure goroutine overhead (<5% of total time)
- [ ] Zero allocations in hot path (use `go test -benchmem`)
- [ ] 🆕 Use `synctest` for concurrent correctness testing

#### Scaling Test (Go 1.24+)

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

#### Benchmark (Go 1.24+)

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

#### What You Should Discover

- Goroutines are cheap but not free
- Work stealing beats static partitioning
- False sharing kills parallel performance
- CPU cache lines are 64 bytes (measure with `-cpu=1` vs `-cpu=8`)
- 🆕 `synctest` makes concurrency tests reproducible

---

## PHASE 3: QUERY ENGINE

### Lesson 3.1: Parse Cypher Queries (Week 7)

**Your Mission:** Turn this string into a data structure:
```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)
WHERE a.age > 30
RETURN b.name
```

#### Investigation Questions

1. **Tokenization:**
   ```go
   input := "MATCH (a:Person) WHERE a.age > 30"
   tokens := Tokenize(input)
   // Should produce: [MATCH, LPAREN, IDENT("a"), COLON, IDENT("Person"), ...]
   ```
   - How do you handle whitespace?
   - How do you distinguish `MATCH` (keyword) from `match` (identifier)?
   - What about string literals with special chars: `"foo\"bar"`?

2. **Parser Choice:**
   
   **Option A: Use ANTLR4**
   ```bash
   # Install ANTLR4 for Go
   go get github.com/antlr4-go/antlr/v4
   # Write Cypher.g4 grammar
   ```
   Pro: Handles complexity. Con: Big dependency, slow compilation.
   
   **Option B: Recursive Descent by Hand**
   ```go
   func (p *Parser) parseMatch() (*MatchClause, error) {
       p.expect(MATCH)
       pattern := p.parsePattern()
       // ...
   }
   ```
   Pro: Full control, fast. Con: You write more code.
   
   Try both. Which feels better for learning?

3. **AST Design:**
   ```go
   type Query struct {
       Match  *MatchClause
       Where  *WhereClause
       Return *ReturnClause
   }
   
   type MatchClause struct {
       Patterns []Pattern
   }
   
   type Pattern struct {
       // (a)-[:KNOWS]->(b)
       // How do you represent this?
   }
   ```
   Draw the AST for 5 different queries. What common structure emerges?

#### Challenge

- [ ] Parse 20 test queries correctly
- [ ] Produce useful error messages: "Expected RETURN, got WHERE on line 3"
- [ ] Handle nested expressions: `a.age > 30 AND (b.city = "NYC" OR b.city = "LA")`
- [ ] Pretty-print the AST for debugging

#### Test Cases

```go
tests := []string{
    "MATCH (n) RETURN n",
    "MATCH (a)-[r]->(b) WHERE a.age > 30 RETURN b.name",
    "MATCH (a)-[:KNOWS*1..3]->(b) RETURN count(b)",
    // Add 17 more edge cases
}
```

#### What You Should Discover

- Parsing is 20% lexing, 80% handling edge cases
- Good error messages are harder than parsing itself
- ASTs need methods: `node.String()`, `node.Validate()`

---

### Lesson 3.2: Query Planning (Week 8)

**Your Mission:** Turn an AST into an execution plan.

**The Problem:**
```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)-[:LIKES]->(c:Movie)
WHERE a.age > 30 AND c.year = 2020
RETURN c.title
```

**Many valid execution orders:**
1. Scan all Persons → filter age → traverse KNOWS → traverse LIKES → filter year
2. Scan Movies → filter year → reverse LIKES → reverse KNOWS → filter age
3. Hash join Person and Movie on intermediate LIKES edges

Which is fastest?

#### Investigation Questions

1. **Cost Estimation:**
   ```go
   type Operator interface {
       EstimatedCost() float64
       EstimatedCardinality() int
   }
   
   type SeqScan struct {
       Table string
       Filter Expr
   }
   
   func (s *SeqScan) EstimatedCost() float64 {
       // ???
   }
   ```
   - How do you estimate selectivity of `age > 30`?
   - Research: Histograms vs sampling vs bloom filters
   - What if you have no statistics yet?

2. **Plan Enumeration:**
   Three patterns: (a)->(b)->(c)
   
   **Possible join orders:**
   - Left-deep: ((a⋈b)⋈c)
   - Right-deep: (a⋈(b⋈c))
   - Bushy: Various combinations
   
   How many possible plans for N patterns? (It's factorial)
   
   **Optimization:**
   - Dynamic programming (PostgreSQL approach)
   - Greedy heuristic (pick smallest intermediate result)
   - Random sampling (try 100 plans, pick best)

3. **Predicate Pushdown:**
   ```go
   // Before optimization:
   SeqScan(Person) → Filter(age > 30) → Join(KNOWS)
   
   // After pushdown:
   SeqScan(Person, filter: age > 30) → Join(KNOWS)
   ```
   - When can you push filters down?
   - What about `WHERE a.age + b.age > 60`? (Needs both sides)

#### 🆕 Go 1.23: Intern Plan Operator Types

```go
import "unique"

type LogicalPlan struct {
    OpType unique.Handle[string]  // "SeqScan", "HashJoin", etc.
    // Most queries use same 10-20 operator types
}

// Fast operator comparison
func (p1 *LogicalPlan) SameOperator(p2 *LogicalPlan) bool {
    return p1.OpType == p2.OpType  // Pointer comparison!
}
```

#### Your Implementation

- [ ] Generate at least 3 different plans for same query
- [ ] Choose plan based on estimated cost
- [ ] Show plan with `EXPLAIN` command
- [ ] Measure actual vs estimated cardinality

#### Visualization

```
EXPLAIN MATCH (a:Person)-[:KNOWS]->(b)
  WHERE a.age > 30 RETURN b.name;

PhysicalPlan:
  Project(b.name)
  └─ HashJoin(a.id = KNOWS.src)
     ├─ SeqScan(Person, filter: age > 30)  [est: 100K rows]
     └─ IndexScan(KNOWS.src)                [est: 1M rows]
     
Estimated cost: 1234.5
```

#### What You Should Discover

- Query optimization is NP-hard
- Cardinality estimates are often wildly wrong
- Good plans are 100x faster than bad plans

---

### Lesson 3.3: Join Algorithms (Week 9)

**Your Mission:** Implement three join algorithms, benchmark them.

**The Query:**
```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)
```
Translation: Join Person table with KNOWS edge list on Person.id = KNOWS.dst

#### Investigation Questions

1. **Hash Join:**
   ```go
   // Phase 1: Build hash table on smaller side
   hashTable := make(map[NodeID]*Person, len(persons))  // 🆕 Pre-size!
   for i := range persons {
       hashTable[persons[i].ID] = &persons[i]
   }
   
   // Phase 2: Probe with larger side
   for _, edge := range edges {
       if person, found := hashTable[edge.Dst]; found {
           yield(edge.Src, person)
       }
   }
   ```
   - What if the hash table doesn't fit in memory?
   - Research: Grace hash join (partition-based)
   - When does Go's map become slow? (Measure at 1M, 10M, 100M entries)

2. **Sort-Merge Join:**
   ```go
   // Sort both sides
   sort.Slice(persons, func(i, j int) bool {
       return persons[i].ID < persons[j].ID
   })
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].Dst < edges[j].Dst
   })
   
   // Merge with two pointers
   i, j := 0, 0
   for i < len(persons) && j < len(edges) {
       // Your code here
   }
   ```
   - What if data is already sorted? (CSR is sorted!)
   - What about duplicate keys?
   - Measure: sort cost vs probe cost

3. **Worst-Case Optimal Join (Triangle Query):**
   ```cypher
   MATCH (a)-[:KNOWS]->(b)-[:KNOWS]->(c)-[:KNOWS]->(a)
   RETURN count(*)
   ```
   This is a cycle—hash joins explode.
   
   **Binary hash join:** O(E^1.5) intermediate results
   **Multiway intersection:** O(E) time
   
   ```go
   func intersectThree(a, b, c []NodeID) []NodeID {
       // How do you do this efficiently?
       // Hint: They're sorted
   }
   ```

#### 🆕 Go 1.24: Swiss Tables Performance

Hash joins got 30-35% faster!

```go
func BenchmarkHashJoin(b *testing.B) {
    persons := generatePersons(1_000_000)
    edges := generateEdges(10_000_000)
    
    b.Run("Build", func(b *testing.B) {
        for b.Loop() {
            // 🆕 35% faster with pre-sized map in Go 1.24
            hashTable := make(map[NodeID]*Person, len(persons))
            for i := range persons {
                hashTable[persons[i].ID] = &persons[i]
            }
        }
    })
    
    b.Run("Probe", func(b *testing.B) {
        hashTable := buildHashTable(persons)
        
        for b.Loop() {
            // 🆕 30% faster lookups in Go 1.24
            count := 0
            for _, edge := range edges {
                if _, found := hashTable[edge.Dst]; found {
                    count++
                }
            }
        }
    })
}
```

**Investigation:** Measure performance improvement from Go 1.23 to 1.24.

#### Your Implementation

- [ ] Hash join for 1-to-many
- [ ] Sort-merge for many-to-many
- [ ] Multiway intersection for cycles
- [ ] Adaptive: pick algorithm based on input size
- [ ] 🆕 Benchmark shows Swiss Tables improvement

#### Benchmark (Go 1.24+)

```go
// Dataset: 10K persons, 100K KNOWS edges
func BenchmarkJoins(b *testing.B) {
    scenarios := []struct{
        name string
        persons int
        edges int
    }{
        {"small", 1000, 10000},
        {"medium", 10000, 100000},
        {"large", 100000, 1000000},
    }
    
    for _, s := range scenarios {
        b.Run(s.name + "/hash", func(b *testing.B) {
            // hash join
        })
        b.Run(s.name + "/merge", func(b *testing.B) {
            // sort-merge join
        })
    }
}
```

#### What You Should Discover

- Hash joins are usually fastest (if memory allows)
- Sort-merge shines when data is pre-sorted
- Cyclic queries need special treatment (WCO joins)
- 🆕 Swiss Tables make hash joins 30-35% faster

---

### Lesson 3.4: Execution Engine (Week 10)

**Your Mission:** Execute plans with pipelining.

#### 🆕 Go 1.23: Iterator-Based Execution

Modern approach using range-over-func:

```go
import "iter"

type Operator interface {
    Execute() iter.Seq[Row]
}

type SeqScan struct {
    table *NodeTable
}

func (s *SeqScan) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for i := 0; i < s.table.NumRows(); i++ {
            if !yield(s.table.Row(i)) {
                return  // Consumer stopped
            }
        }
    }
}

// Composable pipeline:
type Filter struct {
    child Operator
    predicate func(Row) bool
}

func (f *Filter) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for row := range f.child.Execute() {
            if f.predicate(row) {
                if !yield(row) {
                    return
                }
            }
        }
    }
}

// Usage:
scan := &SeqScan{table: persons}
filter := &Filter{child: scan, predicate: func(r Row) bool {
    return r.Age > 30
}}

for row := range filter.Execute() {
    process(row)
}
```

#### Investigation Questions

1. **Iterator vs Vectorized:**
   ```go
   // Iterator: one row at a time
   for row := range operator.Execute() {
       process(row)
   }
   
   // Vectorized: batch of rows
   for {
       batch := operator.NextBatch()  // 1000 rows
       if batch == nil { break }
       for _, row := range batch {
           process(row)
       }
   }
   ```
   Implement both. Measure virtual function call overhead.

2. **Pipeline Breakers:**
   Some operators can't pipeline:
   - Sort (needs all input before outputting)
   - Hash join build phase (needs all input for hash table)
   - Aggregation (GROUP BY needs all rows)
   
   ```go
   type SortOperator struct {
       child Operator
       buffer []Row  // Must accumulate here
   }
   
   func (s *SortOperator) Execute() iter.Seq[Row] {
       return func(yield func(Row) bool) {
           // Collect all rows
           for row := range s.child.Execute() {
               s.buffer = append(s.buffer, row)
           }
           
           // Sort
           sort.Slice(s.buffer, ...)
           
           // Yield sorted
           for _, row := range s.buffer {
               if !yield(row) {
                   return
               }
           }
       }
   }
   ```
   - How do you detect pipeline breaks?
   - Can you parallelize across pipeline boundaries?

3. **Parallel Execution:**
   ```go
   type ParallelSeqScan struct {
       table string
       workers int
       chunks chan NodeGroupID
   }
   
   func (p *ParallelSeqScan) Execute() iter.Seq[Row] {
       return func(yield func(Row) bool) {
           results := make(chan Row, 1000)
           
           // Start workers
           var wg sync.WaitGroup
           for i := 0; i < p.workers; i++ {
               wg.Add(1)
               go func() {
                   defer wg.Done()
                   for chunk := range p.chunks {
                       for row := range scanChunk(chunk) {
                           results <- row
                       }
                   }
               }()
           }
           
           // Close results when done
           go func() {
               wg.Wait()
               close(results)
           }()
           
           // Yield results
           for row := range results {
               if !yield(row) {
                   // Early exit - drain channels
                   return
               }
           }
       }
   }
   ```
   - Each worker scans different NodeGroups
   - Results must be thread-safe
   - Use sync.WaitGroup or channels?

#### Your Implementation

- [ ] Iterator-based execution (🆕 Go 1.23)
- [ ] Support 10+ operator types (scan, filter, join, project, limit)
- [ ] Pipeline as much as possible
- [ ] Parallelize scans across goroutines
- [ ] Profile: what % of time in each operator?

#### Query Example

```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)
WHERE a.age > 30 AND b.city = "NYC"
RETURN b.name, count(*) as friends
LIMIT 10
```

**Execution plan:**
```
Limit(10)
└─ Aggregate(b.name, count(*))
   └─ Filter(b.city = "NYC")
      └─ HashJoin(a.id = KNOWS.dst)
         ├─ Filter(a.age > 30)
         │  └─ SeqScan(Person)  [parallel: 8 workers]
         └─ SeqScan(KNOWS)     [parallel: 8 workers]
```

#### What You Should Discover

- Iterator model is simple but has overhead
- Vectorization makes modern CPUs happy
- Parallelism is easy; correct parallelism is hard
- 🆕 Range-over-func provides clean composition

---

## PHASE 4: TRANSACTIONS

### Lesson 4.1: Locking Protocols (Week 11)

**Your Mission:** Two transactions run concurrently. Prevent anomalies.

**The Scenario:**
```go
// Transaction 1:
tx1.Execute("MATCH (a {id: 1}) SET a.balance = a.balance - 100")

// Transaction 2 (concurrent):
tx2.Execute("MATCH (a {id: 1}) RETURN a.balance")
```

What should tx2 see?

#### Investigation Questions

1. **Two-Phase Locking:**
   ```go
   type LockManager struct {
       locks map[NodeID]*sync.RWMutex
   }
   
   func (tx *Transaction) ReadNode(id NodeID) {
       tx.locks.RLock(id)  // Shared lock
       defer tx.locks.RUnlock(id)
       // But when do you unlock?
   }
   ```
   - Growing phase: Acquire locks
   - Shrinking phase: Release locks
   - What if you release early? (Lost updates)
   - What if you never release? (Deadlock)

2. **Deadlock Detection:**
   ```
   TX1: Lock(A) → waiting for Lock(B)
   TX2: Lock(B) → waiting for Lock(A)
   ```
   Build a wait-for graph:
   ```go
   type WaitForGraph struct {
       edges map[TxnID][]TxnID
   }
   
   func (g *WaitForGraph) HasCycle() bool {
       // How do you detect cycles efficiently?
   }
   ```
   - Run cycle detection every N seconds?
   - Or use timeout (pessimistic)?
   - Who do you abort when cycle found?

3. **Isolation Levels:**
   ```go
   type IsolationLevel int
   const (
       ReadUncommitted
       ReadCommitted
       RepeatableRead
       Serializable
   )
   ```
   Implement each level. What locking behavior changes?
   - Read Uncommitted: No read locks
   - Read Committed: Short-duration read locks
   - Repeatable Read: Hold read locks until commit
   - Serializable: Predicate locks

#### 🆕 Go 1.25: Deterministic Deadlock Testing

```go
import "testing/synctest"

func TestDeadlockDetection(t *testing.T) {
    synctest.Run(func() {
        db := NewDB()
        deadlockDetected := atomic.Bool{}
        
        // Transaction 1: Lock A → B
        go func() {
            tx := db.Begin()
            tx.Lock(NodeID(1))
            time.Sleep(10 * time.Millisecond)  // Fake time!
            
            if err := tx.Lock(NodeID(2)); err == ErrDeadlock {
                deadlockDetected.Store(true)
            }
        }()
        
        // Transaction 2: Lock B → A
        go func() {
            tx := db.Begin()
            tx.Lock(NodeID(2))
            time.Sleep(10 * time.Millisecond)
            
            if err := tx.Lock(NodeID(1)); err == ErrDeadlock {
                deadlockDetected.Store(true)
            }
        }()
        
        synctest.Wait()
        
        if !deadlockDetected.Load() {
            t.Fatal("Deadlock was not detected")
        }
    })
}
```

This test runs in **fake time** and is **deterministic**—no more flaky tests!

#### Your Implementation

- [ ] Acquire locks in consistent order (prevent deadlock)
- [ ] Detect deadlocks via timeout or graph analysis
- [ ] Support different isolation levels
- [ ] Show current locks: `SHOW LOCKS;`
- [ ] 🆕 Use `synctest` for deterministic testing

#### Test Cases (Go 1.25+)

```go
func TestConcurrentTransactions(t *testing.T) {
    synctest.Run(func() {
        db := NewDB()
        
        // Lost update test
        var balance atomic.Int64
        balance.Store(1000)
        
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                tx := db.Begin()
                // Read, increment, write
                current := balance.Load()
                balance.Store(current + 10)
                tx.Commit()
            }()
        }
        wg.Wait()
        
        // Should be: 1000 + (100 * 10) = 2000
        // Without locking: race condition!
        if balance.Load() != 2000 {
            t.Errorf("Lost updates: got %d, want 2000", balance.Load())
        }
    })
}
```

#### What You Should Discover

- 2PL prevents anomalies but reduces concurrency
- Deadlocks are inevitable with random lock order
- Isolation levels trade correctness for performance
- 🆕 `synctest` makes concurrency bugs reproducible

---

### Lesson 4.2: MVCC (Multi-Version Concurrency Control) (Week 12)

**Your Mission:** Readers don't block writers. Writers don't block readers.

**The Idea:**
```go
type Node struct {
    ID       NodeID
    Versions []Version
}

type Version struct {
    TxnID     uint64
    Timestamp uint64
    Data      []byte
    Deleted   bool
}
```

Each update creates a new version instead of overwriting.

#### Investigation Questions

1. **Version Visibility:**
   ```go
   func (tx *Transaction) ReadNode(id NodeID) *Version {
       node := db.GetNode(id)
       // Which version should this transaction see?
       // Options:
       // - Latest committed before tx.StartTS
       // - Latest committed <= tx.SnapshotTS
       // - Something else?
   }
   ```
   - How do you assign timestamps?
   - What if clocks are not synchronized?
   - Research: Hybrid Logical Clocks

2. **Garbage Collection:**
   ```go
   // Node has 1000 versions from old transactions
   type Node struct {
       Versions []Version  // Growing unbounded!
   }
   ```
   - When can you delete old versions?
   - How do you know no transaction needs version 42?
   - Background GC vs inline GC?

3. **Write-Write Conflicts:**
   ```
   TX1: Read node A (version 5)
   TX2: Read node A (version 5)
   TX1: Write node A (create version 6) → commit
   TX2: Write node A (create version ?) → should this succeed?
   ```
   - First-committer-wins rule
   - Track read set vs write set
   - Implement optimistic concurrency control

#### 🆕 Go 1.24: Weak Pointers for Old Versions

```go
import "weak"

type VersionChain struct {
    current *Version
    old     []weak.Pointer[*Version]  // Can be GC'd
}

func (vc *VersionChain) GetVersion(timestamp uint64) *Version {
    if vc.current.Timestamp <= timestamp {
        return vc.current
    }
    
    // Check old versions (might be GC'd)
    for i := len(vc.old) - 1; i >= 0; i-- {
        if v := vc.old[i].Value(); v != nil {
            if v.Timestamp <= timestamp {
                return v
            }
        }
    }
    
    return nil  // Too old, was GC'd
}

func (vc *VersionChain) AddVersion(v *Version) {
    // Move current to old versions as weak pointer
    vc.old = append(vc.old, weak.Make(vc.current))
    vc.current = v
}
```

**Investigation:** Does weak pointer GC reduce memory usage? Measure!

#### 🆕 Go 1.25: Test Version Visibility

```go
func TestSnapshotIsolation(t *testing.T) {
    synctest.Run(func() {
        db := NewMVCCDB()
        
        // TX1: Read at T=0
        tx1 := db.BeginAt(0)
        val1 := tx1.Read(NodeID(1))
        
        // TX2: Write at T=5
        tx2 := db.BeginAt(5)
        tx2.Write(NodeID(1), "new value")
        tx2.Commit()
        
        // TX1 should still see old value
        val2 := tx1.Read(NodeID(1))
        
        if val1 != val2 {
            t.Error("Snapshot isolation violated")
        }
    })
}
```

#### Your Implementation

- [ ] Each write creates new version
- [ ] Reads see snapshot at transaction start
- [ ] Detect write-write conflicts
- [ ] GC old versions safely
- [ ] Benchmark: concurrent reads while writing
- [ ] 🆕 Optional: Use weak pointers for memory efficiency

#### Benchmark (Go 1.24+)

```go
func BenchmarkMVCC(b *testing.B) {
    db := NewDB("test.db")
    nodeID := NodeID(42)
    
    // One writer
    go func() {
        for {
            tx := db.Begin()
            tx.Update(nodeID, generateValue())
            tx.Commit()
            time.Sleep(1 * time.Millisecond)
        }
    }()
    
    // Many readers (should not block)
    for b.Loop() {
        tx := db.Begin()
        _ = tx.Read(nodeID)
        tx.Commit()
    }
}
```
Readers should not slow down during writes.

#### What You Should Discover

- MVCC trades disk space for concurrency
- Version chains grow quickly under write load
- Snapshot isolation is not serializable (write skew)
- 🆕 Weak pointers allow memory-aware version GC

---

## BONUS CHALLENGES

### Challenge A: Query Optimizer Showdown

**Implement competing optimization strategies:**

1. **Exhaustive enumeration** (try all join orders)
2. **Greedy heuristic** (always pick smallest intermediate)
3. **Genetic algorithm** (evolve query plans)

Run 100 random queries. Which wins most often?

### Challenge B: Adaptive Indexing

**Auto-create indexes based on query patterns:**
```go
type QueryMonitor struct {
    slowQueries []Query
}

func (m *QueryMonitor) ShouldIndex() (table, column string) {
    // Analyze which predicates appear frequently
    // E.g., "WHERE age > X" appears 1000 times
    // → Create index on Person.age
}
```

### Challenge C: Distributed Execution

**Shard the graph across 3 Go processes:**
- Node 1-1M on server A
- Node 1M-2M on server B  
- Node 2M-3M on server C

How do you execute cross-shard queries?

### 🆕 Challenge D: Experimental GC Evaluation

**Test Go 1.25's new garbage collector:**

```bash
# Build with experimental GC
GOEXPERIMENT=greenteagc go build

# Run your graph traversal benchmarks
./kuzu-go benchmark --workload=graph-traversal

# Compare with standard GC
go build  # without GOEXPERIMENT
./kuzu-go benchmark --workload=graph-traversal
```

**Questions:**
- Does the new GC reduce cache misses?
- What's the throughput improvement?
- When does it help most?

---

## YOUR PROGRESS TRACKING

After each lesson, write a document answering:

1. **What surprised you?** (The thing you didn't expect)
2. **What was harder than expected?** (Where you got stuck)
3. **What optimization made the biggest difference?** (The 10x win)
4. **What would you do differently next time?** (Lessons learned)
5. **🆕 How did Go 1.23+ features help?** (Iterator benefits, `unique` savings, etc.)

These reflections are where the real learning happens.

---

## RECOMMENDED RESOURCES

### Books
- **Database Internals** by Alex Petrov - Best modern database book
- **Designing Data-Intensive Applications** by Martin Kleppmann - Systems thinking
- **The Go Programming Language** by Donovan & Kernighan - Go fundamentals

### Online Courses
- **CMU 15-445: Database Systems** (YouTube) - Andy Pavlo's legendary course
- **CMU 15-721: Advanced Database Systems** - Deep dives into modern techniques

### Papers
- Kuzu's CIDR 2023 paper - Learn their design decisions
- "Worst-Case Optimal Join Algorithms" - Understand WCO joins
- "ARIES: A Transaction Recovery Method" - WAL and recovery

### Codebases to Study
- **DuckDB** - Similar embedded analytical DB in C++, clean code
- **SQLite** - The gold standard for embedded databases
- **Dgraph** - Production graph database written in Go

### Go-Specific Resources
- **Go 1.23 Release Notes** - https://go.dev/doc/go1.23
- **Go 1.24 Release Notes** - https://go.dev/doc/go1.24
- **Go 1.25 Release Notes** - https://go.dev/doc/go1.25
- **Range-over-func Proposal** - Understanding iterators
- **Swiss Tables Deep Dive** - How Go's new map works

### Tools
- `go tool pprof` - CPU and memory profiling
- `go tool trace` - Visualize goroutine execution
- `go test -race` - Detect data races
- `go test -benchmem` - Measure allocations
- `perf` (Linux) - Hardware performance counters
- `GOEXPERIMENT` - Enable experimental features

---

## PROJECT STRUCTURE RECOMMENDATION

```
kuzu-go/
├── cmd/
│   └── kuzu/
│       └── main.go              # CLI interface
├── storage/
│   ├── page.go                  # Page abstraction
│   ├── buffer_pool.go           # Buffer manager
│   ├── wal.go                   # Write-ahead log
│   └── storage_test.go
├── graph/
│   ├── csr.go                   # Compressed Sparse Row
│   ├── column_store.go          # Columnar properties
│   ├── node_group.go            # NodeGroup parallelism
│   └── graph_test.go
├── query/
│   ├── lexer.go                 # Tokenizer
│   ├── parser.go                # AST builder
│   ├── planner.go               # Query optimizer
│   ├── executor.go              # Execution engine
│   ├── operators.go             # Physical operators
│   └── query_test.go
├── index/
│   ├── hash.go                  # Hash index
│   └── btree.go                 # B-tree (optional)
├── transaction/
│   ├── lock_manager.go          # 2PL implementation
│   ├── mvcc.go                  # Multi-version CC
│   └── transaction_test.go
├── go.mod                       # go 1.25
├── go.sum
└── README.md
```

---

## GETTING STARTED CHECKLIST

### Day 1: Setup
- [ ] Install Go 1.25
- [ ] Set up project structure
- [ ] Verify `go test` runs
- [ ] Enable experimental GC (optional): `export GOEXPERIMENT=greenteagc`

### Week 1: Storage Foundation
- [ ] Complete Lesson 1.1 (Pages)
- [ ] Complete Lesson 1.2 (Buffer Pool)
- [ ] Complete Lesson 1.3 (WAL)

### Week 2-3: Graph Structures
- [ ] Complete Lesson 2.1 (CSR) - **Use iterators!**
- [ ] Complete Lesson 2.2 (Columns) - **Try `unique` package**
- [ ] Complete Lesson 2.3 (Parallelism) - **Test with `synctest`**

### Week 4-5: Query Engine
- [ ] Complete Lesson 3.1 (Parser)
- [ ] Complete Lesson 3.2 (Planner)
- [ ] Complete Lesson 3.3 (Joins) - **Measure Swiss Tables improvement**
- [ ] Complete Lesson 3.4 (Executor) - **Iterator-based operators**

### Week 6: Transactions
- [ ] Complete Lesson 4.1 (Locking) - **Use `synctest` for deadlock tests**
- [ ] Complete Lesson 4.2 (MVCC) - **Try weak pointers**

### Week 7+: Advanced Topics
- [ ] Choose 1-2 bonus challenges
- [ ] Benchmark with experimental GC
- [ ] Profile and optimize hot paths
- [ ] Write a comprehensive README

---

## VERSION COMPATIBILITY MATRIX

| Feature | Go 1.23 | Go 1.24 | Go 1.25 | Impact |
|---------|---------|---------|---------|--------|
| Range-over-func iterators | ✅ | ✅ | ✅ | High |
| `unique` package | ✅ | ✅ | ✅ | High |
| Timer improvements | ✅ | ✅ | ✅ | Low |
| Swiss Tables maps | ❌ | ✅ | ✅ | High |
| Weak pointers | ❌ | ✅ | ✅ | Medium |
| `testing.B.Loop` | ❌ | ✅ | ✅ | Low |
| `testing/synctest` | ❌ | 🧪 | ✅ | High |
| Container-aware GOMAXPROCS | ❌ | ❌ | ✅ | Medium |
| Experimental new GC | ❌ | ❌ | 🧪 | High |

**Legend:**
- ✅ Stable
- 🧪 Experimental (requires GOEXPERIMENT)
- ❌ Not available

---

## PERFORMANCE EXPECTATIONS

### Expected Improvements by Version

**Lesson 2.1 (CSR Traversal):**
- Go 1.23: Baseline with iterators
- Go 1.24: +0-5% (Swiss Tables for metadata)
- Go 1.25 + GreenTeaGC: **+10-35%** (better cache locality)

**Lesson 3.3 (Hash Join):**
- Go 1.23: Baseline
- Go 1.24: **+30-35%** 🔥 (Swiss Tables)
- Go 1.25: Same as 1.24

**Lesson 4.2 (MVCC):**
- Go 1.23: Baseline
- Go 1.24: +5-10% (weak pointers)
- Go 1.25 + GreenTeaGC: **+15-25%** (faster GC)

---

## FINAL THOUGHTS

**This curriculum uses cutting-edge Go features** to teach database concepts. You'll learn:
- Modern Go idioms (iterators, generics, weak pointers)
- Database internals (storage, indexing, query processing)
- Concurrent programming (with deterministic testing!)
- Performance engineering (profiling, benchmarking, optimization)

**Start with Lesson 1.1 and build incrementally.** Each lesson provides new challenges without giving solutions. You'll discover the answers through measurement and experimentation.

**The learning path adapts to your Go version:**
- Using Go 1.23? Focus on iterators and `unique`
- Using Go 1.24? Leverage Swiss Tables and weak pointers
- Using Go 1.25? Test everything with `synctest` and the new GC

**Most importantly: Have fun!** Building a database from scratch is one of the most rewarding learning experiences in computer science.

---

**Ready to start? Begin with Lesson 1.1!** 🚀

```bash
mkdir storage
touch storage/page.go
go test ./storage
# Let the journey begin...
```
