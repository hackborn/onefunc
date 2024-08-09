package geo

import (
	"fmt"
	"strings"
)

func Fmt(a any, params any) fmt.Stringer {
	switch t := a.(type) {
	case SegF:
		return SegFmtF{Seg: t}
	case []SegF:
		return SliceSegFmtF{Segs: t}
	default:
		return noFmt{}
	}
}

type noFmt struct {
}

func (f noFmt) String() string {
	return "unknown"
}

type SegFmtF struct {
	Seg SegF
}

func (f SegFmtF) String() string {
	return fmt.Sprintf("(%.4f, %.4f) - (%.4f, %.4f)", f.Seg.A.X, f.Seg.A.Y, f.Seg.B.X, f.Seg.B.Y)
}

type SliceSegFmtF struct {
	Segs []SegF
}

func (f SliceSegFmtF) String() string {
	sb := &strings.Builder{}
	sb.WriteString("[\n")
	for _, seg := range f.Segs {
		sb.WriteString("\t")
		sb.WriteString(Fmt(seg, nil).String())
		sb.WriteString("\n")
	}
	sb.WriteString("]")
	return sb.String()
}
