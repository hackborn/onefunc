//go:build !darwin

package sys

import (
	"fmt"
	"runtime"
)

func get(...string) (Info, error) {
	return Info{}, fmt.Errorf("platform.Get() unimplemented for %v", runtime.GOOS)
}
