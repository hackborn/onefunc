package rasterizing

// PixelFunc handles a single pixel in a rasterization.
type PixelFunc func(Pixel)

// Part of a currently on-hold experiement to iterate blocks
// of pixels at a time. The point is performance, but it's not really faster.
type PixelsFunc func([]Pixel)

type Pixel struct {
	X, Y   int
	Amount float64
}
