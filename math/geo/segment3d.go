package geo

// Seg is shorthand for creating a segment from two points.
func Seg3d[T Number](ax, ay, az, bx, by, bz T) Segment3d[T] {
	return Segment3d[T]{A: Pt3d(ax, ay, az), B: Pt3d(bx, by, bz)}
}

// Segment represents a line segment with start and end points
type Segment3d[T Number] struct {
	A Point3d[T]
	B Point3d[T]
}
