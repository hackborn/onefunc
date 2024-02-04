package nodes

import (
	"github.com/hackborn/onefunc/pipeline"
)

func init() {
	pipeline.RegisterNode("fmt", func() pipeline.Runner {
		return &FmtNode{}
	})
	pipeline.RegisterNode("load", func() pipeline.Runner {
		return &LoadFileNode{}
	})
	pipeline.RegisterNode("regexp", func() pipeline.Runner {
		return &RegexpNode{}
	})
	pipeline.RegisterNode("save", func() pipeline.Runner {
		return &SaveFileNode{}
	})
	pipeline.RegisterNode("struct", func() pipeline.Runner {
		return &StructNode{}
	})
}
