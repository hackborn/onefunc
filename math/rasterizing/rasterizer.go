package rasterizing

type Rasterizer interface {
	// Rasterize the supplied shape to pixels.
	// shape can be any type the rasterizer understands,
	// but default types are:
	// * geo.SegmentI for lines
	// * []geo.PointF64 for polygons
	Rasterize(shape any, fn PixelFunc)
}

// This was supposed to be an optimization -- buffer up events
// and process them in batches. It is faster, but by like 0.25%.
// Shockingly, disappointingly similar. Not enough I think to use
// at this time, unless I can find a way to speed it up more.
type Rasterizer2 interface {
	// Rasterize the supplied shape to pixels.
	// shape can be any type the rasterizer understands,
	// but default types are:
	// * geo.SegmentI for lines
	// * []geo.PointF64 for polygons
	Rasterize2(shape any, fn PixelsFunc)
}
