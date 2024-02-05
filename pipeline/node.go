package pipeline

// Node is a single processor in the pipeline graph.
// Optional interfaces:
// * Flusher
type Node interface {
	Runner
}

// Runner processes input.
type Runner interface {
	// Run the supplied pins, producing output pins or an error.
	Run(*State, RunInput) (*RunOutput, error)
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
