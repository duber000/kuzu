# Project 1: CLI Tool with Flags

## Overview
Build a command-line tool that processes CSV data with filtering and aggregation capabilities.

## Concepts Covered
- `flag` package for CLI argument parsing
- File I/O operations
- Error handling with wrapped errors
- Table-driven tests
- Benchmark testing

## Requirements

### Core Functionality
1. **CLI Arguments** - Accept the following flags:
   - `-input` or `-i`: Input CSV file path (required)
   - `-output` or `-o`: Output file path (required)
   - `-filter`: Column to filter on (optional)
   - `-value`: Value to filter for (optional)
   - `-aggregate`: Column to aggregate (optional)
   - `-operation`: Aggregation operation (sum, avg, count, min, max)

2. **CSV Reading**
   - Read and parse CSV data
   - Handle headers correctly
   - Support both comma and tab-delimited files

3. **Filtering**
   - Filter rows based on column value
   - Support exact match and pattern matching
   - Handle missing values gracefully

4. **Aggregation**
   - Support numeric aggregations (sum, avg, min, max)
   - Support count operation
   - Group by support (optional)

5. **Output**
   - Write results to specified output file
   - Maintain CSV format
   - Include summary statistics

### Error Handling
- Proper error wrapping with context
- Validate file paths before processing
- Handle malformed CSV data
- Clear error messages for users

### Testing Requirements
- Table-driven tests for each function
- Test edge cases (empty files, malformed data)
- Test error conditions
- Achieve 80%+ code coverage

### Benchmark Requirements
- Benchmark I/O performance
- Benchmark filtering operations
- Benchmark aggregation operations
- Compare performance with different file sizes

## Example Usage

```bash
# Basic filtering
./csvtool -i data.csv -o results.csv -filter "age" -value "25"

# Aggregation
./csvtool -i sales.csv -o summary.csv -aggregate "amount" -operation "sum"

# Combined filtering and aggregation
./csvtool -i sales.csv -o results.csv -filter "region" -value "west" -aggregate "amount" -operation "avg"
```

## Sample Data Format

**Input (sales.csv):**
```csv
date,region,product,amount
2024-01-01,west,widget,100
2024-01-02,east,gadget,150
2024-01-03,west,widget,200
2024-01-04,south,tool,75
```

**Output (filtered by region=west):**
```csv
date,region,product,amount
2024-01-01,west,widget,100
2024-01-03,west,widget,200
Summary: 2 rows, avg amount: 150
```

## Getting Started

1. Initialize the Go module:
   ```bash
   go mod init csvtool
   ```

2. Run the starter code:
   ```bash
   go run main.go -i sample.csv -o output.csv
   ```

3. Run tests:
   ```bash
   go test -v
   go test -cover
   ```

4. Run benchmarks:
   ```bash
   go test -bench=. -benchmem
   ```

## Implementation Hints

1. **Flag Package**
   - Use `flag.String()` for string arguments
   - Use `flag.Parse()` after defining flags
   - Check for required flags after parsing

2. **CSV Reading**
   - Use `encoding/csv` package
   - Handle headers separately from data rows
   - Use `csv.Reader` for efficient reading

3. **Error Handling**
   - Use `fmt.Errorf()` with `%w` for wrapping
   - Return errors up the call stack
   - Log errors at the appropriate level

4. **Testing**
   - Create test fixtures with sample CSV data
   - Use `t.Helper()` for test helper functions
   - Test both success and failure cases

## Stretch Goals

1. **Advanced Filtering**
   - Support multiple filter conditions (AND/OR)
   - Regular expression matching
   - Numeric comparisons (<, >, <=, >=)

2. **Performance Optimizations**
   - Stream processing for large files
   - Parallel CSV parsing
   - Memory-mapped file I/O

3. **Additional Features**
   - Support for JSON output
   - Column projection (select specific columns)
   - Sorting capabilities
   - Progress bar for large files

## Learning Outcomes

After completing this project, you should understand:
- How to build robust CLI tools in Go
- Proper error handling patterns
- Table-driven testing methodology
- Performance benchmarking techniques
- File I/O best practices

## Time Estimate
6-8 hours for core functionality
2-4 hours for comprehensive tests and benchmarks
2-3 hours for stretch goals (optional)
