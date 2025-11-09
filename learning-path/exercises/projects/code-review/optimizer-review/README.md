# Code Review Exercise 2: Query Optimizer Review

## Overview
Review a query optimizer implementation for correctness, performance, and code quality issues.

**Duration:** 2-3 hours
**Difficulty:** Hard

## Issues to Find

### 1. Incorrect Cost Estimates
- Off-by-one errors in cardinality
- Missing selectivity factors
- Wrong cost formulas

### 2. Missing Optimizations
- Cross products not detected
- Filter pushdown not applied
- Projection pruning missing

### 3. Corner Cases
- Empty tables
- Null statistics
- Single-row tables
- Cross products

### 4. Performance Issues
- Exponential join enumeration
- Redundant cost calculations
- Missing memoization

## Expected Findings
- 3-5 correctness bugs
- 2-3 missing optimizations
- 2-4 performance issues
- 5-8 code quality issues

## Time Estimate
2-3 hours
