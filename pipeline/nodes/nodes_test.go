package nodes

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hackborn/onefunc/jacl"
	"github.com/hackborn/onefunc/pipeline"
)

// ---------------------------------------------------------
// TEST-LOAD-FILE
func TestLoadFile(t *testing.T) {
	table := []struct {
		pipeline string
		cmp      []string
		wantErr  error
	}{
		{`graph (loadfile(Glob="` + testDataShortGlob + `"))`, []string{`0/Payload/Data=a`, `1/Payload/Data=b`, `2/Payload/Data=c`}, nil},
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
			t.Fatalf("TestLoadFile %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestLoadFile %v has no error but exptected %v", i, v.wantErr)
		} else if cmpErr != nil {
			t.Fatalf("TestLoadFile %v comparison error: %v", i, cmpErr)
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
		{`graph (loadfile(Glob="` + testDataDomainGlob + `") -> struct)`, []string{`0/Payload/Name=Company`, `1/Payload/Name=Filing`}, nil},
	}
	for i, v := range table {
		p, err := pipeline.Compile(v.pipeline)
		if err != nil {
			t.Fatalf("TestLoadFile %v compile err %v", i, err)
		}
		output, haveErr := pipeline.Run(p, nil)
		var cmpErr error
		if output != nil {
			cmpErr = jacl.Run(output.Pins, v.cmp...)
		}
		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestLoadFile %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestLoadFile %v has no error but exptected %v", i, v.wantErr)
		} else if cmpErr != nil {
			t.Fatalf("TestLoadFile %v comparison error: %v", i, cmpErr)
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
