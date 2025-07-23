package cache_go

// LoaderFunc is a function type that loads a value for a key.
// It returns the loaded value and a boolean indicating whether the value was found.
type LoaderFunc[K comparable, V any] func(key K) (V, bool)

// CacheWithLoader extends the Cache with a loader function that can fetch values from a backend.
type CacheWithLoader[K comparable, V any] struct {
	*Cache[K, V]
	loader LoaderFunc[K, V]
}

// NewWithLoader creates a new cache with the specified size and loader function.
func NewWithLoader[K comparable, V any](size uint, loader LoaderFunc[K, V]) (*CacheWithLoader[K, V], error) {
	cache, err := New[K, V](size)
	if err != nil {
		return nil, err
	}

	return &CacheWithLoader[K, V]{
		Cache:  cache,
		loader: loader,
	}, nil
}

// Get retrieves a value from the cache. If the value is not in the cache,
// it attempts to load it using the loader function.
func (cl *CacheWithLoader[K, V]) Get(key K) (V, bool) {
	// Try to get from cache first
	value, found := cl.Cache.Get(key)
	if found {
		return value, true
	}

	// If not in cache, try to load from backend
	if cl.loader != nil {
		value, found := cl.loader(key)
		if found {
			// Store the loaded value in the cache
			cl.Cache.Set(key, value)
			return value, true
		}
	}

	// Return zero value if not found
	var zero V
	return zero, false
}