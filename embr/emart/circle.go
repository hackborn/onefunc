package emart

import (
	"image/color"

	"github.com/hackborn/onefunc/embr"
	"github.com/hackborn/onefunc/math/geo"
)

// ---------------------------------------------------------
// MACROS

func NewCircle(outsideColor, insideColor color.NRGBA) *Circle {
	return &Circle{OutsideColor: outsideColor,
		InsideColor: insideColor,
	}
}

func NewCircleOutline(c color.NRGBA, width float64) *CircleOutline {
	return &CircleOutline{Color: c, Width: width}
}

// Circle draws a circle at the size of the image.
type Circle struct {
	OutsideColor color.NRGBA
	InsideColor  color.NRGBA

	distance *embr.Distance
}

func (c *Circle) Setup(args embr.SetupArgs) error {
	c.distance = args.MustStep(&embr.Distance{})
	return nil
}

func (c *Circle) Pixel(embr.PixelArgs) color.NRGBA {
	if c.distance.Distance < c.distance.MaxDistance {
		if c.distance.Aa >= 1. {
			return c.InsideColor
		}
		return embr.Interp(c.OutsideColor, c.InsideColor, 1.-c.distance.Aa)
	} else {
		return c.OutsideColor
	}
}

// CircleOutline draws an outline of a circle at the size of the
// image. It does not paint all pixels.
type CircleOutline struct {
	Color color.NRGBA
	Width float64

	distance *embr.Distance
}

func (c *CircleOutline) Setup(args embr.SetupArgs) error {
	c.distance = args.MustStep(&embr.Distance{})
	return nil
}

func (c *CircleOutline) Pixel(args embr.PixelArgs) color.NRGBA {
	// TODO: Be nice if there was a cache step for this,
	// but I don't have any setup that happens after the Stepper set up.
	rng := geo.RngF{Min: c.distance.MaxDistance - c.Width, Max: c.distance.MaxDistance}
	if rng.Contains(c.distance.Distance) {
		if c.distance.Distance > rng.Max-1 {
			// Antialias outside edge
			aa := c.distance.Distance - (rng.Max - 1)
			return embr.Interp(args.Color, c.Color, aa)
		} else if c.distance.Distance < rng.Min+1 {
			// Antialias inside edge
			aa := (rng.Min + 1) - c.distance.Distance
			return embr.Interp(args.Color, c.Color, aa)
		} else {
			return c.Color
		}
	} else {
		return args.Color
	}
}
