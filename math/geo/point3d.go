package geo

import "math"

// Pt3d is shorthand for creating a point3d from three values.
func Pt3d[T Number](x, y, z T) Point3d[T] {
	return Point3d[T]{X: x, Y: y, Z: z}
}

type Point3d[T Number] struct {
	X T
	Y T
	Z T
}

func (a Point3d[T]) Add(b Point3d[T]) Point3d[T] {
	return Point3d[T]{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

func (a Point3d[T]) Sub(b Point3d[T]) Point3d[T] {
	return Point3d[T]{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func (a Point3d[T]) Magnitude() float64 {
	// √(a2 + b2 + c2)
	x, y, z := float64(a.X), float64(a.Y), float64(a.Z)
	return math.Sqrt(x*x + y*y + z*z)
}

func (a Point3d[T]) Normalize() Point3d[T] {
	mag := a.Magnitude()
	x, y, z := float64(a.X)/mag, float64(a.Y)/mag, float64(a.Z)/mag
	return Point3d[T]{X: T(x), Y: T(y), Z: T(z)}
}

func (a Point3d[T]) DotProduct(b Point3d[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

type Pt3dF = Point3d[float64]
type Pt3dI = Point3d[int]

type Pt3dF64 = Point3d[float64]
type Pt3dI64 = Point3d[int64]
type Pt3dUI64 = Point3d[uint64]
