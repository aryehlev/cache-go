package structures

type SimpleQueue[T any] struct {
	First    *Node[T]
	Last     *Node[T]
	capacity int
	length   int
}

func (s *SimpleQueue[T]) Put(val T, hash uint64) (evicted *Node[T], wasEviction bool) {
	newNode := &Node[T]{
		v:    val,
		Hash: hash,
	}

	return s.PutNode(newNode)
}

func (s *SimpleQueue[T]) PutNode(newNode *Node[T]) (evicted *Node[T], wasEviction bool) {
	if s.length == 0 {
		s.First = newNode
		s.Last = s.First

		return
	}

	newNode.Prev = s.Last
	s.Last = newNode

	if s.length >= s.capacity {
		evicted = s.First
		wasEviction = true

		s.First = s.First.Next
	}

	return
}

func (s *SimpleQueue[T]) Pop() T {
	val := s.First.v
	s.First = s.First.Next
	return val
}

func (s *SimpleQueue[T]) Delete(node *Node[T]) {
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
