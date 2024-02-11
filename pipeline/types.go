package pipeline

// State provides current state data to a running node.
type State struct {
	// The data set in Starter.StartNode.
	NodeData any
}

func (s *State) SetNodeData(nd any) {
	s.NodeData = nd
}
