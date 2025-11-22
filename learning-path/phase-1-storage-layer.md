# Phase 1: Storage Layer

**Prerequisites:** Complete Phase 0 (Go Fundamentals)

**Time Estimate:** 3 weeks
**Go Version:** 1.23+

In this phase, you'll build the foundational storage layer for your graph database. You'll learn how databases actually store data on disk, manage memory efficiently, and ensure data survives crashes.

---

## Table of Contents

- [Lesson 1.1: The Page Abstraction](#lesson-11-the-page-abstraction-week-1)
- [Lesson 1.2: The Buffer Pool](#lesson-12-the-buffer-pool-week-2)
- [Lesson 1.3: Write-Ahead Log](#lesson-13-write-ahead-log-week-3)

---

## Lesson 1.1: The Page Abstraction (Week 1)

**Your Mission:** Store 1 million integers on disk and read them back.

**Constraints:**
- File size must be exactly 4MB (1M integers × 4 bytes)
- Read them back in under 100ms
- Prove data survives program crash

### Investigation Questions

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

### Your Implementation Must

- [ ] Store exactly N pages with zero wasted space
- [ ] Read any page in O(1) time using page ID
- [ ] Survive `kill -9` without data loss (after sync)
- [ ] Print stats: total pages, bytes per page, storage efficiency

### Verification

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

### What You Should Discover

- OS page cache is doing a lot of work for you
- `mmap` trades simplicity for control
- Alignment matters more than you think

---

## Lesson 1.2: The Buffer Pool (Week 2)

**Your Mission:** You can only keep 10 pages in memory, but need to access 1000 pages randomly.

**The Problem:**
```go
// This will OOM with 1M pages
cache := make(map[uint32]*Page)
for i := 0; i < 1000000; i++ {
    cache[i] = readPageFromDisk(i)
}
```

### Investigation Questions

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

### 🆕 Go 1.24: Weak Reference Cache (Optional Advanced Challenge)

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

### Challenge Requirements

- [ ] Support 1000 concurrent goroutines accessing pages
- [ ] Never exceed memory budget (configurable)
- [ ] Achieve >90% hit rate on Zipfian workload
- [ ] Background flusher writes dirty pages every 1s
- [ ] Zero data races (`go test -race` passes)

### Verification Benchmark (Go 1.24+)

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

### What You Should Discover

- LRU needs a doubly-linked list + hashmap
- Pin counts prevent use-after-free bugs
- Write-back caching is complex but essential
- 🆕 Go 1.24: Swiss Tables make large maps 30% faster
- 🆕 Go 1.24: Weak pointers enable memory-aware caching

---

## Lesson 1.3: Write-Ahead Log (Week 3)

**Your Mission:** Crash your program randomly during writes. Always recover correctly.

**The Scenario:**
```go
// This code must be crash-safe:
db.CreateNode(id: 1, name: "Alice")
db.CreateEdge(from: 1, to: 2)
// CRASH HERE - what's on disk?
```

### Investigation Questions

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

### 🆕 Go 1.23: Timer Improvements

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

### Your Implementation

- [ ] Log every write operation before modifying data pages
- [ ] `fsync()` the log on commit (measure this latency)
- [ ] Implement recovery that's idempotent
- [ ] Add chaos testing: inject crashes randomly

### Chaos Test

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

### What You Should Discover

- `fsync()` is SLOW (5-10ms on SSDs)
- Group commit amortizes fsync cost
- ARIES algorithm has three phases: Analysis, Redo, Undo

---

## Phase 1 Complete!

**What you've built:**
- ✅ Page-based storage system
- ✅ Buffer pool with eviction policies
- ✅ Crash-safe write-ahead logging

**Key Insights:**
- Databases are just files with careful organization
- Memory management is critical for performance
- Durability requires careful synchronization

**Next:** Phase 2 - Graph Structure, where you'll build compressed graph representations that can handle millions of nodes and edges efficiently.

---

## Resources

### Papers to Read
- "The Five-Minute Rule for Trading Memory for Disk Accesses" - Jim Gray
- "ARIES: A Transaction Recovery Method" - Mohan et al.

### Reference Implementations
- SQLite's pager module
- PostgreSQL's buffer manager
- DuckDB's buffer pool

### Go-Specific
- `syscall` package documentation
- Understanding Go's memory management
- Profiling with `pprof`
