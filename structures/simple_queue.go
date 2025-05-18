package structures

type SimpleQueue[V any] struct {
	data       []V
	capacity   int
	head, tail int
	size       int
}

func NewSimpleQueue[V any](capacity int) *SimpleQueue[V] {
	return &SimpleQueue[V]{
		data:     make([]V, capacity),
		capacity: capacity,
		head:     0,
		tail:     0,
		size:     0,
	}
}
func (sq *SimpleQueue[V]) Enqueue(value V) (evicted V, wasEviction bool) {
	if sq.size == sq.capacity {
		evicted = sq.data[sq.tail]
		wasEviction = true
		sq.data[sq.tail] = value
		sq.tail = (sq.tail + 1) % sq.capacity
		return evicted, true
	}

	sq.data[sq.head] = value
	sq.head = (sq.head + 1) % sq.capacity
	sq.size++
	return
}

func (sq *SimpleQueue[V]) Dequeue() (value V, ok bool) {
	if sq.size == 0 {
		return
	}
	value = sq.data[sq.tail]
	sq.tail = (sq.tail + 1) % sq.capacity
	sq.size--

	return value, true
}
