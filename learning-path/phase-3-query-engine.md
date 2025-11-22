# Phase 3: Query Engine

**Prerequisites:** Complete Phase 2 (Graph Structure)

**Time Estimate:** 4 weeks
**Go Version:** 1.23+

In this phase, you'll build a complete query engine that parses Cypher queries, optimizes query plans, implements join algorithms, and executes queries efficiently.

---

## Table of Contents

- [Lesson 3.1: Parse Cypher Queries](#lesson-31-parse-cypher-queries-week-7)
- [Lesson 3.2: Query Planning](#lesson-32-query-planning-week-8)
- [Lesson 3.3: Join Algorithms](#lesson-33-join-algorithms-week-9)
- [Lesson 3.4: Execution Engine](#lesson-34-execution-engine-week-10)

---

## Lesson 3.1: Parse Cypher Queries (Week 7)

**Your Mission:** Turn this string into a data structure:
```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)
WHERE a.age > 30
RETURN b.name
```

### Investigation Questions

1. **Tokenization:**
   ```go
   input := "MATCH (a:Person) WHERE a.age > 30"
   tokens := Tokenize(input)
   // Should produce: [MATCH, LPAREN, IDENT("a"), COLON, IDENT("Person"), ...]
   ```
   - How do you handle whitespace?
   - How do you distinguish `MATCH` (keyword) from `match` (identifier)?
   - What about string literals with special chars: `"foo\"bar"`?

2. **Parser Choice:**

   **Option A: Use ANTLR4**
   ```bash
   # Install ANTLR4 for Go
   go get github.com/antlr4-go/antlr/v4
   # Write Cypher.g4 grammar
   ```
   Pro: Handles complexity. Con: Big dependency, slow compilation.

   **Option B: Recursive Descent by Hand**
   ```go
   func (p *Parser) parseMatch() (*MatchClause, error) {
       p.expect(MATCH)
       pattern := p.parsePattern()
       // ...
   }
   ```
   Pro: Full control, fast. Con: You write more code.

   Try both. Which feels better for learning?

3. **AST Design:**
   ```go
   type Query struct {
       Match  *MatchClause
       Where  *WhereClause
       Return *ReturnClause
   }

   type MatchClause struct {
       Patterns []Pattern
   }

   type Pattern struct {
       // (a)-[:KNOWS]->(b)
       // How do you represent this?
   }
   ```
   Draw the AST for 5 different queries. What common structure emerges?

### Challenge

- [ ] Parse 20 test queries correctly
- [ ] Produce useful error messages: "Expected RETURN, got WHERE on line 3"
- [ ] Handle nested expressions: `a.age > 30 AND (b.city = "NYC" OR b.city = "LA")`
- [ ] Pretty-print the AST for debugging

### Test Cases

```go
tests := []string{
    "MATCH (n) RETURN n",
    "MATCH (a)-[r]->(b) WHERE a.age > 30 RETURN b.name",
    "MATCH (a)-[:KNOWS*1..3]->(b) RETURN count(b)",
    // Add 17 more edge cases
}
```

### What You Should Discover

- Parsing is 20% lexing, 80% handling edge cases
- Good error messages are harder than parsing itself
- ASTs need methods: `node.String()`, `node.Validate()`

---

## Lesson 3.2: Query Planning (Week 8)

**Your Mission:** Turn an AST into an execution plan.

**The Problem:**
```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)-[:LIKES]->(c:Movie)
WHERE a.age > 30 AND c.year = 2020
RETURN c.title
```

**Many valid execution orders:**
1. Scan all Persons → filter age → traverse KNOWS → traverse LIKES → filter year
2. Scan Movies → filter year → reverse LIKES → reverse KNOWS → filter age
3. Hash join Person and Movie on intermediate LIKES edges

Which is fastest?

### Investigation Questions

1. **Cost Estimation:**
   ```go
   type Operator interface {
       EstimatedCost() float64
       EstimatedCardinality() int
   }

   type SeqScan struct {
       Table string
       Filter Expr
   }

   func (s *SeqScan) EstimatedCost() float64 {
       // ???
   }
   ```
   - How do you estimate selectivity of `age > 30`?
   - Research: Histograms vs sampling vs bloom filters
   - What if you have no statistics yet?

2. **Plan Enumeration:**
   Three patterns: (a)->(b)->(c)

   **Possible join orders:**
   - Left-deep: ((a⋈b)⋈c)
   - Right-deep: (a⋈(b⋈c))
   - Bushy: Various combinations

   How many possible plans for N patterns? (It's factorial)

   **Optimization:**
   - Dynamic programming (PostgreSQL approach)
   - Greedy heuristic (pick smallest intermediate result)
   - Random sampling (try 100 plans, pick best)

3. **Predicate Pushdown:**
   ```go
   // Before optimization:
   SeqScan(Person) → Filter(age > 30) → Join(KNOWS)

   // After pushdown:
   SeqScan(Person, filter: age > 30) → Join(KNOWS)
   ```
   - When can you push filters down?
   - What about `WHERE a.age + b.age > 60`? (Needs both sides)

### 🆕 Go 1.23: Intern Plan Operator Types

```go
import "unique"

type LogicalPlan struct {
    OpType unique.Handle[string]  // "SeqScan", "HashJoin", etc.
    // Most queries use same 10-20 operator types
}

// Fast operator comparison
func (p1 *LogicalPlan) SameOperator(p2 *LogicalPlan) bool {
    return p1.OpType == p2.OpType  // Pointer comparison!
}
```

### Your Implementation

- [ ] Generate at least 3 different plans for same query
- [ ] Choose plan based on estimated cost
- [ ] Show plan with `EXPLAIN` command
- [ ] Measure actual vs estimated cardinality

### Visualization

```
EXPLAIN MATCH (a:Person)-[:KNOWS]->(b)
  WHERE a.age > 30 RETURN b.name;

PhysicalPlan:
  Project(b.name)
  └─ HashJoin(a.id = KNOWS.src)
     ├─ SeqScan(Person, filter: age > 30)  [est: 100K rows]
     └─ IndexScan(KNOWS.src)                [est: 1M rows]

Estimated cost: 1234.5
```

### What You Should Discover

- Query optimization is NP-hard
- Cardinality estimates are often wildly wrong
- Good plans are 100x faster than bad plans

---

## Lesson 3.3: Join Algorithms (Week 9)

**Your Mission:** Implement three join algorithms, benchmark them.

**The Query:**
```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)
```
Translation: Join Person table with KNOWS edge list on Person.id = KNOWS.dst

### Investigation Questions

1. **Hash Join:**
   ```go
   // Phase 1: Build hash table on smaller side
   hashTable := make(map[NodeID]*Person, len(persons))  // 🆕 Pre-size!
   for i := range persons {
       hashTable[persons[i].ID] = &persons[i]
   }

   // Phase 2: Probe with larger side
   for _, edge := range edges {
       if person, found := hashTable[edge.Dst]; found {
           yield(edge.Src, person)
       }
   }
   ```
   - What if the hash table doesn't fit in memory?
   - Research: Grace hash join (partition-based)
   - When does Go's map become slow? (Measure at 1M, 10M, 100M entries)

2. **Sort-Merge Join:**
   ```go
   // Sort both sides
   sort.Slice(persons, func(i, j int) bool {
       return persons[i].ID < persons[j].ID
   })
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].Dst < edges[j].Dst
   })

   // Merge with two pointers
   i, j := 0, 0
   for i < len(persons) && j < len(edges) {
       // Your code here
   }
   ```
   - What if data is already sorted? (CSR is sorted!)
   - What about duplicate keys?
   - Measure: sort cost vs probe cost

3. **Worst-Case Optimal Join (Triangle Query):**
   ```cypher
   MATCH (a)-[:KNOWS]->(b)-[:KNOWS]->(c)-[:KNOWS]->(a)
   RETURN count(*)
   ```
   This is a cycle—hash joins explode.

   **Binary hash join:** O(E^1.5) intermediate results
   **Multiway intersection:** O(E) time

   ```go
   func intersectThree(a, b, c []NodeID) []NodeID {
       // How do you do this efficiently?
       // Hint: They're sorted
   }
   ```

### 🆕 Go 1.24: Swiss Tables Performance

Hash joins got 30-35% faster!

```go
func BenchmarkHashJoin(b *testing.B) {
    persons := generatePersons(1_000_000)
    edges := generateEdges(10_000_000)

    b.Run("Build", func(b *testing.B) {
        for b.Loop() {
            // 🆕 35% faster with pre-sized map in Go 1.24
            hashTable := make(map[NodeID]*Person, len(persons))
            for i := range persons {
                hashTable[persons[i].ID] = &persons[i]
            }
        }
    })

    b.Run("Probe", func(b *testing.B) {
        hashTable := buildHashTable(persons)

        for b.Loop() {
            // 🆕 30% faster lookups in Go 1.24
            count := 0
            for _, edge := range edges {
                if _, found := hashTable[edge.Dst]; found {
                    count++
                }
            }
        }
    })
}
```

**Investigation:** Measure performance improvement from Go 1.23 to 1.24.

### Your Implementation

- [ ] Hash join for 1-to-many
- [ ] Sort-merge for many-to-many
- [ ] Multiway intersection for cycles
- [ ] Adaptive: pick algorithm based on input size
- [ ] 🆕 Benchmark shows Swiss Tables improvement

### Benchmark (Go 1.24+)

```go
// Dataset: 10K persons, 100K KNOWS edges
func BenchmarkJoins(b *testing.B) {
    scenarios := []struct{
        name string
        persons int
        edges int
    }{
        {"small", 1000, 10000},
        {"medium", 10000, 100000},
        {"large", 100000, 1000000},
    }

    for _, s := range scenarios {
        b.Run(s.name + "/hash", func(b *testing.B) {
            // hash join
        })
        b.Run(s.name + "/merge", func(b *testing.B) {
            // sort-merge join
        })
    }
}
```

### What You Should Discover

- Hash joins are usually fastest (if memory allows)
- Sort-merge shines when data is pre-sorted
- Cyclic queries need special treatment (WCO joins)
- 🆕 Swiss Tables make hash joins 30-35% faster

---

## Lesson 3.4: Execution Engine (Week 10)

**Your Mission:** Execute plans with pipelining.

### 🆕 Go 1.23: Iterator-Based Execution

Modern approach using range-over-func:

```go
import "iter"

type Operator interface {
    Execute() iter.Seq[Row]
}

type SeqScan struct {
    table *NodeTable
}

func (s *SeqScan) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for i := 0; i < s.table.NumRows(); i++ {
            if !yield(s.table.Row(i)) {
                return  // Consumer stopped
            }
        }
    }
}

// Composable pipeline:
type Filter struct {
    child Operator
    predicate func(Row) bool
}

func (f *Filter) Execute() iter.Seq[Row] {
    return func(yield func(Row) bool) {
        for row := range f.child.Execute() {
            if f.predicate(row) {
                if !yield(row) {
                    return
                }
            }
        }
    }
}

// Usage:
scan := &SeqScan{table: persons}
filter := &Filter{child: scan, predicate: func(r Row) bool {
    return r.Age > 30
}}

for row := range filter.Execute() {
    process(row)
}
```

### Investigation Questions

1. **Iterator vs Vectorized:**
   ```go
   // Iterator: one row at a time
   for row := range operator.Execute() {
       process(row)
   }

   // Vectorized: batch of rows
   for {
       batch := operator.NextBatch()  // 1000 rows
       if batch == nil { break }
       for _, row := range batch {
           process(row)
       }
   }
   ```
   Implement both. Measure virtual function call overhead.

2. **Pipeline Breakers:**
   Some operators can't pipeline:
   - Sort (needs all input before outputting)
   - Hash join build phase (needs all input for hash table)
   - Aggregation (GROUP BY needs all rows)

   ```go
   type SortOperator struct {
       child Operator
       buffer []Row  // Must accumulate here
   }

   func (s *SortOperator) Execute() iter.Seq[Row] {
       return func(yield func(Row) bool) {
           // Collect all rows
           for row := range s.child.Execute() {
               s.buffer = append(s.buffer, row)
           }

           // Sort
           sort.Slice(s.buffer, ...)

           // Yield sorted
           for _, row := range s.buffer {
               if !yield(row) {
                   return
               }
           }
       }
   }
   ```
   - How do you detect pipeline breaks?
   - Can you parallelize across pipeline boundaries?

3. **Parallel Execution:**
   ```go
   type ParallelSeqScan struct {
       table string
       workers int
       chunks chan NodeGroupID
   }

   func (p *ParallelSeqScan) Execute() iter.Seq[Row] {
       return func(yield func(Row) bool) {
           results := make(chan Row, 1000)

           // Start workers
           var wg sync.WaitGroup
           for i := 0; i < p.workers; i++ {
               wg.Add(1)
               go func() {
                   defer wg.Done()
                   for chunk := range p.chunks {
                       for row := range scanChunk(chunk) {
                           results <- row
                       }
                   }
               }()
           }

           // Close results when done
           go func() {
               wg.Wait()
               close(results)
           }()

           // Yield results
           for row := range results {
               if !yield(row) {
                   // Early exit - drain channels
                   return
               }
           }
       }
   }
   ```
   - Each worker scans different NodeGroups
   - Results must be thread-safe
   - Use sync.WaitGroup or channels?

### Your Implementation

- [ ] Iterator-based execution (🆕 Go 1.23)
- [ ] Support 10+ operator types (scan, filter, join, project, limit)
- [ ] Pipeline as much as possible
- [ ] Parallelize scans across goroutines
- [ ] Profile: what % of time in each operator?

### Query Example

```cypher
MATCH (a:Person)-[:KNOWS]->(b:Person)
WHERE a.age > 30 AND b.city = "NYC"
RETURN b.name, count(*) as friends
LIMIT 10
```

**Execution plan:**
```
Limit(10)
└─ Aggregate(b.name, count(*))
   └─ Filter(b.city = "NYC")
      └─ HashJoin(a.id = KNOWS.dst)
         ├─ Filter(a.age > 30)
         │  └─ SeqScan(Person)  [parallel: 8 workers]
         └─ SeqScan(KNOWS)     [parallel: 8 workers]
```

### What You Should Discover

- Iterator model is simple but has overhead
- Vectorization makes modern CPUs happy
- Parallelism is easy; correct parallelism is hard
- 🆕 Range-over-func provides clean composition

---

## Phase 3 Complete!

**What you've built:**
- ✅ Cypher query parser
- ✅ Query optimizer with cost estimation
- ✅ Multiple join algorithms
- ✅ Iterator-based execution engine

**Key Insights:**
- Parsing is tedious but straightforward
- Query optimization is an art and a science
- Join algorithm choice dramatically affects performance
- Iterators make operator composition elegant

**Next:** Phase 4 - Transactions, where you'll implement locking protocols and MVCC to support concurrent access.

---

## Resources

### Papers to Read
- "Access Path Selection in a Relational Database" - Selinger et al. (query optimization)
- "Worst-Case Optimal Join Algorithms" - Ngo et al.
- "MonetDB/X100: Hyper-Pipelining Query Execution"

### Reference Implementations
- PostgreSQL's query planner
- DuckDB's execution engine
- Kuzu's query optimizer

### Go-Specific
- Writing parsers in Go
- `iter` package for iterators (Go 1.23+)
- Performance profiling with `pprof`
