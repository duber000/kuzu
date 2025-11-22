# Project 2.3: Parallel Graph Algorithms

## Overview
Implement parallel graph algorithms using worker pools, Go 1.25 `testing/synctest` for deterministic testing, and measure scalability.

**Duration:** 15-20 hours
**Difficulty:** Hard

## Learning Objectives
- Implement parallel graph algorithms
- Master worker pool patterns
- Use Go 1.25 `testing/synctest` for deterministic concurrency testing
- Measure parallel speedup and scalability
- Handle load balancing
- Avoid race conditions

## Algorithms to Implement

### 1. Parallel BFS
- Level-synchronous BFS
- Frontier-based approach
- Work stealing for load balancing

### 2. PageRank
- Iterative computation
- Parallel edge traversal
- Convergence detection

### 3. Triangle Counting
- Parallel edge intersection
- Work distribution strategies
- Atomic counters

### 4. Connected Components
- Union-find with path compression
- Parallel edge processing
- Lock-free optimizations

## API Design

```go
package parallelalgo

// Parallel BFS
func ParallelBFS(g *CSRGraph, source NodeID, workers int) []int

// PageRank
func PageRank(g *CSRGraph, iterations int, workers int) []float64

// Triangle counting
func CountTriangles(g *CSRGraph, workers int) int64

// Connected components
func ConnectedComponents(g *CSRGraph, workers int) []int
```

## Key Concepts

### Worker Pool Pattern
```go
type WorkerPool struct {
	workers   int
	workQueue chan Task
	results   chan Result
	wg        sync.WaitGroup
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.workQueue {
		result := task.Execute()
		p.results <- result
	}
}
```

### Testing with synctest (Go 1.25)
```go
func TestParallelBFS_Deterministic(t *testing.T) {
	synctest.Run(func() {
		// Test runs deterministically
		// All goroutine interleavings explored
	})
}
```

## Performance Goals

- BFS speedup: >3x with 4 cores
- PageRank: >5x with 8 cores
- Triangle counting: >4x with 4 cores
- No race conditions (test with -race)
- Load balancing: <10% variance

## Time Estimate
Core: 10-12 hours, Testing: 3-4 hours, Optimization: 3-4 hours
