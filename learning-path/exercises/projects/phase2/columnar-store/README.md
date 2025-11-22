# Project 2.2: Columnar Property Store

## Overview
Implement a memory-efficient columnar property storage system using Go 1.23's `unique` package for string interning, compression techniques, and efficient scanning.

**Duration:** 12-15 hours
**Difficulty:** Medium

## Learning Objectives
- Understand columnar vs row-oriented storage
- Master Go 1.23 `unique.Handle` for string interning
- Implement compression techniques (dictionary encoding, bit packing)
- Build efficient scan and filter operations
- Profile memory usage and optimize
- Compare columnar vs row storage

## Concepts Covered
- Columnar data layout
- String interning with `unique.Handle`
- Dictionary encoding
- Bit packing for integers
- NULL bitmap representation
- Vectorized scanning
- Compression algorithms

## Requirements

### Core Functionality

#### 1. Column Types
- **IntColumn** - Integer values with bit packing
- **StringColumn** - Strings with dictionary encoding
- **BoolColumn** - Booleans with bitmap
- **FloatColumn** - Float64 values
- Support NULL values with bitmap

#### 2. Property Store
- Store multiple columns per entity type
- Add/remove columns dynamically
- Efficient row-wise and column-wise access
- Memory-efficient storage

#### 3. Compression Techniques
- Dictionary encoding for strings
- Bit packing for small integers
- NULL bitmap (1 bit per value)
- Run-length encoding (optional)

#### 4. Query Operations
- **Scan** - Iterate all values
- **Filter** - Select rows by predicate
- **Project** - Select specific columns
- **Aggregate** - Sum, count, avg, etc.

## API Design

```go
package columnarstore

import "unique"

// Column interface for different column types
type Column interface {
	Append(value any) error
	Get(index int) (any, bool)  // (value, isNull)
	Scan() iter.Seq2[int, any]
	MemoryUsage() int64
	RowCount() int
}

// IntColumn stores integers with bit packing
type IntColumn struct {
	values   []byte       // bit-packed values
	nulls    *Bitmap      // NULL bitmap
	bitWidth int          // bits per value
	rowCount int
}

// StringColumn stores strings with dictionary encoding
type StringColumn struct {
	dict    []unique.Handle[string]  // dictionary
	indices []uint32                 // indices into dictionary
	nulls   *Bitmap
}

// PropertyStore stores columns for entities
type PropertyStore struct {
	columns map[string]Column
	rowCount int
}

// Create new property store
func NewPropertyStore() *PropertyStore

// Add column
func (ps *PropertyStore) AddColumn(name string, col Column) error

// Append row
func (ps *PropertyStore) AppendRow(values map[string]any) error

// Get value
func (ps *PropertyStore) Get(row int, col string) (any, bool, error)

// Scan column
func (ps *PropertyStore) Scan(col string) iter.Seq2[int, any]

// Filter rows
func (ps *PropertyStore) Filter(pred func(map[string]any) bool) []int

// Memory usage
func (ps *PropertyStore) MemoryUsage() int64
```

## Implementation Hints

### String Interning with unique.Handle
```go
import "unique"

type StringColumn struct {
	dict     []unique.Handle[string]
	indices  []uint32
	nulls    *Bitmap
	dictMap  map[unique.Handle[string]]uint32
	rowCount int
}

func (c *StringColumn) Append(value any) error {
	if value == nil {
		c.nulls.Set(c.rowCount)
		c.indices = append(c.indices, 0)
		c.rowCount++
		return nil
	}

	str := value.(string)
	handle := unique.Make(str)

	// Check if already in dictionary
	if idx, found := c.dictMap[handle]; found {
		c.indices = append(c.indices, idx)
	} else {
		// Add to dictionary
		idx := uint32(len(c.dict))
		c.dict = append(c.dict, handle)
		c.dictMap[handle] = idx
		c.indices = append(c.indices, idx)
	}

	c.rowCount++
	return nil
}

func (c *StringColumn) Get(index int) (any, bool) {
	if c.nulls.Test(index) {
		return nil, true
	}

	dictIdx := c.indices[index]
	return c.dict[dictIdx].Value(), false
}
```

### Bit Packing for Integers
```go
type IntColumn struct {
	values   []byte
	nulls    *Bitmap
	bitWidth int
	minValue int64
	rowCount int
}

func (c *IntColumn) packValue(value int64, index int) {
	// Store value - minValue with bitWidth bits
	normalized := uint64(value - c.minValue)

	bitOffset := index * c.bitWidth
	byteOffset := bitOffset / 8
	bitPos := bitOffset % 8

	// Write bits (handle multi-byte case)
	for i := 0; i < c.bitWidth; i++ {
		if normalized&(1<<i) != 0 {
			c.values[byteOffset] |= (1 << bitPos)
		}
		bitPos++
		if bitPos == 8 {
			bitPos = 0
			byteOffset++
		}
	}
}

func (c *IntColumn) unpackValue(index int) int64 {
	bitOffset := index * c.bitWidth
	byteOffset := bitOffset / 8
	bitPos := bitOffset % 8

	var value uint64
	for i := 0; i < c.bitWidth; i++ {
		if c.values[byteOffset]&(1<<bitPos) != 0 {
			value |= (1 << i)
		}
		bitPos++
		if bitPos == 8 {
			bitPos = 0
			byteOffset++
		}
	}

	return int64(value) + c.minValue
}
```

### NULL Bitmap
```go
type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {
	return &Bitmap{
		bits: make([]byte, (size+7)/8),
		size: size,
	}
}

func (b *Bitmap) Set(pos int) {
	b.bits[pos/8] |= (1 << (pos % 8))
}

func (b *Bitmap) Clear(pos int) {
	b.bits[pos/8] &= ^(1 << (pos % 8))
}

func (b *Bitmap) Test(pos int) bool {
	return (b.bits[pos/8] & (1 << (pos % 8))) != 0
}

func (b *Bitmap) CountOnes() int {
	// Use popcount for efficiency
	count := 0
	for _, b := range b.bits {
		count += bits.OnesCount8(b)
	}
	return count
}
```

## Performance Goals

- String memory: < 50% of row format (via deduplication)
- Integer memory: < 25% of row format (via bit packing)
- Scan speed: > 1GB/s
- Filter speed: > 500MB/s
- Append speed: > 1M rows/sec
- Memory overhead: < 10% for metadata

## Stretch Goals

### 1. Run-Length Encoding
- Compress repeated values
- Mixed encoding (RLE + dictionary)
- Adaptive encoding selection

### 2. Bloom Filters
- Per-column bloom filters
- Fast membership testing
- Estimate cardinality

### 3. Column Statistics
- Min/max values
- Distinct count (HyperLogLog)
- Histogram
- Use for query optimization

### 4. Vectorized Operations
- SIMD operations
- Batch processing
- Predicate pushdown

## Validation Checklist

- [ ] Pass all unit tests
- [ ] Support NULL values
- [ ] Implement all column types
- [ ] Use unique.Handle for strings
- [ ] Achieve >50% memory savings
- [ ] Scan at >500MB/s
- [ ] No memory leaks

## Time Estimate
- Core implementation: 8-10 hours
- Compression optimization: 2-3 hours
- Testing and profiling: 2-3 hours
