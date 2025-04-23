package metadata

import (
	"s3fifo/queues"
	"s3fifo/structures"

	"github.com/dolthub/maphash"
)

type S3Fifo[K comparable, V any] struct {
	main   structures.NodeQueue[V]
	small  structures.NodeQueue[V]
	hasher maphash.Hasher[K]

	ghost queues.Ghost[V]

	data map[K]*structures.Node[V]
}

func (c S3Fifo[K, V]) Set(key K, v V) {
	iterations := 0

	evicted, needEviction := c.setInitial(key, v)
	for needEviction && iterations < 5 {
		iterations++
		switch evicted.EvictedPlcmt() {
		case structures.Small:
			evicted, needEviction = c.small.PutNode(evicted)
		case structures.Main:
			evicted, needEviction = c.main.PutNode(evicted)
		case structures.Ghost:
			c.ghost.Add(evicted)
		case structures.None:
			return
		}
	}
}

func (c S3Fifo[K, V]) Get(key K) {

}

func (c S3Fifo[K, V]) setInitial(key K, v V) (*structures.Node[V], bool) {
	hash := c.hasher.Hash(key)

	var evicted *structures.Node[V]
	var needEviction bool

	if c.ghost.Exists(hash) {
		evicted, needEviction = c.main.Put(v, hash)
	} else {
		evicted, needEviction = c.small.Put(v, hash)
	}

	return evicted, needEviction
}
