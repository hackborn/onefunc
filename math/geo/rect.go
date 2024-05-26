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

func (r RectT[T]) Translate(pt Point[T]) RectT[T] {
	return RectT[T]{L: r.L + pt.X, T: r.T + pt.Y, R: r.R + pt.X, B: r.B + pt.Y}
}

func (r RectT[T]) Size() Point[T] {
	return Point[T]{X: r.R - r.L, Y: r.B - r.T}
}

func (r RectT[T]) WithSize(pt Point[T]) RectT[T] {
	return RectT[T]{L: r.L, T: r.T, R: r.L + pt.X, B: r.T + pt.Y}
}

func (r RectT[T]) Area() T {
	lt := r.LT()
	rb := r.RB()
	return (rb.X - lt.X) * (rb.Y - lt.Y)
}

func (r1 RectT[T]) Union(r2 RectT[T]) RectT[T] {
	r1.L = min(r1.L, r2.L)
	r1.T = min(r1.T, r2.T)
	r1.R = max(r1.R, r2.R)
	r1.B = max(r1.B, r2.B)
	return r1
}

// Expand adds the value to all edges.
func (r RectT[T]) WithExpand(v T) RectT[T] {
	return r.Add(-v, -v, v, v)
}

func (r1 RectT[T]) Add(l, t, r, b T) RectT[T] {
	return RectT[T]{L: r1.L + l,
		T: r1.T + t,
		R: r1.R + r,
		B: r1.B + b,
	}
}

func ConvertRect[A Number, B Number](a RectT[A]) RectT[B] {
	return RectT[B]{L: B(a.L),
		T: B(a.T),
		R: B(a.R),
		B: B(a.B)}
}

type RectF32 = RectT[float32]
type RectF64 = RectT[float64]
type RectI = RectT[int]
type RectI64 = RectT[int64]
type RectUI64 = RectT[uint64]
