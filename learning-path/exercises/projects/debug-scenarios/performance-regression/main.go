package perfregression

import "sync"

// SLOW: Query executor with performance issues

type Executor struct {
	data []Row
}

type Row map[string]interface{}

// SLOW: Allocates on every call
func (e *Executor) Execute() []Row {
	results := make([]Row, 0)  // Should pre-allocate!

	for _, row := range e.data {
		// SLOW: Expensive operation in tight loop
		if e.slowFilter(row) {
			// SLOW: Append causes many allocations
			results = append(results, e.transform(row))
		}
	}

	return results
}

func (e *Executor) slowFilter(row Row) bool {
	// SLOW: Unnecessary computation
	for i := 0; i < 1000; i++ {
		_ = i * i
	}
	return true
}

func (e *Executor) transform(row Row) Row {
	// SLOW: Creates new map every time
	newRow := make(Row)
	for k, v := range row {
		newRow[k] = v
	}
	return newRow
}

// SLOW: Cache with lock contention
type Cache struct {
	data map[string]interface{}
	mu   sync.Mutex  // Global lock causes contention!
}

func (c *Cache) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[key]
}

func (c *Cache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// SLOW: Join with wrong algorithm
func HashJoin(left, right []Row) []Row {
	results := make([]Row, 0)

	// SLOW: O(n*m) nested loop instead of hash join!
	for _, l := range left {
		for _, r := range right {
			if l["id"] == r["id"] {
				results = append(results, merge(l, r))
			}
		}
	}

	return results
}

func merge(left, right Row) Row {
	result := make(Row)
	for k, v := range left {
		result[k] = v
	}
	for k, v := range right {
		result[k] = v
	}
	return result
}

// TODO: Fix these performance issues!
// Hints:
// 1. Pre-allocate slices with capacity
// 2. Avoid allocations in hot paths
// 3. Use sharded locks to reduce contention
// 4. Choose correct algorithm (hash join vs nested loop)
// 5. Profile with pprof to verify improvements
