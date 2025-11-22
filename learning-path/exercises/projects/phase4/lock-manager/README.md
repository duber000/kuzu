# Project 4.1: Lock Manager

## Overview
Implement a lock manager with shared/exclusive locks, deadlock detection, and two-phase locking using Go 1.25 `testing/synctest`.

**Duration:** 15-18 hours
**Difficulty:** Hard

## Core Features
- Shared and exclusive locks
- Lock upgrades/downgrades
- Deadlock detection with wait-for graph
- Two-phase locking enforcement
- Lock timeout handling

## API Design
```go
type LockManager struct {
	locks map[ResourceID]*LockTable
	waitForGraph *WaitForGraph
}

func (lm *LockManager) AcquireLock(txn TxnID, resource ResourceID, mode LockMode) error
func (lm *LockManager) ReleaseLock(txn TxnID, resource ResourceID) error
func (lm *LockManager) DetectDeadlock() ([]TxnID, bool)
```

## Lock Modes
- **Shared (S)** - Read lock
- **Exclusive (X)** - Write lock
- **IntentionShared (IS)** - Intent to read
- **IntentionExclusive (IX)** - Intent to write

## Deadlock Detection
- Wait-for graph construction
- Cycle detection algorithm
- Victim selection for abort

## Testing with synctest (Go 1.25)
```go
func TestDeadlock_Deterministic(t *testing.T) {
	synctest.Run(func() {
		// Test deadlock scenarios deterministically
	})
}
```

## Performance Goals
- Lock acquisition: <1µs (uncontended)
- Deadlock detection: <10ms
- Support 10K concurrent transactions
- No race conditions

## Time Estimate
Core: 10-12 hours, Deadlock detection: 3-4 hours, Testing: 2-3 hours
