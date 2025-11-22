# Phase 0: Go Fundamentals

**Welcome to Go!** This phase takes you from zero to comfortable with Go basics. No prior Go experience required.

**Time Estimate:** 2-3 weeks
**Go Version:** 1.23+ recommended

---

## Table of Contents

- [Lesson 0.1: Hello World & Basic Syntax](#lesson-01-hello-world--basic-syntax)
- [Lesson 0.2: Functions, Packages, and Modules](#lesson-02-functions-packages-and-modules)
- [Lesson 0.3: Structs, Interfaces, and Methods](#lesson-03-structs-interfaces-and-methods)
- [Lesson 0.4: Concurrency Basics](#lesson-04-concurrency-basics)
- [Lesson 0.5: Error Handling and Testing](#lesson-05-error-handling-and-testing)
- [Lesson 0.6: File I/O and System Calls](#lesson-06-file-io-and-system-calls)

---

## Lesson 0.1: Hello World & Basic Syntax

**Your Mission:** Write and run your first Go program.

### Setup

1. **Install Go 1.23+:**
   ```bash
   # Download from https://go.dev/dl/
   # Or use a version manager like gvm

   # Verify installation:
   go version  # Should show go1.23 or later
   ```

2. **Set up your workspace:**
   ```bash
   mkdir -p ~/learning-go
   cd ~/learning-go
   ```

### Hello World

Create `hello.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**Run it:**
```bash
go run hello.go
# Output: Hello, World!
```

**Build it:**
```bash
go build hello.go
./hello
# Output: Hello, World!
```

### Investigation Questions

1. **What does `package main` mean?**
   - Try changing it to `package foo` - what happens?
   - What's special about package `main`?

2. **Imports:**
   ```go
   import "fmt"
   import "time"

   // vs

   import (
       "fmt"
       "time"
   )
   ```
   - Which style is preferred?
   - What happens if you import something you don't use?

3. **Variables and Types:**
   ```go
   package main

   import "fmt"

   func main() {
       // Explicit type declaration
       var name string = "Alice"
       var age int = 30

       // Type inference
       var city = "NYC"

       // Short declaration (inside functions only)
       country := "USA"

       fmt.Println(name, age, city, country)
   }
   ```
   - What's the difference between `:=` and `=`?
   - What's the zero value of `int`, `string`, `bool`?

### Your Tasks

- [ ] Print "Hello, World!" to the console
- [ ] Create variables of types: `int`, `string`, `bool`, `float64`
- [ ] Print their values using `fmt.Printf` with format specifiers
- [ ] Experiment with constants using `const`

### Challenge: Basic Calculator

```go
package main

import "fmt"

func main() {
    var a int = 10
    var b int = 3

    fmt.Printf("%d + %d = %d\n", a, b, a+b)
    fmt.Printf("%d - %d = %d\n", a, b, a-b)
    fmt.Printf("%d * %d = %d\n", a, b, a*b)
    fmt.Printf("%d / %d = %d\n", a, b, a/b)
    fmt.Printf("%d %% %d = %d\n", a, b, a%b)
}
```

**Extend it:**
- Accept user input using `fmt.Scan`
- Add error handling for division by zero
- Support floating-point operations

### Control Flow

```go
// If statements
if age >= 18 {
    fmt.Println("Adult")
} else {
    fmt.Println("Minor")
}

// For loops (Go only has 'for', no 'while')
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style loop
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}

// Infinite loop
for {
    // break to exit
    break
}

// Switch statements
switch day := "Monday"; day {
case "Monday":
    fmt.Println("Start of the week")
case "Friday":
    fmt.Println("Almost weekend!")
default:
    fmt.Println("Midweek")
}
```

### Arrays and Slices

```go
// Arrays (fixed size)
var arr [5]int
arr[0] = 1
fmt.Println(arr) // [1 0 0 0 0]

// Slices (dynamic size - used more commonly)
slice := []int{1, 2, 3}
slice = append(slice, 4, 5)
fmt.Println(slice) // [1 2 3 4 5]
fmt.Println(len(slice)) // 5
fmt.Println(cap(slice)) // capacity

// Slice operations
subSlice := slice[1:4] // [2 3 4]
```

### Maps

```go
// Create a map
ages := make(map[string]int)
ages["Alice"] = 30
ages["Bob"] = 25

// Or initialize with values
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}

// Check if key exists
age, exists := ages["Alice"]
if exists {
    fmt.Println("Alice is", age)
}

// Iterate over map
for name, age := range ages {
    fmt.Printf("%s is %d years old\n", name, age)
}

// Delete a key
delete(ages, "Bob")
```

### What You Should Discover

- Go is statically typed but has type inference
- Go's syntax is minimal and consistent
- Zero values prevent uninitialized variable bugs
- Slices are more flexible than arrays
- Maps are built-in and efficient

---

## Lesson 0.2: Functions, Packages, and Modules

**Your Mission:** Organize code into reusable functions and packages.

### Functions Basics

```go
package main

import "fmt"

// Simple function
func greet(name string) {
    fmt.Println("Hello,", name)
}

// Function with return value
func add(a int, b int) int {
    return a + b
}

// Shorthand for same type parameters
func multiply(a, b int) int {
    return a * b
}

// Multiple return values
func divide(a, b int) (int, int) {
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}

// Named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // naked return
}

func main() {
    greet("Alice")

    sum := add(3, 5)
    fmt.Println("Sum:", sum)

    q, r := divide(10, 3)
    fmt.Println("10 / 3 =", q, "remainder", r)

    // Ignore a return value with _
    quot, _ := divide(10, 3)
    fmt.Println("Quotient:", quot)
}
```

### Variadic Functions

```go
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

func main() {
    fmt.Println(sum(1, 2, 3))        // 6
    fmt.Println(sum(1, 2, 3, 4, 5))  // 15
}
```

### Creating a Module

```bash
mkdir mathutil
cd mathutil
go mod init github.com/yourusername/mathutil
```

Create `mathutil/calc.go`:

```go
package mathutil

// Add returns the sum of two integers
func Add(a, b int) int {
    return a + b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
    return a * b
}

// Private function (lowercase first letter)
func helper() {
    // Only accessible within this package
}
```

**Note:** Exported functions/types start with uppercase, private ones with lowercase.

### Using Your Package

Create `main.go` in the parent directory:

```go
package main

import (
    "fmt"
    "github.com/yourusername/mathutil"
)

func main() {
    result := mathutil.Add(10, 20)
    fmt.Println("Result:", result)

    product := mathutil.Multiply(5, 6)
    fmt.Println("Product:", product)
}
```

### Investigation Questions

1. **What's the difference between a package and a module?**
   - A package is a directory of `.go` files
   - A module is a collection of packages (defined by `go.mod`)

2. **Package naming conventions:**
   ```go
   // Good
   package stringutil

   // Bad
   package string_util
   package stringUtil
   ```
   - Go prefers short, lowercase, single-word package names

3. **Import paths:**
   ```go
   import (
       "fmt"                              // Standard library
       "math/rand"                        // Standard library sub-package
       "github.com/yourusername/package"  // External package
   )
   ```

### Your Tasks

- [ ] Create a `mathutil` package with functions for basic math operations
- [ ] Create a `stringutil` package with functions like `Reverse`, `IsPalindrome`
- [ ] Import and use these packages in a main program
- [ ] Understand the difference between exported and unexported identifiers

### Challenge: String Utilities

Create a package that provides:

```go
package stringutil

// Reverse returns the reversed string
func Reverse(s string) string {
    // Your implementation
}

// IsPalindrome checks if a string is a palindrome
func IsPalindrome(s string) bool {
    // Your implementation
}

// CountVowels returns the number of vowels
func CountVowels(s string) int {
    // Your implementation
}
```

### What You Should Discover

- Functions are first-class citizens in Go
- Multiple return values are idiomatic (especially for error handling)
- Package organization makes code reusable
- Exported vs unexported controls API surface

---

## Lesson 0.3: Structs, Interfaces, and Methods

**Your Mission:** Learn Go's approach to object-oriented programming.

### Structs

```go
package main

import "fmt"

// Define a struct
type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    // Create a struct
    alice := Person{
        Name: "Alice",
        Age:  30,
        City: "NYC",
    }

    fmt.Println(alice)
    fmt.Println(alice.Name) // Access fields

    // Shorthand (order matters)
    bob := Person{"Bob", 25, "LA"}

    // Pointer to struct
    p := &alice
    p.Age = 31 // Automatically dereferenced
    fmt.Println(alice.Age) // 31
}
```

### Methods

Go doesn't have classes, but you can define methods on types:

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Method with pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}

    fmt.Println("Area:", rect.Area()) // 50

    rect.Scale(2)
    fmt.Println("New dimensions:", rect.Width, rect.Height) // 20 10
}
```

**Value vs Pointer Receivers:**
- Value receiver: Method operates on a copy (read-only)
- Pointer receiver: Method can modify the original (read-write)

### Interfaces

```go
// Interface defines behavior
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

// Function accepts any Shape
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}

    PrintShapeInfo(rect)
    PrintShapeInfo(circle)
}
```

**Key Concept:** Interfaces are satisfied implicitly. No `implements` keyword needed.

### Empty Interface

```go
func PrintAnything(v interface{}) {
    fmt.Println(v)
}

func main() {
    PrintAnything(42)
    PrintAnything("hello")
    PrintAnything([]int{1, 2, 3})
}
```

`interface{}` (or `any` in Go 1.18+) can hold any type.

### Type Assertions

```go
var i interface{} = "hello"

// Type assertion
s := i.(string)
fmt.Println(s) // hello

// Type assertion with check
s, ok := i.(string)
if ok {
    fmt.Println("String:", s)
}

// Type switch
switch v := i.(type) {
case int:
    fmt.Println("Integer:", v)
case string:
    fmt.Println("String:", v)
default:
    fmt.Println("Unknown type")
}
```

### Investigation Questions

1. **When to use value vs pointer receivers?**
   - Use pointer if you need to modify the receiver
   - Use pointer for large structs (avoid copying)
   - Be consistent within a type

2. **Common interfaces in Go:**
   ```go
   // io.Reader
   type Reader interface {
       Read(p []byte) (n int, err error)
   }

   // io.Writer
   type Writer interface {
       Write(p []byte) (n int, err error)
   }

   // fmt.Stringer
   type Stringer interface {
       String() string
   }
   ```

### Your Tasks

- [ ] Create a `BankAccount` struct with `Deposit` and `Withdraw` methods
- [ ] Create an `Animal` interface with `Speak()` method
- [ ] Implement `Animal` for `Dog` and `Cat` types
- [ ] Write a function that accepts any `Animal` and makes it speak

### Challenge: Stack Data Structure

```go
type Stack struct {
    items []int
}

func (s *Stack) Push(item int) {
    // Add item to stack
}

func (s *Stack) Pop() (int, bool) {
    // Remove and return top item
    // Return false if stack is empty
}

func (s *Stack) Peek() (int, bool) {
    // Return top item without removing
}

func (s *Stack) IsEmpty() bool {
    // Check if stack is empty
}
```

Implement these methods and test them.

### What You Should Discover

- Go uses composition over inheritance
- Interfaces are implicit and flexible
- Methods can be defined on any type (not just structs)
- Pointer receivers are common for modifying state

---

## Lesson 0.4: Concurrency Basics

**Your Mission:** Learn goroutines and channels - Go's superpowers.

### Goroutines

```go
package main

import (
    "fmt"
    "time"
)

func say(s string) {
    for i := 0; i < 5; i++ {
        fmt.Println(s)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Run in a separate goroutine
    go say("world")

    // Run in main goroutine
    say("hello")
}
```

**Note:** Goroutines are lightweight threads managed by the Go runtime.

### Channels

```go
package main

import "fmt"

func main() {
    // Create a channel
    messages := make(chan string)

    // Send to channel (in goroutine)
    go func() {
        messages <- "ping"
    }()

    // Receive from channel
    msg := <-messages
    fmt.Println(msg) // ping
}
```

### Buffered Channels

```go
// Unbuffered channel (blocks until receiver ready)
ch := make(chan int)

// Buffered channel (blocks only when buffer full)
ch := make(chan int, 2)

ch <- 1
ch <- 2
// ch <- 3 // Would block (buffer full)

fmt.Println(<-ch) // 1
fmt.Println(<-ch) // 2
```

### Channel Synchronization

```go
func worker(done chan bool) {
    fmt.Println("Working...")
    time.Sleep(1 * time.Second)
    fmt.Println("Done")

    done <- true
}

func main() {
    done := make(chan bool, 1)
    go worker(done)

    <-done // Wait for worker to finish
}
```

### Select Statement

```go
func main() {
    c1 := make(chan string)
    c2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        c1 <- "one"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        c2 <- "two"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("Received", msg1)
        case msg2 := <-c2:
            fmt.Println("Received", msg2)
        }
    }
}
```

### WaitGroups

```go
import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // Decrement counter when done

    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1) // Increment counter
        go worker(i, &wg)
    }

    wg.Wait() // Wait for all workers
    fmt.Println("All workers done")
}
```

### Mutexes (for shared state)

```go
import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    counter := Counter{}

    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
    fmt.Println("Final value:", counter.Value()) // 1000
}
```

### Investigation Questions

1. **What's the difference between goroutines and OS threads?**
   - Goroutines are much lighter (~2KB stack vs 1-2MB for threads)
   - Go runtime multiplexes goroutines onto OS threads
   - Can have millions of goroutines in one program

2. **When to use channels vs mutexes?**
   - Channels: For communication between goroutines
   - Mutexes: For protecting shared state
   - Slogan: "Share memory by communicating, don't communicate by sharing memory"

3. **Common concurrency patterns:**
   - Worker pools
   - Pipeline processing
   - Fan-out/fan-in
   - Rate limiting

### Your Tasks

- [ ] Create 10 goroutines that each print their ID
- [ ] Use a channel to send numbers from one goroutine to another
- [ ] Implement a worker pool (5 workers processing 100 jobs)
- [ ] Use `sync.Mutex` to protect a shared counter

### Challenge: Parallel Sum

```go
// Calculate sum of slice using multiple goroutines
func parallelSum(numbers []int, numWorkers int) int {
    // Split work among goroutines
    // Use channels to collect partial sums
    // Return total sum
}

func main() {
    numbers := make([]int, 1000000)
    for i := range numbers {
        numbers[i] = i + 1
    }

    result := parallelSum(numbers, 4)
    fmt.Println("Sum:", result)
}
```

### What You Should Discover

- Goroutines are cheap and plentiful
- Channels make concurrent code easier to reason about
- `select` enables complex coordination patterns
- Race conditions are real - use `-race` flag to detect them

---

## Lesson 0.5: Error Handling and Testing

**Your Mission:** Write robust code with proper error handling and tests.

### Error Handling

Go doesn't have exceptions. Errors are values.

```go
package main

import (
    "errors"
    "fmt"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)

    // This will error
    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Custom Error Types

```go
type DivisionError struct {
    Dividend float64
    Divisor  float64
}

func (e *DivisionError) Error() string {
    return fmt.Sprintf("cannot divide %.2f by %.2f", e.Dividend, e.Divisor)
}

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, &DivisionError{Dividend: a, Divisor: b}
    }
    return a / b, nil
}
```

### Error Wrapping (Go 1.13+)

```go
import (
    "fmt"
    "errors"
)

func readConfig() error {
    err := readFile("config.json")
    if err != nil {
        return fmt.Errorf("failed to read config: %w", err)
    }
    return nil
}

// Check for specific error
func main() {
    err := readConfig()
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("Config file doesn't exist")
    }
}
```

### Defer, Panic, and Recover

```go
// Defer: Execute at end of function
func main() {
    defer fmt.Println("World")
    fmt.Println("Hello")
}
// Output:
// Hello
// World

// Common use: Resource cleanup
func processFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close() // Always closes, even if error occurs

    // Process file...
    return nil
}

// Panic: Unrecoverable error
func mustConnect(addr string) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        panic("cannot connect to " + addr)
    }
    defer conn.Close()
}

// Recover: Catch panic
func safeExecute() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    panic("something went wrong")
}
```

### Writing Tests

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
    expected := 5

    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}

func TestMultiply(t *testing.T) {
    tests := []struct {
        a, b     int
        expected int
    }{
        {2, 3, 6},
        {0, 5, 0},
        {-2, 3, -6},
    }

    for _, tt := range tests {
        result := Multiply(tt.a, tt.b)
        if result != tt.expected {
            t.Errorf("Multiply(%d, %d) = %d; want %d",
                tt.a, tt.b, result, tt.expected)
        }
    }
}
```

**Run tests:**
```bash
go test
go test -v        # Verbose output
go test -cover    # Show coverage
go test -race     # Detect race conditions
```

### Benchmarks

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

func BenchmarkMultiply(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Multiply(2, 3)
    }
}
```

**Run benchmarks:**
```bash
go test -bench=.
go test -bench=. -benchmem  # Include memory stats
```

### Your Tasks

- [ ] Write functions that return errors
- [ ] Use `defer` for cleanup operations
- [ ] Write table-driven tests
- [ ] Write benchmarks for your functions
- [ ] Achieve >80% test coverage

### Challenge: Safe Calculator

```go
type Calculator struct{}

func (c *Calculator) Divide(a, b float64) (float64, error) {
    // Return error if b is zero
}

func (c *Calculator) Sqrt(x float64) (float64, error) {
    // Return error if x is negative
}

func (c *Calculator) Factorial(n int) (int, error) {
    // Return error if n is negative
    // Use defer/panic/recover for overflow protection
}
```

Write comprehensive tests for each function.

### What You Should Discover

- Explicit error handling makes code more robust
- `defer` is essential for cleanup
- Table-driven tests are idiomatic in Go
- Benchmarks should be part of your workflow

---

## Lesson 0.6: File I/O and System Calls

**Your Mission:** Work with files, directories, and system resources.

### Reading Files

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Read entire file
    data, err := os.ReadFile("test.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(data))
}
```

### Writing Files

```go
func main() {
    content := []byte("Hello, File!")

    err := os.WriteFile("output.txt", content, 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
}
```

### Working with File Objects

```go
import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // Open file
    file, err := os.Open("test.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Read line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
    }
}
```

### Creating and Writing Files

```go
func main() {
    // Create file
    file, err := os.Create("output.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // Write to file
    writer := bufio.NewWriter(file)
    writer.WriteString("Line 1\n")
    writer.WriteString("Line 2\n")
    writer.Flush() // Important!
}
```

### Working with Directories

```go
import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Create directory
    os.Mkdir("testdir", 0755)

    // Create nested directories
    os.MkdirAll("parent/child/grandchild", 0755)

    // List directory contents
    entries, err := os.ReadDir(".")
    if err != nil {
        panic(err)
    }

    for _, entry := range entries {
        fmt.Println(entry.Name(), entry.IsDir())
    }

    // Walk directory tree
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        fmt.Println(path, info.Size())
        return nil
    })
}
```

### File Information

```go
func main() {
    fileInfo, err := os.Stat("test.txt")
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("File doesn't exist")
        }
        return
    }

    fmt.Println("Name:", fileInfo.Name())
    fmt.Println("Size:", fileInfo.Size())
    fmt.Println("Mode:", fileInfo.Mode())
    fmt.Println("Modified:", fileInfo.ModTime())
    fmt.Println("Is Directory:", fileInfo.IsDir())
}
```

### System Calls (Low-Level I/O)

```go
import (
    "fmt"
    "syscall"
)

func main() {
    // Open file with syscall
    fd, err := syscall.Open("test.txt", syscall.O_RDONLY, 0)
    if err != nil {
        panic(err)
    }
    defer syscall.Close(fd)

    // Read into buffer
    buffer := make([]byte, 1024)
    n, err := syscall.Read(fd, buffer)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Read %d bytes: %s\n", n, string(buffer[:n]))
}
```

### Memory-Mapped Files

```go
import (
    "fmt"
    "os"
    "syscall"
)

func main() {
    file, err := os.Open("data.bin")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    fileInfo, _ := file.Stat()
    size := int(fileInfo.Size())

    // Memory-map the file
    data, err := syscall.Mmap(
        int(file.Fd()),
        0,
        size,
        syscall.PROT_READ,
        syscall.MAP_SHARED,
    )
    if err != nil {
        panic(err)
    }
    defer syscall.Munmap(data)

    // Access file as byte slice
    fmt.Println(data[:10])
}
```

### Investigation Questions

1. **Buffered vs Unbuffered I/O:**
   - When to use `bufio`?
   - What's the performance difference?
   - Measure: Read 1MB file line-by-line with and without buffering

2. **File Permissions:**
   ```go
   os.Create("file.txt")           // 0666
   os.OpenFile("file.txt", flags, 0644)
   ```
   - What do permission bits mean? (read=4, write=2, execute=1)
   - Owner, group, others

3. **When to use syscalls directly?**
   - Performance-critical code
   - Direct I/O (bypassing OS cache)
   - Memory-mapped files

### Your Tasks

- [ ] Read a text file and count lines, words, characters
- [ ] Copy a file (implement `cp` command)
- [ ] List all `.go` files in a directory recursively
- [ ] Create a simple log file system with timestamps

### Challenge: File-Based Key-Value Store

```go
type KVStore struct {
    filename string
    data     map[string]string
}

func NewKVStore(filename string) (*KVStore, error) {
    // Load from file if exists
}

func (kv *KVStore) Set(key, value string) error {
    // Update in-memory map
    // Write to file
}

func (kv *KVStore) Get(key string) (string, bool) {
    // Read from in-memory map
}

func (kv *KVStore) Delete(key string) error {
    // Remove from map
    // Update file
}

func (kv *KVStore) Save() error {
    // Persist all data to file
}
```

**Requirements:**
- Persist data in a simple format (JSON or custom)
- Handle concurrent access with mutex
- Recover from crashes (file should always be valid)

### What You Should Discover

- Go's `os` package mirrors Unix file operations
- Buffered I/O is usually faster for small reads/writes
- Memory-mapped files are powerful for large datasets
- Error handling is crucial for I/O operations

---

## Phase 0 Complete!

**Congratulations!** You now know enough Go to start building a graph database.

**What you've learned:**
- ✅ Go syntax and basic types
- ✅ Functions, packages, and modules
- ✅ Structs, interfaces, and methods (Go's OOP)
- ✅ Goroutines and channels (concurrency)
- ✅ Error handling and testing
- ✅ File I/O and system calls

**Next Steps:**
- Move on to **Phase 1: Storage Layer** where you'll apply these concepts
- You'll build page-based storage, buffer pools, and write-ahead logs
- Everything from here builds on what you learned in Phase 0

**Ready?** Let's build a database! 🚀

---

## Quick Reference

### Common Commands

```bash
# Run code
go run main.go

# Build executable
go build

# Run tests
go test
go test -v
go test -cover
go test -race

# Run benchmarks
go test -bench=.
go test -bench=. -benchmem

# Format code
go fmt ./...

# Vet code (static analysis)
go vet ./...

# Get dependencies
go get package-name

# Module management
go mod init module-name
go mod tidy
go mod download
```

### Useful Resources

- **Official Tour:** https://go.dev/tour/
- **Go by Example:** https://gobyexample.com/
- **Effective Go:** https://go.dev/doc/effective_go
- **Standard Library:** https://pkg.go.dev/std
- **Go Playground:** https://go.dev/play/

Happy coding! 🎉
