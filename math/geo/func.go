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

// RadialDist answers the distance between two unit values,
// accounting for wrapping. Andswer 0 - 1, with 1 being
// the two values are closest.
// Examples:
// a=.1, b=.9, distance is .2, i.e. it wraps around 1.
func RadialDist(a, b float64) (float64, Orientation) {
	if a == b {
		return 1., Collinear
	}
	a1 := a - b
	a2 := b - a
	if a1 < 0. {
		a1 = 1. + a1
	}
	if a2 < 0. {
		a2 = 1. + a2
	}
	// This results in a range of 0. - .5, so scale up
	if a1 < a2 {
		return 1. - (a1 * 2.), Clockwise
	} else {
		return 1. - (a2 * 2.), CounterClockwise
	}
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
	if pt.X < 0 || pt.X >= width || pt.Y < 0 || pt.Y >= height {
		return 0, outOfBoundsErr
	}
	return XYToIndexFast(pt.X, pt.Y, width), nil
	/*
			idx := (pt.Y * width) + pt.X
			if idx < 0 || idx >= (width*height) {
				return 0, outOfBoundsErr
			}
		return idx, nil
	*/
}

// XYToIndexFast converts an X, Y to a flat index without any bounds checking.
func XYToIndexFast[T constraints.Integer](x, y, width T) T {
	return (y * width) + x
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

// Orient answers whether C is left / clockwise to direct
// segment AB. Response:
// 0: C is collinear
// < 0: C is left / clockwise
// > 0: C is right / counterclockwise
func Orient[T Number](a, b, c Point[T]) T {
	return b.Sub(a).Cross(c.Sub(a))
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
// TODO: Yes, builtin min() and max(), and they are highly
// optimized, perform the same as this. So remove these.
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
