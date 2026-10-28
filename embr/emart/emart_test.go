package emart

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hackborn/onefunc/embr"
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
func drawImage(t *testing.T, name string, width, height int, layers ...embr.LayerFx) {
	t.Helper()
	img, err := embr.Make(width, height, layers...)
	oferrors.Panic(err)

	saveImage(filepath.Join("testdata", "draw_"+name+".png"), img)
}

func TestDrawCircle1(t *testing.T) {
	w, h := 100, 100
	outside := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	inside := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	drawImage(t, "circle1", w, h, NewCircle(outside, inside))
}

func TestDrawCircle2(t *testing.T) {
	w, h := 100, 100
	outside := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	inside := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	drawImage(t, "circle2", w, h, NewCircle(outside, inside))
}

func TestDrawCircleOutline1(t *testing.T) {
	w, h := 100, 100
	outside := color.NRGBA{R: 0, G: 255, B: 0, A: 0}
	inside := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	outline := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	drawImage(t, "circleoutline1", w, h, NewCircle(outside, inside), NewCircleOutline(outline, 4))
}

func TestDrawCircleOutline2(t *testing.T) {
	w, h := 100, 100
	outside := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	inside := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	outline := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	drawImage(t, "circleoutline2", w, h, NewCircle(outside, inside), NewCircleOutline(outline, 4))
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
