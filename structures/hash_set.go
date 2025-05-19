package structures

type HashSet struct {
	buckets []uint64
	used    []bool
}

func NewHashSet(capacity uint) *HashSet {
	return &HashSet{
		buckets: make([]uint64, capacity),
		used:    make([]bool, capacity),
	}
}

func (s *HashSet) Add(h uint64) {
	idx := int(h % uint64(len(s.buckets)))
	s.buckets[idx] = h
	s.used[idx] = true
}

func (s *HashSet) Contains(h uint64) bool {
	idx := int(h % uint64(len(s.buckets)))
	return s.used[idx] && s.buckets[idx] == h
}

func (s *HashSet) Delete(h uint64) {
	idx := int(h % uint64(len(s.buckets)))
	if s.used[idx] && s.buckets[idx] == h {
		s.used[idx] = false
		s.buckets[idx] = 0
	}
}
