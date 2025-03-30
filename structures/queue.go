package structures

type SimpleQueue[T any] struct {
	data                 []T
	start, end, capacity int
}

func (s *SimpleQueue[T]) Put(val T) {
	s.data[s.end] = val
	s.end = (s.end + 1) % s.capacity
}

func (s *SimpleQueue[T]) Pop() T {
	toReturn := s.data[s.start]
	s.start = (s.start + 1) % s.capacity

	return toReturn
}

func (s *SimpleQueue[T]) IsFull() bool {
	return (s.end+1)%s.capacity == s.start
}
