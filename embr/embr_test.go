package embr

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	oferrors "github.com/hackborn/onefunc/errors"
)

// go test . -golden
var (
	golden = flag.Bool("golden", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

// isGolden answers true if we are currently updating the golden files.
func isGolden() bool {
	return golden != nil && *golden == true
}

// ---------------------------------------------------------
// TEST-DRAW

// draw is just a place to draw some images.
func drawImage(t *testing.T, name string, width, height int, fn PixelFunc, args ...Arg) {
	t.Helper()
	img, err := Make(width, height, fn, args...)
	oferrors.Panic(err)

	saveImage(filepath.Join("testdata", "draw_"+name+".png"), img)
}

func TestDrawCircle(t *testing.T) {
	w, h := 100, 100
	distance := Distance{}
	fn := func(args PixelArgs) color.NRGBA {
		if distance.Distance < distance.MaxDistance {
			a := uint8(distance.Aa * 255)
			return color.NRGBA{R: 0, G: 0, B: 0, A: a}
		} else {
			return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		}
	}
	drawImage(t, "circle", w, h, fn, &distance)
}

// ---------------------------------------------------------
// IO

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	oferrors.Panic(err)
	defer f.Close()
	err = png.Encode(f, img)
	oferrors.Panic(err)
}
