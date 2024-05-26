package geo

func PolyBB[T Float](pts []Point[T]) PolygonBB[T] {
	return PolygonBB[T]{Pts: pts,
		BB: PolygonBounds(pts)}
}

// PolygonBB is a closed polygon that includes
// the computed bounding box.
type PolygonBB[T Number] struct {
	Pts []Point[T]
	BB  RectT[T]
}

type PolygonBBF64 = PolygonBB[float64]

func PolygonBounds[T Float](pts []Point[T]) RectT[T] {
	if len(pts) < 1 {
		return RectT[T]{}
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
