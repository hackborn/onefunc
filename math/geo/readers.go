package geo

import (
	"io"
)

// ---------------------------------------------------------
// SLICE-READER

// SliceReader wraps a point slice into a PointReader.
type SliceReader struct {
	Pts []PtF

	current int
}

func (r *SliceReader) ReadPoint() (PtF, error) {
	if r.current >= len(r.Pts) {
		return PtF{}, io.EOF
	}
	i := r.current
	r.current++
	return r.Pts[i], nil
}

// ---------------------------------------------------------
// INTERPOLATOR-READER

// InterpolatorReader converts a PointInterpolator into a reader,
// supplying all interpolated points from 0 - 1 based on step.
type InterpolatorReader struct {
	Source PointInterpolator
	Step   float64

	current float64
	done    bool
}

// SetSource will also reset the state.
func (r *InterpolatorReader) SetSource(src PointInterpolator) {
	r.Source = src
	r.Reset()
}

func (r *InterpolatorReader) Reset() {
	r.current = 0
	r.done = false
}

func (r *InterpolatorReader) ReadPoint() (PtF, error) {
	if r.done || r.Source == nil {
		return PtF{}, io.EOF
	}
	if r.current >= 1. {
		r.current = 1.
		r.done = true
	}
	pt := r.Source.PointAt(r.current)
	r.current += r.Step
	return pt, nil
}

// ---------------------------------------------------------
// SLICE-READER-3D

// SliceReader3d wraps a point slice into a PointReader.
type SliceReader3d struct {
	Pts []PtF

	current int
}

func (r *SliceReader3d) ReadPoint() (PtF, error) {
	if r.current >= len(r.Pts) {
		return PtF{}, io.EOF
	}
	i := r.current
	r.current++
	return r.Pts[i], nil
}

// ---------------------------------------------------------
// INTERPOLATOR-READER-3D

// InterpolatorReader3d converts a PointInterpolator3d into a reader,
// supplying all interpolated points from 0 - 1 based on step.
type InterpolatorReader3d struct {
	Source PointInterpolator3d
	Step   float64

	current float64
	done    bool
}

// SetSource will also reset the state.
func (r *InterpolatorReader3d) SetSource(src PointInterpolator3d) {
	r.Source = src
	r.Reset()
}

func (r *InterpolatorReader3d) Reset() {
	r.current = 0
	r.done = false
}

func (r *InterpolatorReader3d) ReadPoint() (Pt3dF, error) {
	if r.done || r.Source == nil {
		return Pt3dF{}, io.EOF
	}
	if r.current >= 1. {
		r.current = 1.
		r.done = true
	}
	pt := r.Source.PointAt(r.current)
	r.current += r.Step
	return pt, nil
}
