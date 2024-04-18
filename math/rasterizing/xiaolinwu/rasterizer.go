package xiaolinwu

import (
	"math"

	"github.com/hackborn/onefunc/math/geo"
	"github.com/hackborn/onefunc/math/rasterizing"
)

func NewRasterizer() rasterizing.Rasterizer {
	return &rasterizer{}
}

type rasterizer struct {
}

// Rasterize draws an antialiased line using Xiaolin Wu’s algorithm.
// The output of each step is provided to the out func.
// Modified from https://www.geeksforgeeks.org/anti-aliased-line-xiaolin-wus-algorithm/
func (r *rasterizer) Rasterize(shape any, fn rasterizing.PixelFunc) {
	x0, y0, x1, y1, ok := r.getSegment(shape)
	if !ok {
		return
	}

	steep := math.Abs(float64(y1-y0)) > math.Abs(float64(x1-x0))
	inc := 1
	cmp := func(a, b int) bool {
		return a <= b
	}
	// swap the co-ordinates if slope > 1 or we draw backwards
	if steep {
		x0, y0 = swap(x0, y0)
		x1, y1 = swap(x1, y1)
	}
	if x0 > x1 {
		inc = -1
		cmp = func(a, b int) bool {
			return a >= b
		}
	}

	// compute the slope
	dx := float64(x1 - x0)
	dy := float64(y1 - y0)
	gradient := dy / dx
	if dx == 0.0 {
		gradient = 1.0
	}
	if inc < 0 {
		gradient = -gradient
	}

	xpxl1 := x0
	xpxl2 := x1
	intersectY := float64(y0)

	// main loop
	if steep {
		for x := xpxl1; cmp(x, xpxl2); x += inc {
			// pixel coverage is determined by fractional
			// part of y co-ordinate
			if amount := rfPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: iPartOfNumber(intersectY), Y: x, Amount: amount}
				fn(args)
			}
			if amount := fPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: iPartOfNumber(intersectY) + 1, Y: x, Amount: amount}
				fn(args)
			}
			intersectY += gradient
		}
	} else {
		for x := xpxl1; cmp(x, xpxl2); x += inc {
			// pixel coverage is determined by fractional
			// part of y co-ordinate
			if amount := rfPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: x, Y: iPartOfNumber(intersectY), Amount: amount}
				fn(args)
			}
			if amount := fPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: x, Y: iPartOfNumber(intersectY) + 1, Amount: amount}
				fn(args)
			}
			intersectY += gradient
		}
	}
}

func (r *rasterizer) getSegment(shape any) (int, int, int, int, bool) {
	switch s := shape.(type) {
	case geo.SegmentI:
		return s.A.X, s.A.Y, s.B.X, s.B.Y, true
	case geo.SegmentF64:
		return int(s.A.X), int(s.A.Y), int(s.B.X), int(s.B.Y), true

	}
	return 0, 0, 0, 0, false
}

// DrawLine draws an antialiased line using Xiaolin Wu’s algorithm.
// The output of each step is provided to the out func.
// From https://www.geeksforgeeks.org/anti-aliased-line-xiaolin-wus-algorithm/
func DrawLine(x0, y0, x1, y1 int, out rasterizing.PixelFunc) {
	steep := math.Abs(float64(y1-y0)) > math.Abs(float64(x1-x0))
	inc := 1
	cmp := func(a, b int) bool {
		return a <= b
	}
	// swap the co-ordinates if slope > 1 or we draw backwards
	if steep {
		x0, y0 = swap(x0, y0)
		x1, y1 = swap(x1, y1)
	}
	if x0 > x1 {
		inc = -1
		cmp = func(a, b int) bool {
			return a >= b
		}
	}

	// compute the slope
	dx := float64(x1 - x0)
	dy := float64(y1 - y0)
	gradient := dy / dx
	if dx == 0.0 {
		gradient = 1.0
	}
	if inc < 0 {
		gradient = -gradient
	}

	xpxl1 := x0
	xpxl2 := x1
	intersectY := float64(y0)

	// main loop
	if steep {
		for x := xpxl1; cmp(x, xpxl2); x += inc {
			// pixel coverage is determined by fractional
			// part of y co-ordinate
			if amount := rfPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: iPartOfNumber(intersectY), Y: x, Amount: amount}
				out(args)
			}
			if amount := fPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: iPartOfNumber(intersectY) + 1, Y: x, Amount: amount}
				out(args)
			}
			intersectY += gradient
		}
	} else {
		for x := xpxl1; cmp(x, xpxl2); x += inc {
			// pixel coverage is determined by fractional
			// part of y co-ordinate
			if amount := rfPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: x, Y: iPartOfNumber(intersectY), Amount: amount}
				out(args)
			}
			if amount := fPartOfNumberClamped(intersectY); amount > 0.0 {
				args := rasterizing.PixelArgs{X: x, Y: iPartOfNumber(intersectY) + 1, Amount: amount}
				out(args)
			}
			intersectY += gradient
		}
	}
}

// swaps two numbers
func swap(a, b int) (int, int) {
	return b, a
}

// returns integer part of a floating point number
func iPartOfNumber(x float64) int {
	return int(x)
}

// rounds off a number
func roundNumber(x float64) int {
	return iPartOfNumber(x + 0.5)
}

// returns fractional part of a number
func fPartOfNumberClamped(x float64) float64 {
	v := fPartOfNumber(x)
	if v > 1.0 {
		return 1.0
	}
	return v
}

// returns fractional part of a number
func fPartOfNumber(x float64) float64 {
	if x > 0 {
		return x - float64(iPartOfNumber(x))
	} else {
		return x - float64(iPartOfNumber(x)+1)
	}
}

// returns 1 - fractional part of number
func rfPartOfNumberClamped(x float64) float64 {
	v := rfPartOfNumber(x)
	if v > 1.0 {
		return 1.0
	}
	return v
}

// returns 1 - fractional part of number
func rfPartOfNumber(x float64) float64 {
	return 1 - fPartOfNumber(x)
}
