package s3fifo

import (
	"github.com/aryehlev/s3fifo/structures"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicSetAndGet(t *testing.T) {
	cache := New[string, int](2)

	cache.Set("a", 1)
	cache.Set("b", 2)

	val, ok := cache.Get("a")
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	val, ok = cache.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	_, ok = cache.Get("c")
	assert.False(t, ok)
}

func TestOverwriteValue(t *testing.T) {
	cache := New[string, int](2)

	cache.Set("a", 1)
	cache.Set("a", 2)

	val, ok := cache.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
}

func TestEvictionAndGhostInsertion(t *testing.T) {
	cache := New[string, int](3)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)
	cache.Set("d", 4) // should cause an eviction

	_, ok := cache.Get("a")
	_, ok2 := cache.Get("b")
	_, ok3 := cache.Get("c")
	_, ok4 := cache.Get("d")
	assert.True(t, ok4)
	assert.False(t, ok2)
	assert.False(t, ok3)
	assert.False(t, ok)
	// Only 3 items can live at once.
	//assert.Equal(t, 3, boolToInt(ok)+boolToInt(ok2)+boolToInt(ok3)+boolToInt(ok4))
}

func TestGhostPromotion(t *testing.T) {
	cache := New[string, int](2)

	cache.Set("a", 1) // put in small
	cache.Set("b", 2) // put in small

	assert.Equal(t, structures.Ghost, cache.Where("a"))
	// Suppose "a" was evicted.
	_, ok := cache.Get("a") // is in ghost
	assert.False(t, ok)

	cache.Set("a", 100) // reinsert "a", should promote to main
	//cache.Set("b", 100) // should still be in small

	assert.Equal(t, structures.Main, cache.Where("a"))
	assert.Equal(t, structures.Small, cache.Where("b"))

	cache.Set("c", 1000)
	val, ok := cache.Get("a") // should still be in main

	assert.Equal(t, structures.Main, cache.Where("a"))
	assert.Equal(t, structures.Small, cache.Where("c"))
	//assert.Equal(t, structures.None, cache.Where("b"))

	assert.True(t, ok)
	assert.Equal(t, 100, val)

	val, ok = cache.Get("b") // should be evicted

	assert.False(t, ok)
}

func TestIterationCap(t *testing.T) {
	cache := New[string, int](2)

	// Insert many elements to hit iterations cap (eviction loop max 5)
	for i := 0; i < 10; i++ {
		cache.Set(strconv.Itoa(i), i)
	}

	count := 0
	for i := 0; i < 10; i++ {
		_, ok := cache.Get(strconv.Itoa(i))
		if ok {
			count++
		}
	}

	// Should still maintain size limits roughly (can't have >1 active)
	assert.LessOrEqual(t, count, 2)
}

func TestEmptyCache(t *testing.T) {
	cache := New[string, int](2)

	cache.Set("a", 1)
	val, ok := cache.Get("a")
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestRepeatedInsertions(t *testing.T) {
	cache := New[string, int](2)

	for i := 0; i < 100; i++ {
		cache.Set("k", i)
		val, ok := cache.Get("k")
		assert.True(t, ok)
		assert.Equal(t, i, val)
	}
}

// Helper
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
