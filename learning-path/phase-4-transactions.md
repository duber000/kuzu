# Phase 4: Transactions

**Prerequisites:** Complete Phase 3 (Query Engine)

**Time Estimate:** 2 weeks
**Go Version:** 1.23+

In this final phase, you'll implement transaction support with locking protocols and multi-version concurrency control (MVCC). This ensures your database can handle concurrent users safely and efficiently.

---

## Table of Contents

- [Lesson 4.1: Locking Protocols](#lesson-41-locking-protocols-week-11)
- [Lesson 4.2: MVCC (Multi-Version Concurrency Control)](#lesson-42-mvcc-multi-version-concurrency-control-week-12)

---

## Lesson 4.1: Locking Protocols (Week 11)

**Your Mission:** Two transactions run concurrently. Prevent anomalies.

**The Scenario:**
```go
// Transaction 1:
tx1.Execute("MATCH (a {id: 1}) SET a.balance = a.balance - 100")

// Transaction 2 (concurrent):
tx2.Execute("MATCH (a {id: 1}) RETURN a.balance")
```

What should tx2 see?

### Investigation Questions

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

### 🆕 Go 1.25: Deterministic Deadlock Testing

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

### Your Implementation

- [ ] Acquire locks in consistent order (prevent deadlock)
- [ ] Detect deadlocks via timeout or graph analysis
- [ ] Support different isolation levels
- [ ] Show current locks: `SHOW LOCKS;`
- [ ] 🆕 Use `synctest` for deterministic testing

### Test Cases (Go 1.25+)

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

### What You Should Discover

- 2PL prevents anomalies but reduces concurrency
- Deadlocks are inevitable with random lock order
- Isolation levels trade correctness for performance
- 🆕 `synctest` makes concurrency bugs reproducible

---

## Lesson 4.2: MVCC (Multi-Version Concurrency Control) (Week 12)

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

### Investigation Questions

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

### 🆕 Go 1.24: Weak Pointers for Old Versions

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

### 🆕 Go 1.25: Test Version Visibility

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

### Your Implementation

- [ ] Each write creates new version
- [ ] Reads see snapshot at transaction start
- [ ] Detect write-write conflicts
- [ ] GC old versions safely
- [ ] Benchmark: concurrent reads while writing
- [ ] 🆕 Optional: Use weak pointers for memory efficiency

### Benchmark (Go 1.24+)

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

### What You Should Discover

- MVCC trades disk space for concurrency
- Version chains grow quickly under write load
- Snapshot isolation is not serializable (write skew)
- 🆕 Weak pointers allow memory-aware version GC

---

## Phase 4 Complete!

**Congratulations!** You've built a complete graph database from scratch!

**What you've built:**
- ✅ Locking-based transaction support
- ✅ MVCC for high concurrency
- ✅ Deadlock detection
- ✅ Multiple isolation levels

**Key Insights:**
- Concurrency control is a fundamental tradeoff
- Locks are simple but limit parallelism
- MVCC provides better concurrency at cost of complexity
- Testing concurrent code is hard (but `synctest` helps!)

---

## What's Next?

### Bonus Challenges

**Challenge A: Query Optimizer Showdown**

Implement competing optimization strategies:

1. **Exhaustive enumeration** (try all join orders)
2. **Greedy heuristic** (always pick smallest intermediate)
3. **Genetic algorithm** (evolve query plans)

Run 100 random queries. Which wins most often?

**Challenge B: Adaptive Indexing**

Auto-create indexes based on query patterns:
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

**Challenge C: Distributed Execution**

Shard the graph across 3 Go processes:
- Node 1-1M on server A
- Node 1M-2M on server B
- Node 2M-3M on server C

How do you execute cross-shard queries?

**🆕 Challenge D: Experimental GC Evaluation**

Test Go 1.25's new garbage collector:

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

## Reflection

After completing all phases, write a document answering:

1. **What surprised you?** (The thing you didn't expect)
2. **What was harder than expected?** (Where you got stuck)
3. **What optimization made the biggest difference?** (The 10x win)
4. **What would you do differently next time?** (Lessons learned)
5. **🆕 How did Go 1.23+ features help?** (Iterator benefits, `unique` savings, etc.)

These reflections are where the real learning happens.

---

## Final Project: Real-World Application

Build a social network on your database:

```cypher
// Create users
CREATE (alice:Person {name: "Alice", age: 30})
CREATE (bob:Person {name: "Bob", age: 25})

// Create relationships
MATCH (a:Person {name: "Alice"}), (b:Person {name: "Bob"})
CREATE (a)-[:KNOWS]->(b)

// Friend recommendations (friends of friends)
MATCH (me:Person {name: "Alice"})-[:KNOWS]->(friend)-[:KNOWS]->(recommendation)
WHERE NOT (me)-[:KNOWS]->(recommendation)
RETURN recommendation.name, count(*) as mutual_friends
ORDER BY mutual_friends DESC
LIMIT 10
```

**Requirements:**
- Support 10K users
- Handle 100K friendships
- Execute friend recommendations in <100ms
- Support concurrent queries
- Persist data across restarts

---

## Resources

### Papers to Read
- "A Critique of ANSI SQL Isolation Levels" - Berenson et al.
- "Granularity of Locks and Degrees of Consistency" - Gray et al.
- "Serializable Snapshot Isolation" - Cahill et al.

### Reference Implementations
- PostgreSQL's MVCC
- MySQL's InnoDB locking
- CockroachDB's transaction layer

### Go-Specific
- `testing/synctest` for deterministic concurrency tests (Go 1.25+)
- `weak` package for memory management (Go 1.24+)
- Advanced concurrency patterns in Go

---

## Congratulations! 🎉

You've completed the entire learning path! You now understand:

- **Storage:** Pages, buffer pools, WAL
- **Graph structures:** CSR, columnar storage, parallelism
- **Query processing:** Parsing, optimization, joins, execution
- **Transactions:** Locking, MVCC, isolation levels
- **Modern Go:** Iterators, unique handles, weak pointers, synctest

Most importantly, you've learned **how databases really work** by building one from scratch.

**Next steps:**
- Study production databases (DuckDB, PostgreSQL, Kuzu)
- Contribute to open source database projects
- Apply these concepts to your own projects
- Share your learnings with others!

Happy hacking! 🚀
