# Challenge 3: Distributed Transactions

## Overview
Implement distributed transactions using Two-Phase Commit (2PC) protocol with transaction coordinator, failure recovery, and distributed deadlock detection.

**Duration:** 50-60 hours
**Difficulty:** Extremely Hard

## Concepts

### Two-Phase Commit (2PC)
```
Coordinator                Participants
    |                          |
    |--- PREPARE ------------->|
    |                          | (vote YES/NO)
    |<-- VOTE ----------------|
    |                          |
    |--- COMMIT/ABORT -------->|
    |                          | (commit/abort)
    |<-- ACK -----------------|
```

### Roles
- **Coordinator:** Manages transaction lifecycle
- **Participants:** Execute operations and vote
- **Transaction Manager:** Coordinates across nodes

## Protocol Phases

### Phase 1: Prepare
1. Coordinator sends PREPARE to all participants
2. Participants write to WAL and reply YES/NO
3. If any NO or timeout, abort
4. If all YES, proceed to Phase 2

### Phase 2: Commit/Abort
1. Coordinator writes decision to log
2. Coordinator sends COMMIT or ABORT
3. Participants commit/abort and reply ACK
4. Coordinator completes transaction

## Failure Scenarios

### Participant Failure
- Coordinator times out on PREPARE
- Abort transaction
- Participant recovers and rolls back

### Coordinator Failure
- Participants block waiting for decision
- Recovery protocol reads coordinator log
- Re-send commit/abort to participants

### Network Partition
- Use timeouts and retries
- Eventually abort if can't reach majority

## API Design

```go
type TransactionCoordinator struct {
	participants []Participant
	txnLog       *TxnLog
}

type Participant interface {
	Prepare(txnID TxnID, operations []Operation) (Vote, error)
	Commit(txnID TxnID) error
	Abort(txnID TxnID) error
}

type Vote int
const (
	VoteYes Vote = iota
	VoteNo
	VoteAbort
)

// Begin distributed transaction
func (tc *TransactionCoordinator) Begin() (TxnID, error)

// Execute operation on participant
func (tc *TransactionCoordinator) Execute(txnID TxnID, participantID int, op Operation) error

// Commit distributed transaction
func (tc *TransactionCoordinator) Commit(txnID TxnID) error

// Abort distributed transaction
func (tc *TransactionCoordinator) Abort(txnID TxnID) error
```

## Distributed Deadlock Detection

### Global Wait-For Graph
- Each node maintains local wait-for graph
- Periodically exchange graphs with other nodes
- Detect cycles in global graph
- Select victim and abort

### Phantom Deadlock
- False deadlock due to stale information
- Use timestamps to avoid

## Optimizations

### 1. Presumed Abort
If coordinator crashes after Phase 1, assume ABORT
- Reduces log writes
- Faster recovery

### 2. Read-Only Optimization
If participant only reads, skip Phase 2
- Reduces messages
- Faster commits

### 3. Three-Phase Commit (3PC)
Add CanCommit phase to avoid blocking
- More complex
- Better availability

## Test Cases

### Correctness
- All participants commit
- One participant fails -> all abort
- Coordinator fails -> recovery works
- Network partition -> safe abort
- Concurrent transactions

### Stress Tests
- 1000 concurrent transactions
- Random participant failures
- Network delays and partitions
- Distributed deadlocks

## Performance Goals

- Transaction latency: <50ms (3 nodes)
- Throughput: >1000 txns/sec
- Recovery time: <5s
- Deadlock detection: <100ms

## Implementation Hints

### Transaction Log
```go
type TxnLog struct {
	entries map[TxnID]*LogEntry
	mu      sync.RWMutex
}

type LogEntry struct {
	TxnID        TxnID
	State        TxnState  // PREPARING, PREPARED, COMMITTED, ABORTED
	Participants []int
	Decision     Decision
}
```

### Timeout Handling
```go
func (tc *TransactionCoordinator) Prepare(txnID TxnID) error {
	votes := make(chan Vote, len(tc.participants))

	for _, p := range tc.participants {
		go func(participant Participant) {
			vote, err := participant.Prepare(txnID, ops)
			if err != nil {
				votes <- VoteAbort
			} else {
				votes <- vote
			}
		}(p)
	}

	// Wait for votes with timeout
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for i := 0; i < len(tc.participants); i++ {
		select {
		case vote := <-votes:
			if vote != VoteYes {
				return ErrVoteNo
			}
		case <-timer.C:
			return ErrTimeout
		}
	}

	return nil
}
```

## Stretch Goals

### 1. Paxos Commit
Use Paxos for coordinator consensus
- Higher availability
- No single point of failure

### 2. Calvin
Deterministic transaction ordering
- No 2PC overhead
- Better performance

### 3. Spanner-style
Use TrueTime for global ordering
- External consistency
- Distributed reads

## Learning Outcomes
- Distributed consensus protocols
- Failure recovery
- Network partition handling
- Distributed deadlock detection
- System design at scale

## References
- "Transaction Processing" by Jim Gray
- "Designing Data-Intensive Applications" by Martin Kleppmann
- Google Spanner paper
- Calvin paper (Yale)

## Time Estimate
Core 2PC: 25-30 hours
Failure recovery: 10-12 hours
Deadlock detection: 8-10 hours
Testing: 7-10 hours
