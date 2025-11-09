package columnarstore

import (
	"errors"
	"iter"
	"math/bits"
	"unique"
)

// Errors
var (
	ErrColumnNotFound = errors.New("column not found")
	ErrTypeMismatch   = errors.New("type mismatch")
	ErrInvalidRow     = errors.New("invalid row index")
)

// Column interface for different column types
type Column interface {
	Append(value any) error
	Get(index int) (any, bool)
	Scan() iter.Seq2[int, any]
	MemoryUsage() int64
	RowCount() int
}

// Bitmap for NULL values and boolean columns
type Bitmap struct {
	bits []byte
	size int
}

// NewBitmap creates a new bitmap
func NewBitmap(size int) *Bitmap {
	return &Bitmap{
		bits: make([]byte, (size+7)/8),
		size: size,
	}
}

func (b *Bitmap) Set(pos int) {
	// TODO: Implement set bit
}

func (b *Bitmap) Clear(pos int) {
	// TODO: Implement clear bit
}

func (b *Bitmap) Test(pos int) bool {
	// TODO: Implement test bit
	return false
}

func (b *Bitmap) CountOnes() int {
	// TODO: Implement popcount
	// Use bits.OnesCount8
	return 0
}

// IntColumn stores integers with bit packing
type IntColumn struct {
	values   []byte
	nulls    *Bitmap
	bitWidth int
	minValue int64
	rowCount int
}

// NewIntColumn creates a new integer column
func NewIntColumn(bitWidth int, minValue int64) *IntColumn {
	return &IntColumn{
		values:   make([]byte, 0),
		nulls:    NewBitmap(0),
		bitWidth: bitWidth,
		minValue: minValue,
	}
}

func (c *IntColumn) Append(value any) error {
	// TODO: Implement integer append with bit packing
	return nil
}

func (c *IntColumn) Get(index int) (any, bool) {
	// TODO: Implement get with unpacking
	return nil, false
}

func (c *IntColumn) Scan() iter.Seq2[int, any] {
	// TODO: Implement scan iterator
	return nil
}

func (c *IntColumn) MemoryUsage() int64 {
	// TODO: Calculate memory usage
	return 0
}

func (c *IntColumn) RowCount() int {
	return c.rowCount
}

// StringColumn stores strings with dictionary encoding
type StringColumn struct {
	dict     []unique.Handle[string]
	indices  []uint32
	nulls    *Bitmap
	dictMap  map[unique.Handle[string]]uint32
	rowCount int
}

// NewStringColumn creates a new string column
func NewStringColumn() *StringColumn {
	return &StringColumn{
		dict:    make([]unique.Handle[string], 0),
		indices: make([]uint32, 0),
		nulls:   NewBitmap(0),
		dictMap: make(map[unique.Handle[string]]uint32),
	}
}

func (c *StringColumn) Append(value any) error {
	// TODO: Implement string append with dictionary encoding
	// Use unique.Make for string interning
	return nil
}

func (c *StringColumn) Get(index int) (any, bool) {
	// TODO: Implement get from dictionary
	return nil, false
}

func (c *StringColumn) Scan() iter.Seq2[int, any] {
	// TODO: Implement scan iterator
	return nil
}

func (c *StringColumn) MemoryUsage() int64 {
	// TODO: Calculate memory usage
	return 0
}

func (c *StringColumn) RowCount() int {
	return c.rowCount
}

func (c *StringColumn) DistinctCount() int {
	// TODO: Return dictionary size
	return len(c.dict)
}

// FloatColumn stores float64 values
type FloatColumn struct {
	values   []float64
	nulls    *Bitmap
	rowCount int
}

// NewFloatColumn creates a new float column
func NewFloatColumn() *FloatColumn {
	return &FloatColumn{
		values: make([]float64, 0),
		nulls:  NewBitmap(0),
	}
}

func (c *FloatColumn) Append(value any) error {
	// TODO: Implement float append
	return nil
}

func (c *FloatColumn) Get(index int) (any, bool) {
	// TODO: Implement get
	return nil, false
}

func (c *FloatColumn) Scan() iter.Seq2[int, any] {
	// TODO: Implement scan iterator
	return nil
}

func (c *FloatColumn) MemoryUsage() int64 {
	return int64(len(c.values)*8 + len(c.nulls.bits))
}

func (c *FloatColumn) RowCount() int {
	return c.rowCount
}

// PropertyStore stores columns for entities
type PropertyStore struct {
	columns  map[string]Column
	rowCount int
}

// NewPropertyStore creates a new property store
func NewPropertyStore() *PropertyStore {
	return &PropertyStore{
		columns: make(map[string]Column),
	}
}

// AddColumn adds a column to the store
func (ps *PropertyStore) AddColumn(name string, col Column) error {
	// TODO: Implement column addition
	return nil
}

// AppendRow appends a row with values for each column
func (ps *PropertyStore) AppendRow(values map[string]any) error {
	// TODO: Implement row append
	// Append to each column (use nil for missing values)
	return nil
}

// Get retrieves a value at a specific row and column
func (ps *PropertyStore) Get(row int, col string) (any, bool, error) {
	// TODO: Implement get
	return nil, false, nil
}

// Scan returns an iterator over a column's values
func (ps *PropertyStore) Scan(col string) iter.Seq2[int, any] {
	// TODO: Implement scan
	return nil
}

// Filter returns row indices matching the predicate
func (ps *PropertyStore) Filter(pred func(map[string]any) bool) []int {
	// TODO: Implement filter
	// For each row, build map of values and test predicate
	return nil
}

// MemoryUsage returns total memory usage in bytes
func (ps *PropertyStore) MemoryUsage() int64 {
	// TODO: Sum memory usage of all columns
	return 0
}

// RowCount returns the number of rows
func (ps *PropertyStore) RowCount() int {
	return ps.rowCount
}
