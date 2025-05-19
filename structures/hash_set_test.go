package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashSet_AddAndContainsAndDelete(t *testing.T) {
	hs := NewHashSet(10)
	values := []uint64{2, 12, 22}
	for _, v := range values {
		assert.False(t, hs.ContainsAndDelete(v), "ContainsAndDelete(%d) should be false before Add", v)

		hs.Add(v)
		assert.True(t, hs.ContainsAndDelete(v), "ContainsAndDelete(%d) should be true after Add", v)

		assert.False(t, hs.ContainsAndDelete(v), "ContainsAndDelete(%d) should be false after deletion", v)
	}
}

func TestHashSet_ContainsAndDelete_NonExistent(t *testing.T) {
	hs := NewHashSet(5)
	hs.Add(7)
	assert.False(t, hs.ContainsAndDelete(8), "ContainsAndDelete(8) should be false for non-existent element")

	assert.True(t, hs.ContainsAndDelete(7), "ContainsAndDelete(7) should be true after checking non-existent element")
}

func TestHashSet_CollisionBehavior(t *testing.T) {
	hs := NewHashSet(2)
	hs.Add(1)

	assert.False(t, hs.ContainsAndDelete(3), "ContainsAndDelete(3) should be false for colliding but absent element")

	assert.True(t, hs.ContainsAndDelete(1), "ContainsAndDelete(1) should be true after colliding check")
}

func TestHashSet_EmptySet(t *testing.T) {
	hs := NewHashSet(3)
	for _, v := range []uint64{0, 1, 2} {
		assert.False(t, hs.ContainsAndDelete(v), "ContainsAndDelete(%d) should be false on empty set", v)
	}
}

func TestHashSet_CapacityOne(t *testing.T) {
	hs := NewHashSet(1)
	hs.Add(5)
	assert.True(t, hs.ContainsAndDelete(5), "ContainsAndDelete(5) should be true after Add for capacity 1")

	assert.False(t, hs.ContainsAndDelete(6), "ContainsAndDelete(6) should be false for absent element in capacity 1")
}
