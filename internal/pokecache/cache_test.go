package pokecache

import (
	"testing"
	"time"

	"internal/pokecache"
)

func TestCacheEviction(t *testing.T) {
	// Use a short interval to speed up the test
	interval := 1 * time.Second

	// Create a cache with the above interval
	cache := pokecache.NewCache(interval)

	// Add an item
	key := "testKey"
	value := []byte("testValue")
	cache.Add(key, value)

	// Immediately retrieve it; it should be present
	if got, ok := cache.Get(key); !ok || string(got) != "testValue" {
		t.Errorf("Expected %q, got %q (ok=%v)", value, got, ok)
	}

	// Wait a bit beyond 'interval' so the item should be evicted
	time.Sleep(interval + 500*time.Millisecond)

	// Try retrieving again; it should NOT be found
	if got, ok := cache.Get(key); ok {
		t.Errorf("Expected item to be evicted, but got %q (ok=%v)", string(got), ok)
	}
}

