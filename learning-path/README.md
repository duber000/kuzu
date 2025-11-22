# Complete Go Learning Path: From Zero to Graph Database

**A Structured Journey from Complete Beginner to Database Implementation**

## Overview

This learning path takes you from zero Go programming knowledge to building a complete graph database implementation. It's designed to work alongside the `complete-graph-database-learning-path-go1.25.md` curriculum, providing the Go knowledge foundation you need to succeed.

**Total Duration:** 18-20 weeks for complete beginners | 13-14 weeks for experienced programmers

## Who Is This For?

### Track A: Complete Beginners (You)
- Little to no programming experience, or
- Experience in other languages but new to systems programming
- **Start here:** Pre-work Week 1

### Track B: Some Programming Experience
- Comfortable with at least one programming language
- Understand basic data structures and algorithms
- **Start here:** Pre-work Week 3 (skip to intermediate)

### Track C: Experienced Programmers
- Strong background in systems programming (C, C++, Rust, etc.)
- Understand concurrency and I/O
- **Start here:** Week 5-6 (just learn Go-specific idioms)

## Learning Path Structure

```
┌─────────────────────────────────────────────────────────────┐
│                    PRE-WORK (Weeks 1-6)                     │
│              Build Go Foundation Before Lesson 1.1           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              PHASE 1: Storage Layer (Weeks 7-9)             │
│         Pages → Buffer Pool → Write-Ahead Log               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│           PHASE 2: Graph Structure (Weeks 10-12)            │
│              CSR → Columnar → Parallelism                   │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│            PHASE 3: Query Engine (Weeks 13-16)              │
│         Parser → Planner → Joins → Executor                 │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│           PHASE 4: Transactions (Weeks 17-18)               │
│              Locking Protocols → MVCC                        │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│          BONUS: Advanced Topics (Weeks 19-20)               │
│      Optimization, Distribution, Production Hardening        │
└─────────────────────────────────────────────────────────────┘
```

## Pre-Work Phase (Weeks 1-6)

**Critical:** Do NOT skip this if you're new to Go. The main curriculum assumes Go proficiency.

| Week | Module | What You'll Learn | Practice Project |
|------|--------|-------------------|------------------|
| 1-2 | [Go Fundamentals](pre-work/week-1-2-go-fundamentals.md) | Syntax, types, functions, structs, interfaces | CLI TODO list app |
| 3-4 | [Intermediate Concepts](pre-work/week-3-4-intermediate-concepts.md) | Slices, maps, pointers, testing, benchmarking | Key-value store with persistence |
| 5-6 | [Concurrency & I/O](pre-work/week-5-6-concurrency-io.md) | Goroutines, channels, mutexes, file I/O | Concurrent file downloader |

**Checkpoint:** After Week 6, you should be able to:
- ✅ Build a CLI application with file persistence
- ✅ Write tests and benchmarks
- ✅ Use goroutines without data races
- ✅ Understand pointers and memory layout

## Main Curriculum Phases

Each phase has two components:
1. **Go Preparation Lessons** (this repository) - Learn the Go concepts needed
2. **Database Implementation** (main curriculum) - Build the actual database features

### Phase 1: Storage Layer (Weeks 7-9)

| Lesson | Go Prep | Database Implementation | Hours |
|--------|---------|------------------------|-------|
| 1.1 | [Pages Prep](phase-1-storage/go-prep-lesson-1-1-pages.md) | Pages & Memory-Mapped I/O | 15-20 |
| 1.2 | [Buffer Pool Prep](phase-1-storage/go-prep-lesson-1-2-buffer-pool.md) | LRU Cache & Concurrency | 15-20 |
| 1.3 | [WAL Prep](phase-1-storage/go-prep-lesson-1-3-wal.md) | Write-Ahead Log & Durability | 15-20 |

**New Go Concepts:**
- `syscall.Mmap` for memory-mapped files
- Doubly-linked lists with maps (LRU)
- `atomic` package for lock-free counters
- **Go 1.23:** Timer improvements
- **Go 1.24:** Weak pointers (optional)

### Phase 2: Graph Structure (Weeks 10-12)

| Lesson | Go Prep | Database Implementation | Hours |
|--------|---------|------------------------|-------|
| 2.1 | [CSR Prep](phase-2-graph/go-prep-lesson-2-1-csr.md) | Compressed Sparse Row & Iterators | 15-20 |
| 2.2 | [Columnar Prep](phase-2-graph/go-prep-lesson-2-2-columnar.md) | Column Storage & Compression | 15-20 |
| 2.3 | [Parallelism Prep](phase-2-graph/go-prep-lesson-2-3-parallelism.md) | Work Stealing & Morsel-Driven | 15-20 |

**New Go Concepts (CRITICAL):**
- **Go 1.23:** `iter.Seq` and range-over-func iterators ⭐⭐⭐
- **Go 1.23:** `unique` package for string interning ⭐⭐
- Work stealing with channels
- **Go 1.25:** `testing/synctest` for deterministic tests ⭐⭐

### Phase 3: Query Engine (Weeks 13-16)

| Lesson | Go Prep | Database Implementation | Hours |
|--------|---------|------------------------|-------|
| 3.1 | [Parser Prep](phase-3-query/go-prep-lesson-3-1-parser.md) | Cypher Parser & AST | 15-20 |
| 3.2 | [Planner Prep](phase-3-query/go-prep-lesson-3-2-planner.md) | Query Optimization | 15-20 |
| 3.3 | [Joins Prep](phase-3-query/go-prep-lesson-3-3-joins.md) | Hash Join & Sort-Merge | 15-20 |
| 3.4 | [Executor Prep](phase-3-query/go-prep-lesson-3-4-executor.md) | Pipeline Execution | 15-20 |

**New Go Concepts:**
- **Go 1.24:** Swiss Tables (30-35% faster hash joins) ⭐⭐⭐
- Iterator-based operator composition
- **Go 1.24:** `testing.B.Loop` for benchmarks

### Phase 4: Transactions (Weeks 17-18)

| Lesson | Go Prep | Database Implementation | Hours |
|--------|---------|------------------------|-------|
| 4.1 | [Locking Prep](phase-4-transactions/go-prep-lesson-4-1-locking.md) | Two-Phase Locking | 15-20 |
| 4.2 | [MVCC Prep](phase-4-transactions/go-prep-lesson-4-2-mvcc.md) | Multi-Version Concurrency | 15-20 |

**New Go Concepts:**
- Deadlock detection algorithms
- **Go 1.25:** `synctest` for deadlock testing ⭐⭐
- **Go 1.24:** `weak.Pointer` for version chains

## How to Use This Learning Path

### Daily Workflow

**For each lesson:**

1. **Morning (1-2 hours): Study Go concepts**
   - Read the "Go Prep" lesson
   - Complete the code examples
   - Run the experiments

2. **Afternoon (2-3 hours): Build the database feature**
   - Follow the main curriculum lesson
   - Apply the Go concepts you learned
   - Implement the required features

3. **Evening (1 hour): Verify & Reflect**
   - Run tests and benchmarks
   - Write your learning journal
   - Complete the checkpoint questions

### Weekly Pattern

**Monday-Wednesday:** New lesson
**Thursday-Friday:** Finish implementation and experiments
**Weekend:** Review, debug, and prepare for next week

## Required Setup

### Install Go 1.25

```bash
# Download from https://go.dev/dl/
# Or use version manager
go install golang.org/dl/go1.25@latest
go1.25 download

# Verify
go version  # Should show go1.25 or later
```

### Optional: Enable Experimental Features

```bash
# Add to your shell profile (~/.bashrc or ~/.zshrc)
export GOEXPERIMENT=greenteagc  # Experimental GC for better graph performance
```

### IDE Setup

Choose one:
- **VS Code** with Go extension (recommended for beginners)
- **GoLand** by JetBrains (powerful but paid)
- **Vim/Neovim** with vim-go (for advanced users)

### Create Your Project

```bash
cd ~
mkdir kuzu-go-learning
cd kuzu-go-learning
go mod init github.com/yourusername/kuzu-go

# Create directory structure
mkdir -p pre-work/{week1-2,week3-4,week5-6}
mkdir -p {storage,graph,query,index,transaction}
mkdir -p tests
```

## Progress Tracking

### Self-Assessment Checklist

After **Pre-work** (Week 6):
- [ ] I can write a CLI application in Go
- [ ] I understand slices, maps, and pointers
- [ ] I can write tests with `go test`
- [ ] I can use goroutines without data races
- [ ] I can benchmark code with `go test -bench`

After **Phase 1** (Week 9):
- [ ] I understand memory-mapped I/O
- [ ] I can implement concurrent data structures
- [ ] I know when to use mutexes vs atomic operations
- [ ] I understand durability and fsync

After **Phase 2** (Week 12):
- [ ] I can write custom iterators with `iter.Seq`
- [ ] I understand string interning with `unique`
- [ ] I can write parallel code that scales
- [ ] I can use `synctest` for deterministic tests

After **Phase 3** (Week 16):
- [ ] I understand Swiss Tables performance improvements
- [ ] I can design pipeline-based systems
- [ ] I can profile and optimize hot paths

After **Phase 4** (Week 18):
- [ ] I can implement complex concurrency protocols
- [ ] I understand Go's memory model deeply
- [ ] I can write production-quality concurrent code

### Learning Journal Template

After each lesson, write:

```markdown
# Lesson X.Y - [Topic] - [Date]

## What I Learned
- Key concept 1
- Key concept 2

## What Surprised Me
[The thing you didn't expect]

## What Was Hard
[Where you got stuck and how you solved it]

## Performance Insights
[What optimization made the biggest difference]

## Go 1.23+ Features
[How modern Go features helped]

## Next Steps
[What to focus on tomorrow]
```

## Additional Resources

### Go Learning Resources

**Fundamentals:**
- [Tour of Go](https://go.dev/tour/) - Interactive tutorial (4-6 hours)
- [Go by Example](https://gobyexample.com/) - Quick reference
- [Effective Go](https://go.dev/doc/effective_go) - Idiomatic patterns

**Concurrency:**
- [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs) - Rob Pike
- [Advanced Go Concurrency](https://www.youtube.com/watch?v=QDDwwePbDtw) - Sameer Ajmani

**Modern Go (1.23-1.25):**
- [Go 1.23 Release Notes](https://go.dev/doc/go1.23)
- [Go 1.24 Release Notes](https://go.dev/doc/go1.24)
- [Go 1.25 Release Notes](https://go.dev/doc/go1.25)
- [Range-over-func Proposal](https://go.dev/wiki/RangefuncExperiment)

**Performance:**
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/gopherchina-2019.html)

### Database Learning Resources

**Books:**
- *Database Internals* by Alex Petrov
- *Designing Data-Intensive Applications* by Martin Kleppmann

**Courses:**
- [CMU 15-445: Database Systems](https://15445.courses.cs.cmu.edu/)
- [CMU 15-721: Advanced Database Systems](https://15721.courses.cs.cmu.edu/)

**Codebases:**
- [DuckDB](https://github.com/duckdb/duckdb) - Similar embedded DB
- [Dgraph](https://github.com/dgraph-io/dgraph) - Production graph DB in Go

### Community & Support

**Go Community:**
- [Gophers Slack](https://invite.slack.golangbridge.org/)
- [r/golang](https://reddit.com/r/golang)
- [Go Forum](https://forum.golangbridge.org/)

**Database Community:**
- [Database Internals Discord](https://discord.gg/database-internals)
- [r/databasedevelopment](https://reddit.com/r/databasedevelopment)

## Troubleshooting & FAQs

### "I'm stuck on pre-work Week 1. Is this too hard?"

No! Week 1-2 introduces many new concepts. Common issues:
- Spend more time on the Tour of Go
- Watch beginner YouTube tutorials
- Join Gophers Slack #newbies channel
- It's OK to take 3-4 weeks on fundamentals

### "Should I skip the pre-work?"

**Only if you can:**
- Write a working CLI tool in Go from scratch
- Explain the difference between `make([]int, 5)` and `make([]int, 0, 5)`
- Use goroutines and channels correctly
- Write tests and benchmarks

Test yourself: Try the Week 5-6 practice project. If you can do it in 4-6 hours, skip pre-work.

### "Which Go version should I use?"

**Go 1.25** - Get the latest features (iterators, synctest, etc.)

If you must use older versions:
- Go 1.24: Skip synctest examples, use older testing patterns
- Go 1.23: Skip Swiss Tables benefits, weak pointers
- Go 1.22 or older: Not recommended, you'll miss critical features

### "How much time should I dedicate?"

**Minimum:** 10 hours/week (will take 24-30 weeks)
**Recommended:** 15 hours/week (18-20 weeks)
**Intensive:** 20-25 hours/week (12-14 weeks)

Less than 10 hours/week: You'll forget concepts between sessions.

### "I'm an experienced programmer. Can I go faster?"

Yes! Suggested accelerated path:
- Week 1: Complete all pre-work (40 hours)
- Weeks 2-4: Phase 1 (45 hours)
- Weeks 5-7: Phase 2 (45 hours)
- Weeks 8-11: Phase 3 (60 hours)
- Weeks 12-13: Phase 4 (30 hours)
- Week 14: Polish and bonus challenges

Total: 14 weeks at 15-20 hours/week

## Success Stories & Testimonials

> "I started with zero Go knowledge. After 6 weeks of pre-work, I felt confident starting the storage layer. By week 12, I built my first working query engine. This is challenging but incredibly rewarding."
> — Learning Path Beta Tester

> "The Go 1.23 iterator sections were mind-blowing. I finally understood why functional programming patterns are powerful. The CSR lesson clicked once I mastered iterators."
> — Experienced Developer New to Go

> "Synctest made concurrency testing actually fun. No more flaky tests!"
> — Systems Programmer

## Next Steps

**Ready to start?**

1. **Complete Beginner:** Start with [Week 1-2: Go Fundamentals](pre-work/week-1-2-go-fundamentals.md)
2. **Some Experience:** Jump to [Week 3-4: Intermediate Concepts](pre-work/week-3-4-intermediate-concepts.md)
3. **Experienced:** Test yourself with [Week 5-6: Concurrency & I/O](pre-work/week-5-6-concurrency-io.md)

**Questions?**
- Open an issue in this repository
- Ask in Gophers Slack #database-learning channel
- Check the troubleshooting section above

---

**Last Updated:** November 2025
**Go Version:** 1.25+
**Curriculum Version:** Based on `complete-graph-database-learning-path-go1.25.md`

**License:** MIT - Use freely for learning

---

Let's build something amazing! 🚀
