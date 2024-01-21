package pipeline

import (
	"fmt"
)

type Pipeline struct {
	active []*runningNode

	roots []*runningNode
	pins  []*runningPin
	nodes []*runningNode
}

func (p *Pipeline) prepareForRun() error {
	if len(p.roots) < 1 {
		return fmt.Errorf("No roots")
	}
	p.active = make([]*runningNode, 0, len(p.roots))
	for _, n := range p.roots {
		p.active = append(p.active, n)
		n.input = nil
	}
	for _, n := range p.nodes {
		n.inputCount = 0
	}
	for _, p := range p.pins {
		p.ready = false
	}
	return nil
}

type runningPin struct {
	inName string
	toNode *runningNode
	ready  bool
}

type runningNode struct {
	node          Node
	inputCount    int
	maxInputCount int
	input         []*runningPin
	output        []*runningPin
}

func (n *runningNode) ready() bool {
	return n.inputCount >= n.maxInputCount
	for _, pin := range n.input {
		if !pin.ready {
			return false
		}
	}
	return true
}
