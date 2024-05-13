package geo

func Rect[T Number](left, top, right, bottom T) RectT[T] {
	return RectT[T]{L: left, T: top, R: right, B: bottom}
}

type RectT[T Number] struct {
	L, T, R, B T
}

func (r RectT[T]) LT() Point[T] {
	return Pt(r.L, r.T)
}

func (r RectT[T]) RB() Point[T] {
	return Pt(r.R, r.B)
}

func (r RectT[T]) Add(pt Point[T]) RectT[T] {
	return RectT[T]{L: r.L + pt.X, T: r.T + pt.Y, R: r.R + pt.X, B: r.B + pt.Y}
}

func (r RectT[T]) Size() Point[T] {
	lt := r.LT()
	rb := r.RB()
	return Point[T]{X: rb.X - lt.X, Y: rb.Y - lt.Y}
}

func (r RectT[T]) Area() T {
	lt := r.LT()
	rb := r.RB()
	return (rb.X - lt.X) * (rb.Y - lt.Y)
}

type RectF32 = RectT[float32]
type RectF64 = RectT[float64]
type RectI = RectT[int]
type RectI64 = RectT[int64]
type RectUI64 = RectT[uint64]
