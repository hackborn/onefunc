package geo

import (
	"math"
)

// Pt3 is shorthand for creating a point3d from three values.
func Pt3[T Number](x, y, z T) Point3d[T] {
	return Point3d[T]{X: x, Y: y, Z: z}
}

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

func (a Point3d[T]) MultScalar(t T) Point3d[T] {
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

func (a Point3d[T]) LengthSq() float64 {
	return float64(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
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

func (a Point3d[T]) DotProduct(b Point3d[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// TODO: Replace with DotProduct
func (a Point3d[T]) Dot(b Point3d[T]) T {
	return a.DotProduct(b)
}

// Cross calculates the cross product of two vectors and returns a new vector
func (a Point3d[T]) Cross(b Point3d[T]) Point3d[T] {
	return Point3d[T]{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (a Point3d[T]) ClampFast(rng Range[T]) Point3d[T] {
	return Point3d[T]{X: rng.ClampFast(a.X),
		Y: rng.ClampFast(a.Y),
		Z: rng.ClampFast(a.Z),
	}
}

func (a Point3d[T]) Distance(b Point3d[T]) float64 {
	return math.Sqrt(float64(a.DistanceSquared(b)))
}

func (a Point3d[T]) DistanceSquared(b Point3d[T]) T {
	dx := b.X - a.X
	dy := b.Y - a.Y
	dz := b.Z - a.Z
	return dx*dx + dy*dy + dz*dz
}

// Distance2 is the Distance function but only using 2 components.
func (a Point3d[T]) Distance2(b Point3d[T]) float64 {
	return math.Sqrt(float64(a.DistanceSquared2(b)))
}

// DistanceSquared2 is the DistanceSquared2 function but only using 2 components.
func (a Point3d[T]) DistanceSquared2(b Point3d[T]) T {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return dx*dx + dy*dy
}

// RotateX rotates a Vector3 around the X-axis by a given angle (in radians).
func (v Point3d[T]) RotateX(angle float64) Point3d[T] {
	sinTheta := T(math.Sin(angle))
	cosTheta := T(math.Cos(angle))
	newY := v.Y*cosTheta - v.Z*sinTheta
	newZ := v.Y*sinTheta + v.Z*cosTheta
	return Point3d[T]{X: v.X, Y: newY, Z: newZ}
}

// RotateY rotates a Vector3 around the Y-axis by a given angle (in radians).
func (v Point3d[T]) RotateY(angle float64) Point3d[T] {
	sinTheta := T(math.Sin(angle))
	cosTheta := T(math.Cos(angle))
	newX := v.X*cosTheta + v.Z*sinTheta
	newZ := -v.X*sinTheta + v.Z*cosTheta
	return Point3d[T]{X: newX, Y: v.Y, Z: newZ}
}

// RotateZ rotates a Vector3 around the Z-axis by a given angle (in radians).
func (v Point3d[T]) RotateZ(angle float64) Point3d[T] {
	sinTheta := T(math.Sin(angle))
	cosTheta := T(math.Cos(angle))
	newX := v.X*cosTheta - v.Y*sinTheta
	newY := v.X*sinTheta + v.Y*cosTheta
	return Point3d[T]{X: newX, Y: newY, Z: v.Z}
}

func Pt3dNearZero(v Pt3dF, tolerance float64) bool {
	return math.Abs(v.X) <= tolerance && math.Abs(v.Y) <= tolerance && math.Abs(v.Z) <= tolerance
}

func Float64NearZero(v float64, tolerance float64) bool {
	return math.Abs(v) <= tolerance
}

func ConvertPoint3[A Number, B Number](a Point3d[A]) Point3d[B] {
	return Point3d[B]{X: B(a.X), Y: B(a.Y), Z: B(a.Z)}
}
