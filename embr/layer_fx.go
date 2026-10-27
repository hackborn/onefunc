package embr

import (
	"image/color"
)

// LayerFx is a single layer of color on the image.
// It is run once per pixel.
type LayerFx interface {
	// Setup is used to install any DataFx used
	// by this layer.
	Setup(LayerSetupArgs) error
	// Pixel provides a color for the given pixel.
	Pixel(PixelArgs) color.NRGBA
}

type LayerSetupArgs struct {
	Width, Height int
	ptr           *ptrSetupArgs
}

type PixelArgs struct {
	X, Y  int
	Color color.NRGBA
}

// MustData answers the DataFx type if it's already in the
// list, otherwise adds ifMissing.
func (a *LayerSetupArgs) MustData[T DataFx](ifMissing T) T {
	if a.ptr == nil {
		return ifMissing
	}
	for _, _s := range a.ptr.datafx {
		if s, ok := _s.(T); ok {
			return s
		}
	}
	a.ptr.datafx = append(a.ptr.datafx, ifMissing)
	return ifMissing
}

// ptrSetupArgs exists solely so the API
// can pass a non-pointer setup. A little
// excessive, but I hate junking APIs with
// things that look like they need nil checks.
type ptrSetupArgs struct {
	datafx []DataFx
}
