package geo

type CircleHitTest struct {
	center        PointF64
	radiusSquared float64
	bb            RectF64
}

func (t *CircleHitTest) Set(center PointF64, radius float64) {
	t.center = center
	t.radiusSquared = radius * radius
	t.bb = Rect(center.X-radius, center.Y-radius, center.X+radius, center.Y+radius)
}

func (t *CircleHitTest) Hit(pt PointF64) bool {
	if pt.X < t.bb.L || pt.Y < t.bb.T || pt.X > t.bb.R || pt.Y > t.bb.B {
		return false
	}
	return t.center.DistSquared(pt) <= t.radiusSquared
}
