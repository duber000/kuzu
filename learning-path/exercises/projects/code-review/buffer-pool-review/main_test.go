package bufferpoolreview

import (
	"sync"
	"testing"
)

// Write tests that expose the bugs
func TestConcurrentFetch(t *testing.T) {
	// TODO: Test that should fail due to race conditions
	bp := NewBufferPool(10)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, _ = bp.FetchPage(PageID(id % 5))
		}(i)
	}
	wg.Wait()
}

func TestUnpinNegative(t *testing.T) {
	// TODO: Test that pin count can go negative (bug)
	t.Skip("not implemented")
}

func TestEvictionLeak(t *testing.T) {
	// TODO: Test that free list doesn't grow after eviction (bug)
	t.Skip("not implemented")
}

// Run with -race to find concurrency issues
