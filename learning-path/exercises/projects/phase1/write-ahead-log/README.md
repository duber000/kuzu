# Project 1.3: Write-Ahead Log

## Overview
Implement a crash-safe Write-Ahead Log (WAL) system with append-only logging, group commit optimization, and recovery mechanisms.

**Duration:** 10-12 hours
**Difficulty:** Medium-Hard

## Learning Objectives
- Understand WAL principles for crash recovery
- Implement append-only file operations
- Master fsync and durability guarantees
- Build recovery and replay mechanisms
- Optimize with group commits
- Implement log truncation and archival

## Concepts Covered
- Write-Ahead Logging (WAL)
- Append-only data structures
- File synchronization (fsync)
- Recovery and replay
- Log sequence numbers (LSN)
- Group commit optimization
- Idempotent operations
- Checkpointing

## Requirements

### Core Functionality

#### 1. WAL System
- Append-only log file
- Atomic log record writes
- fsync after write (configurable)
- Log sequence numbers (LSN)
- Recovery from log
- Log truncation after checkpoint

#### 2. Log Record Structure
```
┌────────────────────────────────┐
│  Log Record Header             │
│  - LSN (8 bytes)               │
│  - Record Type (1 byte)        │
│  - Transaction ID (8 bytes)    │
│  - Record Length (4 bytes)     │
│  - Checksum (4 bytes)          │
├────────────────────────────────┤
│  Record Data (variable)        │
│  - Operation-specific payload  │
│                                │
└────────────────────────────────┘
```

#### 3. Record Types
- **Begin** - Transaction start
- **Commit** - Transaction commit
- **Abort** - Transaction rollback
- **Update** - Data modification (before/after images)
- **Checkpoint** - Recovery point

#### 4. Log Operations
- **Append(record)** - Write log record
- **Flush()** - fsync log to disk
- **Recover()** - Replay log after crash
- **Checkpoint()** - Create recovery point
- **Truncate(LSN)** - Remove old log entries

#### 5. Group Commit
- Buffer log records in memory
- Flush multiple records together
- Reduce fsync overhead
- Configurable flush interval

## Getting Started

```bash
# Initialize module
cd write-ahead-log
go mod init wal

# Run tests
go test -v
go test -race
go test -cover

# Run benchmarks
go test -bench=. -benchmem

# Chaos testing
go test -run=TestCrashRecovery -count=100
```

## API Design

```go
package wal

import "io"

type LSN uint64
type TxnID uint64

type RecordType byte

const (
	RecordBegin RecordType = iota
	RecordCommit
	RecordAbort
	RecordUpdate
	RecordCheckpoint
)

// LogRecord represents a WAL record
type LogRecord struct {
	LSN      LSN
	Type     RecordType
	TxnID    TxnID
	Data     []byte
	Checksum uint32
}

// WAL is the write-ahead log
type WAL struct {
	file        *os.File
	currentLSN  atomic.Uint64
	flushLSN    atomic.Uint64
	buffer      *LogBuffer
	flusher     *GroupCommitFlusher
}

// Options for WAL configuration
type WALOptions struct {
	FilePath      string
	BufferSize    int
	FlushInterval time.Duration
	SyncOnCommit  bool
}

// Create new WAL
func New(opts WALOptions) (*WAL, error)

// Append a log record
func (w *WAL) Append(record *LogRecord) (LSN, error)

// Flush all buffered records to disk
func (w *WAL) Flush() error

// Recover from log file
func (w *WAL) Recover(handler RecoveryHandler) error

// Create checkpoint
func (w *WAL) Checkpoint() (LSN, error)

// Truncate log up to LSN
func (w *WAL) Truncate(lsn LSN) error

// Close WAL
func (w *WAL) Close() error

// RecoveryHandler is called during recovery for each record
type RecoveryHandler interface {
	OnBegin(txnID TxnID, lsn LSN) error
	OnCommit(txnID TxnID, lsn LSN) error
	OnAbort(txnID TxnID, lsn LSN) error
	OnUpdate(txnID TxnID, lsn LSN, data []byte) error
	OnCheckpoint(lsn LSN) error
}
```

## Test Cases

### Correctness Tests
- **TestAppend** - Append log records
- **TestSequentialLSN** - LSNs are monotonically increasing
- **TestFlush** - Flush writes to disk
- **TestRecovery** - Recover from log file
- **TestCheckpoint** - Create and recover from checkpoint
- **TestTruncate** - Truncate old log entries
- **TestChecksum** - Detect corrupted records

### Crash Simulation Tests
- **TestCrashDuringWrite** - Partial record writes
- **TestCrashBeforeFlush** - Lost uncommitted data
- **TestCrashAfterCommit** - Committed data preserved
- **TestMultipleCrashes** - Repeated crash recovery
- **TestCorruptedLog** - Handle corrupted log files

### Concurrency Tests
- **TestConcurrentAppend** - Multiple writers
- **TestGroupCommit** - Batch flushing
- **TestConcurrentRecovery** - Prevent concurrent recovery

### Performance Tests
- **TestThroughput** - Append performance
- **TestGroupCommitLatency** - Commit latency reduction
- **TestRecoveryTime** - Recovery speed
- **TestLargeLog** - Handle large log files

## Benchmarks

```go
BenchmarkAppend             - Single record append
BenchmarkAppendNoSync       - Append without fsync
BenchmarkGroupCommit        - Group commit throughput
BenchmarkRecovery           - Recovery performance
BenchmarkCheckpoint         - Checkpoint overhead
```

## Implementation Hints

### Log Record Encoding
```go
func (r *LogRecord) Encode() []byte {
	buf := make([]byte, 25+len(r.Data))

	// Header
	binary.LittleEndian.PutUint64(buf[0:8], uint64(r.LSN))
	buf[8] = byte(r.Type)
	binary.LittleEndian.PutUint64(buf[9:17], uint64(r.TxnID))
	binary.LittleEndian.PutUint32(buf[17:21], uint32(len(r.Data)))

	// Data
	copy(buf[25:], r.Data)

	// Checksum (over entire record)
	r.Checksum = crc32.ChecksumIEEE(buf[0:25+len(r.Data)])
	binary.LittleEndian.PutUint32(buf[21:25], r.Checksum)

	return buf
}

func DecodeLogRecord(data []byte) (*LogRecord, error) {
	if len(data) < 25 {
		return nil, ErrInvalidRecord
	}

	record := &LogRecord{
		LSN:   LSN(binary.LittleEndian.Uint64(data[0:8])),
		Type:  RecordType(data[8]),
		TxnID: TxnID(binary.LittleEndian.Uint64(data[9:17])),
	}

	dataLen := binary.LittleEndian.Uint32(data[17:21])
	checksum := binary.LittleEndian.Uint32(data[21:25])

	if len(data) < int(25+dataLen) {
		return nil, ErrTruncatedRecord
	}

	record.Data = make([]byte, dataLen)
	copy(record.Data, data[25:25+dataLen])

	// Verify checksum
	computed := crc32.ChecksumIEEE(data[0:25+dataLen])
	if computed != checksum {
		return nil, ErrChecksumMismatch
	}

	record.Checksum = checksum
	return record, nil
}
```

### Group Commit Flusher
```go
type GroupCommitFlusher struct {
	wal       *WAL
	interval  time.Duration
	stopCh    chan struct{}
	doneCh    chan struct{}
	commitCh  chan chan error  // requests for sync
}

func (f *GroupCommitFlusher) Start() {
	go func() {
		ticker := time.NewTicker(f.interval)
		defer ticker.Stop()
		defer close(f.doneCh)

		var waiters []chan error

		for {
			select {
			case waiter := <-f.commitCh:
				waiters = append(waiters, waiter)

			case <-ticker.C:
				if len(waiters) > 0 {
					err := f.wal.flushInternal()
					for _, ch := range waiters {
						ch <- err
						close(ch)
					}
					waiters = waiters[:0]
				}

			case <-f.stopCh:
				// Final flush
				if len(waiters) > 0 {
					err := f.wal.flushInternal()
					for _, ch := range waiters {
						ch <- err
						close(ch)
					}
				}
				return
			}
		}
	}()
}

func (f *GroupCommitFlusher) Commit() error {
	waiter := make(chan error, 1)
	f.commitCh <- waiter
	return <-waiter
}
```

### Recovery Algorithm
```go
func (w *WAL) Recover(handler RecoveryHandler) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Reset to beginning of log
	if _, err := w.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader := bufio.NewReader(w.file)
	maxLSN := LSN(0)

	for {
		// Read record length
		var header [25]byte
		if _, err := io.ReadFull(reader, header[:]); err != nil {
			if err == io.EOF {
				break // End of log
			}
			return err
		}

		dataLen := binary.LittleEndian.Uint32(header[17:21])

		// Read full record
		fullRecord := make([]byte, 25+dataLen)
		copy(fullRecord, header[:])
		if _, err := io.ReadFull(reader, fullRecord[25:]); err != nil {
			// Partial record - truncate here
			break
		}

		// Decode and validate
		record, err := DecodeLogRecord(fullRecord)
		if err != nil {
			// Corrupted record - stop recovery
			break
		}

		// Track highest LSN
		if record.LSN > maxLSN {
			maxLSN = record.LSN
		}

		// Call handler
		if err := w.handleRecord(handler, record); err != nil {
			return err
		}
	}

	// Set current LSN to continue after recovery
	w.currentLSN.Store(uint64(maxLSN) + 1)

	return nil
}

func (w *WAL) handleRecord(handler RecoveryHandler, record *LogRecord) error {
	switch record.Type {
	case RecordBegin:
		return handler.OnBegin(record.TxnID, record.LSN)
	case RecordCommit:
		return handler.OnCommit(record.TxnID, record.LSN)
	case RecordAbort:
		return handler.OnAbort(record.TxnID, record.LSN)
	case RecordUpdate:
		return handler.OnUpdate(record.TxnID, record.LSN, record.Data)
	case RecordCheckpoint:
		return handler.OnCheckpoint(record.LSN)
	default:
		return ErrUnknownRecordType
	}
}
```

### Crash Testing Helper
```go
// CrashSimulator helps test crash recovery
type CrashSimulator struct {
	wal        *WAL
	crashAfter int  // crash after N records
	count      int
}

func (cs *CrashSimulator) Append(record *LogRecord) (LSN, error) {
	cs.count++
	if cs.count == cs.crashAfter {
		// Simulate crash: close without flush
		cs.wal.file.Close()
		return 0, ErrSimulatedCrash
	}
	return cs.wal.Append(record)
}
```

## Performance Goals

- Append (buffered): < 1µs
- Append + fsync: < 1ms (SSD)
- Group commit throughput: > 10K commits/sec
- Recovery: > 1M records/sec
- Checkpoint: < 100ms
- Log space overhead: < 10% vs data size

## Stretch Goals

### 1. Parallel Recovery
- Partition log by transaction
- Replay independent transactions in parallel
- Measure speedup

### 2. Log Compression
- Compress log records
- Benchmark compression ratio vs speed
- Consider streaming compression

### 3. Remote Replication
- Stream log to remote server
- Asynchronous replication
- Monitor replication lag

### 4. Log Archival
- Archive old log files
- Restore from archive
- Integration with checkpoints

## Common Pitfalls

1. **Partial Writes**
   - Always write complete records
   - Detect and handle truncated records
   - Use checksums to validate

2. **fsync Errors**
   - Check fsync return value
   - Handle disk full scenarios
   - Retry logic for transient errors

3. **LSN Gaps**
   - Ensure LSNs are sequential
   - Handle wrap-around (if using limited bits)
   - Atomic LSN generation

4. **Recovery Bugs**
   - Test with various crash points
   - Verify idempotency
   - Handle partially committed transactions

5. **Resource Leaks**
   - Close file handles
   - Stop background goroutines
   - Flush buffers on close

## Debugging Tips

```bash
# Inspect log file
hexdump -C wal.log | less

# Verify checksums
go run tools/verify_log.go wal.log

# Simulate crashes
go test -run=TestCrash -count=100

# Profile recovery
go test -bench=BenchmarkRecovery -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## Validation Checklist

Your implementation should:
- [ ] Pass all unit tests
- [ ] Survive crash simulations
- [ ] Verify all checksums
- [ ] Handle corrupted logs gracefully
- [ ] Support concurrent appends
- [ ] Implement group commit
- [ ] Recover all committed transactions
- [ ] Lose no committed data
- [ ] Truncate log correctly
- [ ] No data races (go test -race)

## Learning Outcomes

After completing this project, you will understand:
- Write-Ahead Logging principles
- Durability and crash recovery
- fsync and disk persistence
- Group commit optimization
- Log sequence numbers
- Idempotent recovery
- Checkpointing strategies
- Chaos engineering for databases

## Time Estimate
- Core implementation: 6-8 hours
- Testing and crash scenarios: 2-3 hours
- Benchmarking and optimization: 2-3 hours
- Stretch goals: 5-7 hours (optional)

## Next Steps
After completing this project, move on to **Phase 2: Graph Structure Projects** which build on storage concepts with advanced data structures.
