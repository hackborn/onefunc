package nodes

import (
	"os"
	"path/filepath"

	"github.com/hackborn/onefunc/pipeline"
)

type LoadFile struct {
	Glob string
}

func (n *LoadFile) Run(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	filenames, err := n.filenamesGlob(n.Glob)
	if err != nil {
		return nil, err
	}
	output := pipeline.RunOutput{}
	for _, fn := range filenames {
		dat, err := os.ReadFile(fn)
		if err != nil {
			return nil, err
		}
		output.Pins = append(output.Pins, &pipeline.ContentData{Name: fn, Data: string(dat)})
	}
	return &output, nil
}

func (n *LoadFile) filenamesGlob(glob string) ([]string, error) {
	return filepath.Glob(glob)
}
