package embr

import (
	"github.com/hackborn/onefunc/math/geo"
)

// DataFx stores information that might be shared
// across layers. It is run once per pixel.
type DataFx interface {
	// Setup sets any initial information about the data.
	Setup(DataSetupArgs)
	// Step is updated on each pixel.
	Step(StepArgs)
}

type DataSetupArgs struct {
	Width, Height int
}

type StepArgs struct {
	X, Y int
}

// ---------------------------------------------------------
// DISTANCE

// Distance provides the distance from the center to the
// current pixel. It assumes a circular shape, and sets
// the max distance to the nearest edge.
type Distance struct {
	// Distance is the distance from the current pixel
	// to the center of the image.
	Distance float64
	// Aa is the antialias value if we're on the edge
	// of the circle. It will be  <= 0 if outside,
	// 1 if inside.
	Aa          float64
	MaxDistance float64
	center      geo.PtF
	w, h        int
}

func (a *Distance) Setup(setup DataSetupArgs) {
	a.w, a.h = setup.Width, setup.Height
	a.center = geo.PtF{X: (float64(a.w) / 2) - 0, Y: (float64(a.h) / 2) - 0}
	maxX := geo.PtF{X: 0, Y: a.center.Y}.Dist(a.center)
	maxY := geo.PtF{X: a.center.X, Y: 0}.Dist(a.center)
	a.MaxDistance = min(maxX, maxY)
}

func (a *Distance) Step(step StepArgs) {
	pt := geo.PtF{X: float64(step.X) + .5, Y: float64(step.Y) + .5}
	a.Distance = pt.Dist(a.center)
	if a.Distance >= a.MaxDistance {
		a.Aa = 0.
	} else if a.Distance <= a.MaxDistance-1 {
		a.Aa = 1.
	} else {
		a.Aa = (a.MaxDistance - a.Distance)
	}
}

// ---------------------------------------------------------
// EVAL

func collapseDataFx(fxs []DataFx) DataFx {
	switch count := len(fxs); {
	case count <= 0:
		return &dfx0{}
	case count == 1:
		return fxs[0]
	// Basic unroll. I could do more of this but I can't
	// imagine it makes the slightest difference for any
	// case I'll use it for.
	case count == 2:
		return &dfx2{a: fxs[0], b: fxs[1]}
	default:
		return &dfxN{n: fxs}
	}
}

type dfx0 struct {
}

func (d dfx0) Setup(DataSetupArgs) {}

func (d dfx0) Step(StepArgs) {}

type dfx2 struct {
	a, b DataFx
}

func (d *dfx2) Setup(args DataSetupArgs) {
	d.a.Setup(args)
	d.b.Setup(args)
}

func (d *dfx2) Step(args StepArgs) {
	d.a.Step(args)
	d.b.Step(args)
}

type dfxN struct {
	n []DataFx
}

func (d *dfxN) Setup(args DataSetupArgs) {
	for _, fx := range d.n {
		fx.Setup(args)
	}
}

func (d *dfxN) Step(args StepArgs) {
	for _, fx := range d.n {
		fx.Step(args)
	}
}
