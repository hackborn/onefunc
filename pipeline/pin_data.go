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
	// The Type will be unknown for complicated types (maps, etc.).
	// In that case, you can use the RawType.
	// And yes, this is a bad design. It is a result of pushing
	// unhandled cases into this node instead of letting clients
	// decide for themselves. Future clients should just use RawType
	// and expect that at some point RawType will be renamed to Type.
	Type    string
	RawType string
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

const (
	UnknownType = "unknown"
)
