package io

import (
	"fmt"
	"io/fs"
)

// ReadString returns the first matching file as a string.
func ReadString(fsys fs.FS, globPattern string) (string, error) {
	d, err := Read(fsys, globPattern)
	if d == nil || err != nil {
		return "", err
	}
	return string(d), nil
}

// ReadString returns the contents of the first matching file.
func Read(fsys fs.FS, globPattern string) ([]byte, error) {
	matches, err := fs.Glob(fsys, globPattern)
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		dat, err := fs.ReadFile(fsys, match)
		if err != nil {
			return nil, err
		}
		return dat, nil
	}
	return nil, fmt.Errorf("No match for \"%v\"", globPattern)
}
