package geo

// Tri3d is shorthand for creating a triangle.
func Tri3d[T Number](a, b, c Point3d[T]) Triangle3d[T] {
	return Triangle3d[T]{A: a, B: b, C: c}
}

type Triangle3d[T Number] struct {
	A Point3d[T]
	B Point3d[T]
	C Point3d[T]
}

type Tri3dF = Triangle3d[float64]
