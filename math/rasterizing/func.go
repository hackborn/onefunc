package rasterizing

// PixelFunc handles a single pixel in a rasterization.
type PixelFunc func(PixelArgs)

type PixelArgs struct {
	X, Y   int
	Amount float64
}
