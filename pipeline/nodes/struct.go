package nodes

import (
	astpkg "go/ast"
	parserpkg "go/parser"
	tokenpkg "go/token"

	"github.com/hackborn/onefunc/pipeline"
)

// StructNode takes ContentData with source code and converts it
// to StructData.
type StructNode struct {
}

func (n *StructNode) Run(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	output := pipeline.RunOutput{}
	for _, pin := range input.Pins {
		switch t := pin.(type) {
		case *pipeline.ContentData:
			err := n.runContent(t, &output)
			if err != nil {
				return nil, err
			}
		}
	}
	return &output, nil
}

func (n *StructNode) runContent(pin *pipeline.ContentData, output *pipeline.RunOutput) error {
	fset := tokenpkg.NewFileSet()
	f, err := parserpkg.ParseFile(fset, "", pin.Data, 0)
	if err != nil {
		return err
	}
	return n.readAst(f, output)
}

func (n *StructNode) readAst(file *astpkg.File, output *pipeline.RunOutput) error {
	astpkg.Inspect(file, func(n astpkg.Node) bool {
		// Handle struct declarations
		if t, ok := n.(*astpkg.TypeSpec); ok {
			//				fmt.Printf("READ %v type %T type type %v %T\n", t, t, t.Type, t.Type)
			switch tt := t.Type.(type) {
			case *astpkg.StructType:
				output.Pins = append(output.Pins, newSructData(t, tt))
			}
		}
		return true
	})
	return nil
}

func newSructData(spec *astpkg.TypeSpec, structType *astpkg.StructType) *pipeline.StructData {
	name := spec.Name.Name
	fields := make([]pipeline.StructField, 0, len(structType.Fields.List))
	//	taggedKeys := []keyField{}

	// Iterate over struct fields
	for _, field := range structType.Fields.List {
		sf := pipeline.StructField{Name: field.Names[0].Name,
			Type: field.Type.(*astpkg.Ident).Name}
		//		fieldName := field.Names[0].Name
		tag := ""
		//		tagName := fieldName
		//		keyName := ""

		if field.Tag != nil {
			tag = field.Tag.Value
			/*
				_name, _key := parseTag(field.Tag.Value)
				if _name != "" {
					tagName = _name
				}
				if _key != "" {
					taggedKeys = append(taggedKeys, keyField{Name: tagName, Tag: _key})
				}
			*/
		}
		sf.Tag = tag
		fields = append(fields, sf)
	}
	// Create an ordered list of keys for the request.
	/*
		if len(taggedKeys) > 0 {
			sort.Slice(taggedKeys, func(i, j int) bool {
				return taggedKeys[i].Tag < taggedKeys[j].Tag
			})
			keys := []any{}
			for _, f := range taggedKeys {
				keys = append(keys, f.Name)
			}
			pin.keys = keys
		}
	*/
	return pipeline.NewStructData(name, fields)
}
