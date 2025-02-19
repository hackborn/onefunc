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

// Pop answers the last element of T, removing from the slice.
func Pop[T any](s []T) ([]T, T) {
	if len(s) < 1 {
		var t T
		return s, t
	}
	t := s[len(s)-1]
	s = s[0 : len(s)-1]
	return s, t
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

// SetLenInit returns a slice with the given length,
// reusing the supplied slice if it has enough capacity.
// If the slice is reused then all values are initialized.
func SetLenInit[T any](s []T, area int) []T {
	if cap(s) < area {
		return make([]T, area, area)
	} else {
		newS := s[:area]
		for i := range area {
			var t T
			newS[i] = t
		}
		return newS
	}
}

// SetLenCopy returns a slice with the given length,
// reusing the supplied slice if it has enough capacity.
// If the slice is created then the original values are copied.
func SetLenCopy[T any](s []T, area int) []T {
	if cap(s) < area {
		newS := make([]T, area, area)
		for i := 0; i < min(len(s), area); i++ {
			newS[i] = s[i]
		}
		return newS
	} else {
		return s[:area]
	}
}
