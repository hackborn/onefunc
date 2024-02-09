package pipeline

// State is a placeholder for any global data passed through
// the pipeline.
type State struct {
	// nodeState stores any per-run working state nodes want
	// to persist between runs.
	nodeState map[Node]any
}

// GetNodeState answers or creates a struct for the given
// node key. This should be used by any nodes that want data
// to persist beyond a single call to run, so they can be
// thread-safe.
func GetNodeState[T any](key Node, s *State) *T {
	if s.nodeState != nil {
		if ns := s.nodeState[key]; ns != nil {
			if cast, ok := ns.(*T); ok {
				return cast
			}
		}
	}
	value := new(T)
	setNodeState(key, value, s)
	return value
}

func setNodeState(key Node, value any, s *State) {
	if value == nil {
		if s.nodeState != nil {
			delete(s.nodeState, key)
		}
	} else {
		if s.nodeState == nil {
			s.nodeState = make(map[Node]any)
		}
		s.nodeState[key] = value
	}
}
