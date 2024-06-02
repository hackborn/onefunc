package pipeline

import (
	"fmt"

	"github.com/hackborn/onefunc/values"
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
	if err != nil {
		return nil, err
	}

	//	fmt.Println("pipeline running, roots", len(p.roots), "active", len(active))
	state := &State{}
	runOutput := &RunOutput{}
	outputPins := make([]Pin, 0, 8)
	finalOutput := RunOutput{}
	for len(running) > 0 {
		for _, rn := range running {
			state.NodeData = rn.nodeData
			runOutput.Pins = outputPins[:0]
			err := rn.cn.node.Run(state, rn.input, runOutput)
			if err != nil {
				return nil, fmt.Errorf("Pipeline: %T run err: %w", rn.cn.node, err)
			}
			// This node is done processing, flush it.
			err = flush(state, rn.cn.flusher, runOutput)
			if err != nil {
				return nil, fmt.Errorf("Pipeline: %T flush err: %w", rn.cn.node, err)
			}

			if len(rn.output) > 0 {
				fanout := runFanOut{}
				for _, topin := range rn.output {
					topin.toNode.inputCount++
					fanout.on(topin, runOutput)
					if topin.toNode.ready() {
						topin.toNode.isReady = true
					}
				}
			} else if len(runOutput.Pins) > 0 {
				finalOutput.Pins = append(finalOutput.Pins, runOutput.Pins...)
			}
		}
		running = running[:0]
		for _, nextrn := range build.running {
			if nextrn.isReady {
				nextrn.isReady = false
				running = append(running, nextrn)
			}
		}
	}
	return &finalOutput, nil
}

func flush(state *State, flusher Flusher, output *RunOutput) error {
	if flusher == nil {
		return nil
	}
	return flusher.Flush(state, output)
}

// runFanOut is responsible for adding runOutput to the input of
// destination nodes, cloning according to policy.
type runFanOut struct {
}

func (f *runFanOut) on(topin *runningPin, runOutput *RunOutput) {
	if len(runOutput.Pins) < 1 {
		return
	}
	tonode := topin.toNode
	first := -1
	for _, pin := range runOutput.Pins {
		first++
		if pin.Payload == nil {
			tonode.input.Pins = append(tonode.input.Pins, pin)
			continue
		}
		newpin := pin
		switch pin.Policy {
		case AlwaysClone:
			newpin.Payload = pin.Payload.Clone()
		case NeverClone:
		default:
			if first > 0 {
				newpin.Payload = pin.Payload.Clone()
			}
		}
		tonode.input.Pins = append(tonode.input.Pins, newpin)
	}
}

// initFanOut is responsible for adding runoutput to the input of
// destination nodes, cloning according to policy.
type initFanOut struct {
}

func (f *initFanOut) on(index int, src []Pin, dst []Pin) []Pin {
	for _, srcpin := range src {
		if srcpin.Payload == nil {
			dst = append(dst, srcpin)
			continue
		}
		dstpin := srcpin
		switch srcpin.Policy {
		case AlwaysClone:
			dstpin.Payload = srcpin.Payload.Clone()
		case NeverClone:
		default:
			if index > 0 {
				dstpin.Payload = srcpin.Payload.Clone()
			}
		}
		dst = append(dst, dstpin)
	}
	return dst
}

// buildRun takes the compiled nodes and wraps them in
// running nodes.
type buildRun struct {
	compiled []*compiledNode
	running  map[*compiledNode]*runningNode
}

func newBuildRun(compiled []*compiledNode) *buildRun {
	running := make(map[*compiledNode]*runningNode, len(compiled))
	b := &buildRun{compiled: compiled,
		running: running}
	for _, cn := range compiled {
		rn := &runningNode{cn: cn}
		b.running[cn] = rn
	}
	return b
}

func (b *buildRun) buildPipeline(p *Pipeline, input *RunInput, env map[string]any) ([]*runningNode, error) {
	if len(p.roots) < 1 {
		return nil, fmt.Errorf("Pipeline: No roots")
	}
	for _, rn := range b.running {
		err := b.buildNode(rn, env)
		if err != nil {
			return nil, fmt.Errorf("Pipeline: build node err %w", err)
		}
		err = b.buildPins(rn)
		if err != nil {
			return nil, fmt.Errorf("Pipeline: build pins err %w", err)
		}
	}
	// Queue up the roots
	running := make([]*runningNode, 0, len(p.roots))
	for _, cn := range p.roots {
		if rn, ok := b.running[cn]; ok {
			running = append(running, rn)
		} else {
			return nil, fmt.Errorf("Pipeline: Missing running node for compiled node")
		}
	}
	if input != nil && len(input.Pins) > 0 {
		fanout := initFanOut{}
		for i, n := range running {
			if n.input.Pins == nil {
				n.input.Pins = make([]Pin, 0, len(input.Pins))
			}
			//			n.input.Pins = append(n.input.Pins, input.Pins...)
			n.input.Pins = fanout.on(i, input.Pins, n.input.Pins)
		}
	}
	return running, nil
}

func (b *buildRun) buildNode(rn *runningNode, env map[string]any) error {
	rn.input.Pins = nil
	rn.inputCount = 0
	if rn.cn.starter != nil {
		si := &_startInput{}
		err := rn.cn.starter.Start(si)
		if err != nil {
			return err
		}
		rn.nodeData = si.nodeData
	}
	// This is questionable -- the env vars will be applied to
	// the rn.nodeData. However, in the event there is no nodeData,
	// then I set the nodeData to the node itself, and they are applied
	// there. But that is not thread-safe... so, it's a trade off.
	// Slightly less boilerplate code if someone is writing a node
	// that truly doesn't need to be thread safe, but things will blow
	// up if the node is actually used concurrently.
	// Probably... probably should not do this and you just don't get
	// env vars if you don't supply node data. Kind of annoying, but safe.
	if rn.nodeData == nil {
		rn.nodeData = rn.cn.node
	}
	// Apply env vars
	envlen := len(rn.cn.envVars)
	if envlen > 0 {
		req := values.SetRequest{}
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
		return values.Set(req, rn.nodeData)
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

// RunNode is a convenience to run a single node.
func RunNode(node Node, input RunInput, output *RunOutput) error {
	state := &State{}
	if starter, ok := node.(Starter); ok {
		si := _startInput{}
		if err := starter.Start(&si); err != nil {
			return err
		}
		state.NodeData = si.nodeData
	}
	return node.Run(state, input, output)
}
