package geo

import (
	"math"
)

// Pt is shorthand for creating a point from two values.
func Pt[T Number](x, y T) Point[T] {
	return Point[T]{X: x, Y: y}
}

// Pts creates a slice of Points from x y values.
func Pts[T Number](xys ...T) []Point[T] {
	pts := make([]Point[T], 0, len(xys)/2)
	for i := 1; i < len(xys); i += 2 {
		pts = append(pts, Point[T]{X: xys[i-1], Y: xys[i]})
	}
	return pts
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

func (a Point[T]) Cross(b Point[T]) T {
	return a.X*b.Y - a.Y*b.X
}

func (a Point[T]) AddValue(v T) Point[T] {
	return Point[T]{X: a.X + v, Y: a.Y + v}
}

func (a Point[T]) SubValue(v T) Point[T] {
	return Point[T]{X: a.X - v, Y: a.Y - v}
}

// Normalize normalizs the point to a unit.
// Can be negative if they are generated from a negative point,
// i.e. a negative direction vector.
func (a Point[T]) Normalize() Point[T] {
	mag := a.Magnitude()
	if mag <= 0. {
		return Point[T]{X: 0, Y: 0}
	}
	return Point[T]{X: T(float64(a.X) / mag),
		Y: T(float64(a.Y) / mag)}
}

func (a Point[T]) Dot(b Point[T]) float64 {
	return float64(a.X*b.X + a.Y*b.Y)
}

// Length is the magnitude of the point.
// I *think* this is just more standard
// nomenclature than Magnitude.
func (a Point[T]) Length() float64 {
	x, y := float64(a.X), float64(a.Y)
	return math.Sqrt(x*x + y*y)
}

// Magnitude is idental to Length, but I think
// that's a more common term.
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

// Rotate the point. Rotation will be:
// In Quandrant 4 this rotates CW for positive angles and CCW for negative.
// TODO: API might be confusing, after not using it I've started thinking
// of the a point as the center and the point being passed in as the point to
// rotate. so might switch those.
func (a Point[T]) Rotate(origin Point[T], angleInRads float64) Point[T] {
	// from https://stackoverflow.com/questions/2259476/rotating-a-point-about-another-point-2d

	// translate to origin
	a.X -= origin.X
	a.Y -= origin.Y

	af := ConvertPoint[T, float64](a)

	// rotate
	s := math.Sin(angleInRads)
	c := math.Cos(angleInRads)
	xnew := af.X*c - af.Y*s
	ynew := af.X*s + af.Y*c

	// translate back
	return Point[T]{X: origin.X + T(xnew),
		Y: origin.Y + T(ynew)}
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

func PtBounds[T Number](pts []Point[T]) Rectangle[T] {
	r := Rectangle[T]{L: 0, T: 0, R: 0, B: 0}
	for i, pt := range pts {
		if i == 0 {
			r.L, r.T, r.R, r.B = pt.X, pt.Y, pt.X, pt.Y
		} else {
			r.L, r.T = min(r.L, pt.X), min(r.T, pt.Y)
			r.R, r.B = max(r.R, pt.X), max(r.B, pt.Y)
		}
	}
	return r
}
