package geo

type Rect[T Number] struct {
	LT Point[T]
	RB Point[T]
}

func (p Rect[T]) Area() T {
	return (p.RB.X - p.LT.X) * (p.RB.Y - p.LT.Y)
}

type RectF64 = Rect[float64]

type RectI = Rect[int]

type RectI64 = Rect[int64]

type RectUI64 = Rect[uint64]