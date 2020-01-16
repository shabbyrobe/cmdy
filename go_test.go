package cmdy_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const expectedMod = `
module github.com/shabbyrobe/cmdy

go 1.12
`

func TestNoDeps(t *testing.T) {
	if os.Getenv("CMDY_SKIP_MOD") != "" {
		// Use this to avoid this check if you need to use spew.Dump in tests:
		t.Skip()
	}

	{
		bts, err := ioutil.ReadFile("go.mod")
		if os.IsNotExist(err) {
			t.Skip()
		}
		if !bytes.Equal([]byte(strings.TrimSpace(expectedMod)), bytes.TrimSpace(bts)) {
			t.Fatal("go.mod contains unexpected content")
		}
	}

	{
		bts, err := ioutil.ReadFile("go.sum")
		if os.IsNotExist(err) {
			t.Skip()
		}
		if len(bts) != 0 {
			t.Fatal("go.sum contains unexpected content")
		}
	}
}
