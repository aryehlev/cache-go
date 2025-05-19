package queues

import "github.com/aryehlev/cache-go/structures"

type Ghost struct {
	data *structures.HashSet
}

func NewGhost(size uint) *Ghost {
	return &Ghost{
		data: structures.NewHashSet(size),
	}
}

func (g Ghost) Put(hash uint64) {
	g.data.Add(hash)
}

func (g Ghost) GetAndDel(key uint64) bool {
	if g.data.ContainsAndDelete(key) {
		return true
	}

	return false
}

func (g Ghost) Get(key uint64) bool {
	return g.data.Contains(key)
}
