# Phase 2 Lesson 2.2: Go Prep - Columnar Storage

**Prerequisites:** Lesson 2.1 complete (CSR + Iterators)
**Time:** 3-4 hours Go prep + 20-25 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 2.2

## Overview

Columnar storage stores data by column instead of by row, enabling faster analytics. Before implementing it, master these Go concepts:
- **Go 1.23:** `unique` package for string interning ⭐⭐ **IMPORTANT!**
- `unique.Handle[T]` for low-cardinality columns
- Memory profiling with `pprof`
- Bit manipulation for compression
- Validity bitmaps for NULL handling

**The `unique` package can reduce memory usage by 90% for string columns!**

## Go Concepts for This Lesson

### 1. Go 1.23 unique Package: String Interning

**New in Go 1.23: Automatic string deduplication!**

```go
package main

import (
    "fmt"
    "unique"
)

func main() {
    // Old way: Each string allocates memory
    s1 := "hello"
    s2 := "hello"
    s3 := "hello"

    fmt.Printf("s1 addr: %p\n", &s1)  // Different addresses
    fmt.Printf("s2 addr: %p\n", &s2)
    fmt.Printf("s3 addr: %p\n", &s3)

    // New way: Interned strings share memory!
    h1 := unique.Make("hello")
    h2 := unique.Make("hello")
    h3 := unique.Make("hello")

    fmt.Printf("h1 == h2: %v\n", h1 == h2)  // true (pointer equality!)
    fmt.Printf("h2 == h3: %v\n", h2 == h3)  // true

    // Get the value back
    fmt.Printf("Value: %s\n", h1.Value())

    // Fast comparisons (pointer comparison, not string comparison!)
    if h1 == h2 {
        fmt.Println("Equal (compared in O(1)!)")
    }
}
```

**Key insight:** `unique.Handle` stores a pointer to the interned string. Equality check is O(1) pointer comparison!

### 2. String Columns with unique.Handle

**Reduce memory by 90% for low-cardinality columns!**

```go
package main

import (
    "fmt"
    "runtime"
    "unique"
)

// Old way: Each row stores full string (wasteful!)
type OldPersonTable struct {
    names []string
    cities []string
}

// New way: Intern strings with unique.Handle
type NewPersonTable struct {
    names  []unique.Handle[string]
    cities []unique.Handle[string]
}

func measureMemory(label string, f func()) {
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)

    f()

    runtime.GC()
    runtime.ReadMemStats(&m2)

    fmt.Printf("%s: %d KB\n", label, (m2.Alloc-m1.Alloc)/1024)
}

func main() {
    // Only 5 unique cities repeated 100,000 times
    cities := []string{"NYC", "SF", "LA", "Chicago", "Boston"}

    measureMemory("Old approach", func() {
        table := OldPersonTable{
            cities: make([]string, 100000),
        }
        for i := 0; i < 100000; i++ {
            table.cities[i] = cities[i%5]  // Duplicates!
        }
        _ = table
    })

    measureMemory("New approach (unique)", func() {
        table := NewPersonTable{
            cities: make([]unique.Handle[string], 100000),
        }
        for i := 0; i < 100000; i++ {
            table.cities[i] = unique.Make(cities[i%5])  // Interned!
        }
        _ = table
    })
}
```

**Typical results:**
- Old: ~800 KB (100,000 string headers + duplicated data)
- New: ~800 KB for handles + ~20 bytes for actual strings = ~800 KB total, but...
- **With unique:** Only 5 strings stored once! ~800 KB vs ~80 KB = 90% savings!

### 3. Fast String Operations with unique.Handle

**O(1) equality, fast GROUP BY, fast DISTINCT!**

```go
package main

import (
    "fmt"
    "unique"
)

type Person struct {
    Name unique.Handle[string]
    City unique.Handle[string]
    Age  int
}

type Table struct {
    persons []Person
}

// Count distinct cities (fast with unique!)
func (t *Table) DistinctCities() int {
    seen := make(map[unique.Handle[string]]struct{})
    for _, p := range t.persons {
        seen[p.City] = struct{}{}
    }
    return len(seen)
}

// Group by city (fast with unique!)
func (t *Table) GroupByCity() map[string]int {
    counts := make(map[unique.Handle[string]]int)
    for _, p := range t.persons {
        counts[p.City]++
    }

    // Convert handles back to strings for display
    result := make(map[string]int)
    for handle, count := range counts {
        result[handle.Value()] = count
    }
    return result
}

// Filter by city (O(1) comparison!)
func (t *Table) FilterByCity(city string) []Person {
    cityHandle := unique.Make(city)  // Intern once
    result := make([]Person, 0)

    for _, p := range t.persons {
        if p.City == cityHandle {  // Pointer comparison!
            result = append(result, p)
        }
    }

    return result
}

func main() {
    table := Table{
        persons: []Person{
            {Name: unique.Make("Alice"), City: unique.Make("NYC"), Age: 30},
            {Name: unique.Make("Bob"), City: unique.Make("SF"), Age: 25},
            {Name: unique.Make("Charlie"), City: unique.Make("NYC"), Age: 35},
            {Name: unique.Make("Diana"), City: unique.Make("SF"), Age: 28},
        },
    }

    fmt.Printf("Distinct cities: %d\n", table.DistinctCities())

    fmt.Println("Group by city:")
    for city, count := range table.GroupByCity() {
        fmt.Printf("  %s: %d\n", city, count)
    }

    fmt.Println("People in NYC:")
    for _, p := range table.FilterByCity("NYC") {
        fmt.Printf("  %s\n", p.Name.Value())
    }
}
```

### 4. Columnar Layout

**Store data by column, not by row!**

```go
package main

import (
    "fmt"
    "unique"
)

// Row-oriented (traditional)
type RowStore struct {
    persons []Person
}

type Person struct {
    Name unique.Handle[string]
    City unique.Handle[string]
    Age  int
}

// Column-oriented (better for analytics!)
type ColumnStore struct {
    names  []unique.Handle[string]
    cities []unique.Handle[string]
    ages   []int
}

func (cs *ColumnStore) NumRows() int {
    return len(cs.names)
}

func (cs *ColumnStore) GetPerson(idx int) Person {
    return Person{
        Name: cs.names[idx],
        City: cs.cities[idx],
        Age:  cs.ages[idx],
    }
}

// Analytical query: Average age per city
func (cs *ColumnStore) AvgAgeByCity() map[string]float64 {
    sums := make(map[unique.Handle[string]]int)
    counts := make(map[unique.Handle[string]]int)

    // Only scan cities and ages columns!
    for i := 0; i < len(cs.ages); i++ {
        city := cs.cities[i]
        sums[city] += cs.ages[i]
        counts[city]++
    }

    result := make(map[string]float64)
    for city, sum := range sums {
        result[city.Value()] = float64(sum) / float64(counts[city])
    }

    return result
}

func main() {
    // Build column store
    cs := ColumnStore{
        names:  []unique.Handle[string]{
            unique.Make("Alice"),
            unique.Make("Bob"),
            unique.Make("Charlie"),
        },
        cities: []unique.Handle[string]{
            unique.Make("NYC"),
            unique.Make("SF"),
            unique.Make("NYC"),
        },
        ages: []int{30, 25, 35},
    }

    fmt.Println("Average age by city:")
    for city, avg := range cs.AvgAgeByCity() {
        fmt.Printf("  %s: %.1f\n", city, avg)
    }
}
```

**Key advantage:** Query only scans cities + ages columns, skipping names!

### 5. Compression: Bit-Packing for Integers

**Store small integers in fewer bits!**

```go
package main

import (
    "fmt"
)

// Store ages (0-127) in 7 bits instead of 64 bits!
type BitPackedColumn struct {
    data      []byte  // Packed bits
    numValues int
    bitsPerValue int
}

func NewBitPackedColumn(bitsPerValue int) *BitPackedColumn {
    return &BitPackedColumn{
        data:         make([]byte, 0),
        bitsPerValue: bitsPerValue,
    }
}

func (c *BitPackedColumn) Append(value uint64) {
    bitOffset := c.numValues * c.bitsPerValue
    byteOffset := bitOffset / 8
    bitInByte := bitOffset % 8

    // Extend data if needed
    requiredBytes := (bitOffset + c.bitsPerValue + 7) / 8
    for len(c.data) < requiredBytes {
        c.data = append(c.data, 0)
    }

    // Write bits (simplified - assumes value fits in one byte)
    c.data[byteOffset] |= byte(value << bitInByte)

    c.numValues++
}

func (c *BitPackedColumn) Get(idx int) uint64 {
    bitOffset := idx * c.bitsPerValue
    byteOffset := bitOffset / 8
    bitInByte := bitOffset % 8

    mask := uint64((1 << c.bitsPerValue) - 1)
    value := uint64(c.data[byteOffset]) >> bitInByte
    return value & mask
}

func main() {
    // Store ages (max 127) in 7 bits
    ages := NewBitPackedColumn(7)

    ages.Append(25)
    ages.Append(30)
    ages.Append(35)
    ages.Append(40)

    fmt.Printf("Stored %d values in %d bytes\n", 4, len(ages.data))
    fmt.Printf("vs %d bytes with uint64\n", 4*8)

    for i := 0; i < 4; i++ {
        fmt.Printf("Age %d: %d\n", i, ages.Get(i))
    }
}
```

**Savings:** 7 bits vs 64 bits = 89% reduction!

### 6. NULL Handling with Validity Bitmaps

**Use 1 bit per value to track NULLs!**

```go
package main

import (
    "fmt"
)

type NullableColumn struct {
    values   []int
    validity []byte  // Bitmap: 1 = valid, 0 = NULL
}

func NewNullableColumn() *NullableColumn {
    return &NullableColumn{
        values:   make([]int, 0),
        validity: make([]byte, 0),
    }
}

func (c *NullableColumn) Append(value *int) {
    idx := len(c.values)

    if value != nil {
        c.values = append(c.values, *value)
        c.setBit(idx, true)
    } else {
        c.values = append(c.values, 0)  // Placeholder
        c.setBit(idx, false)
    }
}

func (c *NullableColumn) Get(idx int) *int {
    if !c.isValid(idx) {
        return nil
    }
    value := c.values[idx]
    return &value
}

func (c *NullableColumn) setBit(idx int, value bool) {
    byteIdx := idx / 8
    bitIdx := idx % 8

    // Extend validity bitmap if needed
    for len(c.validity) <= byteIdx {
        c.validity = append(c.validity, 0)
    }

    if value {
        c.validity[byteIdx] |= 1 << bitIdx
    } else {
        c.validity[byteIdx] &^= 1 << bitIdx
    }
}

func (c *NullableColumn) isValid(idx int) bool {
    byteIdx := idx / 8
    bitIdx := idx % 8

    if byteIdx >= len(c.validity) {
        return false
    }

    return (c.validity[byteIdx] & (1 << bitIdx)) != 0
}

func main() {
    col := NewNullableColumn()

    val1 := 10
    val2 := 20

    col.Append(&val1)
    col.Append(nil)  // NULL
    col.Append(&val2)
    col.Append(nil)  // NULL

    for i := 0; i < 4; i++ {
        if val := col.Get(i); val != nil {
            fmt.Printf("Row %d: %d\n", i, *val)
        } else {
            fmt.Printf("Row %d: NULL\n", i)
        }
    }

    fmt.Printf("\nValidity bitmap uses %d bytes for 4 values\n", len(col.validity))
}
```

### 7. Memory Profiling with pprof

**Measure memory usage of your columns!**

```go
package main

import (
    "fmt"
    "os"
    "runtime"
    "runtime/pprof"
    "unique"
)

func buildTable(size int) []unique.Handle[string] {
    cities := []string{"NYC", "SF", "LA", "Chicago", "Boston"}
    column := make([]unique.Handle[string], size)

    for i := 0; i < size; i++ {
        column[i] = unique.Make(cities[i%len(cities)])
    }

    return column
}

func main() {
    // Start memory profile
    f, _ := os.Create("mem.prof")
    defer f.Close()

    runtime.GC()  // Clean start

    // Build large table
    table := buildTable(1000000)

    runtime.GC()  // Force GC before profile

    pprof.WriteHeapProfile(f)

    fmt.Printf("Built table with %d rows\n", len(table))
    fmt.Println("Memory profile written to mem.prof")
    fmt.Println("Analyze with: go tool pprof mem.prof")

    _ = table
}
```

Analyze with:
```bash
go run main.go
go tool pprof -http=:8080 mem.prof
```

## Pre-Implementation Exercises

### Exercise 1: String Column with unique.Handle

```go
package main

import (
    "unique"
)

type StringColumn struct {
    values []unique.Handle[string]
}

func NewStringColumn() *StringColumn {
    return &StringColumn{
        values: make([]unique.Handle[string], 0),
    }
}

func (c *StringColumn) Append(value string) {
    // TODO: Intern and append
}

func (c *StringColumn) Get(idx int) string {
    // TODO: Return value
    return ""
}

func (c *StringColumn) Distinct() int {
    // TODO: Count distinct values (fast with unique!)
    return 0
}

func (c *StringColumn) Filter(predicate string) []int {
    // TODO: Return indices where value == predicate
    return nil
}

func main() {
    col := NewStringColumn()
    col.Append("Alice")
    col.Append("Bob")
    col.Append("Alice")
    col.Append("Charlie")

    println("Distinct:", col.Distinct())  // Should be 3
    println("Rows with 'Alice':", len(col.Filter("Alice")))  // Should be 2
}
```

### Exercise 2: Columnar Table

```go
package main

import (
    "unique"
)

type Table struct {
    names  []unique.Handle[string]
    ages   []int
    cities []unique.Handle[string]
}

func NewTable() *Table {
    // TODO: Initialize
    return nil
}

func (t *Table) AddRow(name string, age int, city string) {
    // TODO: Append to each column
}

func (t *Table) NumRows() int {
    // TODO: Return number of rows
    return 0
}

func (t *Table) GetName(idx int) string {
    // TODO: Return name at index
    return ""
}

func (t *Table) AverageAge() float64 {
    // TODO: Compute average age
    return 0
}

func (t *Table) FilterByCity(city string) []int {
    // TODO: Return indices of rows with given city
    return nil
}

func main() {
    table := NewTable()
    table.AddRow("Alice", 30, "NYC")
    table.AddRow("Bob", 25, "SF")
    table.AddRow("Charlie", 35, "NYC")

    println("Average age:", table.AverageAge())
    println("Rows in NYC:", len(table.FilterByCity("NYC")))
}
```

### Exercise 3: Bit-Packed Column

```go
package main

type BitPackedColumn struct {
    data         []byte
    numValues    int
    bitsPerValue int
}

func NewBitPackedColumn(bitsPerValue int) *BitPackedColumn {
    // TODO: Initialize
    return nil
}

func (c *BitPackedColumn) Append(value uint64) {
    // TODO: Pack value into bits
}

func (c *BitPackedColumn) Get(idx int) uint64 {
    // TODO: Unpack value from bits
    return 0
}

func (c *BitPackedColumn) BytesUsed() int {
    // TODO: Return bytes used
    return len(c.data)
}

func main() {
    // Store ages (0-127) in 7 bits
    ages := NewBitPackedColumn(7)

    for i := 0; i < 100; i++ {
        ages.Append(uint64(20 + i%50))
    }

    println("Stored 100 ages in", ages.BytesUsed(), "bytes")
    println("vs", 100*8, "bytes with uint64")
    println("First age:", ages.Get(0))
    println("Last age:", ages.Get(99))
}
```

### Exercise 4: Nullable Column

```go
package main

type NullableIntColumn struct {
    values   []int
    validity []byte
}

func NewNullableIntColumn() *NullableIntColumn {
    // TODO: Initialize
    return nil
}

func (c *NullableIntColumn) Append(value *int) {
    // TODO: Append value or NULL
}

func (c *NullableIntColumn) Get(idx int) *int {
    // TODO: Return value or nil
    return nil
}

func (c *NullableIntColumn) Sum() int {
    // TODO: Sum non-NULL values
    return 0
}

func (c *NullableIntColumn) CountNulls() int {
    // TODO: Count NULL values
    return 0
}

func main() {
    col := NewNullableIntColumn()

    val1 := 10
    val2 := 20

    col.Append(&val1)
    col.Append(nil)
    col.Append(&val2)
    col.Append(nil)

    println("Sum:", col.Sum())  // 30
    println("NULLs:", col.CountNulls())  // 2
}
```

### Exercise 5: Benchmark String Storage

```go
package main

import (
    "testing"
    "unique"
)

func BenchmarkRegularStrings(b *testing.B) {
    cities := []string{"NYC", "SF", "LA", "Chicago", "Boston"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        column := make([]string, 10000)
        for j := 0; j < 10000; j++ {
            column[j] = cities[j%5]
        }
        _ = column
    }
}

func BenchmarkUniqueHandles(b *testing.B) {
    cities := []string{"NYC", "SF", "LA", "Chicago", "Boston"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        column := make([]unique.Handle[string], 10000)
        for j := 0; j < 10000; j++ {
            column[j] = unique.Make(cities[j%5])
        }
        _ = column
    }
}
```

Run with:
```bash
go test -bench=. -benchmem
```

## Performance Benchmarks

### Benchmark 1: Memory Usage

Measure memory for 1M rows with low-cardinality strings:

```go
func benchmarkMemory(useUnique bool) {
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)

    if useUnique {
        // unique.Handle approach
    } else {
        // Regular string approach
    }

    runtime.GC()
    runtime.ReadMemStats(&m2)
    fmt.Printf("Memory: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
}
```

### Benchmark 2: String Comparison Speed

```go
func BenchmarkStringEquals(b *testing.B) {
    s1 := "New York City"
    s2 := "New York City"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = s1 == s2  // Byte-by-byte comparison
    }
}

func BenchmarkHandleEquals(b *testing.B) {
    h1 := unique.Make("New York City")
    h2 := unique.Make("New York City")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = h1 == h2  // Pointer comparison (faster!)
    }
}
```

## Common Gotchas to Avoid

### Gotcha 1: Using unique for High-Cardinality Data

```go
// WRONG: unique doesn't help with unique values!
for i := 0; i < 1000000; i++ {
    col.Append(unique.Make(fmt.Sprintf("unique-%d", i)))
}
// No memory savings, only overhead!

// RIGHT: Use unique for low-cardinality data
cities := []string{"NYC", "SF", "LA"}
for i := 0; i < 1000000; i++ {
    col.Append(unique.Make(cities[i%3]))
}
// Huge memory savings!
```

### Gotcha 2: Forgetting Column Alignment

```go
// WRONG: Columns have different lengths!
table.names = append(table.names, unique.Make("Alice"))
table.ages = append(table.ages, 30)
// Forgot cities! Columns misaligned!

// RIGHT: Append to all columns together
func (t *Table) AddRow(name string, age int, city string) {
    t.names = append(t.names, unique.Make(name))
    t.ages = append(t.ages, age)
    t.cities = append(t.cities, unique.Make(city))
}
```

### Gotcha 3: Bit-Packing Overflow

```go
// WRONG: Value doesn't fit in bits!
col := NewBitPackedColumn(7)  // Max value: 127
col.Append(200)  // OVERFLOW!

// RIGHT: Check value range
func (c *BitPackedColumn) Append(value uint64) error {
    maxValue := uint64(1<<c.bitsPerValue) - 1
    if value > maxValue {
        return fmt.Errorf("value %d exceeds max %d", value, maxValue)
    }
    // ... append ...
}
```

### Gotcha 4: Not Checking NULL

```go
// WRONG: Dereferencing NULL!
val := col.Get(0)
fmt.Println(*val)  // PANIC if NULL!

// RIGHT: Check for NULL
if val := col.Get(0); val != nil {
    fmt.Println(*val)
} else {
    fmt.Println("NULL")
}
```

## Checklist Before Starting Lesson 2.2

- [ ] I understand Go 1.23 `unique.Make()` and `unique.Handle[T]`
- [ ] I know when to use unique (low-cardinality data)
- [ ] I understand columnar vs row-oriented storage
- [ ] I can implement bit-packing for integer compression
- [ ] I know how to use validity bitmaps for NULLs
- [ ] I've profiled memory usage with pprof
- [ ] I can measure memory savings from interning
- [ ] I understand cache benefits of columnar layout
- [ ] I know how to handle column alignment
- [ ] I can benchmark string comparison performance

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 2.2

You'll implement:
- Columnar storage for node/edge properties
- String interning with unique.Handle
- Bit-packed integer columns
- NULL handling with validity bitmaps
- Compression benchmarks
- Analytical query performance tests

**Time estimate:** 20-25 hours for full implementation

Good luck! 🚀
