package s3fifo

import (
	"github.com/aryehlev/s3fifo/queues"
	"github.com/aryehlev/s3fifo/structures"
	"github.com/dolthub/maphash"
)

type Cache[K comparable, V any] struct {
	hasher maphash.Hasher[K]

	main  queues.Main[V]
	small queues.Small[V]
	ghost *queues.Ghost

	data map[uint64]*structures.Node[V]

	cap int
}

func New[K comparable, V any](size int) Cache[K, V] {
	smallSize := size / 10
	mainSize := size - smallSize

	return Cache[K, V]{
		data:   make(map[uint64]*structures.Node[V]),
		main:   queues.NewMain[V](mainSize),
		small:  queues.NewSmall[V](smallSize),
		ghost:  queues.NewGhost(mainSize),
		hasher: maphash.NewHasher[K](),
		cap:    size,
	}
}

func (sf Cache[K, V]) Where(key K) structures.QueuePlcmt {
	hash := sf.hasher.Hash(key)
	if sf.ghost.Get(hash) {
		return structures.Ghost
	}

	if n, ok := sf.data[hash]; ok {
		return n.CurrentQueuePlcmt
	}

	return structures.None
}

func (sf Cache[K, V]) Set(key K, v V) {
	hash := sf.hasher.Hash(key)

	var evicted *structures.Node[V]
	var needEviction bool

	node := structures.NewNode(v, hash)
	sf.data[hash] = node

	if sf.ghost.GetAndDel(hash) {
		evicted, needEviction = sf.main.Put(node)
	} else {
		evicted, needEviction = sf.small.Put(node)
	}

	iterations := 0
	for needEviction {
		iterations++
		switch evicted.EvictionPlacement() {
		case structures.Small:
			evicted, needEviction = sf.small.Put(evicted)
		case structures.Main:
			evicted, needEviction = sf.main.Put(evicted)
			if needEviction && iterations > 20 {
				evicted.CurrentQueuePlcmt = structures.None
			}
		case structures.Ghost:
			sf.ghost.Put(evicted.Hash)
			fallthrough
		case structures.None:
			delete(sf.data, evicted.Hash)
			needEviction = false
		}
	}
}

func (sf Cache[K, V]) Get(key K) (v V, ok bool) {
	hash := sf.hasher.Hash(key)
	if node, okMap := sf.data[hash]; okMap {
		node.Hit()
		v, ok = node.GetVal(), true
		return
	}

	return
}

func (sf Cache[K, V]) Size() int {
	return len(sf.data)
}
