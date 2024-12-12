package io

import (
	"cmp"
	"os"
	"path/filepath"
	"strings"

	"github.com/hackborn/onefunc/sys"
)

// MustExpandPath answers the path, resolving any sys variables.
// It ignores any errors, which for some clients can't happen.
func MustExpandPath(path string) string {
	path, _ = ExpandPath(path)
	return path
}

// ExpandPath answers the path, resolving any sys variables.
func ExpandPath(path string) (string, error) {
	var err error
	// Replace variables
	for _, v := range sysvars {
		if strings.Contains(path, v) {
			if s, err2 := sys.GetString(strings.Trim(v, "$")); err == nil {
				path = strings.ReplaceAll(path, v, s)
			} else {
				err = cmp.Or(err, err2)
			}
		}
	}
	return path, err
}

// MkExpandPath answers the path, resolving any sys variables and
// guaranteeing the parent dir exists. It assumes the leaf is
// a file.
// Variables: Anything defined in sys.
func MkExpandPath(path string) (string, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return path, err
	}
	// Guarantee parent dirs
	parent := filepath.Dir(path)
	return path, os.MkdirAll(parent, 0700)
}
