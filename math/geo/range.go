package geo

func Rng[T Number](min, max T) Range[T] {
	return Range[T]{Min: min, Max: max}
}

type Range[T Number] struct {
	Min T
	Max T
}

// Contains returns true if the value is contained in the range.
func (p Range[T]) Contains(value T) bool {
	min, max := p.Min, p.Max
	if p.Max < p.Min {
		min, max = p.Max, p.Min
	}
	return value >= min && value <= max
}

// Clip returns the value clipped to my range.
// TODO: Clamp() seems more common so switch to that.
func (p Range[T]) Clip(value T) T {
	return p.Clamp(value)
}

// Clamp returns the value clipped to my range.
func (p Range[T]) Clamp(value T) T {
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

// ClampFast returns the value clipped to my range.
// It assumes Min is less than Max.
func (r Range[T]) ClampFast(value T) T {
	if value <= r.Min {
		return r.Min
	} else if value >= r.Max {
		return r.Max
	} else {
		return value
	}
}

// Midpoint returns the center of the range.
func (p Range[T]) Midpoint() T {
	return T((p.Min / 2) + (p.Max / 2))
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
	if p.Min == p.Max {
		return float64(p.Min)
	}
	min, max := p.Min, p.Max
	if p.Max < p.Min {
		min, max = p.Max, p.Min
	}
	//	fmt.Println("value", value, "min", min, "max", max)
	var v float64
	if value <= min {
		v = 0.0
	} else if value >= max {
		v = 1.0
	} else {
		v = float64(value-min) / float64(max-min)
	}
	if p.Max < p.Min {
		return 1. - v
	}
	return v
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
