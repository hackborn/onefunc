package geo

import (
	"math"
)

// Pt is shorthand for creating a point from two values.
func Pt[T Number](x, y T) Point[T] {
	return Point[T]{X: x, Y: y}
}

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

// Dist2 is the distance without the square root.
func (a Point[T]) DistSquared(b Point[T]) T {
	return Sqr(a.X-b.X) + Sqr(a.Y-b.Y)
}

func (a Point[T]) Add(b Point[T]) Point[T] {
	return Point[T]{X: a.X + b.X, Y: a.Y + b.Y}
}

func (a Point[T]) Sub(b Point[T]) Point[T] {
	return Point[T]{X: a.X - b.X, Y: a.Y - b.Y}
}

func (a Point[T]) Mult(b Point[T]) Point[T] {
	return Point[T]{X: a.X * b.X, Y: a.Y * b.Y}
}

// ??? What? This doesn't look like a normalize, it
// needs to divide by the magnitude. Is anyone using this?
// Well hold on maybe this is valid (with tweaks), from the C#
// source:
// https://www.dotnetframework.org/default.aspx/Net/Net/3@5@50727@3053/DEVDIV/depot/DevDiv/releases/Orcas/SP/wpf/src/Core/CSharp/System/Windows/Media3D/Vector3D@cs/1/Vector3D@cs
/*
public void Normalize()
       {
           // Computation of length can overflow easily because it
           // first computes squared length, so we first divide by
           // the largest coefficient.
           double m = Math.Abs(_x);
           double absy = Math.Abs(_y);
           double absz = Math.Abs(_z);
           if (absy > m)
           {
               m = absy;
           }
           if (absz > m)
           {
               m = absz;
           }

           _x /= m;
           _y /= m;
           _z /= m;

           double length = Math.Sqrt(_x * _x + _y * _y + _z * _z);
           this /= length;
       }
*/
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

// Degrees finds the angle of the segment with this point as the origin.
// Degrees will be 0-360, with 0/360 on the right, proceeding clockwise.
func (a Point[T]) Degrees(b Point[T]) float64 {
	angle := math.Atan2(float64(b.Y-a.Y), float64(b.X-a.X))
	v := angle*180/math.Pi + 180
	// Rotate to put 0 at hard right
	v += 180
	if v >= 360 {
		v -= 360
	}
	return v
}

func (a Point[T]) Inside(r Rectangle[T]) bool {
	return a.X >= r.L && a.Y >= r.T && a.X < r.R && a.Y < r.B
}

// Given slope m and distance, project the positive and negative points on the line.
// Uses upper-left coordinate system.
func (a Point[T]) Project(m, dist float64) (Point[T], Point[T]) {
	af := PtF{X: float64(a.X), Y: float64(a.Y)}

	// Special case -- vertical lines as reported by Slope()
	if m == math.MaxFloat64 {
		return Point[T]{X: a.X, Y: T(af.Y + dist)},
			Point[T]{X: a.X, Y: T(af.Y - dist)}
	}

	denom := math.Sqrt(1. + m*m)
	x1 := dist * (1. / denom)
	y1 := dist * (m / denom)
	if m < 0. {
		x1 = -x1
	} else if m > 0. {
		y1 = -y1
	}

	return Point[T]{X: T(af.X + x1), Y: T(af.Y + y1)},
		Point[T]{X: T(af.X - x1), Y: T(af.Y - y1)}
}

// Given slope m and distance, project the positive and negative points on the line.
// https://stackoverflow.com/questions/1250419/finding-points-on-a-line-with-a-given-distance
// A bit confused because I wasn't paying attention to the coordinate
// system, but I was using this to check against my real algo for awhile.
func (a Point[T]) projectAlgo2(m, dist float64) (Point[T], Point[T]) {
	n := PtF{}
	if m >= math.MaxFloat64 {
		n = PtF{X: 0., Y: 1.}
	} else if FloatsEqualTol(m, 0., 0.00001) {
		n = PtF{X: 1., Y: 0.}
	} else {
		magnitude := math.Pow(1*1+m*m, 1./2.)
		// n = PtF{X: m / magnitude, Y: 1. / magnitude}
		n = PtF{X: 1. / magnitude, Y: m / magnitude}
	}
	n.X *= dist
	n.Y *= dist
	af := PtF{X: float64(a.X), Y: float64(a.Y)}
	pp := af.Add(n)
	pn := af.Sub(n)
	return Point[T]{X: T(pp.X), Y: T(pp.Y)}, Point[T]{X: T(pn.X), Y: T(pn.Y)}
}

// projectAlgo3 is a different implementation of Project found via gemini.
// I assume the Sqrt makes it a little slower, but it also has no branching
// so who can say. In comparisons they both work the same.
// Update, doesn't work right, possibly because of not paying attention to
// the coord system.
/*
func (a Point[T]) projectAlgo3(m, dist float64) (Point[T], Point[T]) {
	dyP := dist / math.Sqrt(1+m*m) // Δy based on distance and slope
	dyN := -dyP
	dxP := math.Sqrt(dist*dist - dyP*dyP)
	dxN := -dxP
	af := PtF{X: float64(a.X), Y: float64(a.Y)}
	return Point[T]{X: T(af.X + dxP), Y: T(af.Y + dyP)}, Point[T]{X: T(af.X + dxN), Y: T(af.Y + dyN)}
}
*/

// ProjectDegree takes a degree and distance and projects a new point.
func (a Point[T]) ProjectDegree(deg, dist float64) Point[T] {
	radians := DegreesToRadians(deg)
	return Point[T]{X: a.X + T(dist*math.Cos(radians)),
		Y: a.Y + T(dist*math.Sin(radians))}
}

// ToIndex converts the xy point into an index into a flat
// array as represented by this point.
func (p Point[T]) ToIndex(xy Point[T]) T {
	return (xy.Y * p.X) + xy.X
}

func ConvertPoint[A Number, B Number](a Point[A]) Point[B] {
	return Point[B]{X: B(a.X), Y: B(a.Y)}
}

type PtF = Point[float64]
type PtI = Point[int]

type PtF32 = Point[float32]
type PtF64 = Point[float64]
type PtI64 = Point[int64]
