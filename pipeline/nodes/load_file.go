package nodes

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hackborn/onefunc/pipeline"
)

type LoadFileNode struct {
	// Name of the filesystem to load from. If this is an
	// empty string, the local filesystem is used. Otherwise
	// the name must have been previously registered with
	// RegisterFs.
	Fs string

	// Glob pattern used to select which files are loaded.
	Glob string
}

func (n *LoadFileNode) Run(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	filenames, err := n.getFilenames(n.Fs, n.Glob)
	if err != nil {
		return nil, err
	}
	output := pipeline.RunOutput{}
	for _, fn := range filenames {
		dat, err := os.ReadFile(fn)
		if err != nil {
			return nil, err
		}
		base := filepath.Base(fn)
		output.Pins = append(output.Pins, pipeline.Pin{Payload: &pipeline.ContentData{Name: base, Data: string(dat)}})
	}
	return &output, nil
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
