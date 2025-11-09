# Project 3.3: Hash Join Implementation

## Overview
Build efficient join operators using Go 1.24 Swiss Tables with pre-sizing, sort-merge join, and index nested loop join.

**Duration:** 12-15 hours
**Difficulty:** Medium-Hard

## Join Algorithms
1. **Hash Join** - Build hash table, probe
2. **Sort-Merge Join** - Sort both sides, merge
3. **Index Nested Loop** - Use index for inner table

## Go 1.24 Features
- Swiss Tables with pre-sizing for performance
- Efficient hash functions
- Memory-efficient hash tables

## API
```go
func HashJoin(left, right []Row, leftKey, rightKey func(Row) Key) []Row
func SortMergeJoin(left, right []Row) []Row
func IndexNestedLoopJoin(outer []Row, index Index) []Row
```

## Test Cases
- Duplicate keys
- NULL handling
- Empty inputs
- Large joins (1M x 1M)

## Performance Goals
- Hash join: >1M tuples/sec
- Pre-sizing: 20% speedup vs no pre-sizing
- Memory: <100 bytes per tuple overhead

## Time Estimate
Core: 8-10 hours, Optimization: 2-3 hours, Testing: 2-3 hours
