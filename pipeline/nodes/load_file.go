package nodes

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	oferrors "github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/pipeline"
)

type LoadFileNode struct {
	// Name of the filesystem to read from. If this is an
	// empty string, the local filesystem is used, otherwise
	// the name must have been previously registered with
	// RegisterFs.
	Fs string

	// Glob pattern used to select which files are loaded.
	Glob string
}

func (n *LoadFileNode) Run(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	if n.Fs == "" {
		return n.runLocal(s, input)
	}
	return n.runFs(s, input)
}

func (n *LoadFileNode) runFs(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	fsys, ok := pipeline.FindFs(n.Fs)
	if !ok {
		return nil, fmt.Errorf("LoadFileNode: no registered filesystem named \"%v\"", n.Fs)
	}
	eb := &oferrors.FirstBlock{}
	matches, err := fs.Glob(fsys, n.Glob)
	eb.AddError(err)

	output := pipeline.RunOutput{}
	for _, fn := range matches {
		dat, err := fs.ReadFile(fsys, fn)
		eb.AddError(err)

		base := path.Base(fn)
		output.Pins = append(output.Pins, pipeline.Pin{Payload: &pipeline.ContentData{Name: base, Data: string(dat)}})
	}
	return &output, eb.Err
}

func (n *LoadFileNode) runLocal(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	eb := &oferrors.FirstBlock{}
	matches, err := filepath.Glob(filepath.FromSlash(n.Glob))
	eb.AddError(err)

	output := pipeline.RunOutput{}
	for _, fn := range matches {
		dat, err := os.ReadFile(fn)
		eb.AddError(err)

		base := filepath.Base(fn)
		output.Pins = append(output.Pins, pipeline.Pin{Payload: &pipeline.ContentData{Name: base, Data: string(dat)}})
	}
	return &output, eb.Err
}

func (n *LoadFileNode) getFilenames(fsname string, glob string) ([]string, error) {
	if fsname == "" {
		return filepath.Glob(filepath.FromSlash(glob))
	}
	fsys, ok := pipeline.FindFs(fsname)
	if !ok {
		return nil, fmt.Errorf("LoadFileNode: no registered filesystem named \"%v\"", fsname)
	}
	return fs.Glob(fsys, glob)
}
