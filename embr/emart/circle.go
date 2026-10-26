package emart

import (
	"image/color"

	"github.com/hackborn/onefunc/embr"
)

func NewCircle(outsideColor, insideColor color.NRGBA) *Circle {
	return &Circle{OutsideColor: outsideColor,
		InsideColor: insideColor,
	}
}

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
