# Project 3: Key-Value Store

## Overview
Build an in-memory key-value store with persistence, thread-safety, and a simple query language.

## Concepts Covered
- Maps and hash tables
- File I/O and serialization
- Thread-safe operations with mutexes
- JSON and binary encoding
- Benchmarking and performance testing
- Write-ahead logging (optional)

## Requirements

### Core Functionality
1. **Basic Operations**
   - `GET(key)` - Retrieve value for key
   - `SET(key, value)` - Store key-value pair
   - `DELETE(key)` - Remove key
   - `EXISTS(key)` - Check if key exists
   - `KEYS(pattern)` - List keys matching pattern

2. **Persistence**
   - Snapshot to disk (JSON or binary format)
   - Load from disk on startup
   - Automatic periodic snapshots
   - Manual snapshot command

3. **Thread Safety**
   - Concurrent read/write support
   - Use `sync.RWMutex` for readers-writer lock
   - Atomic operations where appropriate
   - Pass race detector tests

4. **Query Language**
   - Simple command parser
   - Support for basic operations
   - Pattern matching for keys
   - Batch operations (optional)

### Advanced Features (Optional)
1. **Write-Ahead Log (WAL)**
   - Log all mutations before applying
   - Replay log on startup
   - Log compaction/truncation
   - Crash recovery

2. **Data Types**
   - String values
   - Integer values
   - List values
   - Hash values (nested maps)

3. **Time-to-Live (TTL)**
   - Expiration for keys
   - Automatic cleanup
   - TTL refresh on access

## Example Usage

### As a Library
```go
store := kvstore.New()
store.Set("name", "Alice")
val, _ := store.Get("name")
fmt.Println(val) // Output: Alice

// Save to disk
store.Snapshot("data.json")

// Load from disk
store.Load("data.json")
```

### CLI Interface
```bash
# Start the server
./kvstore -file data.json

> SET user:1 Alice
OK

> GET user:1
Alice

> SET user:2 Bob
OK

> KEYS user:*
user:1
user:2

> DELETE user:1
OK

> SNAPSHOT
Saved snapshot to data.json

> EXIT
Goodbye!
```

## Architecture

```
┌─────────────────────────────────┐
│         KVStore                 │
├─────────────────────────────────┤
│  - data: map[string]string      │
│  - mu: sync.RWMutex             │
│  - wal: *WAL (optional)         │
├─────────────────────────────────┤
│  + Get(key) string              │
│  + Set(key, value)              │
│  + Delete(key)                  │
│  + Snapshot(file)               │
│  + Load(file)                   │
└─────────────────────────────────┘
```

## Data Format

### JSON Format (simpler)
```json
{
  "version": 1,
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "user:1": "Alice",
    "user:2": "Bob",
    "counter": "42"
  }
}
```

### Binary Format (more efficient)
```
Magic Number: 0x4B56DB (4 bytes)
Version: 1 (4 bytes)
Timestamp: Unix timestamp (8 bytes)
Entry Count: N (4 bytes)
Entries:
  - Key Length (4 bytes)
  - Key Data (variable)
  - Value Length (4 bytes)
  - Value Data (variable)
```

## Getting Started

1. Initialize the Go module:
   ```bash
   go mod init kvstore
   ```

2. Run the store:
   ```bash
   go run main.go -file data.json
   ```

3. Run tests:
   ```bash
   go test -v
   go test -race
   go test -cover
   ```

4. Run benchmarks:
   ```bash
   go test -bench=. -benchmem
   go test -bench=. -cpuprofile=cpu.prof
   go tool pprof cpu.prof
   ```

## Implementation Hints

### Thread-Safe Map Operations
```go
type Store struct {
    data map[string]string
    mu   sync.RWMutex
}

func (s *Store) Get(key string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    val, ok := s.data[key]
    return val, ok
}

func (s *Store) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
}
```

### JSON Persistence
```go
func (s *Store) SaveJSON(filename string) error {
    s.mu.RLock()
    defer s.mu.RUnlock()

    data, err := json.MarshalIndent(s.data, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(filename, data, 0644)
}
```

### Pattern Matching
```go
func (s *Store) Keys(pattern string) []string {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var keys []string
    for k := range s.data {
        if match, _ := filepath.Match(pattern, k); match {
            keys = append(keys, k)
        }
    }
    return keys
}
```

## Testing Strategy

### Unit Tests
- Test each operation in isolation
- Test concurrent access
- Test error conditions
- Test edge cases (empty store, large values)

### Integration Tests
- Test snapshot and restore
- Test concurrent operations
- Test WAL replay (if implemented)
- Test crash scenarios

### Benchmark Tests
- GET performance
- SET performance
- Concurrent operations
- Large dataset performance
- Snapshot performance

## Performance Goals

- Single-threaded GET: > 1M ops/sec
- Single-threaded SET: > 500K ops/sec
- Concurrent reads: Scale with CPU cores
- Memory usage: Reasonable overhead (<2x data size)
- Snapshot: < 1 second for 100K entries

## Stretch Goals

1. **Transactions**
   - MULTI/EXEC for atomic operations
   - Optimistic locking
   - Rollback support

2. **Replication**
   - Master-slave replication
   - Raft consensus (very advanced)
   - Event streaming

3. **Network Interface**
   - TCP server
   - Redis protocol compatible
   - HTTP REST API

4. **Advanced Data Structures**
   - Sorted sets
   - HyperLogLog
   - Bloom filters

5. **Compression**
   - Compress values
   - Compress snapshots
   - Trade-off CPU for storage

## Common Pitfalls

1. **Race Conditions**
   - Always run tests with `-race`
   - Lock before reading/writing shared state
   - Be careful with compound operations

2. **Deadlocks**
   - Avoid nested locks
   - Use defer for unlocking
   - Consider lock-free alternatives

3. **Memory Leaks**
   - Clean up expired keys
   - Limit memory usage
   - Profile with pprof

4. **Data Corruption**
   - Validate data on load
   - Use checksums
   - Atomic file writes

## Learning Outcomes

After completing this project, you should understand:
- How to build thread-safe data structures
- Persistence and serialization techniques
- Performance benchmarking and optimization
- Write-ahead logging concepts
- Testing concurrent code
- Memory management and profiling

## Time Estimate
8-10 hours for core functionality
4-6 hours for comprehensive tests and benchmarks
6-8 hours for WAL implementation (optional)
10-15 hours for stretch goals (optional)

## Additional Resources
- [Redis Design](https://redis.io/docs/about/)
- [Bolt DB](https://github.com/etcd-io/bbolt) (pure Go KV store)
- [BadgerDB](https://github.com/dgraph-io/badger) (LSM-based store)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
