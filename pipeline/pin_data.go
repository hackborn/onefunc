package pipeline

import (
// "go/ast"
)

// PinData is a single piece of data assigned to a pin.
// It supplies the name of the owning pin, if any.
type PinData interface {
	PinName() string
}

func NewStructData(name string, fields []StructField) *StructData {
	return &StructData{name: name, fields: fields}
}

// StructData provides information about a single struct from
// source data.
type StructData struct {
	pinName string
	name    string
	fields  []StructField
	// keys    []any
}

func (s *StructData) PinName() string {
	return s.pinName
}

func (s *StructData) Name() string {
	return s.name
}

func (s *StructData) Fields() []StructField {
	return s.fields
}

/*
// Keys is an ordered (if specified in the tag) list of field key names.
func (s *StructData) Keys() []any {
	return s.keys
}
*/
/*
func (s *StructData) Spec() *ast.TypeSpec {
	return s.spec
}
*/

/*
func (s *StructPin) FieldFor(tagName string) *Field {
	for _, f := range s.fields {
		if f.TagName == tagName {
			return &f
		}
	}
	return nil
}
*/

type StructField struct {
	// The name of the field in the original source.
	Name string
	// The Go type of the field (string, float64, etc.)
	Type string
	// Tag data for the field.
	Tag string
	// The name assigned by the tag (or the Name, if no tag).
	TagName string
}

// ContentData provides a generic content string.
type ContentData struct {
	pinName string
	Name    string
	Data    string
}

func (s *ContentData) PinName() string {
	return s.pinName
}
