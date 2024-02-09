package pipeline

type Pipeline struct {
	roots []*compiledNode
	nodes []*compiledNode
	env   map[string]any
}

// Env answers the contents of the env() term in the initial
// expression, if any. This has no functional impact, but serves
// as a form of documentation to let clients discover what env vars
// the pipeline supports.
func (p Pipeline) Env() map[string]any {
	// There is an experimental maps package that has a Clone,
	// but until that's in std do this.
	env := make(map[string]any)
	for k, v := range p.env {
		env[k] = v
	}
	return env
}

type compiledPin struct {
	inName string
	toNode *compiledNode
}

type compiledNode struct {
	node    Runner
	flusher Flusher
	starter Starter

	// nodeState will be either the result of StartNodeState() or the
	// node itself, but never nil.
	nodeState any
	// hasNodeData is true if nodeState came from StartNodeState()
	hasStartNodeState bool

	maxInputCount int
	output        []*compiledPin
	envVars       map[string]string
}

type runningPin struct {
	cp     *compiledPin
	toNode *runningNode
}

type runningNode struct {
	cn         *compiledNode
	inputCount int
	input      RunInput
	output     []*runningPin
}

func (n *runningNode) ready() bool {
	return n.inputCount >= n.cn.maxInputCount
}
