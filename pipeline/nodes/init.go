package nodes

import (
	"github.com/hackborn/onefunc/pipeline"
)

func init() {
	pipeline.RegisterNode("fmt", func() pipeline.Node {
		return &FmtNode{}
	})
	pipeline.RegisterNode("load", func() pipeline.Node {
		return &LoadFileNode{}
	})
	pipeline.RegisterNode("regexp", func() pipeline.Node {
		return &RegexpNode{}
	})
	pipeline.RegisterNode("save", func() pipeline.Node {
		return &SaveFileNode{}
	})
	pipeline.RegisterNode("struct", func() pipeline.Node {
		return &StructNode{}
	})
}
