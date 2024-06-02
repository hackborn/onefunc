package sys

import (
	"github.com/hackborn/onefunc/math/geo"
)

func Get(keys ...string) (Info, error) {
	return get(keys...)
}

type Info struct {
	// Path to application data folder.
	AppDataPath string

	// The system-reported platform DPI. Note that this might
	// not be a final value: In some cases it might need to
	// be multipled by the scale.
	Dpi geo.PointF64

	// The current screen scaling. Will be 1 for no scale.
	Scale float64
}

const (
	AppDataPath = "appdatapath"
	Dpi         = "dpi"
	Scale       = "scale"
)
