package s3fifo

import (
	"s3fifo/queues"
	"s3fifo/structures"

	"github.com/dolthub/maphash"
)

type S3Fifo[K comparable, V any] struct {
	main   *structures.NodeQueue[V]
	small  *structures.NodeQueue[V]
	hasher maphash.Hasher[K]

	ghost *queues.Ghost

	data map[uint64]*structures.Node[V]
}

func New[K comparable, V any](size int) S3Fifo[K, V] {
	smallSize := size / 10
	mainSize := size - smallSize

	return S3Fifo[K, V]{
		data:   make(map[uint64]*structures.Node[V]),
		main:   structures.NewNodeQueue[V](mainSize),
		small:  structures.NewNodeQueue[V](smallSize),
		ghost:  queues.NewGhost(size),
		hasher: maphash.NewHasher[K](),
	}
}

func (sf S3Fifo[K, V]) Set(key K, v V) {
	hash := sf.hasher.Hash(key)
	if node, ok := sf.data[hash]; ok {
		node.Hit()
		node.SetVal(v)
		return
	}

	iterations := 0

	var evicted *structures.Node[V]
	var needEviction bool

	node := structures.NewNode(v, hash)
	sf.data[hash] = node

	if sf.ghost.GetAndDel(hash) {
		node.CurrentQueuePlcmt = structures.Main
		evicted, needEviction = sf.main.PutNode(node)
	} else {
		node.CurrentQueuePlcmt = structures.Small
		evicted, needEviction = sf.small.PutNode(node)
	}

	for needEviction && iterations < 5 {
		iterations++

		switch evicted.EvictedPlcmt() {
		case structures.Small:
			evicted.CurrentQueuePlcmt = structures.Small
			evicted, needEviction = sf.small.PutNode(evicted)
		case structures.Main:
			evicted.CurrentQueuePlcmt = structures.Main
			evicted, needEviction = sf.main.PutNode(evicted)
		case structures.Ghost:
			evicted.CurrentQueuePlcmt = structures.Ghost
			delete(sf.data, evicted.Hash)
			sf.ghost.Put(evicted.Hash)
		case structures.None:
			evicted.CurrentQueuePlcmt = structures.None
			delete(sf.data, evicted.Hash)
			break
		}

	}
}

func (sf S3Fifo[K, V]) Get(key K) (v V, ok bool) {
	hash := sf.hasher.Hash(key)
	if node, okMap := sf.data[hash]; okMap {
		node.Hit()
		v, ok = node.GetVal(), true
		return
	}

	return
}
