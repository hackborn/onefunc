package nodes

import (
	"os"
	"path/filepath"

	"github.com/hackborn/onefunc/pipeline"
)

type LoadFileNode struct {
	Glob string
}

func (n *LoadFileNode) Run(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
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
		base := filepath.Base(fn)
		output.Pins = append(output.Pins, pipeline.Pin{Payload: &pipeline.ContentData{Name: base, Data: string(dat)}})
	}
	return &output, nil
}

func (n *LoadFileNode) filenamesGlob(glob string) ([]string, error) {
	return filepath.Glob(filepath.FromSlash(glob))
}
