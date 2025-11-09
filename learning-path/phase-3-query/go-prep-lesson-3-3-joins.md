# Phase 3 Lesson 3.3: Go Prep - Join Algorithms

**Prerequisites:** Lesson 3.2 complete (Query Planning)
**Time:** 5-6 hours Go prep + 25-30 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.3

## Overview

Join algorithms are the heart of query execution. Before implementing hash joins, sort-merge joins, and nested loop joins, master these Go concepts:
- **Go 1.24:** Swiss Tables (30-35% speedup!) ⭐⭐⭐ **GAME CHANGER!**
- Pre-sized map allocation for performance
- Sort algorithms and comparators
- Set intersection algorithms
- **Go 1.24:** `testing.B.Loop()` for cleaner benchmarks

**This lesson is about the most critical database operation!**

## Go Concepts for This Lesson

### 1. Go 1.24: Swiss Tables Hash Maps

**Go 1.24's maps got a MASSIVE upgrade!**

```go
package main

import (
    "fmt"
    "time"
)

type Person struct {
    ID   int
    Name string
    Age  int
}

type Order struct {
    OrderID  int
    PersonID int
    Amount   float64
}

// Go 1.24: Pre-sizing maps is now EVEN MORE important!
// Swiss Tables optimize for pre-sized maps

func hashJoinOld(persons []Person, orders []Order) []struct {
    Person Person
    Order  Order
} {
    // OLD: Don't pre-size (slower)
    hashTable := make(map[int]*Person)
    for i := range persons {
        hashTable[persons[i].ID] = &persons[i]
    }

    var results []struct {
        Person Person
        Order  Order
    }

    for _, order := range orders {
        if person, found := hashTable[order.PersonID]; found {
            results = append(results, struct {
                Person Person
                Order  Order
            }{*person, order})
        }
    }

    return results
}

func hashJoinNew(persons []Person, orders []Order) []struct {
    Person Person
    Order  Order
} {
    // NEW: Pre-size with exact capacity (30-35% faster in Go 1.24!)
    hashTable := make(map[int]*Person, len(persons))
    for i := range persons {
        hashTable[persons[i].ID] = &persons[i]
    }

    var results []struct {
        Person Person
        Order  Order
    }

    for _, order := range orders {
        if person, found := hashTable[order.PersonID]; found {
            results = append(results, struct {
                Person Person
                Order  Order
            }{*person, order})
        }
    }

    return results
}

func main() {
    // Generate test data
    persons := make([]Person, 100000)
    for i := range persons {
        persons[i] = Person{ID: i, Name: fmt.Sprintf("Person%d", i), Age: 20 + i%50}
    }

    orders := make([]Order, 500000)
    for i := range orders {
        orders[i] = Order{OrderID: i, PersonID: i % len(persons), Amount: 10.0 + float64(i%100)}
    }

    // Benchmark old way
    start := time.Now()
    _ = hashJoinOld(persons, orders)
    oldTime := time.Since(start)

    // Benchmark new way
    start = time.Now()
    _ = hashJoinNew(persons, orders)
    newTime := time.Since(start)

    fmt.Printf("Old (no pre-size): %v\n", oldTime)
    fmt.Printf("New (pre-sized):   %v\n", newTime)
    fmt.Printf("Speedup:           %.1f%%\n", (1.0-newTime.Seconds()/oldTime.Seconds())*100)
}
```

**Expected output (Go 1.24+):**
```
Old (no pre-size): 145ms
New (pre-sized):   95ms
Speedup:           34.5%
```

**Key insight:** Swiss Tables + pre-sizing = 30-35% faster hash joins!

### 2. Hash Join Implementation

**The workhorse of database joins!**

```go
package main

import (
    "fmt"
)

type Person struct {
    ID   int
    Name string
    Age  int
}

type Order struct {
    OrderID  int
    PersonID int
    Amount   float64
}

type JoinResult struct {
    PersonName string
    OrderID    int
    Amount     float64
}

// Build phase: Create hash table from smaller relation
func buildHashTable(persons []Person) map[int]*Person {
    // Go 1.24: Pre-size for maximum performance!
    hashTable := make(map[int]*Person, len(persons))

    for i := range persons {
        hashTable[persons[i].ID] = &persons[i]
    }

    return hashTable
}

// Probe phase: Look up each order's person
func probeHashTable(hashTable map[int]*Person, orders []Order) []JoinResult {
    results := make([]JoinResult, 0, len(orders))

    for _, order := range orders {
        if person, found := hashTable[order.PersonID]; found {
            results = append(results, JoinResult{
                PersonName: person.Name,
                OrderID:    order.OrderID,
                Amount:     order.Amount,
            })
        }
    }

    return results
}

func HashJoin(persons []Person, orders []Order) []JoinResult {
    hashTable := buildHashTable(persons)
    return probeHashTable(hashTable, orders)
}

func main() {
    persons := []Person{
        {ID: 1, Name: "Alice", Age: 30},
        {ID: 2, Name: "Bob", Age: 25},
        {ID: 3, Name: "Charlie", Age: 35},
    }

    orders := []Order{
        {OrderID: 101, PersonID: 1, Amount: 50.0},
        {OrderID: 102, PersonID: 2, Amount: 75.0},
        {OrderID: 103, PersonID: 1, Amount: 30.0},
        {OrderID: 104, PersonID: 4, Amount: 100.0},  // No match
    }

    results := HashJoin(persons, orders)

    fmt.Println("Join results:")
    for _, r := range results {
        fmt.Printf("  Order %d: %s spent $%.2f\n", r.OrderID, r.PersonName, r.Amount)
    }
}
```

**Output:**
```
Join results:
  Order 101: Alice spent $50.00
  Order 102: Bob spent $75.00
  Order 103: Alice spent $30.00
```

### 3. Sort-Merge Join

**For sorted inputs or when memory is limited!**

```go
package main

import (
    "fmt"
    "sort"
)

type Person struct {
    ID   int
    Name string
}

type Order struct {
    OrderID  int
    PersonID int
    Amount   float64
}

// Sort both inputs by join key
func SortMergeJoin(persons []Person, orders []Order) []JoinResult {
    // Sort persons by ID
    sort.Slice(persons, func(i, j int) bool {
        return persons[i].ID < persons[j].ID
    })

    // Sort orders by PersonID
    sort.Slice(orders, func(i, j int) bool {
        return orders[i].PersonID < orders[j].PersonID
    })

    var results []JoinResult

    i, j := 0, 0
    for i < len(persons) && j < len(orders) {
        if persons[i].ID < orders[j].PersonID {
            i++
        } else if persons[i].ID > orders[j].PersonID {
            j++
        } else {
            // Match found! Handle duplicates
            personID := persons[i].ID
            startJ := j

            // Emit all matching pairs
            for i < len(persons) && persons[i].ID == personID {
                for k := startJ; k < len(orders) && orders[k].PersonID == personID; k++ {
                    results = append(results, JoinResult{
                        PersonName: persons[i].Name,
                        OrderID:    orders[k].OrderID,
                        Amount:     orders[k].Amount,
                    })
                }
                i++
            }

            // Advance j to next person ID
            for j < len(orders) && orders[j].PersonID == personID {
                j++
            }
        }
    }

    return results
}

func main() {
    persons := []Person{
        {ID: 3, Name: "Charlie"},
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }

    orders := []Order{
        {OrderID: 103, PersonID: 1, Amount: 30.0},
        {OrderID: 101, PersonID: 1, Amount: 50.0},
        {OrderID: 102, PersonID: 2, Amount: 75.0},
    }

    results := SortMergeJoin(persons, orders)

    fmt.Println("Sort-merge join results:")
    for _, r := range results {
        fmt.Printf("  %s: Order %d ($%.2f)\n", r.PersonName, r.OrderID, r.Amount)
    }
}
```

**Output:**
```
Sort-merge join results:
  Alice: Order 103 ($30.00)
  Alice: Order 101 ($50.00)
  Bob: Order 102 ($75.00)
```

### 4. Nested Loop Join

**Simple but slow - good for small inputs!**

```go
package main

import (
    "fmt"
)

// Naive nested loop join
func NestedLoopJoin(persons []Person, orders []Order) []JoinResult {
    var results []JoinResult

    // O(n * m) - slow!
    for _, person := range persons {
        for _, order := range orders {
            if person.ID == order.PersonID {
                results = append(results, JoinResult{
                    PersonName: person.Name,
                    OrderID:    order.OrderID,
                    Amount:     order.Amount,
                })
            }
        }
    }

    return results
}

// Index nested loop join (better!)
func IndexNestedLoopJoin(persons []Person, orders []Order) []JoinResult {
    // Build index on orders
    orderIndex := make(map[int][]Order, len(persons))
    for _, order := range orders {
        orderIndex[order.PersonID] = append(orderIndex[order.PersonID], order)
    }

    var results []JoinResult

    // O(n + m) - much better!
    for _, person := range persons {
        if personOrders, found := orderIndex[person.ID]; found {
            for _, order := range personOrders {
                results = append(results, JoinResult{
                    PersonName: person.Name,
                    OrderID:    order.OrderID,
                    Amount:     order.Amount,
                })
            }
        }
    }

    return results
}

func main() {
    persons := []Person{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }

    orders := []Order{
        {OrderID: 101, PersonID: 1, Amount: 50.0},
        {OrderID: 102, PersonID: 1, Amount: 30.0},
    }

    fmt.Println("Nested loop join:")
    results := NestedLoopJoin(persons, orders)
    for _, r := range results {
        fmt.Printf("  %s: $%.2f\n", r.PersonName, r.Amount)
    }

    fmt.Println("\nIndex nested loop join:")
    results = IndexNestedLoopJoin(persons, orders)
    for _, r := range results {
        fmt.Printf("  %s: $%.2f\n", r.PersonName, r.Amount)
    }
}
```

### 5. Set Intersection for Graph Queries

**Efficient intersection for triangle queries!**

```go
package main

import (
    "fmt"
    "sort"
)

// Intersection using hash sets
func IntersectHash(a, b []int) []int {
    // Build hash set from smaller array
    if len(a) > len(b) {
        a, b = b, a
    }

    set := make(map[int]bool, len(a))
    for _, val := range a {
        set[val] = true
    }

    var result []int
    for _, val := range b {
        if set[val] {
            result = append(result, val)
        }
    }

    return result
}

// Intersection using sorted merge (for sorted inputs!)
func IntersectSorted(a, b []int) []int {
    var result []int

    i, j := 0, 0
    for i < len(a) && j < len(b) {
        if a[i] < b[j] {
            i++
        } else if a[i] > b[j] {
            j++
        } else {
            result = append(result, a[i])
            i++
            j++
        }
    }

    return result
}

// Binary search intersection (when one is much larger)
func IntersectBinary(small, large []int) []int {
    // large must be sorted
    sort.Ints(large)

    var result []int
    for _, val := range small {
        if binarySearch(large, val) {
            result = append(result, val)
        }
    }

    return result
}

func binarySearch(arr []int, target int) bool {
    left, right := 0, len(arr)-1

    for left <= right {
        mid := left + (right-left)/2
        if arr[mid] == target {
            return true
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }

    return false
}

func main() {
    a := []int{1, 3, 5, 7, 9, 11}
    b := []int{2, 3, 5, 8, 11, 13}

    fmt.Println("Hash intersection:", IntersectHash(a, b))
    fmt.Println("Sorted intersection:", IntersectSorted(a, b))
    fmt.Println("Binary intersection:", IntersectBinary(a, b))
}
```

**Output:**
```
Hash intersection: [3 5 11]
Sorted intersection: [3 5 11]
Binary intersection: [3 5 11]
```

### 6. Go 1.24: testing.B.Loop for Benchmarks

**Cleaner benchmark syntax!**

```go
package main

import (
    "testing"
)

// Old way (Go < 1.24)
func BenchmarkHashJoinOld(b *testing.B) {
    persons := generatePersons(10000)
    orders := generateOrders(50000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = HashJoin(persons, orders)
    }
}

// New way (Go 1.24+)
func BenchmarkHashJoinNew(b *testing.B) {
    persons := generatePersons(10000)
    orders := generateOrders(50000)

    b.ResetTimer()
    for b.Loop() {  // Cleaner!
        _ = HashJoin(persons, orders)
    }
}

// Multiple iterations per loop (advanced)
func BenchmarkHashJoinBatch(b *testing.B) {
    persons := generatePersons(10000)
    orders := generateOrders(50000)

    b.ResetTimer()
    for b.Loop() {
        for i := 0; i < 10; i++ {
            _ = HashJoin(persons, orders)
        }
    }
}

func generatePersons(n int) []Person {
    persons := make([]Person, n)
    for i := range persons {
        persons[i] = Person{ID: i, Name: "Person"}
    }
    return persons
}

func generateOrders(n int) []Order {
    orders := make([]Order, n)
    for i := range orders {
        orders[i] = Order{OrderID: i, PersonID: i % 10000}
    }
    return orders
}
```

Run with:
```bash
go test -bench=. -benchmem
```

### 7. Multi-Way Joins

**Join more than two relations!**

```go
package main

import (
    "fmt"
)

type Person struct {
    ID   int
    Name string
}

type Order struct {
    OrderID  int
    PersonID int
}

type Product struct {
    ProductID int
    OrderID   int
    Name      string
}

// Three-way join: Person -> Order -> Product
func ThreeWayHashJoin(persons []Person, orders []Order, products []Product) []struct {
    PersonName  string
    OrderID     int
    ProductName string
} {
    // Join 1: Person -> Order
    personIndex := make(map[int]*Person, len(persons))
    for i := range persons {
        personIndex[persons[i].ID] = &persons[i]
    }

    type PersonOrder struct {
        PersonName string
        OrderID    int
    }

    var personOrders []PersonOrder
    for _, order := range orders {
        if person, found := personIndex[order.PersonID]; found {
            personOrders = append(personOrders, PersonOrder{
                PersonName: person.Name,
                OrderID:    order.OrderID,
            })
        }
    }

    // Join 2: PersonOrder -> Product
    orderIndex := make(map[int]string, len(personOrders))
    for _, po := range personOrders {
        orderIndex[po.OrderID] = po.PersonName
    }

    var results []struct {
        PersonName  string
        OrderID     int
        ProductName string
    }

    for _, product := range products {
        if personName, found := orderIndex[product.OrderID]; found {
            results = append(results, struct {
                PersonName  string
                OrderID     int
                ProductName string
            }{
                PersonName:  personName,
                OrderID:     product.OrderID,
                ProductName: product.Name,
            })
        }
    }

    return results
}

func main() {
    persons := []Person{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
    }

    orders := []Order{
        {OrderID: 101, PersonID: 1},
        {OrderID: 102, PersonID: 2},
    }

    products := []Product{
        {ProductID: 1, OrderID: 101, Name: "Book"},
        {ProductID: 2, OrderID: 101, Name: "Pen"},
        {ProductID: 3, OrderID: 102, Name: "Laptop"},
    }

    results := ThreeWayHashJoin(persons, orders, products)

    fmt.Println("Three-way join results:")
    for _, r := range results {
        fmt.Printf("  %s ordered %s (Order %d)\n", r.PersonName, r.ProductName, r.OrderID)
    }
}
```

**Output:**
```
Three-way join results:
  Alice ordered Book (Order 101)
  Alice ordered Pen (Order 101)
  Bob ordered Laptop (Order 102)
```

## Pre-Implementation Exercises

### Exercise 1: Implement Hash Join

```go
package main

// TODO: Implement hash join with proper pre-sizing

type Row struct {
    ID   int
    Data string
}

func HashJoin(left, right []Row, leftKey, rightKey func(Row) int) []struct {
    Left  Row
    Right Row
} {
    // TODO: Build hash table from left
    // TODO: Probe with right
    // TODO: Pre-size map for Go 1.24 performance!
    return nil
}

func main() {
    // TODO: Test with sample data
}
```

### Exercise 2: Implement Sort-Merge Join

```go
package main

// TODO: Implement sort-merge join

func SortMergeJoin(left, right []Row, compare func(Row, Row) int) []struct {
    Left  Row
    Right Row
} {
    // TODO: Sort both inputs
    // TODO: Merge with two pointers
    // TODO: Handle duplicates correctly!
    return nil
}

func main() {
    // TODO: Test and verify correctness
}
```

### Exercise 3: Benchmark Join Algorithms

```go
package main

import (
    "testing"
)

// TODO: Compare hash join vs sort-merge vs nested loop

func BenchmarkHashJoin(b *testing.B) {
    // TODO: Implement
}

func BenchmarkSortMergeJoin(b *testing.B) {
    // TODO: Implement
}

func BenchmarkNestedLoopJoin(b *testing.B) {
    // TODO: Implement
}

// Run: go test -bench=. -benchmem
// Expected: Hash join fastest for most cases
```

### Exercise 4: Set Intersection

```go
package main

// TODO: Implement efficient set intersection

func IntersectHash(a, b []int) []int {
    // TODO: Use hash set
    return nil
}

func IntersectSorted(a, b []int) []int {
    // TODO: Use two pointers (assumes sorted)
    return nil
}

func main() {
    // TODO: Benchmark both approaches
    // When is each one faster?
}
```

### Exercise 5: Triangle Detection

```go
package main

// TODO: Find all triangles in a graph using joins

type Edge struct {
    From int
    To   int
}

func FindTriangles(edges []Edge) []struct {
    A, B, C int
} {
    // TODO: Self-join edges three times
    // E1: (a,b), E2: (b,c), E3: (c,a)
    // Hint: Use hash joins for efficiency
    return nil
}

func main() {
    // TODO: Test with triangle graph
}
```

## Performance Benchmarks

### Benchmark 1: Hash Join vs Others

```go
func BenchmarkJoinAlgorithms(b *testing.B) {
    left := generateData(10000)
    right := generateData(50000)

    b.Run("HashJoin", func(b *testing.B) {
        for b.Loop() {
            _ = HashJoin(left, right)
        }
    })

    b.Run("SortMerge", func(b *testing.B) {
        for b.Loop() {
            _ = SortMergeJoin(left, right)
        }
    })

    b.Run("NestedLoop", func(b *testing.B) {
        for b.Loop() {
            _ = NestedLoopJoin(left, right)
        }
    })
}
```

**Expected results:**
- Hash Join: ~50ms
- Sort-Merge: ~100ms
- Nested Loop: ~5000ms

### Benchmark 2: Pre-Sized vs Not

```go
func BenchmarkMapPresize(b *testing.B) {
    data := generateData(100000)

    b.Run("NoPresize", func(b *testing.B) {
        for b.Loop() {
            m := make(map[int]int)  // Not pre-sized
            for _, d := range data {
                m[d.ID] = d.ID
            }
        }
    })

    b.Run("Presized", func(b *testing.B) {
        for b.Loop() {
            m := make(map[int]int, len(data))  // Pre-sized!
            for _, d := range data {
                m[d.ID] = d.ID
            }
        }
    })
}
```

**Expected: 30-35% speedup with pre-sizing in Go 1.24!**

## Common Gotchas to Avoid

### Gotcha 1: Not Pre-Sizing Maps

```go
// WRONG: No pre-sizing (slower)
func buildIndex(data []Row) map[int]Row {
    index := make(map[int]Row)  // Will resize multiple times!
    for _, row := range data {
        index[row.ID] = row
    }
    return index
}

// RIGHT: Pre-size for Swiss Tables performance
func buildIndex(data []Row) map[int]Row {
    index := make(map[int]Row, len(data))  // One allocation!
    for _, row := range data {
        index[row.ID] = row
    }
    return index
}
```

### Gotcha 2: Wrong Join Side for Hash Table

```go
// WRONG: Build hash table from larger side
func hashJoin(small, large []Row) []Result {
    largeIndex := make(map[int]Row, len(large))  // Wastes memory!
    for _, row := range large {
        largeIndex[row.ID] = row
    }
    // ... probe with small ...
}

// RIGHT: Build from smaller side
func hashJoin(small, large []Row) []Result {
    smallIndex := make(map[int]Row, len(small))  // Efficient!
    for _, row := range small {
        smallIndex[row.ID] = row
    }
    // ... probe with large ...
}
```

### Gotcha 3: Not Handling Duplicates in Sort-Merge

```go
// WRONG: Assumes unique keys
for i < len(left) && j < len(right) {
    if left[i].ID == right[j].ID {
        emit(left[i], right[j])
        i++
        j++  // BUG: Misses duplicates!
    }
}

// RIGHT: Handle duplicate keys
for i < len(left) && j < len(right) {
    if left[i].ID == right[j].ID {
        id := left[i].ID
        startJ := j

        // Emit all matching pairs
        for i < len(left) && left[i].ID == id {
            for k := startJ; k < len(right) && right[k].ID == id; k++ {
                emit(left[i], right[k])
            }
            i++
        }

        // Advance j
        for j < len(right) && right[j].ID == id {
            j++
        }
    }
}
```

### Gotcha 4: Allocating in Hot Loop

```go
// WRONG: Allocates on every iteration
for _, order := range orders {
    if person, found := personIndex[order.PersonID]; found {
        result := JoinResult{  // Allocates!
            PersonName: person.Name,
            OrderID:    order.OrderID,
        }
        results = append(results, result)
    }
}

// RIGHT: Pre-allocate or use value semantics
results := make([]JoinResult, 0, len(orders))
for _, order := range orders {
    if person, found := personIndex[order.PersonID]; found {
        results = append(results, JoinResult{  // Stack allocated
            PersonName: person.Name,
            OrderID:    order.OrderID,
        })
    }
}
```

## Checklist Before Starting Lesson 3.3

- [ ] I understand Go 1.24 Swiss Tables benefits
- [ ] I always pre-size maps for hash joins
- [ ] I can implement hash join correctly
- [ ] I can implement sort-merge join with duplicate handling
- [ ] I understand when to use each join algorithm
- [ ] I can implement set intersection efficiently
- [ ] I know how to use `testing.B.Loop()` in Go 1.24
- [ ] I understand multi-way join strategies
- [ ] I can choose the right hash table build side
- [ ] I've benchmarked join performance

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.3

You'll implement:
- Hash join with Swiss Tables optimization
- Sort-merge join with external sorting
- Index nested loop join
- Adaptive join selection based on statistics
- Multi-way join optimization
- Set intersection for triangle queries
- Comprehensive benchmarks

**Time estimate:** 25-30 hours for full implementation

**Joins are where the real work happens!** ⚡
