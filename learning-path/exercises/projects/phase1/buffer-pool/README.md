# Project 1.2: Buffer Pool Implementation

## Overview
Build a thread-safe buffer pool with LRU eviction, pin/unpin mechanism, and background flushing for high-performance page caching.

**Duration:** 12-15 hours
**Difficulty:** Medium-Hard

## Learning Objectives
- Implement thread-safe concurrent data structures
- Master LRU eviction policies
- Understand pin/unpin reference counting
- Implement background worker patterns
- Handle dirty page management
- Profile and optimize lock contention

## Concepts Covered
- Concurrent data structures with mutexes
- Atomic operations for counters
- LRU cache eviction
- Reference counting (pin/unpin)
- Background goroutines
- Channel-based communication
- Lock-free techniques
- Memory management

## Requirements

### Core Functionality

#### 1. Buffer Pool
- Thread-safe page cache
- Configurable pool size (number of frames)
- LRU eviction policy
- Pin/unpin reference counting
- Dirty page tracking
- Background flusher goroutine

#### 2. Frame Structure
```
Frame = In-memory slot for a page
┌────────────────────────────────┐
│  Frame Metadata                │
│  - Frame ID                    │
│  - Page ID (or -1 if empty)    │
│  - Pin Count (atomic)          │
│  - Dirty Flag                  │
│  - Last Access Time            │
├────────────────────────────────┤
│  Page Data (4KB)               │
│                                │
└────────────────────────────────┘
```

#### 3. Buffer Pool Operations
- **FetchPage(pageID)** - Get page from pool or disk
- **UnpinPage(pageID, dirty)** - Release page reference
- **FlushPage(pageID)** - Write dirty page to disk
- **FlushAll()** - Write all dirty pages
- **NewPage()** - Allocate new page
- **DeletePage(pageID)** - Remove page from pool and disk

#### 4. Eviction Policy
- LRU (Least Recently Used)
- Cannot evict pinned pages
- Flush dirty pages before eviction
- Clock-sweep algorithm (optional optimization)

#### 5. Background Flusher
- Periodically flush dirty pages
- Configurable flush interval
- Graceful shutdown on Close()
- Metrics collection

## Getting Started

```bash
# Initialize module
cd buffer-pool
go mod init bufferpool

# Run tests
go test -v
go test -race
go test -cover

# Run benchmarks
go test -bench=. -benchmem
go test -bench=. -benchtime=10s

# Check for race conditions
go test -race -count=10
```

## API Design

```go
package bufferpool

type BufferPool struct {
    frames      []*Frame
    pageTable   map[PageID]FrameID  // page -> frame mapping
    freeList    []FrameID           // available frames
    replacer    *LRUReplacer        // eviction policy
    diskManager DiskManager         // underlying storage
    mu          sync.RWMutex        // protects pool structures
    flusher     *BackgroundFlusher  // async dirty page writer
}

type Frame struct {
    frameID   FrameID
    pageID    PageID
    data      [PageSize]byte
    pinCount  atomic.Int32
    dirty     atomic.Bool
    mu        sync.RWMutex
}

// Create new buffer pool
func New(diskManager DiskManager, poolSize int) *BufferPool

// Fetch page from pool or disk
func (bp *BufferPool) FetchPage(pageID PageID) (*Frame, error)

// Unpin page and mark dirty if modified
func (bp *BufferPool) UnpinPage(pageID PageID, dirty bool) error

// Flush specific page to disk
func (bp *BufferPool) FlushPage(pageID PageID) error

// Flush all dirty pages to disk
func (bp *BufferPool) FlushAll() error

// Allocate new page
func (bp *BufferPool) NewPage() (PageID, *Frame, error)

// Delete page from pool and disk
func (bp *BufferPool) DeletePage(pageID PageID) error

// Get pool statistics
func (bp *BufferPool) Stats() PoolStats

// Close pool and cleanup
func (bp *BufferPool) Close() error
```

## Test Cases

### Correctness Tests
- **TestFetchPage** - Fetch same page multiple times
- **TestPinUnpin** - Verify pin count tracking
- **TestEviction** - Evict unpinned pages when pool is full
- **TestDirtyPage** - Dirty pages are written before eviction
- **TestNewPage** - Allocate new pages
- **TestDeletePage** - Delete pages from pool
- **TestConcurrentAccess** - Multiple goroutines accessing pool
- **TestPinnedNotEvicted** - Pinned pages cannot be evicted
- **TestFlushAll** - All dirty pages written to disk

### Race Condition Tests
- **TestConcurrentFetch** - Many goroutines fetching same pages
- **TestConcurrentPinUnpin** - Concurrent pin/unpin operations
- **TestConcurrentEviction** - Race-free eviction
- **TestBackgroundFlusher** - No races with flusher goroutine

### Edge Cases
- **TestFullPool** - All frames pinned, cannot evict
- **TestEmptyPool** - Operations on empty pool
- **TestInvalidPageID** - Handle invalid page IDs
- **TestDoubleFetch** - Fetch already-pinned page
- **TestUnpinNotFetched** - Unpin page that wasn't fetched

### Integration Tests
- **TestWithRealDisk** - Use actual disk manager
- **TestRecovery** - Flush and reload pages
- **TestLargeWorkload** - 10K pages, 100K operations

## Benchmarks

```go
BenchmarkFetchPage          - Fetch cached page
BenchmarkFetchPageCold      - Fetch uncached page
BenchmarkPinUnpin           - Pin/unpin overhead
BenchmarkConcurrentFetch    - Concurrent access scalability
BenchmarkEviction           - Eviction performance
BenchmarkFlush              - Flush performance
BenchmarkNewPage            - Page allocation speed
```

## Implementation Hints

### Frame Structure
```go
type FrameID int32
type PageID int64

type Frame struct {
    frameID  FrameID
    pageID   PageID
    data     [PageSize]byte
    pinCount atomic.Int32
    dirty    atomic.Bool
    mu       sync.RWMutex
}

func (f *Frame) Pin() {
    f.pinCount.Add(1)
}

func (f *Frame) Unpin() {
    if f.pinCount.Add(-1) < 0 {
        panic("unpin of unpinned frame")
    }
}

func (f *Frame) IsPinned() bool {
    return f.pinCount.Load() > 0
}
```

### LRU Replacer
```go
type LRUReplacer struct {
    capacity  int
    frames    map[FrameID]*list.Element
    lruList   *list.List
    mu        sync.Mutex
}

// Record frame access
func (r *LRUReplacer) RecordAccess(frameID FrameID)

// Get victim frame for eviction
func (r *LRUReplacer) Victim() (FrameID, bool)

// Remove frame from replacer
func (r *LRUReplacer) Remove(frameID FrameID)

// Check size
func (r *LRUReplacer) Size() int
```

### Background Flusher
```go
type BackgroundFlusher struct {
    pool       *BufferPool
    interval   time.Duration
    stopCh     chan struct{}
    doneCh     chan struct{}
}

func (f *BackgroundFlusher) Start() {
    go func() {
        ticker := time.NewTicker(f.interval)
        defer ticker.Stop()
        defer close(f.doneCh)

        for {
            select {
            case <-ticker.C:
                f.flushDirtyPages()
            case <-f.stopCh:
                f.flushDirtyPages() // final flush
                return
            }
        }
    }()
}

func (f *BackgroundFlusher) Stop() {
    close(f.stopCh)
    <-f.doneCh  // wait for shutdown
}
```

### Fetch Page Algorithm
```go
func (bp *BufferPool) FetchPage(pageID PageID) (*Frame, error) {
    bp.mu.Lock()

    // Check if page already in pool
    if frameID, found := bp.pageTable[pageID]; found {
        frame := bp.frames[frameID]
        bp.mu.Unlock()

        frame.Pin()
        bp.replacer.Remove(frameID)
        return frame, nil
    }

    // Find victim frame
    var frameID FrameID
    if len(bp.freeList) > 0 {
        frameID = bp.freeList[len(bp.freeList)-1]
        bp.freeList = bp.freeList[:len(bp.freeList)-1]
    } else {
        victimID, found := bp.replacer.Victim()
        if !found {
            bp.mu.Unlock()
            return nil, ErrNoVictimFrame
        }
        frameID = victimID

        // Evict old page
        oldFrame := bp.frames[frameID]
        if oldFrame.dirty.Load() {
            // Flush before eviction
            bp.diskManager.WritePage(oldFrame.pageID, oldFrame.data[:])
        }
        delete(bp.pageTable, oldFrame.pageID)
    }

    frame := bp.frames[frameID]
    frame.pageID = pageID
    bp.pageTable[pageID] = frameID
    bp.mu.Unlock()

    // Load page from disk
    if err := bp.diskManager.ReadPage(pageID, frame.data[:]); err != nil {
        return nil, err
    }

    frame.dirty.Store(false)
    frame.Pin()

    return frame, nil
}
```

## Performance Goals

- Cached page fetch: < 500ns
- Uncached page fetch: < 50µs (SSD)
- Pin/Unpin: < 100ns
- Eviction: < 10µs
- Background flush: < 100ms for 1000 dirty pages
- Lock contention: < 10% under concurrent load
- Throughput: > 100K ops/sec with 8 threads

## Stretch Goals

### 1. Clock-Sweep Eviction
- Implement Clock (Second-Chance) algorithm
- Compare to LRU performance
- Measure CPU usage difference

### 2. Lock-Free Optimizations
- Use atomic operations where possible
- Read-write locks for better concurrency
- Lock-free page table (sync.Map)

### 3. Adaptive Flushing
- Monitor dirty page ratio
- Adjust flush frequency dynamically
- Burst detection and handling

### 4. Prefetching
- Sequential scan detection
- Read-ahead pages
- Measure prefetch effectiveness

### 5. Different Eviction Policies
- LRU-K (track K recent accesses)
- ARC (Adaptive Replacement Cache)
- LIRS (Low Inter-reference Recency Set)
- Benchmark comparisons

## Common Pitfalls

1. **Deadlocks**
   - Always acquire locks in consistent order
   - Avoid holding pool lock during I/O
   - Use defer for unlocking

2. **Pin Count Bugs**
   - Every Fetch must have matching Unpin
   - Use defer to ensure unpin happens
   - Check for negative pin counts

3. **Race Conditions**
   - Protect all shared state
   - Use atomic operations for counters
   - Test with -race flag extensively

4. **Eviction Failures**
   - Handle case where all pages are pinned
   - Return appropriate error
   - Don't panic on full pool

5. **Memory Leaks**
   - Properly stop background goroutines
   - Clear page table entries
   - Release resources in Close()

## Debugging Tips

```bash
# Race detection
go test -race -run=TestConcurrent

# CPU profiling
go test -bench=BenchmarkConcurrentFetch -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof

# Mutex contention profiling
go test -bench=. -mutexprofile=mutex.prof
go tool pprof mutex.prof

# Trace analysis
go test -trace=trace.out
go tool trace trace.out
```

## Validation Checklist

Your implementation should:
- [ ] Pass all unit tests
- [ ] Pass race detector with no warnings
- [ ] Achieve >85% code coverage
- [ ] Meet performance goals
- [ ] Handle all edge cases correctly
- [ ] Properly manage pin counts
- [ ] Flush dirty pages before eviction
- [ ] Support concurrent access
- [ ] Clean shutdown of background workers
- [ ] No goroutine leaks (check with pprof)

## Learning Outcomes

After completing this project, you will understand:
- Thread-safe concurrent data structure design
- LRU cache implementation and variants
- Pin/unpin reference counting patterns
- Background worker goroutine patterns
- Lock granularity and contention optimization
- Atomic operations for lock-free code
- Profiling concurrent Go programs
- Common concurrency pitfalls and solutions

## Time Estimate
- Core implementation: 8-10 hours
- Testing and race condition fixes: 2-3 hours
- Benchmarking and optimization: 2-3 hours
- Stretch goals: 6-8 hours (optional)

## Next Steps
After completing this project, move on to **Project 1.3: Write-Ahead Log** which implements crash-safe transaction logging that works with your buffer pool.
