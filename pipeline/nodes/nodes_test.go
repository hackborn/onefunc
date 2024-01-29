package nodes

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hackborn/onefunc/jacl"
	"github.com/hackborn/onefunc/pipeline"
)

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
		output, haveErr := pipeline.Run(p, nil)
		fmt.Println("output", output)
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
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Name=Company`, `1/Payload/Name=Filing`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Fields/0/Name=Id`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Fields/0/Tag="doc:''id, key''"`}, nil},
		{`graph (load(Glob="` + testDataDomainGlob + `") -> struct(Tag=doc))`, []string{`0/Payload/Fields/0/Tag="id, key"`}, nil},
	}
	for i, v := range table {
		p, err := pipeline.Compile(v.pipeline)
		if err != nil {
			t.Fatalf("TestPipeline %v compile err %v", i, err)
		}
		output, haveErr := pipeline.Run(p, nil)
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

// Globs
var (
	testDataDomainGlob = filepath.Join(".", "test_data", "domain_*")
	testDataShortGlob  = filepath.Join(".", "test_data", "short_*")
)
