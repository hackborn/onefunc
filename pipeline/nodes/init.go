package nodes

import (
	"github.com/hackborn/onefunc/pipeline"
)

func init() {
	pipeline.RegisterNode("loadfile", func() pipeline.Node {
		return &LoadFile{}
	})
}
