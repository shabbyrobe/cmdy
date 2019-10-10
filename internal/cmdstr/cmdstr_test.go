package cmdstr

import (
	"reflect"
	"testing"
)

func assertYep(t *testing.T, in string, out ...string) {
	t.Helper()
	result, err := SplitString(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, out) {
		t.Fatal("result", result, "!=", out)
	}
}

func TestCmdStr(t *testing.T) {
	assertYep(t, `foo`, `foo`)

	// intra-unqouted quotes are passed through:
	assertYep(t, `foo"bar"baz`, `foo"bar"baz`)
	assertYep(t, `foo'bar'baz`, `foo'bar'baz`)

	// intra-unqouted quotes don't have to be closed:
	assertYep(t, `foo'barbaz`, `foo'barbaz`)
}
