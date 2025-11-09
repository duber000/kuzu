package pagemanager

import (
	"os"
	"sync"
)

// PageManager manages pages on disk with caching
type PageManager struct {
	file       *os.File
	pageSize   int
	cache      *LRUCache
	freeBitmap *Bitmap
	mu         sync.RWMutex
	nextPageID PageID
}

// New creates a new page manager
func New(filename string, cacheSize int) (*PageManager, error) {
	// TODO: Implement initialization
	// - Open or create file
	// - Load or initialize bitmap
	// - Create cache
	// - Read file header

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	pm := &PageManager{
		file:       file,
		pageSize:   PageSize,
		cache:      NewLRUCache(cacheSize),
		freeBitmap: NewBitmap(1000), // Initial size
		nextPageID: 0,
	}

	return pm, nil
}

// AllocatePage allocates a new page and returns its ID
func (pm *PageManager) AllocatePage() (PageID, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// TODO: Implement page allocation
	// - Find free page in bitmap
	// - If none, grow file
	// - Mark page as allocated
	// - Return page ID

	pageID := pm.nextPageID
	pm.nextPageID++

	return pageID, nil
}

// FreePage marks a page as free
func (pm *PageManager) FreePage(pageID PageID) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// TODO: Implement page freeing
	// - Validate page ID
	// - Evict from cache if present
	// - Mark as free in bitmap

	return nil
}

// ReadPage reads a page from disk (may come from cache)
func (pm *PageManager) ReadPage(pageID PageID) (*Page, error) {
	// TODO: Implement page reading
	// - Check cache first
	// - If not in cache, read from disk
	// - Add to cache
	// - Update LRU

	return nil, nil
}

// WritePage writes a page to disk (may be cached)
func (pm *PageManager) WritePage(page *Page) error {
	// TODO: Implement page writing
	// - Mark page as dirty
	// - Add to cache
	// - Optionally flush immediately

	return nil
}

// Flush writes all dirty pages to disk
func (pm *PageManager) Flush() error {
	// TODO: Implement flushing
	// - Iterate through cache
	// - Write all dirty pages
	// - Clear dirty flags

	return nil
}

// Close flushes and closes the page manager
func (pm *PageManager) Close() error {
	if err := pm.Flush(); err != nil {
		return err
	}
	return pm.file.Close()
}

// readPageFromDisk reads a page from disk at the given offset
func (pm *PageManager) readPageFromDisk(pageID PageID) (*Page, error) {
	// TODO: Implement disk read
	// - Calculate offset: pageID * pageSize
	// - Seek to position
	// - Read PageSize bytes
	// - Unmarshal into Page

	return nil, nil
}

// writePageToDisk writes a page to disk
func (pm *PageManager) writePageToDisk(page *Page) error {
	// TODO: Implement disk write
	// - Calculate offset
	// - Marshal page
	// - Seek and write
	// - Sync to disk

	return nil
}
