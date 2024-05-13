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

// Normalize returns the value clipped to my range and normalized to 0-1.
func (p Range[T]) Normalize(value T) float64 {
	min, max := p.Min, p.Max
	if p.Max < p.Min {
		min, max = p.Max, p.Min
	}
	//	fmt.Println("value", value, "min", min, "max", max)
	if min == max {
		return float64(min)
	} else if value <= min {
		return 0.0
	} else if value >= max {
		return 1.0
	} else {
		return float64(value-min) / float64(max-min)
	}
}

// MapNormal takes a normalized (0-1) value and maps it to my range.
func (p Range[T]) MapNormal(normal T) T {
	if normal < 0 {
		normal = 0
	} else if normal > 1 {
		normal = 1
	}
	return ((1.0 - normal) * p.Min) + (normal * p.Max)
}

type RangeF64 = Range[float64]
type RangeI = Range[int]
type RangeI64 = Range[int64]
