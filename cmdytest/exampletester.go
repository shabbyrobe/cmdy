package cmdytest

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/internal/cmdstr"
)

type ExampleTester struct {
	TestName string
	Builder  cmdy.Builder
	Setup    func(c cmdy.Command)
	Cleanup  func(c cmdy.Command)
}

func (e *ExampleTester) wrapBuilder(example cmdy.Example) cmdy.Builder {
	return func() cmdy.Command {
		cmd := e.Builder()
		return &wrappedCommand{cmd, e, example}
	}
}

func (e *ExampleTester) Examples() cmdy.Examples {
	return e.Builder().Help().Examples
}

func (e *ExampleTester) TestExamples(t *testing.T) {
	t.Helper()

	for _, example := range e.Examples() {
		slug := slugify(example.Desc)
		name := fmt.Sprintf("%s/%s", e.TestName, slug)
		t.Run(name, func(t *testing.T) {
			if err := e.RunExample(example); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func (e *ExampleTester) RunExample(example cmdy.Example) error {
	args, err := cmdstr.ParseString(example.Command, "")
	if err != nil {
		return err
	}

	ctx := context.Background()
	runner := cmdy.NewBufferedRunner()
	runner.StdinBuffer.WriteString(example.Input)
	builder := e.wrapBuilder(example)
	runErr := runner.Run(ctx, "tester", args, builder)
	code := cmdy.ErrCode(runErr)
	if code != example.Code {
		return fmt.Errorf("unexpected code %d, expected %d: %w", code, example.Code, runErr)
	}

	// FIXME: output

	return nil
}

type wrappedCommand struct {
	cmdy.Command
	tester  *ExampleTester
	example cmdy.Example
}

func (w *wrappedCommand) Run(ctx cmdy.Context) error {
	if w.example.TestMode == cmdy.ExampleParseOnly {
		return nil
	}

	if w.tester.Setup != nil {
		w.tester.Setup(w.Command)
	}
	ret := w.Command.Run(ctx)
	if w.tester.Cleanup != nil {
		w.tester.Cleanup(w.Command)
	}
	return ret
}

var spacePtn = regexp.MustCompile(`[\s\-]+`)
var unslugPtn = regexp.MustCompile(`[^\s\-\pL\pN]`)

func slugify(v string) string {
	out := v
	out = unslugPtn.ReplaceAllString(out, "")
	out = spacePtn.ReplaceAllString(out, "-")
	out = strings.Trim(out, "-")
	out = strings.ToLower(out)
	return fmt.Sprintf("%s-%x", out, fnv1(v))
}

func fnv1(s string) uint32 {
	var x uint32
	slen := len(s)
	for i := 0; i < slen; i++ {
		x = x*16777619 ^ uint32(s[i])
	}
	return x
}
