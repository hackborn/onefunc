package geo

type Point3d[T Number] struct {
	X T
	Y T
}

func (p Point3d[T]) Area() T {
	return p.X * p.Y
}

type Point3dF64 = Point3d[float64]

type Point3dI = Point3d[int]

type Point3dI64 = Point3d[int64]

type Point3dUI64 = Point3d[uint64]
