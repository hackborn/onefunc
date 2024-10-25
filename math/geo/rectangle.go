package geo

import (
	"fmt"
)

func Rect[T Number](left, top, right, bottom T) Rectangle[T] {
	return Rectangle[T]{L: left, T: top, R: right, B: bottom}
}

type Rectangle[T Number] struct {
	L, T, R, B T
}

func (r Rectangle[T]) Width() T {
	return r.R - r.L
}

func (r Rectangle[T]) Height() T {
	return r.B - r.T
}

func (r Rectangle[T]) LT() Point[T] {
	return Pt(r.L, r.T)
}

func (r Rectangle[T]) RB() Point[T] {
	return Pt(r.R, r.B)
}

func (r Rectangle[T]) Translate(pt Point[T]) Rectangle[T] {
	return Rectangle[T]{L: r.L + pt.X, T: r.T + pt.Y, R: r.R + pt.X, B: r.B + pt.Y}
}

func (r Rectangle[T]) Size() Point[T] {
	return Point[T]{X: r.R - r.L, Y: r.B - r.T}
}

func (r Rectangle[T]) WithSize(pt Point[T]) Rectangle[T] {
	return Rectangle[T]{L: r.L, T: r.T, R: r.L + pt.X, B: r.T + pt.Y}
}

func (r Rectangle[T]) Area() T {
	lt := r.LT()
	rb := r.RB()
	return (rb.X - lt.X) * (rb.Y - lt.Y)
}

func (r Rectangle[T]) Center() Point[T] {
	return Point[T]{X: (r.L + r.R) / 2, Y: (r.T + r.B) / 2}
}

func (r1 Rectangle[T]) Contains(r2 Rectangle[T]) bool {
	return r2.L >= r1.L && r2.T >= r1.T &&
		r2.R <= r1.R && r2.B <= r1.B
}

// ContainsPoint returns true if the point is inside the
// half-closed rectangle.
func (r Rectangle[T]) ContainsPoint(pt Point[T]) bool {
	return pt.X >= r.L && pt.X < r.R && pt.Y >= r.T && pt.Y < r.B
}

func (r1 Rectangle[T]) Overlaps(r2 Rectangle[T]) bool {
	if r1.R < r2.L || r2.R < r1.L {
		return false
	}
	if r1.B < r2.T || r2.B < r1.T {
		return false
	}
	return true
}

func (r1 Rectangle[T]) Union(r2 Rectangle[T]) Rectangle[T] {
	r1.L = min(r1.L, r2.L)
	r1.T = min(r1.T, r2.T)
	r1.R = max(r1.R, r2.R)
	r1.B = max(r1.B, r2.B)
	return r1
}

func (r Rectangle[T]) String() string {
	return fmt.Sprintf("l=%v t=%v r=%v b=%v", r.L, r.T, r.R, r.B)
}

// Expand adds the value to all edges.
func (r Rectangle[T]) WithExpand(v T) Rectangle[T] {
	return r.Add(-v, -v, v, v)
}

func (r1 Rectangle[T]) Add(l, t, r, b T) Rectangle[T] {
	return Rectangle[T]{L: r1.L + l,
		T: r1.T + t,
		R: r1.R + r,
		B: r1.B + b,
	}
}

// MergePoint will expand the bounds to include the point.
func (r1 Rectangle[T]) MergePoint(x, y T) Rectangle[T] {
	if x < r1.L {
		r1.L = x
	} else if x > r1.R {
		r1.R = x
	}
	if y < r1.T {
		r1.T = y
	} else if y > r1.B {
		r1.B = y
	}
	return r1
}

func ConvertRect[A Number, B Number](a Rectangle[A]) Rectangle[B] {
	return Rectangle[B]{L: B(a.L),
		T: B(a.T),
		R: B(a.R),
		B: B(a.B)}
}

type RectF = Rectangle[float64]
type RectI = Rectangle[int]

type RectF32 = Rectangle[float32]
type RectF64 = Rectangle[float64]
type RectI64 = Rectangle[int64]
