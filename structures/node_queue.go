package structures

type NodeQueue[T any] struct {
	First    *Node[T]
	Last     *Node[T]
	capacity int
	length   int
}

func NewNodeQueue[T any](capacity int) *NodeQueue[T] {
	return &NodeQueue[T]{
		capacity: capacity,
	}
}

func (s *NodeQueue[T]) Put(val T, hash uint64) (evicted *Node[T], wasEviction bool) {
	newNode := &Node[T]{
		v:    val,
		Hash: hash,
	}

	return s.PutNode(newNode)
}

func (s *NodeQueue[T]) PutNode(newNode *Node[T]) (evicted *Node[T], wasEviction bool) {
	if s.length == 0 {
		s.First = newNode
		s.Last = newNode

		s.length++
		return
	}

	s.length++
	s.Last.Next = newNode
	newNode.Prev = s.Last
	s.Last = newNode

	if s.length > s.capacity {
		evicted = s.First
		wasEviction = true
		s.length--
		s.First = s.First.Next
	}

	return
}

func (s *NodeQueue[T]) Pop() T {
	val := s.First.v
	s.First = s.First.Next
	s.length--
	return val
}

func (s *NodeQueue[T]) Delete(node *Node[T]) {
	s.length--
	if node.Prev == nil {
		s.First = s.First.Next
		return
	}

	if node.Next == nil {
		s.Last = s.Last.Prev
		s.Last.Next = nil
	}

	node.Prev.Next = node.Next
}

func (s *NodeQueue[T]) Len() int {
	return s.length
}
