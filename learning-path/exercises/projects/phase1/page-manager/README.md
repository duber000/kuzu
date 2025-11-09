# Project 1.1: Simple Page Manager

## Overview
Implement a page manager that allocates, reads, and writes fixed-size pages to disk with a simple LRU cache.

**Duration:** 8-10 hours
**Difficulty:** Medium

## Learning Objectives
- Understand page-based storage
- Implement file I/O at page granularity
- Build a simple caching layer
- Track free pages with bitmaps
- Handle data integrity

## Concepts Covered
- File I/O with `os` package
- Byte manipulation
- LRU cache implementation
- Bitmap data structure
- Page layout and alignment
- Cache hit rate measurement

## Requirements

### Core Functionality

#### 1. Page Manager
- Fixed page size: 4KB (4096 bytes)
- Allocate new pages
- Read pages from disk
- Write pages to disk
- Free pages (mark as available)
- Track free pages with bitmap

#### 2. Page Structure
```
┌────────────────────────────────┐
│  Page Header (64 bytes)        │
├────────────────────────────────┤
│  - Page ID (8 bytes)           │
│  - Checksum (8 bytes)          │
│  - Free Space (8 bytes)        │
│  - Flags (8 bytes)             │
│  - LSN (8 bytes)               │
│  - Reserved (24 bytes)         │
├────────────────────────────────┤
│  Page Data (4032 bytes)        │
│                                │
│  ... user data ...             │
│                                │
└────────────────────────────────┘
```

#### 3. LRU Cache
- Configurable cache size (number of pages)
- Least Recently Used eviction
- Pin/unpin mechanism (prevent eviction of in-use pages)
- Track cache hit rate
- Write-back policy (flush dirty pages)

#### 4. Free Page Bitmap
- Track which pages are allocated/free
- Efficient free page lookup
- Persist bitmap to disk
- Support for growing file

### File Layout
```
┌─────────────────────────────────┐
│  File Header (1 page)           │
│  - Magic number                 │
│  - Page size                    │
│  - Total pages                  │
│  - Free page count              │
├─────────────────────────────────┤
│  Free Page Bitmap (N pages)     │
│  - 1 bit per page               │
├─────────────────────────────────┤
│  Data Pages                     │
│  - Page 0                       │
│  - Page 1                       │
│  - ...                          │
└─────────────────────────────────┘
```

## Getting Started

```bash
# Initialize module
go mod init pagemanager

# Run tests
go test -v
go test -race
go test -cover

# Run benchmarks
go test -bench=. -benchmem
```

## API Design

```go
type PageManager struct {
    file      *os.File
    pageSize  int
    cache     *LRUCache
    freeBitmap *Bitmap
}

// Create new page manager
func New(filename string, pageSize int, cacheSize int) (*PageManager, error)

// Allocate a new page, returns page ID
func (pm *PageManager) AllocatePage() (PageID, error)

// Free a page (mark as available)
func (pm *PageManager) FreePage(pageID PageID) error

// Read page from disk (may come from cache)
func (pm *PageManager) ReadPage(pageID PageID) (*Page, error)

// Write page to disk (may be cached)
func (pm *PageManager) WritePage(page *Page) error

// Flush all dirty pages to disk
func (pm *PageManager) Flush() error

// Get cache statistics
func (pm *PageManager) CacheStats() CacheStats

// Close and cleanup
func (pm *PageManager) Close() error
```

## Test Cases

### Basic Operations
- Allocate and write 1000 pages
- Read pages in random order
- Verify data integrity
- Free pages and reallocate
- Test page overflow (beyond file size)

### Cache Tests
- Measure cache hit rate
- Test LRU eviction
- Test pin/unpin preventing eviction
- Test dirty page flushing
- Fill cache beyond capacity

### Bitmap Tests
- Find free page efficiently
- Handle full bitmap (all allocated)
- Test bitmap persistence
- Verify bitmap correctness

### Integrity Tests
- Checksum validation
- Corrupted page detection
- Partial write handling
- Verify data after restart

## Benchmarks

```
BenchmarkAllocatePage       - Page allocation speed
BenchmarkReadPage           - Read with cache hits
BenchmarkReadPageCold       - Read with cache misses
BenchmarkWritePage          - Write performance
BenchmarkSequentialRead     - Sequential access pattern
BenchmarkRandomRead         - Random access pattern
```

## Implementation Hints

### Page Structure
```go
type PageID uint64

type Page struct {
    ID       PageID
    Data     [4032]byte
    Dirty    bool
    Pinned   bool
    Checksum uint64
}

func (p *Page) ComputeChecksum() uint64 {
    // Use hash/crc64 or similar
}
```

### LRU Cache
```go
type LRUCache struct {
    capacity int
    pages    map[PageID]*list.Element
    lru      *list.List
    mu       sync.RWMutex
}

type cacheEntry struct {
    pageID PageID
    page   *Page
}
```

### Bitmap
```go
type Bitmap struct {
    bits []byte
    size int
}

func (b *Bitmap) Set(n int) {
    b.bits[n/8] |= (1 << (n % 8))
}

func (b *Bitmap) Clear(n int) {
    b.bits[n/8] &= ^(1 << (n % 8))
}

func (b *Bitmap) Test(n int) bool {
    return (b.bits[n/8] & (1 << (n % 8))) != 0
}

func (b *Bitmap) FindFirstZero() int {
    // Efficiently find first free page
}
```

## Performance Goals

- Page allocation: < 1µs
- Cached read: < 100ns
- Uncached read: < 50µs (SSD)
- Write (cached): < 200ns
- Cache hit rate: > 90% for sequential access
- Cache hit rate: > 70% for random access

## Stretch Goals

### 1. Page Compression
- Compress pages before writing to disk
- Decompress on read
- Variable-length compressed pages
- Benchmark compression ratio vs speed

### 2. Wear Leveling
- Distribute writes across file
- Track write counts per page
- Implement wear-aware allocation

### 3. Crash Recovery
- WAL integration
- Redo logging
- Atomic page writes
- Verify recovery correctness

### 4. Advanced Caching
- 2Q or ARC eviction policy
- Prefetching (read-ahead)
- Adaptive cache sizing
- Per-access-pattern optimization

## Common Pitfalls

1. **Alignment Issues**
   - Ensure pages are properly aligned
   - Use `syscall.Ftruncate` for file sizing
   - Be careful with offsets

2. **Concurrency**
   - Protect shared structures with mutexes
   - Avoid holding locks during I/O
   - Test with `-race` flag

3. **Resource Leaks**
   - Always close files
   - Clean up goroutines
   - Free memory in cache

4. **Data Corruption**
   - Validate checksums
   - Use atomic writes
   - Flush before closing

## Debugging Tips

```bash
# Check file structure
hexdump -C data.db | head -20

# Profile I/O operations
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Check for races
go test -race

# Memory profiling
go test -memprofile=mem.prof
go tool pprof mem.prof
```

## Validation

Your implementation should:
- [ ] Pass all unit tests
- [ ] Pass race detector
- [ ] Achieve >80% code coverage
- [ ] Meet performance goals
- [ ] Handle edge cases (full file, no free pages, etc.)
- [ ] Correctly compute and verify checksums
- [ ] Persist data across restarts

## Learning Outcomes

After completing this project, you will understand:
- How databases organize data in pages
- The importance of caching for performance
- LRU cache implementation
- Bitmap usage for tracking resources
- Data integrity with checksums
- File I/O optimization techniques

## Time Estimate
- Core implementation: 6-8 hours
- Testing: 2-3 hours
- Benchmarking and optimization: 2-3 hours
- Stretch goals: 4-6 hours (optional)

## Next Steps
After completing this project, move on to **Project 1.2: Buffer Pool Implementation** which builds on these concepts with more sophisticated concurrency and eviction policies.
