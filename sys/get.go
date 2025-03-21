package sys

import (
	"fmt"

	"github.com/hackborn/onefunc/math/geo"
)

func Get(keys ...string) (Info, error) {
	return get(keys...)
}

func GetString(key string) (string, error) {
	info, err := get([]string{key}...)
	if err != nil {
		return "", err
	}
	switch key {
	case AppPath:
		return info.AppPath, nil
	case AppDocumentsPath:
		return info.AppDocumentsPath, nil
	case AppCachePath:
		return info.AppCachePath, nil
	case HardwareModel:
		return info.HardwareModel, nil
	default:
		return "", fmt.Errorf("Unknown key: \"%v\"", key)
	}
}

type Info struct {
	// Path to application data folder (not necessarily where the app is
	// running, but the system-defined location where the app stores data).
	// Uses filepath separator. Root path for all app paths.
	AppPath string

	// Path to application user data folder. Uses filepath separator.
	// Hardcoded on some platforms, generated on others.
	AppDocumentsPath string

	// Path to application cache folder. Uses filepath separator.
	// Hardcoded on some platforms, generated on others.
	AppCachePath string

	// HadrwareModel of this device.
	HardwareModel string

	// The system-reported platform DPI. Note that this might
	// not be a final value: In some cases it might need to
	// be multipled by the scale.
	Dpi geo.PtF

	// The current screen scaling. Will be 1 for no scale.
	Scale float64
}
