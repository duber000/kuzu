# Debug Scenario 3: Performance Regression

## Overview
Identify and fix performance regressions using profiling tools and optimization techniques.

**Duration:** 4-6 hours
**Difficulty:** Medium-Hard

## Scenario
A query that used to run in 100ms now takes 2 seconds. Find the bottleneck and fix it.

## Regression Sources
1. **Inefficient query plan** - Wrong join order
2. **Lock contention** - Too much locking
3. **Memory allocation** - Excessive allocations
4. **Cache thrashing** - Poor cache utilization

## Debugging Steps

### 1. Profile CPU
```bash
go test -bench=BenchmarkQuery -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

### 2. Profile Memory
```bash
go test -bench=BenchmarkQuery -memprofile=mem.prof
go tool pprof -alloc_objects mem.prof
```

### 3. Profile Locks
```bash
go test -bench=BenchmarkQuery -mutexprofile=mutex.prof
go tool pprof mutex.prof
```

### 4. Trace Execution
```bash
go test -bench=BenchmarkQuery -trace=trace.out
go tool trace trace.out
```

## Common Issues

### Issue 1: Hot Path Allocation
**Problem:**
```go
// SLOW: allocates every call
func (e *Executor) Execute() []Row {
	results := make([]Row, 0)  // allocation
	for row := range e.scan() {
		results = append(results, row)
	}
	return results
}
```

**Fix:**
```go
// FAST: pre-allocate
func (e *Executor) Execute() []Row {
	results := make([]Row, 0, e.estimatedRows)
	// ...
}
```

### Issue 2: Lock Contention
**Problem:**
```go
// SLOW: single global lock
var mu sync.Mutex
func (c *Cache) Get(key string) {
	mu.Lock()
	defer mu.Unlock()
	// ...
}
```

**Fix:**
```go
// FAST: shard locks
type Cache struct {
	shards [16]shard
}
func (c *Cache) Get(key string) {
	shard := &c.shards[hash(key) % 16]
	shard.mu.Lock()
	defer shard.mu.Unlock()
	// ...
}
```

### Issue 3: Inefficient Algorithm
**Problem:** Wrong join algorithm selected
**Fix:** Use cost-based optimizer

## Tools Reference

### pprof
```
top10              # Top CPU consumers
list FuncName      # Source code view
web               # Call graph
pdf               # PDF report
```

### trace
- View goroutine timeline
- Identify blocking operations
- Find GC pauses

### benchstat
```bash
go test -bench=. > old.txt
# Make changes
go test -bench=. > new.txt
benchstat old.txt new.txt
```

## Performance Goals
After fixes:
- Query time: <100ms (20x faster)
- CPU usage: <50% of before
- Allocations: <1000 per query
- Lock contention: <5%

## Learning Objectives
- Profile Go programs effectively
- Identify performance bottlenecks
- Apply optimization techniques
- Validate improvements with benchmarks

## Time Estimate
Investigation: 2-3 hours, Fixes: 2-3 hours
