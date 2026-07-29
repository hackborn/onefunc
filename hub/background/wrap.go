package background

// getTask will walk the wrapped task chain, answering
// the first task with a matching type.
func getTask[T any](src Task) (T, bool) {
	for src != nil {
		switch f := src.(type) {
		case T:
			return f, true
		}
		src = unwrap(src)
	}
	var t T
	return t, false
}

// unwrap answers the wrapped task, or nil.
func unwrap(t Task) Task {
	switch x := t.(type) {
	case interface{ Unwrap() Task }:
		return x.Unwrap()
	default:
		return nil
	}
}
