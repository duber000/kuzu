package optimizerreview

// REVIEW THIS CODE

type TableStats struct {
	RowCount int64
}

type Optimizer struct {
	stats map[string]*TableStats
}

// Issue 1: Wrong cost formula
func (o *Optimizer) EstimateJoinCost(left, right string) float64 {
	leftRows := o.stats[left].RowCount
	rightRows := o.stats[right].RowCount

	// WRONG: Should be leftRows * rightRows, not sum
	return float64(leftRows + rightRows)
}

// Issue 2: Missing selectivity
func (o *Optimizer) EstimateFilter(table string, predicate string) int64 {
	// WRONG: Always returns 50%, should use statistics
	return o.stats[table].RowCount / 2
}

// Issue 3: Doesn't check for cross product
func (o *Optimizer) OptimizeJoinOrder(tables []string) []string {
	// WRONG: Just returns as-is, no optimization
	return tables
}

// Issue 4: Exponential complexity
func (o *Optimizer) EnumeratePlans(tables []string) int {
	if len(tables) == 1 {
		return 1
	}

	count := 0
	for i := range tables {
		// WRONG: Generates all permutations (n!)
		remaining := append(tables[:i], tables[i+1:]...)
		count += o.EnumeratePlans(remaining)
	}

	return count
}

// TODO: Find all optimization bugs
