package pipeline

func Run(p *Pipeline, input *RunInput) (*RunOutput, error) {
	err := p.prepareForRun(input)
	if err != nil {
		return nil, err
	}
	state := &State{}
	finalOutput := RunOutput{}
	for len(p.active) > 0 {
		var nextNodesMap map[*runningNode]struct{}
		for _, n := range p.active {
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
		p.active = make([]*runningNode, 0, len(nextNodesMap))
		for k, _ := range nextNodesMap {
			p.active = append(p.active, k)
		}
	}
	return &finalOutput, nil
}
