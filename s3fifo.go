package s3fifo

import (
	"github.com/aryehlev/s3fifo/queues"
	"github.com/aryehlev/s3fifo/structures"
	"hash/maphash"
	"sync"
)

type Cache[K comparable, V any] struct {
	hasher maphash.Seed

	main  queues.Main[V]
	small queues.Small[V]
	ghost *queues.Ghost

	data  map[uint64]*structures.Node[V]
	mutex sync.RWMutex

	smallCap int
	mainCap  int
}

func New[K comparable, V any](size int) *Cache[K, V] {
	smallSize := size / 10
	mainSize := size - smallSize

	return &Cache[K, V]{
		data:     make(map[uint64]*structures.Node[V]),
		main:     queues.NewMain[V](mainSize),
		small:    queues.NewSmall[V](smallSize),
		ghost:    queues.NewGhost(mainSize),
		hasher:   maphash.MakeSeed(),
		smallCap: smallSize,
		mainCap:  mainSize,
	}
}

func (sf *Cache[K, V]) Where(key K) structures.QueuePlcmt {
	hash := maphash.Comparable(sf.hasher, key)
	if sf.ghost.Get(hash) {
		return structures.Ghost
	}

	if n, ok := sf.data[hash]; ok {
		return n.CurrentQueuePlcmt
	}

	return structures.None
}

func (sf *Cache[K, V]) Set(key K, v V) {
	hash := maphash.Comparable(sf.hasher, key)

	var evicted *structures.Node[V]
	var needEviction bool

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	if node, ok := sf.data[hash]; ok {
		node.SetVal(v)
		return
	}

	node := structures.NewNode(v, hash)
	sf.data[hash] = node

	if sf.ghost.GetAndDel(hash) {
		evicted, needEviction = sf.main.Put(node)
	} else {
		evicted, needEviction = sf.small.Put(node)
	}

	roomInMain := sf.main.Len() < sf.mainCap

	iterations := 0
	for needEviction {
		iterations++
		switch evicted.EvictionPlacement(roomInMain) {
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

func (sf *Cache[K, V]) Get(key K) (v V, ok bool) {
	hash := maphash.Comparable(sf.hasher, key)
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()

	if node, okMap := sf.data[hash]; okMap {
		node.Hit()
		v, ok = node.GetVal(), true
		return
	}

	return
}

func (sf *Cache[K, V]) Size() int {
	return len(sf.data)
}

func (sf *Cache[K, V]) Delete(key K) bool {
	hash := maphash.Comparable(sf.hasher, key)

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	inGhost := sf.ghost.GetAndDel(hash)
	if inGhost {
		return true
	}

	node, ok := sf.data[hash]
	if !ok {
		return false
	}

	delete(sf.data, hash)
	switch node.CurrentQueuePlcmt {
	case structures.Main:
		sf.main.Delete(node)
	case structures.Small:
		sf.small.Delete(node)
	}

	return true
}
