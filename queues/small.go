package queues

import "github.com/aryehlev/s3fifo/structures"

type Small[V any] struct {
	queue *structures.NodeQueue[V]
}

func NewSmall[V any](capacity int) Small[V] {
	return Small[V]{
		queue: structures.NewNodeQueue[V](capacity),
	}
}

func (s Small[V]) Put(newNode *structures.Node[V]) (evicted *structures.Node[V], wasEviction bool) {
	newNode.CurrentQueuePlcmt = structures.Small
	evicted, wasEviction = s.queue.PutNode(newNode)
	return
}

func (s Small[V]) Delete(node *structures.Node[V]) {
	node.CurrentQueuePlcmt = structures.None
	s.queue.Delete(node)
}
