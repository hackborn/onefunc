package geo

func PolyBB[T Float](pts []Point[T]) PolygonBB[T] {
	return PolygonBB[T]{Pts: pts,
		BB: PolygonBounds(pts)}
}

// PolygonBB is a closed polygon that includes
// the computed bounding box.
type PolygonBB[T Number] struct {
	Pts []Point[T]
	BB  Rectangle[T]
}

type PolygonBBF64 = PolygonBB[float64]

func PolygonBounds[T Float](ptss ...[]Point[T]) Rectangle[T] {
	bounds := Rectangle[T]{}
	for i, pts := range ptss {
		for ii, pt := range pts {
			if i == 0 && ii == 0 {
				bounds = Rect(pt.X, pt.Y, pt.X, pt.Y)
			} else {
				bounds.L = Min(bounds.L, pt.X)
				bounds.T = Min(bounds.T, pt.Y)
				bounds.R = Max(bounds.R, pt.X)
				bounds.B = Max(bounds.B, pt.Y)
			}
		}
	}
	return bounds
}
