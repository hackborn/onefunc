package nodes

import (
	astpkg "go/ast"
	parserpkg "go/parser"
	tokenpkg "go/token"
	"strings"

	"github.com/hackborn/onefunc/pipeline"
)

// StructNode takes ContentData with source code and converts it
// to StructData.
type StructNode struct {
	structData
}

type structData struct {
	// Optional. Supply a name to extract from the tag field. For example,
	// if a field as the tag `json:"..."` and you supply a Tag of "json" then
	// the StructData field will have a tag value of "...".
	Tag string
}

func (n *StructNode) Start(state *pipeline.State) {
	data := n.structData
	state.NodeData = &data
}

func (n *StructNode) Run(state *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	data := state.NodeData.(*structData)
	output := pipeline.RunOutput{}
	for _, pin := range input.Pins {
		switch t := pin.Payload.(type) {
		case *pipeline.ContentData:
			err := n.runContent(data, t, &output)
			if err != nil {
				return nil, err
			}
		}
	}
	return &output, nil
}

func (n *StructNode) runContent(data *structData, pin *pipeline.ContentData, output *pipeline.RunOutput) error {
	fset := tokenpkg.NewFileSet()
	f, err := parserpkg.ParseFile(fset, "", pin.Data, 0)
	if err != nil {
		return err
	}
	return n.readAst(data, f, output)
}

func (n *StructNode) readAst(data *structData, file *astpkg.File, output *pipeline.RunOutput) error {
	astpkg.Inspect(file, func(node astpkg.Node) bool {
		// Handle struct declarations
		if t, ok := node.(*astpkg.TypeSpec); ok {
			//				fmt.Printf("READ %v type %T type type %v %T\n", t, t, t.Type, t.Type)
			switch tt := t.Type.(type) {
			case *astpkg.StructType:
				output.Pins = append(output.Pins, pipeline.Pin{Payload: newStructData(t, tt, data.Tag)})
			}
		}
		return true
	})
	return nil
}

func newStructData(spec *astpkg.TypeSpec, structType *astpkg.StructType, tagFilter string) *pipeline.StructData {
	name := spec.Name.Name
	fields := make([]pipeline.StructField, 0, len(structType.Fields.List))
	unexportedFields := make([]pipeline.StructField, 0, len(structType.Fields.List))

	// Iterate over struct fields
	for _, field := range structType.Fields.List {
		sf := pipeline.StructField{Name: field.Names[0].Name,
			Type: field.Type.(*astpkg.Ident).Name}
		tag := ""

		if field.Tag != nil {
			// Not totally sure why it supplies it with the backticks, hopefully
			// this doesn't breaking something.
			tag = strings.Trim(field.Tag.Value, "`")
			if tag != "" && tagFilter != "" {
				tag = filterTag(tag, tagFilter)
			}
		}
		sf.Tag = tag
		if astpkg.IsExported(sf.Name) {
			fields = append(fields, sf)
		} else {
			unexportedFields = append(unexportedFields, sf)
		}
	}
	return pipeline.NewStructData(name, fields, unexportedFields)
}

// filterTag takes a tag string and finds the filter
// keywords, returning only the items in the keyword.
func filterTag(tag, filter string) string {
	filter += ":"
	idx := strings.Index(tag, filter)
	if idx < 0 {
		return ""
	}
	start := idx + len(filter) + 1
	nextIdx := strings.Index(tag[start:], `"`)
	if nextIdx < 0 {
		return ""
	}
	nextIdx = start + nextIdx
	return tag[start:nextIdx]
}
