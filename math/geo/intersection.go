package geo

import (
	"math"
)

// FindIntersection finds the intersection between two line segments.
// Note: I am finding cases where this doesn't work, where an
// intersection is found on some segments but if you scale one of
// the points further out, no intersection.
// From https://github.com/vlecomte/cp-geo/blob/master/basics/segment.tex
// Further note: This doesn't seem to find intersections exactly on
// end points. Further review of that page shows yes this is correct,
// there's a different algorithm to account for endpoints.
func FindIntersection[T Number](s1, s2 Segment[T]) (Point[T], bool) {
	oa := Orient(s2.A, s2.B, s1.A)
	ob := Orient(s2.A, s2.B, s1.B)
	oc := Orient(s1.A, s1.B, s2.A)
	od := Orient(s1.A, s1.B, s2.B)
	// Proper intersection exists if opposite signs
	if oa*ob < 0 && oc*od < 0 {
		pt := Point[T]{X: (s1.A.X*ob - s1.B.X*oa) / (ob - oa),
			Y: (s1.A.Y*ob - s1.B.Y*oa) / (ob - oa)}
		return pt, true
	}
	return Point[T]{}, false
}

// Here's the full intersection function, where properInter is FindIntersection above
/*
struct cmpX {
    bool operator()(pt a, pt b) {
        return make_pair(a.x, a.y) < make_pair(b.x, b.y);
    }
};

set<pt,cmpX> inters(pt a, pt b, pt c, pt d) {
    pt out;
    if (properInter(a,b,c,d,out)) return {out};
    set<pt,cmpX> s;
    if (onSegment(c,d,a)) s.insert(a);
    if (onSegment(c,d,b)) s.insert(b);
    if (onSegment(a,b,c)) s.insert(c);
    if (onSegment(a,b,d)) s.insert(d);
    return s;
}
*/

// PerpendicularIntersection finds the intersection of point to the line
// segment by drawing a perpendicular line from point to segment.
// Note: this is from Gemini. Seems to work but not heavily tested.
func PerpendicularIntersection[T Float](seg Segment[T], pt Point[T]) (Point[T], bool) {
	// Check for vertical line (infinite slope)
	if seg.A.X == seg.B.X {
		// Point is on the line if px == x1, return the point itself
		if pt.X == seg.A.X {
			return pt, true
		} else {
			// Point is not on the line, return None (or indicate error)
			return Point[T]{}, false
		}
	}

	// Calculate the slope of the line
	m := (seg.B.Y - seg.A.Y) / (seg.B.X - seg.A.X)

	// Calculate the negative reciprocal of the slope (perpendicular slope)
	m_perp := -1 / m

	// Calculate the x-coordinate of the intersection point (ix)
	ix := (m*pt.X - pt.Y + m_perp*seg.A.X - m_perp*seg.A.Y) / (m - m_perp)

	// Calculate the y-coordinate of the intersection point (iy) using either the point or line equation
	iy := m_perp*(ix-pt.X) + pt.Y // Using point equation

	// Return the intersection point
	return Pt(ix, iy), true
}

func PtToSegIntersection(pt, direction PtF, seg SegF) (PtF, bool) {
	// Create a new segment and check the intersection.
	scale := Max(pt.DistSquared(seg.A), pt.DistSquared(seg.B)) * 10.

	newPt := PtF{X: pt.X + direction.X*scale, Y: pt.Y + direction.Y*scale}
	return FindIntersection(seg, SegF{A: pt, B: newPt})
}

func PtToSegIntersection_Off(pt, direction PtF, seg SegF) (PtF, bool) {
	p := pt
	d := direction
	a := seg.A
	b := seg.B
	//	fmt.Println("dir", direction)
	// Handle parallel cases
	denominator := d.X*(b.Y-a.Y) - d.Y*(b.X-a.X)
	if denominator == 0 {
		return PtF{}, false
	}
	t := ((p.X-a.X)*(b.Y-a.Y) - (p.Y-a.Y)*(b.X-a.X)) / denominator
	u := ((p.X-a.X)*d.Y - (p.Y-a.Y)*d.X) / denominator

	if t >= 0. && 0. <= u && u <= 1 {
		return Pt(p.X+t*d.X, p.Y+t*d.Y), true
	}
	return PtF{}, false
}

func PtToSegIntersection_Off2(pt, direction PtF, seg SegF) (PtF, bool) {
	scale := 1000000.
	newPt := PtF{X: pt.X + direction.X*scale, Y: pt.Y + direction.Y*scale}
	return intersectLines(pt, newPt, seg.A, seg.B, Segments)
}

func intersectLines(a1, a2, b1, b2 PtF, mode Mode) (PtF, bool) {
	// Eh I see this is a shortcut
	if mode == Segments {
		ar := rectFromSeg(a1, a2)
		br := rectFromSeg(b1, b2)
		if !ar.Overlaps(br) {
			return PtF{}, false
		}
	}

	cc := (hp(a1).Cross(hp(a2))).Cross(hp(b1).Cross(hp(b2)))
	if math.Abs(cc.Z) < 1.e-6 {
		return PtF{}, false
	}
	mult := 1. / cc.Z
	x := Pt(cc.X*mult, cc.Y*mult)

	fail := true
	switch mode {
	case Rays:
		fail = !testInter(x, a1, a2, false) || !testInter(x, b1, b2, false)
	case RayLine:
		fail = !testInter(x, a1, a2, false)
	case RaySegment:
		fail = !testInter(x, a1, a2, false) || !testInter(x, b1, b2, true)
	case Segments:
		fail = !testInter(x, a1, a2, true) || !testInter(x, b1, b2, true)
	}
	if fail {
		return PtF{}, false
	}
	/*
		ok := false

		if(mode switch {
		  Mode.Rays => !test(x, a1, a2) || !test(x, b1, b2),
		  Mode.RayLine => !test(x, a1, a2),
		  Mode.RaySegment => !test(x, a1, a2) || !test(x, b1, b2, bidi: true),
		  Mode.Segments => !test(x, a1, a2, bidi: true) || !test(x, b1, b2, bidi: true),
		  _ => false // lines can't fail at this
		}) return false;
	*/
	return x, true
}

type ptfField func(PtF) float64

func ptfFieldX(pt PtF) float64 {
	return pt.X
}
func ptfFieldY(pt PtF) float64 {
	return pt.Y
}

func testInter(p, a, b PtF, bidi bool) bool {
	fn := ptfFieldX
	if math.Abs(b.X-a.X) < math.Abs(b.Y-a.Y) {
		fn = ptfFieldY
	}
	n, d := fn(p)-fn(a), fn(b)-fn(a)
	if bidi && math.Abs(n) > math.Abs(d) {
		return false
	}
	return (n >= 0.) == (d >= 0.)
}

func hp(pt PtF) Pt3dF {
	return Pt3d(pt.X, pt.Y, 1.)
}

func rectFromSeg(a, b PtF) RectF {
	return Rect(Min(a.X, b.X), Min(a.Y, b.Y), Max(a.X, b.X), Max(a.Y, b.Y))
	// var min = Vectx.LeastOf(a, b); // build your own, example is in the post
	// return new Rect(min, Vectx.MostOf(a, b) - min);
}

type Mode = uint8

const (
	Lines Mode = iota
	Rays
	RayLine
	RaySegment
	Segments
)

/*
// https://forum.unity.com/threads/need-a-good-2d-line-segment-intersector-grab-it-while-its-hot.1341557/

// compute intersects like a boss you most certainly are™
static bool intersectLines(Vector2 a1, Vector2 a2, Vector2 b1, Vector2 b2, out Vector2 intersect, Mode mode = Mode.Segments) {
  intersect = new Vector2(float.NaN, float.NaN);

  if(mode == Mode.Segments) {
    var ar = rectFromSeg(a1, a2);
    var br = rectFromSeg(b1, b2);
    if(!ar.Overlaps(br)) return false;
  }

  var cc = (hp(a1).Cross(hp(a2))).Cross(hp(b1).Cross(hp(b2)));
  if(cc.z.Abs() < 1E-6f) return false;

  var x = ((Vector2)cc) * (1f / cc.z);

  if(mode switch {
    Mode.Rays => !test(x, a1, a2) || !test(x, b1, b2),
    Mode.RayLine => !test(x, a1, a2),
    Mode.RaySegment => !test(x, a1, a2) || !test(x, b1, b2, bidi: true),
    Mode.Segments => !test(x, a1, a2, bidi: true) || !test(x, b1, b2, bidi: true),
    _ => false // lines can't fail at this
  }) return false;

  intersect = x;
  return true;

  // local functions
  static Vector3 hp(Vector2 p) => new Vector3(p.x, p.y, 1f);

  static bool test(Vector2 p, Vector2 a, Vector2 b, bool bidi = false) {
    int i = (b.x - a.x).Abs() < (b.y - a.y).Abs()? 1 : 0;
    float n = p[i] - a[i], d = b[i] - a[i];
    if(bidi && n.Abs() > d.Abs()) return false;
    return n >= 0f == d >= 0f;
  }

}

static Rect rectFromSeg(Vector2 a, Vector2 b) {
  var min = Vectx.LeastOf(a, b); // build your own, example is in the post
  return new Rect(min, Vectx.MostOf(a, b) - min);
}

enum Mode {
  Lines,
  Rays,
  RayLine,
  RaySegment,
  Segments
}
Usage example
Code (csharp):
if(Inters
*/

/*
p = np.array(ray_origin)
  d = np.array(ray_direction)
  a = np.array(line_start)
  b = np.array(line_end)

  # Handle parallel cases
  denominator = d[0] * (b[1] - a[1]) - d[1] * (b[0] - a[0])
  if denominator == 0:
      return None, False

  t = ((p[0] - a[0]) * (b[1] - a[1]) - (p[1] - a[1]) * (b[0] - a[0])) / denominator
  u = ((p[0] - a[0]) * d[1] - (p[1] - a[1]) * d[0]) / denominator

  if t >= 0 and 0 <= u <= 1:
      intersection_point = p + t * d
      return intersection_point, True

  return None, False
*/

// https://stackoverflow.com/questions/34415671/intersection-of-a-line-with-a-line-segment-in-c
// line segment p-q intersect with line A-B.
// This looks a lot more efficient than the current seg-seg test I'm doing,
// but also I actually want a seg-seg test because this gives false positives.
// So... maybe I'll switch to it if I think of a way for it to make sense,
// but probably not because this is just a different operation.
func LineToSegIntersection(A, B PtF, pq SegF) (PtF, bool) {
	a := B.Y - A.Y
	b := A.X - B.X
	c := B.X*A.Y - A.X*B.Y
	u := math.Abs(a*pq.A.X + b*pq.A.Y + c)
	v := math.Abs(a*pq.B.X + b*pq.B.Y + c)
	newPt := Pt((pq.A.X*v+pq.B.X*u)/(u+v), (pq.A.Y*v+pq.B.Y*u)/(u+v))
	return newPt, true
}
