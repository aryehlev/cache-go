package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimpleQueue(t *testing.T) {
	q := NewSimpleQueue[int](5)
	assert.NotNil(t, q)
	assert.Equal(t, 5, q.capacity)
	assert.Equal(t, 0, q.size)
}

func TestEnqueueDequeueBasic(t *testing.T) {
	q := NewSimpleQueue[int](3)
	evicted, wasEviction := q.Enqueue(1)
	assert.False(t, wasEviction)
	assert.Equal(t, 0, evicted)

	evicted, wasEviction = q.Enqueue(2)
	assert.False(t, wasEviction)
	assert.Equal(t, 0, evicted)

	value, ok := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	value, ok = q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	// Dequeue on empty queue
	value, ok = q.Dequeue()
	assert.False(t, ok)
	assert.Equal(t, 0, value)
}

func TestEviction(t *testing.T) {
	q := NewSimpleQueue[int](2)
	q.Enqueue(10)
	q.Enqueue(20)

	evicted, wasEviction := q.Enqueue(30)
	assert.True(t, wasEviction)
	assert.Equal(t, 10, evicted)

	value, ok := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 20, value)

	value, ok = q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 30, value)
}

func TestWrapAround(t *testing.T) {
	q := NewSimpleQueue[int](3)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	q.Dequeue()

	evicted, wasEviction := q.Enqueue(4)
	assert.False(t, wasEviction)
	assert.Equal(t, 0, evicted)

	expected := []int{2, 3, 4}
	for _, exp := range expected {
		value, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, exp, value)
	}
}

func TestZeroCapacity(t *testing.T) {
	q := NewSimpleQueue[int](0)
	evicted, wasEviction := q.Enqueue(1)
	assert.False(t, wasEviction)
	assert.Equal(t, 0, evicted)

	value, ok := q.Dequeue()
	assert.False(t, ok)
	assert.Equal(t, 0, value)
}

func TestGenericType(t *testing.T) {
	q := NewSimpleQueue[string](2)
	q.Enqueue("a")
	q.Enqueue("b")

	evicted, wasEviction := q.Enqueue("c")
	assert.True(t, wasEviction)
	assert.Equal(t, "a", evicted)

	value, ok := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, "b", value)

	value, ok = q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, "c", value)
}
