package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashSet(t *testing.T) {
	capacity := uint(10)
	hs := NewHashSet(capacity)
	assert.NotNil(t, hs)
	assert.Equal(t, capacity, uint(len(hs.buckets)))
	assert.Equal(t, capacity, uint(len(hs.used)))
}

func TestHashSetAddAndContains(t *testing.T) {
	hs := NewHashSet(5)

	value := uint64(42)
	hs.Add(value)
	assert.True(t, hs.Contains(value))

	assert.False(t, hs.Contains(99))
}

func TestHashSetDelete(t *testing.T) {
	hs := NewHashSet(5)

	value := uint64(42)
	hs.Add(value)
	assert.True(t, hs.Contains(value))

	hs.Delete(value)
	assert.False(t, hs.Contains(value))

	hs.Delete(99)
}

func TestHashSetCollision(t *testing.T) {
	hs := NewHashSet(2)

	value1 := uint64(1)
	value2 := uint64(3)

	hs.Add(value1)
	assert.True(t, hs.Contains(value1))

	hs.Add(value2)
	assert.True(t, hs.Contains(value2))
	assert.False(t, hs.Contains(value1))
}

func TestHashSetOverwrite(t *testing.T) {
	hs := NewHashSet(5)

	// Test overwriting a value
	value := uint64(42)
	hs.Add(value)
	hs.Add(value) // Add same value again
	assert.True(t, hs.Contains(value))
}

func TestHashSetZeroValue(t *testing.T) {
	hs := NewHashSet(5)

	// Test adding and checking zero value
	hs.Add(0)
	assert.True(t, hs.Contains(0))

	hs.Delete(0)
	assert.False(t, hs.Contains(0))
}

func TestHashSetLargeValues(t *testing.T) {
	hs := NewHashSet(5)

	largeValue := uint64(1<<63 - 1)
	hs.Add(largeValue)
	assert.True(t, hs.Contains(largeValue))

	hs.Delete(largeValue)
	assert.False(t, hs.Contains(largeValue))
}

func TestHashSetMultipleOperations(t *testing.T) {
	hs := NewHashSet(10)

	values := []uint64{1, 2, 3, 4, 5}

	for _, v := range values {
		hs.Add(v)
	}

	for _, v := range values {
		assert.True(t, hs.Contains(v))
	}

	hs.Delete(2)
	hs.Delete(4)

	assert.True(t, hs.Contains(1))
	assert.False(t, hs.Contains(2))
	assert.True(t, hs.Contains(3))
	assert.False(t, hs.Contains(4))
	assert.True(t, hs.Contains(5))
}
