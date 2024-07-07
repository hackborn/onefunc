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

func PtToSegIntersection(pt, direction PtF, seg SegF) (PtF, bool) {
	// Create a new segment and check the intersection.
	scale := Max(math.Abs(pt.X-seg.A.X), math.Abs(pt.X-seg.B.X))
	scale = Max(scale, math.Abs(pt.Y-seg.A.Y))
	scale = Max(scale, math.Abs(pt.Y-seg.B.Y))
	scale *= 10

	newPt := PtF{X: pt.X + direction.X*scale, Y: pt.Y + direction.Y*scale}
	return FindIntersection(seg, SegF{A: pt, B: newPt})
}

// https://stackoverflow.com/questions/34415671/intersection-of-a-line-with-a-line-segment-in-c
// line segment p-q intersect with line A-B.
// This looks a lot more efficient than the current seg-seg test I'm doing,
// but also I actually want a seg-seg test because this gives false positives.
// So... maybe I'll switch to it if I think of a way for it to make sense,
// but probably not because this is just a different operation.
func LineToSegIntersection(A, B PtF, pq SegF) (PtF, bool) {
	a := B.Y - A.Y
	b := A.X - B.X
	c := B.X*A.Y - A.X*B.Y
	u := math.Abs(a*pq.A.X + b*pq.A.Y + c)
	v := math.Abs(a*pq.B.X + b*pq.B.Y + c)
	newPt := Pt((pq.A.X*v+pq.B.X*u)/(u+v), (pq.A.Y*v+pq.B.Y*u)/(u+v))
	return newPt, true
}
