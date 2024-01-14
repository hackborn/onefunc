package pipeline

// Node is a processing node in the pipeline.
// Nodes can also implement io.Closer to handle
// any post-run behaviour.
type Node interface {
	// Run the supplied pins, producing output pins or an error.
	Run(*State, RunInput) (*RunOutput, error)
}

func NewInput(pins ...PinData) RunInput {
	return RunInput{Pins: pins}
}

// RunInput trapnsports pin data into a Run function.
type RunInput struct {
	Pins []PinData
}

func (r RunInput) NewOutput(pins []PinData) *RunOutput {
	return &RunOutput{Pins: pins}
}

// RunOutput transports pin data out of a Run function.
type RunOutput struct {
	Pins []PinData
}

// NewInput creates a new input object on all the output pin data.
func (r RunOutput) newInput() RunInput {
	return RunInput{Pins: r.Pins}
}
