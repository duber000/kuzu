# Phase 2 Lesson 2.3: Go Prep - Parallelism

**Prerequisites:** Lessons 2.1-2.2 complete (CSR + Columnar Storage)
**Time:** 4-5 hours Go prep + 25-30 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 2.3

## Overview

Parallel graph processing is essential for performance. Before implementing it, master these Go concepts:
- Worker pool pattern with channels
- Work stealing for load balancing
- `runtime.GOMAXPROCS` and core utilization
- **Go 1.25:** `testing/synctest` for deterministic tests ⭐⭐ **GAME CHANGER!**
- False sharing and cache line optimization
- Lock-free result aggregation with atomics

**The `testing/synctest` package makes testing concurrent code 100x easier!**

## Go Concepts for This Lesson

### 1. Worker Pool Pattern

**Distribute work across multiple goroutines!**

```go
package main

import (
    "fmt"
    "sync"
)

type Job struct {
    ID     int
    NodeID uint32
}

type Result struct {
    JobID       int
    NeighborCount int
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        // Simulate work (count neighbors)
        neighborCount := job.NodeID * 2  // Fake calculation

        results <- Result{
            JobID:         job.ID,
            NeighborCount: int(neighborCount),
        }
    }
}

func main() {
    numWorkers := 4
    numJobs := 100

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            worker(workerID, jobs, results)
        }(i)
    }

    // Send jobs
    go func() {
        for i := 0; i < numJobs; i++ {
            jobs <- Job{ID: i, NodeID: uint32(i)}
        }
        close(jobs)
    }()

    // Close results when all workers done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for result := range results {
        fmt.Printf("Job %d: %d neighbors\n", result.JobID, result.NeighborCount)
    }
}
```

**Key pattern:** Fixed number of workers, buffered channels for throughput.

### 2. Work Stealing for Load Balancing

**Better than static partitioning!**

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type WorkQueue struct {
    tasks   chan int
    counter atomic.Int64
}

func NewWorkQueue(tasks []int) *WorkQueue {
    wq := &WorkQueue{
        tasks: make(chan int, len(tasks)),
    }

    for _, task := range tasks {
        wq.tasks <- task
    }
    close(wq.tasks)

    return wq
}

func worker(id int, queue *WorkQueue, results *[]int, mu *sync.Mutex) {
    processed := 0

    for task := range queue.tasks {
        // Simulate variable work (some tasks take longer)
        result := task * task
        processed++

        // Steal work from queue dynamically
        mu.Lock()
        *results = append(*results, result)
        mu.Unlock()
    }

    fmt.Printf("Worker %d processed %d tasks\n", id, processed)
}

func main() {
    // Tasks with varying difficulty
    tasks := make([]int, 100)
    for i := range tasks {
        tasks[i] = i
    }

    queue := NewWorkQueue(tasks)

    var results []int
    var mu sync.Mutex
    var wg sync.WaitGroup

    numWorkers := 4

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(id, queue, &results, &mu)
        }(i)
    }

    wg.Wait()

    fmt.Printf("Processed %d results\n", len(results))
}
```

**Key insight:** Workers steal from shared queue → automatic load balancing!

### 3. runtime.GOMAXPROCS

**Control parallelism!**

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func cpuIntensiveWork() {
    // Burn CPU
    sum := 0
    for i := 0; i < 100000000; i++ {
        sum += i
    }
}

func benchmarkParallelism(numGoroutines int) time.Duration {
    start := time.Now()

    var wg sync.WaitGroup
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            cpuIntensiveWork()
        }()
    }

    wg.Wait()
    return time.Since(start)
}

func main() {
    fmt.Printf("CPU cores: %d\n", runtime.NumCPU())

    // Test with different GOMAXPROCS values
    for _, procs := range []int{1, 2, 4, 8} {
        runtime.GOMAXPROCS(procs)

        elapsed := benchmarkParallelism(8)
        fmt.Printf("GOMAXPROCS=%d: %v\n", procs, elapsed)
    }

    // Reset to default
    runtime.GOMAXPROCS(runtime.NumCPU())
}
```

**Typical output:**
```
CPU cores: 8
GOMAXPROCS=1: 2.5s  (no parallelism)
GOMAXPROCS=2: 1.3s  (2x speedup)
GOMAXPROCS=4: 700ms (3.5x speedup)
GOMAXPROCS=8: 400ms (6x speedup)
```

### 4. Go 1.25: testing/synctest (REVOLUTIONARY!)

**Deterministic testing of concurrent code!**

```go
package main

import (
    "sync/atomic"
    "testing"
    "testing/synctest"
    "time"
)

func TestConcurrentCounter(t *testing.T) {
    synctest.Run(func() {
        var counter atomic.Int32

        // Spawn 100 goroutines incrementing counter
        for i := 0; i < 100; i++ {
            go func() {
                time.Sleep(1 * time.Millisecond)  // Simulated delay
                counter.Add(1)
            }()
        }

        // Wait for all goroutines (deterministic!)
        synctest.Wait()

        if counter.Load() != 100 {
            t.Errorf("Expected 100, got %d", counter.Load())
        }
    })
}

func TestWorkStealing(t *testing.T) {
    synctest.Run(func() {
        queue := make(chan int, 100)
        results := make([]atomic.Int32, 4)

        // Send jobs
        for i := 0; i < 100; i++ {
            queue <- i
        }
        close(queue)

        // Start 4 workers
        for i := 0; i < 4; i++ {
            workerID := i
            go func() {
                for range queue {
                    time.Sleep(1 * time.Millisecond)  // Fake time!
                    results[workerID].Add(1)
                }
            }()
        }

        synctest.Wait()  // Deterministic wait!

        total := 0
        for i := 0; i < 4; i++ {
            count := results[i].Load()
            t.Logf("Worker %d processed %d jobs", i, count)
            total += int(count)
        }

        if total != 100 {
            t.Errorf("Expected 100 jobs processed, got %d", total)
        }
    })
}
```

**Key features:**
- **Deterministic execution:** Same result every time!
- **Fast:** `time.Sleep()` is simulated (instant in tests)!
- **Automatic deadlock detection:** Test fails if goroutines deadlock!

**Before synctest:** Flaky tests, real sleeps, race conditions.
**After synctest:** Fast, reliable, deterministic tests!

### 5. False Sharing and Cache Lines

**CRITICAL for performance!**

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "sync/atomic"
    "time"
)

// Bad: False sharing (counters on same cache line)
type BadCounters struct {
    counter1 atomic.Int64  // 8 bytes
    counter2 atomic.Int64  // 8 bytes (same cache line!)
    counter3 atomic.Int64
    counter4 atomic.Int64
}

// Good: Cache line padding (64 bytes apart)
type GoodCounters struct {
    counter1 atomic.Int64
    _        [7]int64  // 56 bytes padding (64 - 8 = 56)
    counter2 atomic.Int64
    _        [7]int64
    counter3 atomic.Int64
    _        [7]int64
    counter4 atomic.Int64
}

func benchmarkFalseSharing(padded bool) time.Duration {
    var wg sync.WaitGroup
    numWorkers := 4

    if padded {
        counters := &GoodCounters{}
        start := time.Now()

        for i := 0; i < numWorkers; i++ {
            wg.Add(1)
            go func(id int) {
                defer wg.Done()
                var counter *atomic.Int64
                switch id {
                case 0:
                    counter = &counters.counter1
                case 1:
                    counter = &counters.counter2
                case 2:
                    counter = &counters.counter3
                case 3:
                    counter = &counters.counter4
                }

                for j := 0; j < 10000000; j++ {
                    counter.Add(1)
                }
            }(i)
        }

        wg.Wait()
        return time.Since(start)
    } else {
        counters := &BadCounters{}
        start := time.Now()

        for i := 0; i < numWorkers; i++ {
            wg.Add(1)
            go func(id int) {
                defer wg.Done()
                var counter *atomic.Int64
                switch id {
                case 0:
                    counter = &counters.counter1
                case 1:
                    counter = &counters.counter2
                case 2:
                    counter = &counters.counter3
                case 3:
                    counter = &counters.counter4
                }

                for j := 0; j < 10000000; j++ {
                    counter.Add(1)
                }
            }(i)
        }

        wg.Wait()
        return time.Since(start)
    }
}

func main() {
    runtime.GOMAXPROCS(4)

    bad := benchmarkFalseSharing(false)
    good := benchmarkFalseSharing(true)

    fmt.Printf("False sharing: %v\n", bad)
    fmt.Printf("Cache-line padded: %v\n", good)
    fmt.Printf("Speedup: %.2fx\n", float64(bad)/float64(good))
}
```

**Typical results:** 2-4x speedup with padding!

**Why?** Without padding, each atomic operation invalidates other cores' cache lines.

### 6. Lock-Free Result Aggregation

**Aggregate results without locks!**

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type Stats struct {
    count atomic.Int64
    sum   atomic.Int64
    max   atomic.Int64
}

func (s *Stats) Add(value int64) {
    s.count.Add(1)
    s.sum.Add(value)

    // Update max (lock-free!)
    for {
        oldMax := s.max.Load()
        if value <= oldMax {
            break
        }
        if s.max.CompareAndSwap(oldMax, value) {
            break
        }
    }
}

func (s *Stats) Average() float64 {
    count := s.count.Load()
    if count == 0 {
        return 0
    }
    return float64(s.sum.Load()) / float64(count)
}

func main() {
    stats := &Stats{}
    var wg sync.WaitGroup

    // 100 goroutines adding values
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(val int64) {
            defer wg.Done()
            stats.Add(val)
        }(int64(i))
    }

    wg.Wait()

    fmt.Printf("Count: %d\n", stats.count.Load())
    fmt.Printf("Sum: %d\n", stats.sum.Load())
    fmt.Printf("Max: %d\n", stats.max.Load())
    fmt.Printf("Average: %.2f\n", stats.Average())
}
```

**Key technique:** `CompareAndSwap` for lock-free updates!

## Pre-Implementation Exercises

### Exercise 1: Parallel Graph Traversal

```go
package main

import (
    "iter"
    "sync"
)

type NodeID uint32

type Graph struct {
    offsets []uint64
    targets []NodeID
}

func (g *Graph) Neighbors(node NodeID) iter.Seq[NodeID] {
    // TODO: From Lesson 2.1
    return nil
}

// TODO: Implement parallel BFS
func (g *Graph) ParallelBFS(start NodeID, numWorkers int) []NodeID {
    visited := make(map[NodeID]bool)
    result := make([]NodeID, 0)

    // TODO:
    // 1. Create work queue with starting node
    // 2. Spawn workers to process nodes
    // 3. Each worker explores neighbors, adds new nodes to queue
    // 4. Track visited nodes (need mutex!)
    // 5. Return all visited nodes

    return result
}

func main() {
    // TODO: Test parallel BFS
}
```

### Exercise 2: Worker Pool

```go
package main

import (
    "sync"
)

type Job struct {
    NodeID uint32
}

type Result struct {
    NodeID   uint32
    Degree   int
}

type WorkerPool struct {
    numWorkers int
    jobs       chan Job
    results    chan Result
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    // TODO: Initialize worker pool
    return nil
}

func (wp *WorkerPool) Start() {
    // TODO: Start worker goroutines
}

func (wp *WorkerPool) Submit(job Job) {
    // TODO: Send job to workers
}

func (wp *WorkerPool) Results() <-chan Result {
    // TODO: Return results channel
    return nil
}

func (wp *WorkerPool) Close() {
    // TODO: Close jobs channel and wait for workers
}

func main() {
    // TODO: Test worker pool
}
```

### Exercise 3: Deterministic Tests with synctest

```go
package main

import (
    "sync/atomic"
    "testing"
    "testing/synctest"
    "time"
)

func TestParallelIncrement(t *testing.T) {
    synctest.Run(func() {
        var counter atomic.Int64

        // TODO:
        // 1. Spawn 10 goroutines
        // 2. Each increments counter 100 times
        // 3. Use synctest.Wait()
        // 4. Assert counter == 1000
    })
}

func TestWorkerPool(t *testing.T) {
    synctest.Run(func() {
        // TODO:
        // 1. Create worker pool with 4 workers
        // 2. Submit 100 jobs
        // 3. Collect results
        // 4. Assert all jobs completed
    })
}
```

### Exercise 4: Cache-Line Padded Counters

```go
package main

import (
    "sync/atomic"
)

type PaddedCounters struct {
    // TODO: Add 4 counters with cache-line padding
    // Each counter should be 64 bytes apart
}

func NewPaddedCounters() *PaddedCounters {
    // TODO: Initialize
    return nil
}

func (pc *PaddedCounters) Increment(id int) {
    // TODO: Increment counter by ID (0-3)
}

func (pc *PaddedCounters) Get(id int) int64 {
    // TODO: Get counter value
    return 0
}

func main() {
    // TODO: Benchmark with and without padding
}
```

### Exercise 5: Lock-Free Statistics

```go
package main

import (
    "sync/atomic"
)

type LockFreeStats struct {
    count atomic.Int64
    sum   atomic.Int64
    min   atomic.Int64
    max   atomic.Int64
}

func NewLockFreeStats() *LockFreeStats {
    stats := &LockFreeStats{}
    stats.min.Store(1<<63 - 1)  // Max int64
    stats.max.Store(-1 << 63)   // Min int64
    return stats
}

func (s *LockFreeStats) Add(value int64) {
    // TODO:
    // 1. Increment count
    // 2. Add to sum
    // 3. Update min with CompareAndSwap
    // 4. Update max with CompareAndSwap
}

func (s *LockFreeStats) Average() float64 {
    // TODO: Return average
    return 0
}

func (s *LockFreeStats) Range() int64 {
    // TODO: Return max - min
    return 0
}

func main() {
    // TODO: Test with concurrent updates
}
```

## Performance Benchmarks

### Benchmark 1: Sequential vs Parallel

```go
func BenchmarkSequentialBFS(b *testing.B) {
    g := makeGraph(10000, 10)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        g.BFS(0)
    }
}

func BenchmarkParallelBFS(b *testing.B) {
    g := makeGraph(10000, 10)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        g.ParallelBFS(0, 4)
    }
}
```

### Benchmark 2: Scaling with Core Count

```go
func BenchmarkParallelScaling(b *testing.B) {
    g := makeGraph(100000, 10)

    for _, numWorkers := range []int{1, 2, 4, 8, 16} {
        b.Run(fmt.Sprintf("workers-%d", numWorkers), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                g.ParallelBFS(0, numWorkers)
            }
        })
    }
}
```

### Benchmark 3: False Sharing Impact

```go
func BenchmarkFalseSharingCounters(b *testing.B) {
    counters := &BadCounters{}

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counters.counter1.Add(1)
        }
    })
}

func BenchmarkPaddedCounters(b *testing.B) {
    counters := &GoodCounters{}

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counters.counter1.Add(1)
        }
    })
}
```

## Common Gotchas to Avoid

### Gotcha 1: Too Many Goroutines

```go
// WRONG: Spawning millions of goroutines!
for i := 0; i < 1000000; i++ {
    go func(id int) {
        processNode(id)
    }(i)
}
// Overhead > benefit!

// RIGHT: Worker pool with fixed goroutines
numWorkers := runtime.NumCPU()
jobs := make(chan int, 1000)

for i := 0; i < numWorkers; i++ {
    go worker(jobs)
}

for i := 0; i < 1000000; i++ {
    jobs <- i
}
```

### Gotcha 2: Forgetting synctest.Wait()

```go
// WRONG: Test exits before goroutines finish
func TestConcurrent(t *testing.T) {
    synctest.Run(func() {
        go func() {
            // Do work...
        }()
        // Forgot synctest.Wait()!
    })
}

// RIGHT: Always wait
func TestConcurrent(t *testing.T) {
    synctest.Run(func() {
        go func() {
            // Do work...
        }()
        synctest.Wait()  // Wait for goroutines!
    })
}
```

### Gotcha 3: False Sharing

```go
// WRONG: Counters on same cache line
type Counters struct {
    a atomic.Int64  // 8 bytes
    b atomic.Int64  // 8 bytes (same cache line!)
}

// Multiple cores writing to a and b = cache thrashing!

// RIGHT: Pad to 64 bytes
type Counters struct {
    a atomic.Int64
    _ [7]int64  // Padding
    b atomic.Int64
    _ [7]int64
}
```

### Gotcha 4: Not Checking Channel Closure

```go
// WRONG: Sending on closed channel
close(jobs)
jobs <- Job{ID: 1}  // PANIC!

// RIGHT: Check if channel closed
select {
case jobs <- Job{ID: 1}:
    // Sent
default:
    // Channel full or closed
}
```

### Gotcha 5: Inefficient CompareAndSwap Loop

```go
// WRONG: Busy-waiting without yield
for {
    old := value.Load()
    new := old + 1
    if value.CompareAndSwap(old, new) {
        break
    }
    // Spins forever under contention!
}

// BETTER: Add backoff
for {
    old := value.Load()
    new := old + 1
    if value.CompareAndSwap(old, new) {
        break
    }
    runtime.Gosched()  // Yield to other goroutines
}
```

## Checklist Before Starting Lesson 2.3

- [ ] I can implement worker pool pattern
- [ ] I understand work stealing with channels
- [ ] I know how to use `runtime.GOMAXPROCS`
- [ ] I can write deterministic tests with `testing/synctest`
- [ ] I understand false sharing and cache lines
- [ ] I can use cache-line padding for performance
- [ ] I know how to use `atomic.CompareAndSwap`
- [ ] I can aggregate results lock-free with atomics
- [ ] I understand when to use locks vs atomics
- [ ] I can benchmark parallel speedup

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 2.3

You'll implement:
- Parallel graph traversal (BFS/DFS)
- Worker pool for query processing
- Lock-free result aggregation
- Cache-line optimized data structures
- Deterministic concurrency tests with synctest
- Scaling benchmarks (1-16 cores)

**Time estimate:** 25-30 hours for full implementation

**Go 1.25's synctest makes this lesson way easier than before!**

Good luck! 🚀
