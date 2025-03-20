package sys

import (
	"fmt"
	"os"
	"testing"
)

// This code exists mostly to help development, much of
// what it accesses is device-dependent.

func TestMain(m *testing.M) {
	Set(SetAppName("test"))

	code := m.Run()
	os.Exit(code)
}

// ---------------------------------------------------------
// TEST-GET
func TestGet(t *testing.T) {
	f := func(key string, want string) {
		t.Helper()

		info, err := Get(key)
		if err != nil {
			panic(err)
		}
		fmt.Println(info)
		// t.Fatalf("no")
	}

	f(HardwareModel, "not a real test")
}
