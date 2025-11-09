# Project 3.2: Query Optimizer

## Overview
Implement a cost-based query optimizer with join order optimization, filter pushdown, and cost estimation.

**Duration:** 20-25 hours
**Difficulty:** Hard

## Core Features
- Join order optimization (dynamic programming)
- Filter pushdown
- Projection pushdown
- Cost estimation with statistics
- Plan comparison and selection

## API
```go
type Optimizer struct {
	stats StatisticsCollector
}

func (o *Optimizer) Optimize(plan LogicalPlan) PhysicalPlan
func (o *Optimizer) EstimateCost(plan PhysicalPlan) Cost
```

## Algorithms
- Dynamic programming for join ordering
- Rule-based transformations
- Cost model (CPU + I/O)
- Cardinality estimation

## Time Estimate
Core: 12-15 hours, Testing: 4-5 hours, Extensions: 4-5 hours
