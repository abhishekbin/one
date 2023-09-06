package onecache

import (
	"context"
	"testing"
	"time"
)

// Helpers:

type TestMemoryCache struct {
	cache *MemoryCache
	ctx   context.Context
	t     *testing.T
}

func (c TestMemoryCache) Get(key string) (interface{}, bool) {
	val, ok, err := c.cache.Get(c.ctx, key)
	assertNoError(c.t, err)
	return val, ok
}

func (c TestMemoryCache) Set(key string, val interface{}, expireAfter time.Duration) {
	assertNoError(c.t, c.cache.Set(c.ctx, key, val, expireAfter))
}

func (c TestMemoryCache) Clear(key string) bool {
	ok, err := c.cache.Clear(c.ctx, key)
	assertNoError(c.t, err)
	return ok
}

func (c TestMemoryCache) assertNotInCache(key string) {
	val, ok := c.Get(key)
	assertNot(c.t, ok, "Expected %s to not be in cache but was", key)
	assertNil(c.t, val)
}

func (c TestMemoryCache) assertCachedEquals(key string, expected interface{}) {
	val, ok := c.Get(key)
	assert(c.t, ok, "Expected %s to be in cache but wasn't", key)
	assertEquals(c.t, val, expected)
}

// Tests:

func TestMemoryCache_Basic(t *testing.T) {
	t.Parallel()

	// Without worrying about capacity or expiry, should be able to get, set, and clear.

	cache := TestMemoryCache{
		cache: &MemoryCache{},
		ctx:   context.Background(),
		t:     t,
	}

	// Initially, the cache should contain nothing.

	cache.assertNotInCache("a")
	cache.assertNotInCache("b")
	cache.assertNotInCache("c")

	// We should be able to add items to the cache, then retrieve them.

	a := "A"
	b := "B"

	cache.Set("a", a, 0)
	cache.Set("b", b, 0)

	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("b", b)

	// Should be able to keep reading the items; they shouldn't be evicted.

	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("b", b)

	// We should be able to clear items from the cache too.
	// Other items should remain in the cache.

	assert(t, cache.Clear("a"), "Clear should return true when item in cache")
	assertNot(t, cache.Clear("a"), "Clear should return false when item not in cache")

	cache.assertNotInCache("a")
	cache.assertCachedEquals("b", b)
}

func TestMemoryCache_Capacity(t *testing.T) {
	t.Parallel()

	// We should be able to limit the cache to a fixed memory capacity.
	// Adding items to the cache beyond the capacity should purge the
	// least recently read items.

	cache := TestMemoryCache{
		cache: &MemoryCache{
			MaxItems: 3,
		},
		ctx: context.Background(),
		t:   t,
	}

	// Add items up to the capacity. All items should be retained.

	a := "A"
	b := "B"
	c := "C"

	cache.Set("a", a, 0)
	cache.Set("b", b, 0)
	cache.Set("c", c, 0)

	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("b", b)
	cache.assertCachedEquals("c", c)

	// Access the items in some different order.

	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("b", b)
	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("b", b) // b is least recently read
	cache.assertCachedEquals("c", c)
	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("c", c) // c is most recently read

	// Add one more item to the cache. The least recently read item should be gone.

	d := "D"

	cache.Set("d", d, 0)

	cache.assertNotInCache("b")

	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("c", c)
	cache.assertCachedEquals("d", d)

	// Access the current items in some different order again.

	cache.assertCachedEquals("c", c) // c is least recently read
	cache.assertCachedEquals("d", d)
	cache.assertCachedEquals("a", a) // a is most recently read

	// Add *two* items to the cache. The *two* least recently read items should be gone.

	e := "E"

	cache.Set("b", b, 0) // b again
	cache.Set("e", e, 0)

	cache.assertNotInCache("c")
	cache.assertNotInCache("d")

	cache.assertCachedEquals("a", a)
	cache.assertCachedEquals("b", b)
	cache.assertCachedEquals("e", e)
}

func TestMemoryCache_Expiry(t *testing.T) {
	t.Parallel()

	t.Skip("TODO: Implement!")
}
