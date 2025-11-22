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

### Hello World - Your First Graph!

Create `hello_graph.go`:

```go
package main

import "fmt"

func main() {
    // Your first graph database!
    fmt.Println("=== My First Graph Database ===")

    // Nodes (people in a social network)
    nodes := []string{"Alice", "Bob", "Charlie", "David"}

    // Edges (who knows whom)
    edges := [][2]string{
        {"Alice", "Bob"},
        {"Bob", "Charlie"},
        {"Charlie", "David"},
        {"Alice", "Charlie"},
    }

    fmt.Println("\nNodes (People):")
    for i, node := range nodes {
        fmt.Printf("  %d: %s\n", i, node)
    }

    fmt.Println("\nEdges (Friendships):")
    for _, edge := range edges {
        fmt.Printf("  %s -> %s\n", edge[0], edge[1])
    }
}
```

**Run it:**
```bash
go run hello_graph.go
```

**Output:**
```
=== My First Graph Database ===

Nodes (People):
  0: Alice
  1: Bob
  2: Charlie
  3: David

Edges (Friendships):
  Alice -> Bob
  Bob -> Charlie
  Charlie -> David
  Alice -> Charlie
```

**Build it:**
```bash
go build hello_graph.go
./hello_graph
```

🎉 **Congratulations!** You just represented graph data in Go! This is the foundation of what Kuzu does.

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

### Challenge: Graph Query - Find Friends

Build on your graph to answer queries:

```go
package main

import "fmt"

func main() {
    // Your graph from before
    nodes := []string{"Alice", "Bob", "Charlie", "David"}
    edges := [][2]string{
        {"Alice", "Bob"},
        {"Bob", "Charlie"},
        {"Charlie", "David"},
        {"Alice", "Charlie"},
    }

    // Query: Who are Alice's friends?
    person := "Alice"
    fmt.Printf("\n%s's friends:\n", person)

    for _, edge := range edges {
        if edge[0] == person {
            fmt.Printf("  - %s\n", edge[1])
        }
    }
}
```

**Expected Output:**
```
Alice's friends:
  - Bob
  - Charlie
```

**Extend it:**
- Find friends of different people (change `person` variable)
- Count total friends for each person
- Find if two people are directly connected

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

**Database Connection:**
You just represented a graph using basic Go data structures! This is how Kuzu stores graph data internally (though more efficiently). You learned:
- Nodes can be represented as slices
- Edges can be represented as pairs (tuples)
- Simple queries are just loops and conditionals

**Next:** Let's organize these queries into reusable functions...

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

### Creating a Module - Graph Query Package

Let's organize your graph queries into a reusable package!

```bash
mkdir graphutil
cd graphutil
go mod init github.com/yourusername/graphutil
```

Create `graphutil/queries.go`:

```go
package graphutil

// FindNeighbors returns all nodes connected to the given node
// This is like Kuzu's neighbor traversal!
func FindNeighbors(node string, edges [][2]string) []string {
    var neighbors []string
    for _, edge := range edges {
        if edge[0] == node {
            neighbors = append(neighbors, edge[1])
        }
    }
    return neighbors
}

// CountDegree returns the number of outgoing edges from a node
func CountDegree(node string, edges [][2]string) int {
    count := 0
    for _, edge := range edges {
        if edge[0] == node {
            count++
        }
    }
    return count
}

// HasEdge checks if there's a direct edge from src to dst
func HasEdge(src, dst string, edges [][2]string) bool {
    for _, edge := range edges {
        if edge[0] == src && edge[1] == dst {
            return true
        }
    }
    return false
}

// helper is a private function (lowercase first letter)
func helper() {
    // Only accessible within this package
}
```

**Note:** Exported functions (uppercase) can be used outside the package. Private functions (lowercase) cannot.

### Using Your Graph Package

Create `main.go` in the parent directory:

```go
package main

import (
    "fmt"
    "github.com/yourusername/graphutil"
)

func main() {
    // Your graph data
    edges := [][2]string{
        {"Alice", "Bob"},
        {"Alice", "Charlie"},
        {"Bob", "David"},
        {"Charlie", "David"},
    }

    // Use your graph utilities!
    friends := graphutil.FindNeighbors("Alice", edges)
    fmt.Println("Alice's friends:", friends)

    degree := graphutil.CountDegree("Alice", edges)
    fmt.Printf("Alice has %d friends\n", degree)

    connected := graphutil.HasEdge("Alice", "Bob")
    fmt.Printf("Alice -> Bob exists: %v\n", connected)
}
```

**Output:**
```
Alice's friends: [Bob Charlie]
Alice has 2 friends
Alice -> Bob exists: true
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

- [ ] Create a `graphutil` package with graph query functions
- [ ] Implement `FindNeighbors`, `CountDegree`, and `HasEdge` functions
- [ ] Import and use these packages in a main program
- [ ] Understand the difference between exported and unexported identifiers
- [ ] Write tests for your graph query functions

### Challenge: Advanced Graph Queries

Extend your `graphutil` package with these functions:

```go
package graphutil

// GetAllNodes returns unique nodes from edge list
func GetAllNodes(edges [][2]string) []string {
    // Your implementation
    // Hint: Use a map to track unique nodes
}

// FindPath checks if path exists from src to dst (1 or 2 hops)
func FindPath(src, dst string, edges [][2]string) bool {
    // Your implementation
    // Check direct edge, then check through one intermediate node
}

// GetMostConnected returns the node with most outgoing edges
func GetMostConnected(edges [][2]string) string {
    // Your implementation
    // Count degrees for all nodes, return the max
}
```

**Test your functions:**
```go
edges := [][2]string{
    {"Alice", "Bob"},
    {"Bob", "Charlie"},
    {"Alice", "Charlie"},
}

nodes := graphutil.GetAllNodes(edges)
// Should return: ["Alice", "Bob", "Charlie"]

exists := graphutil.FindPath("Alice", "Charlie", edges)
// Should return: true (direct edge exists)

popular := graphutil.GetMostConnected(edges)
// Should return: "Alice" (2 outgoing edges)
```

### What You Should Discover

- Functions are first-class citizens in Go
- Multiple return values are idiomatic (especially for error handling)
- Package organization makes code reusable
- Exported vs unexported controls API surface

**Database Connection:**
You just built a graph query library! These functions are the foundation of database operations:
- `FindNeighbors` → Kuzu's graph traversal
- `CountDegree` → Analytical queries
- `HasEdge` → Edge existence checks

Real databases optimize these operations, but the logic is the same.

**Next:** Your queries work, but they're using strings for nodes. Let's create proper data structures with IDs and properties...

---

## Lesson 0.3: Structs, Interfaces, and Methods

**Your Mission:** Build a real graph database with proper data structures!

### Structs - Building Your Graph Database

Time to upgrade from simple strings to a real data model:

```go
package main

import "fmt"

// Node represents a vertex in the graph
type Node struct {
    ID         int
    Label      string // e.g., "Person", "Movie"
    Properties map[string]interface{} // Flexible properties
}

// Edge represents a relationship between nodes
type Edge struct {
    From  int    // Source node ID
    To    int    // Destination node ID
    Label string // e.g., "KNOWS", "LIKES"
}

func main() {
    // Create nodes
    alice := Node{
        ID:         1,
        Label:      "Person",
        Properties: map[string]interface{}{
            "name": "Alice",
            "age":  30,
        },
    }

    bob := Node{
        ID:         2,
        Label:      "Person",
        Properties: map[string]interface{}{
            "name": "Bob",
            "age":  25,
        },
    }

    // Create an edge: Alice knows Bob
    friendship := Edge{
        From:  alice.ID,
        To:    bob.ID,
        Label: "KNOWS",
    }

    fmt.Printf("Node: %s (id=%d)\n", alice.Properties["name"], alice.ID)
    fmt.Printf("Edge: %d -[%s]-> %d\n", friendship.From, friendship.Label, friendship.To)
}
```

**Output:**
```
Node: Alice (id=1)
Edge: 1 -[KNOWS]-> 2
```

🎉 **This is exactly how graph databases model data!**

### Methods - Adding Behavior to Your Graph

Now let's create a `Graph` type with methods:

```go
package main

import "fmt"

// Graph is your in-memory graph database!
type Graph struct {
    nodes map[int]*Node
    edges []*Edge
}

// NewGraph creates an empty graph
func NewGraph() *Graph {
    return &Graph{
        nodes: make(map[int]*Node),
        edges: make([]*Edge, 0),
    }
}

// AddNode inserts a node into the graph
func (g *Graph) AddNode(id int, label string) *Node {
    node := &Node{
        ID:         id,
        Label:      label,
        Properties: make(map[string]interface{}),
    }
    g.nodes[id] = node
    return node
}

// AddEdge creates a relationship between two nodes
func (g *Graph) AddEdge(from, to int, label string) *Edge {
    edge := &Edge{
        From:  from,
        To:    to,
        Label: label,
    }
    g.edges = append(g.edges, edge)
    return edge
}

// GetNeighbors returns all nodes connected to the given node
func (g *Graph) GetNeighbors(nodeID int) []*Node {
    var neighbors []*Node
    for _, edge := range g.edges {
        if edge.From == nodeID {
            if node, exists := g.nodes[edge.To]; exists {
                neighbors = append(neighbors, node)
            }
        }
    }
    return neighbors
}

// NodeCount returns the number of nodes
func (g *Graph) NodeCount() int {
    return len(g.nodes)
}

func main() {
    // Create your graph database!
    g := NewGraph()

    // Add nodes
    alice := g.AddNode(1, "Person")
    alice.Properties["name"] = "Alice"
    alice.Properties["age"] = 30

    bob := g.AddNode(2, "Person")
    bob.Properties["name"] = "Bob"

    charlie := g.AddNode(3, "Person")
    charlie.Properties["name"] = "Charlie"

    // Add edges
    g.AddEdge(1, 2, "KNOWS")
    g.AddEdge(1, 3, "KNOWS")

    // Query the graph!
    friends := g.GetNeighbors(1)
    fmt.Printf("Alice's friends:\n")
    for _, friend := range friends {
        fmt.Printf("  - %s\n", friend.Properties["name"])
    }

    fmt.Printf("\nTotal nodes: %d\n", g.NodeCount())
}
```

**Output:**
```
Alice's friends:
  - Bob
  - Charlie

Total nodes: 3
```

**Pointer Receivers:**
All Graph methods use `(g *Graph)` because we're modifying the graph. This avoids copying the entire graph on every method call!

🚀 **You just built an in-memory graph database!**

### Interfaces - Flexible Graph Storage

Interfaces let you define different storage backends for your graph:

```go
package main

import "fmt"

// GraphStorage defines how to store/retrieve graph data
type GraphStorage interface {
    SaveGraph(g *Graph) error
    LoadGraph() (*Graph, error)
}

// MemoryStorage keeps graph in RAM
type MemoryStorage struct {
    data *Graph
}

func (m *MemoryStorage) SaveGraph(g *Graph) error {
    m.data = g
    fmt.Println("Graph saved to memory")
    return nil
}

func (m *MemoryStorage) LoadGraph() (*Graph, error) {
    return m.data, nil
}

// Later you could add:
// type DiskStorage struct { ... }
// type CloudStorage struct { ... }

func BackupGraph(storage GraphStorage, g *Graph) {
    storage.SaveGraph(g)
}

func main() {
    g := NewGraph()
    g.AddNode(1, "Person")

    // Use any storage that implements GraphStorage
    memory := &MemoryStorage{}
    BackupGraph(memory, g)

    // Could swap to disk storage without changing BackupGraph!
}
```

**Key Concept:** Interfaces are satisfied implicitly. No `implements` keyword needed. This is how Kuzu supports different storage backends!

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

- [ ] Build the `Node`, `Edge`, and `Graph` structs shown above
- [ ] Add methods: `AddNode`, `AddEdge`, `GetNeighbors`, `NodeCount`, `EdgeCount`
- [ ] Create a social network with at least 5 people and connections
- [ ] Query the graph to find someone's friends
- [ ] Test your graph with different queries

### Challenge: Graph Analytics

Extend your Graph with analytical methods:

```go
// GetDegree returns the number of outgoing edges from a node
func (g *Graph) GetDegree(nodeID int) int {
    // Your implementation
}

// FindCommonNeighbors returns nodes that both n1 and n2 are connected to
func (g *Graph) FindCommonNeighbors(n1, n2 int) []*Node {
    // Your implementation
    // Hint: Find neighbors of both, return intersection
}

// GetMostConnected returns the node with the highest degree
func (g *Graph) GetMostConnected() *Node {
    // Your implementation
}

// HasPath checks if there's a path from src to dst (up to 2 hops)
func (g *Graph) HasPath(src, dst int) bool {
    // Your implementation
    // Check direct connection, then check through intermediates
}
```

**Test your analytics:**
```go
g := NewGraph()
alice := g.AddNode(1, "Person")
alice.Properties["name"] = "Alice"
bob := g.AddNode(2, "Person")
bob.Properties["name"] = "Bob"
charlie := g.AddNode(3, "Person")
charlie.Properties["name"] = "Charlie"

g.AddEdge(1, 2, "KNOWS")
g.AddEdge(1, 3, "KNOWS")
g.AddEdge(2, 3, "KNOWS")

degree := g.GetDegree(1) // Should be 2
mostConnected := g.GetMostConnected() // Should be Alice
hasPath := g.HasPath(1, 3) // Should be true
```

### What You Should Discover

- Go uses composition over inheritance
- Interfaces are implicit and flexible
- Methods can be defined on any type (not just structs)
- Pointer receivers are common for modifying state

**Database Connection:**
🎉 **HUGE WIN!** You just built a working in-memory graph database with:
- Node and Edge data model (same as Kuzu!)
- Graph storage and query methods
- Property support (flexible schema)
- Neighbor traversal (graph queries)

Your code can handle thousands of nodes. But what about millions? We need:
- Parallel queries (next lesson!)
- Disk persistence (Lesson 0.6)
- Optimized storage (Phase 1!)

**Next:** Your graph works great, but queries are sequential. Let's make them run in parallel...

---

## Lesson 0.4: Concurrency Basics

**Your Mission:** Make your graph database queries run in parallel - just like Kuzu!

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

- [ ] Add `ParallelScan` method to your Graph (processes all nodes concurrently)
- [ ] Add `ParallelFilter` method to find nodes matching a condition
- [ ] Test parallel queries on a graph with 1000+ nodes
- [ ] Compare performance: sequential vs parallel

### Challenge: Parallel Graph Analytics

Add these methods to your Graph struct from Lesson 0.3:

```go
// ParallelScan processes all nodes concurrently
func (g *Graph) ParallelScan(fn func(*Node)) {
    var wg sync.WaitGroup

    for _, node := range g.nodes {
        wg.Add(1)
        go func(n *Node) {
            defer wg.Done()
            fn(n)  // Process node in parallel
        }(node)
    }

    wg.Wait()
}

// ParallelFilter finds nodes matching a predicate (in parallel!)
func (g *Graph) ParallelFilter(predicate func(*Node) bool) []*Node {
    resultsCh := make(chan *Node, len(g.nodes))
    var wg sync.WaitGroup

    // Process nodes in parallel
    for _, node := range g.nodes {
        wg.Add(1)
        go func(n *Node) {
            defer wg.Done()
            if predicate(n) {
                resultsCh <- n
            }
        }(node)
    }

    // Close channel when all workers done
    go func() {
        wg.Wait()
        close(resultsCh)
    }()

    // Collect results
    var results []*Node
    for node := range resultsCh {
        results = append(results, node)
    }

    return results
}
```

**Test your parallel queries:**
```go
g := NewGraph()

// Create a social network
for i := 1; i <= 100; i++ {
    node := g.AddNode(i, "Person")
    node.Properties["age"] = 20 + (i % 50)
}

// Find all adults (age > 30) in parallel!
adults := g.ParallelFilter(func(n *Node) bool {
    if age, ok := n.Properties["age"].(int); ok {
        return age > 30
    }
    return false
})

fmt.Printf("Found %d adults\n", len(adults))

// Print all nodes in parallel
g.ParallelScan(func(n *Node) {
    fmt.Printf("Node %d: age=%v\n", n.ID, n.Properties["age"])
})
```

**This is exactly what Kuzu does for parallel query execution!**

### What You Should Discover

- Goroutines are cheap and plentiful
- Channels make concurrent code easier to reason about
- `select` enables complex coordination patterns
- Race conditions are real - use `-race` flag to detect them

**Database Connection:**
You just implemented parallel query execution! Your Graph can now:
- Scan all nodes concurrently (parallel table scan)
- Filter nodes in parallel (parallel predicate evaluation)
- Process 100K+ nodes efficiently

This is the foundation of Kuzu's **morsel-driven parallelism**. Phase 2 will teach you even more advanced patterns!

**Performance win:** On a 4-core machine, parallel queries can be 3-4x faster!

**Next:** Parallel queries are powerful, but what if a query fails? We need error handling...

---

## Lesson 0.5: Error Handling and Testing

**Your Mission:** Make your graph database robust with error handling and comprehensive tests!

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

### Custom Error Types - Graph Errors

Add error handling to your Graph from Lesson 0.3:

```go
package main

import (
    "errors"
    "fmt"
)

// Graph errors
var (
    ErrNodeNotFound = errors.New("node not found")
    ErrEdgeNotFound = errors.New("edge not found")
    ErrDuplicateNode = errors.New("node already exists")
)

// QueryError provides detailed error information
type QueryError struct {
    Operation string
    NodeID    int
    Err       error
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("%s failed for node %d: %v", e.Operation, e.NodeID, e.Err)
}

// GetNode retrieves a node by ID with error handling
func (g *Graph) GetNode(id int) (*Node, error) {
    node, exists := g.nodes[id]
    if !exists {
        return nil, fmt.Errorf("%w: id=%d", ErrNodeNotFound, id)
    }
    return node, nil
}

// AddNodeSafe creates a node with duplicate checking
func (g *Graph) AddNodeSafe(id int, label string) (*Node, error) {
    if _, exists := g.nodes[id]; exists {
        return nil, fmt.Errorf("%w: id=%d", ErrDuplicateNode, id)
    }

    node := &Node{
        ID:         id,
        Label:      label,
        Properties: make(map[string]interface{}),
    }
    g.nodes[id] = node
    return node, nil
}

// AddEdgeSafe creates an edge with validation
func (g *Graph) AddEdgeSafe(from, to int, label string) (*Edge, error) {
    // Validate source node exists
    if _, exists := g.nodes[from]; !exists {
        return nil, &QueryError{
            Operation: "AddEdge",
            NodeID:    from,
            Err:       ErrNodeNotFound,
        }
    }

    // Validate destination node exists
    if _, exists := g.nodes[to]; !exists {
        return nil, &QueryError{
            Operation: "AddEdge",
            NodeID:    to,
            Err:       ErrNodeNotFound,
        }
    }

    edge := &Edge{From: from, To: to, Label: label}
    g.edges = append(g.edges, edge)
    return edge, nil
}
```

**Usage with error handling:**
```go
g := NewGraph()

// This will succeed
_, err := g.AddNodeSafe(1, "Person")
if err != nil {
    fmt.Println("Error:", err)
}

// This will fail - node already exists
_, err = g.AddNodeSafe(1, "Person")
if errors.Is(err, ErrDuplicateNode) {
    fmt.Println("Node already exists!")
}

// This will fail - node 2 doesn't exist yet
_, err = g.AddEdgeSafe(1, 2, "KNOWS")
if err != nil {
    fmt.Println("Error:", err)
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

### Writing Tests - Testing Your Graph

Create `graph_test.go` to test your graph database:

```go
package main

import (
    "errors"
    "testing"
)

func TestAddNode(t *testing.T) {
    g := NewGraph()
    node := g.AddNode(1, "Person")

    if node.ID != 1 {
        t.Errorf("Expected ID 1, got %d", node.ID)
    }

    if node.Label != "Person" {
        t.Errorf("Expected label 'Person', got %s", node.Label)
    }
}

func TestGetNeighbors(t *testing.T) {
    g := NewGraph()
    g.AddNode(1, "Person")
    g.AddNode(2, "Person")
    g.AddNode(3, "Person")

    g.AddEdge(1, 2, "KNOWS")
    g.AddEdge(1, 3, "KNOWS")

    neighbors := g.GetNeighbors(1)

    if len(neighbors) != 2 {
        t.Errorf("Expected 2 neighbors, got %d", len(neighbors))
    }
}

func TestAddEdgeSafe(t *testing.T) {
    // Table-driven tests
    tests := []struct {
        name      string
        from      int
        to        int
        wantError error
    }{
        {
            name:      "both nodes exist",
            from:      1,
            to:        2,
            wantError: nil,
        },
        {
            name:      "source node missing",
            from:      99,
            to:        2,
            wantError: ErrNodeNotFound,
        },
        {
            name:      "dest node missing",
            from:      1,
            to:        99,
            wantError: ErrNodeNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            g := NewGraph()
            g.AddNode(1, "Person")
            g.AddNode(2, "Person")

            _, err := g.AddEdgeSafe(tt.from, tt.to, "KNOWS")

            if tt.wantError != nil {
                if err == nil {
                    t.Errorf("Expected error %v, got nil", tt.wantError)
                } else if !errors.Is(err, tt.wantError) {
                    t.Errorf("Expected error %v, got %v", tt.wantError, err)
                }
            } else if err != nil {
                t.Errorf("Expected no error, got %v", err)
            }
        })
    }
}

func TestParallelFilter(t *testing.T) {
    g := NewGraph()

    // Create 100 nodes
    for i := 1; i <= 100; i++ {
        node := g.AddNode(i, "Person")
        node.Properties["age"] = 20 + (i % 50)
    }

    // Find adults (age > 30)
    adults := g.ParallelFilter(func(n *Node) bool {
        if age, ok := n.Properties["age"].(int); ok {
            return age > 30
        }
        return false
    })

    // Should find nodes with age 31-70
    if len(adults) == 0 {
        t.Error("Expected to find some adults")
    }
}
```

**Run tests:**
```bash
go test
go test -v        # Verbose output
go test -cover    # Show coverage
go test -race     # Detect race conditions (important for parallel code!)
```

### Benchmarks - Measure Performance

Add benchmarks to `graph_test.go`:

```go
func BenchmarkAddNode(b *testing.B) {
    g := NewGraph()
    for i := 0; i < b.N; i++ {
        g.AddNode(i, "Person")
    }
}

func BenchmarkGetNeighbors(b *testing.B) {
    g := NewGraph()
    // Setup: create graph with 1000 nodes
    for i := 0; i < 1000; i++ {
        g.AddNode(i, "Person")
    }
    for i := 0; i < 999; i++ {
        g.AddEdge(i, i+1, "KNOWS")
    }

    b.ResetTimer() // Don't count setup time

    for i := 0; i < b.N; i++ {
        g.GetNeighbors(500)
    }
}

func BenchmarkParallelFilter(b *testing.B) {
    g := NewGraph()
    for i := 0; i < 10000; i++ {
        node := g.AddNode(i, "Person")
        node.Properties["age"] = 20 + (i % 50)
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        g.ParallelFilter(func(n *Node) bool {
            if age, ok := n.Properties["age"].(int); ok {
                return age > 30
            }
            return false
        })
    }
}
```

**Run benchmarks:**
```bash
go test -bench=.
go test -bench=. -benchmem  # Include memory stats

# Compare sequential vs parallel:
go test -bench=BenchmarkParallelFilter
```

### Your Tasks

- [ ] Add error handling to all Graph methods (`GetNode`, `AddEdgeSafe`, etc.)
- [ ] Define custom error types for graph operations
- [ ] Write comprehensive tests for your graph database
- [ ] Write benchmarks comparing sequential vs parallel queries
- [ ] Run `go test -race` to ensure no race conditions
- [ ] Achieve >80% test coverage

### Challenge: Complete Test Suite

Build a comprehensive test suite for your graph database:

```go
// Test basic operations
func TestGraphOperations(t *testing.T) {
    // Test adding nodes, edges
    // Test querying neighbors
    // Test error conditions
}

// Test concurrent access
func TestConcurrentAccess(t *testing.T) {
    g := NewGraph()
    var wg sync.WaitGroup

    // Add nodes from multiple goroutines
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            g.AddNode(id, "Person")
        }(i)
    }

    wg.Wait()

    // Should have 10 nodes with no data races
    if g.NodeCount() != 10 {
        t.Errorf("Expected 10 nodes, got %d", g.NodeCount())
    }
}

// Benchmark sequential vs parallel filter
func BenchmarkSequentialFilter(b *testing.B) {
    g := setupLargeGraph(10000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Sequential filter implementation
        var results []*Node
        for _, node := range g.nodes {
            if predicate(node) {
                results = append(results, node)
            }
        }
    }
}

func BenchmarkParallelFilter(b *testing.B) {
    g := setupLargeGraph(10000)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        g.ParallelFilter(predicate)
    }
}
```

**Run your tests:**
```bash
go test -v
go test -race
go test -cover
go test -bench=. -benchmem
```

**Goal:** See parallel filter 3-4x faster than sequential on 4-core machine!

### What You Should Discover

- Explicit error handling makes code more robust
- `defer` is essential for cleanup
- Table-driven tests are idiomatic in Go
- Benchmarks should be part of your workflow

**Database Connection:**
Your graph database now has production-quality error handling and tests! You've learned:
- Custom error types for database operations
- Comprehensive test coverage (unit + integration + benchmarks)
- Race detection for concurrent code
- Performance measurement and comparison

Real databases are tested extensively - Kuzu has thousands of tests!

**Next:** Your graph lives in memory. When the program ends, data is lost. Let's add persistence...

---

## Lesson 0.6: File I/O and System Calls

**Your Mission:** Save your graph database to disk and load it back - data persistence!

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

### Graph Persistence - Save/Load Your Database

Add persistence methods to your Graph from Lesson 0.3:

```go
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
)

// Save writes the graph to a JSON file
func (g *Graph) Save(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    // Create a serializable format
    data := struct {
        Nodes []*Node `json:"nodes"`
        Edges []*Edge `json:"edges"`
    }{
        Nodes: make([]*Node, 0, len(g.nodes)),
        Edges: g.edges,
    }

    // Convert map to slice for JSON
    for _, node := range g.nodes {
        data.Nodes = append(data.Nodes, node)
    }

    // Write JSON
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(data); err != nil {
        return fmt.Errorf("failed to encode graph: %w", err)
    }

    return nil
}

// Load reads a graph from a JSON file
func Load(filename string) (*Graph, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    var data struct {
        Nodes []*Node `json:"nodes"`
        Edges []*Edge `json:"edges"`
    }

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&data); err != nil {
        return nil, fmt.Errorf("failed to decode graph: %w", err)
    }

    // Reconstruct graph
    g := NewGraph()
    for _, node := range data.Nodes {
        g.nodes[node.ID] = node
    }
    g.edges = data.Edges

    return g, nil
}

// SaveCSV exports edges in CSV format (like Kuzu's input format!)
func (g *Graph) SaveCSV(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    // Write header
    fmt.Fprintln(writer, "from,to,label")

    // Write edges
    for _, edge := range g.edges {
        fmt.Fprintf(writer, "%d,%d,%s\n", edge.From, edge.To, edge.Label)
    }

    return nil
}
```

**Usage:**
```go
func main() {
    // Create graph
    g := NewGraph()
    alice := g.AddNode(1, "Person")
    alice.Properties["name"] = "Alice"
    alice.Properties["age"] = 30

    bob := g.AddNode(2, "Person")
    bob.Properties["name"] = "Bob"

    g.AddEdge(1, 2, "KNOWS")

    // Save to disk
    if err := g.Save("graph.json"); err != nil {
        fmt.Println("Save error:", err)
        return
    }
    fmt.Println("Graph saved to graph.json")

    // Export edges to CSV
    if err := g.SaveCSV("edges.csv"); err != nil {
        fmt.Println("CSV export error:", err)
        return
    }
    fmt.Println("Edges exported to edges.csv")

    // Load from disk
    g2, err := Load("graph.json")
    if err != nil {
        fmt.Println("Load error:", err)
        return
    }

    fmt.Printf("Loaded graph with %d nodes and %d edges\n",
        g2.NodeCount(), len(g2.edges))
}
```

**Output:**
```
Graph saved to graph.json
Edges exported to edges.csv
Loaded graph with 2 nodes and 1 edges
```

**graph.json:**
```json
{
  "nodes": [
    {
      "ID": 1,
      "Label": "Person",
      "Properties": {
        "age": 30,
        "name": "Alice"
      }
    },
    {
      "ID": 2,
      "Label": "Person",
      "Properties": {
        "name": "Bob"
      }
    }
  ],
  "edges": [
    {
      "From": 1,
      "To": 2,
      "Label": "KNOWS"
    }
  ]
}
```

**edges.csv:**
```csv
from,to,label
1,2,KNOWS
```

### Your Tasks

- [ ] Add `Save` and `Load` methods to your Graph
- [ ] Add `SaveCSV` method to export edges
- [ ] Test saving and loading a graph with 100+ nodes
- [ ] Handle errors when file doesn't exist or is corrupted
- [ ] Verify loaded graph has same data as original

### Challenge: CSV Import

Implement CSV import for your graph database (just like Kuzu!):

```go
// LoadCSV imports edges from CSV file
func LoadCSV(filename string) (*Graph, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    g := NewGraph()
    scanner := bufio.NewScanner(file)

    // Skip header line
    scanner.Scan()

    lineNum := 1
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()

        var from, to int
        var label string

        // Parse CSV line: "from,to,label"
        _, err := fmt.Sscanf(line, "%d,%d,%s", &from, &to, &label)
        if err != nil {
            fmt.Printf("Warning: skipping malformed line %d: %s\n", lineNum, line)
            continue
        }

        // Auto-create nodes if they don't exist
        if _, exists := g.nodes[from]; !exists {
            g.AddNode(from, "Node")
        }
        if _, exists := g.nodes[to]; !exists {
            g.AddNode(to, "Node")
        }

        g.AddEdge(from, to, label)
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return g, nil
}
```

**Test with real data:**
Create `social_network.csv`:
```csv
from,to,label
1,2,KNOWS
1,3,KNOWS
2,4,KNOWS
3,4,KNOWS
4,5,KNOWS
```

Load and query:
```go
g, _ := LoadCSV("social_network.csv")
fmt.Printf("Loaded %d nodes and %d edges\n", g.NodeCount(), len(g.edges))

neighbors := g.GetNeighbors(1)
fmt.Printf("Node 1 has %d neighbors\n", len(neighbors))
```

### What You Should Discover

- Go's `os` package mirrors Unix file operations
- Buffered I/O is usually faster for small reads/writes
- JSON is flexible but slower than binary formats
- CSV is a universal interchange format
- Error handling is crucial for I/O operations

**Database Connection:**
🎉 **FINAL WIN!** Your graph database now has complete persistence!
- Save/load graph in JSON format
- Export/import edges in CSV (just like Kuzu!)
- Handle file errors gracefully
- Compatible with standard data formats

**Limitations discovered:**
- JSON files are large (3-4x data size)
- Loading is slow for large graphs (seconds for 1M nodes)
- No crash recovery - if program crashes during save, data is lost

**This motivates Phase 1:** Kuzu uses binary pages, buffer pools, and WAL for:
- 100x faster loading (memory-mapped files)
- Efficient disk usage (compression)
- Crash recovery (write-ahead logging)

---

## Phase 0 Complete: You Built a Graph Database!

**Congratulations!** 🎉 You didn't just learn Go - you built a **working graph database**!

### What You Built

Your graph database has:

```
my-graph-db/
├── graph.go              # Node, Edge, Graph structs
├── graph_test.go         # Comprehensive tests
├── queries.go            # GetNeighbors, GetDegree, etc.
├── parallel.go           # ParallelScan, ParallelFilter
├── errors.go             # Error types and validation
└── persistence.go        # Save/Load (JSON and CSV)
```

### Capabilities

Your database can:
- ✅ **Create** nodes and edges with properties (schema-free!)
- ✅ **Query** neighbors and traverse the graph
- ✅ **Filter** nodes by properties (age > 30, etc.)
- ✅ **Run parallel queries** (3-4x faster on multi-core!)
- ✅ **Handle errors** gracefully (duplicate nodes, missing edges)
- ✅ **Persist to disk** (JSON format + CSV import/export)
- ✅ **Test coverage** (unit, integration, benchmarks)
- ✅ **Handle ~100K nodes** efficiently in memory

### What You Learned

**Go Skills:**
- ✅ Go syntax and basic types
- ✅ Functions, packages, and modules
- ✅ Structs, interfaces, and methods
- ✅ Goroutines and channels (concurrency)
- ✅ Error handling and testing
- ✅ File I/O and persistence

**Database Skills:**
- ✅ Graph data modeling (nodes, edges, properties)
- ✅ Query execution (traversal, filtering)
- ✅ Parallel query processing
- ✅ Data persistence and serialization
- ✅ Error handling and validation
- ✅ Performance benchmarking

### Try It Out!

Load a real social network:

```go
// Create graph
g := NewGraph()

// Load from CSV
g, _ = LoadCSV("social_network.csv")
fmt.Printf("Loaded %d nodes, %d edges\n", g.NodeCount(), len(g.edges))

// Query
friends := g.GetNeighbors(1)
fmt.Printf("Person 1 has %d friends\n", len(friends))

// Parallel analytics
adults := g.ParallelFilter(func(n *Node) bool {
    age, ok := n.Properties["age"].(int)
    return ok && age >= 18
})
fmt.Printf("Found %d adults\n", len(adults))

// Save
g.Save("backup.json")
```

### Where It Breaks Down

Your database is great for small to medium graphs, but has limitations:

**Problem 1: Memory constraints**
- Can't load 1M+ nodes (out of memory)
- No way to query without loading entire graph

**Problem 2: Slow persistence**
- JSON is 3-4x larger than necessary
- Loading 100K nodes takes seconds
- No crash recovery (losing data on crash)

**Problem 3: Query performance**
- Scanning all edges is O(n) - slow for large graphs
- No indexes for fast lookups
- Parallel queries help, but still scanning all data

### Transition to Phase 1: Professional Storage

**Phase 1 solves these problems with Kuzu's architecture:**

**Pages (Lesson 1.1):**
- Fixed 4KB blocks instead of JSON
- Memory-mapped I/O (100x faster than JSON)
- Binary encoding (1/4 the size)

**Buffer Pool (Lesson 1.2):**
- LRU cache for hot pages
- Query graphs larger than RAM
- Concurrent access with latching

**Write-Ahead Log (Lesson 1.3):**
- Crash recovery
- ACID guarantees
- Fast writes

### The Bridge

Your current graph:
```go
type Graph struct {
    nodes map[int]*Node
    edges []*Edge
}

g.Save("graph.json")  // Slow, large files
```

Phase 1 upgrade:
```go
type Graph struct {
    pages *PageManager      // Disk-backed pages
    buffer *BufferPool      // LRU cache
    wal *WriteAheadLog      // Crash recovery
}

g.pages.LoadPage(pageID)  // Fast, memory-mapped
```

**Same concepts, professional implementation!**

### Next Steps

You're ready for **Phase 1: Storage Layer**!

**What to expect:**
- Same Go skills, applied to real database internals
- Build on your Graph code (upgrade, don't replace)
- Learn techniques used by PostgreSQL, SQLite, and Kuzu
- Handle million-node graphs efficiently

**Preview - Lesson 1.1:**
> "Your graph works great for 10K nodes. Let's see what happens at 1M...
>
> [Tries to load 1M nodes]
> - JSON file: 450MB
> - Load time: 12 seconds
> - Memory: 1.2GB
>
> With pages and mmap:
> - Binary file: 120MB (3.7x smaller!)
> - Load time: 0.3 seconds (40x faster!)
> - Memory: 50MB (24x less!)
>
> Let's rebuild your persistence layer with pages..."

**Ready to level up?** 🚀

→ Continue to [Phase 1: Storage Layer](phase-1-storage-layer.md)

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
