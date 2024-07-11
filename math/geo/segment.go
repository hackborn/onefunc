package geo

import (
	"math"
)

// Seg is shorthand for creating a segment from two points.
func Seg[T Number](ax, ay, bx, by T) Segment[T] {
	return Segment[T]{A: Pt(ax, ay), B: Pt(bx, by)}
}

// Segment represents a line segment with start and end points
type Segment[T Number] struct {
	A Point[T]
	B Point[T]
}

// Slope answers the slope of this line segment.
func (s Segment[T]) Slope() Slope {
	if s.A.X == s.B.X {
		return VerticalSlope
	} else if s.A.Y == s.B.Y {
		return HorizontalSlope
	}
	// Account for reverse-Y coords
	x1, y1 := float64(s.A.X), float64(s.A.Y)
	x2, y2 := float64(s.B.X), float64(s.B.Y)
	return Slope{Angle: Oblique, M: (y1 - y2) / (x2 - x1)}
}

// PerpendicularSlope answers the perpendicular of Slope.
func (s Segment[T]) PerpendicularSlope() Slope {
	m := s.Slope()
	if m.Angle == Horizontal {
		return VerticalSlope
	} else if m.Angle == Vertical {
		return HorizontalSlope
	}
	return Slope{Angle: Oblique, M: -1.0 / m.M}
}

// AsArray returns my points in a slice.
func (s Segment[T]) AsArray() []Point[T] {
	return []Point[T]{
		s.A,
		s.B,
	}
}

// Len answers the length of this segment.
func (s Segment[T]) Len() T {
	return T(s.A.Dist(s.B))
}

// Dir answers the direction vector of this segment.
func (s Segment[T]) Dir() Point[T] {
	return s.B.Sub(s.A)
}

// Degrees finds the angle of the segment with A as the origin.
// Degrees will be 0-360, with 0/360 on the right, proceeding clockwise.
func (s Segment[T]) Degrees() float64 {
	return s.A.Degrees(s.B)
}

func (s Segment[T]) Midpoint() Point[T] {
	return Point[T]{X: (s.A.X / 2) + (s.B.X / 2),
		Y: (s.A.Y / 2) + (s.B.Y / 2)}
}

// Interp answers a new point at the unit position on this segment,
// where 0. = A and 1. = B.
// Note this really only works on floats, need a way to narrow that constraint.
func (s Segment[T]) Interp(unit T) Point[T] {
	return Point[T]{X: (s.A.X * (1 - unit)) + (s.B.X * unit),
		Y: (s.A.Y * (1 - unit)) + (s.B.Y * unit)}
}

// IsCollinear checks if three points are collinear
func IsCollinear[T Number](p1, p2, p3 Point[T]) bool {
	a := float64(p2.X - p1.X)
	b := float64(p3.Y - p1.Y)
	c := float64(p2.Y - p1.Y)
	d := float64(p3.X - p1.X)
	return math.Abs(a*b-c*d) < 1e-6
}

// OnSegment checks if a point lies on a line segment
func OnSegment[T Number](p, start, end Point[T]) bool {
	sx := float64(start.X)
	ex := float64(end.X)
	sy := float64(start.Y)
	ey := float64(end.Y)
	return (float64(p.X) >= math.Min(sx, ex) &&
		float64(p.X) <= math.Max(sx, ex) &&
		float64(p.Y) >= math.Min(sy, ey) &&
		float64(p.Y) <= math.Max(sy, ey))
}

// DistSquared answers the squared distance from the point to the segment,
// as well as the point found on the segment.
// From https://stackoverflow.com/questions/849211/shortest-distance-between-a-point-and-a-line-segment
func DistSquared(seg SegF, p PtF) (float64, PtF) {
	l2 := seg.A.DistSquared(seg.B)
	if l2 == 0 {
		return p.DistSquared(seg.A), seg.A
	}
	t := ((p.X-seg.A.X)*(seg.B.X-seg.A.X) + (p.Y-seg.A.Y)*(seg.B.Y-seg.A.Y)) / l2
	t = math.Max(0, math.Min(1, t))
	newP := PtF{X: seg.A.X + t*(seg.B.X-seg.A.X),
		Y: seg.A.Y + t*(seg.B.Y-seg.A.Y)}
	return p.DistSquared(newP), newP
}

// DistPointToSegment answers the distance from the point to the segment,
// as well as the point found on the segment.
// From https://stackoverflow.com/questions/849211/shortest-distance-between-a-point-and-a-line-segment
func DistPointToSegment(seg SegF, p PtF) (float64, PtF) {
	d, newP := DistSquared(seg, p)
	return math.Sqrt(d), newP
}

// DistanceFromPointToSegment calculates the distance from a point to a line segment
// Ugh, what a mess. I guess this is obsoleted by DistPointToSegment, although
// I don't see why this couldn't answer return the point. And then I suppose
// I should see if there are differenecs in result and speed to decide which to keep.
// Although looking at the implementation, it sorta looks like it's just a distance
// to the closest endpoint? Which seems completely wrong?
func DistanceFromPointToSegment[T Number](p Point[T], s Segment[T]) float64 {
	p1, p2 := s.A, s.B

	// Check if the point lies on the line segment
	if IsCollinear(p, p1, p2) && OnSegment(p, p1, p2) {
		return minDistanceToPointOnLine(p, p1, p2)
	}

	// Calculate distances to the two line endpoints
	d1 := p.Dist(p1)
	d2 := p.Dist(p2)

	// Find the closer endpoint
	closerPoint := p1
	if d2 < d1 {
		closerPoint = p2
	}

	// **Use the closer endpoint for distance calculation:**
	return p.Dist(closerPoint)
}

// XAtY answers the X value for this segment at the given Y
// value, or false if the line does not intersect y.
func XAtY(s SegF, y float64) (float64, bool) {
	if s.B.Y-s.A.Y == 0 {
		return 0, false
	}
	miny, maxy := s.A.Y, s.B.Y
	if s.B.Y < s.A.Y {
		miny, maxy = maxy, miny
	}
	if maxy-miny == 0 {
		return 0, false
	} else if !(y >= miny && y <= maxy) {
		return 0, false
	}
	return s.A.X + (((y - s.A.Y) * (s.B.X - s.A.X)) / (s.B.Y - s.A.Y)), true
}

// projectPointOnLine projects a point onto a line
func projectPointOnLine[T Number](p, p1, p2 Point[T]) Point[T] {
	v1 := Point[T]{p2.X - p1.X, p2.Y - p1.Y}
	v2 := Point[T]{p.X - p1.X, p.Y - p1.Y}
	t := dotProduct(v1, v2) / dotProduct(v1, v1)
	return Point[T]{X: T(float64(p1.X) + t*float64(v1.X)),
		Y: T(float64(p1.Y) + t*float64(v1.Y))}
}

// minDistanceToPointOnLine calculates the minimum distance from a point to a line
func minDistanceToPointOnLine[T Number](p, p1, p2 Point[T]) float64 {
	v1 := Point[T]{p2.X - p1.X, p2.Y - p1.Y}
	v2 := Point[T]{p.X - p1.X, p.Y - p1.Y}
	return math.Abs(dotProduct(v1, v2)) / math.Sqrt(dotProduct(v1, v1))
}

// dotProduct calculates the dot product of two vectors
func dotProduct[T Number](v1, v2 Point[T]) float64 {
	return float64(v1.X*v2.X + v1.Y*v2.Y)
}

func ConvertSegment[A Number, B Number](seg Segment[A]) Segment[B] {
	a := ConvertPoint[A, B](seg.A)
	b := ConvertPoint[A, B](seg.B)
	return Segment[B]{A: a, B: b}
}

type SegF = Segment[float64]
type SegI = Segment[int]

type SegF32 = Segment[float32]
type SegF64 = Segment[float64]
type SegI64 = Segment[int64]
