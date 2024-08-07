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

type Orientation uint8

const (
	Collinear Orientation = iota
	Clockwise
	CounterClockwise
)

var (
	HorizontalSlope = Slope{Angle: Horizontal, M: 0.}
	VerticalSlope   = Slope{Angle: Vertical, M: math.MaxFloat64}
)
