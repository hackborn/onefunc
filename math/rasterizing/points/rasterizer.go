package points

import (
	"github.com/hackborn/onefunc/math/geo"
	"github.com/hackborn/onefunc/math/rasterizing"
)

// NewRasterizer answers a new rasterizer for drawing points.
// Args:
// radius: Radius (in pixels) of each rasterized point.
// Accepts:
// []geo.PtF
func NewRasterizer(radius float64) rasterizing.Rasterizer {
	return &rasterizer{radius: radius, kernel: makeKernel(radius)}
}

type rasterizer struct {
	radius float64
	kernel []kernelValue
}

func (r *rasterizer) Rasterize(shape any, fn rasterizing.PixelFunc) {
	pts := r.getPts(shape)
	for _, pt := range pts {
		pti := geo.ConvertPoint[float64, int](pt)
		for _, k := range r.kernel {
			args := rasterizing.Pixel{X: pti.X + k.x, Y: pti.Y + k.y, Amount: k.value}
			fn(args)
		}
	}
}

func (r *rasterizer) getPts(shape any) []geo.PtF {
	switch s := shape.(type) {
	case []geo.PtF:
		return s
	}
	return nil
}

func makeKernel(radius float64) []kernelValue {
	ri := int(radius)
	if ri < 1 {
		return []kernelValue{{0, 0, 1.}}
	}
	kv := make([]kernelValue, ri*ri)
	idx := -1
	cen := geo.Pt(0., 0.)
	for y := -ri; y <= ri; y++ {
		for x := -ri; x <= ri; x++ {
			idx++
			value := 0.
			d := cen.Dist(geo.Pt(float64(x), float64(y)))
			if d < radius {
				value = 1. - (d / radius)
			}
			kv = append(kv, kernelValue{x: x, y: y, value: value})
		}
	}
	return kv
}

type kernelValue struct {
	x, y  int
	value float64
}
