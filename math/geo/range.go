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

// Overlaps returns true if the ranges overlap.
func (a Range[T]) Overlaps(b Range[T]) bool {
	min, max := a.Min, a.Max
	if min > max {
		min, max = max, min
	}
	return b.Min >= min && b.Min <= max
}

// Intersection returns the intersection of A and B ranges.
func (a Range[T]) Intersection(b Range[T]) Range[T] {
	a.Min = Max(a.Min, b.Min)
	a.Max = Min(a.Max, b.Max)
	return a
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
func (p Range[T]) MapNormal(normal float64) T {
	if normal < 0 {
		normal = 0
	} else if normal > 1 {
		normal = 1
	}
	ans := T(((1.0 - normal) * float64(p.Min)) + (normal * float64(p.Max)))
	// Don't think we need to clip this but not sure
	//	min, max := p.Min, p.Max
	//	if min > max {
	//	min, max = max, min
	//	if ans <= min {
	//		return min
	//	} else if ans >= max {
	//		return max
	//	}
	return ans
}

type RangeF64 = Range[float64]
type RangeI = Range[int]
type RangeI64 = Range[int64]
