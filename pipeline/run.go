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
	build := newBuildRun(p.nodes)
	running, err := build.buildPipeline(p, input, env)
	build = nil
	if err != nil {
		return nil, err
	}
	//	fmt.Println("pipeline running, roots", len(p.roots), "active", len(active))
	state := &State{}
	finalOutput := RunOutput{}
	for len(running) > 0 {
		var nextNodesMap map[*runningNode]struct{}
		for _, rn := range running {
			state.NodeData = rn.nodeData
			output, err := rn.cn.node.Run(state, rn.input)
			if err != nil {
				return nil, err
			}
			// This node is done processing, flush it.
			output, err = flush(state, rn.cn.flusher, output)
			if err != nil {
				return nil, err
			}

			if len(rn.output) > 0 {
				for _, topin := range rn.output {
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
		running = make([]*runningNode, 0, len(nextNodesMap))
		for k, _ := range nextNodesMap {
			running = append(running, k)
		}
	}
	return &finalOutput, nil
}

func flush(state *State, flusher Flusher, output *RunOutput) (*RunOutput, error) {
	if flusher == nil {
		return output, nil
	}
	flushOutput, err := flusher.Flush(state)
	if err != nil {
		return output, err
	}
	if flushOutput != nil && len(flushOutput.Pins) > 0 {
		if output == nil || len(output.Pins) < 1 {
			output = flushOutput
		} else {
			output.Pins = append(output.Pins, flushOutput.Pins...)
		}
	}
	return output, nil
}

// buildRun takes the compiled nodes and wraps them in
// running nodes.
type buildRun struct {
	compiled []*compiledNode
	running  map[*compiledNode]*runningNode
}

func newBuildRun(compiled []*compiledNode) *buildRun {
	b := &buildRun{compiled: compiled,
		running: make(map[*compiledNode]*runningNode)}
	for _, cn := range compiled {
		rn := &runningNode{cn: cn}
		b.running[cn] = rn
	}
	return b
}

func (b *buildRun) buildPipeline(p *Pipeline, input *RunInput, env map[string]any) ([]*runningNode, error) {
	if len(p.roots) < 1 {
		return nil, fmt.Errorf("No roots")
	}
	for _, rn := range b.running {
		err := b.buildNode(rn, env)
		if err != nil {
			return nil, err
		}
		err = b.buildPins(rn)
		if err != nil {
			return nil, err
		}
	}
	// Queue up the roots
	running := make([]*runningNode, 0, len(p.roots))
	for _, cn := range p.roots {
		if rn, ok := b.running[cn]; ok {
			running = append(running, rn)
		} else {
			return nil, fmt.Errorf("missing running node for compiled node")
		}
	}
	if input != nil && len(input.Pins) > 0 {
		for _, n := range running {
			if n.input.Pins == nil {
				n.input.Pins = make([]Pin, 0, len(input.Pins))
			}
			n.input.Pins = append(n.input.Pins, input.Pins...)
		}
	}
	return running, nil
}

func (b *buildRun) buildNode(rn *runningNode, env map[string]any) error {
	rn.input.Pins = nil
	rn.inputCount = 0
	if rn.cn.starter != nil {
		startState := &State{}
		rn.cn.starter.Start(startState)
		rn.nodeData = startState.NodeData
	}
	// This is questionable -- the env vars will be applied to
	// the rn.nodeData. However, in the event there is no nodeData,
	// then I set the nodeData to the node itself, and they are applied
	// there. But that is not thread-safe... so, it's a trade off.
	// Slightly less boilerplate code if someone is writing a node
	// that truly doesn't need to be thread safe, but things will blow
	// up if the node is actually used concurrently.
	if rn.nodeData == nil {
		rn.nodeData = rn.cn.node
	}
	// Apply env vars
	envlen := len(rn.cn.envVars)
	if envlen > 0 {
		req := assign.ValuesRequest{}
		req.FieldNames = make([]string, envlen, envlen)
		req.NewValues = make([]any, envlen, envlen)
		i := -1
		for k, v := range rn.cn.envVars {
			i++
			req.FieldNames[i] = k
			if vv, ok := mapAt(v, env); ok {
				req.NewValues[i] = vv
			} else {
				return fmt.Errorf("Missing environment variable \"%v\"", v)
			}
		}
		return assign.Values(req, rn.nodeData)
	}
	return nil
}

func (b *buildRun) buildPins(rn *runningNode) error {
	rn.output = make([]*runningPin, 0, len(rn.cn.output))
	for _, cp := range rn.cn.output {
		if target, ok := b.running[cp.toNode]; ok {
			rp := &runningPin{cp: cp, toNode: target}
			rn.output = append(rn.output, rp)
		} else {
			return fmt.Errorf("no running pin for compiled")
		}
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
