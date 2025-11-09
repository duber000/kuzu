# Pre-work Week 5-6: Concurrency & I/O

**Duration:** 2 weeks | **Time Commitment:** 15-20 hours/week | **Difficulty:** Advanced

## Overview

This is the **most critical** pre-work module for database development. Databases are inherently concurrent (multiple queries running simultaneously) and I/O intensive (reading/writing to disk).

By the end of Week 6, you'll:
- ✅ Master goroutines and channels
- ✅ Understand mutexes and avoid race conditions
- ✅ Work with file I/O and system calls
- ✅ Build a concurrent file downloader
- ✅ Be ready to start the storage layer implementation

## Week 5: Concurrency

### Day 1-2: Goroutines

Goroutines are lightweight threads managed by the Go runtime.

#### Your First Goroutine

```go
package main

import (
    "fmt"
    "time"
)

func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    // Start a goroutine
    go sayHello()

    // Main goroutine
    fmt.Println("Hello from main!")

    // Wait a bit (not ideal, we'll learn better ways)
    time.Sleep(time.Second)
}
```

**Important:** The `main` function is itself a goroutine. When it exits, all other goroutines are terminated immediately.

#### Multiple Goroutines

```go
package main

import (
    "fmt"
    "time"
)

func count(id int) {
    for i := 1; i <= 5; i++ {
        fmt.Printf("Goroutine %d: %d\n", id, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Start 3 goroutines
    for i := 1; i <= 3; i++ {
        go count(i)
    }

    // Wait for them to finish
    time.Sleep(time.Second)
    fmt.Println("Done!")
}
```

**Observation:** The output is interleaved! Goroutines run concurrently.

#### The Problem with Sleep

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("Goroutine finished")
    }()

    // If we sleep too little, goroutine doesn't finish
    time.Sleep(50 * time.Millisecond)
    // Program exits, goroutine is killed!
}
```

**Solution:** We need proper synchronization!

### Day 3-4: Channels

Channels are Go's way of communicating between goroutines.

#### Basic Channel Usage

```go
package main

import "fmt"

func main() {
    // Create a channel
    messages := make(chan string)

    // Send value in a goroutine
    go func() {
        messages <- "Hello"  // Send
    }()

    // Receive value
    msg := <-messages  // Receive (blocks until value available)
    fmt.Println(msg)
}
```

**Key concept:**
- `ch <- value` sends to channel
- `value := <-ch` receives from channel
- Channels block until both sender and receiver are ready

#### Buffered Channels

```go
package main

import "fmt"

func main() {
    // Buffered channel (capacity 2)
    ch := make(chan int, 2)

    // Can send 2 values without blocking
    ch <- 1
    ch <- 2

    // ch <- 3  // This would block!

    fmt.Println(<-ch)  // 1
    fmt.Println(<-ch)  // 2
}
```

#### Channel Patterns

**Pattern 1: Pipeline**

```go
package main

import "fmt"

func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)  // Important: close when done
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {  // Receive until channel is closed
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    // Pipeline: generate -> square
    nums := generator(1, 2, 3, 4)
    squares := square(nums)

    for s := range squares {
        fmt.Println(s)  // 1, 4, 9, 16
    }
}
```

**Pattern 2: Fan-Out (Multiple workers)**

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 10)
    results := make(chan int, 10)

    var wg sync.WaitGroup

    // Start 3 workers
    for w := 1; w <= 3; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    // Send 9 jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Wait for all workers to finish
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for r := range results {
        fmt.Println("Result:", r)
    }
}
```

**Pattern 3: Select (Multiple Channels)**

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "from channel 1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "from channel 2"
    }()

    // Select waits on multiple channels
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println(msg1)
        case msg2 := <-ch2:
            fmt.Println(msg2)
        }
    }
}
```

#### Interactive Exercise 1: Concurrent Sum

```go
package main

import (
    "fmt"
    "sync"
)

// TODO: Implement concurrentSum
// Split the slice into chunks, sum each chunk in a goroutine,
// combine results using channels
func concurrentSum(numbers []int, workers int) int {
    // Your code here
    return 0
}

func main() {
    numbers := make([]int, 1000000)
    for i := range numbers {
        numbers[i] = i + 1
    }

    result := concurrentSum(numbers, 4)
    fmt.Println("Sum:", result)  // Should be 500000500000
}
```

**Hint:**
1. Create a channel for partial sums
2. Split the slice into `workers` chunks
3. Start `workers` goroutines, each sums a chunk
4. Collect partial sums and combine

### Day 5-6: Synchronization with Mutexes

Sometimes channels are overkill. Mutexes protect shared data.

#### The Race Condition Problem

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    counter := 0
    var wg sync.WaitGroup

    // 1000 goroutines incrementing counter
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++  // RACE CONDITION!
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter)  // Not 1000!
}
```

Run with race detector:
```bash
go run -race main.go
```

You'll see: `WARNING: DATA RACE`

#### Solution 1: Mutex

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    counter := 0
    var mu sync.Mutex
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter)  // Now 1000!
}
```

#### Solution 2: Atomic Operations

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var counter int64  // Must be int32 or int64 for atomic
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)  // Atomic increment
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", atomic.LoadInt64(&counter))
}
```

**When to use what:**
- **Channels:** Communicating between goroutines, passing ownership
- **Mutex:** Protecting shared state
- **Atomic:** Simple counters and flags

#### Read-Write Mutex

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type SafeCache struct {
    mu    sync.RWMutex
    cache map[string]string
}

func (c *SafeCache) Get(key string) (string, bool) {
    c.mu.RLock()  // Read lock (multiple readers OK)
    defer c.mu.RUnlock()
    val, ok := c.cache[key]
    return val, ok
}

func (c *SafeCache) Set(key, value string) {
    c.mu.Lock()  // Write lock (exclusive)
    defer c.mu.Unlock()
    c.cache[key] = value
}

func main() {
    cache := &SafeCache{
        cache: make(map[string]string),
    }

    // Many concurrent readers
    for i := 0; i < 100; i++ {
        go func(id int) {
            for j := 0; j < 1000; j++ {
                cache.Get("key")
                time.Sleep(time.Microsecond)
            }
        }(i)
    }

    // Few writers
    for i := 0; i < 10; i++ {
        go func(id int) {
            for j := 0; j < 100; j++ {
                cache.Set("key", fmt.Sprintf("value%d", id))
                time.Sleep(time.Millisecond)
            }
        }(i)
    }

    time.Sleep(2 * time.Second)
}
```

#### Interactive Exercise 2: Thread-Safe Counter

```go
package main

import (
    "fmt"
    "sync"
)

type Counter struct {
    // TODO: Add fields for mutex and value
}

func (c *Counter) Increment() {
    // TODO: Thread-safe increment
}

func (c *Counter) Decrement() {
    // TODO: Thread-safe decrement
}

func (c *Counter) Value() int {
    // TODO: Thread-safe read
    return 0
}

func main() {
    counter := &Counter{}
    var wg sync.WaitGroup

    // 100 goroutines incrementing
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                counter.Increment()
            }
        }()
    }

    // 50 goroutines decrementing
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                counter.Decrement()
            }
        }()
    }

    wg.Wait()
    fmt.Println("Final value:", counter.Value())  // Should be 50000

    // Verify no race conditions
    // go run -race main.go
}
```

### Day 7: Advanced Patterns

#### Worker Pool Pattern

```go
package main

import (
    "fmt"
    "time"
)

type Job struct {
    ID int
}

type Result struct {
    JobID int
    Value int
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(time.Second)  // Simulate work
        results <- Result{JobID: job.ID, Value: job.ID * 2}
    }
}

func main() {
    numJobs := 10
    numWorkers := 3

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    // Start workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j}
    }
    close(jobs)

    // Collect results
    for r := 1; r <= numJobs; r++ {
        result := <-results
        fmt.Printf("Result: Job %d -> %d\n", result.JobID, result.Value)
    }
}
```

## Week 6: File I/O and System Calls

### Day 1-2: Basic File Operations

#### Reading Files

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    // Method 1: Read entire file
    content, err := os.ReadFile("test.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(content))

    // Method 2: Open and read
    file, err := os.Open("test.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    buffer := make([]byte, 1024)
    for {
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }
        fmt.Print(string(buffer[:n]))
    }
}
```

#### Writing Files

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Method 1: Write entire file
    data := []byte("Hello, World!\n")
    err := os.WriteFile("output.txt", data, 0644)
    if err != nil {
        panic(err)
    }

    // Method 2: Open and write
    file, err := os.Create("output2.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    _, err = file.WriteString("Line 1\n")
    if err != nil {
        panic(err)
    }

    // Write bytes
    _, err = file.Write([]byte("Line 2\n"))
    if err != nil {
        panic(err)
    }
}
```

#### Buffered I/O

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // Buffered reading (much faster for large files)
    file, _ := os.Open("large_file.txt")
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineCount := 0

    for scanner.Scan() {
        lineCount++
        // scanner.Text() gives you the line
    }

    fmt.Println("Lines:", lineCount)

    // Buffered writing
    outFile, _ := os.Create("output.txt")
    defer outFile.Close()

    writer := bufio.NewWriter(outFile)
    defer writer.Flush()  // Important: flush buffered data

    for i := 0; i < 1000; i++ {
        fmt.Fprintf(writer, "Line %d\n", i)
    }
}
```

### Day 3-4: System Calls and Performance

#### fsync for Durability

```go
package main

import (
    "os"
    "time"
)

func main() {
    file, _ := os.Create("important.dat")
    defer file.Close()

    // Write data
    file.WriteString("Critical data")

    // Without Sync(): data might still be in OS buffer!
    // If power fails, data could be lost

    start := time.Now()
    file.Sync()  // Force write to disk
    elapsed := time.Since(start)

    println("Sync took:", elapsed.Milliseconds(), "ms")
    // On SSD: 1-10ms
    // On HDD: 5-20ms
}
```

**This is critical for databases:** Data must survive crashes!

#### Direct I/O (Advanced)

```go
package main

import (
    "fmt"
    "os"
    "syscall"
)

func main() {
    // Open file with O_DIRECT flag (bypasses OS cache)
    // This is used in databases for predictable performance

    fd, err := syscall.Open(
        "data.db",
        syscall.O_CREAT|syscall.O_RDWR|syscall.O_DIRECT,
        0644,
    )
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer syscall.Close(fd)

    // O_DIRECT requires aligned buffers
    // This is complex - we'll cover it in the storage layer lessons
}
```

#### Memory-Mapped Files

```go
package main

import (
    "fmt"
    "os"
    "syscall"
)

func main() {
    // Create a file
    file, _ := os.Create("mmap_test.dat")
    file.Truncate(4096)  // 4KB
    defer file.Close()

    // Memory map it
    data, err := syscall.Mmap(
        int(file.Fd()),
        0,
        4096,
        syscall.PROT_READ|syscall.PROT_WRITE,
        syscall.MAP_SHARED,
    )
    if err != nil {
        panic(err)
    }
    defer syscall.Munmap(data)

    // Write to memory = write to file!
    copy(data, []byte("Hello, mmap!"))

    // Force to disk
    syscall.Msync(data, syscall.MS_SYNC)

    fmt.Println("Data written:", string(data[:12]))
}
```

**Why databases use mmap:**
- Fast random access
- OS handles paging
- Less syscall overhead

### Day 5-7: Project - Concurrent File Downloader

**Build a tool that downloads multiple URLs concurrently!**

#### Requirements

- Download multiple files in parallel
- Progress reporting
- Error handling and retries
- Rate limiting (max N concurrent downloads)
- Save files to disk

#### Implementation

Create `downloader/downloader.go`:

```go
package downloader

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "sync"
)

type Download struct {
    URL      string
    Filename string
}

type Result struct {
    Download Download
    Error    error
    Bytes    int64
}

type Downloader struct {
    MaxConcurrent int
    OutputDir     string
}

func New(maxConcurrent int, outputDir string) *Downloader {
    return &Downloader{
        MaxConcurrent: maxConcurrent,
        OutputDir:     outputDir,
    }
}

func (d *Downloader) Download(downloads []Download) []Result {
    // Create semaphore channel for rate limiting
    sem := make(chan struct{}, d.MaxConcurrent)
    results := make(chan Result, len(downloads))
    var wg sync.WaitGroup

    // Start downloads
    for _, dl := range downloads {
        wg.Add(1)
        go func(download Download) {
            defer wg.Done()

            // Acquire semaphore
            sem <- struct{}{}
            defer func() { <-sem }()

            // Perform download
            result := d.downloadOne(download)
            results <- result
        }(dl)
    }

    // Wait and close results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    var allResults []Result
    for result := range results {
        allResults = append(allResults, result)
    }

    return allResults
}

func (d *Downloader) downloadOne(dl Download) Result {
    // Create output file path
    outputPath := filepath.Join(d.OutputDir, dl.Filename)

    // Download
    resp, err := http.Get(dl.URL)
    if err != nil {
        return Result{Download: dl, Error: err}
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return Result{
            Download: dl,
            Error:    fmt.Errorf("bad status: %s", resp.Status),
        }
    }

    // Create file
    file, err := os.Create(outputPath)
    if err != nil {
        return Result{Download: dl, Error: err}
    }
    defer file.Close()

    // Copy data
    bytes, err := io.Copy(file, resp.Body)
    if err != nil {
        return Result{Download: dl, Error: err}
    }

    fmt.Printf("Downloaded %s (%d bytes)\n", dl.Filename, bytes)

    return Result{
        Download: dl,
        Bytes:    bytes,
    }
}
```

Create `main.go`:

```go
package main

import (
    "fmt"
    "os"

    "downloader"
)

func main() {
    // Create output directory
    os.MkdirAll("downloads", 0755)

    // Define downloads
    downloads := []downloader.Download{
        {
            URL:      "https://go.dev/robots.txt",
            Filename: "go-robots.txt",
        },
        {
            URL:      "https://www.google.com/robots.txt",
            Filename: "google-robots.txt",
        },
        {
            URL:      "https://github.com/robots.txt",
            Filename: "github-robots.txt",
        },
        // Add more URLs...
    }

    // Create downloader (max 3 concurrent)
    dl := downloader.New(3, "downloads")

    fmt.Printf("Downloading %d files...\n", len(downloads))

    // Download
    results := dl.Download(downloads)

    // Report
    successful := 0
    failed := 0
    totalBytes := int64(0)

    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("❌ %s: %v\n", result.Download.Filename, result.Error)
            failed++
        } else {
            fmt.Printf("✓ %s: %d bytes\n", result.Download.Filename, result.Bytes)
            successful++
            totalBytes += result.Bytes
        }
    }

    fmt.Printf("\nResults: %d successful, %d failed\n", successful, failed)
    fmt.Printf("Total downloaded: %d bytes\n", totalBytes)
}
```

#### Enhancement Challenges

1. **Progress bar:** Show download progress for each file
2. **Retry logic:** Retry failed downloads up to N times
3. **Bandwidth limiting:** Limit total download speed
4. **Resume support:** Resume partial downloads
5. **Checksum verification:** Verify file integrity with SHA256
6. **Timeout handling:** Cancel slow downloads

## Week 5-6 Checkpoint

### Self-Assessment

Can you answer these?

1. What's the difference between a goroutine and a thread?
2. How do unbuffered channels differ from buffered channels?
3. When should you use a mutex vs a channel?
4. What does `go run -race` detect?
5. Why is `file.Sync()` important for databases?
6. What's the purpose of memory-mapped files?

### Practical Test

**Build a concurrent URL checker:**

```go
// Check if 100 URLs are accessible
// Use a worker pool pattern
// Report:
//   - How many are up (200 OK)
//   - How many are down
//   - Average response time
// Requirements:
//   - Max 10 concurrent requests
//   - Use channels for communication
//   - No race conditions (verify with -race)
//   - Timeout after 5 seconds per URL
```

If you can build this in 4-5 hours, you're ready for Phase 1!

## What You've Mastered

- ✅ Goroutines and concurrent programming
- ✅ Channels for communication
- ✅ Mutexes and atomic operations
- ✅ Race condition detection
- ✅ File I/O and buffering
- ✅ System calls (fsync, mmap)
- ✅ Worker pool pattern
- ✅ Concurrent file downloading

## What's Next?

**You're now ready for the main curriculum!**

Start with: **[Phase 1 Lesson 1.1: Go Prep - Pages](../phase-1-storage/go-prep-lesson-1-1-pages.md)**

The storage layer will build on everything you've learned:
- Concurrency from Week 5-6
- File I/O and system calls from Week 6
- Testing and benchmarking from Week 3-4
- Pointers and memory from Week 3

## Resources

**Concurrency Deep Dive:**
- [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs) - Rob Pike (Must watch!)
- [Advanced Go Concurrency](https://www.youtube.com/watch?v=QDDwwePbDtw) - Sameer Ajmani
- [Go Memory Model](https://go.dev/ref/mem) - Official docs

**I/O and Systems:**
- [Linux System Call Table](https://man7.org/linux/man-pages/man2/syscalls.2.html)
- [Understanding mmap](https://stackoverflow.com/questions/258091/when-should-i-use-mmap-for-file-access)

**Practice:**
- Build a web crawler with worker pools
- Create a concurrent chat server
- Implement a rate limiter

**Congratulations!** You've completed the pre-work. You're ready to build a database! 🚀
