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
	loadFileData
}

type loadFileData struct {
	// Name of the filesystem to read from. If this is an
	// empty string, the local filesystem is used, otherwise
	// the name must have been previously registered with
	// RegisterFs.
	Fs string

	// Glob pattern used to select which files are loaded.
	Glob string
}

func (n *LoadFileNode) Start(input pipeline.StartInput) error {
	data := n.loadFileData
	input.SetNodeData(&data)
	return nil
}

func (n *LoadFileNode) Run(state *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	data := state.NodeData.(*loadFileData)
	if data.Fs == "" {
		return n.runLocal(data, input)
	}
	return n.runFs(data, input)
}

func (n *LoadFileNode) runFs(data *loadFileData, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	fsys, ok := pipeline.FindFs(data.Fs)
	if !ok {
		return nil, fmt.Errorf("LoadFileNode: no registered filesystem named \"%v\"", data.Fs)
	}
	eb := &oferrors.FirstBlock{}
	matches, err := fs.Glob(fsys, data.Glob)
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

func (n *LoadFileNode) runLocal(data *loadFileData, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	eb := &oferrors.FirstBlock{}
	matches, err := filepath.Glob(filepath.FromSlash(data.Glob))
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
