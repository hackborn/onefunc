package io

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// A resource is an io.ReadCloser. Provide some
// conveniences for generating them.

// FileResource answers a new resource on the file. If
// there is an error, the response will also be an error
// type (or the error will just be returned when you try to use it).
//
// You might be rightly thinking this should just return an
// error argument. I wouldn't argue. It's handled this way to
// eliminate an error check -- all code should be checking
// whatever is using the resource for an error, making it
// redundant. I don't love it, but error checking in Go is
// kinda crazy and I am elways experimenting with ways to
// reduce the checks.
func FileResource(filename string) io.ReadCloser {
	file, err := os.Open(filename)
	if err != nil {
		return &errorRes{err: err.Error()}
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

type errorRes struct {
	err string
}

func (r *errorRes) Error() string {
	return r.err
}

func (r *errorRes) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("%s", r.err)
}

func (r *errorRes) Close() error {
	return nil
}
