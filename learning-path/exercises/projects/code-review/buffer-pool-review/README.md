# Code Review Exercise 1: Buffer Pool Review

## Overview
Review a buffer pool implementation and identify issues in concurrency, correctness, performance, and code quality.

**Duration:** 2-3 hours
**Difficulty:** Medium

## Review Checklist

### 1. Race Conditions
- [ ] Pin count updates are atomic/protected
- [ ] Page table access is synchronized
- [ ] Free list modifications are safe
- [ ] LRU replacer is thread-safe

### 2. Memory Leaks
- [ ] Evicted pages are properly freed
- [ ] Background goroutines are cleaned up
- [ ] No circular references
- [ ] Resources released on Close()

### 3. Performance Issues
- [ ] Lock granularity is appropriate
- [ ] No unnecessary allocations in hot paths
- [ ] LRU updates are efficient
- [ ] Cache-friendly data layout

### 4. Missing Edge Cases
- [ ] All frames pinned (no eviction possible)
- [ ] Invalid page IDs handled
- [ ] Concurrent flush and eviction
- [ ] Double unpin detected

### 5. Code Quality
- [ ] Clear variable names
- [ ] Proper error handling
- [ ] Comments for complex logic
- [ ] Consistent style

## Common Issues to Find

### Issue 1: Pin Count Race
```go
// WRONG
frame.pinCount++

// CORRECT
frame.pinCount.Add(1)
```

### Issue 2: Holding Lock During I/O
```go
// WRONG: Holds lock during disk read
func (bp *BufferPool) FetchPage(id PageID) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	bp.disk.Read(id, frame.data)  // I/O with lock!
}

// CORRECT: Release lock before I/O
func (bp *BufferPool) FetchPage(id PageID) {
	bp.mu.Lock()
	frame := bp.getVictim()
	bp.mu.Unlock()

	bp.disk.Read(id, frame.data)  // I/O without lock
}
```

### Issue 3: Missing Error Checks
```go
// WRONG
page, _ := bp.FetchPage(id)  // Ignores error
page.Write(data)

// CORRECT
page, err := bp.FetchPage(id)
if err != nil {
	return err
}
page.Write(data)
```

## Review Questions
1. Is the implementation thread-safe?
2. Can it deadlock?
3. Are there any memory leaks?
4. What happens if all frames are pinned?
5. Is the LRU policy correctly implemented?
6. Are there performance bottlenecks?

## Expected Findings
- 2-3 race conditions
- 1-2 memory leaks
- 2-3 missing edge cases
- 3-5 performance issues
- 5-10 code quality improvements

## Time Estimate
Initial review: 1-2 hours, Detailed analysis: 1-2 hours
