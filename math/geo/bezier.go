package geo

// Bez is shorthand for creating a quadratic bezier curve.
func Bez(x0, y0, x1, y1, x2, y2 float64) QuadraticBezier {
	return QuadraticBez(x0, y0, x1, y1, x2, y2)
}

// QuadraticBez is shorthand for creating a quadratic bezier curve.
func QuadraticBez(x0, y0, x1, y1, x2, y2 float64) QuadraticBezier {
	return QuadraticBezier{P0: Pt(x0, y0),
		P1: Pt(x1, y1),
		P2: Pt(x2, y2),
	}
}

// CubicBez is shorthand for creating a cubic bezier curve.
func CubicBez(x0, y0, x1, y1, x2, y2, x3, y3 float64) CubicBezier {
	return CubicBezier{P0: Pt(x0, y0),
		P1: Pt(x1, y1),
		P2: Pt(x2, y2),
		P3: Pt(x3, y3),
	}
}

// ---------------------------------------------------------
// QUADRATIC BEZIER

// QuadraticBezier represents a quadratic Bezier curve.
type QuadraticBezier struct {
	P0, P1, P2 PtF
}

// At evaluates the Bezier curve at a given parameter t.
func (b *QuadraticBezier) At(t float64) PtF {
	x := (1-t)*(1-t)*b.P0.X + 2*(1-t)*t*b.P1.X + t*t*b.P2.X
	y := (1-t)*(1-t)*b.P0.Y + 2*(1-t)*t*b.P1.Y + t*t*b.P2.Y
	return PtF{X: x, Y: y}
}

// PointBounds is the bounding box for my control points. It does
// not include any point that lies outside the controls.
func (b *QuadraticBezier) PointBounds() RectF {
	r := RectF{}
	r.L, r.R = min(b.P0.X, b.P1.X), max(b.P0.X, b.P1.X)
	r.L, r.R = min(r.L, b.P2.X), max(r.R, b.P2.X)
	r.T, r.B = min(b.P0.Y, b.P1.Y), max(b.P0.Y, b.P1.Y)
	r.T, r.B = min(r.T, b.P2.Y), max(r.B, b.P2.Y)
	return r
}

// ---------------------------------------------------------
// CUBIC BEZIER

// CubicBezier represents a cubic Bezier curve.
type CubicBezier struct {
	P0, P1, P2, P3 PtF
}

// At evaluates the Bezier curve at a given parameter t.
func (bc *CubicBezier) At(t float64) PtF {
	bt := 1 - t
	bt2 := bt * bt
	bt3 := bt * bt2
	t2 := t * t
	t3 := t * t2

	x := bt3*bc.P0.X + 3*t*bt2*bc.P1.X + 3*t2*bt*bc.P2.X + t3*bc.P3.X
	y := bt3*bc.P0.Y + 3*t*bt2*bc.P1.Y + 3*t2*bt*bc.P2.Y + t3*bc.P3.Y

	return PtF{X: x, Y: y}
}

// PointBounds is the bounding box for my control points. It does
// not include any point that lies outside the controls.
func (bc *CubicBezier) PointBounds() RectF {
	r := RectF{}
	r.L, r.R = min(bc.P0.X, bc.P1.X), max(bc.P0.X, bc.P1.X)
	r.L, r.R = min(r.L, bc.P2.X), max(r.R, bc.P2.X)
	r.L, r.R = min(r.L, bc.P3.X), max(r.R, bc.P3.X)
	r.T, r.B = min(bc.P0.Y, bc.P1.Y), max(bc.P0.Y, bc.P1.Y)
	r.T, r.B = min(r.T, bc.P2.Y), max(r.B, bc.P2.Y)
	r.T, r.B = min(r.T, bc.P3.Y), max(r.B, bc.P3.Y)
	return r
}

type BezF = QuadraticBezier

/*
https://stackoverflow.com/questions/2742610/closest-point-on-a-cubic-bezier-curve

import numpy as np

# Bezier Class representing a CUBIC bezier defined by four
# control points.
#
# at(t):            gets a point on the curve at t
# distance2(pt)      returns the closest distance^2 of
#                   pt and the curve
# closest(pt)       returns the point on the curve
#                   which is closest to pt
# maxes(pt)         plots the curve using matplotlib
class Bezier(object):
    exp3 = np.array([[3, 3], [2, 2], [1, 1], [0, 0]], dtype=np.float32)
    exp3_1 = np.array([[[3, 3], [2, 2], [1, 1], [0, 0]]], dtype=np.float32)
    exp4 = np.array([[4], [3], [2], [1], [0]], dtype=np.float32)
    boundaries = np.array([0, 1], dtype=np.float32)

    # Initialize the curve by assigning the control points.
    # Then create the coefficients.
    def __init__(self, points):
        assert isinstance(points, np.ndarray)
        assert points.dtype == np.float32
        self.points = points
        self.create_coefficients()

    # Create the coefficients of the bezier equation, bringing
    # the bezier in the form:
    # f(t) = a * t^3 + b * t^2 + c * t^1 + d
    #
    # The coefficients have the same dimensions as the control
    # points.
    def create_coefficients(self):
        points = self.points
        a = - points[0] + 3*points[1] - 3*points[2] + points[3]
        b = 3*points[0] - 6*points[1] + 3*points[2]
        c = -3*points[0] + 3*points[1]
        d = points[0]
        self.coeffs = np.stack([a, b, c, d]).reshape(-1, 4, 2)

    # Return a point on the curve at the parameter t.
    def at(self, t):
        if type(t) != np.ndarray:
            t = np.array(t)
        pts = self.coeffs * np.power(t, self.exp3_1)
        return np.sum(pts, axis = 1)

    # Return the closest DISTANCE (squared) between the point pt
    # and the curve.
    def distance2(self, pt):
        points, distances, index = self.measure_distance(pt)
        return distances[index]

    # Return the closest POINT between the point pt
    # and the curve.
    def closest(self, pt):
        points, distances, index = self.measure_distance(pt)
        return points[index]

    # Measure the distance^2 and closest point on the curve of
    # the point pt and the curve. This is done in a few steps:
    # 1     Define the distance^2 depending on the pt. I am
    #       using the squared distance because it is sufficient
    #       for comparing distances and doesn't have the over-
    #       head of an additional root operation.
    #       D(t) = (f(t) - pt)^2
    # 2     Get the roots of D'(t). These are the extremes of
    #       D(t) and contain the closest points on the unclipped
    #       curve. Only keep the minima by checking if
    #       D''(roots) > 0 and discard imaginary roots.
    # 3     Calculate the distances of the pt to the minima as
    #       well as the start and end of the curve and return
    #       the index of the shortest distance.
    #
    # This desmos graph is a helpful visualization.
    # https://www.desmos.com/calculator/ktglugn1ya
    def measure_distance(self, pt):
        coeffs = self.coeffs

        # These are the coefficients of the derivatives d/dx and d/(d/dx).
        da = 6*np.sum(coeffs[0][0]*coeffs[0][0])
        db = 10*np.sum(coeffs[0][0]*coeffs[0][1])
        dc = 4*(np.sum(coeffs[0][1]*coeffs[0][1]) + 2*np.sum(coeffs[0][0]*coeffs[0][2]))
        dd = 6*(np.sum(coeffs[0][0]*(coeffs[0][3]-pt)) + np.sum(coeffs[0][1]*coeffs[0][2]))
        de = 2*(np.sum(coeffs[0][2]*coeffs[0][2])) + 4*np.sum(coeffs[0][1]*(coeffs[0][3]-pt))
        df = 2*np.sum(coeffs[0][2]*(coeffs[0][3]-pt))

        dda = 5*da
        ddb = 4*db
        ddc = 3*dc
        ddd = 2*dd
        dde = de
        dcoeffs = np.stack([da, db, dc, dd, de, df])
        ddcoeffs = np.stack([dda, ddb, ddc, ddd, dde]).reshape(-1, 1)

        # Calculate the real extremes, by getting the roots of the first
        # derivativ of the distance function.
        extrema = np_real_roots(dcoeffs)
        # Remove the roots which are out of bounds of the clipped range [0, 1].
        # [future reference] https://stackoverflow.com/questions/47100903/deleting-every-3rd-element-of-a-tensor-in-tensorflow
        dd_clip = (np.sum(ddcoeffs * np.power(extrema, self.exp4)) >= 0) & (extrema > 0) & (extrema < 1)
        minima = extrema[dd_clip]

        # Add the start and end position as possible positions.
        potentials = np.concatenate((minima, self.boundaries))

        # Calculate the points at the possible parameters t and
        # get the index of the closest
        points = self.at(potentials.reshape(-1, 1, 1))
        distances = np.sum(np.square(points - pt), axis = 1)
        index = np.argmin(distances)

        return points, distances, index


    # Point the curve to a matplotlib figure.
    # maxes         ... the axes of a matplotlib figure
    def plot(self, maxes):
        import matplotlib.path as mpath
        import matplotlib.patches as mpatches
        Path = mpath.Path
        pp1 = mpatches.PathPatch(
            Path(self.points, [Path.MOVETO, Path.CURVE4, Path.CURVE4, Path.CURVE4]),
            fc="none")#, transform=ax.transData)
        pp1.set_alpha(1)
        pp1.set_color('#00cc00')
        pp1.set_fill(False)
        pp2 = mpatches.PathPatch(
            Path(self.points, [Path.MOVETO, Path.LINETO , Path.LINETO , Path.LINETO]),
            fc="none")#, transform=ax.transData)
        pp2.set_alpha(0.2)
        pp2.set_color('#666666')
        pp2.set_fill(False)

        maxes.scatter(*zip(*self.points), s=4, c=((0, 0.8, 1, 1), (0, 1, 0.5, 0.8), (0, 1, 0.5, 0.8),
                                                  (0, 0.8, 1, 1)))
        maxes.add_patch(pp2)
        maxes.add_patch(pp1)

# Wrapper around np.roots, but only returning real
# roots and ignoring imaginary results.
def np_real_roots(coefficients, EPSILON=1e-6):
    r = np.roots(coefficients)
    return r.real[abs(r.imag) < EPSILON]

*/
