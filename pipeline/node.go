package pipeline

// Node is a single processor in the pipeline graph. It's basic
// responsibility is to run an operation on some input, but it has
// additional optional behaviour.
//
// Note that node implementations should be immutable and thread-safe.
// If you have working state to persist beyond a single run, use
// GetNodeState().
//
// Optional interfaces:
// * Starter
// * Flusher
type Node interface {
	Runner
}

// Runner processes input.
type Runner interface {
	// Run the supplied pins, producing output pins or an error.
	Run(*State, RunInput) (*RunOutput, error)
}

// Starter is called at the start of a pipeline run. Implementing
// starter, and placing all run state in a node data object,
// will make your node thread-safe.
type Starter interface {
	// One-time notification that the node is starting a run.
	// Clients that want to use the NodeData pattern can assign
	// the NodeData field and it will be available during the Run.
	StartNode(*State)
}

// Flusher implements a flush operation.
type Flusher interface {
	// Flush any data in the node.
	Flush(*State) (*RunOutput, error)
}

func NewInput(pins ...Pin) RunInput {
	return RunInput{Pins: pins}
}

// RunInput trapnsports pin data into a Run function.
type RunInput struct {
	Pins []Pin
}

func (r RunInput) NewOutput(pins []Pin) *RunOutput {
	return &RunOutput{Pins: pins}
}

// RunOutput transports pin data out of a Run function.
type RunOutput struct {
	Pins []Pin
}
