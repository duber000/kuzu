# Practice Projects and Exercises

This document contains hands-on projects to reinforce concepts from each phase of the learning path. Complete these exercises to solidify your understanding before moving to the next phase.

## Pre-Work Practice Projects

### Project 1: CLI Tool with Flags
**Concepts:** flag package, file I/O, error handling

Build a command-line tool that:
- Accepts flags for input/output files
- Reads CSV data
- Performs filtering and aggregation
- Writes results to output file

**Requirements:**
- Use `flag` package for CLI arguments
- Proper error handling with wrapped errors
- Table-driven tests
- Benchmark I/O performance

### Project 2: Concurrent Web Scraper
**Concepts:** Goroutines, channels, sync, HTTP

Build a concurrent web scraper that:
- Scrapes multiple URLs in parallel
- Uses worker pool pattern
- Respects rate limits
- Aggregates results

**Requirements:**
- Worker pool with configurable size
- Use channels for work distribution
- Context for cancellation
- Proper error handling

### Project 3: Key-Value Store
**Concepts:** Maps, file I/O, serialization

Build an in-memory key-value store with:
- GET, SET, DELETE operations
- Persistence to disk (JSON or binary)
- Snapshot and restore
- Simple query language

**Requirements:**
- Thread-safe operations
- Write-ahead log (optional)
- Benchmark performance
- 90%+ test coverage

---

## Phase 1: Storage Layer Projects

### Project 1.1: Simple Page Manager
**Duration:** 8-10 hours
**Concepts:** File I/O, byte manipulation, caching

Implement a page manager that:
- Allocates 4KB pages
- Reads/writes pages to disk
- Tracks free pages with a bitmap
- Implements a simple LRU cache

**Test cases:**
- Allocate and write 1000 pages
- Read pages in random order
- Verify data integrity
- Measure cache hit rate

**Stretch goals:**
- Add page compression
- Implement wear leveling
- Add crash recovery

### Project 1.2: Buffer Pool Implementation
**Duration:** 12-15 hours
**Concepts:** LRU, mutexes, atomic operations

Build a thread-safe buffer pool:
- LRU eviction policy
- Pin/unpin mechanism
- Dirty page tracking
- Background flushing

**Test cases:**
- Concurrent access from multiple goroutines
- Pin count correctness
- Eviction correctness
- Race detector must pass

**Benchmarks:**
- Compare to no caching
- Measure lock contention
- Test different pool sizes

### Project 1.3: Write-Ahead Log
**Duration:** 10-12 hours
**Concepts:** Append-only files, fsync, recovery

Implement a WAL system:
- Append-only log entries
- Group commit optimization
- Crash-safe recovery
- Log truncation

**Test cases:**
- Simulate crashes at random points
- Verify idempotent recovery
- Test with large log files
- Measure fsync overhead

**Chaos testing:**
- Random crashes during writes
- Disk full scenarios
- Corrupted log entries

---

## Phase 2: Graph Structure Projects

### Project 2.1: CSR Graph Implementation
**Duration:** 15-18 hours
**Concepts:** Go 1.23 iterators, memory layout, cache locality

Build a CSR graph with:
- Efficient neighbor iteration using `iter.Seq`
- 2-hop query support
- Degree computation
- Iterator composition (filter, map)

**Test cases:**
- Build graph with 1M nodes, 10M edges
- Verify neighbor iteration correctness
- Test early exit with break
- Compare memory usage vs adjacency list

**Benchmarks:**
- Iteration speed vs slice return
- Cache miss rate with perf tools
- 2-hop query performance

### Project 2.2: Columnar Property Store
**Duration:** 12-15 hours
**Concepts:** Go 1.23 unique, compression, memory profiling

Implement columnar property storage:
- String interning with `unique.Handle`
- Bit-packed integers
- NULL bitmap
- Dictionary encoding

**Test cases:**
- Store 1M rows with mixed types
- Measure memory savings
- Verify compression correctness
- Distinct value counting

**Benchmarks:**
- Scan performance
- Filtering speed
- Memory usage vs row format

### Project 2.3: Parallel Graph Algorithms
**Duration:** 15-20 hours
**Concepts:** Worker pools, Go 1.25 synctest, parallelism

Implement parallel algorithms:
- Parallel BFS
- PageRank
- Triangle counting
- Connected components

**Test cases:**
- Use `testing/synctest` for determinism
- Verify correctness vs sequential
- Test with various graph sizes
- Check for race conditions

**Benchmarks:**
- Speedup vs sequential
- Scalability with core count
- Load balancing effectiveness

---

## Phase 3: Query Engine Projects

### Project 3.1: Expression Parser
**Duration:** 10-12 hours
**Concepts:** Recursive descent, AST, error messages

Build an expression parser for:
- Arithmetic expressions
- Comparison operators
- Boolean logic (AND, OR, NOT)
- Function calls

**Test cases:**
- Operator precedence
- Parentheses grouping
- Error messages with position
- Malformed input handling

**Stretch goals:**
- Add CASE expressions
- String operations
- Aggregate functions

### Project 3.2: Query Optimizer
**Duration:** 20-25 hours
**Concepts:** Cost estimation, DP, visitor pattern

Implement a cost-based optimizer:
- Join order optimization
- Filter pushdown
- Projection pushdown
- Cost estimation with statistics

**Test cases:**
- Compare plans for same query
- Verify optimizations are applied
- Test with different table sizes
- Edge cases (cross products, etc.)

**Benchmarks:**
- Optimization time
- Plan quality measurement
- Compare to unoptimized plans

### Project 3.3: Hash Join Implementation
**Duration:** 12-15 hours
**Concepts:** Go 1.24 Swiss Tables, hashing, join algorithms

Build efficient join operators:
- Hash join with pre-sizing
- Sort-merge join
- Index nested loop join
- Join type selection

**Test cases:**
- Correctness with duplicate keys
- NULL handling
- Empty input handling
- Large join verification

**Benchmarks:**
- Go 1.24 pre-sizing speedup
- Join algorithm comparison
- Memory usage profiling

### Project 3.4: Pipelined Executor
**Duration:** 18-22 hours
**Concepts:** Go 1.23 iterators, pipelines, profiling

Implement execution engine:
- Iterator-based operators
- Pipeline breakers (sort, hash join)
- Operator profiling
- Adaptive execution

**Test cases:**
- Complex query plans
- Early exit verification
- Pipeline vs materialization
- Memory usage tracking

**Benchmarks:**
- End-to-end query performance
- Operator-level profiling
- Compare to other engines (optional)

---

## Phase 4: Transactions Projects

### Project 4.1: Lock Manager
**Duration:** 15-18 hours
**Concepts:** Locks, deadlock detection, Go 1.25 synctest

Implement a lock manager:
- Shared and exclusive locks
- Lock upgrades/downgrades
- Deadlock detection with wait-for graph
- Two-phase locking enforcement

**Test cases:**
- Use `testing/synctest` for determinism
- Deadlock detection accuracy
- Lock compatibility matrix
- Concurrent transaction workloads

**Stress tests:**
- 1000 concurrent transactions
- High contention scenarios
- Deadlock resolution

### Project 4.2: MVCC Implementation
**Duration:** 25-30 hours
**Concepts:** Version chains, Go 1.24 weak pointers, GC

Build an MVCC system:
- Version chain management
- Snapshot isolation
- Write conflict detection
- Garbage collection with `weak.Pointer`

**Test cases:**
- Concurrent read/write workloads
- Snapshot consistency
- Write conflict detection
- GC correctness

**Benchmarks:**
- MVCC vs locking throughput
- Read-heavy vs write-heavy workloads
- GC overhead measurement

---

## Integration Projects

### Integration Project 1: Mini Graph Database
**Duration:** 40-50 hours
**Combines:** All phases

Build a minimal graph database:
- Node and edge storage (Phase 1)
- CSR graph representation (Phase 2)
- Simple query language (Phase 3)
- Transaction support (Phase 4)

**Features:**
- CREATE nodes and edges
- MATCH pattern queries
- WHERE filtering
- RETURN projection
- BEGIN/COMMIT transactions

**Test cases:**
- End-to-end query tests
- Concurrent transaction tests
- Crash recovery
- Performance benchmarks

### Integration Project 2: Social Network Analyzer
**Duration:** 30-40 hours
**Combines:** Phases 2-3

Build a social network analyzer:
- Load graph from CSV
- Friend recommendation (2-hop)
- Community detection
- Influence analysis (PageRank)

**Features:**
- Interactive CLI
- Query statistics
- Visualization (optional)
- Export results

### Integration Project 3: Performance Benchmarking Suite
**Duration:** 20-25 hours
**Combines:** All phases

Create comprehensive benchmarks:
- Micro-benchmarks for each component
- End-to-end query benchmarks
- Scalability tests
- Comparison to other systems

**Metrics:**
- Throughput (queries/sec)
- Latency (p50, p95, p99)
- Memory usage
- Scalability graphs

---

## Debug Scenarios

### Scenario 1: Race Condition Hunt
**Duration:** 3-5 hours
**Concepts:** Race detector, synctest

Fix intentional race conditions in:
- Buffer pool implementation
- Concurrent hash join
- Lock manager

**Tools:**
- `go test -race`
- `testing/synctest`
- pprof for goroutine profiling

### Scenario 2: Memory Leak Detection
**Duration:** 3-5 hours
**Concepts:** Memory profiling, weak pointers

Find and fix memory leaks in:
- Version chain GC
- Buffer pool eviction
- Query result sets

**Tools:**
- `pprof` memory profiling
- Heap dump analysis
- GC trace analysis

### Scenario 3: Performance Regression
**Duration:** 4-6 hours
**Concepts:** Benchmarking, profiling, optimization

Identify why a query got slower:
- Profile with pprof
- Identify hot paths
- Optimize bottlenecks
- Verify improvements

---

## Code Review Exercises

### Exercise 1: Buffer Pool Review
Review this buffer pool implementation and identify:
- Race conditions
- Memory leaks
- Performance issues
- Missing edge cases

### Exercise 2: Query Optimizer Review
Review this optimizer and identify:
- Incorrect cost estimates
- Missing optimizations
- Corner cases
- Code quality issues

### Exercise 3: MVCC Review
Review this MVCC implementation and identify:
- Visibility bugs
- Write conflict issues
- GC problems
- Concurrency issues

---

## Challenge Problems

### Challenge 1: Implement Adaptive Radix Tree
**Difficulty:** Hard
**Duration:** 30-40 hours

Implement an ART index:
- Node types (Node4, Node16, Node48, Node256)
- Lazy expansion
- Prefix compression
- Iterator support

### Challenge 2: Vectorized Execution
**Difficulty:** Hard
**Duration:** 25-35 hours

Implement vectorized operators:
- Batch-at-a-time processing
- SIMD operations (optional)
- Type-specific code generation
- Benchmark improvements

### Challenge 3: Distributed Transactions
**Difficulty:** Very Hard
**Duration:** 50-60 hours

Implement 2PC for distributed transactions:
- Transaction coordinator
- Prepare/commit protocol
- Failure recovery
- Distributed deadlock detection

---

## How to Use These Exercises

### For Self-Study:
1. Complete projects in order
2. Aim for 90%+ test coverage
3. Run benchmarks and profile
4. Compare to reference implementations

### For Instructors:
1. Assign projects as homework
2. Use debug scenarios in class
3. Code review exercises for discussion
4. Challenge problems for advanced students

### For Teams:
1. Pair programming on projects
2. Code review each other's work
3. Competition on benchmark performance
4. Collaborate on integration projects

---

## Reference Solutions

Reference solutions are available separately. Try to complete exercises without looking at solutions first!

**Learning path:** Try → Struggle → Research → Implement → Review → Refine

**Remember:** The goal is learning, not just completing. Take your time and understand each concept deeply!

---

## Additional Resources

### Books:
- "Database Internals" by Alex Petrov
- "Designing Data-Intensive Applications" by Martin Kleppmann
- "The Go Programming Language" by Donovan & Kernighan

### Papers:
- "The Case for Shared Nothing" (Stonebraker)
- "Access Path Selection in a Relational Database Management System" (Selinger et al.)
- "Serializable Snapshot Isolation in PostgreSQL" (Ports & Grittner)

### Online Courses:
- CMU 15-445/645: Database Systems
- Stanford CS245: Principles of Data-Intensive Systems
- MIT 6.824: Distributed Systems

### Open Source Projects to Study:
- PostgreSQL (transaction system)
- SQLite (B-tree implementation)
- DuckDB (vectorized execution)
- MemGraph (graph database in C++)

---

**Happy coding! Remember: The best way to learn is by doing!** 🚀
