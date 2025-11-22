# Debug Scenario 1: Race Condition Hunt

## Overview
Fix intentional race conditions in buffer pool, hash join, and lock manager implementations.

**Duration:** 3-5 hours
**Difficulty:** Medium

## Scenarios

### Scenario 1: Buffer Pool Race
**Bug:** Pin count updated without proper synchronization
```go
// BUGGY CODE
func (bp *BufferPool) FetchPage(pageID PageID) (*Frame, error) {
	frame := bp.getFrame(pageID)
	frame.pinCount++  // RACE!
	return frame, nil
}
```

**Symptoms:**
- Incorrect pin counts
- Premature page eviction
- Crashes under load

**Tools:**
- `go test -race`
- `testing/synctest`
- Print debugging

### Scenario 2: Hash Join Race
**Bug:** Concurrent writes to result slice
```go
// BUGGY CODE
func HashJoin(left, right []Row) []Row {
	results := make([]Row, 0)
	var wg sync.WaitGroup

	for _, row := range left {
		wg.Add(1)
		go func(r Row) {
			defer wg.Done()
			if match := probe(r); match != nil {
				results = append(results, match)  // RACE!
			}
		}(row)
	}
	wg.Wait()
	return results
}
```

**Symptoms:**
- Missing join results
- Panic: concurrent map writes
- Non-deterministic output

### Scenario 3: Lock Manager Race
**Bug:** Wait-for graph updated without lock
```go
// BUGGY CODE
func (lm *LockManager) AcquireLock(txn TxnID, resource ResourceID) error {
	if !lm.canGrant(txn, resource) {
		lm.waitForGraph[txn] = append(lm.waitForGraph[txn], holder)  // RACE!
		lm.wait(txn, resource)
	}
}
```

**Symptoms:**
- Deadlock detection failures
- Missing edges in wait-for graph
- Panics

## Debugging Workflow

1. **Reproduce:** Run with `-race` flag
2. **Locate:** Identify racy variables
3. **Fix:** Add proper synchronization
4. **Verify:** Re-run with `-race`
5. **Test:** Add test to prevent regression

## Tools

```bash
# Detect races
go test -race

# Use synctest for determinism
go test -run=TestRace

# Stress test
go test -race -count=100

# Profile goroutines
go test -trace=trace.out
go tool trace trace.out
```

## Learning Objectives
- Identify common race conditions
- Use Go race detector effectively
- Apply proper synchronization patterns
- Write race-free concurrent code

## Time Estimate
Each scenario: 1-2 hours
