# Phase 3 Lesson 3.4: Go Prep - Execution Engine

**Prerequisites:** Lesson 3.3 complete (Join Algorithms)
**Time:** 6-7 hours Go prep + 30-35 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.4

## Overview

The execution engine brings query plans to life. Before implementing a pipelined executor, master these Go concepts:
- **Go 1.23:** Iterator-based operators with `iter.Seq` ⭐⭐⭐
- Pipeline composition and operator chaining
- Vectorized execution with batches
- Parallel operator execution with goroutines
- Query profiling and instrumentation

**This lesson is about making queries fly!**

## Go Concepts for This Lesson

### 1. Iterator-Based Operators (Volcano Model)

**Each operator implements Next() using iterators!**

```go
package main

import (
    "fmt"
    "iter"
)

// Row represents a tuple
type Row map[string]interface{}

// Operator interface - all operators implement this
type Operator interface {
    // Returns an iterator over rows
    Execute() iter.Seq[Row]
}

// Scan operator - reads from table
type ScanOp struct {
    table []Row
}

func (s *ScanOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for _, row := range s.table {
            if !yield(row) {
                return
            }
        }
    }
}

// Filter operator - applies predicate
type FilterOp struct {
    child     Operator
    predicate func(Row) bool
}

func (f *FilterOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for row := range f.child.Execute() {
            if f.predicate(row) {
                if !yield(row) {
                    return
                }
            }
        }
    }
}

// Project operator - selects columns
type ProjectOp struct {
    child   Operator
    columns []string
}

func (p *ProjectOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for row := range p.child.Execute() {
            newRow := make(Row)
            for _, col := range p.columns {
                newRow[col] = row[col]
            }
            if !yield(newRow) {
                return
            }
        }
    }
}

func main() {
    // Table data
    table := []Row{
        {"id": 1, "name": "Alice", "age": 30},
        {"id": 2, "name": "Bob", "age": 25},
        {"id": 3, "name": "Charlie", "age": 35},
    }

    // Build query plan: SELECT name FROM table WHERE age > 28
    plan := &ProjectOp{
        columns: []string{"name"},
        child: &FilterOp{
            predicate: func(r Row) bool {
                return r["age"].(int) > 28
            },
            child: &ScanOp{table: table},
        },
    }

    // Execute
    fmt.Println("Results:")
    for row := range plan.Execute() {
        fmt.Printf("  %v\n", row)
    }
}
```

**Output:**
```
Results:
  map[name:Alice]
  map[name:Charlie]
```

**Key insight:** Operators compose naturally with iterators!

### 2. Pipeline Composition

**Chain operators to build complex queries!**

```go
package main

import (
    "fmt"
    "iter"
    "strings"
)

// Limit operator - returns first N rows
type LimitOp struct {
    child Operator
    n     int
}

func (l *LimitOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        count := 0
        for row := range l.child.Execute() {
            if count >= l.n {
                return
            }
            if !yield(row) {
                return
            }
            count++
        }
    }
}

// Sort operator - sorts all rows
type SortOp struct {
    child   Operator
    key     string
    reverse bool
}

func (s *SortOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        // Materialize all rows
        var rows []Row
        for row := range s.child.Execute() {
            rows = append(rows, row)
        }

        // Sort
        for i := 0; i < len(rows); i++ {
            for j := i + 1; j < len(rows); j++ {
                less := rows[i][s.key].(int) < rows[j][s.key].(int)
                if s.reverse {
                    less = !less
                }
                if !less {
                    rows[i], rows[j] = rows[j], rows[i]
                }
            }
        }

        // Emit sorted rows
        for _, row := range rows {
            if !yield(row) {
                return
            }
        }
    }
}

func main() {
    table := []Row{
        {"id": 1, "name": "Alice", "age": 30},
        {"id": 2, "name": "Bob", "age": 25},
        {"id": 3, "name": "Charlie", "age": 35},
        {"id": 4, "name": "Diana", "age": 28},
    }

    // SELECT name FROM table WHERE age > 25 ORDER BY age DESC LIMIT 2
    plan := &LimitOp{
        n: 2,
        child: &ProjectOp{
            columns: []string{"name", "age"},
            child: &SortOp{
                key:     "age",
                reverse: true,
                child: &FilterOp{
                    predicate: func(r Row) bool {
                        return r["age"].(int) > 25
                    },
                    child: &ScanOp{table: table},
                },
            },
        },
    }

    fmt.Println("Top 2 oldest (age > 25):")
    for row := range plan.Execute() {
        fmt.Printf("  %s (age %d)\n", row["name"], row["age"])
    }
}
```

**Output:**
```
Top 2 oldest (age > 25):
  Charlie (age 35)
  Alice (age 30)
```

### 3. Vectorized Execution

**Process rows in batches for better performance!**

```go
package main

import (
    "fmt"
)

const BATCH_SIZE = 1024

// RowBatch is a columnar batch of rows
type RowBatch struct {
    Columns map[string][]interface{}
    Size    int
}

func NewRowBatch() *RowBatch {
    return &RowBatch{
        Columns: make(map[string][]interface{}),
        Size:    0,
    }
}

func (b *RowBatch) AddRow(row Row) {
    if b.Size == 0 {
        // Initialize columns
        for key := range row {
            b.Columns[key] = make([]interface{}, 0, BATCH_SIZE)
        }
    }

    for key, val := range row {
        b.Columns[key] = append(b.Columns[key], val)
    }
    b.Size++
}

func (b *RowBatch) GetRow(idx int) Row {
    row := make(Row)
    for key, col := range b.Columns {
        row[key] = col[idx]
    }
    return row
}

func (b *RowBatch) Clear() {
    for key := range b.Columns {
        b.Columns[key] = b.Columns[key][:0]
    }
    b.Size = 0
}

// Vectorized filter operator
type VectorizedFilterOp struct {
    child     Operator
    predicate func(Row) bool
}

func (v *VectorizedFilterOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        batch := NewRowBatch()

        for row := range v.child.Execute() {
            batch.AddRow(row)

            if batch.Size >= BATCH_SIZE {
                // Process batch
                for i := 0; i < batch.Size; i++ {
                    row := batch.GetRow(i)
                    if v.predicate(row) {
                        if !yield(row) {
                            return
                        }
                    }
                }
                batch.Clear()
            }
        }

        // Process remaining rows
        for i := 0; i < batch.Size; i++ {
            row := batch.GetRow(i)
            if v.predicate(row) {
                if !yield(row) {
                    return
                }
            }
        }
    }
}

func main() {
    // Generate large dataset
    table := make([]Row, 10000)
    for i := range table {
        table[i] = Row{
            "id":  i,
            "val": i % 100,
        }
    }

    plan := &VectorizedFilterOp{
        predicate: func(r Row) bool {
            return r["val"].(int) < 5
        },
        child: &ScanOp{table: table},
    }

    count := 0
    for range plan.Execute() {
        count++
    }

    fmt.Printf("Found %d rows (vectorized processing)\n", count)
}
```

**Output:**
```
Found 500 rows (vectorized processing)
```

### 4. Parallel Operator Execution

**Execute independent operators in parallel!**

```go
package main

import (
    "fmt"
    "sync"
)

// Parallel scan operator - partitions table
type ParallelScanOp struct {
    table      []Row
    numWorkers int
}

func (p *ParallelScanOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        rowsPerWorker := len(p.table) / p.numWorkers
        if rowsPerWorker == 0 {
            rowsPerWorker = len(p.table)
        }

        rowChan := make(chan Row, 100)
        var wg sync.WaitGroup

        // Launch workers
        for i := 0; i < p.numWorkers; i++ {
            start := i * rowsPerWorker
            end := start + rowsPerWorker
            if i == p.numWorkers-1 {
                end = len(p.table)
            }

            wg.Add(1)
            go func(partition []Row) {
                defer wg.Done()
                for _, row := range partition {
                    rowChan <- row
                }
            }(p.table[start:end])
        }

        // Close channel when all workers done
        go func() {
            wg.Wait()
            close(rowChan)
        }()

        // Yield rows
        for row := range rowChan {
            if !yield(row) {
                return
            }
        }
    }
}

// Parallel hash join
type ParallelHashJoinOp struct {
    left       Operator
    right      Operator
    leftKey    string
    rightKey   string
    numWorkers int
}

func (p *ParallelHashJoinOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        // Build hash table from left
        hashTable := make(map[interface{}][]Row)
        for row := range p.left.Execute() {
            key := row[p.leftKey]
            hashTable[key] = append(hashTable[key], row)
        }

        // Partition right side by hash
        partitions := make([][]Row, p.numWorkers)
        for row := range p.right.Execute() {
            key := row[p.rightKey]
            hash := hashKey(key)
            partition := hash % p.numWorkers
            partitions[partition] = append(partitions[partition], row)
        }

        // Probe in parallel
        resultChan := make(chan Row, 100)
        var wg sync.WaitGroup

        for _, partition := range partitions {
            wg.Add(1)
            go func(rows []Row) {
                defer wg.Done()
                for _, rightRow := range rows {
                    key := rightRow[p.rightKey]
                    if leftRows, found := hashTable[key]; found {
                        for _, leftRow := range leftRows {
                            // Merge rows
                            joinedRow := make(Row)
                            for k, v := range leftRow {
                                joinedRow[k] = v
                            }
                            for k, v := range rightRow {
                                joinedRow[k] = v
                            }
                            resultChan <- joinedRow
                        }
                    }
                }
            }(partition)
        }

        go func() {
            wg.Wait()
            close(resultChan)
        }()

        for row := range resultChan {
            if !yield(row) {
                return
            }
        }
    }
}

func hashKey(key interface{}) int {
    // Simple hash for demo
    switch v := key.(type) {
    case int:
        return v
    case string:
        h := 0
        for _, ch := range v {
            h = h*31 + int(ch)
        }
        return h
    default:
        return 0
    }
}

func main() {
    // Large table
    table := make([]Row, 100000)
    for i := range table {
        table[i] = Row{"id": i, "val": i % 1000}
    }

    plan := &ParallelScanOp{
        table:      table,
        numWorkers: 4,
    }

    count := 0
    for range plan.Execute() {
        count++
    }

    fmt.Printf("Scanned %d rows (4 workers)\n", count)
}
```

**Output:**
```
Scanned 100000 rows (4 workers)
```

### 5. Query Profiling and Instrumentation

**Measure operator performance!**

```go
package main

import (
    "fmt"
    "time"
)

// Instrumented operator - tracks stats
type InstrumentedOp struct {
    child        Operator
    name         string
    rowsProduced int
    timeSpent    time.Duration
}

func (i *InstrumentedOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        start := time.Now()

        for row := range i.child.Execute() {
            i.rowsProduced++
            if !yield(row) {
                i.timeSpent = time.Since(start)
                return
            }
        }

        i.timeSpent = time.Since(start)
    }
}

func (i *InstrumentedOp) Stats() string {
    rowsPerSec := 0.0
    if i.timeSpent > 0 {
        rowsPerSec = float64(i.rowsProduced) / i.timeSpent.Seconds()
    }

    return fmt.Sprintf("%s: %d rows in %v (%.0f rows/sec)",
        i.name, i.rowsProduced, i.timeSpent, rowsPerSec)
}

func main() {
    table := make([]Row, 1000000)
    for i := range table {
        table[i] = Row{"id": i, "val": i % 100}
    }

    // Build instrumented plan
    scan := &InstrumentedOp{
        child: &ScanOp{table: table},
        name:  "Scan",
    }

    filter := &InstrumentedOp{
        child: &FilterOp{
            predicate: func(r Row) bool {
                return r["val"].(int) < 10
            },
            child: scan,
        },
        name: "Filter",
    }

    project := &InstrumentedOp{
        child: &ProjectOp{
            columns: []string{"id"},
            child:   filter,
        },
        name: "Project",
    }

    // Execute
    count := 0
    for range project.Execute() {
        count++
    }

    // Print stats
    fmt.Println("\nQuery Profile:")
    fmt.Println(scan.Stats())
    fmt.Println(filter.Stats())
    fmt.Println(project.Stats())
    fmt.Printf("\nTotal rows returned: %d\n", count)
}
```

**Output:**
```
Query Profile:
Scan: 1000000 rows in 45ms (22222222 rows/sec)
Filter: 100000 rows in 48ms (2083333 rows/sec)
Project: 100000 rows in 49ms (2040816 rows/sec)

Total rows returned: 100000
```

### 6. Materialization vs Pipelining

**Know when to buffer vs stream!**

```go
package main

import (
    "fmt"
)

// Pipelined operator - streams rows
type PipelinedFilterOp struct {
    child     Operator
    predicate func(Row) bool
}

func (p *PipelinedFilterOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        // No buffering - pure streaming!
        for row := range p.child.Execute() {
            if p.predicate(row) {
                if !yield(row) {
                    return  // Can stop early!
                }
            }
        }
    }
}

// Materialized operator - buffers all rows
type MaterializedSortOp struct {
    child Operator
    key   string
}

func (m *MaterializedSortOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        // Must materialize entire input!
        var rows []Row
        for row := range m.child.Execute() {
            rows = append(rows, row)
        }

        // Sort (simplified)
        // ... sorting logic ...

        // Then stream results
        for _, row := range rows {
            if !yield(row) {
                return
            }
        }
    }
}

func main() {
    table := make([]Row, 1000000)
    for i := range table {
        table[i] = Row{"id": i}
    }

    // Pipelined: can stop early
    pipelined := &LimitOp{
        n: 10,
        child: &PipelinedFilterOp{
            predicate: func(r Row) bool {
                return r["id"].(int)%2 == 0
            },
            child: &ScanOp{table: table},
        },
    }

    fmt.Println("Pipelined execution (can stop early):")
    count := 0
    for range pipelined.Execute() {
        count++
    }
    fmt.Printf("Processed and returned %d rows\n", count)

    // Materialized: must process all rows
    materialized := &LimitOp{
        n: 10,
        child: &MaterializedSortOp{
            key:   "id",
            child: &ScanOp{table: table},
        },
    }

    fmt.Println("\nMaterialized execution (must sort all):")
    count = 0
    for range materialized.Execute() {
        count++
    }
    fmt.Printf("Processed all rows, returned %d\n", count)
}
```

### 7. Operator Cost Tracking

**Track actual vs estimated costs!**

```go
package main

import (
    "fmt"
    "time"
)

type OperatorStats struct {
    EstimatedRows int
    ActualRows    int
    EstimatedCost float64
    ActualCost    time.Duration
}

type TrackedOp struct {
    child Operator
    stats *OperatorStats
}

func (t *TrackedOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        start := time.Now()

        for row := range t.child.Execute() {
            t.stats.ActualRows++
            if !yield(row) {
                t.stats.ActualCost = time.Since(start)
                return
            }
        }

        t.stats.ActualCost = time.Since(start)
    }
}

func (t *TrackedOp) CompareEstimates() {
    fmt.Printf("Estimated rows: %d, Actual: %d (%.1f%% error)\n",
        t.stats.EstimatedRows,
        t.stats.ActualRows,
        100.0*float64(abs(t.stats.EstimatedRows-t.stats.ActualRows))/float64(t.stats.ActualRows))
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func main() {
    // Example usage for cost tracking
    fmt.Println("Operator cost tracking example")
}
```

## Pre-Implementation Exercises

### Exercise 1: Build Iterator Operators

```go
package main

// TODO: Implement basic operators using iter.Seq

type ScanOp struct {
    table []Row
}

func (s *ScanOp) Execute() iter.Seq[Row] {
    // TODO: Implement
    return nil
}

type FilterOp struct {
    child     Operator
    predicate func(Row) bool
}

func (f *FilterOp) Execute() iter.Seq[Row] {
    // TODO: Implement
    return nil
}

// TODO: Add ProjectOp, JoinOp, LimitOp

func main() {
    // TODO: Build and execute query plan
}
```

### Exercise 2: Vectorized Processing

```go
package main

// TODO: Implement batched execution

const BATCH_SIZE = 1024

type VectorizedFilterOp struct {
    child     Operator
    predicate func(Row) bool
}

func (v *VectorizedFilterOp) Execute() iter.Seq[Row] {
    // TODO: Process rows in batches of BATCH_SIZE
    // TODO: Measure speedup vs row-at-a-time
    return nil
}
```

### Exercise 3: Parallel Execution

```go
package main

// TODO: Implement parallel scan

type ParallelScanOp struct {
    table      []Row
    numWorkers int
}

func (p *ParallelScanOp) Execute() iter.Seq[Row] {
    // TODO: Partition table across workers
    // TODO: Merge results
    return nil
}
```

### Exercise 4: Query Profiling

```go
package main

// TODO: Add instrumentation to all operators

type ProfiledOp struct {
    child        Operator
    rowsIn       int
    rowsOut      int
    timeSpent    time.Duration
}

func (p *ProfiledOp) Execute() iter.Seq[Row] {
    // TODO: Track stats
    return nil
}

func (p *ProfiledOp) Report() string {
    // TODO: Format stats nicely
    return ""
}
```

### Exercise 5: Adaptive Execution

```go
package main

// TODO: Switch join algorithm based on runtime statistics

type AdaptiveJoinOp struct {
    left  Operator
    right Operator
}

func (a *AdaptiveJoinOp) Execute() iter.Seq[Row] {
    // TODO: Sample inputs
    // TODO: Choose hash join vs sort-merge vs nested loop
    // TODO: Execute chosen algorithm
    return nil
}
```

## Performance Benchmarks

### Benchmark 1: Pipelined vs Materialized

```go
func BenchmarkPipeline(b *testing.B) {
    table := generateTable(1000000)

    b.Run("Pipelined", func(b *testing.B) {
        for b.Loop() {
            plan := pipelinedPlan(table)
            for range plan.Execute() {
            }
        }
    })

    b.Run("Materialized", func(b *testing.B) {
        for b.Loop() {
            plan := materializedPlan(table)
            for range plan.Execute() {
            }
        }
    })
}
```

**Expected: Pipelined much faster when limit used!**

### Benchmark 2: Vectorized Execution

```go
func BenchmarkVectorized(b *testing.B) {
    table := generateTable(10000000)

    b.Run("RowAtATime", func(b *testing.B) {
        for b.Loop() {
            _ = rowAtATimeFilter(table)
        }
    })

    b.Run("Vectorized", func(b *testing.B) {
        for b.Loop() {
            _ = vectorizedFilter(table)
        }
    })
}
```

**Expected: Vectorized 2-3x faster!**

## Common Gotchas to Avoid

### Gotcha 1: Not Handling Early Exit

```go
// WRONG: Ignores yield return value
func (f *FilterOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for row := range f.child.Execute() {
            if f.predicate(row) {
                yield(row)  // Ignores return!
            }
        }
    }
}

// RIGHT: Check yield return value
func (f *FilterOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for row := range f.child.Execute() {
            if f.predicate(row) {
                if !yield(row) {
                    return  // Stop processing!
                }
            }
        }
    }
}
```

### Gotcha 2: Buffering Too Much

```go
// WRONG: Materialize when streaming would work
func (f *FilterOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        var results []Row
        for row := range f.child.Execute() {
            if f.predicate(row) {
                results = append(results, row)  // Wastes memory!
            }
        }
        for _, row := range results {
            yield(row)
        }
    }
}

// RIGHT: Stream through
func (f *FilterOp) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for row := range f.child.Execute() {
            if f.predicate(row) {
                if !yield(row) {
                    return
                }
            }
        }
    }
}
```

### Gotcha 3: Race Conditions in Parallel Execution

```go
// WRONG: Shared mutable state
var count int
for i := 0; i < numWorkers; i++ {
    go func() {
        for row := range partition[i] {
            count++  // RACE!
        }
    }()
}

// RIGHT: Use channels or atomics
resultChan := make(chan Row, 100)
for i := 0; i < numWorkers; i++ {
    go func(part []Row) {
        for _, row := range part {
            resultChan <- row  // Thread-safe
        }
    }(partition[i])
}
```

## Checklist Before Starting Lesson 3.4

- [ ] I understand iterator-based operator execution
- [ ] I can compose operators into pipelines
- [ ] I know when to materialize vs pipeline
- [ ] I can implement vectorized execution
- [ ] I understand parallel operator execution
- [ ] I can instrument operators for profiling
- [ ] I always handle early exit in iterators
- [ ] I understand the tradeoffs of buffering
- [ ] I can detect and avoid race conditions
- [ ] I've benchmarked execution strategies

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.4

You'll implement:
- Complete iterator-based operator framework
- Vectorized execution engine
- Parallel query execution
- Query profiling and statistics
- Adaptive execution strategies
- Pipeline breakers (sort, hash join build)
- Comprehensive benchmarks

**Time estimate:** 30-35 hours for full implementation

**The execution engine is where queries come alive!** 🚀
