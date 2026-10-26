package embr

import (
	"image/color"
)

// Layer is a single layer of color on the image.
type Layer interface {
	// Setup is used to install any Steppers used
	// by this layer.
	Setup(SetupArgs) error
	// Pixel provides a color for the given pixel.
	Pixel(PixelArgs) color.NRGBA
}

type SetupArgs struct {
	Width, Height int
	ptr           *ptrSetupArgs
}

type PixelArgs struct {
	X, Y  int
	Color color.NRGBA
}

// MustStep answers the step type if it's already in the
// steppers list, otherwise adds it.
func (a *SetupArgs) MustStep[T Stepper](ifMissing T) T {
	if a.ptr == nil {
		return ifMissing
	}
	for _, _s := range a.ptr.steppers {
		if s, ok := _s.(T); ok {
			return s
		}
	}
	a.ptr.steppers = append(a.ptr.steppers, ifMissing)
	return ifMissing
}

// ptrSetupArgs exists solely so the API
// can pass a non-pointer setup. A little
// excessive, but I hate junking APIs with
// things that look like they need nil checks.
type ptrSetupArgs struct {
	steppers []Stepper
}
