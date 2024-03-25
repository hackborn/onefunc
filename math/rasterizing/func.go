package rasterizing

// PixelFunc handles a single pixel in a rasterization.
type PixelFunc func(x, y int, amount float64)
