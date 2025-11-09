package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadCSV(t *testing.T) {
	tests := []struct {
		name        string
		csvContent  string
		wantRecords int
		wantHeaders []string
		wantErr     bool
	}{
		{
			name: "valid csv",
			csvContent: `name,age,city
Alice,30,NYC
Bob,25,LA`,
			wantRecords: 2,
			wantHeaders: []string{"name", "age", "city"},
			wantErr:     false,
		},
		{
			name: "empty csv",
			csvContent: `name,age,city`,
			wantRecords: 0,
			wantHeaders: []string{"name", "age", "city"},
			wantErr:     false,
		},
		// TODO: Add more test cases
		// - CSV with missing values
		// - CSV with extra columns
		// - Malformed CSV
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpfile := createTempCSV(t, tt.csvContent)
			defer os.Remove(tmpfile)

			records, headers, err := readCSV(tmpfile)

			if (err != nil) != tt.wantErr {
				t.Errorf("readCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(records) != tt.wantRecords {
				t.Errorf("readCSV() got %d records, want %d", len(records), tt.wantRecords)
			}

			if len(headers) != len(tt.wantHeaders) {
				t.Errorf("readCSV() got %d headers, want %d", len(headers), len(tt.wantHeaders))
			}
		})
	}
}

func TestFilterRecords(t *testing.T) {
	records := []Record{
		{"name": "Alice", "age": "30", "city": "NYC"},
		{"name": "Bob", "age": "25", "city": "LA"},
		{"name": "Charlie", "age": "30", "city": "SF"},
	}

	tests := []struct {
		name      string
		column    string
		value     string
		wantCount int
	}{
		{
			name:      "filter by age",
			column:    "age",
			value:     "30",
			wantCount: 2,
		},
		{
			name:      "filter by city",
			column:    "city",
			value:     "LA",
			wantCount: 1,
		},
		{
			name:      "no matches",
			column:    "city",
			value:     "Boston",
			wantCount: 0,
		},
		// TODO: Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, err := filterRecords(records, tt.column, tt.value)
			if err != nil {
				t.Errorf("filterRecords() error = %v", err)
				return
			}

			if len(filtered) != tt.wantCount {
				t.Errorf("filterRecords() got %d records, want %d", len(filtered), tt.wantCount)
			}
		})
	}
}

func TestAggregateRecords(t *testing.T) {
	records := []Record{
		{"amount": "100"},
		{"amount": "200"},
		{"amount": "300"},
	}

	tests := []struct {
		name      string
		column    string
		operation string
		want      float64
		wantErr   bool
	}{
		{"sum", "amount", "sum", 600.0, false},
		{"avg", "amount", "avg", 200.0, false},
		{"min", "amount", "min", 100.0, false},
		{"max", "amount", "max", 300.0, false},
		{"count", "amount", "count", 3.0, false},
		{"invalid op", "amount", "invalid", 0.0, true},
		// TODO: Add more test cases
		// - Empty records
		// - Non-numeric values
		// - Mixed numeric/non-numeric
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aggregateRecords(records, tt.column, tt.operation)

			if (err != nil) != tt.wantErr {
				t.Errorf("aggregateRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.want {
				t.Errorf("aggregateRecords() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestWriteCSV(t *testing.T) {
	records := []Record{
		{"name": "Alice", "age": "30"},
		{"name": "Bob", "age": "25"},
	}
	headers := []string{"name", "age"}

	tmpfile := filepath.Join(t.TempDir(), "output.csv")

	err := writeCSV(tmpfile, records, headers, "Summary: 2 rows")
	if err != nil {
		t.Errorf("writeCSV() error = %v", err)
		return
	}

	// Verify file was created and has content
	if _, err := os.Stat(tmpfile); os.IsNotExist(err) {
		t.Error("writeCSV() did not create output file")
	}

	// TODO: Read back and verify content
}

// Helper function to create temporary CSV file
func createTempCSV(t *testing.T, content string) string {
	t.Helper()

	tmpfile := filepath.Join(t.TempDir(), "test.csv")
	if err := os.WriteFile(tmpfile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	return tmpfile
}

// TODO: Add benchmarks
func BenchmarkReadCSV(b *testing.B) {
	// Create test file with many rows
	content := "name,age,city\n"
	for i := 0; i < 1000; i++ {
		content += "Alice,30,NYC\n"
	}

	tmpfile := filepath.Join(b.TempDir(), "bench.csv")
	if err := os.WriteFile(tmpfile, []byte(content), 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := readCSV(tmpfile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFilterRecords(b *testing.B) {
	// TODO: Implement benchmark
	// Create large record set
	// Benchmark filtering operations
}

func BenchmarkAggregateRecords(b *testing.B) {
	// TODO: Implement benchmark
	// Create large record set
	// Benchmark different aggregation operations
}
