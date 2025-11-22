package racehunt

import (
	"sync"
	"testing"
)

func TestBufferPoolRace(t *testing.T) {
	// TODO: Write test that triggers the race
	// Run with: go test -race
	t.Skip("not implemented")
}

func TestHashJoinRace(t *testing.T) {
	// TODO: Write test that triggers the race
	t.Skip("not implemented")
}

func TestLockManagerRace(t *testing.T) {
	// TODO: Write test that triggers the race
	t.Skip("not implemented")
}

// Example test that should fail with -race
func TestConcurrentPinUnpin(t *testing.T) {
	bp := &BufferPool{
		frames: make(map[int]*Frame),
	}

	bp.frames[1] = &Frame{data: make([]byte, 100)}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			bp.FetchPage(1)
		}()
		go func() {
			defer wg.Done()
			bp.UnpinPage(1)
		}()
	}
	wg.Wait()
}
