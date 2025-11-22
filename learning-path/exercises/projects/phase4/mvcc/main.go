package mvcc

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Type definitions
type Key string
type Value []byte
type TxnID uint64
type Timestamp uint64

// Errors
var (
	ErrWriteConflict = errors.New("write conflict")
	ErrKeyNotFound   = errors.New("key not found")
	ErrTxnAborted    = errors.New("transaction aborted")
)

// Version represents a single version of a value
type Version struct {
	data    Value
	beginTS Timestamp
	endTS   *Timestamp  // nil if latest version
	txnID   TxnID
	prev    *Version    // previous version (could use weak.Pointer in Go 1.24)
}

// VersionChain is a linked list of versions
type VersionChain struct {
	latest *Version
	mu     sync.RWMutex
}

// Transaction represents an active transaction
type Transaction struct {
	id       TxnID
	snapshot Timestamp
	writeSet map[Key]*Version
	readSet  map[Key]Timestamp
	mu       sync.Mutex
}

// MVCCStore implements multi-version concurrency control
type MVCCStore struct {
	data        map[Key]*VersionChain
	transactions map[TxnID]*Transaction
	clock       atomic.Uint64
	mu          sync.RWMutex
	gc          *GarbageCollector
}

// NewMVCCStore creates a new MVCC store
func NewMVCCStore() *MVCCStore {
	store := &MVCCStore{
		data:        make(map[Key]*VersionChain),
		transactions: make(map[TxnID]*Transaction),
	}
	store.gc = NewGarbageCollector(store)
	return store
}

// BeginTransaction starts a new transaction
func (s *MVCCStore) BeginTransaction() *Transaction {
	// TODO: Implement transaction start
	// Assign snapshot timestamp
	return nil
}

// Read reads a value at the transaction's snapshot
func (s *MVCCStore) Read(txn *Transaction, key Key) (Value, error) {
	// TODO: Implement MVCC read
	// 1. Get version chain
	// 2. Find visible version based on snapshot
	// 3. Add to read set
	return nil, nil
}

// Write writes a value in the transaction
func (s *MVCCStore) Write(txn *Transaction, key Key, value Value) error {
	// TODO: Implement MVCC write
	// 1. Add to write set (don't write yet)
	// 2. Will be applied at commit
	return nil
}

// Commit commits a transaction
func (s *MVCCStore) Commit(txn *Transaction) error {
	// TODO: Implement commit
	// 1. Check for write conflicts
	// 2. Apply write set
	// 3. Create new versions
	// 4. Update timestamps
	return nil
}

// Abort aborts a transaction
func (s *MVCCStore) Abort(txn *Transaction) error {
	// TODO: Implement abort
	// Discard write set
	return nil
}

// isVisible checks if a version is visible to a transaction
func (s *MVCCStore) isVisible(version *Version, snapshot Timestamp) bool {
	// TODO: Implement visibility check
	// version.beginTS <= snapshot && (version.endTS == nil || version.endTS > snapshot)
	return false
}

// detectConflict checks for write-write conflicts
func (s *MVCCStore) detectConflict(txn *Transaction) error {
	// TODO: Implement conflict detection
	// Check if any written keys were modified after snapshot
	return nil
}

// GarbageCollector removes old versions
type GarbageCollector struct {
	store *MVCCStore
	stopCh chan struct{}
	doneCh chan struct{}
}

// NewGarbageCollector creates a garbage collector
func NewGarbageCollector(store *MVCCStore) *GarbageCollector {
	gc := &GarbageCollector{
		store:  store,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
	go gc.run()
	return gc
}

func (gc *GarbageCollector) run() {
	defer close(gc.doneCh)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gc.collect()
		case <-gc.stopCh:
			return
		}
	}
}

func (gc *GarbageCollector) collect() {
	// TODO: Implement garbage collection
	// 1. Find oldest active transaction snapshot
	// 2. Remove versions older than oldest snapshot
}

func (gc *GarbageCollector) Stop() {
	close(gc.stopCh)
	<-gc.doneCh
}
