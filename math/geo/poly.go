package geo

// Pts is shorthand for creating a slice of 2D points.
func Pts[T Number](xys ...T) []Point[T] {
	pts := []Point[T]{}
	for i := 1; i < len(xys); i += 2 {
		pts = append(pts, Point[T]{X: xys[i-1], Y: xys[i]})
	}
	return pts
}

type Poly[T Number] struct {
	Pts []Point[T]
}

type PolyF64 = Poly[float64]
type PolyI = Poly[int]
type PolyI64 = Poly[int64]
type PolyUI64 = Poly[uint64]
