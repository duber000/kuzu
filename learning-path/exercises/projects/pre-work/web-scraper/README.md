# Project 2: Concurrent Web Scraper

## Overview
Build a concurrent web scraper that fetches and processes multiple URLs in parallel using Go's concurrency primitives.

## Concepts Covered
- Goroutines for concurrent execution
- Channels for communication
- sync.WaitGroup for synchronization
- Worker pool pattern
- Context for cancellation
- HTTP client usage
- Rate limiting

## Requirements

### Core Functionality
1. **URL Management**
   - Accept list of URLs from file or CLI
   - Validate URLs before processing
   - Handle HTTP/HTTPS schemes

2. **Concurrent Fetching**
   - Use worker pool pattern
   - Configurable number of workers
   - Distribute work via channels
   - Collect results concurrently

3. **Rate Limiting**
   - Respect robots.txt (optional)
   - Implement delay between requests
   - Per-domain rate limiting
   - Global rate limiting

4. **Result Aggregation**
   - Collect response status codes
   - Extract page titles
   - Count word frequencies (optional)
   - Save results to JSON file

5. **Error Handling**
   - Handle network errors gracefully
   - Retry failed requests (optional)
   - Timeout configuration
   - Context cancellation support

### Performance Requirements
- Process at least 100 URLs concurrently
- Configurable worker pool size
- Efficient memory usage
- Graceful shutdown

### Testing Requirements
- Unit tests for core components
- Integration tests with test server
- Mock HTTP responses
- Test timeout and cancellation
- 80%+ code coverage

## Example Usage

```bash
# Basic scraping with 5 workers
./scraper -urls urls.txt -workers 5 -output results.json

# With rate limiting
./scraper -urls urls.txt -workers 10 -delay 100ms -output results.json

# With timeout
./scraper -urls urls.txt -workers 5 -timeout 30s -output results.json
```

## Architecture

```
┌──────────┐
│  Main    │
│  Thread  │
└────┬─────┘
     │
     ├──> URL Channel ──> Worker 1 ──┐
     ├──> URL Channel ──> Worker 2 ──┤
     ├──> URL Channel ──> Worker 3 ──┼──> Results Channel ──> Aggregator
     ├──> URL Channel ──> Worker 4 ──┤
     └──> URL Channel ──> Worker 5 ──┘
```

## Sample URLs File

```
https://example.com
https://golang.org
https://github.com
https://stackoverflow.com
https://news.ycombinator.com
```

## Sample Output (results.json)

```json
{
  "total_urls": 5,
  "successful": 4,
  "failed": 1,
  "duration_seconds": 2.5,
  "results": [
    {
      "url": "https://example.com",
      "status_code": 200,
      "title": "Example Domain",
      "size_bytes": 1256,
      "duration_ms": 145
    },
    {
      "url": "https://invalid.example",
      "error": "no such host"
    }
  ]
}
```

## Getting Started

1. Initialize the Go module:
   ```bash
   go mod init webscraper
   ```

2. Run the scraper:
   ```bash
   go run main.go -urls urls.txt -workers 5
   ```

3. Run tests:
   ```bash
   go test -v
   go test -race
   ```

4. Run benchmarks:
   ```bash
   go test -bench=. -benchmem
   ```

## Implementation Hints

### Worker Pool Pattern
```go
func worker(id int, jobs <-chan string, results chan<- Result) {
    for url := range jobs {
        // Fetch and process URL
        result := fetch(url)
        results <- result
    }
}
```

### Using Context for Cancellation
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := client.Do(req)
```

### Rate Limiting
```go
limiter := time.NewTicker(100 * time.Millisecond)
defer limiter.Stop()

for url := range urls {
    <-limiter.C  // Wait for rate limit
    go fetch(url)
}
```

## Stretch Goals

1. **Advanced Features**
   - Extract links from pages
   - Download images
   - Depth-limited crawling
   - Respect robots.txt

2. **Performance Optimizations**
   - HTTP connection pooling
   - Compression support
   - Resume capability
   - Progress tracking

3. **Monitoring**
   - Real-time progress display
   - Prometheus metrics
   - Request/response logging
   - Performance statistics

4. **Persistence**
   - Save intermediate results
   - Resume from checkpoint
   - Cache responses
   - Deduplication

## Common Pitfalls

1. **Goroutine Leaks**
   - Always close channels when done
   - Use WaitGroup correctly
   - Handle context cancellation

2. **Race Conditions**
   - Protect shared state with mutexes
   - Use atomic operations when appropriate
   - Run tests with `-race` flag

3. **Resource Exhaustion**
   - Limit number of concurrent connections
   - Set appropriate timeouts
   - Close HTTP response bodies

4. **Error Handling**
   - Don't ignore errors from channels
   - Handle network timeouts
   - Validate input URLs

## Learning Outcomes

After completing this project, you should understand:
- How to use goroutines and channels effectively
- Worker pool pattern implementation
- Context for cancellation and timeouts
- Rate limiting techniques
- Concurrent error handling
- Testing concurrent code

## Time Estimate
8-10 hours for core functionality
3-4 hours for comprehensive tests
2-3 hours for stretch goals (optional)

## Additional Resources
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Context Package](https://pkg.go.dev/context)
- [HTTP Client Best Practices](https://pkg.go.dev/net/http)
