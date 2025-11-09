# Phase 1 Lesson 1.2: Go Prep - Buffer Pool

**Prerequisites:** Lesson 1.1 complete (Page Abstraction)
**Time:** 3-4 hours Go prep + 20-25 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 1.2

## Overview

A buffer pool is a cache of pages in memory. Before implementing it, master these Go concepts:
- LRU cache implementation with map + doubly-linked list
- Thread-safety with `sync.Mutex` vs `sync.RWMutex`
- Lock-free counters with `atomic` package
- **Go 1.24:** Weak pointers for cache optimization
- **Go 1.24:** Swiss Tables (30-35% faster map operations!)
- Race condition detection

## Go Concepts for This Lesson

### 1. LRU Cache with container/list

The standard library provides a doubly-linked list perfect for LRU!

```go
package main

import (
    "container/list"
    "fmt"
)

type LRUCache struct {
    capacity int
    cache    map[int]*list.Element
    lru      *list.List
}

type entry struct {
    key   int
    value string
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[int]*list.Element),
        lru:      list.New(),
    }
}

func (c *LRUCache) Get(key int) (string, bool) {
    if elem, found := c.cache[key]; found {
        // Move to front (most recently used)
        c.lru.MoveToFront(elem)
        return elem.Value.(*entry).value, true
    }
    return "", false
}

func (c *LRUCache) Put(key int, value string) {
    // Already exists? Update and move to front
    if elem, found := c.cache[key]; found {
        c.lru.MoveToFront(elem)
        elem.Value.(*entry).value = value
        return
    }

    // Evict if at capacity
    if c.lru.Len() >= c.capacity {
        oldest := c.lru.Back()
        if oldest != nil {
            c.lru.Remove(oldest)
            delete(c.cache, oldest.Value.(*entry).key)
        }
    }

    // Add new entry
    elem := c.lru.PushFront(&entry{key, value})
    c.cache[key] = elem
}

func main() {
    cache := NewLRUCache(3)

    cache.Put(1, "one")
    cache.Put(2, "two")
    cache.Put(3, "three")

    fmt.Println(cache.Get(1))  // "one", true - now most recent

    cache.Put(4, "four")  // Evicts 2 (least recently used)

    fmt.Println(cache.Get(2))  // "", false - evicted!
    fmt.Println(cache.Get(1))  // "one", true - still there
}
```

**Key insight:** Map gives O(1) lookup, list gives O(1) move-to-front and eviction.

### 2. Thread-Safety: sync.Mutex vs sync.RWMutex

**Use `RWMutex` when reads far outnumber writes!**

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type SafeCache struct {
    mu    sync.RWMutex
    cache map[int]string
}

func (c *SafeCache) Get(key int) (string, bool) {
    c.mu.RLock()         // Multiple readers allowed
    defer c.mu.RUnlock()

    value, found := c.cache[key]
    return value, found
}

func (c *SafeCache) Put(key int, value string) {
    c.mu.Lock()          // Exclusive write lock
    defer c.mu.Unlock()

    c.cache[key] = value
}

func main() {
    cache := &SafeCache{cache: make(map[int]string)}

    // Simulate read-heavy workload
    var wg sync.WaitGroup

    // 10 readers
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                cache.Get(j % 100)
            }
        }()
    }

    // 1 writer
    wg.Add(1)
    go func() {
        defer wg.Done()
        for j := 0; j < 100; j++ {
            cache.Put(j, fmt.Sprintf("value-%d", j))
            time.Sleep(1 * time.Millisecond)
        }
    }()

    wg.Wait()
}
```

**Performance comparison:**

```go
// Regular Mutex: ALL operations exclusive (slow for reads!)
type SlowCache struct {
    mu    sync.Mutex
    cache map[int]string
}

// RWMutex: Reads can run in parallel (fast!)
type FastCache struct {
    mu    sync.RWMutex
    cache map[int]string
}
```

**Rule of thumb:** If 80%+ operations are reads, use `RWMutex`.

### 3. Lock-Free Counters with atomic

Pin counts need to be lock-free for performance!

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type Page struct {
    ID   uint32
    Data [4096]byte
    pins atomic.Int32  // Thread-safe without locks!
}

func (p *Page) Pin() {
    p.pins.Add(1)
}

func (p *Page) Unpin() {
    p.pins.Add(-1)
}

func (p *Page) PinCount() int32 {
    return p.pins.Load()
}

func main() {
    page := &Page{ID: 1}

    var wg sync.WaitGroup

    // 100 goroutines pinning/unpinning concurrently
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                page.Pin()
                // Do work...
                page.Unpin()
            }
        }()
    }

    wg.Wait()

    fmt.Println("Final pin count:", page.PinCount())  // Should be 0
}
```

**Why atomic instead of mutex?**

```go
// With mutex: ~30ns per operation
type SlowPage struct {
    mu   sync.Mutex
    pins int32
}

func (p *SlowPage) Pin() {
    p.mu.Lock()
    p.pins++
    p.mu.Unlock()
}

// With atomic: ~5ns per operation (6x faster!)
type FastPage struct {
    pins atomic.Int32
}

func (p *FastPage) Pin() {
    p.pins.Add(1)  // Single CPU instruction!
}
```

### 4. Go 1.24 Feature: weak.Pointer

**New in Go 1.24:** Weak pointers don't prevent garbage collection!

```go
package main

import (
    "fmt"
    "runtime"
    "weak"
)

type CachedData struct {
    Value string
}

type WeakCache struct {
    cache map[int]weak.Pointer[*CachedData]
}

func NewWeakCache() *WeakCache {
    return &WeakCache{
        cache: make(map[int]weak.Pointer[*CachedData]),
    }
}

func (c *WeakCache) Put(key int, data *CachedData) {
    c.cache[key] = weak.Make(data)
}

func (c *WeakCache) Get(key int) *CachedData {
    if wp, found := c.cache[key]; found {
        // Value() returns nil if GC collected the data
        return wp.Value()
    }
    return nil
}

func main() {
    cache := NewWeakCache()

    // Create data
    data := &CachedData{Value: "important"}
    cache.Put(1, data)

    fmt.Println("Before GC:", cache.Get(1))  // Found

    // Remove strong reference
    data = nil

    // Force GC
    runtime.GC()

    fmt.Println("After GC:", cache.Get(1))  // Likely nil (GC'd)
}
```

**Use case for buffer pools:** Cache clean (non-dirty) pages with weak pointers. If memory pressure occurs, GC can reclaim them without manual eviction!

### 5. Go 1.24: Swiss Tables (30-35% Speedup!)

**Go 1.24 automatically uses Swiss Tables for maps!** No code changes needed.

```go
package main

import (
    "fmt"
    "testing"
)

// Before Go 1.24: ~100ns per lookup
// Go 1.24+: ~65ns per lookup (35% faster!)

func BenchmarkMapLookup(b *testing.B) {
    m := make(map[uint32]int, 1000000)
    for i := uint32(0); i < 1000000; i++ {
        m[i] = int(i)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m[uint32(i%1000000)]
    }
}

func main() {
    result := testing.Benchmark(BenchmarkMapLookup)
    fmt.Printf("Map lookup: %v\n", result.NsPerOp())
}
```

**Key optimization:** Pre-size maps!

```go
// Bad: Causes multiple rehashes
m := make(map[int]string)
for i := 0; i < 1000000; i++ {
    m[i] = "value"
}

// Good: Pre-sized, no rehashing (35% faster in Go 1.24!)
m := make(map[int]string, 1000000)
for i := 0; i < 1000000; i++ {
    m[i] = "value"
}
```

### 6. Race Detection

**CRITICAL:** Always test concurrent code with `-race` flag!

```go
package main

import (
    "sync"
    "testing"
)

type UnsafeCounter struct {
    count int
}

func (c *UnsafeCounter) Increment() {
    c.count++  // RACE CONDITION!
}

func TestRaceCondition(t *testing.T) {
    counter := &UnsafeCounter{}
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
}
```

Run with:
```bash
go test -race
```

Output:
```
==================
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  counter.Increment()

Previous write at 0x... by goroutine 6:
  counter.Increment()
==================
```

**Fix:**

```go
type SafeCounter struct {
    count atomic.Int32
}

func (c *SafeCounter) Increment() {
    c.count.Add(1)  // Thread-safe!
}
```

## Pre-Implementation Exercises

Complete these BEFORE starting the main lesson:

### Exercise 1: Thread-Safe LRU Cache

```go
package main

import (
    "container/list"
    "sync"
)

type LRUCache struct {
    capacity int
    mu       sync.RWMutex
    cache    map[uint32]*list.Element
    lru      *list.List
}

type entry struct {
    key   uint32
    value *Page
}

type Page struct {
    ID    uint32
    Data  [4096]byte
    Dirty bool
}

func NewLRUCache(capacity int) *LRUCache {
    // TODO: Initialize cache
    return nil
}

func (c *LRUCache) Get(pageID uint32) (*Page, bool) {
    // TODO: Thread-safe get with RLock
    // Move to front if found
    return nil, false
}

func (c *LRUCache) Put(pageID uint32, page *Page) *Page {
    // TODO: Thread-safe put with Lock
    // Return evicted page (if any)
    return nil
}

func (c *LRUCache) Remove(pageID uint32) {
    // TODO: Remove specific page
}
```

**Test it:**

```go
func main() {
    cache := NewLRUCache(3)

    p1 := &Page{ID: 1}
    p2 := &Page{ID: 2}
    p3 := &Page{ID: 3}
    p4 := &Page{ID: 4}

    cache.Put(1, p1)
    cache.Put(2, p2)
    cache.Put(3, p3)

    // Access p1 (moves to front)
    cache.Get(1)

    // This should evict p2 (least recently used)
    evicted := cache.Put(4, p4)

    println("Evicted page ID:", evicted.ID)  // Should be 2
}
```

### Exercise 2: Pin Count Management

```go
package main

import (
    "sync"
    "sync/atomic"
    "testing"
)

type Page struct {
    ID    uint32
    Data  [4096]byte
    pins  atomic.Int32
    Dirty bool
}

func (p *Page) Pin() {
    // TODO: Increment pin count atomically
}

func (p *Page) Unpin() {
    // TODO: Decrement pin count atomically
}

func (p *Page) PinCount() int32 {
    // TODO: Load pin count atomically
    return 0
}

func (p *Page) CanEvict() bool {
    // TODO: Return true if pin count is 0
    return false
}

func TestPinCounts(t *testing.T) {
    page := &Page{ID: 1}
    var wg sync.WaitGroup

    // 100 goroutines each pin/unpin 1000 times
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                page.Pin()
                // Simulate work...
                page.Unpin()
            }
        }()
    }

    wg.Wait()

    if page.PinCount() != 0 {
        t.Errorf("Expected 0 pins, got %d", page.PinCount())
    }
}
```

Run with race detector:
```bash
go test -race
```

### Exercise 3: Buffer Pool with Eviction Policy

```go
package main

import (
    "container/list"
    "errors"
    "sync"
    "sync/atomic"
)

var ErrNoEvictablePages = errors.New("no pages can be evicted (all pinned)")

type BufferPool struct {
    mu       sync.RWMutex
    capacity int
    cache    map[uint32]*list.Element
    lru      *list.List
}

type Page struct {
    ID    uint32
    Data  [4096]byte
    pins  atomic.Int32
    Dirty bool
}

type entry struct {
    pageID uint32
    page   *Page
}

func NewBufferPool(capacity int) *BufferPool {
    // TODO: Initialize buffer pool
    return nil
}

func (bp *BufferPool) FetchPage(pageID uint32) (*Page, error) {
    // TODO:
    // 1. Check if page in cache (RLock)
    // 2. If found, pin it and move to front
    // 3. If not found, check if cache is full
    // 4. If full, evict unpinned page (Lock)
    // 5. Load page from disk (simulated)
    // 6. Add to cache and pin it
    return nil, nil
}

func (bp *BufferPool) UnpinPage(pageID uint32, isDirty bool) error {
    // TODO:
    // 1. Find page in cache (RLock)
    // 2. Unpin it
    // 3. Update dirty flag
    return nil
}

func (bp *BufferPool) FlushPage(pageID uint32) error {
    // TODO:
    // 1. Find page in cache (RLock)
    // 2. If dirty, write to disk (simulated)
    // 3. Clear dirty flag
    return nil
}

func (bp *BufferPool) evictPage() (*Page, error) {
    // TODO: Find an unpinned page to evict (start from LRU)
    // Return ErrNoEvictablePages if all pages are pinned
    return nil, nil
}
```

### Exercise 4: Benchmark Buffer Pool Operations

```go
package main

import (
    "testing"
)

func BenchmarkBufferPoolFetch(b *testing.B) {
    bp := NewBufferPool(1000)

    // Pre-populate
    for i := uint32(0); i < 1000; i++ {
        bp.FetchPage(i)
        bp.UnpinPage(i, false)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pageID := uint32(i % 1000)
        page, _ := bp.FetchPage(pageID)
        bp.UnpinPage(page.ID, false)
    }
}

func BenchmarkConcurrentAccess(b *testing.B) {
    bp := NewBufferPool(1000)

    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            pageID := uint32(i % 1000)
            page, _ := bp.FetchPage(pageID)
            bp.UnpinPage(page.ID, false)
            i++
        }
    })
}
```

Run and analyze:
```bash
go test -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### Exercise 5: Test with Race Detector

```go
package main

import (
    "sync"
    "testing"
)

func TestConcurrentBufferPool(t *testing.T) {
    bp := NewBufferPool(100)
    var wg sync.WaitGroup

    // Multiple goroutines accessing same pages
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                pageID := uint32(j % 50)
                page, err := bp.FetchPage(pageID)
                if err != nil {
                    t.Errorf("Failed to fetch page: %v", err)
                    return
                }

                // Simulate work
                page.Data[0] = byte(id)

                bp.UnpinPage(pageID, true)
            }
        }(i)
    }

    wg.Wait()
}
```

**CRITICAL:** Run this test with:
```bash
go test -race -count=100
```

If you see any race warnings, fix them before proceeding!

## Performance Benchmarks

### Benchmark 1: LRU vs No Eviction

```go
func BenchmarkLRUEviction(b *testing.B) {
    cache := NewLRUCache(100)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        pageID := uint32(i % 200)  // 50% cache misses
        if page, found := cache.Get(pageID); !found {
            page = &Page{ID: pageID}
            cache.Put(pageID, page)
        }
    }
}
```

### Benchmark 2: Mutex vs RWMutex

```go
func BenchmarkMutexReads(b *testing.B) {
    type Cache struct {
        mu    sync.Mutex
        pages map[uint32]*Page
    }

    cache := &Cache{pages: make(map[uint32]*Page)}
    for i := uint32(0); i < 100; i++ {
        cache.pages[i] = &Page{ID: i}
    }

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            cache.mu.Lock()
            _ = cache.pages[50]
            cache.mu.Unlock()
        }
    })
}

func BenchmarkRWMutexReads(b *testing.B) {
    type Cache struct {
        mu    sync.RWMutex
        pages map[uint32]*Page
    }

    cache := &Cache{pages: make(map[uint32]*Page)}
    for i := uint32(0); i < 100; i++ {
        cache.pages[i] = &Page{ID: i}
    }

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            cache.mu.RLock()
            _ = cache.pages[50]
            cache.mu.RUnlock()
        }
    })
}
```

**Expected results:** RWMutex should be 5-10x faster for read-heavy workloads!

### Benchmark 3: atomic vs Mutex for Counters

```go
func BenchmarkAtomicIncrement(b *testing.B) {
    var counter atomic.Int32

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter.Add(1)
        }
    })
}

func BenchmarkMutexIncrement(b *testing.B) {
    var mu sync.Mutex
    var counter int32

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mu.Lock()
            counter++
            mu.Unlock()
        }
    })
}
```

**Expected results:** Atomic should be 3-6x faster!

## Common Gotchas to Avoid

### Gotcha 1: Forgetting to Unpin Pages

```go
// WRONG: Page stays pinned forever!
page, _ := bp.FetchPage(123)
// Forgot to unpin...

// Eventually all pages are pinned, buffer pool deadlocks!
```

**Fix:** Use defer pattern:

```go
// RIGHT: Always unpin
page, _ := bp.FetchPage(123)
defer bp.UnpinPage(page.ID, false)

// Use page...
```

### Gotcha 2: Evicting Pinned Pages

```go
// WRONG: Might evict a page that's still in use!
func (bp *BufferPool) evict() *Page {
    oldest := bp.lru.Back()
    page := oldest.Value.(*Page)
    bp.lru.Remove(oldest)
    return page  // BUG: Didn't check pin count!
}

// RIGHT: Only evict unpinned pages
func (bp *BufferPool) evict() *Page {
    for elem := bp.lru.Back(); elem != nil; elem = elem.Prev() {
        page := elem.Value.(*Page)
        if page.PinCount() == 0 {
            bp.lru.Remove(elem)
            return page
        }
    }
    return nil  // All pages pinned!
}
```

### Gotcha 3: Lock Ordering Issues (Deadlock!)

```go
// WRONG: Can cause deadlock!
func (bp *BufferPool) transferPage(fromID, toID uint32) {
    bp.mu.Lock()
    // ... do some work ...
    bp.mu.Lock()  // DEADLOCK: Already locked!
    // ...
    bp.mu.Unlock()
    bp.mu.Unlock()
}

// RIGHT: Single lock acquisition
func (bp *BufferPool) transferPage(fromID, toID uint32) {
    bp.mu.Lock()
    defer bp.mu.Unlock()
    // ... all work here ...
}
```

### Gotcha 4: Race on Dirty Flag

```go
// WRONG: Dirty flag is not atomic!
type Page struct {
    ID    uint32
    Dirty bool  // RACE CONDITION!
    pins  atomic.Int32
}

// Multiple goroutines doing this = race!
page.Dirty = true

// RIGHT: Use atomic for dirty flag too
type Page struct {
    ID    uint32
    dirty atomic.Bool  // Thread-safe!
    pins  atomic.Int32
}

page.dirty.Store(true)
```

### Gotcha 5: Not Pre-Sizing Maps

```go
// WRONG: Multiple rehashes in Go 1.24
cache := make(map[uint32]*Page)
for i := uint32(0); i < 10000; i++ {
    cache[i] = &Page{ID: i}
}

// RIGHT: Pre-size for 35% speedup!
cache := make(map[uint32]*Page, 10000)
for i := uint32(0); i < 10000; i++ {
    cache[i] = &Page{ID: i}
}
```

## Checklist Before Starting Lesson 1.2

- [ ] I can implement LRU cache with map + doubly-linked list
- [ ] I understand when to use `RWMutex` vs `Mutex`
- [ ] I can use `atomic.Int32` for lock-free counters
- [ ] I've experimented with `weak.Pointer` (Go 1.24+)
- [ ] I know that Go 1.24 uses Swiss Tables automatically
- [ ] I always run `go test -race` on concurrent code
- [ ] I understand pin count semantics
- [ ] I know how to prevent evicting pinned pages
- [ ] I pre-size maps for better performance
- [ ] I can profile and benchmark buffer pool operations

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 1.2

You'll implement:
- Thread-safe buffer pool with LRU eviction
- Pin/unpin mechanism for page safety
- Flush dirty pages to disk
- Concurrent access benchmarks
- Optional: Weak pointer optimization for clean pages

**Time estimate:** 20-25 hours for full implementation

Good luck! 🚀
