package geo

import (
	"math"
)

// Pt3d is shorthand for creating a point3d from three values.
func Pt3d[T Number](x, y, z T) Point3d[T] {
	return Point3d[T]{X: x, Y: y, Z: z}
}

type Point3d[T Number] struct {
	X T
	Y T
	Z T
}

func (a Point3d[T]) XY() Point[T] {
	return Point[T]{X: a.X, Y: a.Y}
}

func (a Point3d[T]) Scale(t T) Point3d[T] {
	return Point3d[T]{X: a.X * t, Y: a.Y * t, Z: a.Z * t}
}

func (a Point3d[T]) Add(b Point3d[T]) Point3d[T] {
	return Point3d[T]{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

func (a Point3d[T]) Sub(b Point3d[T]) Point3d[T] {
	return Point3d[T]{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func (a Point3d[T]) Mult(b Point3d[T]) Point3d[T] {
	return Point3d[T]{X: a.X * b.X, Y: a.Y * b.Y, Z: a.Z * b.Z}
}

func (a Point3d[T]) Div(b Point3d[T]) Point3d[T] {
	return Point3d[T]{X: a.X / b.X, Y: a.Y / b.Y, Z: a.Z / b.Z}
}

func (a Point3d[T]) SubT(t T) Point3d[T] {
	return Point3d[T]{X: a.X - t, Y: a.Y - t, Z: a.Z - t}
}

func (a Point3d[T]) MultT(t T) Point3d[T] {
	return Point3d[T]{X: a.X * t, Y: a.Y * t, Z: a.Z * t}
}

func (a Point3d[T]) DivT(t T) Point3d[T] {
	return Point3d[T]{X: a.X / t, Y: a.Y / t, Z: a.Z / t}
}

func (a Point3d[T]) Magnitude() float64 {
	// √(a2 + b2 + c2)
	x, y, z := float64(a.X), float64(a.Y), float64(a.Z)
	return math.Sqrt(x*x + y*y + z*z)
}

// Mag returns the magnitude (length) of the vector
//func (a Point3d[T]) Mag() float64 {
//	return math.Sqrt(math.Pow(float64(a.X), 2) + math.Pow(float64(a.Y), 2) + math.Pow(float64(a.Z), 2))
//}

func (a Point3d[T]) Normalize() Point3d[T] {
	mag := a.Magnitude()
	if mag > 1e-6 {
		a.X /= T(mag)
		a.Y /= T(mag)
		a.Z /= T(mag)
	}
	return a
}

//func (a Point3d[T]) Normalize() Point3d[T] {
//	mag := a.Magnitude()
//	x, y, z := float64(a.X)/mag, float64(a.Y)/mag, float64(a.Z)/mag
//	return Point3d[T]{X: T(x), Y: T(y), Z: T(z)}
//}

func (a Point3d[T]) Dot(b Point3d[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Point3d[T]) DotProduct(b Point3d[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross calculates the cross product of two vectors and returns a new vector
func (a Point3d[T]) Cross(b Point3d[T]) Point3d[T] {
	return Point3d[T]{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (a Point3d[T]) CLampFast(rng Range[T]) Point3d[T] {
	return Point3d[T]{X: rng.ClampFast(a.X),
		Y: rng.ClampFast(a.Y),
		Z: rng.ClampFast(a.Z),
	}
}

func Pt3dNearZero(v Pt3dF, tolerance float64) bool {
	return math.Abs(v.X) <= tolerance && math.Abs(v.Y) <= tolerance && math.Abs(v.Z) <= tolerance
}

func Float64NearZero(v float64, tolerance float64) bool {
	return math.Abs(v) <= tolerance
}
