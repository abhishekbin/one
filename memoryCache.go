package onecache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type memoryCacheItem struct {
	key   string
	value interface{}
}

type MemoryCache struct {
	// MaxItems specifies the maximum number of items this cache should have
	// at any given time. A value of 0 means unbounded.
	// When adding a new item would put this cache over capacity,
	// the least recently accessed item(s) will be purged to make room.
	MaxItems int

	// To purge the least recently accessed item(s) when over capacity,
	// we implement a standard LRU cache using this doubly linked list + map.
	//
	// The doubly linked list, sorted by most recently read items first,
	// allows us to both query the least recently read items quickly *and*
	// move or add items (on read & write respectively) to the front quickly,
	// while the map allows us to look items up by key quickly.

	mostRecentlyRead *list.List // Elements’ Values are memoryCacheItems
	elementsByKey    map[string]*list.Element

	// Guards against concurrent accesses to structures above, particularly the map.
	mu sync.Mutex

	scheduler Scheduler
}

func (c *MemoryCache) initIfNeeded() {
	if c.mostRecentlyRead == nil {
		c.mostRecentlyRead = list.New()
	}
	if c.elementsByKey == nil {
		c.elementsByKey = map[string]*list.Element{}
	}
	if c.scheduler == nil {
		c.scheduler = DefaultScheduler{}
	}
}

// Get queries this cache for the value under the given key.
// If found (& unexpired), returns the value -- which may be nil -- and true.
// If not found (or expired), returns nil and false.
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initIfNeeded()

	// Do we have this key in our cache?
	elmt, ok := c.elementsByKey[key]
	if !ok {
		return nil, false, nil
	}
	item := elmt.Value.(memoryCacheItem)

	// TODO: Check for expiry, and clear if expired.

	// Mark as most recently read.
	c.mostRecentlyRead.MoveToFront(elmt)

	return item.value, true, nil
}

// Set caches the given value (which may be nil) under the given key,
// optionally expiring after the given duration.
// A duration of 0 means the value should never expire.
func (c *MemoryCache) Set(ctx context.Context, key string, val interface{}, expireAfter time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initIfNeeded()

	// Add item.
	// TODO: Store expiry too, and clear when expired.
	item := memoryCacheItem{key, val}
	elmt := c.mostRecentlyRead.PushFront(item)
	c.elementsByKey[key] = elmt

	// If we're over capacity, evict least recently read items.
	for c.MaxItems > 0 && len(c.elementsByKey) > c.MaxItems {
		oldestItem := c.mostRecentlyRead.Back().Value.(memoryCacheItem)
		c.clearWhenLocked(ctx, oldestItem.key)
	}

	return nil
}

// Clear clears the value, if any, cached under the given key.
// Returns whether a value was cached (and thus cleared).
func (c *MemoryCache) Clear(ctx context.Context, key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initIfNeeded()

	return c.clearWhenLocked(ctx, key)
}

func (c *MemoryCache) clearWhenLocked(ctx context.Context, key string) (bool, error) {
	// Do we have this key in our cache? Noop if not.
	elmt, ok := c.elementsByKey[key]
	if !ok {
		return false, nil
	}

	// Clear item if so.
	c.mostRecentlyRead.Remove(elmt)
	delete(c.elementsByKey, key)

	return true, nil
}

var _ Cache = (*MemoryCache)(nil)
