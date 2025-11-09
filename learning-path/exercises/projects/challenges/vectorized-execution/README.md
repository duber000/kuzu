# Challenge 2: Vectorized Execution Engine

## Overview
Implement a vectorized execution engine with batch-at-a-time processing, SIMD operations, and type-specific code generation.

**Duration:** 25-35 hours
**Difficulty:** Very Hard

## Concepts

### Vectorized Execution
Process data in batches (vectors) instead of row-at-a-time:
- Better CPU cache utilization
- SIMD-friendly
- Reduced interpretation overhead
- Amortized function call costs

### Vector Size
Typical: 1024-2048 rows per batch
- Small enough to fit in L1/L2 cache
- Large enough to amortize overhead

## Architecture

### Vector Representation
```go
type Vector struct {
	data     []interface{}  // or type-specific slice
	nulls    *Bitmap
	size     int
	capacity int
}

type IntVector struct {
	data  []int64
	nulls *Bitmap
	size  int
}
```

### Vectorized Operators

#### Scan
```go
type VectorizedScan struct {
	source    DataSource
	batchSize int
}

func (s *VectorizedScan) Next() *VectorBatch {
	// Return batch of rows
}
```

#### Filter
```go
type VectorizedFilter struct {
	child     Operator
	predicate func(*Vector) *Bitmap  // selection vector
}

func (f *VectorizedFilter) Next() *VectorBatch {
	batch := f.child.Next()
	selection := f.predicate(batch.column(0))
	return batch.ApplySelection(selection)
}
```

#### Project
```go
type VectorizedProject struct {
	child       Operator
	expressions []Expression
}
```

#### HashAggregate
```go
type VectorizedHashAggregate struct {
	child    Operator
	groupBy  []int
	aggFuncs []AggFunc
	hashTable *AggHashTable
}
```

## SIMD Operations

### Example: Vectorized Addition
```go
func AddInt64(a, b, result []int64) {
	// TODO: Use SIMD intrinsics or compiler auto-vectorization
	for i := range a {
		result[i] = a[i] + b[i]
	}
}
```

### Example: Vectorized Filter
```go
func FilterGreaterThan(data []int64, threshold int64, selection *Bitmap) {
	// TODO: SIMD comparison
	for i, v := range data {
		if v > threshold {
			selection.Set(i)
		}
	}
}
```

## Type-Specific Code

### Code Generation
Generate type-specific operators to avoid interface overhead:
```go
// Generic (slow)
func Add(a, b Vector) Vector

// Type-specific (fast)
func AddInt64(a, b []int64) []int64
func AddFloat64(a, b []float64) []float64
```

## API Design

```go
type VectorizedEngine struct {
	operators []VectorOperator
}

type VectorOperator interface {
	Next() *VectorBatch
	Reset()
}

type VectorBatch struct {
	columns []*Vector
	size    int
}

func (e *VectorizedEngine) Execute(plan PhysicalPlan) *VectorBatch
```

## Performance Goals

- Filter: >100M rows/sec
- Scan: >500M rows/sec
- Aggregation: >50M rows/sec
- Join: >10M rows/sec
- 5-10x faster than row-at-a-time
- Memory: <100 bytes per batch overhead

## Implementation Hints

### Selection Vectors
Instead of filtering data, track selected positions:
```go
type SelectionVector struct {
	positions []int  // indices of selected rows
	size      int
}
```

### Batch Processing
```go
func (f *VectorizedFilter) Next() *VectorBatch {
	for {
		batch := f.child.Next()
		if batch == nil {
			return nil
		}

		// Apply filter to entire batch at once
		selection := f.applyPredicate(batch)

		if selection.Count() > 0 {
			return batch.ApplySelection(selection)
		}
	}
}
```

## Test Cases
- Single-column operations
- Multi-column operations
- NULL handling
- Large batches (>1M rows)
- Type conversions
- Complex expressions

## Benchmarks
```go
BenchmarkVectorizedFilter
BenchmarkVectorizedProject
BenchmarkVectorizedAggregate
BenchmarkVectorizedJoin
BenchmarkVsRowOriented
```

## Stretch Goals

### 1. JIT Compilation
Generate native code for expressions

### 2. Adaptive Execution
Switch between vectorized and row-oriented

### 3. GPU Execution
Offload to GPU for large scans

## Time Estimate
Core: 15-20 hours, SIMD: 5-7 hours, Testing: 5-8 hours
