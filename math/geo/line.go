package geo

import (
	"math"
)

type PixelFunc func(x, y int, amount float64)

// DrawLine draws an antialiased line using Xiaolin Wu’s algorithm.
// The output of each step is provided to the out func.
// From https://www.geeksforgeeks.org/anti-aliased-line-xiaolin-wus-algorithm/
func DrawLine2(x0, y0, x1, y1 int, out PixelFunc) {
	steep := math.Abs(float64(y1-y0)) > math.Abs(float64(x1-x0))
	// swap the co-ordinates if slope > 1 or we draw backwards
	if steep {
		x0, y0 = swap(x0, y0)
		x1, y1 = swap(x1, y1)
	}

	if x0 > x1 {
		x0, x1 = swap(x0, x1)
		y0, y1 = swap(y0, y1)
	}

	dx := float64(x1 - x0)
	dy := float64(y1 - y0)
	gradient := dy / dx
	if dx == 0.0 {
		gradient = 1.0
	}

	// handle first endpoint
	xend := round(float64(x0))
	yend := float64(y0) + gradient*float64(xend-x0)
	xgap := rfpart(float64(x0) + 0.5)
	xpxl1 := xend // this will be used in the main loop
	ypxl1 := ipart(yend)
	if steep {
		out(ypxl1, xpxl1, rfpart(yend)*xgap)
		out(ypxl1+1, xpxl1, fpart(yend)*xgap)
	} else {
		out(xpxl1, ypxl1, rfpart(yend)*xgap)
		out(xpxl1, ypxl1+1, fpart(yend)*xgap)
	}
	intery := yend + gradient // first y-intersection for the main loop

	// handle second endpoint
	xend = round(float64(x1))
	yend = float64(y1) + gradient*float64(xend-x1)
	xgap = fpart(float64(x1) + 0.5)
	xpxl2 := xend //this will be used in the main loop
	ypxl2 := ipart(yend)
	if steep {
		out(ypxl2, xpxl2, rfpart(yend)*xgap)
		out(ypxl2+1, xpxl2, fpart(yend)*xgap)
	} else {
		out(xpxl2, ypxl2, rfpart(yend)*xgap)
		out(xpxl2, ypxl2+1, fpart(yend)*xgap)
	}

	// main loop
	if steep {
		for x := xpxl1 + 1; x <= xpxl2-1; x++ {
			out(ipart(intery), x, rfpart(intery))
			out(ipart(intery)+1, x, fpart(intery))
			intery = intery + gradient
		}
	} else {
		for x := xpxl1 + 1; x < xpxl2-1; x++ {
			out(x, ipart(intery), rfpart(intery))
			out(x, ipart(intery)+1, fpart(intery))
			intery = intery + gradient
		}
	}
}

// integer part of x
func ipart(x float64) int {
	return int(math.Floor(x))
}

func round(x float64) int {
	return ipart(x + 0.5)
}

// fractional part of x
func fpart(x float64) float64 {
	return x - float64(ipart(x))
}

func rfpart(x float64) float64 {
	return 1.0 - fpart(x)
}

// DrawLine draws an antialiased line using Xiaolin Wu’s algorithm.
// The output of each step is provided to the out func.
// From https://www.geeksforgeeks.org/anti-aliased-line-xiaolin-wus-algorithm/
func DrawLine(x0, y0, x1, y1 int, out PixelFunc) {
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
				out(iPartOfNumber(intersectY), x, amount)
			}
			if amount := fPartOfNumberClamped(intersectY); amount > 0.0 {
				out(iPartOfNumber(intersectY)+1, x, amount)
			}
			intersectY += gradient
		}
	} else {
		for x := xpxl1; cmp(x, xpxl2); x += inc {
			// pixel coverage is determined by fractional
			// part of y co-ordinate
			if amount := rfPartOfNumberClamped(intersectY); amount > 0.0 {
				out(x, iPartOfNumber(intersectY), amount)
			}
			if amount := fPartOfNumberClamped(intersectY); amount > 0.0 {
				out(x, iPartOfNumber(intersectY)+1, amount)
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
