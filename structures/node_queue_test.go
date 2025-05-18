package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNodeQueue(t *testing.T) {
	q := NewNodeQueue[int](3)
	assert.NotNil(t, q)
	assert.Equal(t, 3, q.capacity)
	assert.Equal(t, 0, q.Len())
}

func TestPutNodeBasic(t *testing.T) {
	q := NewNodeQueue[int](2)
	n1 := &Node[int]{v: 1}
	evicted, wasEviction := q.PutNode(n1)
	assert.False(t, wasEviction)
	assert.Nil(t, evicted)
	assert.Equal(t, 1, q.Len())

	val := q.Pop()
	assert.Equal(t, 1, val)
	assert.Equal(t, 0, q.Len())
}

func TestPutNodeEviction(t *testing.T) {
	q := NewNodeQueue[int](2)
	n1 := &Node[int]{v: 1}
	n2 := &Node[int]{v: 2}
	n3 := &Node[int]{v: 3}
	q.PutNode(n1)
	q.PutNode(n2)
	evicted, wasEviction := q.PutNode(n3)
	assert.True(t, wasEviction)
	assert.Equal(t, n1, evicted)
	assert.Equal(t, 2, q.Len())

	first := q.Pop()
	second := q.Pop()
	assert.Equal(t, 2, first)
	assert.Equal(t, 3, second)
	assert.Equal(t, 0, q.Len())
}

func TestPopMultiple(t *testing.T) {
	q := NewNodeQueue[int](3)
	n1 := &Node[int]{v: 1}
	n2 := &Node[int]{v: 2}
	q.PutNode(n1)
	q.PutNode(n2)

	first := q.Pop()
	second := q.Pop()
	assert.Equal(t, 1, first)
	assert.Equal(t, 2, second)
	assert.Equal(t, 0, q.Len())
}

func TestDeleteHead(t *testing.T) {
	q := NewNodeQueue[int](3)
	n1 := &Node[int]{v: 1}
	n2 := &Node[int]{v: 2}
	q.PutNode(n1)
	q.PutNode(n2)
	q.Delete(n1)

	val := q.Pop()
	assert.Equal(t, 2, val)
	assert.Equal(t, 0, q.Len())
}

func TestDeleteTail(t *testing.T) {
	q := NewNodeQueue[int](3)
	n1 := &Node[int]{v: 1}
	n2 := &Node[int]{v: 2}
	q.PutNode(n1)
	q.PutNode(n2)
	q.Delete(n2)

	val := q.Pop()
	assert.Equal(t, 1, val)
	assert.Equal(t, 0, q.Len())
}

func TestDeleteMiddle(t *testing.T) {
	q := NewNodeQueue[int](3)
	n1 := &Node[int]{v: 1}
	n2 := &Node[int]{v: 2}
	n3 := &Node[int]{v: 3}
	q.PutNode(n1)
	q.PutNode(n2)
	q.PutNode(n3)
	q.Delete(n2)

	first := q.Pop()
	second := q.Pop()
	assert.Equal(t, 1, first)
	assert.Equal(t, 3, second)
	assert.Equal(t, 0, q.Len())
}

func TestLen(t *testing.T) {
	q := NewNodeQueue[int](5)
	assert.Equal(t, 0, q.Len())
	q.PutNode(&Node[int]{v: 42})
	assert.Equal(t, 1, q.Len())
}

func TestDeleteOnlyThenDeleteAgain(t *testing.T) {
	q := NewNodeQueue[int](1)
	n1 := &Node[int]{v: 99}
	q.PutNode(n1)
	q.Delete(n1)
	// queue is now empty
	assert.Equal(t, 0, q.Len())

	q.Delete(n1)
	assert.Equal(t, 0, q.Len())
}

func TestDeleteOnlyThenPop(t *testing.T) {
	q := NewNodeQueue[int](1)
	n1 := &Node[int]{v: 100}
	q.PutNode(n1)
	q.Delete(n1)

	assert.NotPanics(t, func() {
		val := q.Pop()
		assert.Equal(t, 0, val)
	})
	assert.Equal(t, 0, q.Len())
}
