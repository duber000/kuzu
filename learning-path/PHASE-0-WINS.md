# Phase 0: Database-Focused Learning Path

**Goal:** Every lesson should produce a working piece that directly contributes to building Kuzu.

## Current Problem

Phase 0 teaches:
- Hello World → prints to console
- Functions → string utilities, calculators
- Structs → shapes, bank accounts
- Concurrency → parallel sum
- Error handling → calculator errors
- File I/O → generic file operations

**Then Phase 1 suddenly introduces:** pages, binary encoding, buffer pools...

**Gap:** No clear path from "hello world" to "graph database"

---

## Proposed Solution: Progressive Database Wins

Each lesson builds a real piece of the database, starting simple and getting more sophisticated.

### Lesson 0.1: Hello World & Basic Types
**Current:** Print "Hello, World!"
**Database Win:** Read and display simple graph data

```go
// hello_graph.go
package main

import "fmt"

func main() {
    // Your first graph!
    nodes := []string{"Alice", "Bob", "Charlie"}
    edges := [][2]string{
        {"Alice", "Bob"},
        {"Bob", "Charlie"},
    }

    fmt.Println("=== My First Graph ===")
    fmt.Println("Nodes:", nodes)
    fmt.Println("\nEdges:")
    for _, edge := range edges {
        fmt.Printf("  %s -> %s\n", edge[0], edge[1])
    }
}
```

**What they learn:**
- ✅ Go syntax
- ✅ Arrays and slices
- ✅ For loops
- ✅ **Database win:** Representing graph data in memory

**Next step:** "Now let's make this more flexible with functions..."

---

### Lesson 0.2: Functions & Packages
**Current:** String utilities (Reverse, IsPalindrome)
**Database Win:** Graph query utilities

```go
// graph_utils.go
package graphutil

// FindNeighbors returns all nodes connected to the given node
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

// HasPath checks if there's a direct edge from src to dst
func HasPath(src, dst string, edges [][2]string) bool {
    for _, edge := range edges {
        if edge[0] == src && edge[1] == dst {
            return true
        }
    }
    return false
}
```

**What they learn:**
- ✅ Functions and parameters
- ✅ Return values
- ✅ Packages and modules
- ✅ **Database win:** Basic graph traversal algorithms

**Test file shows real queries:**
```go
func TestFindNeighbors(t *testing.T) {
    edges := [][2]string{
        {"Alice", "Bob"},
        {"Alice", "Charlie"},
        {"Bob", "David"},
    }

    neighbors := FindNeighbors("Alice", edges)
    // Alice has 2 friends: Bob and Charlie
    if len(neighbors) != 2 {
        t.Errorf("Expected 2 neighbors, got %d", len(neighbors))
    }
}
```

**Next step:** "These functions work, but we're representing nodes as strings. Let's use proper data structures..."

---

### Lesson 0.3: Structs, Interfaces & Methods
**Current:** Shapes, bank accounts, stacks
**Database Win:** Node, Edge, and Graph data structures

```go
// graph.go
package graph

// Node represents a vertex in the graph
type Node struct {
    ID         int
    Label      string
    Properties map[string]interface{}
}

// Edge represents a relationship between nodes
type Edge struct {
    From       int // Node ID
    To         int // Node ID
    Label      string
    Properties map[string]interface{}
}

// Graph is our in-memory graph database
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

// AddNode adds a node to the graph
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
        From:       from,
        To:         to,
        Label:      label,
        Properties: make(map[string]interface{}),
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

// EdgeCount returns the number of edges
func (g *Graph) EdgeCount() int {
    return len(g.edges)
}
```

**Usage example:**
```go
func main() {
    g := graph.NewGraph()

    // Create social network
    alice := g.AddNode(1, "Person")
    alice.Properties["name"] = "Alice"
    alice.Properties["age"] = 30

    bob := g.AddNode(2, "Person")
    bob.Properties["name"] = "Bob"
    bob.Properties["age"] = 25

    // Alice knows Bob
    friendship := g.AddEdge(1, 2, "KNOWS")
    friendship.Properties["since"] = 2020

    // Query: Who does Alice know?
    friends := g.GetNeighbors(1)
    for _, friend := range friends {
        fmt.Printf("Alice knows %s\n", friend.Properties["name"])
    }
}
```

**What they learn:**
- ✅ Structs for data modeling
- ✅ Methods on structs
- ✅ Pointers vs values
- ✅ Maps for properties
- ✅ **Database win:** Core graph database data model!

**This is HUGE:** They just built an in-memory graph database!

**Next step:** "This works, but it's slow for large graphs. Let's make queries run in parallel..."

---

### Lesson 0.4: Concurrency Basics
**Current:** Parallel sum, worker pools
**Database Win:** Parallel graph scanning

```go
// parallel_scan.go
package graph

import (
    "sync"
)

// ParallelScan processes all nodes concurrently
func (g *Graph) ParallelScan(fn func(*Node)) {
    var wg sync.WaitGroup

    for _, node := range g.nodes {
        wg.Add(1)
        go func(n *Node) {
            defer wg.Done()
            fn(n)
        }(node)
    }

    wg.Wait()
}

// ParallelFilter finds nodes matching a predicate
func (g *Graph) ParallelFilter(predicate func(*Node) bool) []*Node {
    // Create channel for results
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

    // Close channel when done
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

**Usage:**
```go
// Find all people over 25 in parallel
adults := g.ParallelFilter(func(n *Node) bool {
    if age, ok := n.Properties["age"].(int); ok {
        return age > 25
    }
    return false
})
```

**What they learn:**
- ✅ Goroutines
- ✅ WaitGroups
- ✅ Channels
- ✅ **Database win:** Parallel query execution (like Kuzu's morsel-driven parallelism!)

**Next step:** "Parallel queries are fast, but what happens when they fail? We need error handling..."

---

### Lesson 0.5: Error Handling & Testing
**Current:** Calculator errors, defer/panic/recover
**Database Win:** Query validation and error handling

```go
// errors.go
package graph

import (
    "errors"
    "fmt"
)

var (
    ErrNodeNotFound = errors.New("node not found")
    ErrEdgeNotFound = errors.New("edge not found")
    ErrInvalidQuery = errors.New("invalid query")
)

// QueryError contains detailed error information
type QueryError struct {
    Query string
    Err   error
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("query failed: %q - %v", e.Query, e.Err)
}

// GetNode retrieves a node by ID with error handling
func (g *Graph) GetNode(id int) (*Node, error) {
    node, exists := g.nodes[id]
    if !exists {
        return nil, fmt.Errorf("%w: id=%d", ErrNodeNotFound, id)
    }
    return node, nil
}

// AddEdgeWithValidation creates an edge after validating nodes exist
func (g *Graph) AddEdgeWithValidation(from, to int, label string) (*Edge, error) {
    // Validate source node exists
    if _, exists := g.nodes[from]; !exists {
        return nil, fmt.Errorf("source %w: id=%d", ErrNodeNotFound, from)
    }

    // Validate destination node exists
    if _, exists := g.nodes[to]; !exists {
        return nil, fmt.Errorf("destination %w: id=%d", ErrNodeNotFound, to)
    }

    edge := &Edge{
        From:       from,
        To:         to,
        Label:      label,
        Properties: make(map[string]interface{}),
    }
    g.edges = append(g.edges, edge)

    return edge, nil
}

// ExecuteQuery runs a query with proper error handling
func (g *Graph) ExecuteQuery(query string) (interface{}, error) {
    defer func() {
        if r := recover(); r != nil {
            // Convert panic to error
            fmt.Printf("Query panicked: %v\n", r)
        }
    }()

    // Query execution logic here...
    // This is a placeholder
    if query == "" {
        return nil, ErrInvalidQuery
    }

    return nil, nil
}
```

**Tests:**
```go
func TestGetNode(t *testing.T) {
    g := NewGraph()
    g.AddNode(1, "Person")

    // Success case
    node, err := g.GetNode(1)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    // Error case
    _, err = g.GetNode(999)
    if !errors.Is(err, ErrNodeNotFound) {
        t.Errorf("Expected ErrNodeNotFound, got %v", err)
    }
}

func TestAddEdgeWithValidation(t *testing.T) {
    g := NewGraph()
    g.AddNode(1, "Person")

    // Should fail: node 2 doesn't exist
    _, err := g.AddEdgeWithValidation(1, 2, "KNOWS")
    if err == nil {
        t.Error("Expected error for non-existent node")
    }
}
```

**What they learn:**
- ✅ Error types and error wrapping
- ✅ Custom errors
- ✅ Defer/panic/recover
- ✅ Table-driven tests
- ✅ **Database win:** Robust query error handling

**Next step:** "Our graph works in memory, but it disappears when the program ends. Let's persist it to disk..."

---

### Lesson 0.6: File I/O & System Calls
**Current:** Generic file operations, KV store
**Database Win:** Graph serialization and simple storage format

```go
// storage.go
package graph

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
)

// Save writes the graph to a file in JSON format
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

// Load reads a graph from a file
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

// LoadCSV imports edges from CSV
func LoadCSV(filename string) (*Graph, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    g := NewGraph()
    scanner := bufio.NewScanner(file)

    // Skip header
    scanner.Scan()

    // Read edges
    for scanner.Scan() {
        var from, to int
        var label string

        line := scanner.Text()
        _, err := fmt.Sscanf(line, "%d,%d,%s", &from, &to, &label)
        if err != nil {
            continue // Skip malformed lines
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

    return g, scanner.Err()
}
```

**Usage:**
```go
func main() {
    // Create graph
    g := graph.NewGraph()
    alice := g.AddNode(1, "Person")
    alice.Properties["name"] = "Alice"
    bob := g.AddNode(2, "Person")
    bob.Properties["name"] = "Bob"
    g.AddEdge(1, 2, "KNOWS")

    // Save to disk
    g.Save("graph.json")
    g.SaveCSV("edges.csv")

    // Load from disk
    g2, _ := graph.Load("graph.json")
    fmt.Printf("Loaded graph with %d nodes\n", g2.NodeCount())
}
```

**What they learn:**
- ✅ File I/O operations
- ✅ JSON encoding/decoding
- ✅ CSV parsing
- ✅ Buffered I/O
- ✅ **Database win:** Data persistence (foundation for Kuzu's storage layer!)

**Connection to Phase 1:** "JSON is flexible but slow. Next, we'll learn binary formats for 100x faster I/O using pages and memory-mapped files..."

---

## Phase 0 Complete: You Built a Graph Database!

**What the learner has built:**

```
my-graph-db/
├── graph.go           # Node, Edge, Graph structs
├── graph_test.go      # Comprehensive tests
├── operations.go      # Query functions (GetNeighbors, etc.)
├── parallel.go        # Concurrent graph scanning
├── errors.go          # Error handling
└── storage.go         # Persistence (JSON, CSV)
```

**Capabilities:**
- ✅ Create nodes and edges
- ✅ Store properties on nodes/edges
- ✅ Query neighbors and traverse graph
- ✅ Run parallel queries
- ✅ Handle errors gracefully
- ✅ Save/load from disk (JSON and CSV)
- ✅ **~500 lines of working code!**

**Transition to Phase 1:**
> "Your graph database works great for small datasets, but what about 1 million nodes? 1 billion?
>
> JSON files are too slow. Scanning all nodes is inefficient. Loading the entire graph into memory doesn't scale.
>
> **That's where Kuzu's architecture comes in:**
> - **Pages:** Store data in 4KB blocks (like a real database)
> - **Buffer Pool:** Cache hot pages in memory
> - **WAL:** Crash recovery and durability
>
> Let's rebuild your storage layer with professional database techniques..."

---

## Benefits of This Approach

### 1. **Immediate Gratification**
Every lesson produces working code that *does something with graphs*.

### 2. **Progressive Enhancement**
- Lesson 0.1: Arrays and slices → represent graph
- Lesson 0.2: Functions → query graph
- Lesson 0.3: Structs → model graph professionally
- Lesson 0.4: Goroutines → parallel queries
- Lesson 0.5: Errors → robust queries
- Lesson 0.6: File I/O → persist graph

Each builds on the previous!

### 3. **Smooth Transition to Phase 1**
Phase 1 isn't starting from scratch—it's *optimizing* what they built:
- "You stored graphs in JSON → now use binary pages"
- "You loaded entire graph → now use buffer pool"
- "You had no crash recovery → now add WAL"

### 4. **Motivation**
They can show friends: "I built a graph database!" (even if it's simple)

### 5. **Testing Skills**
Every lesson includes tests that verify graph operations work correctly.

---

## Implementation Plan

1. **Create Phase 0 Progressive Examples**
   - Write complete code for each lesson
   - Include tests and benchmarks
   - Add "connection to Kuzu" notes

2. **Update Lesson Structure**
   - Keep current concepts (structs, concurrency, etc.)
   - Replace generic examples with graph database examples
   - Add "What you built" recap at each stage

3. **Create Transition Guide**
   - "From In-Memory to Disk Storage" (Phase 0 → Phase 1)
   - Show side-by-side: JSON persistence vs Pages
   - Benchmark the difference

4. **Add "Mini Project" at End of Phase 0**
   - Load a real dataset (e.g., social network CSV)
   - Run queries (find friends, degrees of separation)
   - Benchmark performance limits
   - Identify pain points that Phase 1 solves

---

## Example: Phase 0 Final Project

**Build a working social network graph database:**

```bash
$ go run main.go load social_network.csv
Loaded 10,000 nodes and 50,000 edges

$ go run main.go query "NEIGHBORS 1234"
Node 1234 (Alice) has 15 neighbors:
  - Bob (id: 567)
  - Charlie (id: 890)
  ...

$ go run main.go stats
Graph Statistics:
  Nodes: 10,000
  Edges: 50,000
  Avg Degree: 5.0

  Performance:
    Neighbor query: 0.1ms
    Full scan: 45ms
    Memory usage: 25MB

$ go run main.go save graph.json
Graph saved to graph.json (3.2MB, took 150ms)
```

**Then in Phase 1:**
> "Your graph works great for 10K nodes, but let's try 1 million...
>
> [Loads 1M node graph]
>
> - JSON file: 500MB
> - Load time: 8 seconds
> - Memory usage: 750MB
> - Query time: still fast (0.1ms)
>
> **Problems:**
> 1. Can't fit 100M nodes in RAM
> 2. Loading 8 seconds is too slow
> 3. No way to update graph efficiently
>
> **Solution: Kuzu's storage architecture**
>
> Let's rebuild with pages, buffer pool, and WAL..."

---

## Questions?

This approach ensures:
- ✅ Every concept is immediately applicable
- ✅ Progressive wins build confidence
- ✅ Smooth transition between phases
- ✅ Real database gets built, not just exercises
- ✅ Learners understand *why* Kuzu is designed this way

Ready to implement this?
