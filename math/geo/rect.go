package geo

func NewRect[T Number](left, top, right, bottom T) Rect[T] {
	return Rect[T]{LT: Point[T]{X: left, Y: top}, RB: Point[T]{X: right, Y: bottom}}
}

type Rect[T Number] struct {
	LT Point[T]
	RB Point[T]
}

func (r Rect[T]) Size() Point[T] {
	return Point[T]{X: r.RB.X - r.LT.X, Y: r.RB.Y - r.LT.Y}
}

func (p Rect[T]) Area() T {
	return (p.RB.X - p.LT.X) * (p.RB.Y - p.LT.Y)
}

type RectF32 = Rect[float32]
type RectF64 = Rect[float64]
type RectI = Rect[int]
type RectI64 = Rect[int64]
type RectUI64 = Rect[uint64]
