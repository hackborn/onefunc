package pipeline

func NewStructData(name string, fields []StructField) *StructData {
	return &StructData{Name: name, Fields: fields}
}

// StructData provides information about a single struct from
// source data.
type StructData struct {
	Name   string
	Fields []StructField
}

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
	Name   string
	Data   string
	Format string
}
