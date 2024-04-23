package testing

import (
	"fmt"
	"math"
	"testing"

	"github.com/hackborn/onefunc/math/geo"
	"github.com/hackborn/onefunc/math/rasterizing"
	"github.com/hackborn/onefunc/math/rasterizing/xiaolinwu"
)

// ---------------------------------------------------------
// TEST-LINE-I
func TestLineI(t *testing.T) {
	factories := rasterizerFactories()
	for name, fn := range factories {
		_testFactoryLineI(t, name, fn)
	}
}

func _testFactoryLineI(t *testing.T, facName string, facFunc newRasterizerFunc) {
	table := []struct {
		seg  geo.SegmentI
		want LineOut
	}{
		{segi(0, 0, 2, 2), newLineOut(0, 0, 1, 1, 2, 2)},
		{segi(0, 2, 2, 0), newLineOut(0, 2, 1, 1, 2, 0)},
		{segi(0, 0, 2, 0), newLineOut(0, 0, 1, 0, 2, 0)},
		{segi(0, 0, -2, 0), newLineOut(0, 0, -1, 0, -2, 0)},
		{segi(0, 0, -2, -2), newLineOut(0, 0, -1, -1, -2, -2)},
	}
	for i, v := range table {
		ras := facFunc()
		have := &LineOut{}
		ras.Rasterize(v.seg, have.Draw)
		if err := have.Cmp(v.want); err != nil {
			have.Print()
			t.Fatalf("TestLine %v %v", i, err)
		}
	}
}

func rasterizerFactories() map[string]newRasterizerFunc {
	m := make(map[string]newRasterizerFunc)
	m["xiaolinwu"] = func() rasterizing.Rasterizer {
		return xiaolinwu.NewRasterizer()
	}
	return m
}

// ---------------------------------------------------------
// FUNC

type newRasterizerFunc func() rasterizing.Rasterizer

// ---------------------------------------------------------
// SUPPORT

// pti is a convenience for creating a geo.PointI
func pti(x, y int) geo.PointI {
	return geo.PointI{X: x, Y: y}
}

func segi(x1, y1, x2, y2 int) geo.SegmentI {
	return geo.SegmentI{A: pti(x1, y1), B: pti(x2, y2)}
}

// Create a new line out on the supplied pixel components.
// Components must always be two ints (the x and y) optionally
// followed by a float64 (the amount).
func newLineOut(pixelComponents ...any) LineOut {
	out := LineOut{}
	pos := 0
	for _, _c := range pixelComponents {
		switch c := _c.(type) {
		case int:
			if pos == 0 {
				out.pixels = append(out.pixels, Pixel{X: c, Amount: -1.0})
				pos = 1
			} else if pos == 1 {
				if len(out.pixels) > 0 {
					idx := len(out.pixels) - 1
					p := out.pixels[idx]
					p.Y = c
					out.pixels[idx] = p
				}
				pos = 0
			}
		case float64:
			if len(out.pixels) > 0 {
				idx := len(out.pixels) - 1
				p := out.pixels[idx]
				p.Amount = c
				out.pixels[idx] = p
			}
			pos = 0
		}
	}
	return out
}

type LineOut struct {
	pixels []Pixel
}

func (l *LineOut) Draw(args rasterizing.PixelArgs) {
	l.pixels = append(l.pixels, Pixel{X: args.X, Y: args.Y, Amount: args.Amount})
}

func (l *LineOut) Print() {
	for i, pixel := range l.pixels {
		fmt.Println("", i, ":", pixel)
	}
}

type Pixel struct {
	X, Y   int
	Amount float64
}

// ---------------------------------------------------------
// CMP

func (a *LineOut) Cmp(b LineOut) error {
	if len(a.pixels) != len(b.pixels) {
		return fmt.Errorf("Line size has %v but wants %v", len(a.pixels), len(b.pixels))
	}
	for i, apix := range a.pixels {
		bpix := b.pixels[i]
		if apix.X != bpix.X {
			return fmt.Errorf("pixel %v has x=%v but wants x=%v", i, apix.X, bpix.X)
		}
		if apix.Y != bpix.Y {
			return fmt.Errorf("pixel %v has y=%v but wants y=%v", i, apix.Y, bpix.Y)
		}
		if bpix.Amount >= 0.0 && !floatsEqual(apix.Amount, bpix.Amount) {
			return fmt.Errorf("pixel %v has amount=%v but wants amount=%v", i, apix.Amount, bpix.Amount)
		}
	}
	return nil
}

func floatsEqual(a, b float64) bool {
	const tolerance = 0.0000001
	diff := math.Abs(a - b)
	return diff < tolerance
}
