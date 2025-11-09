package executor

import (
	"iter"
	"time"
)

// Row represents a database row
type Row map[string]interface{}

// Operator is the base interface for all operators
type Operator interface {
	Execute() iter.Seq[Row]
	Explain() string
	Profile() OperatorStats
}

// OperatorStats tracks operator execution statistics
type OperatorStats struct {
	RowsProduced int64
	ExecutionTime time.Duration
	MemoryUsed int64
}

// ScanOperator scans a table
type ScanOperator struct {
	tableName string
	rows      []Row
	stats     OperatorStats
}

func (o *ScanOperator) Execute() iter.Seq[Row] {
	return func(yield func(Row) bool) {
		// TODO: Implement scan with profiling
	}
}

func (o *ScanOperator) Explain() string {
	return "Scan(" + o.tableName + ")"
}

func (o *ScanOperator) Profile() OperatorStats {
	return o.stats
}

// FilterOperator filters rows
type FilterOperator struct {
	child Operator
	pred  func(Row) bool
	stats OperatorStats
}

func (o *FilterOperator) Execute() iter.Seq[Row] {
	return func(yield func(Row) bool) {
		// TODO: Implement filter
	}
}

// ProjectOperator projects columns
type ProjectOperator struct {
	child   Operator
	columns []string
	stats   OperatorStats
}

func (o *ProjectOperator) Execute() iter.Seq[Row] {
	// TODO: Implement projection
	return nil
}

// HashJoinOperator performs hash join
type HashJoinOperator struct {
	left    Operator
	right   Operator
	leftKey string
	rightKey string
	stats   OperatorStats
}

func (o *HashJoinOperator) Execute() iter.Seq[Row] {
	// TODO: Implement hash join as pipeline operator
	return nil
}
