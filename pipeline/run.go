package pipeline

func Run(p *Pipeline, _input *RunInput) (*RunOutput, error) {
	err := p.prepareForRun()
	if err != nil {
		return nil, err
	}
	input := _input
	if input == nil {
		input = &RunInput{}
	}
	state := &State{}
	finalOutput := RunOutput{}
	for len(p.active) > 0 {
		var nextActive []*runningNode
		for _, n := range p.active {
			output, err := n.node.Run(state, *input)
			if err != nil {
				return nil, err
			}
			if output != nil && len(output.Pins) > 0 {
				finalOutput.Pins = append(finalOutput.Pins, output.Pins...)
			}
		}
		p.active = nextActive
	}
	return &finalOutput, nil
}
