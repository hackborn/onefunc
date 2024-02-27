package geo

import (
	"math"
)

// Segment represents a line segment with start and end points
type Segment[T Number] struct {
	A Point[T]
	B Point[T]
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

// Orientation checks the orientation of three points
func Orientation[T Number](p0, p1, p2 Point[T]) int {
	val := (p1.Y-p0.Y)*(p2.X-p1.X) - (p1.X-p0.X)*(p2.Y-p1.Y)
	if val == 0 {
		return 0 // collinear
	} else if val > 0 {
		return 1 // clockwise
	} else {
		return -1 // counter-clockwise
	}
}

// FindIntersection finds the intersection point of two line segments
func FindIntersection[T Number](s1, s2 Segment[T]) (Point[T], bool) {
	p1, p2 := s1.A, s1.B
	q1, q2 := s2.A, s2.B

	// Check if the lines are parallel
	o1 := Orientation(p1, q1, p2)
	o2 := Orientation(p1, q2, p2)
	if o1 == 0 && o2 == 0 {
		// Check if they are collinear
		if OnSegment(p1, q1, q2) || OnSegment(p2, q1, q2) || OnSegment(q1, p1, p2) || OnSegment(q2, p1, p2) {
			// Line segments are overlapping
			x := math.Min(float64(p1.X), math.Max(float64(p1.X), float64(p2.X)))
			y := math.Min(float64(p1.Y), math.Max(float64(p1.Y), float64(p2.Y)))
			return Point[T]{X: T(x), Y: T(y)}, true
		}
		return Point[T]{}, false // Lines are parallel and don't overlap
	}

	// Check if the segments intersect
	o3 := Orientation(q1, p1, q2)
	o4 := Orientation(q2, p1, q1)

	if o1 != o2 && o3 != o4 {
		// Lines intersect
		x := (q1.X*q2.Y - q1.Y*q2.X - p1.X*p2.Y + p1.Y*p2.X) / (q2.X - q1.X)
		y := (q1.X*p2.Y + p1.X*q1.Y - q1.Y*p2.X - p1.Y*q1.X) / (q2.Y - q1.Y)
		return Point[T]{x, y}, true
	}

	return Point[T]{}, false // Lines don't intersect
}

// DistanceFromPointToLine calculates the distance from a point to a line segment
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

type SegmentF64 = Segment[float64]
type SegmentI = Segment[int]
type SegmentI64 = Segment[int64]
type SegmentUI64 = Segment[uint64]
