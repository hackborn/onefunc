package geo

// ---------------------------------------------------------
// QUADRATIC BEZIER

// QuadraticBezier3D is a 3D quadratic Bézier curve.
// defined by three control points.
type QuadraticBezier3d struct {
	P0, P1, P2 Pt3dF
}

// PointAt evaluates the Bezier curve at a given parameter t (0 <= t <= 1).
// It satisfies the PointInterpolator interface.
func (q QuadraticBezier3d) PointAt(t float64) Pt3dF {
	// TODO: Should I clamp t?

	omt := 1 - t
	omtSquared := omt * omt
	tSquared := t * t
	return Pt3dF{
		X: omtSquared*q.P0.X + 2*omt*t*q.P1.X + tSquared*q.P2.X,
		Y: omtSquared*q.P0.Y + 2*omt*t*q.P1.Y + tSquared*q.P2.Y,
		Z: omtSquared*q.P0.Z + 2*omt*t*q.P1.Z + tSquared*q.P2.Z,
	}
}
