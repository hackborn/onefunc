package geo

import (
	"math"

	"golang.org/x/exp/constraints"
)

type ProcessPtF64Func func(pt PointF64) PointF64

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func RadiansToDegrees(radians float64) float64 {
	return radians * (180.0 / math.Pi)
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
