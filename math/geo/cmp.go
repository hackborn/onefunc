package geo

import (
	"math"
)

func FloatsEqual(a, b float64) bool {
	const eps = 0.000000000000001
	diff := math.Abs(a - b)
	return diff < eps
}

func FloatsEqualTol(a, b, tolerance float64) bool {
	diff := math.Abs(a - b)
	return diff < tolerance
}

func RectsEqual[T Number](a, b Rectangle[T]) bool {
	return a.L == b.L && a.T == b.T && a.R == b.R && a.B == b.B
}
