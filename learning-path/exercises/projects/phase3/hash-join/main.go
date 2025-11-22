package hashjoin

// Row represents a database row
type Row map[string]interface{}

// Key represents a join key
type Key interface{}

// HashJoin performs hash join between left and right tables
func HashJoin(left, right []Row, leftKey, rightKey func(Row) Key) []Row {
	// TODO: Implement hash join
	// 1. Build phase: create hash table from smaller relation
	// 2. Probe phase: lookup each tuple from larger relation
	// 3. Use Go 1.24 pre-sized maps for efficiency
	return nil
}

// SortMergeJoin performs sort-merge join
func SortMergeJoin(left, right []Row, leftKey, rightKey func(Row) Key) []Row {
	// TODO: Implement sort-merge join
	// 1. Sort both inputs
	// 2. Merge with two pointers
	return nil
}

// IndexNestedLoopJoin performs index nested loop join
func IndexNestedLoopJoin(outer []Row, index Index, keyFunc func(Row) Key) []Row {
	// TODO: Implement index nested loop join
	return nil
}

// Index interface for index lookups
type Index interface {
	Lookup(key Key) []Row
}
