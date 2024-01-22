package pipeline

type Pipeline struct {
	roots []*runningNode
	nodes []*runningNode
}

type runningPin struct {
	inName string
	toNode *runningNode
}

type runningNode struct {
	node          Node
	inputCount    int
	maxInputCount int
	input         RunInput
	output        []*runningPin
}

func (n *runningNode) ready() bool {
	return n.inputCount >= n.maxInputCount
}
