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

func (r *SliceReader) NextPoint() (PtF, error) {
	if r.current >= len(r.Pts) {
		return PtF{}, io.EOF
	}
	i := r.current
	r.current++
	return r.Pts[i], nil
}

// ---------------------------------------------------------
// INTERPOLATOR-READER

// InterpolatorReader convers a PointInterpolator into a reader,
// supplying all interpolated points from 0 - 1 based on step.
type InterpolatorReader struct {
	Source PointInterpolator
	Step   float64

	current float64
	done    bool
}

func (r *InterpolatorReader) NextPoint() (PtF, error) {
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
