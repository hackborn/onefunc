package embr

import (
	"image"
	"image/draw"
)

// convertNRGBAToRGBA does what it says. I convenience for clients.
func convertNRGBAToRGBA(nrgbaImg *image.NRGBA) *image.RGBA {
	// 1. Get the bounds of the original image
	bounds := nrgbaImg.Bounds()

	// 2. Initialize a new RGBA canvas with the same bounds
	rgbaImg := image.NewRGBA(bounds)

	// 3. Draw the NRGBA image onto the RGBA canvas.
	// This automatically handles the alpha-premultiplication conversion per pixel.
	draw.Draw(rgbaImg, bounds, nrgbaImg, bounds.Min, draw.Src)

	return rgbaImg
}
