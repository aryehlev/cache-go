package s3fifo

import (
	"github.com/aryehlev/s3fifo/queues"
	"github.com/aryehlev/s3fifo/structures"

	"github.com/dolthub/maphash"
)

type S3Fifo[K comparable, V any] struct {
	hasher maphash.Hasher[K]

	main  queues.Main[V]
	small queues.Small[V]
	ghost queues.Ghost

	data map[uint64]*structures.Node[V]
}

func New[K comparable, V any](size int) S3Fifo[K, V] {
	smallSize := size / 10
	mainSize := size - smallSize

	return S3Fifo[K, V]{
		data:   make(map[uint64]*structures.Node[V]),
		main:   queues.NewMain[V](mainSize),
		small:  queues.NewSmall[V](smallSize),
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
		evicted, needEviction = sf.main.Put(node)
	} else {
		evicted, needEviction = sf.small.Put(node)
	}

	for needEviction && iterations < 5 {
		iterations++

		switch evicted.EvictionPlacement() {
		case structures.Small:
			evicted, needEviction = sf.small.Put(evicted)
		case structures.Main:
			evicted, needEviction = sf.main.Put(evicted)
		case structures.Ghost:
			delete(sf.data, evicted.Hash)
			sf.ghost.Put(evicted.Hash)
		case structures.None:
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
