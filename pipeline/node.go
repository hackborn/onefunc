package pipeline

// ---------------------------------------------------------
// NODE

// Node is a single function in the pipeline graph. It runs an
// operation on some input, optionally providing some output.
//
// Node implementations should be thread-safe. The framework uses
// a NodeData pattern to accomplish this in a simple way: Nodes have
// a parallel NodeData struct that stores their data. A node anonymously
// includes its NodeData, and also sets it to State.NodeData in the
// Start() func. THe data is then passed back in State.Node data
// during the Run() and Flush().
//
// Optional interfaces:
// * Starter
// * Flusher
type Node interface {
	Runner
}

// ---------------------------------------------------------
// RUNNER

// Runner processes input.
type Runner interface {
	// Run the supplied pins, converting the input to output.
	// There's no reason to reallocate or replace the RunOutput.Pins,
	// just append to what's there.
	Run(*State, RunInput, *RunOutput) error
}

// NewRunInput answers a new RunInput on the given pins.
func NewRunInput(pins ...Pin) RunInput {
	return RunInput{Pins: pins}
}

// RunInput transports pin data into a Run function.
type RunInput struct {
	Pins []Pin
}

// NewRunOutput answers a new RunOutput on the given pins.
func (r RunInput) NewRunOutput(pins []Pin) *RunOutput {
	return &RunOutput{Pins: pins}
}

// RunOutput transports pin data out of a Run function.
type RunOutput struct {
	Pins []Pin
}

// ---------------------------------------------------------
// STARTER

// Starter is called at the start of a pipeline run. Implementing
// starter, and placing all run state in a node data object,
// will make your node thread-safe.
type Starter interface {
	// One-time notification that the node is starting a run.
	// Clients that want to use the NodeData pattern can assign
	// the NodeData field and it will be available during the Run.
	Start(StartInput) error
}

type StartInput interface {
	SetNodeData(any)
}

// ---------------------------------------------------------
// FLUSHER

// Flusher implements a flush operation.
type Flusher interface {
	// Flush any data in the node.
	// There's no reason to reallocate or replace the RunOutput.Pins,
	// just append to what's there.
	Flush(*State, *RunOutput) error
}
