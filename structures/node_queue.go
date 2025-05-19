package structures

type NodeQueue[T any] struct {
	First    *Node[T]
	Last     *Node[T]
	capacity int
	length   int
}

func NewNodeQueue[T any](capacity int) *NodeQueue[T] {
	return &NodeQueue[T]{capacity: capacity}
}

func (q *NodeQueue[T]) PutNode(newNode *Node[T]) (evicted *Node[T], wasEviction bool) {
	if q.length == 0 {
		q.First, q.Last = newNode, newNode
		q.length = 1
		return
	}

	q.Last.Next = newNode
	newNode.Prev = q.Last
	q.Last = newNode
	q.length++

	if q.length > q.capacity {
		wasEviction = true
		evicted = q.First

		q.First = evicted.Next
		q.First.Prev = nil
		q.length--

		evicted.Next, evicted.Prev = nil, nil
	}

	return evicted, wasEviction
}

func (q *NodeQueue[T]) Pop() (value T) {
	if q.length == 0 {
		return
	}

	node := q.First
	value = node.v

	if q.length == 1 {
		q.First, q.Last = nil, nil
	} else {
		q.First = node.Next
		q.First.Prev = nil
	}
	q.length--

	node.Next, node.Prev = nil, nil
	return value
}

func (q *NodeQueue[T]) Delete(node *Node[T]) {
	if node == nil || q.length == 0 {
		return
	}

	switch {
	case node == q.First && node == q.Last:
		q.First, q.Last = nil, nil

	case node == q.First:
		q.First = node.Next
		q.First.Prev = nil

	case node == q.Last:
		q.Last = node.Prev
		q.Last.Next = nil

	default:
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}

	q.length--
	node.Next, node.Prev = nil, nil
}

func (q *NodeQueue[T]) Len() int {
	return q.length
}
