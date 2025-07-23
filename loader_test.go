package cache_go

import (
	"testing"
)

func TestCacheWithLoader(t *testing.T) {
	// Create a simple in-memory backend
	backend := map[string]int{
		"key1": 42,
		"key2": 84,
	}

	// Create a loader function that fetches from the backend
	loader := func(key string) (int, bool) {
		value, found := backend[key]
		return value, found
	}

	// Create a cache with the loader
	cache, err := NewWithLoader[string, int](100, loader)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Test getting a value that's not in the cache but is in the backend
	value, found := cache.Get("key1")
	if !found {
		t.Errorf("Expected to find key1 via loader, but it wasn't found")
	}
	if value != 42 {
		t.Errorf("Expected value 42 for key1, got %d", value)
	}

	// Test that the value is now cached
	// We can verify this by checking if the value is still accessible
	// even if we remove it from the backend
	delete(backend, "key1")
	value, found = cache.Get("key1")
	if !found {
		t.Errorf("Expected to find key1 in cache after loading, but it wasn't found")
	}
	if value != 42 {
		t.Errorf("Expected value 42 for key1 from cache, got %d", value)
	}

	// Test getting a value that's not in the cache or backend
	value, found = cache.Get("nonexistent")
	if found {
		t.Errorf("Expected not to find nonexistent key, but it was found with value %d", value)
	}

	// Test that the loader is not used when a value is already in the cache
	// We'll set a value directly in the cache
	cache.Set("key3", 123)
	
	// Then add a different value to the backend
	backend["key3"] = 456
	
	// When we get the value, it should come from the cache, not the backend
	value, found = cache.Get("key3")
	if !found {
		t.Errorf("Expected to find key3 in cache, but it wasn't found")
	}
	if value != 123 {
		t.Errorf("Expected value 123 for key3 from cache, got %d (should not have used loader)", value)
	}
}