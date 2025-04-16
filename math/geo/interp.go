package geo

func LerpPoint[T Number](a, b Point[T], amount float64) Point[T] {
	return Point[T]{X: T((float64(a.X) * (1. - amount)) + (float64(b.X) * amount)),
		Y: T((float64(a.Y) * (1. - amount)) + (float64(b.Y) * amount))}
}

func LerpPointXY[T Number](a, b Point[T], xAmount, yAmount float64) Point[T] {
	return Point[T]{X: T((float64(a.X) * (1. - xAmount)) + (float64(b.X) * xAmount)),
		Y: T((float64(a.Y) * (1. - yAmount)) + (float64(b.Y) * yAmount))}
}

func LerpPoint3d[T Number](a, b Point3d[T], amount float64) Point3d[T] {
	return Point3d[T]{X: T((float64(a.X) * (1. - amount)) + (float64(b.X) * amount)),
		Y: T((float64(a.Y) * (1. - amount)) + (float64(b.Y) * amount)),
		Z: T((float64(a.Z) * (1. - amount)) + (float64(b.Z) * amount))}
}

func LerpPoint3dXYZ[T Number](a, b Point3d[T], xAmount, yAmount, zAmount float64) Point3d[T] {
	return Point3d[T]{X: T((float64(a.X) * (1. - xAmount)) + (float64(b.X) * xAmount)),
		Y: T((float64(a.Y) * (1. - yAmount)) + (float64(b.Y) * yAmount)),
		Z: T((float64(a.Z) * (1. - zAmount)) + (float64(b.Z) * zAmount))}
}
