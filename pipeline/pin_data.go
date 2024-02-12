package pipeline

func NewStructData(name string, fields []StructField, unexportedFields []StructField) *StructData {
	return &StructData{Name: name, Fields: fields, UnexportedFields: unexportedFields}
}

// StructData provides information about a single struct from
// source data.
type StructData struct {
	// The name of the source struct.
	Name string

	// All exported fields in the struct.
	Fields []StructField

	// All unexported fields in the struct. It should be rare
	// for a client to want access to these fields, so we keep
	// them in a separate slice to avoid polluting the fields list.
	UnexportedFields []StructField
}

func (d *StructData) Clone() Cloner {
	dst := *d
	return &dst
}

type StructField struct {
	// The name of the field in the original source.
	Name string
	// The Go type of the field (string, float64, etc.)
	Type string
	// Tag data for the field.
	Tag string
}

// ContentData provides a generic content string.
type ContentData struct {
	Name   string
	Data   string
	Format string
}

func (d *ContentData) Clone() Cloner {
	dst := *d
	return &dst
}
