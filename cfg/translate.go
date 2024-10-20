package cfg

import (
	"image/color"
	"strconv"
	"strings"
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

func hexToUint8(s string, idx int, fallback uint8) uint8 {
	if len(s) >= idx+2 {
		if v, err := strconv.ParseInt(s[idx:idx+2], 16, 32); err == nil {
			return uint8(v)
		}
	}
	return fallback
}
