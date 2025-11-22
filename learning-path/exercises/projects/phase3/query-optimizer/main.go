package optimizer

// LogicalPlan represents a logical query plan
type LogicalPlan interface {
	Children() []LogicalPlan
	String() string
}

// PhysicalPlan represents an executable plan
type PhysicalPlan interface {
	Execute() ResultSet
	Cost() float64
}

// Cost represents plan cost
type Cost struct {
	CPUCost float64
	IOCost  float64
}

// Optimizer optimizes query plans
type Optimizer struct {
	stats map[string]*TableStats
}

// NewOptimizer creates a new optimizer
func NewOptimizer() *Optimizer {
	return &Optimizer{
		stats: make(map[string]*TableStats),
	}
}

// Optimize transforms a logical plan to optimal physical plan
func (o *Optimizer) Optimize(plan LogicalPlan) (PhysicalPlan, error) {
	// TODO: Implement optimization
	// 1. Apply rule-based transformations
	// 2. Enumerate physical plans
	// 3. Cost each plan
	// 4. Select minimum cost plan
	return nil, nil
}

// EstimateCost estimates the cost of a physical plan
func (o *Optimizer) EstimateCost(plan PhysicalPlan) Cost {
	// TODO: Implement cost estimation
	return Cost{}
}

// TableStats stores table statistics
type TableStats struct {
	RowCount     int64
	DistinctVals map[string]int64
}

type ResultSet interface {
	Next() bool
	Values() []interface{}
}
