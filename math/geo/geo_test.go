package geo

import (
	"fmt"
	"math"
	"testing"

	"github.com/hackborn/onefunc/jacl"
)

// ---------------------------------------------------------
// TEST-XY-TO-INDEX
func TestXYToIndex(t *testing.T) {
	table := []struct {
		width   int
		height  int
		pt      PointI
		want    int
		wantErr error
	}{
		{10, 10, PointI{X: 0, Y: 0}, 0, nil},
		{10, 10, PointI{X: 0, Y: 4}, 40, nil},
		{20, 10, PointI{X: 1, Y: 4}, 81, nil},
		// Errors
		{10, 10, PointI{X: -1, Y: 0}, 0, outOfBoundsErr},
		{10, 10, PointI{X: 10, Y: 10}, 0, outOfBoundsErr},
	}
	for i, v := range table {
		have, haveErr := XYToIndex(v.width, v.height, v.pt)

		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestXYToIndex %v expected err %v but has %v", i, v.wantErr, haveErr)
		} else if v.want != have {
			t.Fatalf("TestXYToIndex %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-LINE
func TestLine(t *testing.T) {
	table := []struct {
		a    PointI
		b    PointI
		want LineOut
	}{
		{PointI{X: 0, Y: 0}, PointI{X: 2, Y: 2}, newLineOut(0, 0, 1, 1, 2, 2)},
		{PointI{X: 0, Y: 2}, PointI{X: 2, Y: 0}, newLineOut(0, 2, 1, 1, 2, 0)},
		{PointI{X: 0, Y: 0}, PointI{X: 2, Y: 0}, newLineOut(0, 0, 1, 0, 2, 0)},
		{PointI{X: 0, Y: 0}, PointI{X: -2, Y: 0}, newLineOut(0, 0, -1, 0, -2, 0)},
		{PointI{X: 0, Y: 0}, PointI{X: -2, Y: -2}, newLineOut(0, 0, -1, -1, -2, -2)},
		//		{PointI{X: 0, Y: 0}, PointI{X: 20, Y: 10}, newLineOut(0, 0, -1, -1, -2, -2)},
	}
	for i, v := range table {
		have := &LineOut{}
		DrawLine(v.a.X, v.a.Y, v.b.X, v.b.Y, have.Draw)
		if err := have.Cmp(v.want); err != nil {
			have.Print()
			t.Fatalf("TestLine %v %v", i, err)
		}
	}
}

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

func (l *LineOut) Draw(x, y int, amount float64) {
	l.pixels = append(l.pixels, Pixel{X: x, Y: y, Amount: amount})
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
	tolerance := 0.000001
	diff := math.Abs(a - b)
	return diff < tolerance
}
