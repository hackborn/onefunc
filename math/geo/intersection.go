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

// LineLineIntersectionF provides the intersection of two lines.
// https://github.com/vlecomte/cp-geo/blob/master/basics/line.tex
func LineLineIntersectionF(l1, l2 LnF) (PtF, bool) {
	d := l1.V.Cross(l2.V)
	if d == 0 {
		return PtF{}, false
	}
	x, y := (l2.V.X*l1.C-l1.V.X*l2.C)/d, (l2.V.Y*l1.C-l1.V.Y*l2.C)/d
	return PtF{X: x, Y: y}, true
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

// CubicBezierDistance approximates the distance from a point to a Bezier curve
func CubicBezierDistance(cb CubicBezier, pt PtF) (float64, PtF) {
	d, p := CubicBezierDistanceSquared(cb, pt)
	return math.Sqrt(d), p
}

// CubicBezierDistanceSquared approximates the distance from a point to a Bezier curve
func CubicBezierDistanceSquared(cb CubicBezier, pt PtF) (float64, PtF) {
	// simple binary search-like algo
	// TODO: Just a random experiment. Not optimized or necessarily even right.
	stepRng := Rng(0., 1.)
	count := 0
	change := 9999999.
	ptmid := PtF{}
	distmid := 0.
	for count < 100 && change > .001 {
		// Get the endpoints for the range and find which one is closer
		// to the current midpoint.
		midpoint := stepRng.Midpoint()
		ptmin := cb.At(stepRng.Min)
		distmin := pt.DistSquared(ptmin)
		ptmax := cb.At(stepRng.Max)
		distmax := pt.DistSquared(ptmax)
		ptmid = cb.At(midpoint)
		distmid = pt.DistSquared(ptmid)

		// Could obviously remove a bunch of calculations here by reusing endpoints.
		// It's really only the midpoint that needs to keep getting calculated.
		left, right := math.Abs(distmid-distmin), math.Abs(distmid-distmax)
		if left <= right {
			stepRng.Min, stepRng.Max = stepRng.Min, midpoint
			change = left
		} else {
			stepRng.Min, stepRng.Max = midpoint, stepRng.Max
			change = right
		}

		count++
	}
	//	fmt.Println("count", count, "change", change)
	return distmid, ptmid
}

/*
float dot2( vec2 v ) { return dot(v,v); }
float cro( vec2 a, vec2 b ) { return a.x*b.y-a.y*b.x; }
float cos_acos_3( float x ) { x=sqrt(0.5+0.5*x); return x*(x*(x*(x*-0.008972+0.039071)-0.107074)+0.576975)+0.5; } // https://www.shadertoy.com/view/WltSD7
*/
//func dot2[T Number](pt Point[T]) float64 {
//	return pt.Dot(pt)
//}

// func cos_acos_3( float x ) { x=sqrt(0.5+0.5*x); return x*(x*(x*(x*-0.008972+0.039071)-0.107074)+0.576975)+0.5; } // https://www.shadertoy.com/view/WltSD7
// https://www.shadertoy.com/view/WltSD7
func cos_acos_3(x float64) float64 {
	x = math.Sqrt(.5 + .5*x)
	return x*(x*(x*(x*-0.008972+0.039071)-0.107074)+0.576975) + 0.5
}

func shadersign(v float64) float64 {
	if v < 0. {
		return -1
	} else if v > 0. {
		return 1.
	} else {
		return 0.
	}
}

// QuadraticBezierDistance approximates the distance from a point to a Bezier curve
// TODO: Playing with this but I must have messed up some translation because
// it's not currently right, but I might end up not using it anyway so it's
// on hold for now.
func quadraticBezierDistance(bez QuadraticBezier, pos PtF) (float64, PtF) {
	// signed distance to a quadratic bezier
	//float sdBezier( in vec2 pos, in vec2 A, in vec2 B, in vec2 C, out vec2 outQ )

	ax, ay := bez.P1.X-bez.P0.X, bez.P1.Y-bez.P0.Y
	//	vec2 a = B - A;
	bx, by := bez.P0.X-2.*bez.P1.X+bez.P2.X, bez.P0.Y-2.*bez.P1.Y+bez.P2.Y
	//    vec2 b = A - 2.0*B + C;
	cx, cy := ax*2., ay*2.
	//    vec2 c = a * 2.0;
	dx, dy := bez.P0.X-pos.X, bez.P0.Y-pos.Y
	//   vec2 d = A - pos;
	//	fmt.Println("a", ax, ay, "b", bx, by, "c", cx, cy, "d", dx, dy)

	a, b := PtF{X: ax, Y: ay}, PtF{X: bx, Y: by}
	d := PtF{X: dx, Y: dy}
	// cubic to be solved (kx*=3 and ky*=3)
	kk := 1. / b.Dot(b)
	//	fmt.Println("b", b, "kk dot", b.Dot(b))
	// float kk = 1.0/dot(b,b);
	kx := kk * a.Dot(b)
	// float kx = kk * dot(a,b);
	ky := kk * (2.*a.Dot(a) + d.Dot(b)) / 3.
	// float ky = kk * (2.0*dot(a,a)+dot(d,b))/3.0;
	kz := kk * d.Dot(a)
	//    float kz = kk * dot(d,a);
	res := 0.
	// float res = 0.0;
	sgn := 0.
	//    float sgn = 0.0;

	p := ky - kx*kx
	// float p  = ky - kx*kx;
	q := kx*(2.0*kx*kx-3.0*ky) + kz
	// float q  = kx*(2.0*kx*kx - 3.0*ky) + kz;
	p3 := p * p * p
	// float p3 = p*p*p;
	q2 := q * q
	// float q2 = q*q;
	h := q2 + 4.0*p3
	//    float h  = q2 + 4.0*p3;

	outPt := PtF{}
	if h >= 0.0 {
		// 1 root
		h = math.Sqrt(h)
		xx, xy := (h-q)/2., (-h-q)/2.
		//        vec2 x = (vec2(h,-h)-q)/2.0;
		uv := PtF{X: shadersign(xx) * math.Pow(math.Abs(xx), 1./3.), Y: shadersign(xy) * math.Pow(math.Abs(xy), 1./3.)}
		// vec2 uv = sign(x)*pow(abs(x), vec2(1.0/3.0));
		t := uv.X + uv.Y
		//        t := uv.X + uv.Y

		// from NinjaKoala - single newton iteration to account for cancellation
		t -= (t*(t*t+3.0*p) + q) / (3.0*t*t + 3.0*p)

		t = RngFUnit.ClampFast(t - kx)
		w := PtF{X: dx + (cx+bx*t)*t, Y: dy + (cy+by*t)*t}
		//        vec2  w = d+(c+b*t)*t;
		outPt = w.Add(pos)
		res = w.Dot(w)
		sgn = PtF{X: cx + 2.*bx*t, Y: cy + 2.*by*t}.Cross(w)
		//    	sgn = cro(c+2.0*b*t,w);
	} else { // 3 roots
		z := math.Sqrt(-p)
		m := cos_acos_3(q / (p * z * 2.0))
		//        float m = cos_acos_3( q/(p*z*2.0) );
		n := math.Sqrt(1.0 - m*m)

		n *= math.Sqrt(3.0)
		t := Pt3d((m+m)*z-kx, (-n-m)*z-kx, (n-m)*z-kx).CLampFast(RngFUnit)
		//        vec3  t = clamp( vec3(m+m,-n-m,n-m)*z-kx, 0.0, 1.0 );
		qx := Pt(dx+(cx+bx*t.X)*t.X, dy+(cy+by*t.X)*t.X)
		dx := qx.Dot(qx)
		sx := Pt(ax+bx*t.X, ay+by*t.X).Cross(qx)
		//        vec2  qx=d+(c+b*t.x)*t.x; float dx=dot2(qx), sx=cro(a+b*t.x,qx);
		qy := Pt(dx+(cx+bx*t.Y)*t.Y, dy+(cy+by*t.Y)*t.Y)
		dy := qy.Dot(qy)
		sy := Pt(ax+bx*t.Y, ay+by*t.Y).Cross(qy)
		//        vec2  qy=d+(c+b*t.y)*t.y; float dy=dot2(qy), sy=cro(a+b*t.y,qy);
		if dx < dy {
			res = dx
			sgn = sx
			outPt = qx.Add(pos)
		} else {
			res = dy
			sgn = sy
			outPt = qy.Add(pos)
		}
		//        if( dx<dy ) {res=dx;sgn=sx;outQ=qx+pos;} else {res=dy;sgn=sy;outQ=qy+pos;}
	}

	return math.Sqrt(res) * shadersign(sgn), outPt
	// return sqrt(res) * sign(sgn), outPt
}

// QuadraticBezierDistanceSquared approximates the distance from a point to a Bezier curve
func QuadraticBezierDistanceSquared(cb QuadraticBezier, pt PtF) (float64, PtF) {
	// simple binary search-like algo
	stepRng := Rng(0., 1.)
	count := 0
	change := 9999999.
	ptmid := PtF{}
	distmid := 0.
	for count < 100 && change > .001 {
		// Get the endpoints for the range and find which one is closer
		// to the current midpoint.
		midpoint := stepRng.Midpoint()
		ptmin := cb.At(stepRng.Min)
		distmin := pt.DistSquared(ptmin)
		ptmax := cb.At(stepRng.Max)
		distmax := pt.DistSquared(ptmax)
		ptmid = cb.At(midpoint)
		distmid = pt.DistSquared(ptmid)

		// Could obviously remove a bunch of calculations here by reusing endpoints.
		// It's really only the midpoint that needs to keep getting calculated.
		left, right := math.Abs(distmid-distmin), math.Abs(distmid-distmax)
		if left <= right {
			stepRng.Min, stepRng.Max = stepRng.Min, midpoint
			change = left
		} else {
			stepRng.Min, stepRng.Max = midpoint, stepRng.Max
			change = right
		}

		count++
	}
	//	fmt.Println("count", count, "change", change)
	return distmid, ptmid
}

// https://www.shadertoy.com/view/MlKcDD
// https://iquilezles.org/articles/
/*
// The MIT License
// Copyright © 2018 Inigo Quilez
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software. THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Distance to a quadratic bezier segment, which can be solved analyically with a cubic.

// List of some other 2D distances: https://www.shadertoy.com/playlist/MXdSRf
//
// and iquilezles.org/articles/distfunctions2d


// 0: exact, using a cubic colver
// 1: approximated
#define METHOD 0



float dot2( vec2 v ) { return dot(v,v); }
float cro( vec2 a, vec2 b ) { return a.x*b.y-a.y*b.x; }
float cos_acos_3( float x ) { x=sqrt(0.5+0.5*x); return x*(x*(x*(x*-0.008972+0.039071)-0.107074)+0.576975)+0.5; } // https://www.shadertoy.com/view/WltSD7

#if METHOD==0
// signed distance to a quadratic bezier
float sdBezier( in vec2 pos, in vec2 A, in vec2 B, in vec2 C, out vec2 outQ )
{
    vec2 a = B - A;
    vec2 b = A - 2.0*B + C;
    vec2 c = a * 2.0;
    vec2 d = A - pos;

    // cubic to be solved (kx*=3 and ky*=3)
    float kk = 1.0/dot(b,b);
    float kx = kk * dot(a,b);
    float ky = kk * (2.0*dot(a,a)+dot(d,b))/3.0;
    float kz = kk * dot(d,a);

    float res = 0.0;
    float sgn = 0.0;

    float p  = ky - kx*kx;
    float q  = kx*(2.0*kx*kx - 3.0*ky) + kz;
    float p3 = p*p*p;
    float q2 = q*q;
    float h  = q2 + 4.0*p3;


    if( h>=0.0 )
    {   // 1 root
        h = sqrt(h);
        vec2 x = (vec2(h,-h)-q)/2.0;

        #if 0
        // When p≈0 and p<0, h-q has catastrophic cancelation. So, we do
        // h=√(q²+4p³)=q·√(1+4p³/q²)=q·√(1+w) instead. Now we approximate
        // √ by a linear Taylor expansion into h≈q(1+½w) so that the q's
        // cancel each other in h-q. Expanding and simplifying further we
        // get x=vec2(p³/q,-p³/q-q). And using a second degree Taylor
        // expansion instead: x=vec2(k,-k-q) with k=(1-p³/q²)·p³/q
        if( abs(p)<0.001 )
        {
            float k = p3/q;              // linear approx
          //float k = (1.0-p3/q2)*p3/q;  // quadratic approx
            x = vec2(k,-k-q);
        }
        #endif

        vec2 uv = sign(x)*pow(abs(x), vec2(1.0/3.0));
        float t = uv.x + uv.y;

		// from NinjaKoala - single newton iteration to account for cancellation
        t -= (t*(t*t+3.0*p)+q)/(3.0*t*t+3.0*p);

        t = clamp( t-kx, 0.0, 1.0 );
        vec2  w = d+(c+b*t)*t;
        outQ = w + pos;
        res = dot2(w);
    	sgn = cro(c+2.0*b*t,w);
    }
    else
    {   // 3 roots
        float z = sqrt(-p);
        #if 0
        float v = acos(q/(p*z*2.0))/3.0;
        float m = cos(v);
        float n = sin(v);
        #else
        float m = cos_acos_3( q/(p*z*2.0) );
        float n = sqrt(1.0-m*m);
        #endif
        n *= sqrt(3.0);
        vec3  t = clamp( vec3(m+m,-n-m,n-m)*z-kx, 0.0, 1.0 );
        vec2  qx=d+(c+b*t.x)*t.x; float dx=dot2(qx), sx=cro(a+b*t.x,qx);
        vec2  qy=d+(c+b*t.y)*t.y; float dy=dot2(qy), sy=cro(a+b*t.y,qy);
        if( dx<dy ) {res=dx;sgn=sx;outQ=qx+pos;} else {res=dy;sgn=sy;outQ=qy+pos;}
    }

    return sqrt( res )*sign(sgn);
}
#else

// This method provides just an approximation, and is only usable in
// the very close neighborhood of the curve. Taken and adapted from
// http://research.microsoft.com/en-us/um/people/hoppe/ravg.pdf
float sdBezier( vec2 p, vec2 v0, vec2 v1, vec2 v2, out vec2 outQ )
{
	vec2 i = v0 - v2;
    vec2 j = v2 - v1;
    vec2 k = v1 - v0;
    vec2 w = j-k;

	v0-= p; v1-= p; v2-= p;

	float x = cro(v0, v2);
    float y = cro(v1, v0);
    float z = cro(v2, v1);

	vec2 s = 2.0*(y*j+z*k)-x*i;

    float r =  (y*z-x*x*0.25)/dot2(s);
    float t = clamp( (0.5*x+y+r*dot(s,w))/(x+y+z),0.0,1.0);

    vec2 d = v0+t*(k+k+t*w);
    outQ = d + p;
	return length(d);
}
#endif

float udSegment( in vec2 p, in vec2 a, in vec2 b )
{
	vec2 pa = p - a;
	vec2 ba = b - a;
	float h = clamp( dot(pa,ba)/dot(ba,ba), 0.0, 1.0 );
	return length( pa - ba*h );
}

void mainImage( out vec4 fragColor, in vec2 fragCoord )
{
	vec2 p = (2.0*fragCoord-iResolution.xy)/iResolution.y;
    vec2 m = (2.0*iMouse.xy-iResolution.xy)/iResolution.y;

	vec2 v0 = vec2(1.3,0.9)*cos(iTime*0.5 + vec2(0.0,5.0) );
    vec2 v1 = vec2(1.3,0.9)*cos(iTime*0.6 + vec2(3.0,4.0) );
    vec2 v2 = vec2(1.3,0.9)*cos(iTime*0.7 + vec2(2.0,0.0) );

    vec2 kk;

    float d = sdBezier( p, v0,v1,v2, kk );

    float f = smoothstep(-0.2,0.2,cos(2.0*iTime));
    vec3 col = (d>0.0) ? vec3(0.9,0.6,0.3) : vec3(0.65,0.85,1.0);
	col *= 1.0 - exp(-4.0*abs(d));
	col *= 0.8 + 0.2*cos(110.0*d);
	col = mix( col, vec3(1.0), 1.0-smoothstep(0.0,0.015,abs(d)) );

    if( iMouse.z>0.001 )
    {
    vec2 q;
    d = sdBezier(m, v0,v1,v2, q );
    col = mix(col, vec3(1.0,1.0,0.0), 1.0-smoothstep(0.0, 0.005, abs(length(p-m)-abs(d))-0.0025));
    col = mix(col, vec3(1.0,1.0,0.0), 1.0-smoothstep(0.0, 0.005, length(p-m)-0.015));
    col = mix(col, vec3(1.0,1.0,0.0), 1.0-smoothstep(0.0, 0.005, length(p-q)-0.015));
    }

    if( cos(0.5*iTime)<-0.5 )
    {
        d = min( udSegment(p,v0,v1),
                 udSegment(p,v1,v2) );
        d = min( d, length(p-v0)-0.02 );
        d = min( d, length(p-v1)-0.02 );
        d = min( d, length(p-v2)-0.02 );
        col = mix( col, vec3(1,0,0), 1.0-smoothstep(0.0,0.007,d) );
    }

	fragColor = vec4(col,1.0);
}
*/

// QuadraticBezierDistanceBrute approximates the distance from a point to a Bezier curve
func QuadraticBezierDistanceBrute(cb QuadraticBezier, pt PtF) (float64, PtF) {
	d, p := QuadraticBezierDistanceSquaredBrute(cb, pt)
	return math.Sqrt(d), p
}

// QuadraticBezierDistanceSquaredBrute approximates the distance from a point to a Bezier curve
func QuadraticBezierDistanceSquaredBrute(cb QuadraticBezier, pt PtF) (float64, PtF) {
	// simple binary search-like algo
	stepRng := Rng(0., 1.)
	count := 0
	change := 9999999.
	ptmid := PtF{}
	distmid := 0.
	for count < 100 && change > .001 {
		// Get the endpoints for the range and find which one is closer
		// to the current midpoint.
		midpoint := stepRng.Midpoint()
		ptmin := cb.At(stepRng.Min)
		distmin := pt.DistSquared(ptmin)
		ptmax := cb.At(stepRng.Max)
		distmax := pt.DistSquared(ptmax)
		ptmid = cb.At(midpoint)
		distmid = pt.DistSquared(ptmid)

		// Could obviously remove a bunch of calculations here by reusing endpoints.
		// It's really only the midpoint that needs to keep getting calculated.
		left, right := math.Abs(distmid-distmin), math.Abs(distmid-distmax)
		if left <= right {
			stepRng.Min, stepRng.Max = stepRng.Min, midpoint
			change = left
		} else {
			stepRng.Min, stepRng.Max = midpoint, stepRng.Max
			change = right
		}

		count++
	}
	//	fmt.Println("count", count, "change", change)
	return distmid, ptmid
}

type bezierDistanceFound struct {
	dist float64
	pt   PtF
}
