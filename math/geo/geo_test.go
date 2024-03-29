package geo

import (
	"fmt"
	"testing"

	"github.com/hackborn/onefunc/jacl"
)

// ---------------------------------------------------------
// TEST-SEGMENT-INTERSECTION
func TestSegmentIntersection(t *testing.T) {
	table := []struct {
		s1     SegmentF64
		s2     SegmentF64
		want   PointF64
		wantOk bool
		check  bool // Set to true to fail the test and see the values
	}{
		{segf(5, 0, 5, 10), segf(0, 5, 10, 5), ptf(5, 5), true, false},
		{segf(0, 0, 10, 10), segf(0, 10, 10, 0), ptf(5, 5), true, false},
		{segf(0, 0, 10, 10), segf(0, 5, 20, 0), ptf(4, 4), true, false},
	}
	for i, v := range table {
		have, haveOk := FindIntersection(v.s1, v.s2)

		if v.check {
			have2, haveOk2 := FindIntersectionBAD(v.s1, v.s2)
			t.Fatalf("TestSegmentIntersection %v segment (%v) to (%v) pt1 %v pt2 %v ok1 %v ok2 %v", i, v.s1, v.s2, have, have2, haveOk, haveOk2)
		} else if v.wantOk != haveOk {
			t.Fatalf("TestSegmentIntersection %v has ok %v but expected %v", i, haveOk, v.wantOk)
		} else if !pointsEqual(v.want, have) {
			t.Fatalf("TestSegmentIntersection %v has intersection %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-PROJECT
func TestProject(t *testing.T) {
	table := []struct {
		seg   SegmentF64
		dist  float64
		check bool // Set to true to fail the test and see the values
	}{
		{segf(0, 0, 10, 10), 1, false},
		{segf(0, 0, 20, 10), 1, false},
		{segf(0, 0, 20, 10), 10, false},
		{segf(0, 0, 1, 100), 1, false},
		{segf(0, 0, 0, 10), 1, false},
		{segf(0, 0, 10, 0), 1, false},
	}
	for i, v := range table {
		//		m := v.seg.Slope()
		m := v.seg.PerpendicularSlope()
		haveA, haveB := v.seg.A.Project(m, v.dist)
		distA := v.seg.A.Dist(haveA)
		distB := v.seg.A.Dist(haveB)
		if !FloatsEqual(v.dist, distA) {
			t.Fatalf("TestProject %v has distance a %v but expected %v (pt %v)", i, distA, v.dist, haveA)
		} else if !FloatsEqual(v.dist, distB) {
			t.Fatalf("TestProject %v has distance b %v but expected %v (pt %v)", i, distB, v.dist, haveB)
		} else if v.check {
			t.Fatalf("TestProject %v check: m %v dist %v pt %v haveA %v haveB %v distA %v distB %v", i, m, v.dist, v.seg.A, haveA, haveB, distA, distB)
		}
	}
}

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
// TEST-X-AT-Y
func TestXAtY(t *testing.T) {
	table := []struct {
		seg    SegmentF64
		y      float64
		want   float64
		wantOk bool
		check  bool // Set to true to fail the test and see the values
	}{
		{SegmentF64{A: ptf(0, 0), B: ptf(0, 10)}, 5, 0, true, false},
		{SegmentF64{A: ptf(5, 0), B: ptf(5, 10)}, 5, 5, true, false},
		{SegmentF64{A: ptf(0, 0), B: ptf(10, 10)}, 5, 5, true, false},
		{SegmentF64{A: ptf(0, 0), B: ptf(10, 20)}, 5, 2.5, true, false},
		{SegmentF64{A: ptf(5, 0), B: ptf(5, 10)}, 11, 0, false, false},
	}
	for i, v := range table {
		have, haveOk := XAtY(v.seg, v.y)

		if !FloatsEqual(v.want, have) {
			t.Fatalf("TestXAtY %v has x %v but expected %v", i, have, v.want)
		} else if v.wantOk != haveOk {
			t.Fatalf("TestXAtY %v has ok %v but expected %v", i, haveOk, v.wantOk)
		} else if v.check {
			t.Fatalf("TestXAtY %v check: seg %v y %v wants %v has %v", i, v.seg, v.y, v.want, have)
		}
	}
}

// ---------------------------------------------------------
// SUPPORT

// ptf is a convenience for creating a PointF64
func ptf(x, y float64) PointF64 {
	return PointF64{X: x, Y: y}
}

func segf(x1, y1, x2, y2 float64) SegmentF64 {
	return SegmentF64{A: ptf(x1, y1), B: ptf(x2, y2)}
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
		if bpix.Amount >= 0.0 && !FloatsEqual(apix.Amount, bpix.Amount) {
			return fmt.Errorf("pixel %v has amount=%v but wants amount=%v", i, apix.Amount, bpix.Amount)
		}
	}
	return nil
}

func pointsEqual(a, b PointF64) bool {
	return FloatsEqual(a.X, b.X) && FloatsEqual(a.Y, b.Y)
}
