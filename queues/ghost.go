package queues

import "s3fifo/structures"

type Ghost[V any] struct {
	queue structures.SimpleQueue[uint64]

	data map[uint64]bool
}

func (g *Ghost[V]) Add(evicted *structures.Node[V]) {
	evictedFromGhost, wasEvictedGhost := g.queue.Put(evicted.Hash, evicted.Hash)

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
