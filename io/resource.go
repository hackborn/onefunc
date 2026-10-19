package io

import (
	"io"
	"os"
	"strings"
)

// A resource is an io.ReadCloser. Provide some
// conveniences for generating them.

func FileResource(filename string) io.ReadCloser {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return file
}

func StringResource(s string) io.ReadCloser {
	return &stringRes{r: strings.NewReader(s)}
}

type stringRes struct {
	r io.Reader
}

func (r *stringRes) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *stringRes) Close() error {
	return nil
}
