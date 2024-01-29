package pipeline

import (
	"fmt"

	"github.com/hackborn/onefunc/assign"
)

func RunExpr(expr string, input *RunInput, env map[string]any) (*RunOutput, error) {
	p, err := Compile(expr)
	if err != nil {
		return nil, err
	}
	return Run(p, input, env)
}

func Run(p *Pipeline, input *RunInput, env map[string]any) (*RunOutput, error) {
	active, err := prepareForRun(p, input, env)
	if err != nil {
		return nil, err
	}
	state := &State{}
	flushState := &State{Flush: true}
	finalOutput := RunOutput{}
	for len(active) > 0 {
		var nextNodesMap map[*runningNode]struct{}
		for _, n := range active {
			output, err := n.node.Run(state, n.input)
			if err != nil {
				return nil, err
			}
			// This node is done, flush it.
			flushOutput, _ := n.node.Run(flushState, RunInput{})
			if flushOutput != nil && len(flushOutput.Pins) > 0 {
				if output == nil || len(output.Pins) < 1 {
					output = flushOutput
				} else {
					output.Pins = append(output.Pins, flushOutput.Pins...)
				}
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

func prepareForRun(p *Pipeline, input *RunInput, env map[string]any) ([]*runningNode, error) {
	if len(p.roots) < 1 {
		return nil, fmt.Errorf("No roots")
	}
	active := make([]*runningNode, 0, len(p.roots))
	for _, n := range p.roots {
		active = append(active, n)
		n.input.Pins = nil
	}
	for _, n := range p.nodes {
		err := prepareNodeForRun(n, env)
		if err != nil {
			return nil, err
		}
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

func prepareNodeForRun(rn *runningNode, env map[string]any) error {
	// Reset input accumulator
	rn.inputCount = 0

	// Apply env vars
	envlen := len(rn.envVars)
	if envlen > 0 {
		req := assign.ValuesRequest{}
		req.FieldNames = make([]string, envlen, envlen)
		req.NewValues = make([]any, envlen, envlen)
		i := -1
		for k, v := range rn.envVars {
			i++
			req.FieldNames[i] = k
			if vv, ok := mapAt(v, env); ok {
				req.NewValues[i] = vv
			} else {
				return fmt.Errorf("Missing environment variable \"%v\"", v)
			}
		}
		return assign.Values(req, rn.node)
	}
	return nil
}

func mapAt(key string, m map[string]any) (any, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m[key]
	return v, ok
}
