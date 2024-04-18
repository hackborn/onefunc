package geo

import (
	"math"
)

type Point[T Number] struct {
	X T
	Y T
}

func (p Point[T]) Area() T {
	return p.X * p.Y
}

func (a Point[T]) Dist(b Point[T]) float64 {
	// √((x2 - x1)² + (y2 - y1)²)
	x2 := math.Pow(float64(b.X)-float64(a.X), 2)
	y2 := math.Pow(float64(b.Y)-float64(a.Y), 2)
	return math.Sqrt(x2 + y2)
}

func (a Point[T]) Add(b Point[T]) Point[T] {
	return Point[T]{X: a.X + b.X, Y: a.Y + b.Y}
}

func (a Point[T]) Sub(b Point[T]) Point[T] {
	return Point[T]{X: a.X - b.X, Y: a.Y - b.Y}
}

func (a Point[T]) Normalize() Point[T] {
	if a.X == 0 && a.Y == 0 {
		return a
	}
	max := a.X
	if a.Y > max {
		max = a.Y
	}
	return Point[T]{X: a.X / max, Y: a.Y / max}
}

func (a Point[T]) Magnitude() float64 {
	x, y := float64(a.X), float64(a.Y)
	return math.Sqrt(x*x + y*y)
}

func (a Point[T]) Radians() float64 {
	return math.Atan2(float64(a.Y), float64(a.X))
}

func (a Point[T]) Inside(r Rect[T]) bool {
	return a.X >= r.LT.X && a.Y >= r.LT.Y && a.X < r.RB.X && a.Y < r.RB.Y
}

// Given slope m and distance, project the positive and negative points on the line.
// https://stackoverflow.com/questions/1250419/finding-points-on-a-line-with-a-given-distance
func (a Point[T]) Project(m, dist float64) (Point[T], Point[T]) {
	af := PointF64{X: float64(a.X), Y: float64(a.Y)}
	n := PointF64{}
	if m >= math.MaxFloat64 {
		n = PointF64{X: 0.0, Y: 1.0}
	} else {
		magnitude := math.Pow(1*1+m*m, 1.0/2.0)
		n = PointF64{X: 1.0 / magnitude, Y: m / magnitude}
	}
	n.X *= dist
	n.Y *= dist
	pp := af.Add(n)
	pn := af.Sub(n)
	return Point[T]{X: T(pp.X), Y: T(pp.Y)}, Point[T]{X: T(pn.X), Y: T(pn.Y)}
}

// ToIndex converts the xy point into an index into a flat
// array as represented by this point.
func (p Point[T]) ToIndex(xy Point[T]) T {
	return (xy.Y * p.X) + xy.X
}

func ConvertPoint[A Number, B Number](a Point[A]) Point[B] {
	return Point[B]{X: B(a.X), Y: B(a.Y)}
}

type PointF32 = Point[float32]
type PointF64 = Point[float64]
type PointI = Point[int]
type PointI64 = Point[int64]
type PointUI64 = Point[uint64]
