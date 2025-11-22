# Phase 1 Lesson 1.1: Go Prep - Page Abstraction

**Prerequisites:** Pre-work Weeks 1-6 complete
**Time:** 2-3 hours Go prep + 15-20 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 1.1

## Overview

Before implementing the page abstraction, you need to master these Go-specific concepts:
- Memory-mapped I/O with `syscall.Mmap`
- Binary encoding/decoding
- Fixed-size arrays vs slices
- Memory alignment and padding
- Basic benchmarking for I/O operations

## Go Concepts for This Lesson

### 1. Fixed-Size Arrays (Critical for Pages!)

Pages are fixed size (4096 bytes). Use arrays, not slices.

```go
package main

import (
    "fmt"
    "unsafe"
)

const PageSize = 4096

// Page is a fixed-size structure
type Page struct {
    Data [PageSize]byte  // Fixed-size array
}

func main() {
    var page Page

    fmt.Println("Page size:", unsafe.Sizeof(page))  // Exactly 4096 bytes

    // Arrays are value types - they're copied!
    page2 := page  // Full copy of 4KB
    page2.Data[0] = 42

    fmt.Println(page.Data[0])   // 0
    fmt.Println(page2.Data[0])  // 42

    // To avoid copying, use pointers
    pagePtr := &page
    pagePtr.Data[0] = 99
    fmt.Println(page.Data[0])  // 99
}
```

**Key takeaway:** Pass `*Page` to functions, never `Page` (would copy 4KB every call!).

### 2. Binary Encoding/Decoding

Store integers efficiently in bytes.

```go
package main

import (
    "encoding/binary"
    "fmt"
)

func main() {
    // Write integer to bytes
    var buf [8]byte
    binary.LittleEndian.PutUint64(buf[:], 12345)

    fmt.Printf("Bytes: %v\n", buf)

    // Read integer from bytes
    value := binary.LittleEndian.Uint64(buf[:])
    fmt.Println("Value:", value)  // 12345

    // Why LittleEndian? Most modern CPUs are little-endian
    // BigEndian would also work, but slightly slower
}
```

**Practice Exercise:**

```go
package main

import (
    "encoding/binary"
    "fmt"
)

// TODO: Implement a Page that stores integers
type IntPage struct {
    Data [4096]byte
}

// WriteInt writes an integer at the given offset
func (p *IntPage) WriteInt(offset int, value uint64) {
    // TODO: Write value at offset using binary.LittleEndian
}

// ReadInt reads an integer from the given offset
func (p *IntPage) ReadInt(offset int) uint64 {
    // TODO: Read value from offset
    return 0
}

func main() {
    var page IntPage

    // Page can hold 4096 / 8 = 512 integers
    page.WriteInt(0, 100)
    page.WriteInt(8, 200)

    fmt.Println(page.ReadInt(0))  // Should print 100
    fmt.Println(page.ReadInt(8))  // Should print 200
}
```

### 3. Memory-Mapped I/O

Map a file directly into memory.

```go
package main

import (
    "fmt"
    "os"
    "syscall"
)

func main() {
    // Create a file
    file, err := os.Create("test.dat")
    if err != nil {
        panic(err)
    }

    // Extend to 8KB
    file.Truncate(8192)
    defer file.Close()

    // Memory map the file
    data, err := syscall.Mmap(
        int(file.Fd()),
        0,                    // offset
        8192,                 // length
        syscall.PROT_READ|syscall.PROT_WRITE,  // permissions
        syscall.MAP_SHARED,   // changes written to file
    )
    if err != nil {
        panic(err)
    }
    defer syscall.Munmap(data)

    // Write to memory = write to file!
    copy(data, []byte("Hello, mmap!"))

    // Force write to disk
    syscall.Msync(data, syscall.MS_SYNC)

    fmt.Println("Wrote to file via mmap")

    // Read it back
    fmt.Println("Data:", string(data[:12]))

    // Cleanup
    os.Remove("test.dat")
}
```

**Why databases use mmap:**
- Random access without syscall overhead
- OS handles caching
- Simpler code than manual read/write

**Gotcha:** mmap on Windows is different (uses `CreateFileMapping`). Focus on Linux/Mac for learning.

### 4. File I/O Patterns for Databases

#### Pattern 1: Direct Write with fsync

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func writePageDirect(filename string, pageID int, data []byte) error {
    // Open in read-write mode
    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // Seek to page offset
    offset := int64(pageID) * 4096
    _, err = file.Seek(offset, 0)
    if err != nil {
        return err
    }

    // Write
    _, err = file.Write(data)
    if err != nil {
        return err
    }

    // CRITICAL: Sync to disk (durability!)
    start := time.Now()
    err = file.Sync()
    fmt.Printf("Sync took: %v\n", time.Since(start))

    return err
}

func main() {
    page := make([]byte, 4096)
    copy(page, []byte("Page 0 data"))

    writePageDirect("data.db", 0, page)

    os.Remove("data.db")
}
```

**Benchmark this:** How long does `Sync()` take on your system?

#### Pattern 2: Read Page by ID

```go
package main

import (
    "fmt"
    "os"
)

func readPage(filename string, pageID int) ([]byte, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Seek to page
    offset := int64(pageID) * 4096
    _, err = file.Seek(offset, 0)
    if err != nil {
        return nil, err
    }

    // Read exactly 4096 bytes
    page := make([]byte, 4096)
    n, err := file.Read(page)
    if err != nil {
        return nil, err
    }

    if n != 4096 {
        return nil, fmt.Errorf("short read: got %d bytes", n)
    }

    return page, nil
}
```

### 5. Alignment and Padding

**Critical for performance!**

```go
package main

import (
    "fmt"
    "unsafe"
)

// Bad: Unaligned struct
type BadPage struct {
    ID       uint8   // 1 byte
    Count    uint64  // 8 bytes
    IsActive uint8   // 1 byte
    Data     [4080]byte
}

// Good: Aligned struct
type GoodPage struct {
    Count    uint64  // 8 bytes first
    ID       uint8   // 1 byte
    IsActive uint8   // 1 byte
    _        [6]byte // Explicit padding
    Data     [4080]byte
}

func main() {
    fmt.Println("BadPage size:", unsafe.Sizeof(BadPage{}))   // Likely 4104 due to padding
    fmt.Println("GoodPage size:", unsafe.Sizeof(GoodPage{})) // Exactly 4096

    // Field alignment
    bad := &BadPage{}
    good := &GoodPage{}

    fmt.Printf("Bad Count alignment: %d\n", unsafe.Offsetof(bad.Count))
    fmt.Printf("Good Count alignment: %d\n", unsafe.Offsetof(good.Count))

    // Good: Count at offset 0 (aligned on 8-byte boundary)
    // Bad: Count at offset 1 (unaligned, slower access!)
}
```

**Rule:** Put largest fields first, then sort by size descending.

## Pre-Implementation Exercises

Complete these BEFORE starting the main lesson:

### Exercise 1: Page Manager

```go
package main

import (
    "encoding/binary"
    "os"
)

const PageSize = 4096

type PageManager struct {
    file *os.File
}

func NewPageManager(filename string) (*PageManager, error) {
    // TODO: Open or create file
    return nil, nil
}

func (pm *PageManager) WritePage(pageID int, data []byte) error {
    // TODO: Write page at offset pageID * PageSize
    // Ensure data is exactly PageSize bytes
    // Call Sync() for durability
    return nil
}

func (pm *PageManager) ReadPage(pageID int) ([]byte, error) {
    // TODO: Read page from offset pageID * PageSize
    return nil, nil
}

func (pm *PageManager) Close() error {
    // TODO: Close file
    return nil
}

func main() {
    pm, _ := NewPageManager("test.db")
    defer pm.Close()

    // Write page 0
    page0 := make([]byte, PageSize)
    binary.LittleEndian.PutUint64(page0, 12345)
    pm.WritePage(0, page0)

    // Write page 5
    page5 := make([]byte, PageSize)
    binary.LittleEndian.PutUint64(page5, 67890)
    pm.WritePage(5, page5)

    // Read them back
    p0, _ := pm.ReadPage(0)
    p5, _ := pm.ReadPage(5)

    println("Page 0:", binary.LittleEndian.Uint64(p0))  // 12345
    println("Page 5:", binary.LittleEndian.Uint64(p5))  // 67890

    os.Remove("test.db")
}
```

### Exercise 2: Benchmark Page I/O

```go
package main

import (
    "os"
    "testing"
)

const PageSize = 4096

func BenchmarkPageWrite(b *testing.B) {
    file, _ := os.Create("bench.db")
    defer file.Close()
    defer os.Remove("bench.db")

    page := make([]byte, PageSize)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        file.WriteAt(page, int64(i%100)*PageSize)
    }
}

func BenchmarkPageWriteWithSync(b *testing.B) {
    file, _ := os.Create("bench_sync.db")
    defer file.Close()
    defer os.Remove("bench_sync.db")

    page := make([]byte, PageSize)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        file.WriteAt(page, 0)
        file.Sync()  // SLOW!
    }
}

func BenchmarkPageRead(b *testing.B) {
    file, _ := os.Create("bench_read.db")
    file.Truncate(PageSize * 1000)
    defer file.Close()
    defer os.Remove("bench_read.db")

    page := make([]byte, PageSize)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        file.ReadAt(page, int64(i%1000)*PageSize)
    }
}
```

Run and analyze:
```bash
go test -bench=. -benchmem
```

**Questions to answer:**
1. How fast is a page write without sync?
2. How much does Sync() slow things down?
3. How fast can you read pages?
4. Are there any allocations in the hot path?

### Exercise 3: mmap vs Regular I/O

```go
package main

import (
    "os"
    "syscall"
    "testing"
)

const PageSize = 4096

func BenchmarkMmapWrite(b *testing.B) {
    file, _ := os.Create("mmap.db")
    file.Truncate(PageSize * 1000)
    defer file.Close()
    defer os.Remove("mmap.db")

    data, _ := syscall.Mmap(int(file.Fd()), 0, PageSize*1000,
        syscall.PROT_WRITE, syscall.MAP_SHARED)
    defer syscall.Munmap(data)

    page := make([]byte, PageSize)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        offset := (i % 1000) * PageSize
        copy(data[offset:offset+PageSize], page)
    }
}

func BenchmarkRegularWrite(b *testing.B) {
    file, _ := os.Create("regular.db")
    defer file.Close()
    defer os.Remove("regular.db")

    page := make([]byte, PageSize)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        offset := int64(i%1000) * PageSize
        file.WriteAt(page, offset)
    }
}
```

**Compare:** Which is faster? Why?

## Key Gotchas to Avoid

### Gotcha 1: Forgetting to Sync

```go
// WRONG: Data might not hit disk!
func writePage(file *os.File, data []byte) {
    file.Write(data)
    // If program crashes here, data is lost!
}

// RIGHT: Ensure durability
func writePage(file *os.File, data []byte) error {
    if _, err := file.Write(data); err != nil {
        return err
    }
    return file.Sync()  // Guaranteed on disk
}
```

### Gotcha 2: Copying Large Arrays

```go
// WRONG: Copies 4KB on every call!
func processPage(page Page) {
    // ...
}

// RIGHT: Pass by pointer
func processPage(page *Page) {
    // ...
}
```

### Gotcha 3: Off-by-One in Page Calculation

```go
// WRONG: Integer division truncates
numPages := fileSize / PageSize  // Misses partial page!

// RIGHT: Round up
numPages := (fileSize + PageSize - 1) / PageSize
```

## Checklist Before Starting Lesson 1.1

- [ ] I understand the difference between arrays and slices
- [ ] I can use `binary.LittleEndian` to encode/decode integers
- [ ] I know how to open, seek, read, and write files
- [ ] I understand why `Sync()` is critical
- [ ] I've experimented with `syscall.Mmap`
- [ ] I've benchmarked page I/O on my system
- [ ] I know how to pass pages by pointer
- [ ] I understand memory alignment basics

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 1.1

You'll implement:
- Page structure with fixed-size arrays
- Page manager with O(1) access
- Crash-safe writes with fsync
- Benchmarks comparing different I/O approaches

**Time estimate:** 15-20 hours for full implementation

Good luck! 🚀
