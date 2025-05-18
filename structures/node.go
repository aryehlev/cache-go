package structures

import "sync/atomic"

const (
	None  QueuePlcmt = "none"
	Main  QueuePlcmt = "main"
	Ghost QueuePlcmt = "ghost"
	Small QueuePlcmt = "small"

	maxCount = 3
)

type QueuePlcmt string

type Node[V any] struct {
	v                 V
	count             atomic.Int32
	Hash              uint64
	CurrentQueuePlcmt QueuePlcmt
	Next              *Node[V]
	Prev              *Node[V]
}

func NewNode[V any](v V, hash uint64) *Node[V] {
	return &Node[V]{
		v:    v,
		Hash: hash,
	}
}

func (n *Node[V]) Hit() {
	if n.count.Load() < maxCount {
		n.count.Add(1)
	}
}

func (n *Node[V]) ResetCount() {
	n.count.Store(0)
}

func (n *Node[V]) SetVal(value V) {
	n.v = value
}

func (n *Node[V]) GetVal() V {
	return n.v
}

func (n *Node[V]) EvictionPlacement(roomInMain bool) QueuePlcmt {
	switch n.CurrentQueuePlcmt {
	case Small:
		return n.getFromSmall(roomInMain)
	case Ghost:
		return n.getFromGhost()
	case Main:
		return n.getFromMain()
	default:
		return None
	}
}

func (n *Node[V]) getFromSmall(roomInMain bool) QueuePlcmt {
	count := n.count.Load()
	if roomInMain || count > 0 {
		n.ResetCount()
		return Main
	}

	if count < 0 {
		n.ResetCount()
	}

	return Ghost
}

func (n *Node[V]) getFromMain() QueuePlcmt {
	if n.count.Load() <= 0 {
		return None
	} else {
		n.count.Add(-1)
		return Main
	}
}

func (n *Node[V]) getFromGhost() QueuePlcmt {
	return None
}
