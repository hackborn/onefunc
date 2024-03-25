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
	bounds := RectF64{LT: PointF64{X: pt.X, Y: pt.Y},
		RB: PointF64{X: pt.X, Y: pt.Y}}
	for _, pt := range pts {
		if pt.X < bounds.LT.X {
			bounds.LT.X = pt.X
		} else if pt.X > bounds.RB.X {
			bounds.RB.X = pt.X
		}
		if pt.Y < bounds.LT.Y {
			bounds.LT.Y = pt.Y
		} else if pt.Y > bounds.RB.Y {
			bounds.RB.Y = pt.Y
		}
	}
	return bounds
}
