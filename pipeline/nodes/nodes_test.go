package nodes

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hackborn/onefunc/jacl"
	"github.com/hackborn/onefunc/pipeline"
)

func TestMain(m *testing.M) {
	setupTests()
	code := m.Run()
	os.Exit(code)
}

// ---------------------------------------------------------
// TEST-FILTER-TAG
func TestFilterTag(t *testing.T) {
	table := []struct {
		tag    string
		filter string
		need   string
	}{
		{"", "doc", ""},
		{"a", "doc", ""},
		{`json:"jj"`, "doc", ""},
		{`doc:"id`, "doc", ``},
		{`doc:"id"`, "doc", `id`},
		{`json:"jj" doc:"id"`, "doc", `id`},
		{`doc:"id" json:"jj"`, "doc", `id`},
		{`doc:"id, key"`, "doc", `id, key`},
	}
	for i, v := range table {
		have := filterTag(v.tag, v.filter)
		if have != v.need {
			t.Fatalf("TestFilterTag %v expected %v but got %v", i, v.need, have)
		}
	}
}

// ---------------------------------------------------------
// TEST-LOAD
func TestLoad(t *testing.T) {
	table := []struct {
		pipeline string
		cmp      []string
		wantErr  error
	}{
		{`graph (load(Glob="` + testDataShortGlob + `"))`, []string{`0/Payload/Data=a`, `1/Payload/Data=b`, `2/Payload/Data=c`}, nil},
	}
	for i, v := range table {
		p, err := pipeline.Compile(v.pipeline)
		if err != nil {
			t.Fatalf("TestLoadFile %v compile err %v", i, err)
		}
		output, haveErr := pipeline.Run(p, nil, nil)
		var cmpErr error
		if output != nil {
			cmpErr = jacl.Run(output.Pins, v.cmp...)
		}
		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestLoad %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestLoad %v has no error but exptected %v", i, v.wantErr)
		} else if cmpErr != nil {
			t.Fatalf("TestLoad %v comparison error: %v", i, cmpErr)
		}
	}
}

// ---------------------------------------------------------
// TEST-PIPELINE
func TestPipeline(t *testing.T) {
	table := []struct {
		pipeline string
		cmp      []string
		wantErr  error
	}{
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/{type}="*StructData"`, `1/Payload/{type}="*StructData"`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Fields/0/Type=string`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Name=Company`, `1/Payload/Name=Filing`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Fields/0/Name=Id`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Fields/0/Tag="doc:''id, key''"`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/UnexportedFields/0/Name="_private"`, `0/Payload/UnexportedFields/0/Tag="json:''-''"`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct(Tag=doc))`, []string{`0/Payload/Fields/0/Tag="id, key"`}, nil},
		{`graph (load(Glob="` + testDataShortGlob + `"))`, []string{`0/Payload/{type}="*ContentData"`, `0/Payload/Data="a"`}, nil},
		{`graph (load(Fs="test", Glob="` + testEmbedShortGlob + `"))`, []string{`0/Payload/{type}="*ContentData"`, `0/Payload/Data="a"`}, nil},
		{`graph (anna -> regexp(Target="Content.Name",Expr="be"))`, []string{`0/Payload/{type}="*ContentData"`, `0/Payload/Name="Annath"`}, nil},
		// Errors
		{`graph (load(Fs="-"))`, nil, fmt.Errorf("no filesystem")},
	}
	for i, v := range table {
		p, err := pipeline.Compile(v.pipeline)
		if err != nil {
			t.Fatalf("TestPipeline %v compile err %v", i, err)
		}
		output, haveErr := pipeline.Run(p, nil, nil)
		var cmpErr error
		if output != nil {
			cmpErr = jacl.Run(output.Pins, v.cmp...)
		}
		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestPipeline %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestPipeline %v has no error but exptected %v", i, v.wantErr)
		} else if cmpErr != nil {
			t.Fatalf("TestPipeline %v comparison error: %v", i, cmpErr)
		}
	}
}

// ---------------------------------------------------------
// SUPPORT

// contentSrcNode can be used to generate ContentData.
type contentSrcNode struct {
	data []*pipeline.ContentData
}

func (n *contentSrcNode) Run(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	output := pipeline.RunOutput{}
	for _, cd := range n.data {
		output.Pins = append(output.Pins, pipeline.Pin{Payload: cd})
	}
	return &output, nil
}

// ---------------------------------------------------------
// LIFECYCLE

func setupTests() {
	pipeline.RegisterFs("test", testdataFs)

	pipeline.RegisterNode("anna", func() pipeline.Node {
		n := &contentSrcNode{}
		n.data = append(n.data, &pipeline.ContentData{Name: "Annabeth", Data: "born 2002 of fair skin and stout heart"})
		return n
	})
}

//go:embed test_data/*
var testdataFs embed.FS

// Globs
var (
	testDataDomainGlob = filepath.Join(".", "test_data", "domain_*")
	testDataShortGlob  = filepath.Join(".", "test_data", "short_*")
	testEmbedShortGlob = "test_data/short_*"
)
