# Project 4.2: MVCC Implementation

## Overview
Build a Multi-Version Concurrency Control (MVCC) system with snapshot isolation, version chains, and garbage collection using Go 1.24 `weak.Pointer`.

**Duration:** 25-30 hours
**Difficulty:** Very Hard

## Core Features
- Version chain management
- Snapshot isolation
- Write conflict detection
- Garbage collection with `weak.Pointer`
- Transaction visibility rules

## Key Concepts

### Version Chains
```
Latest -> V3 (txn=103) -> V2 (txn=102) -> V1 (txn=101) -> NULL
Each version has:
- Data
- Begin timestamp
- End timestamp
- Transaction ID
```

### Snapshot Isolation
- Each transaction sees consistent snapshot
- Read from appropriate version
- Write conflicts detected at commit

### Visibility Rules
```go
func isVisible(version *Version, snapshot Timestamp) bool {
	return version.BeginTS <= snapshot &&
	       (version.EndTS == nil || version.EndTS > snapshot)
}
```

## API Design
```go
type MVCCStore struct {
	data     map[Key]*VersionChain
	gcQueue  *GarbageCollector
}

func (s *MVCCStore) Read(key Key, snapshot Timestamp) (Value, error)
func (s *MVCCStore) Write(key Key, value Value, txn TxnID) error
func (s *MVCCStore) Commit(txn TxnID) error
func (s *MVCCStore) GC(olderThan Timestamp)
```

## Go 1.24 weak.Pointer for GC
```go
import "weak"

type Version struct {
	data     Value
	beginTS  Timestamp
	endTS    *Timestamp
	prev     weak.Pointer[*Version]  // GC-friendly
}
```

## Performance Goals
- Read latency: <1µs
- Write latency: <10µs
- GC throughput: >1M versions/sec
- Read-heavy: 10x better than locking
- Write conflicts: <1% in normal workload

## Time Estimate
Core: 15-18 hours, GC: 4-5 hours, Testing: 4-5 hours, Optimization: 2-3 hours
