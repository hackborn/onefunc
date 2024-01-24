package pipeline

// Node is a processing node in the pipeline.
// Nodes can also implement io.Closer to handle
// any post-run behaviour.
type Node interface {
	// Run the supplied pins, producing output pins or an error.
	Run(*State, RunInput) (*RunOutput, error)
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
