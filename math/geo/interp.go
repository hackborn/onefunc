package geo

func LerpPoint[T Number](a, b Point[T], amount float64) Point[T] {
	return Point[T]{X: T((float64(a.X) * (1. - amount)) + (float64(b.X) * amount)),
		Y: T((float64(a.Y) * (1. - amount)) + (float64(b.Y) * amount))}
}
