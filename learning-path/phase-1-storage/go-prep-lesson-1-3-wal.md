# Phase 1 Lesson 1.3: Go Prep - Write-Ahead Log (WAL)

**Prerequisites:** Lessons 1.1-1.2 complete (Pages + Buffer Pool)
**Time:** 3-4 hours Go prep + 20-25 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 1.3

## Overview

A Write-Ahead Log (WAL) ensures durability and crash recovery. Before implementing it, master these Go concepts:
- Sequential file writes and buffering
- `fsync()` for durability guarantees
- Background goroutines with tickers
- **Go 1.23:** Automatic timer cleanup (no more memory leaks!)
- Context-based cancellation
- Idempotent recovery procedures

## Go Concepts for This Lesson

### 1. Sequential File Writes with bufio

WAL requires fast sequential writes. Use `bufio.Writer` for batching!

```go
package main

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "os"
)

type WAL struct {
    file   *os.File
    writer *bufio.Writer
}

func NewWAL(filename string) (*WAL, error) {
    // Open in append mode (O_APPEND ensures atomic appends)
    file, err := os.OpenFile(filename,
        os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    return &WAL{
        file:   file,
        writer: bufio.NewWriter(file),  // 4KB buffer by default
    }, nil
}

func (w *WAL) WriteRecord(data []byte) error {
    // Write length prefix (4 bytes)
    length := uint32(len(data))
    if err := binary.Write(w.writer, binary.LittleEndian, length); err != nil {
        return err
    }

    // Write data
    if _, err := w.writer.Write(data); err != nil {
        return err
    }

    return nil
}

func (w *WAL) Sync() error {
    // Flush buffer to OS
    if err := w.writer.Flush(); err != nil {
        return err
    }

    // Force to disk (CRITICAL for durability!)
    return w.file.Sync()
}

func (w *WAL) Close() error {
    w.Sync()  // Final flush
    return w.file.Close()
}

func main() {
    wal, _ := NewWAL("test.wal")
    defer wal.Close()

    // Write some records
    wal.WriteRecord([]byte("UPDATE page 1"))
    wal.WriteRecord([]byte("UPDATE page 2"))
    wal.WriteRecord([]byte("COMMIT"))

    // Force to disk
    wal.Sync()

    fmt.Println("WAL records written!")
}
```

**Key insight:** `bufio.Writer` batches small writes into larger syscalls (faster!), but you must call `Flush()` before `Sync()`.

### 2. fsync() for Durability

**Without `Sync()`, data can be lost on crash!**

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    file, _ := os.Create("test.dat")
    defer file.Close()

    // Write data
    file.Write([]byte("Important data"))

    // WITHOUT Sync(): Data might only be in OS buffer!
    fmt.Println("Written (but maybe not on disk)")
    time.Sleep(1 * time.Second)
    // If crash here, data lost!

    // WITH Sync(): Guaranteed on disk
    file.Sync()
    fmt.Println("Now definitely on disk")

    os.Remove("test.dat")
}
```

**Benchmark fsync latency:**

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    file, _ := os.Create("bench.dat")
    defer file.Close()
    defer os.Remove("bench.dat")

    data := make([]byte, 4096)

    // Measure write without sync
    start := time.Now()
    for i := 0; i < 100; i++ {
        file.Write(data)
    }
    fmt.Printf("100 writes without sync: %v\n", time.Since(start))

    // Measure write with sync
    start = time.Now()
    for i := 0; i < 100; i++ {
        file.Write(data)
        file.Sync()  // MUCH slower!
    }
    fmt.Printf("100 writes with sync: %v\n", time.Since(start))
}
```

**Typical results:**
- Without sync: ~1ms (buffered in RAM)
- With sync: ~1000ms (waits for disk)

**Solution:** Group commit (covered later).

### 3. Background Goroutines with Timers

**Go 1.23 improvement:** Timers auto-cleanup (no more memory leaks!)

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type BackgroundFlusher struct {
    flushFunc func()
    stopChan  chan struct{}
    wg        sync.WaitGroup
}

func NewBackgroundFlusher(flushFunc func()) *BackgroundFlusher {
    bf := &BackgroundFlusher{
        flushFunc: flushFunc,
        stopChan:  make(chan struct{}),
    }
    bf.start()
    return bf
}

func (bf *BackgroundFlusher) start() {
    bf.wg.Add(1)
    go func() {
        defer bf.wg.Done()

        // Go 1.23: This ticker auto-cleans up!
        ticker := time.NewTicker(100 * time.Millisecond)
        // No need to call ticker.Stop() anymore!

        for {
            select {
            case <-ticker.C:
                bf.flushFunc()
            case <-bf.stopChan:
                return
            }
        }
    }()
}

func (bf *BackgroundFlusher) Stop() {
    close(bf.stopChan)
    bf.wg.Wait()
}

func main() {
    counter := 0
    flusher := NewBackgroundFlusher(func() {
        counter++
        fmt.Printf("Flush #%d\n", counter)
    })

    time.Sleep(1 * time.Second)
    flusher.Stop()
    fmt.Printf("Total flushes: %d\n", counter)
}
```

**Before Go 1.23:** Had to call `ticker.Stop()` to avoid memory leaks.
**Go 1.23+:** Ticker auto-stops when GC'd. Still good practice to stop explicitly for immediate cleanup.

### 4. Context for Graceful Shutdown

Use `context.Context` for proper cancellation:

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type WAL struct {
    file     *os.File
    stopFunc context.CancelFunc
}

func NewWAL() (*WAL, error) {
    file, _ := os.Create("wal.log")

    ctx, cancel := context.WithCancel(context.Background())

    wal := &WAL{
        file:     file,
        stopFunc: cancel,
    }

    go wal.backgroundSync(ctx)

    return wal, nil
}

func (w *WAL) backgroundSync(ctx context.Context) {
    ticker := time.NewTicker(100 * time.Millisecond)

    for {
        select {
        case <-ticker.C:
            w.file.Sync()
            fmt.Println("WAL synced")
        case <-ctx.Done():
            fmt.Println("Stopping background sync")
            return
        }
    }
}

func (w *WAL) Close() error {
    w.stopFunc()  // Cancel context
    time.Sleep(150 * time.Millisecond)  // Wait for final sync
    return w.file.Close()
}

func main() {
    wal, _ := NewWAL()
    time.Sleep(1 * time.Second)
    wal.Close()
}
```

### 5. Group Commit Optimization

Batch multiple writes into a single `fsync()` call!

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "sync"
    "time"
)

type GroupCommitWAL struct {
    file      *os.File
    writer    *bufio.Writer
    mu        sync.Mutex
    pendingCh chan *writeRequest
}

type writeRequest struct {
    data   []byte
    doneCh chan error
}

func NewGroupCommitWAL(filename string) (*GroupCommitWAL, error) {
    file, err := os.OpenFile(filename,
        os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    wal := &GroupCommitWAL{
        file:      file,
        writer:    bufio.NewWriter(file),
        pendingCh: make(chan *writeRequest, 1000),
    }

    go wal.groupCommitLoop()

    return wal, nil
}

func (w *GroupCommitWAL) Write(data []byte) error {
    req := &writeRequest{
        data:   data,
        doneCh: make(chan error, 1),
    }

    w.pendingCh <- req
    return <-req.doneCh  // Wait for commit
}

func (w *GroupCommitWAL) groupCommitLoop() {
    ticker := time.NewTicker(10 * time.Millisecond)

    batch := make([]*writeRequest, 0, 100)

    for {
        select {
        case req := <-w.pendingCh:
            batch = append(batch, req)

            // Collect more writes (non-blocking)
        drainLoop:
            for len(batch) < 100 {
                select {
                case req := <-w.pendingCh:
                    batch = append(batch, req)
                default:
                    break drainLoop
                }
            }

            // Write all at once
            w.mu.Lock()
            for _, req := range batch {
                w.writer.Write(req.data)
            }
            w.writer.Flush()
            err := w.file.Sync()  // Single fsync for all!
            w.mu.Unlock()

            // Notify all waiters
            for _, req := range batch {
                req.doneCh <- err
            }

            batch = batch[:0]  // Reset

        case <-ticker.C:
            // Periodic flush even if no writes
            if len(batch) > 0 {
                // (same commit logic as above)
            }
        }
    }
}

func main() {
    wal, _ := NewGroupCommitWAL("group.wal")

    var wg sync.WaitGroup

    start := time.Now()

    // 100 goroutines writing concurrently
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            wal.Write([]byte(fmt.Sprintf("Write from goroutine %d", id)))
        }(i)
    }

    wg.Wait()

    fmt.Printf("100 writes took: %v\n", time.Since(start))
    fmt.Println("(Much faster than 100 individual fsyncs!)")
}
```

**Performance:** 1 fsync for 100 writes (100x faster!)

### 6. Idempotent Recovery

Crash recovery must be idempotent (safe to replay records).

```go
package main

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "io"
    "os"
)

type LogRecord struct {
    LSN    uint64  // Log Sequence Number
    PageID uint32
    Data   []byte
}

func RecoverFromWAL(filename string) ([]LogRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    records := make([]LogRecord, 0)

    for {
        // Read LSN
        var lsn uint64
        if err := binary.Read(reader, binary.LittleEndian, &lsn); err != nil {
            if err == io.EOF {
                break
            }
            return nil, err
        }

        // Read page ID
        var pageID uint32
        if err := binary.Read(reader, binary.LittleEndian, &pageID); err != nil {
            return nil, err
        }

        // Read data length
        var length uint32
        if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
            return nil, err
        }

        // Read data
        data := make([]byte, length)
        if _, err := io.ReadFull(reader, data); err != nil {
            // Partial record - corruption!
            fmt.Printf("Warning: Partial record at LSN %d\n", lsn)
            break
        }

        records = append(records, LogRecord{
            LSN:    lsn,
            PageID: pageID,
            Data:   data,
        })
    }

    return records, nil
}

func ReplayLog(records []LogRecord) {
    appliedLSNs := make(map[uint64]bool)

    for _, record := range records {
        // Idempotent: Skip if already applied
        if appliedLSNs[record.LSN] {
            fmt.Printf("Skipping already applied LSN %d\n", record.LSN)
            continue
        }

        fmt.Printf("Replaying LSN %d: Update page %d\n", record.LSN, record.PageID)

        // Apply update...

        appliedLSNs[record.LSN] = true
    }
}

func main() {
    records, _ := RecoverFromWAL("test.wal")
    fmt.Printf("Recovered %d records\n", len(records))
    ReplayLog(records)
}
```

## Pre-Implementation Exercises

### Exercise 1: Basic WAL Writer

```go
package main

import (
    "bufio"
    "encoding/binary"
    "os"
)

type WAL struct {
    file   *os.File
    writer *bufio.Writer
    lsn    uint64  // Log Sequence Number
}

func NewWAL(filename string) (*WAL, error) {
    // TODO: Open file in append mode
    // Initialize bufio.Writer
    return nil, nil
}

func (w *WAL) AppendRecord(pageID uint32, data []byte) (uint64, error) {
    // TODO:
    // 1. Increment LSN
    // 2. Write LSN (8 bytes)
    // 3. Write page ID (4 bytes)
    // 4. Write data length (4 bytes)
    // 5. Write data
    return 0, nil
}

func (w *WAL) Sync() error {
    // TODO:
    // 1. Flush buffer
    // 2. Call file.Sync()
    return nil
}

func (w *WAL) Close() error {
    // TODO: Sync and close
    return nil
}

func main() {
    wal, _ := NewWAL("test.wal")
    defer wal.Close()

    lsn1, _ := wal.AppendRecord(1, []byte("Update page 1"))
    lsn2, _ := wal.AppendRecord(2, []byte("Update page 2"))

    println("LSN 1:", lsn1)
    println("LSN 2:", lsn2)

    wal.Sync()
}
```

### Exercise 2: Background Sync with Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type WAL struct {
    // TODO: Add fields
}

func NewWAL() *WAL {
    // TODO: Initialize WAL with background sync goroutine
    return nil
}

func (w *WAL) backgroundSync(ctx context.Context) {
    // TODO:
    // 1. Create ticker for periodic sync (100ms)
    // 2. Loop with select on ticker.C and ctx.Done()
    // 3. Call Sync() on each tick
    // 4. Exit when context cancelled
}

func (w *WAL) Sync() {
    fmt.Println("Syncing WAL...")
    // TODO: Actual sync logic
}

func (w *WAL) Close() {
    // TODO: Cancel context and wait for goroutine
}

func main() {
    wal := NewWAL()
    time.Sleep(1 * time.Second)
    wal.Close()
}
```

### Exercise 3: WAL Reader and Recovery

```go
package main

import (
    "bufio"
    "encoding/binary"
    "io"
    "os"
)

type LogRecord struct {
    LSN    uint64
    PageID uint32
    Data   []byte
}

func ReadWAL(filename string) ([]LogRecord, error) {
    // TODO:
    // 1. Open WAL file
    // 2. Create bufio.Reader
    // 3. Loop reading records until EOF
    // 4. Handle partial records gracefully
    return nil, nil
}

func main() {
    records, _ := ReadWAL("test.wal")
    for _, rec := range records {
        println("LSN:", rec.LSN, "Page:", rec.PageID)
    }
}
```

### Exercise 4: Group Commit Implementation

```go
package main

import (
    "sync"
    "time"
)

type GroupCommitWAL struct {
    mu        sync.Mutex
    pending   [][]byte
    commitCh  chan struct{}
    responses []chan error
}

func NewGroupCommitWAL() *GroupCommitWAL {
    wal := &GroupCommitWAL{
        commitCh: make(chan struct{}, 1),
    }
    go wal.commitLoop()
    return wal
}

func (w *GroupCommitWAL) Write(data []byte) error {
    // TODO:
    // 1. Lock and add data to pending
    // 2. Create response channel
    // 3. Signal commit goroutine
    // 4. Wait for response
    return nil
}

func (w *GroupCommitWAL) commitLoop() {
    ticker := time.NewTicker(10 * time.Millisecond)

    for {
        select {
        case <-w.commitCh:
            // TODO: Commit all pending writes
        case <-ticker.C:
            // TODO: Periodic commit
        }
    }
}

func main() {
    wal := NewGroupCommitWAL()

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            wal.Write([]byte("test"))
        }()
    }
    wg.Wait()
}
```

### Exercise 5: Benchmark fsync Strategies

```go
package main

import (
    "os"
    "testing"
)

func BenchmarkSyncEveryWrite(b *testing.B) {
    file, _ := os.Create("bench1.wal")
    defer file.Close()
    defer os.Remove("bench1.wal")

    data := []byte("test data")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        file.Write(data)
        file.Sync()  // Slow!
    }
}

func BenchmarkSyncBatched(b *testing.B) {
    file, _ := os.Create("bench2.wal")
    defer file.Close()
    defer os.Remove("bench2.wal")

    data := []byte("test data")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        file.Write(data)
        if i%100 == 0 {
            file.Sync()  // Much faster!
        }
    }
}

func BenchmarkBufferedWrites(b *testing.B) {
    file, _ := os.Create("bench3.wal")
    defer file.Close()
    defer os.Remove("bench3.wal")

    writer := bufio.NewWriter(file)
    data := []byte("test data")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        writer.Write(data)
        if i%100 == 0 {
            writer.Flush()
            file.Sync()
        }
    }
}
```

## Performance Benchmarks

### Measure fsync Latency on Your System

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    file, _ := os.Create("fsync_bench.dat")
    defer file.Close()
    defer os.Remove("fsync_bench.dat")

    data := make([]byte, 4096)

    // Measure 100 fsyncs
    start := time.Now()
    for i := 0; i < 100; i++ {
        file.Write(data)
        file.Sync()
    }
    elapsed := time.Since(start)

    fmt.Printf("100 fsyncs took: %v\n", elapsed)
    fmt.Printf("Average per fsync: %v\n", elapsed/100)
}
```

**Typical results:**
- SSD: 0.1-1ms per fsync
- HDD: 5-10ms per fsync (seek time + rotation)
- Cloud storage: 10-100ms (network latency)

### Compare Group Commit

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func sequentialCommit() time.Duration {
    start := time.Now()

    for i := 0; i < 100; i++ {
        // Simulate fsync
        time.Sleep(1 * time.Millisecond)
    }

    return time.Since(start)
}

func groupCommit() time.Duration {
    start := time.Now()

    var wg sync.WaitGroup
    commitCh := make(chan struct{}, 100)

    // Commit goroutine
    go func() {
        for range commitCh {
            // Simulate fsync (once for batch)
            time.Sleep(1 * time.Millisecond)
        }
    }()

    // 100 writes
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            commitCh <- struct{}{}
        }()
    }

    wg.Wait()
    close(commitCh)

    return time.Since(start)
}

func main() {
    seq := sequentialCommit()
    group := groupCommit()

    fmt.Printf("Sequential: %v\n", seq)
    fmt.Printf("Group commit: %v\n", group)
    fmt.Printf("Speedup: %.2fx\n", float64(seq)/float64(group))
}
```

## Common Gotchas to Avoid

### Gotcha 1: Forgetting to Flush Before Sync

```go
// WRONG: Sync doesn't flush bufio!
wal.writer.Write(data)
wal.file.Sync()  // Data still in buffer!

// RIGHT: Flush then sync
wal.writer.Write(data)
wal.writer.Flush()  // Push to OS
wal.file.Sync()     // Push to disk
```

### Gotcha 2: Not Handling Partial Writes

```go
// WRONG: Assumes all bytes written
wal.writer.Write(data)

// RIGHT: Check bytes written
n, err := wal.writer.Write(data)
if err != nil || n != len(data) {
    return fmt.Errorf("partial write: wrote %d/%d bytes", n, len(data))
}
```

### Gotcha 3: Replaying Without Idempotency

```go
// WRONG: Applying twice causes corruption!
func replay(record LogRecord) {
    page.Counter += record.Delta  // BUG if replayed!
}

// RIGHT: Use LSN to track applied records
appliedLSNs := make(map[uint64]bool)

func replay(record LogRecord) {
    if appliedLSNs[record.LSN] {
        return  // Already applied
    }
    page.Counter = record.NewValue  // Idempotent!
    appliedLSNs[record.LSN] = true
}
```

### Gotcha 4: Timer Memory Leaks (Pre-Go 1.23)

```go
// Go 1.22 and earlier: Memory leak!
for {
    ticker := time.NewTicker(1 * time.Second)
    <-ticker.C
    // FORGOT ticker.Stop() - goroutine leaks!
}

// Go 1.23+: Auto-cleanup, but still good practice
for {
    ticker := time.NewTicker(1 * time.Second)
    <-ticker.C
    ticker.Stop()  // Explicit cleanup (optional in 1.23+)
}
```

### Gotcha 5: Not Using O_APPEND

```go
// WRONG: Concurrent appends can interleave!
file, _ := os.OpenFile("wal.log", os.O_WRONLY, 0644)
file.Seek(0, io.SeekEnd)
file.Write(data)  // RACE: Another goroutine might write first!

// RIGHT: O_APPEND makes appends atomic
file, _ := os.OpenFile("wal.log", os.O_APPEND|os.O_WRONLY, 0644)
file.Write(data)  // Atomic append!
```

## Checklist Before Starting Lesson 1.3

- [ ] I understand sequential file writes with `bufio.Writer`
- [ ] I know why `Sync()` is critical for durability
- [ ] I can implement background goroutines with context
- [ ] I understand Go 1.23 timer improvements
- [ ] I know how to implement group commit
- [ ] I can read and parse WAL records
- [ ] I understand idempotent recovery
- [ ] I've measured fsync latency on my system
- [ ] I know the difference between `Flush()` and `Sync()`
- [ ] I can handle partial writes and corrupted records

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 1.3

You'll implement:
- WAL with sequential writes and fsync
- Background sync goroutine with graceful shutdown
- Group commit optimization
- Crash recovery with idempotent replay
- Benchmarks comparing sync strategies
- Chaos testing for crash scenarios

**Time estimate:** 20-25 hours for full implementation

Good luck! 🚀
