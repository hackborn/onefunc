package geo

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Float interface {
	constraints.Float
}

type Number interface {
	constraints.Integer | constraints.Float
}

type Angle uint8

const (
	Oblique Angle = iota
	Horizontal
	Vertical
)

var (
	HorizontalSlope = Slope{Angle: Horizontal, M: 0.}
	VerticalSlope   = Slope{Angle: Vertical, M: math.MaxFloat64}
)
