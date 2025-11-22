# Code Review Exercise 3: MVCC Review

## Overview
Review an MVCC implementation for visibility bugs, write conflicts, GC issues, and concurrency problems.

**Duration:** 3-4 hours
**Difficulty:** Very Hard

## Critical Issues to Find

### 1. Visibility Bugs
- Wrong visibility logic
- Snapshot isolation violations
- Dirty reads

### 2. Write Conflicts
- Missing conflict detection
- Lost updates
- Race conditions in commit

### 3. GC Problems
- Memory leaks in version chains
- Premature version collection
- Missing weak pointers

### 4. Concurrency Issues
- Race conditions in version list
- Deadlocks in commit protocol
- Missing locks

## Expected Findings
- 4-6 correctness bugs
- 2-3 race conditions
- 1-2 memory leaks
- 3-5 performance issues

## Time Estimate
3-4 hours
