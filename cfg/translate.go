package cfg

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/hackborn/onefunc/math/geo"
	"github.com/hackborn/onefunc/reflect"
	"github.com/hackborn/onefunc/sync"
)

// A series of additional translations, separated from the main
// settings file to cut down on some clutter.

// MustNRGBA answers the NRGBA value at the given path or
// fallback if path is absent.
// Path must point to a hex value containing rgba components,
// i.e. "#ffffff" or "#ffffffff"
func (s Settings) MustNRGBA(path string, r, g, b, a uint8) color.NRGBA {
	if v, ok := s.NRGBA(path); ok {
		return v
	}
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// NRGBA answers the NRGBA value at the given path.
// Path must point to a hex value containing rgba components,
// i.e. "#ffffff" or "#ffffffff"
func (s Settings) NRGBA(path string) (color.NRGBA, bool) {
	if v, ok := s.String(path); ok && strings.HasPrefix(v, "#") {
		v = v[1:]
		return color.NRGBA{R: hexToUint8(v, 0, 0),
			G: hexToUint8(v, 2, 0),
			B: hexToUint8(v, 4, 0),
			A: hexToUint8(v, 6, 255),
		}, true
	}
	return color.NRGBA{}, false
}

// MustRectF returns a rect. The path must point to the
// parent of keys that can be "l", "t", "r", "b", or some
// combination (i.e. "lr").
func (s Settings) MustRectF(path string, r geo.RectF) geo.RectF {
	if v, ok := s.RectF(path); ok {
		return v
	}
	return r
}

// MustRectF returns a rect. The path must point to the
// parent of keys that can be "l", "t", "r", "b", or some
// combination (i.e. "lr").
func (s Settings) RectF(path string) (geo.RectF, bool) {
	r := geo.RectF{}
	defer sync.Read(s.rw).Unlock()
	sub := s.lockedSubset(path)
	if len(sub.t) < 1 {
		return r, false
	}
	for k, _v := range sub.t {
		if v, err := reflect.GetFloat64(_v); err == nil {
			for _, rune := range strings.ToLower(k) {
				switch rune {
				case 'l':
					r.L = v
				case 't':
					r.T = v
				case 'r':
					r.R = v
				case 'b':
					r.B = v
				}
			}
		}
	}
	return r, true
}

func hexToUint8(s string, idx int, fallback uint8) uint8 {
	if len(s) >= idx+2 {
		if v, err := strconv.ParseInt(s[idx:idx+2], 16, 32); err == nil {
			return uint8(v)
		}
	}
	return fallback
}
