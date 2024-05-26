package geo

type Poly[T Number] struct {
	Pts []Point[T]
}

type PolyF64 = Poly[float64]
type PolyI = Poly[int]
type PolyI64 = Poly[int64]
type PolyUI64 = Poly[uint64]
