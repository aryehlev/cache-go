package structures

const (
	None  = "none"
	Main  = "main"
	Ghost = "ghost"
	Small = "small"
)

type QueuePlcmt string

type Node[V any] struct {
	v                 V
	Count             uint8
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
	n.Count++
}

func (n *Node[V]) SetVal(value V) {
	n.v = value
}

func (n *Node[V]) GetVal() V {
	return n.v
}

func (n *Node[V]) EvictionPlacement() QueuePlcmt {
	switch n.CurrentQueuePlcmt {
	case Small:
		return n.getFromSmall()
	case Ghost:
		return n.getFromGhost()
	case Main:
		return n.getFromMain()
	default:
		return None
	}
}

func (n *Node[V]) getFromSmall() QueuePlcmt {
	if n.Count == 0 {
		return Ghost
	} else {
		n.Count--
		return Main
	}
}

func (n *Node[V]) getFromMain() QueuePlcmt {
	if n.Count == 0 {
		return None
	} else {
		n.Count--
		return Main
	}
}

func (n *Node[V]) getFromGhost() QueuePlcmt {
	return None
}
