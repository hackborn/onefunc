package geo

import (
	"math"

	"golang.org/x/exp/constraints"
)

type ProcessPtF64Func func(pt PtF) PtF

// HitTest answers true if the point is a hit.
type HitTest func(PtF) bool

func NullHitTest(PtF) bool {
	return false
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func RadiansToDegrees(radians float64) float64 {
	return radians * (180.0 / math.Pi)
}

func Sqr[T Number](x T) T {
	return x * x
}

func Centroid[T Number](pts []Point[T]) Point[T] {
	var x float64 = 0.0
	var y float64 = 0.0
	for _, pt := range pts {
		x += float64(pt.X)
		y += float64(pt.Y)
	}
	x /= float64(len(pts))
	y /= float64(len(pts))
	return Point[T]{X: T(x), Y: T(y)}
}

// IndexToXY converts a flat index into an array into
// an XY position.
func IndexToXY[T constraints.Integer](width, height, index T) (Point[T], error) {
	if index < 0 || index >= width*height {
		return Point[T]{}, outOfBoundsErr
	}
	return Point[T]{X: index % width, Y: index / width}, nil
}

// XYToIndex converts an X, Y to a flat index.
func XYToIndex[T constraints.Integer](width, height T, pt Point[T]) (T, error) {
	// TODO: But this isn't really right, we need to check x and y bounds
	// or else we might get wraparound, right?
	idx := (pt.Y * width) + pt.X
	if idx < 0 || idx >= (width*height) {
		return 0, outOfBoundsErr
	}
	return idx, nil
}

func FloatsEqual(a, b float64) bool {
	const eps = 0.000000000000001
	diff := math.Abs(a - b)
	return diff < eps
}

func FloatsEqualTol(a, b, tolerance float64) bool {
	diff := math.Abs(a - b)
	return diff < tolerance
}

func cross[T Number](a, b Point[T]) T {
	return a.X*b.Y - a.Y*b.X
}

func orient[T Number](a, b, c Point[T]) T {
	return cross(b.Sub(a), c.Sub(a))
}

// Ratio answers a 0-1 value based on a's
// contributing factor to the final value.
func Ratio(a, b float64) float64 {
	if a <= 0 {
		return 0
	}
	sum := a + b
	return a / sum
}

// Surely these are SOMEWHERE?

func Min[T Number](a, b T) T {
	if a <= b {
		return a
	}
	return b
}

func Max[T Number](a, b T) T {
	if a >= b {
		return a
	}
	return b
}
