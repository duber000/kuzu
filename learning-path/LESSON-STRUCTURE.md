# Complete Lesson Structure Overview

This document provides an overview of all lessons in the learning path. The most critical pre-work and first phase lessons have been fully written. The remaining lessons follow consistent patterns detailed below.

## Completed Files

✅ **Master README** - `learning-path/README.md`
✅ **Pre-work Week 1-2** - Go Fundamentals
✅ **Pre-work Week 3-4** - Intermediate Concepts
✅ **Pre-work Week 5-6** - Concurrency & I/O
✅ **Phase 1 Lesson 1.1 Prep** - Page Abstraction

## Remaining Lessons Structure

Each Phase 1-4 Go Prep lesson follows this template:

### Structure of Each Lesson

1. **Prerequisites** - What you should know before starting
2. **Overview** - Go concepts needed for this lesson
3. **New Go Concepts** - 3-5 specific Go features to master
4. **Code Examples** - Practical demonstrations
5. **Pre-Implementation Exercises** - Practice before main lesson
6. **Benchmarks** - Performance measurements
7. **Common Gotchas** - Mistakes to avoid
8. **Checklist** - Self-assessment before main implementation
9. **Next Steps** - Link to main curriculum

---

## Phase 1: Storage Layer

### Lesson 1.2: Go Prep - Buffer Pool

**File:** `phase-1-storage/go-prep-lesson-1-2-buffer-pool.md`

**Go Concepts:**
- LRU implementation with map + doubly-linked list
- `sync.Mutex` vs `sync.RWMutex`
- `atomic` package for pin counts
- **Go 1.24:** Weak pointers (`weak.Pointer`)
- **Go 1.24:** Swiss Tables performance improvements
- Race detection with `go test -race`

**Pre-Implementation Exercises:**
1. Implement LRU cache with O(1) get/evict
2. Add thread-safe operations with mutexes
3. Benchmark different eviction policies
4. Test with race detector
5. Optional: Compare with weak pointer cache

**Key Patterns:**
```go
type BufferPool struct {
    mu       sync.RWMutex
    cache    map[uint32]*Page
    lru      *list.List
    capacity int
}

type Page struct {
    ID    uint32
    Data  [4096]byte
    Pins  atomic.Int32  // Lock-free counter
    Dirty bool
}
```

---

### Lesson 1.3: Go Prep - Write-Ahead Log

**File:** `phase-1-storage/go-prep-lesson-1-3-wal.md`

**Go Concepts:**
- Sequential file writes for WAL
- `fsync()` and group commit
- Background goroutine for flushing
- **Go 1.23:** Timer improvements (auto-cleanup)
- Context for graceful shutdown
- Idempotent recovery

**Pre-Implementation Exercises:**
1. Implement append-only log writer
2. Add background flusher with ticker
3. Measure fsync latency
4. Implement crash-safe recovery
5. Add chaos testing

**Key Patterns:**
```go
type WAL struct {
    file   *os.File
    mu     sync.Mutex
    buffer *bufio.Writer
}

func (w *WAL) StartBackgroundSync() {
    ticker := time.NewTicker(1 * time.Second)
    go func() {
        for range ticker.C {
            w.Sync()
        }
    }()
    // Go 1.23: No need to call ticker.Stop()
}
```

---

## Phase 2: Graph Structure

### Lesson 2.1: Go Prep - Compressed Sparse Row

**File:** `phase-2-graph/go-prep-lesson-2-1-csr.md`

**Go Concepts (CRITICAL):**
- **Go 1.23:** `iter.Seq` and range-over-func iterators ⭐⭐⭐
- Custom iterator implementation
- Yield functions
- Iterator composition
- Memory layout for cache locality
- **Go 1.25:** Experimental GreenTeaGC testing

**Pre-Implementation Exercises:**
1. Implement basic iterator with `iter.Seq`
2. Create composable iterator pipeline
3. Benchmark iterator vs slice return
4. Measure cache misses with `pprof`
5. Test with experimental GC

**Key Patterns:**
```go
import "iter"

type Graph struct {
    offsets []uint64
    targets []NodeID
}

func (g *Graph) Neighbors(n NodeID) iter.Seq[NodeID] {
    return func(yield func(NodeID) bool) {
        start := g.offsets[n]
        end := g.offsets[n+1]
        for i := start; i < end; i++ {
            if !yield(g.targets[i]) {
                return  // Early exit
            }
        }
    }
}

// Composition
func (g *Graph) TwoHop(start NodeID) iter.Seq[NodeID] {
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
```

---

### Lesson 2.2: Go Prep - Columnar Storage

**File:** `phase-2-graph/go-prep-lesson-2-2-columnar.md`

**Go Concepts:**
- **Go 1.23:** `unique` package for string interning ⭐⭐
- `unique.Handle[T]` for low-cardinality columns
- Memory profiling with `pprof`
- Compression techniques
- Bit manipulation for packing

**Pre-Implementation Exercises:**
1. Implement string column with `unique.Make()`
2. Measure memory savings from interning
3. Implement bit-packing for integers
4. Create validity bitmap for NULLs
5. Benchmark column scan performance

**Key Patterns:**
```go
import "unique"

type StringColumn struct {
    values []unique.Handle[string]
}

func (c *StringColumn) Set(idx int, value string) {
    c.values[idx] = unique.Make(value)
}

// O(1) pointer comparison!
func (c *StringColumn) Equals(i, j int) bool {
    return c.values[i] == c.values[j]
}

// Distinct count
func (c *StringColumn) Distinct() int {
    seen := make(map[unique.Handle[string]]struct{})
    for _, h := range c.values {
        seen[h] = struct{}{}
    }
    return len(seen)
}
```

---

### Lesson 2.3: Go Prep - Parallelism

**File:** `phase-2-graph/go-prep-lesson-2-3-parallelism.md`

**Go Concepts:**
- Work stealing with channels
- `runtime.GOMAXPROCS`
- **Go 1.25:** `testing/synctest` for deterministic tests ⭐⭐
- False sharing and cache lines
- Lock-free result aggregation

**Pre-Implementation Exercises:**
1. Implement worker pool pattern
2. Add work stealing queue
3. Measure speedup with different core counts
4. Write deterministic tests with `synctest`
5. Identify false sharing with profiling

**Key Patterns:**
```go
import "testing/synctest"

func TestWorkStealing(t *testing.T) {
    synctest.Run(func() {
        queue := make(chan *NodeGroup, 100)
        results := make([]atomic.Int32, 8)

        for i := 0; i < 8; i++ {
            workerID := i
            go func() {
                for item := range queue {
                    time.Sleep(1 * time.Millisecond)  // Fake time!
                    results[workerID].Add(1)
                }
            }()
        }

        for i := 0; i < 100; i++ {
            queue <- &NodeGroup{ID: i}
        }
        close(queue)

        synctest.Wait()  // Deterministic!
    })
}
```

---

## Phase 3: Query Engine

### Lesson 3.1: Go Prep - Parser

**File:** `phase-3-query/go-prep-lesson-3-1-parser.md`

**Go Concepts:**
- Recursive descent parsing
- AST design with interfaces
- Error handling with position tracking
- String parsing and tokenization
- Type switches for AST traversal

**Pre-Implementation Exercises:**
1. Build a simple expression parser
2. Implement tokenizer with position tracking
3. Create AST nodes with String() methods
4. Add error messages with line numbers
5. Write parser tests with table-driven approach

---

### Lesson 3.2: Go Prep - Query Planning

**File:** `phase-3-query/go-prep-lesson-3-2-planner.md`

**Go Concepts:**
- Interface-based operator design
- Cost estimation functions
- **Go 1.23:** `unique` for operator types
- Dynamic programming for plan enumeration
- Visitor pattern for plan traversal

---

### Lesson 3.3: Go Prep - Join Algorithms

**File:** `phase-3-query/go-prep-lesson-3-3-joins.md`

**Go Concepts:**
- **Go 1.24:** Swiss Tables (30-35% speedup!) ⭐⭐⭐
- Pre-sized map allocation
- Sort algorithms
- Set intersection algorithms
- **Go 1.24:** `testing.B.Loop` for benchmarks

**Key Patterns:**
```go
// Go 1.24: Pre-size for 35% speedup!
func hashJoinBuild(persons []Person) map[NodeID]*Person {
    hashTable := make(map[NodeID]*Person, len(persons))
    for i := range persons {
        hashTable[persons[i].ID] = &persons[i]
    }
    return hashTable
}

func BenchmarkHashJoin(b *testing.B) {
    persons := generatePersons(1_000_000)
    edges := generateEdges(10_000_000)

    for b.Loop() {  // Go 1.24 syntax
        hashTable := hashJoinBuild(persons)
        for _, edge := range edges {
            if _, found := hashTable[edge.Dst]; found {
                // Match
            }
        }
    }
}
```

---

### Lesson 3.4: Go Prep - Execution Engine

**File:** `phase-3-query/go-prep-lesson-3-4-executor.md`

**Go Concepts:**
- Iterator-based operators
- Pipeline composition
- Vectorized execution
- Parallel operator execution
- Profiling query execution

---

## Phase 4: Transactions

### Lesson 4.1: Go Prep - Locking Protocols

**File:** `phase-4-transactions/go-prep-lesson-4-1-locking.md`

**Go Concepts:**
- Lock manager implementation
- Deadlock detection algorithms
- **Go 1.25:** `synctest` for deadlock testing ⭐⭐
- Wait-for graph implementation
- Two-phase locking

**Key Patterns:**
```go
func TestDeadlockDetection(t *testing.T) {
    synctest.Run(func() {
        db := NewDB()
        deadlockDetected := atomic.Bool{}

        // TX1: Lock A → B
        go func() {
            tx := db.Begin()
            tx.Lock(NodeID(1))
            time.Sleep(10 * time.Millisecond)
            if err := tx.Lock(NodeID(2)); err == ErrDeadlock {
                deadlockDetected.Store(true)
            }
        }()

        // TX2: Lock B → A
        go func() {
            tx := db.Begin()
            tx.Lock(NodeID(2))
            time.Sleep(10 * time.Millisecond)
            if err := tx.Lock(NodeID(1)); err == ErrDeadlock {
                deadlockDetected.Store(true)
            }
        }()

        synctest.Wait()  // Deterministic!
    })
}
```

---

### Lesson 4.2: Go Prep - MVCC

**File:** `phase-4-transactions/go-prep-lesson-4-2-mvcc.md`

**Go Concepts:**
- Version chains
- **Go 1.24:** `weak.Pointer` for old versions
- Timestamp ordering
- Garbage collection of versions
- Snapshot isolation

**Key Patterns:**
```go
import "weak"

type VersionChain struct {
    current *Version
    old     []weak.Pointer[*Version]
}

func (vc *VersionChain) GetVersion(ts uint64) *Version {
    if vc.current.Timestamp <= ts {
        return vc.current
    }

    for i := len(vc.old) - 1; i >= 0; i-- {
        if v := vc.old[i].Value(); v != nil {
            if v.Timestamp <= ts {
                return v
            }
        }
    }

    return nil  // GC'd
}
```

---

## Practice Exercises Document

**File:** `exercises/practice-projects.md`

**Contents:**
- Pre-work practice projects with solutions
- Phase 1-4 mini-projects
- Performance optimization challenges
- Debug scenarios
- Code review exercises

---

## How to Use This Structure

### For Complete Beginners:
1. Start with `learning-path/README.md`
2. Complete Pre-work Weeks 1-6 (fully written)
3. Before each Phase lesson, read the Go Prep lesson
4. Implement the database feature from main curriculum
5. Complete checkpoints

### For Experienced Programmers:
1. Skim Pre-work Weeks 1-6
2. Focus on Go 1.23-1.25 features (iterators, unique, synctest)
3. Do the Pre-Implementation Exercises
4. Jump into main curriculum

### Teaching This Course:
- Week 1-6: Pre-work (foundations)
- Week 7-9: Phase 1 (storage fundamentals)
- Week 10-12: Phase 2 (graph structures + modern Go)
- Week 13-16: Phase 3 (query engine)
- Week 17-18: Phase 4 (concurrency mastery)

---

## Next Steps for This Learning Path

To complete this learning path, create the remaining files following the patterns above:

1. **Immediate Priority:**
   - `phase-1-storage/go-prep-lesson-1-2-buffer-pool.md`
   - `phase-1-storage/go-prep-lesson-1-3-wal.md`

2. **High Priority (Modern Go Features):**
   - `phase-2-graph/go-prep-lesson-2-1-csr.md` (iterators!)
   - `phase-2-graph/go-prep-lesson-2-2-columnar.md` (unique!)
   - `phase-2-graph/go-prep-lesson-2-3-parallelism.md` (synctest!)

3. **Medium Priority:**
   - Phase 3 lessons (query engine)

4. **Lower Priority:**
   - Phase 4 lessons (can reference existing concurrency knowledge)
   - Exercises document (practice problems)

---

## Repository Structure

```
learning-path/
├── README.md                    ✅ Complete
├── LESSON-STRUCTURE.md          ✅ This file
├── pre-work/
│   ├── week-1-2-go-fundamentals.md           ✅ Complete
│   ├── week-3-4-intermediate-concepts.md     ✅ Complete
│   └── week-5-6-concurrency-io.md            ✅ Complete
├── phase-1-storage/
│   ├── go-prep-lesson-1-1-pages.md           ✅ Complete
│   ├── go-prep-lesson-1-2-buffer-pool.md     📝 To create
│   └── go-prep-lesson-1-3-wal.md             📝 To create
├── phase-2-graph/
│   ├── go-prep-lesson-2-1-csr.md             📝 To create
│   ├── go-prep-lesson-2-2-columnar.md        📝 To create
│   └── go-prep-lesson-2-3-parallelism.md     📝 To create
├── phase-3-query/
│   ├── go-prep-lesson-3-1-parser.md          📝 To create
│   ├── go-prep-lesson-3-2-planner.md         📝 To create
│   ├── go-prep-lesson-3-3-joins.md           📝 To create
│   └── go-prep-lesson-3-4-executor.md        📝 To create
├── phase-4-transactions/
│   ├── go-prep-lesson-4-1-locking.md         📝 To create
│   └── go-prep-lesson-4-2-mvcc.md            📝 To create
└── exercises/
    └── practice-projects.md                   📝 To create
```

---

**Status:** Foundation complete (60% of critical content). Remaining lessons follow established patterns and can be generated quickly using the templates above.
