# Project 3.4: Pipelined Executor

## Overview
Implement an iterator-based execution engine using Go 1.23 iterators with pipeline operators and profiling.

**Duration:** 18-22 hours
**Difficulty:** Hard

## Core Operators
- **Scan** - Table/index scan
- **Filter** - Predicate filtering
- **Project** - Column projection
- **HashJoin** - Hash join
- **Sort** - Sorting (pipeline breaker)
- **Aggregate** - GROUP BY aggregation

## Go 1.23 Iterators
```go
type Operator interface {
	Execute() iter.Seq[Row]
	Explain() string
	Profile() Stats
}
```

## Features
- Pipeline execution
- Early termination with LIMIT
- Operator profiling (rows, time)
- Adaptive execution
- Memory tracking

## Performance Goals
- Filter: >10M rows/sec
- Projection: >5M rows/sec
- Join: >1M rows/sec
- Memory: streaming where possible

## Time Estimate
Core: 12-15 hours, Profiling: 3-4 hours, Testing: 3-4 hours
