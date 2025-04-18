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

type SignedNumber interface {
	constraints.Signed | constraints.Float
}

type PointInterpolator interface {
	PointAt(float64) PtF
}

type PointReader interface {
	ReadPoint() (PtF, error)
}

type PointInterpolator3d interface {
	PointAt(float64) Pt3dF
}

type PointReader3d interface {
	ReadPoint() (Pt3dF, error)
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

var RngFUnit = RngF{Min: 0., Max: 1.}
