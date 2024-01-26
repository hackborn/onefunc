package pipeline

import (
	"fmt"
)

func RunExpr(expr string, input *RunInput) (*RunOutput, error) {
	p, err := Compile(expr)
	if err != nil {
		return nil, err
	}
	return Run(p, input)
}

func Run(p *Pipeline, input *RunInput) (*RunOutput, error) {
	active, err := prepareForRun(p, input)
	if err != nil {
		return nil, err
	}
	state := &State{}
	finalOutput := RunOutput{}
	for len(active) > 0 {
		var nextNodesMap map[*runningNode]struct{}
		for _, n := range active {
			output, err := n.node.Run(state, n.input)
			if err != nil {
				return nil, err
			}
			if len(n.output) > 0 {
				for _, topin := range n.output {
					topin.toNode.inputCount++
					if output != nil && len(output.Pins) > 0 {
						topin.toNode.input.Pins = append(topin.toNode.input.Pins, output.Pins...)
					}
					if topin.toNode.ready() {
						if nextNodesMap == nil {
							nextNodesMap = make(map[*runningNode]struct{})
						}
						nextNodesMap[topin.toNode] = struct{}{}
					}
				}
			} else if output != nil && len(output.Pins) > 0 {
				finalOutput.Pins = append(finalOutput.Pins, output.Pins...)
			}
		}
		active = make([]*runningNode, 0, len(nextNodesMap))
		for k, _ := range nextNodesMap {
			active = append(active, k)
		}
	}
	return &finalOutput, nil
}

func prepareForRun(p *Pipeline, input *RunInput) ([]*runningNode, error) {
	if len(p.roots) < 1 {
		return nil, fmt.Errorf("No roots")
	}
	active := make([]*runningNode, 0, len(p.roots))
	for _, n := range p.roots {
		active = append(active, n)
		n.input.Pins = nil
	}
	for _, n := range p.nodes {
		n.inputCount = 0
	}
	if input != nil && len(input.Pins) > 0 {
		for _, n := range active {
			if n.input.Pins == nil {
				n.input.Pins = make([]Pin, 0, len(input.Pins))
			}
			n.input.Pins = append(n.input.Pins, input.Pins...)
		}
	}
	return active, nil
}
