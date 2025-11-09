package bufferpoolreview

import "sync"

// REVIEW THIS CODE: Find all issues!

type PageID int
type FrameID int

type Frame struct {
	pageID   PageID
	data     []byte
	pinCount int  // Issue 1: Not thread-safe
	dirty    bool
}

type BufferPool struct {
	frames    map[FrameID]*Frame
	pageTable map[PageID]FrameID
	freeList  []FrameID
	mu        sync.Mutex
}

func NewBufferPool(size int) *BufferPool {
	bp := &BufferPool{
		frames:    make(map[FrameID]*Frame),
		pageTable: make(map[PageID]FrameID),
		freeList:  make([]FrameID, size),
	}

	for i := 0; i < size; i++ {
		bp.freeList[i] = FrameID(i)
		bp.frames[FrameID(i)] = &Frame{
			data: make([]byte, 4096),
		}
	}

	return bp
}

// Issue 2: Holding lock during I/O
func (bp *BufferPool) FetchPage(pageID PageID) (*Frame, error) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	// Check if already in buffer
	if frameID, ok := bp.pageTable[pageID]; ok {
		frame := bp.frames[frameID]
		frame.pinCount++  // Issue 3: Race condition
		return frame, nil
	}

	// Get free frame
	if len(bp.freeList) == 0 {
		return nil, nil  // Issue 4: Wrong error handling
	}

	frameID := bp.freeList[0]
	bp.freeList = bp.freeList[1:]

	frame := bp.frames[frameID]
	frame.pageID = pageID
	frame.pinCount++

	// Issue 5: I/O while holding lock!
	readPageFromDisk(pageID, frame.data)

	bp.pageTable[pageID] = frameID

	return frame, nil
}

// Issue 6: Not thread-safe
func (bp *BufferPool) UnpinPage(pageID PageID, dirty bool) {
	frameID := bp.pageTable[pageID]  // Missing lock
	frame := bp.frames[frameID]
	frame.pinCount--  // Issue 7: Can go negative
	if dirty {
		frame.dirty = true  // Issue 8: Race condition
	}
}

// Issue 9: Memory leak - frames never freed
func (bp *BufferPool) Evict(frameID FrameID) {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	frame := bp.frames[frameID]

	if frame.pinCount > 0 {
		return  // Issue 10: Silent failure
	}

	if frame.dirty {
		writePageToDisk(frame.pageID, frame.data)
	}

	delete(bp.pageTable, frame.pageID)
	// Issue 11: frame not added back to free list!
}

// Issue 12: Background goroutine never stopped
func (bp *BufferPool) StartFlusher() {
	go func() {
		for {  // Issue 13: Infinite loop, no stop channel
			bp.FlushAll()
		}
	}()
}

func (bp *BufferPool) FlushAll() {
	// Issue 14: Wrong locking pattern
	for _, frame := range bp.frames {
		bp.mu.Lock()  // Issue 15: Lock/unlock in loop
		if frame.dirty {
			writePageToDisk(frame.pageID, frame.data)
			frame.dirty = false
		}
		bp.mu.Unlock()
	}
}

func readPageFromDisk(pageID PageID, data []byte) {
	// Simulate I/O
}

func writePageToDisk(pageID PageID, data []byte) {
	// Simulate I/O
}

// TODO: Find and document all issues
// Create a review document listing:
// 1. Issue description
// 2. Location (function + line)
// 3. Severity (critical/major/minor)
// 4. Suggested fix
