package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFetchURL(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		wantStatus     int
		wantTitle      string
		wantError      bool
	}{
		{
			name: "successful fetch",
			serverResponse: `<html><head><title>Test Page</title></head><body>Content</body></html>`,
			serverStatus:   http.StatusOK,
			wantStatus:     http.StatusOK,
			wantTitle:      "Test Page",
			wantError:      false,
		},
		{
			name:           "404 error",
			serverResponse: "Not Found",
			serverStatus:   http.StatusNotFound,
			wantStatus:     http.StatusNotFound,
			wantError:      false,
		},
		// TODO: Add more test cases
		// - Timeout scenario
		// - Context cancellation
		// - Invalid HTML
		// - Large response
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Fetch URL
			client := &http.Client{Timeout: 5 * time.Second}
			result := fetchURL(context.Background(), client, server.URL)

			// Verify results
			if result.StatusCode != tt.wantStatus {
				t.Errorf("fetchURL() status = %d, want %d", result.StatusCode, tt.wantStatus)
			}

			if result.Title != tt.wantTitle {
				t.Errorf("fetchURL() title = %q, want %q", result.Title, tt.wantTitle)
			}

			if (result.Error != "") != tt.wantError {
				t.Errorf("fetchURL() error = %v, wantError %v", result.Error, tt.wantError)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "simple title",
			html: "<html><head><title>Hello World</title></head></html>",
			want: "Hello World",
		},
		{
			name: "no title",
			html: "<html><head></head></html>",
			want: "",
		},
		{
			name: "title with whitespace",
			html: "<html><head><title>  Spaced Title  </title></head></html>",
			want: "Spaced Title",
		},
		// TODO: Add more test cases
		// - Multiple title tags
		// - Malformed HTML
		// - Special characters in title
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTitle(tt.html)
			if got != tt.want {
				t.Errorf("extractTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadURLs(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantURLs int
	}{
		{
			name: "multiple urls",
			content: `https://example.com
https://golang.org
https://github.com`,
			wantURLs: 3,
		},
		{
			name: "with comments",
			content: `# This is a comment
https://example.com
# Another comment
https://golang.org`,
			wantURLs: 2,
		},
		{
			name: "with empty lines",
			content: `https://example.com

https://golang.org

`,
			wantURLs: 2,
		},
		// TODO: Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile := filepath.Join(t.TempDir(), "urls.txt")
			if err := os.WriteFile(tmpfile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			urls, err := readURLs(tmpfile)
			if err != nil {
				t.Errorf("readURLs() error = %v", err)
				return
			}

			if len(urls) != tt.wantURLs {
				t.Errorf("readURLs() got %d URLs, want %d", len(urls), tt.wantURLs)
			}
		})
	}
}

func TestWorkerPool(t *testing.T) {
	// Create test server
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<title>Test</title>"))
	}))
	defer server.Close()

	// Create URLs file
	tmpfile := filepath.Join(t.TempDir(), "urls.txt")
	urls := ""
	for i := 0; i < 10; i++ {
		urls += server.URL + "\n"
	}
	if err := os.WriteFile(tmpfile, []byte(urls), 0644); err != nil {
		t.Fatal(err)
	}

	// Test scraping
	config := &Config{
		URLsFile:   tmpfile,
		Workers:    3,
		Delay:      10 * time.Millisecond,
		Timeout:    5 * time.Second,
		OutputFile: filepath.Join(t.TempDir(), "results.json"),
	}

	results, err := scrape(config)
	if err != nil {
		t.Errorf("scrape() error = %v", err)
		return
	}

	if len(results) != 10 {
		t.Errorf("scrape() got %d results, want 10", len(results))
	}

	// Verify all requests were successful
	for _, r := range results {
		if r.Error != "" {
			t.Errorf("result has error: %v", r.Error)
		}
	}
}

func TestContextCancellation(t *testing.T) {
	// Create slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client := &http.Client{}
	result := fetchURL(ctx, client, server.URL)

	// Should have timeout error
	if result.Error == "" {
		t.Error("expected timeout error, got none")
	}
}

// Benchmarks
func BenchmarkFetchURL(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<title>Benchmark</title>"))
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fetchURL(ctx, client, server.URL)
	}
}

func BenchmarkExtractTitle(b *testing.B) {
	html := `<html><head><title>Test Page Title</title></head><body>Content here</body></html>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractTitle(html)
	}
}

// TODO: Add more benchmarks
// - Benchmark worker pool with different sizes
// - Benchmark with rate limiting
// - Benchmark with different number of URLs
