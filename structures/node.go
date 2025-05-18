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
	Count             atomic.Int32
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

func (n *Node[V]) Small() {
	n.CurrentQueuePlcmt = Small
}

func (n *Node[V]) Main() {
	n.CurrentQueuePlcmt = Main
}

func (n *Node[V]) Hit() {
	if n.Count.Load() < maxCount {
		n.Count.Add(1)
	}
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
	count := n.Count.Load()
	if roomInMain || count > 0 {
		n.Count.Store(0)
		return Main
	}

	if count < 0 {
		n.Count.Store(0)
	}

	return Ghost
}

func (n *Node[V]) getFromMain() QueuePlcmt {
	if n.Count.Load() <= 0 {
		return None
	} else {
		n.Count.Add(-1)
		return Main
	}
}

func (n *Node[V]) getFromGhost() QueuePlcmt {
	return None
}
