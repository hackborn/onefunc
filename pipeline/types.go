package pipeline

// State provides current state data to a running node.
type State struct {
	// The data set in Starter.Start.
	NodeData any
}

type _startInput struct {
	nodeData any
}

func (s *_startInput) SetNodeData(nd any) {
	s.nodeData = nd
}
