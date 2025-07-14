package geo

func Pt3Equal(a, b Pt3dF) bool {
	return FloatsEqual(a.X, b.X) && FloatsEqual(a.Y, b.Y) && FloatsEqual(a.Z, b.Z)
}

func Pt3EquaTol(a, b Pt3dF, tol float64) bool {
	return FloatsEqualTol(a.X, b.X, tol) && FloatsEqualTol(a.Y, b.Y, tol) && FloatsEqualTol(a.Z, b.Z, tol)
}
