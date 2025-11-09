package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Config holds the scraper configuration
type Config struct {
	URLsFile   string
	Workers    int
	Delay      time.Duration
	Timeout    time.Duration
	OutputFile string
}

// Result represents a scraping result
type Result struct {
	URL        string  `json:"url"`
	StatusCode int     `json:"status_code,omitempty"`
	Title      string  `json:"title,omitempty"`
	SizeBytes  int     `json:"size_bytes,omitempty"`
	DurationMS int64   `json:"duration_ms,omitempty"`
	Error      string  `json:"error,omitempty"`
}

// Summary represents the overall scraping summary
type Summary struct {
	TotalURLs       int       `json:"total_urls"`
	Successful      int       `json:"successful"`
	Failed          int       `json:"failed"`
	DurationSeconds float64   `json:"duration_seconds"`
	Results         []Result  `json:"results"`
}

func main() {
	config := parseFlags()

	startTime := time.Now()
	results, err := scrape(config)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	summary := createSummary(results, time.Since(startTime))
	if err := writeSummary(config.OutputFile, summary); err != nil {
		log.Fatalf("Error writing summary: %v", err)
	}

	fmt.Printf("Scraped %d URLs in %.2f seconds\n", summary.TotalURLs, summary.DurationSeconds)
	fmt.Printf("Successful: %d, Failed: %d\n", summary.Successful, summary.Failed)
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.URLsFile, "urls", "", "File containing URLs to scrape (required)")
	flag.IntVar(&config.Workers, "workers", 5, "Number of worker goroutines")
	flag.DurationVar(&config.Delay, "delay", 0, "Delay between requests per worker")
	flag.DurationVar(&config.Timeout, "timeout", 30*time.Second, "HTTP request timeout")
	flag.StringVar(&config.OutputFile, "output", "results.json", "Output JSON file")

	flag.Parse()

	if config.URLsFile == "" {
		fmt.Fprintln(os.Stderr, "Error: -urls flag is required")
		flag.Usage()
		os.Exit(1)
	}

	return config
}

// scrape orchestrates the concurrent scraping process
func scrape(config *Config) ([]Result, error) {
	// Read URLs from file
	urls, err := readURLs(config.URLsFile)
	if err != nil {
		return nil, fmt.Errorf("reading URLs: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout*time.Duration(len(urls)))
	defer cancel()

	// Create channels
	jobs := make(chan string, len(urls))
	results := make(chan Result, len(urls))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, results, config.Delay, config.Timeout, &wg)
	}

	// Send jobs
	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []Result
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults, nil
}

// worker processes URLs from the jobs channel
func worker(ctx context.Context, id int, jobs <-chan string, results chan<- Result, delay time.Duration, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: timeout,
	}

	for url := range jobs {
		// Apply delay for rate limiting
		if delay > 0 {
			time.Sleep(delay)
		}

		result := fetchURL(ctx, client, url)
		results <- result
	}
}

// fetchURL fetches a single URL and extracts information
func fetchURL(ctx context.Context, client *http.Client, url string) Result {
	startTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{URL: url, Error: err.Error()}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{URL: url, Error: err.Error()}
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{
			URL:        url,
			StatusCode: resp.StatusCode,
			Error:      fmt.Sprintf("reading body: %v", err),
		}
	}

	// Extract title
	title := extractTitle(string(body))

	return Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Title:      title,
		SizeBytes:  len(body),
		DurationMS: time.Since(startTime).Milliseconds(),
	}
}

// extractTitle extracts the page title from HTML
func extractTitle(html string) string {
	// TODO: Implement proper HTML parsing
	// This is a simple implementation - use html.Parse for production
	start := strings.Index(html, "<title>")
	if start == -1 {
		return ""
	}
	start += 7 // len("<title>")

	end := strings.Index(html[start:], "</title>")
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(html[start : start+end])
}

// readURLs reads URLs from a file
func readURLs(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var urls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	return urls, nil
}

// createSummary creates a summary from results
func createSummary(results []Result, duration time.Duration) *Summary {
	summary := &Summary{
		TotalURLs:       len(results),
		DurationSeconds: duration.Seconds(),
		Results:         results,
	}

	for _, r := range results {
		if r.Error == "" {
			summary.Successful++
		} else {
			summary.Failed++
		}
	}

	return summary
}

// writeSummary writes the summary to a JSON file
func writeSummary(filename string, summary *Summary) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
