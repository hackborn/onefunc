package pipeline

import (
// "go/ast"
)

// PinData is a single piece of data assigned to a pin.
// It supplies the name of the owning pin, if any.
type PinData interface {
	PinName() string
}

// StructData provides information about a single struct from
// source data.
type StructData struct {
	pinName    string
	structName string
	//	fields     []Field
	keys []any
}

func (s *StructData) PinName() string {
	return s.pinName
}

func (s *StructData) StructName() string {
	return s.structName
}

/*
func (s *StructData) Fields() []Field {
	return s.fields
}
*/

// Keys is an ordered (if specified in the tag) list of field key names.
func (s *StructData) Keys() []any {
	return s.keys
}

/*
func (s *StructData) Spec() *ast.TypeSpec {
	return s.spec
}
*/

/*
func newSrcPin(spec *ast.TypeSpec) *StructData {
//	pin := &StructData{structName: spec.Name.Name, spec: spec}
	pin := &StructData{structName: spec.Name.Name}
	structType := spec.Type.(*ast.StructType)
	taggedKeys := []keyField{}

	// Iterate over struct fields
	for _, field := range structType.Fields.List {
		fieldName := field.Names[0].Name
		fieldType := field.Type.(*ast.Ident).Name
		tagName := fieldName
		//		keyName := ""

		if field.Tag != nil {
			_name, _key := parseTag(field.Tag.Value)
			if _name != "" {
				tagName = _name
			}
			if _key != "" {
				taggedKeys = append(taggedKeys, keyField{Name: tagName, Tag: _key})
			}
		}
		pin.fields = append(pin.fields, Field{StructName: fieldName, TagName: tagName, Type: fieldType})
	}
	// Create an ordered list of keys for the request.
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
	return pin
}

func (s *StructPin) FieldFor(tagName string) *Field {
	for _, f := range s.fields {
		if f.TagName == tagName {
			return &f
		}
	}
	return nil
}
*/

// ContentData provides a generic content string.
type ContentData struct {
	pinName string
	Name    string
	Data    string
}

func (s *ContentData) PinName() string {
	return s.pinName
}
