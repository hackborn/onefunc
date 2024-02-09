package cfg

import (
	"embed"
	"os"
	"testing"

	"github.com/hackborn/onefunc/jacl"
)

func TestMain(m *testing.M) {
	setupTests()
	code := m.Run()
	os.Exit(code)
}

// ---------------------------------------------------------
// TEST-STRING
func TestString(t *testing.T) {
	table := []struct {
		opts    []Option
		subset  []string
		path    string
		want    string
		wantErr error
	}{
		{nil, nil, "", missingString, nil},
		{[]Option{WithFS(dataFs, "test_data/settings_a.json")}, nil, "a", "anna", nil},
		{[]Option{WithFS(dataFs, "test_data/settings_b.json")}, nil, "a", "ava", nil},
		{[]Option{WithFS(dataFs, "test_data/settings_[{a-b}].json")}, nil, "a", "ava", nil},
		{[]Option{WithEnv(EnvPattern("CFG_TESTDATA_*"))}, nil, "CFG_TESTDATA_A", "ant", nil},
		{[]Option{WithEnv(EnvPrefix("CFG_TESTDATA_"))}, nil, "A", "ant", nil},
	}
	for i, v := range table {
		s, haveErr := NewSettings(v.opts...)
		for _, path := range v.subset {
			s = s.Subset(path)
		}
		have := s.MustString(v.path, missingString)
		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestString %v %v", i, err.Error())
		} else if have != v.want {
			t.Fatalf("TestString %v has \"%v\" but wants \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// LIFECYCLE

func setupTests() {
	os.Setenv("CFG_TESTDATA_A", "ant")
	os.Setenv("CFG_TESTDATA_B", "bear")
	os.Setenv("CFG_TESTDATA_C", "cat")
}

// ---------------------------------------------------------
// CONST and VAR

const (
	missingString = "~missing~"
)

//go:embed test_data/*
var dataFs embed.FS
