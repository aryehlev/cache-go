package structures

const (
	Main QueuePlcmt = iota
	Ghost
	Small
	None
)

type QueuePlcmt int8

type Node[V any] struct {
	v                 V
	Count             uint8
	Hash              uint64
	CurrentQueuePlcmt QueuePlcmt
	Next              *Node[V]
	Prev              *Node[V]
}

func (n *Node[V]) Hit() {
	n.Count++
}

func (n *Node[V]) EvictedPlcmt() QueuePlcmt {
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
