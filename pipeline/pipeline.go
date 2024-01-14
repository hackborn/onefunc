package pipeline

type Pipeline struct {
	//	srcNodes []srcNode
	nodes []*runningNode
	pins  []*pin
}

type pin struct {
	inName string
	node   Node
	ready  bool
}

/*
type srcNode struct {
	node   Node
	output []*pin
}
*/

type runningNode struct {
	node  Node
	input []*pin
}

func (n *runningNode) ready() bool {
	for _, pin := range n.input {
		if !pin.ready {
			return false
		}
	}
	return true
}

// OK how does this work -- we need a notion of pins connecting nodes, and to know when a node has had all pins tickled so it can run
