package geo

type Range[T Number] struct {
	Min T
	Max T
}

// Clip returns the value clipped to my range.
func (p Range[T]) Clip(value T) T {
	if value <= p.Min {
		return p.Min
	} else if value >= p.Max {
		return p.Max
	} else {
		return value
	}
}

type RangeF64 = Range[float64]
type RangeI = Range[int]
type RangeI64 = Range[int64]
