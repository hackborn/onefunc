package embr

import (
	"image"
	"image/color"
)

func Make(width, height int, layers ...LayerFx) (*image.NRGBA, error) {
	datafx, err := setup(width, height, layers)
	if err != nil {
		return nil, err
	}
	bounds := image.Rect(0, 0, width, height)
	img := image.NewNRGBA(bounds)
	pa := PixelArgs{}
	for pa.Y = range height {
		for pa.X = range width {
			step := StepArgs{X: pa.X, Y: pa.Y}
			datafx.Step(step)
			pa.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
			for _, a := range layers {
				pa.Color = a.Pixel(pa)
			}
			img.SetNRGBA(pa.X, pa.Y, pa.Color)
		}
	}
	return img, nil
}

func MakeRGBA(width, height int, layers ...LayerFx) (*image.RGBA, error) {
	if img, err := Make(width, height, layers...); err != nil {
		return nil, err
	} else {
		return convertNRGBAToRGBA(img), nil
	}
}

func setup(width, height int, layers []LayerFx) (DataFx, error) {
	setup := LayerSetupArgs{Width: width, Height: height, ptr: &ptrSetupArgs{}}
	for _, a := range layers {
		if err := a.Setup(setup); err != nil {
			return nil, err
		}
	}
	datafx := collapseDataFx(setup.ptr.datafx)
	dataSetup := DataSetupArgs{Width: setup.Width, Height: setup.Height}
	datafx.Setup(dataSetup)
	return datafx, nil
}
