package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Config holds the CLI configuration
type Config struct {
	InputFile  string
	OutputFile string
	Filter     string
	Value      string
	Aggregate  string
	Operation  string
}

// Record represents a CSV row
type Record map[string]string

func main() {
	config := parseFlags()
	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags parses and validates command-line flags
func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.InputFile, "input", "", "Input CSV file path (required)")
	flag.StringVar(&config.InputFile, "i", "", "Input CSV file path (shorthand)")
	flag.StringVar(&config.OutputFile, "output", "", "Output file path (required)")
	flag.StringVar(&config.OutputFile, "o", "", "Output file path (shorthand)")
	flag.StringVar(&config.Filter, "filter", "", "Column to filter on")
	flag.StringVar(&config.Value, "value", "", "Value to filter for")
	flag.StringVar(&config.Aggregate, "aggregate", "", "Column to aggregate")
	flag.StringVar(&config.Operation, "operation", "count", "Aggregation operation (sum, avg, count, min, max)")

	flag.Parse()

	// Validate required flags
	if config.InputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: -input flag is required")
		flag.Usage()
		os.Exit(1)
	}
	if config.OutputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: -output flag is required")
		flag.Usage()
		os.Exit(1)
	}

	return config
}

// run executes the main program logic
func run(config *Config) error {
	// TODO: Implement the main logic
	// 1. Read CSV file
	// 2. Apply filters if specified
	// 3. Perform aggregation if specified
	// 4. Write results to output file

	records, headers, err := readCSV(config.InputFile)
	if err != nil {
		return fmt.Errorf("reading CSV: %w", err)
	}

	// Filter records if filter is specified
	if config.Filter != "" && config.Value != "" {
		records, err = filterRecords(records, config.Filter, config.Value)
		if err != nil {
			return fmt.Errorf("filtering records: %w", err)
		}
	}

	// Perform aggregation if specified
	var summary string
	if config.Aggregate != "" {
		result, err := aggregateRecords(records, config.Aggregate, config.Operation)
		if err != nil {
			return fmt.Errorf("aggregating records: %w", err)
		}
		summary = fmt.Sprintf("Summary: %d rows, %s %s: %.2f",
			len(records), config.Operation, config.Aggregate, result)
	}

	// Write results
	if err := writeCSV(config.OutputFile, records, headers, summary); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	if summary != "" {
		fmt.Println(summary)
	}

	return nil
}

// readCSV reads a CSV file and returns records with headers
func readCSV(filename string) ([]Record, []string, error) {
	// TODO: Implement CSV reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("reading headers: %w", err)
	}

	// Read all records
	var records []Record
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("reading row: %w", err)
		}

		record := make(Record)
		for i, header := range headers {
			if i < len(row) {
				record[header] = row[i]
			}
		}
		records = append(records, record)
	}

	return records, headers, nil
}

// filterRecords filters records based on column and value
func filterRecords(records []Record, column, value string) ([]Record, error) {
	// TODO: Implement filtering logic
	var filtered []Record
	for _, record := range records {
		if record[column] == value {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

// aggregateRecords performs aggregation on a column
func aggregateRecords(records []Record, column, operation string) (float64, error) {
	// TODO: Implement aggregation logic
	if len(records) == 0 {
		return 0, nil
	}

	var values []float64
	for _, record := range records {
		val, err := strconv.ParseFloat(strings.TrimSpace(record[column]), 64)
		if err != nil {
			// Skip non-numeric values
			continue
		}
		values = append(values, val)
	}

	if len(values) == 0 {
		return 0, fmt.Errorf("no numeric values found in column %s", column)
	}

	switch operation {
	case "count":
		return float64(len(values)), nil
	case "sum":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum, nil
	case "avg":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum / float64(len(values)), nil
	case "min":
		min := values[0]
		for _, v := range values[1:] {
			if v < min {
				min = v
			}
		}
		return min, nil
	case "max":
		max := values[0]
		for _, v := range values[1:] {
			if v > max {
				max = v
			}
		}
		return max, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", operation)
	}
}

// writeCSV writes records to a CSV file
func writeCSV(filename string, records []Record, headers []string, summary string) error {
	// TODO: Implement CSV writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("writing headers: %w", err)
	}

	// Write records
	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = record[header]
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("writing row: %w", err)
		}
	}

	// Write summary as comment if present
	if summary != "" {
		// Note: CSV doesn't have comments, so we write it as a row with # prefix
		if err := writer.Write([]string{summary}); err != nil {
			return fmt.Errorf("writing summary: %w", err)
		}
	}

	return nil
}
