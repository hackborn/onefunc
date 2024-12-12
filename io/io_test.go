package io

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// ---------------------------------------------------------
// TEST-EXPAND-PATH
func TestExpandPath(t *testing.T) {
	f := func(path, want string) {
		t.Helper()

		have, haveErr := ExpandPath(path)
		if haveErr != nil {
			t.Fatalf("TestExpandPath has err %v", haveErr)
		} else if have != want {
			t.Fatalf("TestExpandPath wants %v but has %v", want, have)
		}
	}

	f("termy", "termy")
}
