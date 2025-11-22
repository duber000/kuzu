# Learn Go by Building a Graph Database

**A comprehensive, hands-on curriculum for learning Go by building an embeddable graph database inspired by Kuzu DB.**

---

## Overview

This learning path takes you from **"Hello World" in Go** to building a **fully-functional graph database** with transactions, query optimization, and parallel execution. No prior Go experience required!

**What makes this different:**
- **Learn by doing:** Every lesson is a hands-on challenge
- **Complete progression:** From basic syntax to advanced database internals
- **Modern Go:** Uses Go 1.23+ features (iterators, unique package, synctest)
- **Discovery-based:** You figure things out through experimentation, not by reading solutions
- **Real-world project:** Build something substantial you can show off

**Time commitment:** 12-16 weeks (part-time)

---

## Learning Path Structure

The curriculum is organized into 5 phases:

### Phase 0: Go Fundamentals (2-3 weeks)
**Start here if you're new to Go!**

Learn the basics of Go programming:
- Hello World & basic syntax
- Functions, packages, and modules
- Structs, interfaces, and methods
- Concurrency with goroutines and channels
- Error handling and testing
- File I/O and system calls

📖 **[Start Phase 0: Go Fundamentals](phase-0-go-fundamentals.md)**

---

### Phase 1: Storage Layer (3 weeks)
**Prerequisites:** Phase 0

Build the foundation of your database:
- Page-based storage abstraction
- Buffer pool with LRU eviction
- Write-ahead logging (WAL) for crash recovery

**What you'll learn:**
- How databases store data on disk
- Memory management and caching strategies
- Ensuring durability and crash safety

📖 **[Go to Phase 1: Storage Layer](phase-1-storage-layer.md)**

---

### Phase 2: Graph Structure (3 weeks)
**Prerequisites:** Phase 1

Implement efficient graph representations:
- Compressed Sparse Row (CSR) format
- Columnar storage for node properties
- Parallel query execution with NodeGroups

**What you'll learn:**
- Memory-efficient graph representations
- Cache-friendly data structures
- Parallel processing patterns in Go

📖 **[Go to Phase 2: Graph Structure](phase-2-graph-structure.md)**

---

### Phase 3: Query Engine (4 weeks)
**Prerequisites:** Phase 2

Build a complete query processor:
- Parse Cypher queries into AST
- Query optimization and planning
- Join algorithms (hash, sort-merge, WCO)
- Iterator-based execution engine

**What you'll learn:**
- Compiler design (lexing, parsing)
- Query optimization techniques
- Join algorithms and their tradeoffs
- Execution models (iterator vs vectorized)

📖 **[Go to Phase 3: Query Engine](phase-3-query-engine.md)**

---

### Phase 4: Transactions (2 weeks)
**Prerequisites:** Phase 3

Add transaction support:
- Two-phase locking (2PL)
- Deadlock detection
- Multi-version concurrency control (MVCC)
- Isolation levels

**What you'll learn:**
- Concurrency control mechanisms
- ACID properties implementation
- Performance vs correctness tradeoffs
- Testing concurrent systems

📖 **[Go to Phase 4: Transactions](phase-4-transactions.md)**

---

## Quick Start

### Absolute Beginners

**Never used Go before?** Start here:

1. **Install Go 1.23+:** Download from [go.dev/dl](https://go.dev/dl/)
2. **Verify installation:**
   ```bash
   go version  # Should show go1.23 or later
   ```
3. **Start with Phase 0:** [Go Fundamentals](phase-0-go-fundamentals.md)

### Experienced Go Developers

**Already know Go?** Skip to the database content:

1. **Review the prerequisites** in Phase 0 (especially file I/O and concurrency)
2. **Jump to Phase 1:** [Storage Layer](phase-1-storage-layer.md)
3. **Skim or skip** the Go syntax lessons

---

## What You'll Build

By the end of this curriculum, you'll have built:

```
kuzu-go/
├── storage/          # Page storage, buffer pool, WAL
├── graph/            # CSR format, columnar storage
├── query/            # Parser, planner, executor
├── transaction/      # Locking, MVCC
└── cmd/kuzu/         # CLI interface
```

**Capabilities:**
- Store millions of nodes and edges efficiently
- Parse and execute Cypher queries
- Optimize query plans automatically
- Execute queries in parallel across CPU cores
- Support concurrent transactions with ACID guarantees
- Recover from crashes without data loss

---

## Learning Philosophy

This curriculum follows a **discovery-based approach**:

### Pattern: Goal → Investigation → Implementation → Verification → Reflection

1. **Goal:** Each lesson presents a concrete challenge
2. **Investigation:** You explore through questions and experiments
3. **Implementation:** You build the solution yourself
4. **Verification:** You benchmark and test your code
5. **Reflection:** You document what you learned

### Key Principles

- **Measure, don't guess:** Use benchmarks and profiling
- **Break things:** Make it fail, understand why, fix it
- **Iterate:** First version won't be perfect—that's OK
- **No solutions provided:** You discover through experimentation
- **Ask "why":** Understand the reasoning behind design decisions

---

## Modern Go Features

This curriculum uses cutting-edge Go 1.23-1.25 features:

### Go 1.23+
- **Range-over-func iterators:** Clean, composable traversals
- **`unique` package:** Memory-efficient string interning
- **Timer improvements:** Simpler background tasks

### Go 1.24+
- **Swiss Tables:** 30-35% faster hash maps
- **Weak pointers:** Memory-aware caching
- **`testing.B.Loop`:** Cleaner benchmark syntax

### Go 1.25+
- **`testing/synctest`:** Deterministic concurrency tests (no more flaky tests!)
- **Container-aware GOMAXPROCS:** Automatic CPU limit detection
- **Experimental GC:** 10-35% improvement for graph workloads

**Version compatibility matrix:** See individual phase files for details.

---

## Prerequisites

### Required Knowledge
- Basic programming experience (any language)
- Comfortable with command line
- Willingness to learn and experiment

### Optional (Helpful)
- Basic understanding of databases (tables, queries)
- Familiarity with algorithms and data structures
- Some systems programming knowledge

### Tools You'll Need
- Go 1.23 or later
- Text editor or IDE (VS Code, GoLand, etc.)
- Git (for version control)
- Terminal/command line

---

## Recommended Resources

### Go Learning
- **Official Tour:** https://go.dev/tour/
- **Go by Example:** https://gobyexample.com/
- **Effective Go:** https://go.dev/doc/effective_go
- **The Go Programming Language** by Donovan & Kernighan (book)

### Database Internals
- **Database Internals** by Alex Petrov (book)
- **Designing Data-Intensive Applications** by Martin Kleppmann (book)
- **CMU 15-445: Database Systems** (YouTube course)
- **CMU 15-721: Advanced Database Systems** (YouTube course)

### Papers
- Kuzu's CIDR 2023 paper
- "Worst-Case Optimal Join Algorithms"
- "ARIES: A Transaction Recovery Method"

### Reference Implementations
- **Kuzu:** https://github.com/kuzudb/kuzu
- **DuckDB:** https://github.com/duckdb/duckdb
- **SQLite:** https://sqlite.org/
- **Dgraph:** https://github.com/dgraph-io/dgraph (Go graph database)

---

## Project Structure Recommendation

```
kuzu-go/
├── cmd/
│   └── kuzu/
│       └── main.go              # CLI interface
├── storage/
│   ├── page.go                  # Page abstraction
│   ├── buffer_pool.go           # Buffer manager
│   ├── wal.go                   # Write-ahead log
│   └── storage_test.go
├── graph/
│   ├── csr.go                   # Compressed Sparse Row
│   ├── column_store.go          # Columnar properties
│   ├── node_group.go            # NodeGroup parallelism
│   └── graph_test.go
├── query/
│   ├── lexer.go                 # Tokenizer
│   ├── parser.go                # AST builder
│   ├── planner.go               # Query optimizer
│   ├── executor.go              # Execution engine
│   ├── operators.go             # Physical operators
│   └── query_test.go
├── index/
│   ├── hash.go                  # Hash index
│   └── btree.go                 # B-tree (optional)
├── transaction/
│   ├── lock_manager.go          # 2PL implementation
│   ├── mvcc.go                  # Multi-version CC
│   └── transaction_test.go
├── go.mod                       # go 1.23+
├── go.sum
└── README.md
```

---

## Getting Help

### Stuck on a Lesson?

1. **Re-read the Investigation Questions:** They contain hints
2. **Run the experiments:** Measure and observe behavior
3. **Check the benchmark targets:** Are you close?
4. **Review the "What You Should Discover"** section
5. **Study reference implementations:** See how others solved it
6. **Ask the community:** Go forums, Reddit, Discord

### Common Pitfalls

- **Skipping measurements:** Always benchmark before optimizing
- **Ignoring test coverage:** Use `go test -race -cover`
- **Copying code without understanding:** Type it out, experiment with it
- **Perfectionism:** First version doesn't need to be perfect
- **Not taking breaks:** This is a marathon, not a sprint

---

## Progress Tracking

After each lesson, write a brief reflection:

1. **What surprised you?**
2. **What was harder than expected?**
3. **What optimization made the biggest difference?**
4. **What would you do differently next time?**
5. **How did modern Go features help?**

Keep these in a journal or blog. Reflection solidifies learning.

---

## Community

### Share Your Progress

- Blog about your journey
- Share code on GitHub
- Present at local Go meetups
- Help others who are following this path

### Contribute

Found a mistake? Have an improvement?

- This curriculum is open to feedback
- Share your experiences
- Suggest additional challenges
- Create supplementary materials

---

## Frequently Asked Questions

**Q: Do I need to know C++ to understand database internals?**
A: No! This curriculum teaches everything in Go from scratch.

**Q: How long does this take?**
A: 12-16 weeks part-time (10-15 hours/week). Full-time students can complete it faster.

**Q: Can I use an older Go version?**
A: Go 1.23+ is recommended for modern features. Older versions will work but you'll miss iterators, Swiss Tables, etc.

**Q: Is this suitable for production use?**
A: This is a learning project. For production, use established databases like Kuzu, DuckDB, or PostgreSQL.

**Q: What if I get stuck?**
A: Review the investigation questions, run experiments, check reference implementations, ask the community.

**Q: Can I skip phases?**
A: Each phase builds on previous ones. Skipping is not recommended unless you're already familiar with the concepts.

**Q: Do I need a powerful computer?**
A: No. Any modern laptop will work fine. Some benchmarks might be slower, but concepts are the same.

---

## License

This learning curriculum is open source. Feel free to use, modify, and share.

Inspired by:
- Kuzu DB (https://kuzudb.com/)
- "Build Your Own X" philosophy
- Discovery-based learning principles

---

## Ready to Start?

### Beginners: Start with Phase 0
📖 **[Begin Phase 0: Go Fundamentals](phase-0-go-fundamentals.md)**

### Experienced Go Devs: Jump to Phase 1
📖 **[Begin Phase 1: Storage Layer](phase-1-storage-layer.md)**

---

**Happy coding! Let's build a database! 🚀**
