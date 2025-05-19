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

func (s *NodeQueue[T]) Pop() (val T) {
	if s.First == nil {
		return
	}

	val = s.First.v

	if s.First == s.Last {
		s.Last = nil
		s.First = nil
		s.length--
		return
	}

	s.First = s.First.Next
	if s.First != nil {
		s.First.Prev = nil
	}

	s.length--
	return
}

func (s *NodeQueue[T]) Delete(node *Node[T]) {
	if s.First == nil { // queue is empty
		return
	}

	if node.Prev == nil { // the node is the first.
		s.First = s.First.Next
		s.length--
		return
	}

	if node.Next == nil { // the node is last.
		s.Last = s.Last.Prev
		s.Last.Next = nil

		s.length--
		return
	}

	next := node.Next
	node.Prev.Next = next
	node.Next.Prev = node.Prev
	s.length--
}

func (s *NodeQueue[T]) Len() int {
	return s.length
}
