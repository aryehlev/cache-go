package queues

import "github.com/aryehlev/s3fifo/structures"

type Ghost struct {
	queue *structures.SimpleQueue[uint64]

	data map[uint64]bool
}

func NewGhost(size int) Ghost {
	return Ghost{
		data:  make(map[uint64]bool),
		queue: structures.NewSimpleQueue[uint64](size),
	}
}

func (g Ghost) Put(hash uint64) {
	evictedFromGhost, wasEvictedGhost := g.queue.Enqueue(hash)
	g.data[hash] = true
	if wasEvictedGhost {
		delete(g.data, evictedFromGhost)
	}
}

func (g Ghost) GetAndDel(key uint64) bool {
	if g.data[key] {
		delete(g.data, key)
		return true
	}

	return false
}
