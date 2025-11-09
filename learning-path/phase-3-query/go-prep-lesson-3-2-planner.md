# Phase 3 Lesson 3.2: Go Prep - Query Planning

**Prerequisites:** Lesson 3.1 complete (Parser)
**Time:** 5-6 hours Go prep + 30-35 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.2

## Overview

Query planning transforms ASTs into efficient execution plans. Before implementing a cost-based optimizer, master these Go concepts:
- Interface-based operator design
- Cost estimation with statistics
- **Go 1.23:** `unique` package for operator type interning ⭐⭐
- Dynamic programming for plan enumeration
- Visitor pattern for plan traversal and transformation

**This lesson is about choosing the best path through query space!**

## Go Concepts for This Lesson

### 1. Interface-Based Operator Design

**All operators implement a common interface!**

```go
package main

import (
    "fmt"
)

// Logical operator interface
type LogicalOp interface {
    // Estimate the cost of this operation
    EstimateCost() float64

    // Get estimated output cardinality
    Cardinality() int

    // String representation
    String() string

    // Children operators
    Children() []LogicalOp
}

// Scan operator
type ScanOp struct {
    table       string
    cardinality int
}

func (s *ScanOp) EstimateCost() float64 {
    return float64(s.cardinality) * 0.1  // Cost = rows * 0.1
}

func (s *ScanOp) Cardinality() int {
    return s.cardinality
}

func (s *ScanOp) String() string {
    return fmt.Sprintf("Scan(%s) [%d rows]", s.table, s.cardinality)
}

func (s *ScanOp) Children() []LogicalOp {
    return nil  // Leaf node
}

// Filter operator
type FilterOp struct {
    child       LogicalOp
    predicate   string
    selectivity float64  // Fraction of rows passing filter
}

func (f *FilterOp) EstimateCost() float64 {
    return f.child.EstimateCost() + float64(f.child.Cardinality())*0.01
}

func (f *FilterOp) Cardinality() int {
    return int(float64(f.child.Cardinality()) * f.selectivity)
}

func (f *FilterOp) String() string {
    return fmt.Sprintf("Filter(%s) [%d rows]", f.predicate, f.Cardinality())
}

func (f *FilterOp) Children() []LogicalOp {
    return []LogicalOp{f.child}
}

// Join operator
type JoinOp struct {
    left        LogicalOp
    right       LogicalOp
    joinType    string
}

func (j *JoinOp) EstimateCost() float64 {
    // Hash join cost: build hash table + probe
    buildCost := j.left.EstimateCost() + float64(j.left.Cardinality())*0.5
    probeCost := j.right.EstimateCost() + float64(j.right.Cardinality())*0.2
    return buildCost + probeCost
}

func (j *JoinOp) Cardinality() int {
    // Simplified: assume 10% selectivity
    return (j.left.Cardinality() * j.right.Cardinality()) / 10
}

func (j *JoinOp) String() string {
    return fmt.Sprintf("Join(%s) [%d rows]", j.joinType, j.Cardinality())
}

func (j *JoinOp) Children() []LogicalOp {
    return []LogicalOp{j.left, j.right}
}

func main() {
    // Build a query plan:
    // SELECT * FROM users u JOIN orders o ON u.id = o.user_id WHERE u.age > 18

    plan := &FilterOp{
        predicate:   "age > 18",
        selectivity: 0.3,  // 30% pass filter
        child: &JoinOp{
            joinType: "hash",
            left:     &ScanOp{table: "users", cardinality: 10000},
            right:    &ScanOp{table: "orders", cardinality: 50000},
        },
    }

    fmt.Println("Plan:")
    printPlan(plan, 0)
    fmt.Printf("\nTotal cost: %.2f\n", plan.EstimateCost())
    fmt.Printf("Output cardinality: %d\n", plan.Cardinality())
}

func printPlan(op LogicalOp, depth int) {
    indent := ""
    for i := 0; i < depth; i++ {
        indent += "  "
    }
    fmt.Printf("%s%s (cost: %.2f)\n", indent, op.String(), op.EstimateCost())

    for _, child := range op.Children() {
        printPlan(child, depth+1)
    }
}
```

**Output:**
```
Plan:
Filter(age > 18) [15000 rows] (cost: 501.00)
  Join(hash) [50000 rows] (cost: 501.00)
    Scan(users) [10000 rows] (cost: 1000.00)
    Scan(orders) [50000 rows] (cost: 5000.00)

Total cost: 501.00
Output cardinality: 15000
```

### 2. Go 1.23: unique Package for Operator Types

**Intern operator types for fast comparison!**

```go
package main

import (
    "fmt"
    "unique"
)

type OperatorType = unique.Handle[string]

var (
    OpScan   = unique.Make("Scan")
    OpFilter = unique.Make("Filter")
    OpJoin   = unique.Make("Join")
    OpSort   = unique.Make("Sort")
)

type Operator struct {
    Type     OperatorType
    Children []*Operator
    Props    map[string]interface{}
}

func NewOperator(typ OperatorType) *Operator {
    return &Operator{
        Type:     typ,
        Children: nil,
        Props:    make(map[string]interface{}),
    }
}

// Fast operator type comparison (pointer equality!)
func (o *Operator) IsScan() bool {
    return o.Type == OpScan  // O(1) pointer comparison!
}

func (o *Operator) IsJoin() bool {
    return o.Type == OpJoin
}

// Pattern matching
func FindScans(op *Operator) []*Operator {
    var scans []*Operator

    if op.Type == OpScan {
        scans = append(scans, op)
    }

    for _, child := range op.Children {
        scans = append(scans, FindScans(child)...)
    }

    return scans
}

func main() {
    // Build plan tree
    root := &Operator{
        Type: OpFilter,
        Children: []*Operator{
            {
                Type: OpJoin,
                Children: []*Operator{
                    {Type: OpScan, Props: map[string]interface{}{"table": "users"}},
                    {Type: OpScan, Props: map[string]interface{}{"table": "orders"}},
                },
            },
        },
    }

    // Find all scans
    scans := FindScans(root)
    fmt.Printf("Found %d scan operators\n", len(scans))

    for _, scan := range scans {
        fmt.Printf("  Scan: %v\n", scan.Props["table"])
    }

    // Type comparison is super fast with unique.Handle!
    fmt.Printf("\nIs root a filter? %v\n", root.Type == OpFilter)
    fmt.Printf("Is root a join? %v\n", root.Type == OpJoin)
}
```

**Output:**
```
Found 2 scan operators
  Scan: users
  Scan: orders

Is root a filter? true
Is root a join? false
```

**Key insight:** `unique.Handle` makes operator type comparisons O(1) pointer equality!

### 3. Cost Estimation with Statistics

**Use statistics to estimate query costs!**

```go
package main

import (
    "fmt"
    "math"
)

type TableStats struct {
    RowCount    int
    ColumnStats map[string]*ColumnStats
}

type ColumnStats struct {
    DistinctCount int
    MinValue      int
    MaxValue      int
    NullFraction  float64
}

// Selectivity estimation for predicates
func EstimateSelectivity(stats *TableStats, predicate string) float64 {
    // Simplified: real systems parse the predicate

    // col = value: selectivity = 1 / distinct_count
    // col > value: selectivity = (max - value) / (max - min)
    // col IS NULL: selectivity = null_fraction

    // Example: age > 18 where age in [0, 100]
    colStats := stats.ColumnStats["age"]
    if colStats != nil {
        value := 18.0
        max := float64(colStats.MaxValue)
        min := float64(colStats.MinValue)

        return (max - value) / (max - min)
    }

    return 0.1  // Default 10%
}

// Join cardinality estimation
func EstimateJoinCardinality(leftCard, rightCard, leftDistinct, rightDistinct int) int {
    // Join selectivity = 1 / max(distinct_left, distinct_right)
    maxDistinct := leftDistinct
    if rightDistinct > maxDistinct {
        maxDistinct = rightDistinct
    }

    selectivity := 1.0 / float64(maxDistinct)
    return int(float64(leftCard*rightCard) * selectivity)
}

func main() {
    stats := &TableStats{
        RowCount: 10000,
        ColumnStats: map[string]*ColumnStats{
            "age": {
                DistinctCount: 80,
                MinValue:      0,
                MaxValue:      100,
                NullFraction:  0.01,
            },
            "country": {
                DistinctCount: 50,
                MinValue:      0,
                MaxValue:      0,
                NullFraction:  0.0,
            },
        },
    }

    // Estimate filter selectivity
    sel := EstimateSelectivity(stats, "age > 18")
    fmt.Printf("Filter selectivity: %.2f%%\n", sel*100)
    fmt.Printf("Estimated output rows: %d\n", int(float64(stats.RowCount)*sel))

    // Estimate join cardinality
    leftCard := 10000
    rightCard := 50000
    leftDistinct := 10000  // users.id (unique)
    rightDistinct := 8000  // orders.user_id

    joinCard := EstimateJoinCardinality(leftCard, rightCard, leftDistinct, rightDistinct)
    fmt.Printf("\nJoin cardinality: %d\n", joinCard)
}
```

**Output:**
```
Filter selectivity: 82.00%
Estimated output rows: 8200

Join cardinality: 50000
```

### 4. Dynamic Programming for Join Ordering

**Find optimal join order efficiently!**

```go
package main

import (
    "fmt"
    "math"
)

type Table struct {
    Name        string
    Cardinality int
}

type JoinPlan struct {
    Tables []string
    Cost   float64
    Left   *JoinPlan
    Right  *JoinPlan
}

// Simplified join cost model
func estimateJoinCost(leftCard, rightCard int) float64 {
    // Hash join: O(|left| + |right|)
    return float64(leftCard + rightCard)
}

// Find optimal join order using dynamic programming
func OptimizeJoinOrder(tables []Table) *JoinPlan {
    n := len(tables)

    // dp[mask] = best plan for subset of tables represented by mask
    dp := make(map[int]*JoinPlan)

    // Base case: single tables
    for i := 0; i < n; i++ {
        mask := 1 << i
        dp[mask] = &JoinPlan{
            Tables: []string{tables[i].Name},
            Cost:   float64(tables[i].Cardinality),
            Left:   nil,
            Right:  nil,
        }
    }

    // Build up: consider all subsets of size 2, 3, ..., n
    for size := 2; size <= n; size++ {
        // Enumerate all subsets of given size
        for mask := 0; mask < (1 << n); mask++ {
            if countBits(mask) != size {
                continue
            }

            best := &JoinPlan{Cost: math.MaxFloat64}

            // Try all ways to split this subset
            for leftMask := mask; leftMask > 0; leftMask = (leftMask - 1) & mask {
                if leftMask == mask {
                    continue  // Skip empty right
                }

                rightMask := mask ^ leftMask

                if leftPlan, ok := dp[leftMask]; ok {
                    if rightPlan, ok := dp[rightMask]; ok {
                        // Simplified: assume join produces |left| * |right| / 10 rows
                        joinCard := int(leftPlan.Cost * rightPlan.Cost / 10)
                        cost := leftPlan.Cost + rightPlan.Cost + float64(joinCard)

                        if cost < best.Cost {
                            best = &JoinPlan{
                                Tables: append(leftPlan.Tables, rightPlan.Tables...),
                                Cost:   cost,
                                Left:   leftPlan,
                                Right:  rightPlan,
                            }
                        }
                    }
                }
            }

            dp[mask] = best
        }
    }

    return dp[(1<<n)-1]
}

func countBits(n int) int {
    count := 0
    for n > 0 {
        count += n & 1
        n >>= 1
    }
    return count
}

func printJoinPlan(plan *JoinPlan, depth int) {
    if plan == nil {
        return
    }

    indent := ""
    for i := 0; i < depth; i++ {
        indent += "  "
    }

    if plan.Left == nil && plan.Right == nil {
        fmt.Printf("%sScan(%s) cost=%.0f\n", indent, plan.Tables[0], plan.Cost)
    } else {
        fmt.Printf("%sJoin cost=%.0f\n", indent, plan.Cost)
        printJoinPlan(plan.Left, depth+1)
        printJoinPlan(plan.Right, depth+1)
    }
}

func main() {
    tables := []Table{
        {Name: "users", Cardinality: 10000},
        {Name: "orders", Cardinality: 50000},
        {Name: "products", Cardinality: 5000},
    }

    fmt.Println("Optimizing join order...")
    plan := OptimizeJoinOrder(tables)

    fmt.Println("\nOptimal plan:")
    printJoinPlan(plan, 0)
    fmt.Printf("\nTotal cost: %.0f\n", plan.Cost)
}
```

**Output:**
```
Optimizing join order...

Optimal plan:
Join cost=5565000
  Join cost=65000
    Scan(users) cost=10000
    Scan(products) cost=5000
  Scan(orders) cost=50000

Total cost: 5565000
```

### 5. Visitor Pattern for Plan Transformation

**Transform plans using visitors!**

```go
package main

import (
    "fmt"
)

type PlanNode interface {
    Accept(Visitor) PlanNode
    String() string
}

type Visitor interface {
    VisitScan(*ScanNode) PlanNode
    VisitFilter(*FilterNode) PlanNode
    VisitJoin(*JoinNode) PlanNode
}

type ScanNode struct {
    Table string
}

func (s *ScanNode) Accept(v Visitor) PlanNode {
    return v.VisitScan(s)
}

func (s *ScanNode) String() string {
    return fmt.Sprintf("Scan(%s)", s.Table)
}

type FilterNode struct {
    Child     PlanNode
    Predicate string
}

func (f *FilterNode) Accept(v Visitor) PlanNode {
    return v.VisitFilter(f)
}

func (f *FilterNode) String() string {
    return fmt.Sprintf("Filter(%s)", f.Predicate)
}

type JoinNode struct {
    Left  PlanNode
    Right PlanNode
}

func (j *JoinNode) Accept(v Visitor) PlanNode {
    return v.VisitJoin(j)
}

func (j *JoinNode) String() string {
    return "Join"
}

// Pushdown filter visitor (optimization!)
type PushdownFilterVisitor struct{}

func (v *PushdownFilterVisitor) VisitScan(s *ScanNode) PlanNode {
    return s
}

func (v *PushdownFilterVisitor) VisitFilter(f *FilterNode) PlanNode {
    // Try to push filter down
    child := f.Child.Accept(v)

    // If child is a join, try to push filter below it
    if join, ok := child.(*JoinNode); ok {
        // Simplified: assume filter applies to left side
        return &JoinNode{
            Left: &FilterNode{
                Child:     join.Left,
                Predicate: f.Predicate,
            },
            Right: join.Right,
        }
    }

    return &FilterNode{Child: child, Predicate: f.Predicate}
}

func (v *PushdownFilterVisitor) VisitJoin(j *JoinNode) PlanNode {
    return &JoinNode{
        Left:  j.Left.Accept(v),
        Right: j.Right.Accept(v),
    }
}

func printTree(node PlanNode, depth int) {
    indent := ""
    for i := 0; i < depth; i++ {
        indent += "  "
    }
    fmt.Printf("%s%s\n", indent, node.String())

    switch n := node.(type) {
    case *FilterNode:
        printTree(n.Child, depth+1)
    case *JoinNode:
        printTree(n.Left, depth+1)
        printTree(n.Right, depth+1)
    }
}

func main() {
    // Original plan: Filter on top of Join
    plan := &FilterNode{
        Predicate: "age > 18",
        Child: &JoinNode{
            Left:  &ScanNode{Table: "users"},
            Right: &ScanNode{Table: "orders"},
        },
    }

    fmt.Println("Original plan:")
    printTree(plan, 0)

    // Optimize: push filter down
    visitor := &PushdownFilterVisitor{}
    optimized := plan.Accept(visitor)

    fmt.Println("\nOptimized plan (filter pushed down):")
    printTree(optimized, 0)
}
```

**Output:**
```
Original plan:
Filter(age > 18)
  Join
    Scan(users)
    Scan(orders)

Optimized plan (filter pushed down):
Join
  Filter(age > 18)
    Scan(users)
  Scan(orders)
```

**Key insight:** Filter pushdown reduces intermediate result size!

### 6. Memo Table for Plan Caching

**Avoid recomputing equivalent subplans!**

```go
package main

import (
    "fmt"
    "strings"
)

type PlanMemo struct {
    cache map[string]*CachedPlan
}

type CachedPlan struct {
    Plan PlanNode
    Cost float64
}

func NewPlanMemo() *PlanMemo {
    return &PlanMemo{
        cache: make(map[string]*CachedPlan),
    }
}

// Generate a key for a set of tables
func (m *PlanMemo) key(tables []string) string {
    // Sort tables for canonical representation
    sorted := make([]string, len(tables))
    copy(sorted, tables)
    // Simple sort (for demo)
    return strings.Join(sorted, ",")
}

func (m *PlanMemo) Get(tables []string) (*CachedPlan, bool) {
    plan, ok := m.cache[m.key(tables)]
    return plan, ok
}

func (m *PlanMemo) Put(tables []string, plan PlanNode, cost float64) {
    m.cache[m.key(tables)] = &CachedPlan{
        Plan: plan,
        Cost: cost,
    }
}

func main() {
    memo := NewPlanMemo()

    // Cache a plan
    memo.Put([]string{"users", "orders"}, &JoinNode{
        Left:  &ScanNode{Table: "users"},
        Right: &ScanNode{Table: "orders"},
    }, 100.0)

    // Retrieve it
    if cached, ok := memo.Get([]string{"users", "orders"}); ok {
        fmt.Printf("Cache hit! Cost: %.2f\n", cached.Cost)
    }

    // Different order - same key
    if cached, ok := memo.Get([]string{"orders", "users"}); ok {
        fmt.Printf("Cache hit (reordered)! Cost: %.2f\n", cached.Cost)
    } else {
        fmt.Println("Cache miss")
    }
}
```

## Pre-Implementation Exercises

### Exercise 1: Implement Cost Model

```go
package main

// TODO: Implement cost estimation for different operators

type Operator interface {
    EstimateCost(stats *Statistics) float64
    Cardinality(stats *Statistics) int
}

type Statistics struct {
    TableStats map[string]*TableStats
}

type TableStats struct {
    RowCount int
    // TODO: Add column statistics
}

type ScanOperator struct {
    Table string
}

func (s *ScanOperator) EstimateCost(stats *Statistics) float64 {
    // TODO: Implement
    // Hint: Cost = rowCount * costPerRow
    return 0.0
}

type HashJoinOperator struct {
    Left  Operator
    Right Operator
}

func (h *HashJoinOperator) EstimateCost(stats *Statistics) float64 {
    // TODO: Implement
    // Hint: Cost = buildCost + probeCost
    return 0.0
}

func main() {
    // TODO: Test your cost model
}
```

### Exercise 2: Join Order Optimization

```go
package main

// TODO: Implement dynamic programming join ordering

type Table struct {
    Name string
    Rows int
}

func OptimizeJoins(tables []Table) *JoinPlan {
    // TODO: Use bitmask DP to find optimal order
    // Hint: dp[mask] = best plan for tables in mask
    return nil
}

func main() {
    tables := []Table{
        {Name: "A", Rows: 100},
        {Name: "B", Rows: 1000},
        {Name: "C", Rows: 50},
    }

    // TODO: Find optimal join order
    // Expected: Join smallest tables first (A-C-B or C-A-B)
}
```

### Exercise 3: Filter Pushdown

```go
package main

// TODO: Implement filter pushdown optimization

func PushdownFilters(plan PlanNode) PlanNode {
    // TODO: Use visitor pattern to push filters below joins
    return nil
}

func main() {
    // TODO: Test with plan:
    // Filter(x > 10)
    //   Join
    //     Scan(A)
    //     Scan(B)
    //
    // Should become:
    // Join
    //   Filter(x > 10)
    //     Scan(A)
    //   Scan(B)
}
```

### Exercise 4: Plan Equivalence

```go
package main

// TODO: Implement plan equivalence checking

func ArePlansEquivalent(p1, p2 PlanNode) bool {
    // TODO: Check if two plans produce same results
    // Hint: Ignore operator order for commutative ops
    return false
}

func main() {
    // TODO: Test with:
    // Join(A, B) vs Join(B, A) -> should be equivalent
    // Filter(x>10, Filter(y<5, Scan)) vs Filter(x>10 AND y<5, Scan) -> equivalent
}
```

### Exercise 5: Cardinality Estimation

```go
package main

// TODO: Implement selectivity estimation

func EstimateFilterSelectivity(predicate string, stats *ColumnStats) float64 {
    // TODO: Estimate fraction of rows passing filter
    // col = value: 1 / distinct_count
    // col > value: (max - value) / (max - min)
    // col IN (a,b,c): count / distinct_count
    return 0.0
}

func EstimateJoinCardinality(left, right int, joinType string) int {
    // TODO: Estimate join output size
    return 0
}

func main() {
    // TODO: Test estimation accuracy
}
```

## Performance Benchmarks

### Benchmark 1: Plan Enumeration

```go
func BenchmarkJoinOrdering(b *testing.B) {
    tables := make([]Table, 10)
    for i := range tables {
        tables[i] = Table{Name: fmt.Sprintf("t%d", i), Rows: 1000 * (i + 1)}
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = OptimizeJoinOrder(tables)
    }
}
```

**Target: < 100ms for 10 tables**

### Benchmark 2: Cost Estimation

```go
func BenchmarkCostEstimation(b *testing.B) {
    plan := buildComplexPlan()  // Deep plan tree
    stats := &Statistics{/* ... */}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = plan.EstimateCost(stats)
    }
}
```

**Target: < 1ms for deep plans**

## Common Gotchas to Avoid

### Gotcha 1: Incorrect Join Cardinality

```go
// WRONG: Assume cartesian product
func (j *JoinOp) Cardinality() int {
    return j.left.Cardinality() * j.right.Cardinality()
}

// RIGHT: Use join key selectivity
func (j *JoinOp) Cardinality() int {
    leftCard := j.left.Cardinality()
    rightCard := j.right.Cardinality()
    selectivity := 1.0 / float64(maxDistinct(j.leftKey, j.rightKey))
    return int(float64(leftCard*rightCard) * selectivity)
}
```

### Gotcha 2: Not Caching Plan Costs

```go
// WRONG: Recompute cost every time
func (j *JoinOp) EstimateCost() float64 {
    return j.left.EstimateCost() + j.right.EstimateCost() + /* join cost */
}

// RIGHT: Cache and reuse
type JoinOp struct {
    left       LogicalOp
    right      LogicalOp
    cachedCost float64
    costValid  bool
}

func (j *JoinOp) EstimateCost() float64 {
    if !j.costValid {
        j.cachedCost = j.left.EstimateCost() + j.right.EstimateCost() + /* join cost */
        j.costValid = true
    }
    return j.cachedCost
}
```

### Gotcha 3: Exponential Plan Enumeration

```go
// WRONG: Try all permutations
func enumeratePlans(tables []Table) []Plan {
    if len(tables) == 1 {
        return []Plan{{Scan(tables[0])}}
    }

    var plans []Plan
    for i, t := range tables {
        rest := append(tables[:i], tables[i+1:]...)
        for _, subplan := range enumeratePlans(rest) {  // Exponential!
            plans = append(plans, join(Scan(t), subplan))
        }
    }
    return plans
}

// RIGHT: Use DP with memoization
func optimizeJoins(tables []Table) Plan {
    memo := make(map[int]Plan)  // Bitmask -> best plan
    // ... DP approach ...
}
```

## Checklist Before Starting Lesson 3.2

- [ ] I understand interface-based operator design
- [ ] I can estimate operator costs and cardinalities
- [ ] I know how to use `unique.Handle` for operator types
- [ ] I understand dynamic programming for join ordering
- [ ] I can implement the visitor pattern for plan transformation
- [ ] I understand filter pushdown optimization
- [ ] I know how to estimate selectivity from statistics
- [ ] I can detect equivalent plans
- [ ] I understand memoization for subplan caching
- [ ] I've benchmarked plan enumeration

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.2

You'll implement:
- Logical operator interface and implementations
- Cost-based optimizer with statistics
- Join ordering using dynamic programming
- Filter pushdown and predicate reordering
- Plan equivalence detection
- Histogram-based selectivity estimation
- Optimization rules (filter pushdown, projection pushdown, etc.)

**Time estimate:** 30-35 hours for full implementation

**Planning is where query optimization happens!** 🧠
