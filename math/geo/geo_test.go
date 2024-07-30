package geo

import (
	"fmt"
	"math"
	"testing"

	"github.com/hackborn/onefunc/jacl"
)

// go test -bench=.

// ---------------------------------------------------------
// TEST-CIRCLE-HIT-TEST
func TestCircleHitTest(t *testing.T) {
	table := []struct {
		center PtF
		radius float64
		pt     PtF
		want   bool
	}{
		{Pt(5.0, 5.0), 5, Pt(5.0, 4.0), true},
		{Pt(5.0, 5.0), 5, Pt(5.0, 5.0), true},
		{Pt(5.0, 5.0), 5, Pt(5.0, 10.0), true},
		{Pt(5.0, 5.0), 5, Pt(5.0, 10.05), false},
		{Pt(5.0, 5.0), 5, Pt(1.0, 1.0), false},
		{Pt(5.0, 5.0), 5, Pt(2.0, 2.0), true},
	}
	for i, v := range table {
		cht := &CircleHitTest{}
		cht.Set(v.center, v.radius)
		have := cht.Hit(v.pt)
		if have != v.want {
			t.Fatalf("TestCircleHitTest %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-DEGREES
func TestDegrees(t *testing.T) {
	table := []struct {
		center PtF
		pt     PtF
		want   float64
	}{
		{Pt(0.0, 0.0), Pt(10.0, 0.0), 0.0},
		{Pt(0.0, 0.0), Pt(10.0, 10.0), 45.0},
		{Pt(0.0, 0.0), Pt(0.0, 10.0), 90.0},
		{Pt(0.0, 0.0), Pt(-10.0, 10.0), 135.0},
		{Pt(0.0, 0.0), Pt(-10.0, 0.0), 180.0},
		{Pt(0.0, 0.0), Pt(-10.0, -10.0), 225.0},
		{Pt(0.0, 0.0), Pt(0.0, -10.0), 270.0},
		{Pt(0.0, 0.0), Pt(10.0, -10.0), 315.0},
	}
	for i, v := range table {
		have := Seg(v.center.X, v.center.Y, v.pt.X, v.pt.Y).Degrees()
		if !FloatsEqualTol(have, v.want, 0.0001) {
			t.Fatalf("TestDegrees %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-DIST-PT-TO-SEGMENT
func TestDistPointToSegment(t *testing.T) {
	f := func(seg SegF, pt PtF, want float64, wantPt PtF) {
		t.Helper()

		have, havePt := DistPointToSegment(seg, pt)
		if !FloatsEqualTol(have, want, 0.00001) {
			t.Fatalf("TestDistPointToSegment has dist %.6f but wants %.6f", have, want)
		} else if !pointsEqual(havePt, wantPt) {
			t.Fatalf("TestDistPointToSegment has pt %v but wants %v", havePt, wantPt)
		}
	}
	f(Seg(0., 0., 10., 10.), Pt(2., 0.), 1.414214, Pt(1., 1.))
	f(Seg(0., 5., 10., 5.), Pt(2., 0.), 5., Pt(2., 5.))
}

// ---------------------------------------------------------
// TEST-ORIENT
func TestOrient(t *testing.T) {
	f := func(seg SegF, pt PtF, want float64) {
		t.Helper()

		have := Orient(seg.A, seg.B, pt)
		if !FloatsEqualTol(have, want, 0.00001) {
			t.Fatalf("Has %.6f but wants %.6f", have, want)
		}
	}
	f(Seg(5., 5., 10., 5.), Pt(2., 0.), -25.)
	f(Seg(5., 5., 10., 5.), Pt(2., 10.), 25.)
	f(Seg(5., 5., 10., 5.), Pt(2., 5.), 0.)
	// t.Fatalf("no")
}

// ---------------------------------------------------------
// TEST-PERPENDICULAR-INTERSECTION
func TestPerpendicularIntersection(t *testing.T) {
	table := []struct {
		seg    SegF
		pt     PtF
		want   PtF
		wantOk bool
	}{
		{Seg(0.0, 0.0, 20.0, 20.0), Pt(10.0, 0.0), Pt(5.0, 5.0), true},
	}
	for i, v := range table {
		have, haveOk := PerpendicularIntersection(v.seg, v.pt)
		if haveOk != v.wantOk {
			t.Fatalf("TestPerpendicularIntersection %v has %v but expected %v", i, haveOk, v.wantOk)
		} else if !pointsEqual(have, v.want) {
			t.Fatalf("TestPerpendicularIntersection %v has intersection %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-PROJECT-DEGREE
func TestProjectDegree(t *testing.T) {
	table := []struct {
		center   PtF
		degree   float64
		distance float64
		print    bool
	}{
		{Pt(0.0, 0.0), 0, 10, false},
		{Pt(0.0, 0.0), 40, 10, false},
		{Pt(0.0, 0.0), 122, 10, false},
		{Pt(0.0, 0.0), 270, 10, false},
	}
	for i, v := range table {
		pt := v.center.ProjectDegree(v.degree, v.distance)
		have := Seg(v.center.X, v.center.Y, pt.X, pt.Y).Degrees()
		want := v.degree
		if !FloatsEqualTol(have, want, 0.0001) {
			t.Fatalf("TestProjectDegree %v has %v but expected %v for point %v", i, have, want, pt)
		} else if v.print {
			t.Fatalf("TestProjectDegree %v print pt %v degree %v dist %v new point %v", i, v.center, v.degree, v.distance, pt)
		}
	}
}

// ---------------------------------------------------------
// TEST-RANGE-MAP-NORMAL
func TestRangeMapNormal(t *testing.T) {
	table := []struct {
		value float64
		r     RngF
		want  float64
	}{
		{0, Rng(0.0, 1.0), 0},
		{0.1, Rng(0.0, 10.0), 1},
		{0.5, Rng(0.0, 1.0), 0.5},
		{0.5, Rng(0.0, 100.0), 50},
		{1, Rng(0.0, 1.0), 1},
		{10, Rng(0.0, 1.0), 1},
		{-10, Rng(0.0, 1.0), 0},
		{0, Rng(1.0, 0.0), 1},
		{0.1, Rng(10.0, 0.0), 9},
		{1, Rng(100.0, 10.0), 10},
		{100, Rng(100.0, 10.0), 10},
		{-1, Rng(1.0, 0.0), 1},
	}
	for i, v := range table {
		have := v.r.MapNormal(v.value)
		if !FloatsEqualTol(have, v.want, 0.0001) {
			t.Fatalf("TestRangeMapNormal %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-RATIO
func TestRatio(t *testing.T) {
	table := []struct {
		a    float64
		b    float64
		want float64
	}{
		{0.0, 1.0, 0.0},
		{1.0, 0.0, 1.0},
		{2, 2, 0.5},
		{25, 100, 0.2},
		{75, 100, 0.42857},
		{125, 100, 0.5555},
		{200, 100, 0.6666},
		{500, 100, 0.8333},
	}
	for i, v := range table {
		have := Ratio(v.a, v.b)

		if !FloatsEqualTol(have, v.want, 0.0001) {
			t.Fatalf("TestRatio %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-NORMALIZE
func TestNormalize(t *testing.T) {
	table := []struct {
		value float64
		r     RngF
		want  float64
	}{
		{-10, Rng(0.0, 10.0), 0},
		{0, Rng(0.0, 10.0), 0},
		{5, Rng(0.0, 10.0), 0.5},
		{10, Rng(0.0, 10.0), 1},
		{20, Rng(0.0, 10.0), 1},
		{84000, Rng(10000.0, 90000.0), 0.925},
		{0, Rng(1.0, .0), 1},
		{.25, Rng(1.0, .0), .75},
	}
	for i, v := range table {
		have := v.r.Normalize(v.value)
		if !FloatsEqualTol(have, v.want, 0.0001) {
			t.Fatalf("TestNormalize %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-POINT-NORMALIZE
func TestPointNormalize(t *testing.T) {
	f := func(pt PtF, want PtF) {
		t.Helper()

		have := pt.Normalize()
		if !pointsEqual(want, have) {
			t.Fatalf("TestPointNormalize wants %v has %v", want, have)
		}
	}
	f(Pt(10., 0.), Pt(1., 0.))
	f(Pt(0., 10.), Pt(0., 1.))
	f(Pt(10., 10.), Pt(.7071067, .7071067))
}

// ---------------------------------------------------------
// TEST-POINT-ROTATE
func TestPointRotate(t *testing.T) {
	f := func(pt PtF, angleInRads float64, want PtF) {
		t.Helper()

		cen := Pt(0., 0.)
		have := pt.Rotate(cen, angleInRads)
		if !pointsEqual(want, have) {
			t.Fatalf("TestPointRotate wants %v has %v", want, have)
		}
	}
	f(Pt(10., 0.), DegreesToRadians(45), Pt(7.07106781, 7.07106781))
	f(Pt(10., 0.), DegreesToRadians(90), Pt(0., 10.))
	f(Pt(10., 0.), DegreesToRadians(180), Pt(-10., 0.))
}

// ---------------------------------------------------------
// TEST-SEGMENT-INTERSECTION
func TestSegmentIntersection(t *testing.T) {
	table := []struct {
		s1     SegF
		s2     SegF
		want   PtF
		wantOk bool
	}{
		{Seg(5., 0, 5, 10), Seg(0., 5, 10, 5), Pt(5., 5), true},
		{Seg(0., 0, 10, 10), Seg(0., 10, 10, 0), Pt(5., 5), true},
		{Seg(0., 0, 10, 10), Seg(0., 5, 20, 0), Pt(4., 4), true},
	}
	for i, v := range table {
		have, haveOk := FindIntersection(v.s1, v.s2)

		if v.wantOk != haveOk {
			t.Fatalf("TestSegmentIntersection %v has ok %v but expected %v", i, haveOk, v.wantOk)
		} else if haveOk && !pointsEqual(v.want, have) {
			t.Fatalf("TestSegmentIntersection %v has intersection %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-POINT-SEGMENT-INTERSECTION
func TestPointSegmentIntersection(t *testing.T) {
	table := []struct {
		s     SegF
		p     PtF
		check bool // Set to true to fail the test and see the values
	}{
		{Seg(0., 0, 10, 10), Pt(10., 0), false},
		{Seg(0., 0, 10, 10), Pt(20., 10), false},
		{Seg(0., 0, 10, 10), Pt(5., 2), false},
	}
	for i, v := range table {
		have, pt := DistSquared(v.s, v.p)

		if v.check {
			t.Fatalf("TestPointSegmentIntersection %v segment (%v) to pt %v found %v dist %v distsqrt %v", i, v.s, v.p, pt, have, math.Sqrt(have))
			//		} else if !pointsEqual(v.want, have) {
			//	t.Fatalf("TestSegmentIntersection %v has intersection %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-PROJECT
func TestProject(t *testing.T) {
	f := func(pt PtF, m, dist float64, wantP, wantN PtF) {
		t.Helper()

		haveP, haveN := pt.Project(m, dist)
		if !pointsEqual(wantP, haveP) {
			t.Fatalf("TestProject wants positive %v has %v", wantP, haveP)
		} else if !pointsEqual(wantN, haveN) {
			t.Fatalf("TestProject wants negative %v has %v", wantN, haveN)
		}
		/*
			cmpP, cmpN := pt.projectAlgo2(m, dist)
			if !pointsEqual(cmpP, haveP) {
				t.Fatalf("TestProject wants cmp positive %v has %v", cmpP, haveP)
			} else if !pointsEqual(cmpN, haveN) {
				t.Fatalf("TestProject wants cmp negative %v has %v", cmpN, haveN)
			}
		*/
	}
	//	/*
	f(Pt(10., 10.), 1., 5., Pt(13.535533, 6.46446609), Pt(6.464466, 13.5355339))
	f(Pt(10., 10.), 0., 5., Pt(15., 10.), Pt(5., 10.))
	f(Pt(10., 10.), 0.0000000001, 5., Pt(15., 10.), Pt(5., 10.))
	f(Pt(10., 10.), 0.1, 5., Pt(14.9751859, 9.50248140), Pt(5.02481404, 10.4975185))
	f(Pt(10., 10.), math.MaxFloat64, 5., Pt(10., 15.), Pt(10., 5.))
	f(Pt(10., 10.), 9999999999., 5., Pt(10., 5.), Pt(9.9999999, 15.))
	//	*/
	//
	// Real world data
	f(Pt(8., 8.), 2., 5., Pt(10.2360679, 3.5278640), Pt(5.763932022, 12.47213595))
	f(Pt(8., 8.), -2., 5., Pt(5.76393202, 3.52786404), Pt(10.236067977, 12.472135954))
}

// ---------------------------------------------------------
// TEST-SLOPE-AND-PROJECT
func TestSlopeAndProject(t *testing.T) {
	table := []struct {
		seg   SegF
		dist  float64
		check bool // Set to true to fail the test and see the values
	}{
		{Seg(0., 0, 10, 10), 1, false},
		{Seg(0., 0, 20, 10), 1, false},
		{Seg(0., 0, 20, 10), 10, false},
		{Seg(0., 0, 1, 100), 1, false},
		{Seg(0., 0, 0, 10), 1, false},
		{Seg(0., 0, 10, 0), 1, false},
	}
	for i, v := range table {
		//		m := v.seg.Slope()
		m := v.seg.PerpendicularSlope().M
		haveA, haveB := v.seg.A.Project(m, v.dist)
		distA := v.seg.A.Dist(haveA)
		distB := v.seg.A.Dist(haveB)
		if !FloatsEqual(v.dist, distA) {
			t.Fatalf("TestSlopeAndProject %v has distance a %v but expected %v (pt %v)", i, distA, v.dist, haveA)
		} else if !FloatsEqual(v.dist, distB) {
			t.Fatalf("TestSlopeAndProject %v has distance b %v but expected %v (pt %v)", i, distB, v.dist, haveB)
		} else if v.check {
			t.Fatalf("TestSlopeAndProject %v check: m %v dist %v pt %v haveA %v haveB %v distA %v distB %v", i, m, v.dist, v.seg.A, haveA, haveB, distA, distB)
		}
	}
}

// ---------------------------------------------------------
// TEST-SEGMENT-INTERP
func TestSegmentInterp(t *testing.T) {
	f := func(seg SegF, unit float64, want PtF) {
		t.Helper()

		have := seg.Interp(unit)
		if !pointsEqual(want, have) {
			t.Fatalf("TestSegmentInterp wants %v has %v", want, have)
		}
	}
	f(Seg(10., 10., 20., 20.), -1.5, Pt(-5., -5.))
	f(Seg(10., 10., 20., 20.), 0., Pt(10., 10.))
	f(Seg(10., 10., 20., 20.), .5, Pt(15., 15.))
	f(Seg(10., 10., 20., 20.), 1., Pt(20., 20.))
	f(Seg(10., 10., 20., 20.), 1.5, Pt(25., 25.))
}

// ---------------------------------------------------------
// TEST-PT-TO-SEG-INTERSECTION
func TestPtToSegIntersection(t *testing.T) {
	table := []struct {
		pt      PtF
		degrees float64
		seg     SegF
		want    PtF
		wantOk  bool
	}{
		/*
			{Pt(10., 10.), 270., Seg(0., 0., 20., 0.), Pt(10., 0.), true},
			{Pt(10., 10.), 90., Seg(0., 20., 20., 20.), Pt(10., 20.), true},
			{Pt(10., 10.), 0., Seg(0., 0., 0., 20.), Pt(0., 0.), false},
			{Pt(10., 10.), 0., Seg(20., 0., 20., 20.), Pt(20., 10.), true},

			//	{Pt(2.5, 1.5), 135., Seg(0., 1.4, 1.4, 0.), Pt(0., 0.), false},
		*/
		//		{Pt(2.5, 1.5), 135., Seg(0., 4., 0., 1.4), Pt(0., 0.), true},
		// TODO: This doesn't work but should
		//		{Pt(2.5, 1.5), 135., Seg(4., 4., 0., 4.), Pt(0., 0.), true},

		// [{{0 1.4} {1.4 0}} {{1.4 0} {4 0}} {{4 0} {4 4}} {{4 4} {0 4}} {{0 4} {0 1.4}}]
	}
	for i, v := range table {
		dir := DegreesToDirectionCw(v.degrees)
		fmt.Println("dir", dir, "degrees", v.degrees)
		scale := 3.536
		newPt := PtF{X: v.pt.X + dir.X*scale, Y: v.pt.Y + dir.Y*scale}
		fmt.Println("pt", v.pt, "newpt", newPt, "seg", v.seg)
		/*
			segs := []SegF{}
			segs = append(segs, Seg(0, 0, 10, 0.))
			segs = append(segs, Seg(10, 0, 10, 10.))
			segs = append(segs, Seg(10, 10, 0, 10.))
			segs = append(segs, Seg(0, 10, 0, 0.))
			deg := v.degrees
			dir2 := DegreesToDirectionCw(deg)
			fmt.Println("deg", deg, "dir", dir2)
			for _, seg := range segs {
				have, haveOk := PtToSegIntersection(v.pt, dir2, seg)
				if haveOk {
					fmt.Println("have", have, "seg", seg)
				} else {
					fmt.Println("no have seg", seg)
				}
			}

			cmpseg := Seg(-10., 5., 20., 5.)
			fmt.Println("TEST", cmpseg)
			for _, seg := range segs {
				have, haveOk := intersectLines(cmpseg.A, cmpseg.B, seg.A, seg.B, Segments)
				if haveOk {
					fmt.Println("have", have, "seg", seg)
				} else {
					fmt.Println("no have seg", seg)
				}
			}
		*/
		have, haveOk := PtToSegIntersection(v.pt, dir, v.seg)

		/*
			dirpt := v.pt.Add(dir)
			have2, have2Ok := LineToSegIntersection(v.pt, dirpt, v.seg)
			if !pointsEqual(have, have2) {
				t.Fatalf("TestPtToSegIntersection %v NEW CHECK BAD has %v new one %v newHave %v", i, have, have2, have2Ok)
			}
		*/

		if haveOk != v.wantOk {
			t.Fatalf("TestPtToSegIntersection %v expected ok %v but has %v", i, v.wantOk, haveOk)
		} else if haveOk && !pointsEqual(v.want, have) {
			t.Fatalf("TestPtToSegIntersection %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-PT-TO-SEG-INTERSECTIONS
// Harness for some real-world intersection tests.
func __TestPtToSegIntersections(t *testing.T) {
	f := func(pt PtF, degrees float64, segs []SegF, want []PtF) {
		t.Helper()

		dir := DegreesToDirectionCw(degrees)
		have := []PtF{}
		for _, seg := range segs {
			if havePt, ok := PtToSegIntersection(pt, dir, seg); ok {
				have = append(have, havePt)
			}
		}

		if !pointSlicesEqual(want, have) {
			t.Fatalf("TestPtToSegIntersections wants %v has %v", want, have)
		}
	}
	f(Pt(2.5, 1.5), 135., segs(Pt(0, 1.4), Pt(1.4, 0), Pt(1.4, 0), Pt(4., 0), Pt(4., 0), Pt(4., 4), Pt(4., 4), Pt(0., 4), Pt(0., 4), Pt(0., 1.4)), pts(Pt(1., 0.)))
}

// ---------------------------------------------------------
// TEST-TRIANGLE-PLANE-INTERSECTION
func TestTrianglePlaneIntersection(t *testing.T) {
	table := []struct {
		tri      Tri3dF
		triPt    Pt3dF
		planeTri Tri3dF
		planePt  Pt3dF
		want     Pt3dF
		wantOk   bool
	}{
		{Tri3dFlat(1., 0., 2., 2., 1., 2., 1., 2., 4), tri1.Center(), plane1Tri, zeroPt, Pt3d(-1., 2., 0.), true},
	}
	for i, v := range table {
		//		have, haveOk := TrianglePlaneIntersection(v.tri, v.triPt, v.planeTri, v.planePt)
		have, haveOk := TrianglePlaneIntersection(v.tri, v.triPt, v.planeTri, v.planePt)

		if haveOk != v.wantOk {
			t.Fatalf("TestTrianglePlaneIntersection %v expected ok %v but has %v", i, v.wantOk, haveOk)
		} else if !points3dEqual(v.want, have) {
			t.Fatalf("TestTrianglePlaneIntersection %v has %v but expected %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-RADIAL-DIST
func TestRadialDist(t *testing.T) {
	f := func(a, b float64, want float64, wantOrientation Orientation) {
		t.Helper()

		have, haveOrientation := RadialDist(a, b)
		if haveOrientation != wantOrientation {
			t.Fatalf("Has orientation %v but wants %v", haveOrientation, wantOrientation)
		} else if !FloatsEqualTol(have, want, 0.00001) {
			t.Fatalf("Has %.6f but wants %.6f", have, want)
		}
	}
	f(.5, .5, 1., Collinear)
	f(.3, .4, .8, CounterClockwise)
	f(.3, .5, .6, CounterClockwise)
	f(.4, .3, .8, Clockwise)
	f(.1, .9, .6, Clockwise)
	f(.9, .1, .6, CounterClockwise)
}

var zeroPlane = Tri3dFlat(10., 0., 0., 10., 10., 0., 0., 10., 0.)
var zeroPt = Pt3d(0., 0., 0.)
var tri1 = Tri3dFlat(1., 0., 4., 2., 1., 4., 1., 2., 2)
var plane1Tri = Tri3dFlat(1., 0., 0., 2., 1., 0., 1., 2., 0.)

var tri2 = Tri3dFlat(1., 0., 1., 2., 1., 0., 1., 2., 1.)
var tri2Pt = Pt3d(1.5, 0.5, 0)
var plane2Nornal = Pt3d(0., 1., 1.)
var plane2Tri = Tri3dFlat(1., 0., 0., 2., 1., 0., 1., 2., 0.)

// var plane2Normal = Tri3dFlat(1., 0., 0., 2., 1., 0., 1., 2., 0.)
// ---------------------------------------------------------
// TEST-XY-TO-INDEX
func TestXYToIndex(t *testing.T) {
	table := []struct {
		width   int
		height  int
		pt      PtI
		want    int
		wantErr error
	}{
		{10, 10, PtI{X: 0, Y: 0}, 0, nil},
		{10, 10, PtI{X: 0, Y: 4}, 40, nil},
		{20, 10, PtI{X: 1, Y: 4}, 81, nil},
		// Errors
		{10, 10, PtI{X: -1, Y: 0}, 0, outOfBoundsErr},
		{10, 10, PtI{X: 10, Y: 10}, 0, outOfBoundsErr},
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
		seg    SegF
		y      float64
		want   float64
		wantOk bool
		check  bool // Set to true to fail the test and see the values
	}{
		{Seg(0., 0, 0, 10), 5, 0, true, false},
		{Seg(5., 0, 5, 10), 5, 5, true, false},
		{Seg(0., 0, 10, 10), 5, 5, true, false},
		{Seg(0., 0, 10, 20), 5, 2.5, true, false},
		{Seg(5., 0, 5, 10), 11, 0, false, false},
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
// BENCHMARKS

func BenchmarkDistPointToSegment(b *testing.B) {
	seg := Seg(10., 10., 100., 110.)
	pt := Pt(20., 50.)

	for n := 0; n < b.N; n++ {
		DistPointToSegment(seg, pt)
	}
}

func BenchmarkDistSquared(b *testing.B) {
	seg := Seg(10., 10., 100., 110.)
	pt := Pt(20., 50.)

	for n := 0; n < b.N; n++ {
		DistSquared(seg, pt)
	}
}

func BenchmarkOrient(b *testing.B) {
	seg := Seg(10., 10., 100., 110.)
	pt := Pt(20., 50.)

	for n := 0; n < b.N; n++ {
		Orient(seg.A, seg.B, pt)
	}
}

// ---------------------------------------------------------
// CONSTRUCTION

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

func segs(ab ...PtF) []SegF {
	ans := []SegF{}
	for i := 1; i < len(ab); i += 2 {
		a, b := ab[i-1], ab[i]
		ans = append(ans, Seg(a.X, a.Y, b.X, b.Y))
	}
	return ans
}

func pts(pts ...PtF) []PtF {
	return pts
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

func pointsEqual(a, b PtF) bool {
	return FloatsEqualTol(a.X, b.X, 0.000001) && FloatsEqualTol(a.Y, b.Y, 0.000001)
}

func points3dEqual(a, b Pt3dF) bool {
	return FloatsEqual(a.X, b.X) && FloatsEqual(a.Y, b.Y) && FloatsEqual(a.Z, b.Z)
}

func pointSlicesEqual(a, b []PtF) bool {
	if len(a) != len(b) {
		return false
	}
	for i, apt := range a {
		if !pointsEqual(apt, b[i]) {
			return false
		}
	}
	return true
}
