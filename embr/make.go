package embr

import (
	"image"
	"image/color"
)

type PixelArgs struct {
	X, Y int
}

type PixelFunc func(PixelArgs) color.NRGBA

func Make(width, height int, fn PixelFunc, args ...Arg) (*image.NRGBA, error) {
	setup := Setup{Width: width, Height: height}
	for _, a := range args {
		a.Setup(setup)
	}
	bounds := image.Rect(0, 0, width, height)
	img := image.NewNRGBA(bounds)
	pa := PixelArgs{}
	for pa.Y = range height {
		for pa.X = range width {
			step := Step{X: pa.X, Y: pa.Y}
			for _, a := range args {
				a.Step(step)
			}
			img.SetNRGBA(pa.X, pa.Y, fn(pa))
		}
	}
	return img, nil
}

func MakeRGBA(width, height int, fn PixelFunc, args ...Arg) (*image.RGBA, error) {
	if img, err := Make(width, height, fn, args...); err != nil {
		return nil, err
	} else {
		return convertNRGBAToRGBA(img), nil
	}
}
