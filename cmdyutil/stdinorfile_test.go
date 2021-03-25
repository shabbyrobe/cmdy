package cmdyutil

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/shabbyrobe/cmdy"
)

type testContext struct {
	context.Context
	commandRef cmdy.CommandRef
	stdin      io.Reader
	stdout     io.Writer
	stderr     io.Writer
}

func (t *testContext) RawArgs() []string                    { return nil }
func (t *testContext) Stdin() io.Reader                     { return t.stdin }
func (t *testContext) Stdout() io.Writer                    { return t.stdout }
func (t *testContext) Stderr() io.Writer                    { return t.stderr }
func (t *testContext) Runner() *cmdy.Runner                 { return nil }
func (t *testContext) Stack() cmdy.CommandPath              { return nil }
func (t *testContext) Current() cmdy.CommandRef             { return t.commandRef }
func (t *testContext) Push(name string, cmd cmdy.Command)   {}
func (t *testContext) Pop() (name string, cmd cmdy.Command) { return "", nil }

func ctxWithStdin(data []byte) *testContext {
	var stdin bytes.Buffer
	stdin.Write(data)
	var ctx = testContext{stdin: &stdin}
	return &ctx
}

func assertFileContents(t *testing.T, f io.Reader, exp []byte) {
	t.Helper()
	v, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, exp) {
		t.Fatal()
	}
}

func mustOpenStdinOrFile(t *testing.T, ctx cmdy.Context, fname string, flag StdinFlag) io.ReadCloser {
	t.Helper()
	f, err := OpenStdinOrFile(ctx, fname, flag)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestOpenStdinOrFile(t *testing.T) {
	ctx := ctxWithStdin([]byte("data"))
	f := mustOpenStdinOrFile(t, ctx, "", 0)
	assertFileContents(t, f, []byte("data"))
}

func TestOpenStdinOrFileWithHyphenNameTriesToUseFileWhenFlagNotSet(t *testing.T) {
	ctx := ctxWithStdin([]byte("data"))
	_, err := OpenStdinOrFile(ctx, "-", 0)
	if err != errStdinOrFileBoth {
		t.Fatal()
	}
}

func TestOpenStdinOrFileWithHyphenNameUsesInputWhenFlagSet(t *testing.T) {
	ctx := ctxWithStdin([]byte("data"))
	f := mustOpenStdinOrFile(t, ctx, "-", HyphenStdin)
	assertFileContents(t, f, []byte("data"))
}
