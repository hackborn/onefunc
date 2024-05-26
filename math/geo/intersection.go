package geo

import (
	"math"
)

// https://github.com/vlecomte/cp-geo/blob/master/basics/segment.tex
func FindIntersection[T Number](s1, s2 Segment[T]) (Point[T], bool) {
	oa := orient(s2.A, s2.B, s1.A)
	ob := orient(s2.A, s2.B, s1.B)
	oc := orient(s1.A, s1.B, s2.A)
	od := orient(s1.A, s1.B, s2.B)
	// Proper intersection exists if opposite signs
	if oa*ob < 0 && oc*od < 0 {
		pt := Point[T]{X: (s1.A.X*ob - s1.B.X*oa) / (ob - oa),
			Y: (s1.A.Y*ob - s1.B.Y*oa) / (ob - oa)}
		return pt, true
	}
	return Point[T]{}, false
}

// PerpendicularIntersection finds the intersection of point to the line
// segment by drawing a perpendicular line from point to segment.
// Note: this is from Gemini. Seems to work but not heavily tested.
func PerpendicularIntersection[T Float](seg Segment[T], pt Point[T]) (Point[T], bool) {
	// Check for vertical line (infinite slope)
	if seg.A.X == seg.B.X {
		// Point is on the line if px == x1, return the point itself
		if pt.X == seg.A.X {
			return pt, true
		} else {
			// Point is not on the line, return None (or indicate error)
			return Point[T]{}, false
		}
	}

	// Calculate the slope of the line
	m := (seg.B.Y - seg.A.Y) / (seg.B.X - seg.A.X)

	// Calculate the negative reciprocal of the slope (perpendicular slope)
	m_perp := -1 / m

	// Calculate the x-coordinate of the intersection point (ix)
	ix := (m*pt.X - pt.Y + m_perp*seg.A.X - m_perp*seg.A.Y) / (m - m_perp)

	// Calculate the y-coordinate of the intersection point (iy) using either the point or line equation
	iy := m_perp*(ix-pt.X) + pt.Y // Using point equation

	// Return the intersection point
	return Pt(ix, iy), true
}

/*
Actually, the answers everyone have given you so far are not optimal. They are imprecise, and so are not guaranteed to work on integer coordinates. Also, they are way too complicated.

Taken from Victor Lecomte's fabulous handbook, and modified for simplicity, this C++ function properInter returns whether there is an intersection between segments AB and CD:

struct pt { int x, int y };

int cross(pt a, pt b) {
    return a.x*b.y - a.y*b.x;
}

int orient(pt a, pt b, pt c) {
    return cross(b-a, c-a);
}

bool properInter(pt a, pt b, pt c, pt d) {
    int oa = orient(c,d,a),
        ob = orient(c,d,b),
        oc = orient(a,b,c),
        od = orient(a,b,d);
    // Proper intersection exists iff opposite signs
    return (oa*ob < 0 && oc*od < 0);
}
*/
// FindIntersection finds the intersection point of two line segments
func FindIntersectionBAD[T Number](s1, s2 Segment[T]) (Point[T], bool) {
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
