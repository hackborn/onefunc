package nodes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hackborn/onefunc/pipeline"
)

// SaveFileNode handles ContentData by saving to a file.
type SaveFileNode struct {
	saveFileData
}

type saveFileData struct {
	// Path gets prepended to all save files.
	Path string
}

func (n *SaveFileNode) Start(input pipeline.StartInput) error {
	data := n.saveFileData
	input.SetNodeData(&data)
	return nil
}

func (n *SaveFileNode) Run(state *pipeline.State, input pipeline.RunInput, output *pipeline.RunOutput) error {
	data := state.NodeData.(*saveFileData)
	path := filepath.FromSlash(data.Path)
	// fmt.Println("save", path, "inputs", len(input.Pins), "flush", state.Flush)
	err := n.verify(path)
	if err != nil {
		return err
	}
	for _, pin := range input.Pins {
		switch p := pin.Payload.(type) {
		case *pipeline.ContentData:
			err = n.runContentPin(state, p, path)
			if err != nil {
				return err
			}
		}
		output.Pins = append(output.Pins, pin)
	}
	return nil
}

func (n *SaveFileNode) verify(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("SaveFileNode: path \"" + path + "\"does not exist")
	}
	return nil
}

func (n *SaveFileNode) runContentPin(state *pipeline.State, pin *pipeline.ContentData, path string) error {
	if pin.Name == "" {
		return fmt.Errorf("SaveFileNode: pin supplied with no name")
	}
	fn := filepath.Join(path, pin.Name)
	content := []byte(pin.Data)
	return os.WriteFile(fn, content, 0644)
}
