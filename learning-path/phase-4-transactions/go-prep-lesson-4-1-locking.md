# Phase 4 Lesson 4.1: Go Prep - Locking Protocols

**Prerequisites:** Phase 3 complete (Query Engine)
**Time:** 6-7 hours Go prep + 30-35 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 4.1

## Overview

Locking protocols ensure transaction isolation and prevent data races. Before implementing two-phase locking and deadlock detection, master these Go concepts:
- Lock manager with `sync.Mutex` and `sync.RWMutex`
- Deadlock detection with wait-for graphs
- **Go 1.25:** `testing/synctest` for deterministic concurrency testing ⭐⭐⭐ **GAME CHANGER!**
- Timeout-based deadlock prevention
- Two-phase locking (2PL) protocol

**This lesson is about making transactions safe and correct!**

## Go Concepts for This Lesson

### 1. Lock Manager Basics

**Central lock manager coordinates all locks!**

```go
package main

import (
    "fmt"
    "sync"
)

type LockMode int

const (
    LOCK_SHARED LockMode = iota
    LOCK_EXCLUSIVE
)

type LockRequest struct {
    TxnID int
    Mode  LockMode
}

type LockManager struct {
    mu     sync.Mutex
    locks  map[int]*ResourceLock  // ResourceID -> Lock
}

type ResourceLock struct {
    holders []int      // TxnIDs holding lock
    mode    LockMode   // Current lock mode
    waiters []LockRequest
}

func NewLockManager() *LockManager {
    return &LockManager{
        locks: make(map[int]*ResourceLock),
    }
}

func (lm *LockManager) AcquireLock(resourceID, txnID int, mode LockMode) bool {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    lock, exists := lm.locks[resourceID]
    if !exists {
        lock = &ResourceLock{
            holders: []int{txnID},
            mode:    mode,
            waiters: nil,
        }
        lm.locks[resourceID] = lock
        return true
    }

    // Check compatibility
    if mode == LOCK_SHARED && lock.mode == LOCK_SHARED {
        // Multiple shared locks OK
        lock.holders = append(lock.holders, txnID)
        return true
    }

    // If anyone else holds lock, must wait
    if len(lock.holders) > 0 {
        // In real impl, would block here
        lock.waiters = append(lock.waiters, LockRequest{
            TxnID: txnID,
            Mode:  mode,
        })
        return false
    }

    // Grant lock
    lock.holders = append(lock.holders, txnID)
    lock.mode = mode
    return true
}

func (lm *LockManager) ReleaseLock(resourceID, txnID int) {
    lm.mu.Lock()
    defer lm.mu.Unlock()

    lock, exists := lm.locks[resourceID]
    if !exists {
        return
    }

    // Remove from holders
    for i, holder := range lock.holders {
        if holder == txnID {
            lock.holders = append(lock.holders[:i], lock.holders[i+1:]...)
            break
        }
    }

    // If no more holders, grant to waiters
    if len(lock.holders) == 0 && len(lock.waiters) > 0 {
        // Grant to first waiter (simplified)
        waiter := lock.waiters[0]
        lock.holders = append(lock.holders, waiter.TxnID)
        lock.mode = waiter.Mode
        lock.waiters = lock.waiters[1:]
    }
}

func main() {
    lm := NewLockManager()

    // Transaction 1 acquires shared lock
    if lm.AcquireLock(100, 1, LOCK_SHARED) {
        fmt.Println("Txn 1 acquired shared lock on resource 100")
    }

    // Transaction 2 acquires shared lock (compatible!)
    if lm.AcquireLock(100, 2, LOCK_SHARED) {
        fmt.Println("Txn 2 acquired shared lock on resource 100")
    }

    // Transaction 3 wants exclusive lock (must wait)
    if !lm.AcquireLock(100, 3, LOCK_EXCLUSIVE) {
        fmt.Println("Txn 3 must wait for exclusive lock on resource 100")
    }

    // Transaction 1 releases
    lm.ReleaseLock(100, 1)
    fmt.Println("Txn 1 released lock")

    // Transaction 2 releases
    lm.ReleaseLock(100, 2)
    fmt.Println("Txn 2 released lock")
    fmt.Println("Now Txn 3 can proceed")
}
```

**Output:**
```
Txn 1 acquired shared lock on resource 100
Txn 2 acquired shared lock on resource 100
Txn 3 must wait for exclusive lock on resource 100
Txn 1 released lock
Txn 2 released lock
Now Txn 3 can proceed
```

### 2. Wait-For Graph for Deadlock Detection

**Detect cycles to find deadlocks!**

```go
package main

import (
    "fmt"
)

type WaitForGraph struct {
    edges map[int][]int  // TxnID -> list of TxnIDs it's waiting for
}

func NewWaitForGraph() *WaitForGraph {
    return &WaitForGraph{
        edges: make(map[int][]int),
    }
}

func (wfg *WaitForGraph) AddEdge(waiter, holder int) {
    wfg.edges[waiter] = append(wfg.edges[waiter], holder)
}

func (wfg *WaitForGraph) RemoveEdge(waiter, holder int) {
    waitList := wfg.edges[waiter]
    for i, h := range waitList {
        if h == holder {
            wfg.edges[waiter] = append(waitList[:i], waitList[i+1:]...)
            break
        }
    }
}

func (wfg *WaitForGraph) RemoveTransaction(txnID int) {
    delete(wfg.edges, txnID)
    // Also remove from other wait lists
    for waiter := range wfg.edges {
        wfg.RemoveEdge(waiter, txnID)
    }
}

// Detect cycles using DFS
func (wfg *WaitForGraph) DetectDeadlock() []int {
    visited := make(map[int]bool)
    recStack := make(map[int]bool)
    var path []int

    var dfs func(int) bool
    dfs = func(node int) bool {
        visited[node] = true
        recStack[node] = true
        path = append(path, node)

        for _, neighbor := range wfg.edges[node] {
            if !visited[neighbor] {
                if dfs(neighbor) {
                    return true
                }
            } else if recStack[neighbor] {
                // Found cycle!
                return true
            }
        }

        recStack[node] = false
        path = path[:len(path)-1]
        return false
    }

    for txn := range wfg.edges {
        if !visited[txn] {
            if dfs(txn) {
                return path
            }
        }
    }

    return nil
}

func main() {
    wfg := NewWaitForGraph()

    // Txn 1 waits for Txn 2
    wfg.AddEdge(1, 2)

    // Txn 2 waits for Txn 3
    wfg.AddEdge(2, 3)

    // Txn 3 waits for Txn 1 - DEADLOCK!
    wfg.AddEdge(3, 1)

    if cycle := wfg.DetectDeadlock(); cycle != nil {
        fmt.Println("Deadlock detected!")
        fmt.Printf("Cycle: %v\n", cycle)
    } else {
        fmt.Println("No deadlock")
    }
}
```

**Output:**
```
Deadlock detected!
Cycle: [1 2 3]
```

### 3. Go 1.25: synctest for Deterministic Testing

**Test concurrency bugs deterministically!**

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "testing"
    "testing/synctest"
    "time"
)

// Simulate a deadlock scenario
func TestDeadlockDetection(t *testing.T) {
    synctest.Run(func() {
        lm := NewLockManager()
        deadlockDetected := atomic.Bool{}

        // Transaction 1: Lock A then B
        go func() {
            lm.AcquireLock(1, 100, LOCK_EXCLUSIVE)  // Lock A
            time.Sleep(10 * time.Millisecond)       // Simulated work
            if !lm.AcquireLock(2, 100, LOCK_EXCLUSIVE) {  // Try Lock B
                deadlockDetected.Store(true)
            }
        }()

        // Transaction 2: Lock B then A
        go func() {
            lm.AcquireLock(2, 200, LOCK_EXCLUSIVE)  // Lock B
            time.Sleep(10 * time.Millisecond)       // Simulated work
            if !lm.AcquireLock(1, 200, LOCK_EXCLUSIVE) {  // Try Lock A
                deadlockDetected.Store(true)
            }
        }()

        synctest.Wait()  // Wait for all goroutines deterministically!

        if !deadlockDetected.Load() {
            t.Error("Expected deadlock to be detected")
        }
    })
}

// Test two-phase locking
func TestTwoPhaseLocking(t *testing.T) {
    synctest.Run(func() {
        lm := NewLockManager()
        var counter int
        var mu sync.Mutex

        // Multiple transactions incrementing counter
        for i := 0; i < 10; i++ {
            txnID := i
            go func() {
                // Growing phase: acquire locks
                lm.AcquireLock(1, txnID, LOCK_EXCLUSIVE)

                // Critical section
                mu.Lock()
                counter++
                mu.Unlock()

                time.Sleep(1 * time.Millisecond)

                // Shrinking phase: release locks
                lm.ReleaseLock(1, txnID)
            }()
        }

        synctest.Wait()

        if counter != 10 {
            t.Errorf("Expected counter=10, got %d", counter)
        }
    })
}

func main() {
    fmt.Println("Run with: go test -v")
}
```

**Key insight:** `synctest` makes flaky concurrency tests 100% reproducible!

### 4. Lock Upgrade/Downgrade

**Convert lock modes safely!**

```go
package main

import (
    "fmt"
    "sync"
)

type Lock struct {
    mu      sync.Mutex
    mode    LockMode
    holders map[int]LockMode  // TxnID -> lock mode
}

func NewLock() *Lock {
    return &Lock{
        holders: make(map[int]LockMode),
    }
}

func (l *Lock) Acquire(txnID int, mode LockMode) bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    // Already holding this lock?
    if existingMode, held := l.holders[txnID]; held {
        if mode == LOCK_EXCLUSIVE && existingMode == LOCK_SHARED {
            // Upgrade: only if we're the only holder
            if len(l.holders) == 1 {
                l.holders[txnID] = LOCK_EXCLUSIVE
                l.mode = LOCK_EXCLUSIVE
                return true
            }
            return false  // Can't upgrade, others hold shared locks
        }
        return true  // Already have sufficient lock
    }

    // Check compatibility
    if len(l.holders) == 0 {
        // No one holds lock
        l.holders[txnID] = mode
        l.mode = mode
        return true
    }

    if mode == LOCK_SHARED && l.mode == LOCK_SHARED {
        // Multiple shared locks OK
        l.holders[txnID] = mode
        return true
    }

    return false  // Incompatible
}

func (l *Lock) Downgrade(txnID int) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if mode, held := l.holders[txnID]; held && mode == LOCK_EXCLUSIVE {
        l.holders[txnID] = LOCK_SHARED
        l.mode = LOCK_SHARED
    }
}

func (l *Lock) Release(txnID int) {
    l.mu.Lock()
    defer l.mu.Unlock()

    delete(l.holders, txnID)
    if len(l.holders) == 0 {
        l.mode = LOCK_SHARED  // Reset
    }
}

func main() {
    lock := NewLock()

    // Txn 1 gets shared lock
    if lock.Acquire(1, LOCK_SHARED) {
        fmt.Println("Txn 1: Acquired shared lock")
    }

    // Txn 1 tries to upgrade (succeeds - only holder)
    lock.Release(1)
    lock.Acquire(1, LOCK_SHARED)
    if lock.Acquire(1, LOCK_EXCLUSIVE) {
        fmt.Println("Txn 1: Upgraded to exclusive")
    }

    // Txn 1 downgrades
    lock.Downgrade(1)
    fmt.Println("Txn 1: Downgraded to shared")

    // Now Txn 2 can get shared lock
    if lock.Acquire(2, LOCK_SHARED) {
        fmt.Println("Txn 2: Acquired shared lock")
    }
}
```

### 5. Two-Phase Locking (2PL)

**Growing phase then shrinking phase!**

```go
package main

import (
    "fmt"
    "sync"
)

type Transaction struct {
    ID            int
    locksHeld     map[int]LockMode  // ResourceID -> Mode
    phase         string            // "growing" or "shrinking"
    mu            sync.Mutex
}

func NewTransaction(id int) *Transaction {
    return &Transaction{
        ID:        id,
        locksHeld: make(map[int]LockMode),
        phase:     "growing",
    }
}

func (tx *Transaction) Lock(lm *LockManager, resourceID int, mode LockMode) error {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    if tx.phase == "shrinking" {
        return fmt.Errorf("cannot acquire lock in shrinking phase")
    }

    if lm.AcquireLock(resourceID, tx.ID, mode) {
        tx.locksHeld[resourceID] = mode
        return nil
    }

    return fmt.Errorf("failed to acquire lock")
}

func (tx *Transaction) Unlock(lm *LockManager, resourceID int) {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    // First unlock transitions to shrinking phase
    if tx.phase == "growing" {
        tx.phase = "shrinking"
    }

    lm.ReleaseLock(resourceID, tx.ID)
    delete(tx.locksHeld, resourceID)
}

func (tx *Transaction) Commit(lm *LockManager) {
    tx.mu.Lock()
    defer tx.mu.Unlock()

    // Release all locks
    for resourceID := range tx.locksHeld {
        lm.ReleaseLock(resourceID, tx.ID)
    }

    tx.locksHeld = make(map[int]LockMode)
    tx.phase = "committed"
}

func main() {
    lm := NewLockManager()
    tx := NewTransaction(1)

    // Growing phase
    tx.Lock(lm, 100, LOCK_SHARED)
    fmt.Println("Acquired lock on 100 (growing phase)")

    tx.Lock(lm, 200, LOCK_EXCLUSIVE)
    fmt.Println("Acquired lock on 200 (growing phase)")

    // Shrinking phase begins
    tx.Unlock(lm, 100)
    fmt.Println("Released lock on 100 (now in shrinking phase)")

    // This will fail!
    if err := tx.Lock(lm, 300, LOCK_SHARED); err != nil {
        fmt.Printf("Failed to acquire lock on 300: %v\n", err)
    }

    // Commit (releases remaining locks)
    tx.Commit(lm)
    fmt.Println("Transaction committed")
}
```

**Output:**
```
Acquired lock on 100 (growing phase)
Acquired lock on 200 (growing phase)
Released lock on 100 (now in shrinking phase)
Failed to acquire lock on 300: cannot acquire lock in shrinking phase
Transaction committed
```

### 6. Timeout-Based Deadlock Prevention

**Abort transactions that wait too long!**

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func (lm *LockManager) AcquireLockWithTimeout(ctx context.Context, resourceID, txnID int, mode LockMode) error {
    // Try to acquire lock with timeout
    ticker := time.NewTicker(10 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return fmt.Errorf("lock acquisition timeout")

        case <-ticker.C:
            if lm.AcquireLock(resourceID, txnID, mode) {
                return nil
            }
            // Keep trying
        }
    }
}

func main() {
    lm := NewLockManager()

    // Txn 1 holds lock
    lm.AcquireLock(100, 1, LOCK_EXCLUSIVE)
    fmt.Println("Txn 1 acquired exclusive lock")

    // Txn 2 tries with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    if err := lm.AcquireLockWithTimeout(ctx, 100, 2, LOCK_EXCLUSIVE); err != nil {
        fmt.Printf("Txn 2 failed: %v\n", err)
        fmt.Println("Txn 2 aborted to prevent deadlock")
    }
}
```

### 7. Lock Granularity

**Table locks vs row locks vs page locks!**

```go
package main

import (
    "fmt"
    "sync"
)

type LockGranularity int

const (
    LOCK_TABLE LockGranularity = iota
    LOCK_PAGE
    LOCK_ROW
)

type HierarchicalLockManager struct {
    mu         sync.Mutex
    tableLocks map[string]*Lock       // TableName -> Lock
    pageLocks  map[string]*Lock       // TableName:PageID -> Lock
    rowLocks   map[string]*Lock       // TableName:PageID:RowID -> Lock
}

func NewHierarchicalLockManager() *HierarchicalLockManager {
    return &HierarchicalLockManager{
        tableLocks: make(map[string]*Lock),
        pageLocks:  make(map[string]*Lock),
        rowLocks:   make(map[string]*Lock),
    }
}

func (hlm *HierarchicalLockManager) LockRow(table string, pageID, rowID, txnID int, mode LockMode) bool {
    hlm.mu.Lock()
    defer hlm.mu.Unlock()

    // Must acquire intention lock on table and page first
    tableKey := table
    pageKey := fmt.Sprintf("%s:%d", table, pageID)
    rowKey := fmt.Sprintf("%s:%d:%d", table, pageID, rowID)

    // Acquire intention locks (simplified)
    // In real system, would check compatibility

    // Acquire row lock
    rowLock, exists := hlm.rowLocks[rowKey]
    if !exists {
        rowLock = NewLock()
        hlm.rowLocks[rowKey] = rowLock
    }

    return rowLock.Acquire(txnID, mode)
}

func main() {
    hlm := NewHierarchicalLockManager()

    // Txn 1 locks row
    if hlm.LockRow("users", 10, 5, 1, LOCK_SHARED) {
        fmt.Println("Txn 1: Locked users:10:5 (shared)")
    }

    // Txn 2 locks different row in same page
    if hlm.LockRow("users", 10, 6, 2, LOCK_EXCLUSIVE) {
        fmt.Println("Txn 2: Locked users:10:6 (exclusive)")
    }

    fmt.Println("Fine-grained locking allows high concurrency!")
}
```

## Pre-Implementation Exercises

### Exercise 1: Implement Lock Manager

```go
package main

// TODO: Implement lock manager with shared/exclusive locks

type LockManager struct {
    // TODO: Add fields
}

func (lm *LockManager) AcquireLock(resourceID, txnID int, mode LockMode) bool {
    // TODO: Implement
    return false
}

func (lm *LockManager) ReleaseLock(resourceID, txnID int) {
    // TODO: Implement
}

func main() {
    // TODO: Test with multiple transactions
}
```

### Exercise 2: Deadlock Detection

```go
package main

// TODO: Implement wait-for graph and cycle detection

type WaitForGraph struct {
    // TODO: Add fields
}

func (wfg *WaitForGraph) DetectDeadlock() []int {
    // TODO: Use DFS to find cycles
    return nil
}

func main() {
    // TODO: Create deadlock scenario and detect it
}
```

### Exercise 3: Two-Phase Locking

```go
package main

// TODO: Implement 2PL transaction

type Transaction struct {
    phase string  // "growing" or "shrinking"
    // TODO: Add more fields
}

func (tx *Transaction) Lock(resourceID int) error {
    // TODO: Enforce growing phase
    return nil
}

func (tx *Transaction) Unlock(resourceID int) {
    // TODO: Transition to shrinking phase
}

func main() {
    // TODO: Test 2PL enforcement
}
```

### Exercise 4: Test with synctest

```go
package main

import (
    "testing"
    "testing/synctest"
)

func TestConcurrentTransactions(t *testing.T) {
    synctest.Run(func() {
        // TODO: Create deterministic concurrency test
        // TODO: Test for race conditions
        // TODO: Test for deadlocks
    })
}
```

### Exercise 5: Lock Timeout

```go
package main

import (
    "context"
    "time"
)

func (lm *LockManager) AcquireWithTimeout(ctx context.Context, resourceID, txnID int, mode LockMode) error {
    // TODO: Implement timeout-based acquisition
    // TODO: Return error on timeout
    return nil
}

func main() {
    // TODO: Test timeout behavior
}
```

## Performance Benchmarks

### Benchmark 1: Lock Contention

```go
func BenchmarkLockContention(b *testing.B) {
    lm := NewLockManager()

    b.Run("LowContention", func(b *testing.B) {
        // 100 resources, 10 txns
    })

    b.Run("HighContention", func(b *testing.B) {
        // 10 resources, 100 txns
    })
}
```

### Benchmark 2: Deadlock Detection

```go
func BenchmarkDeadlockDetection(b *testing.B) {
    wfg := NewWaitForGraph()
    // Add edges
    for b.Loop() {
        _ = wfg.DetectDeadlock()
    }
}
```

## Common Gotchas to Avoid

### Gotcha 1: Not Checking Lock Compatibility

```go
// WRONG: Always grant lock
func (lm *LockManager) AcquireLock(resourceID, txnID int, mode LockMode) bool {
    lock := lm.locks[resourceID]
    lock.holders = append(lock.holders, txnID)
    return true  // BUG: Doesn't check compatibility!
}

// RIGHT: Check compatibility
func (lm *LockManager) AcquireLock(resourceID, txnID int, mode LockMode) bool {
    lock := lm.locks[resourceID]

    if mode == LOCK_EXCLUSIVE && len(lock.holders) > 0 {
        return false  // Can't get exclusive if anyone holds lock
    }

    if mode == LOCK_SHARED && lock.mode == LOCK_EXCLUSIVE {
        return false  // Can't get shared if someone holds exclusive
    }

    lock.holders = append(lock.holders, txnID)
    return true
}
```

### Gotcha 2: Forgetting to Release Locks

```go
// WRONG: Don't release on error
func doWork(lm *LockManager, txnID int) error {
    lm.AcquireLock(100, txnID, LOCK_EXCLUSIVE)
    if err := process(); err != nil {
        return err  // BUG: Lock not released!
    }
    lm.ReleaseLock(100, txnID)
    return nil
}

// RIGHT: Use defer
func doWork(lm *LockManager, txnID int) error {
    lm.AcquireLock(100, txnID, LOCK_EXCLUSIVE)
    defer lm.ReleaseLock(100, txnID)

    if err := process(); err != nil {
        return err  // Lock will be released
    }
    return nil
}
```

### Gotcha 3: Violating 2PL

```go
// WRONG: Acquire lock after releasing
tx.Lock(100)
tx.Lock(200)
tx.Unlock(100)  // Shrinking phase starts
tx.Lock(300)    // VIOLATION! Can't lock in shrinking phase
```

## Checklist Before Starting Lesson 4.1

- [ ] I understand lock manager design
- [ ] I can implement shared and exclusive locks
- [ ] I know how to detect deadlocks with wait-for graphs
- [ ] I understand two-phase locking protocol
- [ ] I can use `testing/synctest` for deterministic tests
- [ ] I know how to handle lock upgrades/downgrades
- [ ] I understand timeout-based deadlock prevention
- [ ] I can implement hierarchical locking
- [ ] I always release locks (use defer!)
- [ ] I've tested with concurrent workloads

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 4.1

You'll implement:
- Full lock manager with shared/exclusive/intention locks
- Wait-for graph deadlock detection
- Two-phase locking enforcement
- Lock timeout and deadlock prevention
- Hierarchical locking (table/page/row)
- Comprehensive concurrency tests with synctest

**Time estimate:** 30-35 hours for full implementation

**Locking keeps transactions safe and isolated!** 🔒
