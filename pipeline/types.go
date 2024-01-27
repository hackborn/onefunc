package pipeline

// State is a placeholder for any global data passed through
// the pipeline.
type State struct {
	// Flush is true during the flush stage of pipeline Run.
	// Most nodes can ignore this, but source nodes should not
	// generate data.
	Flush bool
}
