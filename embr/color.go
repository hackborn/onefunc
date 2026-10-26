package embr

import (
	"image/color"
)

// Interp intperolates between color a and b, based
// on an amount of 0 - 1.
func Interp(a, b color.NRGBA, amount float64) color.NRGBA {
	return color.NRGBA{R: Mix(a.R, b.R, amount),
		G: Mix(a.G, b.G, amount),
		B: Mix(a.B, b.B, amount),
		A: Mix(a.A, b.A, amount),
	}
}

func Mix(a, b uint8, aAmount float64) uint8 {
	v := (float64(a) * aAmount) + (float64(b) * (1.0 - aAmount))
	if v <= 0.0 {
		return 0
	} else if v >= 255.0 {
		return 255
	} else {
		return uint8(v)
	}
}
