package jsonpathplus

import (
	"container/list"
	"sync"
	"time"
)

// CacheEntry represents a cached compilation result
type CacheEntry struct {
	AST       *astNode
	Timestamp time.Time
	HitCount  int64
}

// Cache interface for compiled JSONPath expressions
type Cache interface {
	Get(path string) (*astNode, bool)
	Put(path string, ast *astNode)
	Clear()
	Size() int
	Stats() CacheStats
}

// CacheStats provides cache statistics
type CacheStats struct {
	Size     int
	Hits     int64
	Misses   int64
	HitRatio float64
}

// LRUCache implements a thread-safe LRU cache
type LRUCache struct {
	mu       sync.RWMutex
	capacity int
	entries  map[string]*list.Element
	order    *list.List
	hits     int64
	misses   int64
}

type lruEntry struct {
	key   string
	value *CacheEntry
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		capacity = 100
	}
	
	return &LRUCache{
		capacity: capacity,
		entries:  make(map[string]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves an entry from the cache
func (c *LRUCache) Get(path string) (*astNode, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if elem, exists := c.entries[path]; exists {
		c.hits++
		c.order.MoveToFront(elem)
		entry := elem.Value.(*lruEntry).value
		entry.HitCount++
		entry.Timestamp = time.Now()
		return entry.AST, true
	}
	
	c.misses++
	return nil, false
}

// Put stores an entry in the cache
func (c *LRUCache) Put(path string, ast *astNode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if elem, exists := c.entries[path]; exists {
		elem.Value.(*lruEntry).value = &CacheEntry{
			AST:       ast,
			Timestamp: time.Now(),
			HitCount:  0,
		}
		c.order.MoveToFront(elem)
		return
	}
	
	if c.order.Len() >= c.capacity {
		// Remove least recently used
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.entries, oldest.Value.(*lruEntry).key)
		}
	}
	
	entry := &lruEntry{
		key: path,
		value: &CacheEntry{
			AST:       ast,
			Timestamp: time.Now(),
			HitCount:  0,
		},
	}
	
	elem := c.order.PushFront(entry)
	c.entries[path] = elem
}

// Clear removes all entries from the cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.entries = make(map[string]*list.Element)
	c.order = list.New()
	c.hits = 0
	c.misses = 0
}

// Size returns the number of entries in the cache
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Stats returns cache statistics
func (c *LRUCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	total := c.hits + c.misses
	hitRatio := 0.0
	if total > 0 {
		hitRatio = float64(c.hits) / float64(total)
	}
	
	return CacheStats{
		Size:     len(c.entries),
		Hits:     c.hits,
		Misses:   c.misses,
		HitRatio: hitRatio,
	}
}

// NoCache is a cache implementation that doesn't cache anything
type NoCache struct{}

func (c *NoCache) Get(path string) (*astNode, bool) { return nil, false }
func (c *NoCache) Put(path string, ast *astNode)    {}
func (c *NoCache) Clear()                           {}
func (c *NoCache) Size() int                        { return 0 }
func (c *NoCache) Stats() CacheStats {
	return CacheStats{Size: 0, Hits: 0, Misses: 0, HitRatio: 0}
}