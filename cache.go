package onecache

import (
	"context"
	"time"
)

type Cache interface {
	// Get queries this cache for the value under the given key.
	// If found (& unexpired), returns the value -- which may be nil -- and true.
	// If not found (or expired), returns nil and false.
	Get(ctx context.Context, key string) (interface{}, bool, error)

	// Set caches the given value (which may be nil) under the given key,
	// optionally expiring after the given duration.
	// A duration of 0 means the value should never expire.
	Set(ctx context.Context, key string, val interface{}, expireAfter time.Duration) error

	// Clear clears the value, if any, cached under the given key.
	// Returns whether a value was cached (and thus cleared).
	Clear(ctx context.Context, key string) (bool, error)
}
