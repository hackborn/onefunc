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

func (n *StructNode) Start(input pipeline.StartInput) error {
	data := n.structData
	input.SetNodeData(&data)
	return nil
}

func (n *StructNode) Run(state *pipeline.State, input pipeline.RunInput, output *pipeline.RunOutput) error {
	data := state.NodeData.(*structData)
	for _, pin := range input.Pins {
		switch t := pin.Payload.(type) {
		case *pipeline.ContentData:
			err := n.runContent(data, t, output)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
	//	fmt.Println("STRUCT", name)
	// Iterate over struct fields
	for _, field := range structType.Fields.List {
		//		fmt.Printf("FIELD %T %v\n", field, field)
		//		fmt.Printf("FIELD names %T %v\n", field.Names, field.Names)
		//		fmt.Printf("FIELD type %T %v\n", field.Type, field.Type)
		typeName := pipeline.UnknownType
		switch t := field.Type.(type) {
		case *astpkg.Ident:
			typeName = t.Name
		}

		sf := pipeline.StructField{Name: field.Names[0].Name,
			Type: typeName,
			Tag:  getTag(field, tagFilter),
		}

		if astpkg.IsExported(sf.Name) {
			fields = append(fields, sf)
		} else {
			unexportedFields = append(unexportedFields, sf)
		}
	}
	return pipeline.NewStructData(name, fields, unexportedFields)
}

func getTag(field *astpkg.Field, tagFilter string) string {
	if field.Tag == nil {
		return ""
	}
	// Not totally sure why it supplies it with the backticks, hopefully
	// this doesn't breaking something.
	tag := strings.Trim(field.Tag.Value, "`")
	if tag != "" && tagFilter != "" {
		tag = filterTag(tag, tagFilter)
	}
	return tag
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
