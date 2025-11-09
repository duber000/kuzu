# Phase 4 Lesson 4.2: Go Prep - MVCC (Multi-Version Concurrency Control)

**Prerequisites:** Lesson 4.1 complete (Locking Protocols)
**Time:** 7-8 hours Go prep + 35-40 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 4.2

## Overview

MVCC allows readers and writers to operate without blocking each other. Before implementing snapshot isolation, master these Go concepts:
- Version chains and timestamp ordering
- **Go 1.24:** `weak.Pointer` for old versions ⭐⭐ **MEMORY SAVER!**
- Visibility determination algorithms
- Garbage collection of old versions
- Snapshot isolation implementation

**This lesson is about maximum concurrency without locks!**

## Go Concepts for This Lesson

### 1. Version Chains Basics

**Each row has multiple versions!**

```go
package main

import (
    "fmt"
    "sync"
)

type Timestamp uint64

type Version struct {
    Data      map[string]interface{}
    TxnID     int
    BeginTS   Timestamp  // When this version was created
    EndTS     Timestamp  // When this version was deleted (0 = current)
    Next      *Version   // Older version
}

type VersionedRow struct {
    RowID   int
    Current *Version
    mu      sync.RWMutex
}

func NewVersionedRow(rowID int, data map[string]interface{}, txnID int, ts Timestamp) *VersionedRow {
    return &VersionedRow{
        RowID: rowID,
        Current: &Version{
            Data:    data,
            TxnID:   txnID,
            BeginTS: ts,
            EndTS:   0,  // Current version
            Next:    nil,
        },
    }
}

func (vr *VersionedRow) Read(readTS Timestamp) *Version {
    vr.mu.RLock()
    defer vr.mu.RUnlock()

    // Find version visible at readTS
    for v := vr.Current; v != nil; v = v.Next {
        if v.BeginTS <= readTS && (v.EndTS == 0 || v.EndTS > readTS) {
            return v
        }
    }

    return nil  // No visible version
}

func (vr *VersionedRow) Update(newData map[string]interface{}, txnID int, ts Timestamp) {
    vr.mu.Lock()
    defer vr.mu.Unlock()

    // Create new version
    newVersion := &Version{
        Data:    newData,
        TxnID:   txnID,
        BeginTS: ts,
        EndTS:   0,
        Next:    vr.Current,
    }

    // Mark old version as deleted
    vr.Current.EndTS = ts

    // New version becomes current
    vr.Current = newVersion
}

func main() {
    // Create initial version
    row := NewVersionedRow(1, map[string]interface{}{
        "name": "Alice",
        "age":  30,
    }, 1, 100)

    fmt.Println("Initial version:")
    if v := row.Read(100); v != nil {
        fmt.Printf("  At TS 100: %v\n", v.Data)
    }

    // Update creates new version
    row.Update(map[string]interface{}{
        "name": "Alice",
        "age":  31,
    }, 2, 200)

    // Read at different timestamps
    fmt.Println("\nAfter update:")
    if v := row.Read(150); v != nil {
        fmt.Printf("  At TS 150: age=%d (old version)\n", v.Data["age"])
    }
    if v := row.Read(250); v != nil {
        fmt.Printf("  At TS 250: age=%d (new version)\n", v.Data["age"])
    }
}
```

**Output:**
```
Initial version:
  At TS 100: map[age:30 name:Alice]

After update:
  At TS 150: age=30 (old version)
  At TS 250: age=31 (new version)
```

**Key insight:** Different transactions see different versions based on timestamp!

### 2. Go 1.24: weak.Pointer for Version Chains

**Old versions can be garbage collected!**

```go
package main

import (
    "fmt"
    "runtime"
    "weak"
)

type Version struct {
    Data    map[string]interface{}
    BeginTS Timestamp
    EndTS   Timestamp
    Next    weak.Pointer[*Version]  // Weak reference to old version!
}

type VersionChain struct {
    Current *Version
    Old     []weak.Pointer[*Version]  // Weak references to old versions
}

func (vc *VersionChain) AddVersion(newVersion *Version) {
    if vc.Current != nil {
        // Keep weak reference to old current
        vc.Old = append(vc.Old, weak.Make(&vc.Current))
    }
    vc.Current = newVersion
}

func (vc *VersionChain) GetVersion(ts Timestamp) *Version {
    // Check current version first
    if vc.Current.BeginTS <= ts && (vc.Current.EndTS == 0 || vc.Current.EndTS > ts) {
        return vc.Current
    }

    // Search old versions (might be GC'd!)
    for i := len(vc.Old) - 1; i >= 0; i-- {
        if v := vc.Old[i].Value(); v != nil {
            version := *v
            if version.BeginTS <= ts && (version.EndTS == 0 || version.EndTS > ts) {
                return version
            }
        }
    }

    return nil  // Version was GC'd or doesn't exist
}

func (vc *VersionChain) GarbageCollect(minActiveTS Timestamp) {
    // Remove weak pointers to versions older than min active transaction
    var kept []weak.Pointer[*Version]

    for _, weakPtr := range vc.Old {
        if v := weakPtr.Value(); v != nil {
            version := *v
            if version.EndTS >= minActiveTS {
                kept = append(kept, weakPtr)
            }
            // Else: version too old, let GC reclaim it
        }
    }

    vc.Old = kept
}

func main() {
    vc := &VersionChain{}

    // Add versions
    vc.AddVersion(&Version{
        Data:    map[string]interface{}{"val": 1},
        BeginTS: 100,
        EndTS:   200,
    })

    vc.AddVersion(&Version{
        Data:    map[string]interface{}{"val": 2},
        BeginTS: 200,
        EndTS:   0,
    })

    fmt.Println("Before GC:")
    if v := vc.GetVersion(150); v != nil {
        fmt.Printf("  At TS 150: val=%d\n", v.Data["val"])
    }

    // Garbage collect old versions
    vc.GarbageCollect(200)
    runtime.GC()  // Force garbage collection

    fmt.Println("\nAfter GC:")
    if v := vc.GetVersion(150); v != nil {
        fmt.Printf("  At TS 150: val=%d\n", v.Data["val"])
    } else {
        fmt.Println("  Version GC'd")
    }

    // Current version still accessible
    if v := vc.GetVersion(250); v != nil {
        fmt.Printf("  At TS 250: val=%d\n", v.Data["val"])
    }
}
```

**Key insight:** Weak pointers let GC reclaim old versions no one is reading!

### 3. Timestamp Ordering

**Assign timestamps to maintain serializability!**

```go
package main

import (
    "fmt"
    "sync/atomic"
)

type TimestampManager struct {
    counter atomic.Uint64
}

func NewTimestampManager() *TimestampManager {
    return &TimestampManager{}
}

func (tm *TimestampManager) Next() Timestamp {
    return Timestamp(tm.counter.Add(1))
}

type Transaction struct {
    ID      int
    StartTS Timestamp  // Read timestamp
    CommitTS Timestamp  // Write timestamp (assigned at commit)
}

type MVCCDatabase struct {
    rows map[int]*VersionedRow
    tm   *TimestampManager
}

func NewMVCCDatabase() *MVCCDatabase {
    return &MVCCDatabase{
        rows: make(map[int]*VersionedRow),
        tm:   NewTimestampManager(),
    }
}

func (db *MVCCDatabase) BeginTransaction(txnID int) *Transaction {
    return &Transaction{
        ID:      txnID,
        StartTS: db.tm.Next(),
    }
}

func (db *MVCCDatabase) Read(tx *Transaction, rowID int) *Version {
    row := db.rows[rowID]
    if row == nil {
        return nil
    }

    // Read version visible at transaction's start timestamp
    return row.Read(tx.StartTS)
}

func (db *MVCCDatabase) Write(tx *Transaction, rowID int, data map[string]interface{}) {
    // In real MVCC, this would buffer writes until commit
    // For now, write directly with transaction's commit timestamp
    if tx.CommitTS == 0 {
        tx.CommitTS = db.tm.Next()
    }

    row := db.rows[rowID]
    if row == nil {
        row = NewVersionedRow(rowID, data, tx.ID, tx.CommitTS)
        db.rows[rowID] = row
    } else {
        row.Update(data, tx.ID, tx.CommitTS)
    }
}

func main() {
    db := NewMVCCDatabase()

    // Transaction 1: Initial write
    tx1 := db.BeginTransaction(1)
    db.Write(tx1, 100, map[string]interface{}{"balance": 1000})
    fmt.Printf("Tx1 (TS %d): Wrote balance=1000\n", tx1.CommitTS)

    // Transaction 2: Read (sees Tx1's write)
    tx2 := db.BeginTransaction(2)
    if v := db.Read(tx2, 100); v != nil {
        fmt.Printf("Tx2 (TS %d): Read balance=%d\n", tx2.StartTS, v.Data["balance"])
    }

    // Transaction 3: Update
    tx3 := db.BeginTransaction(3)
    db.Write(tx3, 100, map[string]interface{}{"balance": 1500})
    fmt.Printf("Tx3 (TS %d): Wrote balance=1500\n", tx3.CommitTS)

    // Transaction 2 still sees old value (snapshot isolation!)
    if v := db.Read(tx2, 100); v != nil {
        fmt.Printf("Tx2 (TS %d): Still reads balance=%d (snapshot isolation)\n",
            tx2.StartTS, v.Data["balance"])
    }

    // New transaction sees new value
    tx4 := db.BeginTransaction(4)
    if v := db.Read(tx4, 100); v != nil {
        fmt.Printf("Tx4 (TS %d): Reads balance=%d (new version)\n",
            tx4.StartTS, v.Data["balance"])
    }
}
```

**Output:**
```
Tx1 (TS 1): Wrote balance=1000
Tx2 (TS 2): Read balance=1000
Tx3 (TS 3): Wrote balance=1500
Tx2 (TS 2): Still reads balance=1000 (snapshot isolation)
Tx4 (TS 4): Reads balance=1500 (new version)
```

### 4. Snapshot Isolation

**Transactions see a consistent snapshot!**

```go
package main

import (
    "fmt"
)

type Snapshot struct {
    SnapshotTS Timestamp
    ActiveTxns map[int]bool  // Transactions active at snapshot time
}

func (db *MVCCDatabase) BeginSnapshotTx(txnID int) (*Transaction, *Snapshot) {
    tx := db.BeginTransaction(txnID)

    snapshot := &Snapshot{
        SnapshotTS: tx.StartTS,
        ActiveTxns: make(map[int]bool),
        // In real impl, would record all active transactions
    }

    return tx, snapshot
}

func (s *Snapshot) IsVisible(version *Version) bool {
    // Version is visible if:
    // 1. It was created before snapshot
    // 2. It's not deleted before snapshot
    // 3. Creating transaction committed before snapshot

    if version.BeginTS > s.SnapshotTS {
        return false  // Created after snapshot
    }

    if version.EndTS != 0 && version.EndTS <= s.SnapshotTS {
        return false  // Deleted before snapshot
    }

    if s.ActiveTxns[version.TxnID] {
        return false  // Creating transaction was active (uncommitted)
    }

    return true
}

func (db *MVCCDatabase) ReadSnapshot(snapshot *Snapshot, rowID int) *Version {
    row := db.rows[rowID]
    if row == nil {
        return nil
    }

    row.mu.RLock()
    defer row.mu.RUnlock()

    // Find visible version
    for v := row.Current; v != nil; v = v.Next {
        if snapshot.IsVisible(v) {
            return v
        }
    }

    return nil
}

func main() {
    db := NewMVCCDatabase()

    // Initial state
    tx1 := db.BeginTransaction(1)
    db.Write(tx1, 100, map[string]interface{}{"name": "Alice", "balance": 1000})

    // Start snapshot transaction
    tx2, snapshot := db.BeginSnapshotTx(2)
    fmt.Printf("Tx2 started snapshot at TS %d\n", snapshot.SnapshotTS)

    // Concurrent update
    tx3 := db.BeginTransaction(3)
    db.Write(tx3, 100, map[string]interface{}{"name": "Alice", "balance": 500})
    fmt.Println("Tx3 updated balance to 500")

    // Tx2 still sees snapshot
    if v := db.ReadSnapshot(snapshot, 100); v != nil {
        fmt.Printf("Tx2 (snapshot): balance=%d (snapshot isolation!)\n", v.Data["balance"])
    }

    // New transaction sees latest
    tx4 := db.BeginTransaction(4)
    if v := db.Read(tx4, 100); v != nil {
        fmt.Printf("Tx4: balance=%d (latest)\n", v.Data["balance"])
    }
}
```

### 5. Write-Write Conflict Detection

**Prevent lost updates!**

```go
package main

import (
    "fmt"
)

var ErrWriteConflict = fmt.Errorf("write-write conflict")

func (db *MVCCDatabase) WriteWithConflictDetection(tx *Transaction, rowID int, data map[string]interface{}) error {
    row := db.rows[rowID]
    if row == nil {
        // No conflict for new row
        db.Write(tx, rowID, data)
        return nil
    }

    row.mu.Lock()
    defer row.mu.Unlock()

    // Check if any version was created after our snapshot
    for v := row.Current; v != nil; v = v.Next {
        if v.BeginTS > tx.StartTS {
            // Someone wrote after we started - conflict!
            return ErrWriteConflict
        }

        if v.BeginTS <= tx.StartTS {
            // This is the version we should have read
            break
        }
    }

    // No conflict, safe to write
    if tx.CommitTS == 0 {
        tx.CommitTS = db.tm.Next()
    }

    row.Update(data, tx.ID, tx.CommitTS)
    return nil
}

func main() {
    db := NewMVCCDatabase()

    // Initial state
    tx1 := db.BeginTransaction(1)
    db.Write(tx1, 100, map[string]interface{}{"balance": 1000})
    fmt.Println("Initial balance: 1000")

    // Two concurrent transactions
    tx2 := db.BeginTransaction(2)
    tx3 := db.BeginTransaction(3)

    // Tx3 writes first
    db.WriteWithConflictDetection(tx3, 100, map[string]interface{}{"balance": 900})
    fmt.Println("Tx3: Updated balance to 900")

    // Tx2 tries to write - conflict!
    if err := db.WriteWithConflictDetection(tx2, 100, map[string]interface{}{"balance": 800}); err != nil {
        fmt.Printf("Tx2: Write failed: %v\n", err)
        fmt.Println("Tx2: Must abort and retry")
    }
}
```

**Output:**
```
Initial balance: 1000
Tx3: Updated balance to 900
Tx2: Write failed: write-write conflict
Tx2: Must abort and retry
```

### 6. Garbage Collection

**Reclaim old versions!**

```go
package main

import (
    "fmt"
)

type GarbageCollector struct {
    db             *MVCCDatabase
    minActiveTS    Timestamp
    cleanupCounter int
}

func NewGarbageCollector(db *MVCCDatabase) *GarbageCollector {
    return &GarbageCollector{
        db: db,
    }
}

func (gc *GarbageCollector) UpdateMinActiveTS(ts Timestamp) {
    gc.minActiveTS = ts
}

func (gc *GarbageCollector) CollectRow(row *VersionedRow) int {
    row.mu.Lock()
    defer row.mu.Unlock()

    collected := 0

    // Keep current version and versions visible to active transactions
    var keep *Version
    for v := row.Current; v != nil; {
        next := v.Next

        if v.EndTS != 0 && v.EndTS < gc.minActiveTS {
            // This version is too old, can be GC'd
            collected++
        } else {
            // Keep this version
            v.Next = keep
            keep = v
        }

        v = next
    }

    row.Current = keep
    return collected
}

func (gc *GarbageCollector) Collect() int {
    total := 0

    for _, row := range gc.db.rows {
        total += gc.CollectRow(row)
    }

    gc.cleanupCounter++
    return total
}

func main() {
    db := NewMVCCDatabase()
    gc := NewGarbageCollector(db)

    // Create many versions
    for i := 0; i < 10; i++ {
        tx := db.BeginTransaction(i)
        db.Write(tx, 100, map[string]interface{}{"version": i})
    }

    fmt.Println("Created 10 versions")

    // Simulate min active transaction at TS 8
    gc.UpdateMinActiveTS(8)

    // Run GC
    collected := gc.Collect()
    fmt.Printf("Garbage collected %d old versions\n", collected)
    fmt.Println("Kept versions visible to active transactions (TS >= 8)")
}
```

### 7. Phantom Reads Prevention

**Use predicate locks or serializable snapshot isolation!**

```go
package main

import (
    "fmt"
    "sync"
)

type PredicateLock struct {
    Predicate func(map[string]interface{}) bool
    TxnID     int
    Mode      LockMode
}

type MVCCWithPredicateLocks struct {
    MVCCDatabase
    predicateLocks []PredicateLock
    mu             sync.Mutex
}

func (db *MVCCWithPredicateLocks) LockPredicate(txnID int, predicate func(map[string]interface{}) bool, mode LockMode) {
    db.mu.Lock()
    defer db.mu.Unlock()

    db.predicateLocks = append(db.predicateLocks, PredicateLock{
        Predicate: predicate,
        TxnID:     txnID,
        Mode:      mode,
    })
}

func (db *MVCCWithPredicateLocks) CheckPredicateConflict(txnID int, data map[string]interface{}) bool {
    db.mu.Lock()
    defer db.mu.Unlock()

    for _, plock := range db.predicateLocks {
        if plock.TxnID != txnID && plock.Predicate(data) {
            return true  // Conflict!
        }
    }

    return false
}

func main() {
    db := &MVCCWithPredicateLocks{
        MVCCDatabase: *NewMVCCDatabase(),
    }

    // Tx1: SELECT * FROM accounts WHERE balance > 1000
    tx1 := db.BeginTransaction(1)
    db.LockPredicate(tx1.ID, func(row map[string]interface{}) bool {
        if balance, ok := row["balance"].(int); ok {
            return balance > 1000
        }
        return false
    }, LOCK_SHARED)
    fmt.Println("Tx1: Locked predicate (balance > 1000)")

    // Tx2: INSERT account with balance=1500 - conflicts!
    tx2 := db.BeginTransaction(2)
    newRow := map[string]interface{}{"balance": 1500}

    if db.CheckPredicateConflict(tx2.ID, newRow) {
        fmt.Println("Tx2: Insert conflicts with Tx1's predicate lock")
        fmt.Println("Prevented phantom read!")
    }
}
```

## Pre-Implementation Exercises

### Exercise 1: Version Chain

```go
package main

// TODO: Implement version chain with reads at different timestamps

type Version struct {
    // TODO: Add fields
}

type VersionedRow struct {
    // TODO: Add fields
}

func (vr *VersionedRow) Read(ts Timestamp) *Version {
    // TODO: Find visible version
    return nil
}

func (vr *VersionedRow) Update(data map[string]interface{}, txnID int, ts Timestamp) {
    // TODO: Create new version
}

func main() {
    // TODO: Test multi-version reads
}
```

### Exercise 2: Snapshot Isolation

```go
package main

// TODO: Implement snapshot isolation

type Snapshot struct {
    SnapshotTS Timestamp
    ActiveTxns map[int]bool
}

func (s *Snapshot) IsVisible(version *Version) bool {
    // TODO: Check visibility rules
    return false
}

func main() {
    // TODO: Test snapshot isolation
}
```

### Exercise 3: Garbage Collection

```go
package main

// TODO: Implement version GC

type GarbageCollector struct {
    minActiveTS Timestamp
}

func (gc *GarbageCollector) CollectOldVersions(row *VersionedRow) int {
    // TODO: Remove versions older than minActiveTS
    return 0
}

func main() {
    // TODO: Test GC
}
```

### Exercise 4: Write Conflict Detection

```go
package main

// TODO: Detect write-write conflicts

var ErrWriteConflict = fmt.Errorf("write conflict")

func (db *MVCCDatabase) WriteWithCheck(tx *Transaction, rowID int, data map[string]interface{}) error {
    // TODO: Check for conflicts
    // TODO: Abort if conflict detected
    return nil
}

func main() {
    // TODO: Test conflict detection
}
```

### Exercise 5: Use weak.Pointer

```go
package main

import "weak"

// TODO: Use weak pointers for old versions

type VersionChain struct {
    Current *Version
    Old     []weak.Pointer[*Version]
}

func (vc *VersionChain) GetVersion(ts Timestamp) *Version {
    // TODO: Check weak pointers (might be GC'd)
    return nil
}

func main() {
    // TODO: Test weak pointer GC
}
```

## Performance Benchmarks

### Benchmark 1: MVCC vs Locking

```go
func BenchmarkConcurrentReads(b *testing.B) {
    b.Run("MVCC", func(b *testing.B) {
        // TODO: Multiple readers, no blocking
    })

    b.Run("Locking", func(b *testing.B) {
        // TODO: Readers block writers
    })
}
```

**Expected: MVCC much faster for read-heavy workloads!**

### Benchmark 2: Garbage Collection

```go
func BenchmarkGC(b *testing.B) {
    // TODO: Benchmark GC with different version chain lengths
}
```

## Common Gotchas to Avoid

### Gotcha 1: Incorrect Visibility

```go
// WRONG: Don't check EndTS
func (vr *VersionedRow) Read(ts Timestamp) *Version {
    for v := vr.Current; v != nil; v = v.Next {
        if v.BeginTS <= ts {
            return v  // BUG: Might return deleted version!
        }
    }
    return nil
}

// RIGHT: Check both BeginTS and EndTS
func (vr *VersionedRow) Read(ts Timestamp) *Version {
    for v := vr.Current; v != nil; v = v.Next {
        if v.BeginTS <= ts && (v.EndTS == 0 || v.EndTS > ts) {
            return v  // Correct!
        }
    }
    return nil
}
```

### Gotcha 2: Not Handling Write Conflicts

```go
// WRONG: Allow lost updates
func (db *MVCCDatabase) Write(tx *Transaction, rowID int, data map[string]interface{}) {
    // Just write without checking - LOST UPDATE BUG!
    row.Update(data, tx.ID, tx.CommitTS)
}

// RIGHT: Detect conflicts
func (db *MVCCDatabase) Write(tx *Transaction, rowID int, data map[string]interface{}) error {
    if hasConflict(row, tx.StartTS) {
        return ErrWriteConflict
    }
    row.Update(data, tx.ID, tx.CommitTS)
    return nil
}
```

### Gotcha 3: Keeping Too Many Versions

```go
// WRONG: Never GC
// Result: Memory leak!

// RIGHT: Periodic GC
go func() {
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        gc.Collect()
    }
}()
```

## Checklist Before Starting Lesson 4.2

- [ ] I understand version chains
- [ ] I can implement visibility determination
- [ ] I know how to use `weak.Pointer` for old versions
- [ ] I understand snapshot isolation
- [ ] I can detect write-write conflicts
- [ ] I understand garbage collection of versions
- [ ] I know how to assign timestamps
- [ ] I can prevent phantom reads
- [ ] I understand MVCC performance characteristics
- [ ] I've compared MVCC to locking

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 4.2

You'll implement:
- Complete MVCC with version chains
- Snapshot isolation
- Write conflict detection
- Garbage collection with weak pointers
- Timestamp management
- Serializable snapshot isolation
- Performance benchmarks vs locking

**Time estimate:** 35-40 hours for full implementation

**MVCC enables maximum concurrency!** 🚀

---

## Congratulations!

You've completed the Go Prep for all four phases! You now have:

✅ **Phase 1:** Storage fundamentals with pages, buffer pools, and WAL
✅ **Phase 2:** Graph structures with CSR, columnar storage, and modern Go features
✅ **Phase 3:** Query engine with parsing, planning, joins, and execution
✅ **Phase 4:** Transactions with locking and MVCC

**You're ready to build a production-grade graph database in Go!**

Next: Complete the practice exercises and start implementing the main curriculum! 🎉
