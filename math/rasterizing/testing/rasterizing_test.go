package testing

import (
	"fmt"
	"math"
	"testing"

	"github.com/hackborn/onefunc/math/geo"
	"github.com/hackborn/onefunc/math/rasterizing"
	"github.com/hackborn/onefunc/math/rasterizing/xiaolinwu"
)

// go test -bench=.

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
		seg  geo.SegI
		want LineOut
	}{
		{geo.Seg(0, 0, 2, 2), newLineOut(0, 0, 1, 1, 2, 2)},
		{geo.Seg(0, 2, 2, 0), newLineOut(0, 2, 1, 1, 2, 0)},
		{geo.Seg(0, 0, 2, 0), newLineOut(0, 0, 1, 0, 2, 0)},
		{geo.Seg(0, 0, -2, 0), newLineOut(0, 0, -1, 0, -2, 0)},
		{geo.Seg(0, 0, -2, -2), newLineOut(0, 0, -1, -1, -2, -2)},
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
// BENCHMARKS

func BenchmarkXiaolinwuShort(b *testing.B) {
	r := xiaolinwu.NewRasterizer()
	runShortLine(b, r)
}

func BenchmarkXiaolinwuLong(b *testing.B) {
	r := xiaolinwu.NewRasterizer()
	runLongLine(b, r)
}

func BenchmarkXiaolinwuVeryLong(b *testing.B) {
	r := xiaolinwu.NewRasterizer()
	runVeryLongLine(b, r)
}

//go:noinline
func runShortLine(b *testing.B, r rasterizing.Rasterizer) {
	seg := geo.Seg(0.0, 0.0, 20.0, 10.0)
	runLine(b, seg, r)
}

//go:noinline
func runLongLine(b *testing.B, r rasterizing.Rasterizer) {
	seg := geo.Seg(0.0, 0.0, 200.0, 100.0)
	runLine(b, seg, r)
}

//go:noinline
func runVeryLongLine(b *testing.B, r rasterizing.Rasterizer) {
	seg := geo.Seg(0.0, 0.0, 800.0, 100.0)
	runLine(b, seg, r)
}

//go:noinline
func runLine(b *testing.B, shape any, r rasterizing.Rasterizer) {
	for n := 0; n < b.N; n++ {
		r.Rasterize(shape, rasterizePixel)
	}
}

func BenchmarkXiaolinwuShort2(b *testing.B) {
	r := xiaolinwu.NewRasterizer2()
	runShortLine2(b, r)
}

func BenchmarkXiaolinwuLong2(b *testing.B) {
	r := xiaolinwu.NewRasterizer2()
	runLongLine2(b, r)
}

func BenchmarkXiaolinwuVeryLong2(b *testing.B) {
	r := xiaolinwu.NewRasterizer2()
	runVeryLongLine2(b, r)
}

//go:noinline
func runShortLine2(b *testing.B, r rasterizing.Rasterizer2) {
	seg := geo.Seg(0.0, 0.0, 20.0, 10.0)
	runLine2(b, seg, r)
}

//go:noinline
func runLongLine2(b *testing.B, r rasterizing.Rasterizer2) {
	seg := geo.Seg(0.0, 0.0, 200.0, 100.0)
	runLine2(b, seg, r)
}

//go:noinline
func runVeryLongLine2(b *testing.B, r rasterizing.Rasterizer2) {
	seg := geo.Seg(0.0, 0.0, 800.0, 100.0)
	runLine2(b, seg, r)
}

//go:noinline
func runLine2(b *testing.B, shape any, r rasterizing.Rasterizer2) {
	for n := 0; n < b.N; n++ {
		r.Rasterize2(shape, rasterizePixels)
	}
}

//go:noinline
func rasterizePixel(rasterizing.Pixel) {
}

//go:noinline
func rasterizePixels([]rasterizing.Pixel) {
}

// ---------------------------------------------------------
// FUNC

type newRasterizerFunc func() rasterizing.Rasterizer

// ---------------------------------------------------------
// SUPPORT

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

func (l *LineOut) Draw(args rasterizing.Pixel) {
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
