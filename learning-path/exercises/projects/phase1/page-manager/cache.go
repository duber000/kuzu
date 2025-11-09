package pagemanager

import (
	"container/list"
	"sync"
)

// LRUCache implements a least-recently-used cache for pages
type LRUCache struct {
	capacity int
	pages    map[PageID]*list.Element
	lru      *list.List
	mu       sync.RWMutex
	hits     uint64
	misses   uint64
}

type cacheEntry struct {
	pageID PageID
	page   *Page
}

// NewLRUCache creates a new LRU cache
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		pages:    make(map[PageID]*list.Element),
		lru:      list.New(),
	}
}

// Get retrieves a page from cache
func (c *LRUCache) Get(pageID PageID) (*Page, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: Implement cache get
	// - Check if page exists
	// - Move to front of LRU list
	// - Update hit/miss stats
	// - Return page and found status

	if elem, ok := c.pages[pageID]; ok {
		c.lru.MoveToFront(elem)
		c.hits++
		return elem.Value.(*cacheEntry).page, true
	}

	c.misses++
	return nil, false
}

// Put adds a page to cache
func (c *LRUCache) Put(page *Page) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: Implement cache put
	// - Check if page already exists, update it
	// - If at capacity, evict LRU page
	// - Add new page to front
	// - If evicted page is dirty, need to flush first

	if elem, ok := c.pages[page.ID]; ok {
		// Update existing
		c.lru.MoveToFront(elem)
		elem.Value.(*cacheEntry).page = page
		return
	}

	// Add new
	if c.lru.Len() >= c.capacity {
		// Evict LRU
		oldest := c.lru.Back()
		if oldest != nil {
			entry := oldest.Value.(*cacheEntry)
			delete(c.pages, entry.pageID)
			c.lru.Remove(oldest)
			// TODO: Flush if dirty
		}
	}

	entry := &cacheEntry{pageID: page.ID, page: page}
	elem := c.lru.PushFront(entry)
	c.pages[page.ID] = elem
}

// Remove removes a page from cache
func (c *LRUCache) Remove(pageID PageID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: Implement removal
	if elem, ok := c.pages[pageID]; ok {
		delete(c.pages, pageID)
		c.lru.Remove(elem)
	}
}

// Evict evicts all dirty pages and returns them for flushing
func (c *LRUCache) Evict() []*Page {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: Implement eviction
	// - Collect all dirty pages
	// - Clear cache
	// - Return dirty pages for flushing

	return nil
}

// Stats returns cache statistics
func (c *LRUCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := 0.0
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return CacheStats{
		Hits:    c.hits,
		Misses:  c.misses,
		HitRate: hitRate,
		Size:    c.lru.Len(),
	}
}

// CacheStats holds cache statistics
type CacheStats struct {
	Hits    uint64
	Misses  uint64
	HitRate float64
	Size    int
}
