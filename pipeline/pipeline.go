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

func (p *Pipeline) prepareForRun(input *RunInput) error {
	if len(p.roots) < 1 {
		return fmt.Errorf("No roots")
	}
	p.active = make([]*runningNode, 0, len(p.roots))
	for _, n := range p.roots {
		p.active = append(p.active, n)
		n.input.Pins = nil
	}
	for _, n := range p.nodes {
		n.inputCount = 0
	}
	if input != nil && len(input.Pins) > 0 {
		for _, n := range p.active {
			if n.input.Pins == nil {
				n.input.Pins = make([]PinData, 0, len(input.Pins))
			}
			n.input.Pins = append(n.input.Pins, input.Pins...)
		}
	}
	return nil
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
