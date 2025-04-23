package queues

import "s3fifo/structures"

type Ghost[V any] struct {
	queue structures.SimpleQueue[uint64]

	data map[uint64]bool
}

func (g *Ghost[V]) Add(hash uint64]) {
	evictedFromGhost, wasEvictedGhost := g.queue.Enqueue(hash)

	if wasEvictedGhost {
		delete(g.data, evictedFromGhost.Hash)
	}
}

func (g *Ghost[V]) Exists(key uint64) bool {
	if g.data[key] {
		delete(g.data, key)
		return true
	}

	return false
}
