package queues

import "github.com/aryehlev/cache-go/structures"

type Main[V any] struct {
	queue *structures.NodeQueue[V]
}

func NewMain[V any](capacity uint) Main[V] {
	return Main[V]{
		queue: structures.NewNodeQueue[V](capacity),
	}
}

func (m Main[V]) Put(newNode *structures.Node[V]) (evicted *structures.Node[V], wasEviction bool) {
	newNode.Main()
	evicted, wasEviction = m.queue.PutNode(newNode)
	return
}

func (m Main[V]) Len() uint {
	return m.queue.Len()
}

func (m Main[V]) Delete(node *structures.Node[V]) {
	node.None()
	m.queue.Delete(node)
}
