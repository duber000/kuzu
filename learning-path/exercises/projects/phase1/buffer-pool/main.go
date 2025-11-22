package bufferpool

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Constants
const (
	PageSize = 4096 // 4KB pages
)

// Type definitions
type FrameID int32
type PageID int64

// Errors
var (
	ErrNoVictimFrame = errors.New("no victim frame available: all pages pinned")
	ErrInvalidPageID = errors.New("invalid page ID")
	ErrPageNotFound  = errors.New("page not found in pool")
)

// DiskManager interface for reading/writing pages to disk
type DiskManager interface {
	ReadPage(pageID PageID, data []byte) error
	WritePage(pageID PageID, data []byte) error
	AllocatePage() (PageID, error)
	DeallocatePage(pageID PageID) error
}

// Frame represents an in-memory slot for a page
type Frame struct {
	frameID  FrameID
	pageID   PageID
	data     [PageSize]byte
	pinCount atomic.Int32
	dirty    atomic.Bool
	mu       sync.RWMutex
}

// Pin increments the pin count
func (f *Frame) Pin() {
	f.pinCount.Add(1)
}

// Unpin decrements the pin count
func (f *Frame) Unpin() {
	if f.pinCount.Add(-1) < 0 {
		panic("unpin of unpinned frame")
	}
}

// IsPinned returns true if the frame is pinned
func (f *Frame) IsPinned() bool {
	return f.pinCount.Load() > 0
}

// MarkDirty marks the frame as dirty
func (f *Frame) MarkDirty() {
	f.dirty.Store(true)
}

// IsDirty returns true if the frame is dirty
func (f *Frame) IsDirty() bool {
	return f.dirty.Load()
}

// Data returns a pointer to the frame's data
func (f *Frame) Data() []byte {
	return f.data[:]
}

// LRUReplacer implements LRU eviction policy
type LRUReplacer struct {
	capacity int
	frames   map[FrameID]*list.Element
	lruList  *list.List
	mu       sync.Mutex
}

// NewLRUReplacer creates a new LRU replacer
func NewLRUReplacer(capacity int) *LRUReplacer {
	return &LRUReplacer{
		capacity: capacity,
		frames:   make(map[FrameID]*list.Element),
		lruList:  list.New(),
	}
}

// RecordAccess records that a frame was accessed
func (r *LRUReplacer) RecordAccess(frameID FrameID) {
	// TODO: Implement LRU access tracking
	// Move frame to front of LRU list
}

// Victim returns a victim frame for eviction
func (r *LRUReplacer) Victim() (FrameID, bool) {
	// TODO: Implement victim selection
	// Return least recently used unpinned frame
	return -1, false
}

// Remove removes a frame from the replacer
func (r *LRUReplacer) Remove(frameID FrameID) {
	// TODO: Implement frame removal
	// Remove frame from LRU tracking
}

// Size returns the number of frames in the replacer
func (r *LRUReplacer) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lruList.Len()
}

// BackgroundFlusher periodically flushes dirty pages
type BackgroundFlusher struct {
	pool     *BufferPool
	interval time.Duration
	stopCh   chan struct{}
	doneCh   chan struct{}
}

// NewBackgroundFlusher creates a new background flusher
func NewBackgroundFlusher(pool *BufferPool, interval time.Duration) *BackgroundFlusher {
	return &BackgroundFlusher{
		pool:     pool,
		interval: interval,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
}

// Start starts the background flusher goroutine
func (f *BackgroundFlusher) Start() {
	// TODO: Implement background flusher
	// Periodically flush dirty pages using ticker
}

// Stop stops the background flusher
func (f *BackgroundFlusher) Stop() {
	// TODO: Implement graceful shutdown
	// Close stopCh and wait for doneCh
}

// flushDirtyPages flushes all dirty pages
func (f *BackgroundFlusher) flushDirtyPages() {
	// TODO: Implement dirty page flushing
	// Iterate through frames and flush dirty ones
}

// PoolStats contains buffer pool statistics
type PoolStats struct {
	TotalFrames  int
	PinnedFrames int
	DirtyFrames  int
	FreeFrames   int
	CacheHits    int64
	CacheMisses  int64
}

// BufferPool manages a pool of page frames
type BufferPool struct {
	frames      []*Frame
	pageTable   map[PageID]FrameID
	freeList    []FrameID
	replacer    *LRUReplacer
	diskManager DiskManager
	mu          sync.RWMutex
	flusher     *BackgroundFlusher
	cacheHits   atomic.Int64
	cacheMisses atomic.Int64
}

// New creates a new buffer pool
func New(diskManager DiskManager, poolSize int) *BufferPool {
	bp := &BufferPool{
		frames:      make([]*Frame, poolSize),
		pageTable:   make(map[PageID]FrameID),
		freeList:    make([]FrameID, poolSize),
		replacer:    NewLRUReplacer(poolSize),
		diskManager: diskManager,
	}

	// Initialize frames and free list
	for i := 0; i < poolSize; i++ {
		bp.frames[i] = &Frame{
			frameID: FrameID(i),
			pageID:  -1,
		}
		bp.freeList[i] = FrameID(i)
	}

	// Start background flusher
	bp.flusher = NewBackgroundFlusher(bp, 5*time.Second)
	bp.flusher.Start()

	return bp
}

// FetchPage fetches a page from the pool or disk
func (bp *BufferPool) FetchPage(pageID PageID) (*Frame, error) {
	// TODO: Implement page fetching
	// 1. Check if page is already in pool (cache hit)
	// 2. If not, find a victim frame (from free list or via eviction)
	// 3. If victim is dirty, flush it
	// 4. Load new page from disk
	// 5. Update page table and pin the frame
	return nil, nil
}

// UnpinPage unpins a page and marks it dirty if modified
func (bp *BufferPool) UnpinPage(pageID PageID, dirty bool) error {
	// TODO: Implement unpinning
	// 1. Find frame in page table
	// 2. Decrement pin count
	// 3. Mark dirty if needed
	// 4. Add to replacer if pin count == 0
	return nil
}

// FlushPage flushes a specific page to disk
func (bp *BufferPool) FlushPage(pageID PageID) error {
	// TODO: Implement single page flush
	// 1. Find frame in page table
	// 2. If dirty, write to disk
	// 3. Clear dirty flag
	return nil
}

// FlushAll flushes all dirty pages to disk
func (bp *BufferPool) FlushAll() error {
	// TODO: Implement flush all
	// Iterate through all frames and flush dirty ones
	return nil
}

// NewPage allocates a new page
func (bp *BufferPool) NewPage() (PageID, *Frame, error) {
	// TODO: Implement page allocation
	// 1. Allocate page from disk manager
	// 2. Fetch the new page into pool
	// 3. Return page ID and frame
	return -1, nil, nil
}

// DeletePage deletes a page from pool and disk
func (bp *BufferPool) DeletePage(pageID PageID) error {
	// TODO: Implement page deletion
	// 1. Remove from page table
	// 2. Add frame to free list
	// 3. Deallocate from disk
	return nil
}

// Stats returns buffer pool statistics
func (bp *BufferPool) Stats() PoolStats {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	stats := PoolStats{
		TotalFrames: len(bp.frames),
		FreeFrames:  len(bp.freeList),
		CacheHits:   bp.cacheHits.Load(),
		CacheMisses: bp.cacheMisses.Load(),
	}

	// Count pinned and dirty frames
	for _, frame := range bp.frames {
		if frame.IsPinned() {
			stats.PinnedFrames++
		}
		if frame.IsDirty() {
			stats.DirtyFrames++
		}
	}

	return stats
}

// Close closes the buffer pool
func (bp *BufferPool) Close() error {
	// TODO: Implement cleanup
	// 1. Stop background flusher
	// 2. Flush all dirty pages
	// 3. Wait for all pins to be released (or timeout)
	return nil
}
