package nodes

import (
	"fmt"

	"github.com/hackborn/onefunc/pipeline"
)

// FmtNode
type FmtNode struct {
	Verbose bool
}

func (n *FmtNode) Run(state *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	flushing := ""
	if state.Flush {
		flushing = " (flush)"
	}
	fmt.Println("fmt input pins:", len(input.Pins), flushing)
	for _, pin := range input.Pins {
		switch p := pin.Payload.(type) {
		case *pipeline.ContentData:
			n.runContentPin(p)
		case *pipeline.StructData:
			n.runStructPin(p)
		default:
			fmt.Printf("unknown pin type: %T\n", p)
		}
	}
	return &pipeline.RunOutput{Pins: input.Pins}, nil
}

func (n *FmtNode) runContentPin(pin *pipeline.ContentData) {
	fmt.Println("ContentData name:", pin.Name, "data:")
	if n.Verbose {
		fmt.Println(pin.Data)
	} else if len(pin.Data) < 40 {
		fmt.Println(pin.Data)
	} else {
		fmt.Println(pin.Data[0:40] + "...")
	}
}

func (n *FmtNode) runStructPin(pin *pipeline.StructData) {
	fmt.Println("StructData name:", pin.Name, "fields:", pin.Fields)
}
