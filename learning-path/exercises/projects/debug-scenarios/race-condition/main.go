package racehunt

import "sync"

// BUGGY: Buffer Pool with race conditions
type BufferPool struct {
	frames map[int]*Frame
	mu     sync.Mutex
}

type Frame struct {
	data     []byte
	pinCount int  // BUGGY: should be atomic or protected
}

func (bp *BufferPool) FetchPage(pageID int) *Frame {
	bp.mu.Lock()
	frame := bp.frames[pageID]
	bp.mu.Unlock()

	// RACE: pinCount incremented without lock!
	frame.pinCount++

	return frame
}

func (bp *BufferPool) UnpinPage(pageID int) {
	bp.mu.Lock()
	frame := bp.frames[pageID]
	bp.mu.Unlock()

	// RACE: pinCount decremented without lock!
	frame.pinCount--
}

// BUGGY: Hash Join with race conditions
type Row map[string]interface{}

func HashJoin(left, right []Row) []Row {
	results := make([]Row, 0)
	var wg sync.WaitGroup

	for _, row := range left {
		wg.Add(1)
		go func(r Row) {
			defer wg.Done()
			// RACE: concurrent appends!
			results = append(results, r)
		}(row)
	}

	wg.Wait()
	return results
}

// BUGGY: Lock Manager with race conditions
type LockManager struct {
	locks map[string][]int  // resource -> txn holders
	waitForGraph map[int][]int  // BUGGY: concurrent access
}

func (lm *LockManager) AcquireLock(txn int, resource string) {
	// RACE: reading and writing waitForGraph without lock!
	if len(lm.locks[resource]) > 0 {
		lm.waitForGraph[txn] = lm.locks[resource]
	}
}

// TODO: Fix these race conditions!
// Hints:
// 1. Use atomic operations for counters
// 2. Protect shared slices/maps with mutex
// 3. Use channels for synchronization
// 4. Test with -race flag
