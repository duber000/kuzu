# Pre-work Week 3-4: Intermediate Concepts

**Duration:** 2 weeks | **Time Commitment:** 15-20 hours/week | **Difficulty:** Intermediate

## Overview

You've learned the basics. Now we go deeper into concepts that are **critical** for building a database:
- Pointers and memory (databases care about memory layout)
- Testing and benchmarking (databases need correctness and speed)
- JSON encoding/decoding (data persistence)
- File I/O (databases live on disk)

By the end of Week 4, you'll build a **persistent key-value store** that survives crashes!

## Week 3: Deep Concepts

### Day 1-2: Pointers and Memory

Understanding pointers is **essential** for database programming.

#### What is a Pointer?

```go
package main

import "fmt"

func main() {
    // Regular variable
    x := 42
    fmt.Println("Value of x:", x)
    fmt.Println("Address of x:", &x)  // & gets the address

    // Pointer variable
    var p *int  // Pointer to an int
    p = &x      // p now points to x

    fmt.Println("Value of p:", p)   // Memory address
    fmt.Println("Value at p:", *p)  // * dereferences (gets value)

    // Modify through pointer
    *p = 100
    fmt.Println("x is now:", x)  // 100!
}
```

**Key concepts:**
- `&variable` - Get the memory address
- `*pointer` - Dereference (get the value at that address)
- `*Type` - Pointer type (e.g., `*int`, `*string`)

#### Pointers with Structs

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func modifyByValue(p Person) {
    p.Age = 100  // Changes the COPY, not the original
}

func modifyByPointer(p *Person) {
    p.Age = 100  // Changes the ORIGINAL
}

func main() {
    alice := Person{Name: "Alice", Age: 30}

    modifyByValue(alice)
    fmt.Println(alice.Age)  // Still 30

    modifyByPointer(&alice)
    fmt.Println(alice.Age)  // Now 100
}
```

**Why this matters for databases:**
- Passing large structs by value copies all the data (slow)
- Passing pointers is fast (just 8 bytes on 64-bit systems)
- Databases need to modify data in place

#### new() and make()

```go
package main

import "fmt"

func main() {
    // new() allocates and returns a pointer
    p := new(int)  // Allocates an int, returns *int
    *p = 42
    fmt.Println(*p)

    // make() is for slices, maps, channels
    s := make([]int, 5)     // slice
    m := make(map[string]int)  // map

    // For structs, you usually use & with literal
    person := &Person{Name: "Bob", Age: 25}
    fmt.Println(person.Name)  // Go auto-dereferences
}

type Person struct {
    Name string
    Age  int
}
```

#### Interactive Exercise 1: Pointer Playground

```go
package main

import "fmt"

type Counter struct {
    Value int
}

// TODO: Implement Increment using a pointer receiver
func (c *Counter) Increment() {
    // Your code here
}

// TODO: Implement Add using a pointer receiver
func (c *Counter) Add(n int) {
    // Your code here
}

func main() {
    counter := Counter{Value: 0}

    counter.Increment()
    fmt.Println(counter.Value)  // Should be 1

    counter.Add(5)
    fmt.Println(counter.Value)  // Should be 6

    // Without pointer receiver, this wouldn't work!
}
```

#### Memory Layout Deep Dive

```go
package main

import (
    "fmt"
    "unsafe"
)

type SmallStruct struct {
    A int8   // 1 byte
    B int64  // 8 bytes
    C int8   // 1 byte
}

type OptimizedStruct struct {
    B int64  // 8 bytes
    A int8   // 1 byte
    C int8   // 1 byte
}

func main() {
    small := SmallStruct{}
    optimized := OptimizedStruct{}

    fmt.Println("SmallStruct size:", unsafe.Sizeof(small))      // 24 bytes (padding!)
    fmt.Println("OptimizedStruct size:", unsafe.Sizeof(optimized))  // 16 bytes

    // Why? Memory alignment!
    // Fields are aligned to their size
}
```

**Key concept:** Struct field order matters for memory efficiency. Put large fields first.

**Why this matters for databases:** A database page might fit 256 optimized structs but only 170 unoptimized ones!

### Day 3-4: Testing

Go has testing built into the language!

#### Your First Test

Create `math.go`:
```go
package math

func Add(a, b int) int {
    return a + b
}

func Multiply(a, b int) int {
    return a * b
}
```

Create `math_test.go`:
```go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

func TestMultiply(t *testing.T) {
    result := Multiply(3, 4)
    expected := 12
    if result != expected {
        t.Errorf("Multiply(3, 4) = %d; want %d", result, expected)
    }
}
```

Run tests:
```bash
go test
go test -v  # Verbose output
```

#### Table-Driven Tests (Idiomatic Go!)

```go
package math

import "testing"

func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"zero", 0, 5, 5},
        {"mixed", -2, 5, 3},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

Run specific test:
```bash
go test -run TestAdd
go test -run TestAdd/positive  # Run specific sub-test
```

#### Testing with Errors

Create `validator.go`:
```go
package validator

import "errors"

func ValidateAge(age int) error {
    if age < 0 {
        return errors.New("age cannot be negative")
    }
    if age > 150 {
        return errors.New("age too large")
    }
    return nil
}
```

Create `validator_test.go`:
```go
package validator

import "testing"

func TestValidateAge(t *testing.T) {
    tests := []struct {
        name      string
        age       int
        shouldErr bool
    }{
        {"valid age", 30, false},
        {"negative age", -1, true},
        {"too large", 200, true},
        {"zero", 0, false},
        {"boundary", 150, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAge(tt.age)
            if tt.shouldErr && err == nil {
                t.Error("expected error, got nil")
            }
            if !tt.shouldErr && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

#### Interactive Exercise 2: Write Tests

Create a `stringutils` package with these functions, then write tests:

```go
package stringutils

// Reverse returns the reversed string
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// IsPalindrome checks if a string is the same forwards and backwards
func IsPalindrome(s string) bool {
    // TODO: Implement this
    return false
}

// CountVowels returns the number of vowels in a string
func CountVowels(s string) int {
    // TODO: Implement this
    return 0
}
```

**Your tasks:**
1. Implement `IsPalindrome` and `CountVowels`
2. Write table-driven tests for all three functions
3. Make sure all tests pass: `go test -v`

### Day 5-6: Benchmarking

Benchmarks measure performance. Critical for database work!

#### Your First Benchmark

Add to `math_test.go`:
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

Run benchmark:
```bash
go test -bench=.
go test -bench=. -benchmem  # Show memory allocations
```

Output:
```
BenchmarkAdd-8   	1000000000	         0.25 ns/op	       0 B/op	       0 allocs/op
```

Meaning:
- `1000000000`: Number of iterations (b.N)
- `0.25 ns/op`: Time per operation
- `0 B/op`: Bytes allocated per operation
- `0 allocs/op`: Number of allocations per operation

#### Comparing Implementations

```go
package main

import "testing"

// Append implementation
func sumAppend(nums []int) int {
    result := 0
    for _, n := range nums {
        result += n
    }
    return result
}

// Pre-allocated implementation
func sumPrealloc(nums []int) int {
    result := 0
    for i := 0; i < len(nums); i++ {
        result += nums[i]
    }
    return result
}

func BenchmarkSumAppend(b *testing.B) {
    nums := make([]int, 1000)
    for i := range nums {
        nums[i] = i
    }

    b.ResetTimer()  // Don't count setup time

    for i := 0; i < b.N; i++ {
        sumAppend(nums)
    }
}

func BenchmarkSumPrealloc(b *testing.B) {
    nums := make([]int, 1000)
    for i := range nums {
        nums[i] = i
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        sumPrealloc(nums)
    }
}
```

Run and compare:
```bash
go test -bench=. -benchmem
```

#### Interactive Exercise 3: Optimize This!

```go
package main

import "testing"

// Slow version: Creates new slice on every append
func filterEvensSlow(nums []int) []int {
    var result []int
    for _, n := range nums {
        if n%2 == 0 {
            result = append(result, n)
        }
    }
    return result
}

// TODO: Fast version - pre-allocate slice
func filterEvensFast(nums []int) []int {
    // Hint: estimate capacity, or use make([]int, 0, len(nums))
    return nil
}

func BenchmarkFilterEvensSlow(b *testing.B) {
    nums := make([]int, 10000)
    for i := range nums {
        nums[i] = i
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        filterEvensSlow(nums)
    }
}

func BenchmarkFilterEvensFast(b *testing.B) {
    nums := make([]int, 10000)
    for i := range nums {
        nums[i] = i
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        filterEvensFast(nums)
    }
}
```

**Your tasks:**
1. Implement `filterEvensFast` with pre-allocation
2. Run benchmarks and compare allocations
3. Aim for 50%+ reduction in allocations

### Day 7: JSON and Serialization

Databases need to save data to disk. JSON is a simple format.

#### Encoding to JSON

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
    City string `json:"city,omitempty"`  // Omit if empty
}

func main() {
    person := Person{
        Name: "Alice",
        Age:  30,
        City: "New York",
    }

    // Marshal to JSON bytes
    jsonBytes, err := json.Marshal(person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(string(jsonBytes))
    // Output: {"name":"Alice","age":30,"city":"New York"}

    // Pretty print
    prettyJSON, _ := json.MarshalIndent(person, "", "  ")
    fmt.Println(string(prettyJSON))
}
```

#### Decoding from JSON

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    jsonStr := `{"name":"Bob","age":25}`

    var person Person
    err := json.Unmarshal([]byte(jsonStr), &person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)
}
```

#### Working with JSON Files

```go
package main

import (
    "encoding/json"
    "os"
)

type Config struct {
    DatabasePath string `json:"database_path"`
    MaxConnections int  `json:"max_connections"`
}

// Save to file
func (c *Config) Save(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(c)
}

// Load from file
func LoadConfig(filename string) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    return &config, err
}

func main() {
    config := Config{
        DatabasePath: "/data/mydb",
        MaxConnections: 100,
    }

    // Save
    config.Save("config.json")

    // Load
    loaded, err := LoadConfig("config.json")
    if err != nil {
        panic(err)
    }

    println(loaded.DatabasePath)
}
```

## Week 4: Building a Key-Value Store

**Project:** Build a persistent, in-memory key-value store with JSON persistence.

### Requirements

- Store string keys → string values
- Operations: `Set()`, `Get()`, `Delete()`, `List()`
- Save to JSON file on command
- Load from JSON file on startup
- Thread-safe (we'll learn this next week, for now single-threaded is OK)
- Handle errors properly
- Full test coverage
- Benchmarks for operations

### Step 1: Basic Structure

Create `kvstore/store.go`:

```go
package kvstore

import (
    "encoding/json"
    "errors"
    "os"
)

var (
    ErrKeyNotFound = errors.New("key not found")
    ErrEmptyKey    = errors.New("key cannot be empty")
)

type Store struct {
    data     map[string]string
    filename string
}

// New creates a new key-value store
func New(filename string) *Store {
    return &Store{
        data:     make(map[string]string),
        filename: filename,
    }
}

// Set stores a key-value pair
func (s *Store) Set(key, value string) error {
    if key == "" {
        return ErrEmptyKey
    }
    s.data[key] = value
    return nil
}

// Get retrieves a value by key
func (s *Store) Get(key string) (string, error) {
    if key == "" {
        return "", ErrEmptyKey
    }
    value, exists := s.data[key]
    if !exists {
        return "", ErrKeyNotFound
    }
    return value, nil
}

// Delete removes a key-value pair
func (s *Store) Delete(key string) error {
    if key == "" {
        return ErrEmptyKey
    }
    if _, exists := s.data[key]; !exists {
        return ErrKeyNotFound
    }
    delete(s.data, key)
    return nil
}

// List returns all keys
func (s *Store) List() []string {
    keys := make([]string, 0, len(s.data))
    for key := range s.data {
        keys = append(keys, key)
    }
    return keys
}

// Size returns the number of entries
func (s *Store) Size() int {
    return len(s.data)
}
```

### Step 2: Add Persistence

```go
// Save writes the store to disk
func (s *Store) Save() error {
    file, err := os.Create(s.filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(s.data)
}

// Load reads the store from disk
func (s *Store) Load() error {
    file, err := os.Open(s.filename)
    if err != nil {
        if os.IsNotExist(err) {
            // File doesn't exist yet, that's OK
            return nil
        }
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    return decoder.Decode(&s.data)
}
```

### Step 3: Write Tests

Create `kvstore/store_test.go`:

```go
package kvstore

import (
    "os"
    "testing"
)

func TestBasicOperations(t *testing.T) {
    store := New("test.json")

    // Test Set
    err := store.Set("name", "Alice")
    if err != nil {
        t.Fatalf("Set failed: %v", err)
    }

    // Test Get
    value, err := store.Get("name")
    if err != nil {
        t.Fatalf("Get failed: %v", err)
    }
    if value != "Alice" {
        t.Errorf("Get returned %s; want Alice", value)
    }

    // Test Size
    if store.Size() != 1 {
        t.Errorf("Size = %d; want 1", store.Size())
    }

    // Test Delete
    err = store.Delete("name")
    if err != nil {
        t.Fatalf("Delete failed: %v", err)
    }

    // Verify deleted
    _, err = store.Get("name")
    if err != ErrKeyNotFound {
        t.Error("Expected ErrKeyNotFound after delete")
    }
}

func TestPersistence(t *testing.T) {
    filename := "test_persist.json"
    defer os.Remove(filename)

    // Create and save
    store1 := New(filename)
    store1.Set("key1", "value1")
    store1.Set("key2", "value2")
    err := store1.Save()
    if err != nil {
        t.Fatalf("Save failed: %v", err)
    }

    // Load in new store
    store2 := New(filename)
    err = store2.Load()
    if err != nil {
        t.Fatalf("Load failed: %v", err)
    }

    // Verify data
    value, err := store2.Get("key1")
    if err != nil || value != "value1" {
        t.Errorf("After load, Get(key1) = %s, %v; want value1, nil", value, err)
    }

    if store2.Size() != 2 {
        t.Errorf("After load, Size = %d; want 2", store2.Size())
    }
}

func TestErrors(t *testing.T) {
    store := New("test.json")

    // Empty key
    err := store.Set("", "value")
    if err != ErrEmptyKey {
        t.Error("Expected ErrEmptyKey for empty key")
    }

    // Get non-existent key
    _, err = store.Get("nonexistent")
    if err != ErrKeyNotFound {
        t.Error("Expected ErrKeyNotFound")
    }

    // Delete non-existent key
    err = store.Delete("nonexistent")
    if err != ErrKeyNotFound {
        t.Error("Expected ErrKeyNotFound")
    }
}
```

Run tests:
```bash
go test -v
```

### Step 4: Add Benchmarks

Add to `store_test.go`:

```go
func BenchmarkSet(b *testing.B) {
    store := New("bench.json")
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        store.Set("key", "value")
    }
}

func BenchmarkGet(b *testing.B) {
    store := New("bench.json")
    store.Set("key", "value")
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        store.Get("key")
    }
}

func BenchmarkSetLarge(b *testing.B) {
    store := New("bench.json")

    // Pre-populate with 10,000 entries
    for i := 0; i < 10000; i++ {
        key := fmt.Sprintf("key%d", i)
        store.Set(key, "value")
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        key := fmt.Sprintf("newkey%d", i)
        store.Set(key, "value")
    }
}
```

Run benchmarks:
```bash
go test -bench=. -benchmem
```

### Step 5: Build a CLI

Create `main.go`:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "kvstore"  // Your package
)

func main() {
    store := kvstore.New("data.json")

    // Load existing data
    if err := store.Load(); err != nil {
        fmt.Println("Warning: Could not load data:", err)
    } else {
        fmt.Printf("Loaded %d entries\n", store.Size())
    }

    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("KV Store - Commands: set, get, delete, list, save, quit")

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }

        input := scanner.Text()
        parts := strings.Fields(input)

        if len(parts) == 0 {
            continue
        }

        command := parts[0]

        switch command {
        case "set":
            if len(parts) < 3 {
                fmt.Println("Usage: set <key> <value>")
                continue
            }
            key := parts[1]
            value := strings.Join(parts[2:], " ")
            if err := store.Set(key, value); err != nil {
                fmt.Println("Error:", err)
            } else {
                fmt.Println("OK")
            }

        case "get":
            if len(parts) < 2 {
                fmt.Println("Usage: get <key>")
                continue
            }
            value, err := store.Get(parts[1])
            if err != nil {
                fmt.Println("Error:", err)
            } else {
                fmt.Println(value)
            }

        case "delete":
            if len(parts) < 2 {
                fmt.Println("Usage: delete <key>")
                continue
            }
            if err := store.Delete(parts[1]); err != nil {
                fmt.Println("Error:", err)
            } else {
                fmt.Println("OK")
            }

        case "list":
            keys := store.List()
            fmt.Printf("Keys (%d):\n", len(keys))
            for _, key := range keys {
                fmt.Printf("  - %s\n", key)
            }

        case "save":
            if err := store.Save(); err != nil {
                fmt.Println("Error:", err)
            } else {
                fmt.Printf("Saved %d entries\n", store.Size())
            }

        case "quit":
            fmt.Println("Saving before exit...")
            if err := store.Save(); err != nil {
                fmt.Println("Error saving:", err)
            }
            fmt.Println("Goodbye!")
            return

        default:
            fmt.Println("Unknown command:", command)
        }
    }
}
```

### Enhancement Challenges

1. **Auto-save:** Automatically save after every write operation
2. **Timestamps:** Track when each key was created/modified
3. **TTL:** Add time-to-live for keys (auto-delete after N seconds)
4. **Search:** Search values by prefix or regex
5. **Transactions:** Batch multiple operations
6. **Backup:** Create timestamped backup files
7. **Stats:** Track operation counts and access patterns

## Week 3-4 Checkpoint

### Self-Assessment

Can you answer these?

1. What's the difference between `&x` and `*x`?
2. When should you use a pointer receiver vs value receiver?
3. How do you run all tests in a package?
4. What does `-benchmem` show you?
5. How do you save a struct to JSON?
6. What's the difference between `json.Marshal` and `json.Encode`?

### Practical Test

**Challenge:** Extend your key-value store with:

1. **Add a `Has(key string) bool` method**
2. **Add a `Clear()` method to delete all entries**
3. **Write tests for both**
4. **Write benchmarks comparing Get vs Has**
5. **Add JSON tags to properly serialize a Stats struct:**

```go
type Stats struct {
    TotalKeys      int
    LastSaveTime   time.Time
    OperationCount map[string]int  // "set", "get", "delete" counts
}
```

If you can complete this in 3-4 hours, you're ready for Week 5-6!

## What's Next?

You've learned:
- ✅ Pointers and memory layout
- ✅ Testing with table-driven tests
- ✅ Benchmarking and performance measurement
- ✅ JSON serialization
- ✅ File I/O
- ✅ Built a complete persistent key-value store

**Next:** [Week 5-6: Concurrency & I/O](week-5-6-concurrency-io.md)

You'll learn:
- Goroutines and channels
- Mutexes and race conditions
- Concurrent data structures
- Advanced file I/O
- Build a concurrent file downloader

## Resources

**Deepen your understanding:**
- [Go Pointers Explained](https://dave.cheney.net/2017/04/26/understand-go-pointers-in-less-than-800-words)
- [Effective Go - Testing](https://go.dev/doc/effective_go#testing)
- [Go Blog - Profiling](https://go.dev/blog/pprof)

**Practice more:**
- [LeetCode in Go](https://leetcode.com) - Practice problems
- [Project Euler](https://projecteuler.net/) - Math + programming challenges

Keep building, keep testing, keep learning!
