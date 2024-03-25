package rasterizing

type Rasterizer interface {
	// Rasterize the supplied shape to pixels.
	// shape can be any type the rasterizer understands,
	// but default types are:
	// * geo.SegmentI for lines
	// * []geo.PointF64 for polygons
	Rasterize(shape any, fn PixelFunc)
}
