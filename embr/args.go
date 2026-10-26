package embr

import (
	"github.com/hackborn/onefunc/math/geo"
)

// Arg stores information about the image,
// updated on each step.
type Arg interface {
	// Setup sets any initial information about the image.
	Setup(Setup)

	// Step is updated on each pixel. Note, this might
	// go away -- you could just have functions that
	// take the current PixelArgs and it'd be the same.
	// The only value I'm seeing here is as a caching
	// mechanism, and I don't know that you'd ever need that.
	Step(Step)
}

type Setup struct {
	Width, Height int
}

type Step struct {
	X, Y int
}

// ---------------------------------------------------------
// DISTANCE-ARG

// Distance provides the distance from the center to the
// current pixel. It assumes a circular shape, and sets
// the max distance to the nearest edge.
type Distance struct {
	Distance float64
	// Aa is the antialias value if we're on the edge
	// of the circle. It will be  <= 0 if outside,
	// 1 if inside.
	Aa          float64
	MaxDistance float64
	Center      geo.PtF
	w, h        int
}

func (a *Distance) Setup(setup Setup) {
	a.w, a.h = setup.Width, setup.Height
	a.Center = geo.PtF{X: (float64(a.w) / 2) - 0, Y: (float64(a.h) / 2) - 0}
	maxX := geo.PtF{X: 0, Y: a.Center.Y}.Dist(a.Center)
	maxY := geo.PtF{X: a.Center.X, Y: 0}.Dist(a.Center)
	a.MaxDistance = min(maxX, maxY)
}

func (a *Distance) Step(step Step) {
	pt := geo.PtF{X: float64(step.X) + .5, Y: float64(step.Y) + .5}
	a.Distance = pt.Dist(a.Center)
	if a.Distance >= a.MaxDistance {
		a.Aa = 0.
	} else if a.Distance <= a.MaxDistance-1 {
		a.Aa = 1.
	} else {
		a.Aa = (a.MaxDistance - a.Distance)
	}
}
