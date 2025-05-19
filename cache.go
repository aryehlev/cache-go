package cache_go

import (
	"errors"
	"hash/maphash"
	"sync"

	"github.com/aryehlev/cache-go/queues"
	"github.com/aryehlev/cache-go/structures"
)

type Cache[K comparable, V any] struct {
	hasher maphash.Seed

	main  queues.Main[V]
	small queues.Small[V]
	ghost *queues.Ghost

	data  map[uint64]*structures.Node[V]
	mutex sync.RWMutex

	smallCap uint
	mainCap  uint
}

func New[K comparable, V any](size uint) (*Cache[K, V], error) {
	smallSize := size / 10
	if smallSize < 1 {
		return nil, errors.New("size must be larger than 10")
	}
	mainSize := size - smallSize

	return &Cache[K, V]{
		data:     make(map[uint64]*structures.Node[V]),
		main:     queues.NewMain[V](mainSize),
		small:    queues.NewSmall[V](smallSize),
		ghost:    queues.NewGhost(mainSize),
		hasher:   maphash.MakeSeed(),
		smallCap: smallSize,
		mainCap:  mainSize,
	}, nil
}

func (sf *Cache[K, V]) Where(key K) structures.QueuePlcmt {
	hash := maphash.Comparable(sf.hasher, key)
	if sf.ghost.Get(hash) {
		return structures.Ghost
	}

	sf.mutex.RLock()
	defer sf.mutex.RUnlock()

	if n, ok := sf.data[hash]; ok {
		return n.CurrentPlcmt()
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
		switch evicted.NextPlacement(roomInMain) {
		case structures.Small:
			evicted, needEviction = sf.small.Put(evicted)
		case structures.Main:
			evicted, needEviction = sf.main.Put(evicted)
			if needEviction && iterations > 20 {
				evicted.None()
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
	sf.mutex.RLock()
	defer sf.mutex.RUnlock()
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

	switch node.CurrentPlcmt() {
	case structures.Main:
		sf.main.Delete(node)
	case structures.Small:
		sf.small.Delete(node)
	}

	return true
}

func (sf *Cache[K, V]) Clear() {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	initialMapCapacity := sf.smallCap + sf.mainCap + ((sf.smallCap + sf.mainCap) / 4)
	sf.data = make(map[uint64]*structures.Node[V], initialMapCapacity)

	sf.main = queues.NewMain[V](sf.mainCap)
	sf.small = queues.NewSmall[V](sf.smallCap)
	sf.ghost = queues.NewGhost(sf.mainCap)
}
