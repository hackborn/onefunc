package rasterizing

// PixelFunc handles a single pixel in a rasterization.
type PixelFunc func(Pixel)

// PixelAaFunc handles a single pixel with antialiasing.
type PixelAaFunc func(PixelAa)

// Part of a currently on-hold experiement to iterate blocks
// of pixels at a time. The point is performance, but it's not really faster.
type PixelsFunc func([]Pixel)

type Pixel struct {
	X, Y   int
	Amount float64
}

type PixelAa struct {
	X, Y int

	// Amount is the amount value for this pixel, 0 - 1. It will
	// typically by 1 unless an algorithm is doing something different.
	Amount float64

	// Aa is the antialiasing amount to apply to this pixel, 0 - 1.
	// If you want to use AA, multiply the Amount by this value to
	// get the final value.
	Aa float64
}
