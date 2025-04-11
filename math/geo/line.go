package geo

// From direction vector v and offset c
func LineFromDir[T Float](v Point[T], c T) Line[T] {
	return Line[T]{V: v, C: c}
}

// From equation ax+by=c
func lineFromEq[T Float](a, b, c T) Line[T] {
	return Line[T]{V: Point[T]{X: b, Y: -a}, C: c}
}

// From points P and Q
func LineFromPts[T Float](p, q Point[T]) Line[T] {
	v := q.Sub(p)
	return Line[T]{V: v, C: v.Cross(p)}
}

// From points P and Q as XYs.
func LineFromXys[T Float](px, py, qx, qy T) Line[T] {
	p, q := Point[T]{X: px, Y: py}, Point[T]{X: qx, Y: qy}
	return LineFromPts(p, q)
}

// https://github.com/vlecomte/cp-geo/blob/master/basics/line.tex
type Line[T Float] struct {
	V Point[T]
	C T
}

// From gemini, check it out
/*
// NearestPointOnLine finds the nearest point on an infinite line to a given point.
func NearestPointOnLine(p, a, b Point) Point {
	ab := Point{X: b.X - a.X, Y: b.Y - a.Y}
	ap := Point{X: p.X - a.X, Y: p.Y - a.Y}

	dotProduct := ap.X*ab.X + ap.Y*ab.Y
	abLengthSquared := ab.X*ab.X + ab.Y*ab.Y

	if abLengthSquared == 0 {
		return a
	}

	t := dotProduct / abLengthSquared
	return Point{X: a.X + t*ab.X, Y: a.Y + t*ab.Y}
}
*/

type LnF = Line[float64]
