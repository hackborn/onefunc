package geo

type Poly[T Number] struct {
	Pts []Point[T]
}

type PolyF64 = Poly[float64]
type PolyI = Poly[int]
type PolyI64 = Poly[int64]
type PolyUI64 = Poly[uint64]

func PolygonBounds(pts []PointF64) RectF64 {
	if len(pts) < 1 {
		return RectF64{}
	}
	pt := pts[0]
	bounds := Rect(pt.X, pt.Y, pt.X, pt.Y)
	for _, pt := range pts {
		if pt.X < bounds.L {
			bounds.L = pt.X
		} else if pt.X > bounds.R {
			bounds.R = pt.X
		}
		if pt.Y < bounds.T {
			bounds.T = pt.Y
		} else if pt.Y > bounds.B {
			bounds.B = pt.Y
		}
	}
	return bounds
}
