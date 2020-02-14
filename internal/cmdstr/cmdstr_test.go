package cmdstr

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

func yep(t *testing.T, in string, out ...string) {
	t.Helper()
	result, err := ParseString(in, "")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, out) {
		t.Fatal("result", result, "!=", out)
	}
}

func nup(t *testing.T, in string, expected error) {
	t.Helper()
	_, err := ParseString(in, "")
	if err == nil {
		t.Fatal()
	}
	if !errors.Is(err, expected) {
		t.Fatal(err, "!=", expected)
	}
}

func TestSplitString(t *testing.T) {
	yep(t, ``)
	yep(t, `  `)
	yep(t, `""`, "")
	yep(t, ` "" `, "")
	yep(t, `''`, "")
	yep(t, ` '' `, "")

	yep(t, `foo`, `foo`)
	yep(t, `foo bar`, `foo`, `bar`)
	yep(t, `foo   bar`, `foo`, `bar`)
	yep(t, `  foo   bar  `, `foo`, `bar`)

	yep(t, `"foo"`, `foo`)
	yep(t, `" foo "`, ` foo `)
	yep(t, `"foo" "bar"`, `foo`, `bar`)
	yep(t, `"foo"   "bar"`, `foo`, `bar`)
	yep(t, `  "foo"   "bar"  `, `foo`, `bar`)

	yep(t, `'foo'`, `foo`)
	yep(t, `' foo '`, ` foo `)
	yep(t, `'foo' 'bar'`, `foo`, `bar`)
	yep(t, `'foo'   'bar'`, `foo`, `bar`)
	yep(t, `  'foo'   'bar'  `, `foo`, `bar`)

	// Escapes in bare arguments:
	yep(t, `foo\nbar`, "foo\nbar")

	// XXX: deviation from shlex. shlex errors with 'unclosed quotation, but we just ignore the escape
	yep(t, `foo\'bar`, `foo\'bar`)

	// Newline escapes in double quotes:
	yep(t, `"\`+"\n"+`foo"`, "foo")
	yep(t, `"f\`+"\n"+`oo"`, "foo")
	yep(t, `"fo\`+"\n"+`o"`, "foo")
	yep(t, `"foo\`+"\n"+`"`, "foo")

	yep(t, `'\nfoo'`, "\nfoo")
	yep(t, `'f\noo'`, "f\noo")
	yep(t, `'fo\no'`, "fo\no")
	yep(t, `'foo\n'`, "foo\n")

	yep(t, `"\\foo"`, "\\foo")
	yep(t, `"f\\oo"`, "f\\oo")
	yep(t, `"fo\\o"`, "fo\\o")
	yep(t, `"foo\\"`, "foo\\")

	// Unknown escape is literal:
	yep(t, `"\s"`, `\s`)

	// intra-unqouted quotes are passed through:
	yep(t, `foo"bar"baz`, `foo"bar"baz`)
	yep(t, `foo'bar'baz`, `foo'bar'baz`)

	// intra-unqouted quotes have to be closed:
	nup(t, `foo'barbaz`, ErrIncompleteCommand)
}

func TestParseIncompleteErr(t *testing.T) {
	tt := assert.WrapTB(t)
	in := "foo'abc\\"
	_, _, err := Parse([]byte(in), "")
	tt.Assert(strings.Contains(err.Error(), "arg > squot > esc"))
}

func TestParseEndSet(t *testing.T) {
	tt := assert.WrapTB(t)
	in := "foo bar [yep"
	out, n, err := Parse([]byte(in), "[")
	tt.MustOK(err)
	tt.MustEqual("[yep", in[n:])
	tt.MustEqual([]string{"foo", "bar"}, out)
}

var BenchSplitResult []string

func BenchmarkSplitString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BenchSplitResult, _ = ParseString(`foo bar 'baz' "qux" \n "yep" `, "")
	}
}
