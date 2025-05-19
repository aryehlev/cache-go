package queues

import "github.com/aryehlev/cache-go/structures"

type Small[V any] struct {
	queue *structures.NodeQueue[V]
}

func NewSmall[V any](capacity uint) Small[V] {
	return Small[V]{
		queue: structures.NewNodeQueue[V](capacity),
	}
}

func (s Small[V]) Put(newNode *structures.Node[V]) (evicted *structures.Node[V], wasEviction bool) {
	newNode.Small()
	evicted, wasEviction = s.queue.PutNode(newNode)
	return
}

func (s Small[V]) Delete(node *structures.Node[V]) {
	node.None()
	s.queue.Delete(node)
}
