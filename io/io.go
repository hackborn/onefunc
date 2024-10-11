package io

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

// ReadString returns the first matching file as a string.
func ReadString(fsys fs.FS, globPattern string) (string, error) {
	d, err := Read(fsys, globPattern)
	if d == nil || err != nil {
		return "", err
	}
	return string(d), nil
}

// Read returns the byte of the first matching file.
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

// ReadJson reads the file and unmarshals to the type as JSON
func ReadJson[T any](path string) (T, error) {
	var t T
	dat, err := os.ReadFile(path)
	err = cmp.Or(err, json.Unmarshal(dat, &t))
	return t, err
}
