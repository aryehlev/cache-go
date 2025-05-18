package queues

import "github.com/aryehlev/s3fifo/structures"

type Main[V any] struct {
	queue *structures.NodeQueue[V]
}

func NewMain[V any](capacity int) Main[V] {
	return Main[V]{
		queue: structures.NewNodeQueue[V](capacity),
	}
}

func (m Main[V]) Put(newNode *structures.Node[V]) (evicted *structures.Node[V], wasEviction bool) {
	newNode.CurrentQueuePlcmt = structures.Main
	evicted, wasEviction = m.queue.PutNode(newNode)
	return
}

func (m Main[V]) Len() int {
	return m.queue.Len()
}

func (m Main[V]) Delete(node *structures.Node[V]) {
	m.queue.Delete(node)
}
