# Debug Scenario 2: Memory Leak Detection

## Overview
Find and fix memory leaks in version chain GC, buffer pool eviction, and query result sets.

**Duration:** 3-5 hours
**Difficulty:** Medium

## Scenarios

### Scenario 1: Version Chain Leak
**Bug:** Old versions not cleaned up in MVCC
```go
// BUGGY CODE
func (s *MVCCStore) CreateVersion(key Key, value Value) {
	newVersion := &Version{
		data: value,
		prev: s.latest[key],  // LEAK: circular reference
	}
	s.latest[key] = newVersion
	// Missing: GC of old versions!
}
```

**Symptoms:**
- Memory grows unbounded
- Heap profiles show version accumulation
- OOM on long-running workloads

### Scenario 2: Buffer Pool Leak
**Bug:** Evicted pages not freed
```go
// BUGGY CODE
func (bp *BufferPool) Evict(frameID int) {
	frame := bp.frames[frameID]
	delete(bp.pageTable, frame.pageID)
	// LEAK: frame still in frames array!
}
```

**Symptoms:**
- Memory usage doesn't decrease
- Cache size grows beyond limit
- Gradual performance degradation

### Scenario 3: Query Result Leak
**Bug:** Result sets not closed
```go
// BUGGY CODE
func (db *DB) ExecuteQuery(sql string) error {
	rs := db.query(sql)
	for rs.Next() {
		// process results
	}
	// LEAK: rs.Close() never called!
	return nil
}
```

**Symptoms:**
- Goroutine leaks
- File descriptor exhaustion
- Memory growth with queries

## Debugging Tools

```bash
# Memory profiling
go test -memprofile=mem.prof
go tool pprof mem.prof

# Heap analysis
go tool pprof -alloc_space mem.prof
go tool pprof -inuse_space mem.prof

# GC trace
GODEBUG=gctrace=1 go test

# Leak detection
go test -run=TestLeak -memprofile=leak.prof
```

## pprof Commands
```
(pprof) top10          # Show top allocators
(pprof) list FuncName  # Show allocations in function
(pprof) web            # Visualize call graph
(pprof) pdf            # Generate PDF report
```

## Fixes
1. Implement proper GC in MVCC
2. Clear references in buffer pool
3. Add defer for cleanup
4. Use weak.Pointer for GC-friendly refs

## Learning Objectives
- Profile memory usage with pprof
- Identify memory leak patterns
- Implement proper cleanup
- Use weak pointers effectively

## Time Estimate
Each scenario: 1-2 hours
