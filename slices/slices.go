package slices

// ArrayFrom takes a slice and a function and answers
// a new slice based on evaluating the function for
// each member of the incoming slice.
// Example:
//
//	return ArrayFrom(d.Fields, func(f structField) string {
//		return f.Field
//	})
func ArrayFrom[T any, U any](s []T, f func(T) U) []U {
	if len(s) < 1 {
		return []U{}
	}
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// SetLen returns a slice with the given length,
// reusing the supplied slice if it has enough capacity.
func SetLen[T any](s []T, area int) []T {
	if cap(s) < area {
		return make([]T, area, area)
	} else {
		return s[:area]
	}
}
