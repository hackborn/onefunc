package geo

func Rng[T Number](min, max T) Range[T] {
	return Range[T]{Min: min, Max: max}
}

type Range[T Number] struct {
	Min T
	Max T
}

// Clip returns the value clipped to my range.
func (p Range[T]) Clip(value T) T {
	min, max := p.Min, p.Max
	if p.Max < p.Min {
		min, max = p.Max, p.Min
	}

	if value <= min {
		return min
	} else if value >= max {
		return max
	} else {
		return value
	}
}

// Map returns the normalized value mapped to my range.
func (p Range[T]) Map(value T) T {
	if value < 0 {
		value = 0
	} else if value > 1 {
		value = 1
	}
	return ((1.0 - value) * p.Min) + (value * p.Max)
}

type RangeF64 = Range[float64]
type RangeI = Range[int]
type RangeI64 = Range[int64]
