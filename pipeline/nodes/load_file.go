package nodes

import (
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

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

	// Optional separator to split the glob into multiple
	// patterns. Seems like you should be able to do this
	// with {} but it's not working for me, and also it's
	// a little confusing.
	Separator string
}

func (n *LoadFileNode) Start(input pipeline.StartInput) error {
	data := n.loadFileData
	input.SetNodeData(&data)
	return nil
}

func (n *LoadFileNode) Run(state *pipeline.State, input pipeline.RunInput, output *pipeline.RunOutput) error {
	data := state.NodeData.(*loadFileData)

	get, read, err := n.prepare(data)
	if err != nil {
		return err
	}

	eb := &oferrors.FirstBlock{}
	matches, err := n.getMatches(data.Glob, data.Separator, get)
	eb.AddError(err)
	//	fmt.Println("fs matches", matches, err)
	for _, fn := range matches {
		dat, err := read(fn)
		eb.AddError(err)

		base := path.Base(fn)
		output.Pins = append(output.Pins, pipeline.Pin{Payload: &pipeline.ContentData{Name: base, Data: string(dat)}})
	}
	return eb.Err
}

func (n *LoadFileNode) prepare(data *loadFileData) (loadGetMatches, loadReadFile, error) {
	if data.Fs == "" {
		get := func(glob string) ([]string, error) {
			return filepath.Glob(filepath.FromSlash(glob))
		}
		read := func(path string) ([]byte, error) {
			return os.ReadFile(path)
		}
		return get, read, nil
	} else {
		fsys, ok := pipeline.FindFs(data.Fs)
		if !ok {
			return nil, nil, fmt.Errorf("LoadFileNode: no registered filesystem named \"%v\"", data.Fs)
		}
		get := func(glob string) ([]string, error) {
			return fs.Glob(fsys, glob)
		}
		read := func(path string) ([]byte, error) {
			return fs.ReadFile(fsys, path)
		}
		return get, read, nil
	}
}

func (n *LoadFileNode) getMatches(glob, sep string, fn loadGetMatches) ([]string, error) {
	if sep == "" {
		return fn(glob)
	}
	var err error
	matches := []string{}
	globs := strings.Split(glob, sep)
	for _, g := range globs {
		g = strings.TrimSpace(g)
		m, e := fn(g)
		err = cmp.Or(err, e)
		if len(m) > 0 {
			matches = append(matches, m...)
		}
	}
	return matches, err
}

// loadGetMatches gets matches for a glob.
type loadGetMatches func(glob string) ([]string, error)

// loadReaedFile reads a filename. It is done this way
// so I can wrap os.ReadFile instead of using os.DirFS,
// which does not handle relative paths.
type loadReadFile func(path string) ([]byte, error)
