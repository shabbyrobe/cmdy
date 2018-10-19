package cmdy

import (
	"bytes"

	"github.com/shabbyrobe/cmdy/args"
)

type testCmd struct {
	synopsis string
	usage    string
	flags    *FlagSet
	args     *args.ArgSet
	run      func(c Context) error
}

func testCmdRunBuilder(r func(c Context) error) func() (Command, Init) {
	return func() (Command, Init) { return &testCmd{run: r}, nil }
}

func testCmdRunInitBuilder(r func(c Context) error, init Init) func() (Command, Init) {
	return func() (Command, Init) { return &testCmd{run: r}, init }
}

func testBuilder(c Command) func() (Command, Init) {
	return func() (Command, Init) { return c, nil }
}

func testInitBuilder(c Command, init Init) func() (Command, Init) {
	return func() (Command, Init) { return c, init }
}

func (t *testCmd) Synopsis() string   { return t.synopsis }
func (t *testCmd) Usage() string      { return t.usage }
func (t *testCmd) Flags() *FlagSet    { return t.flags }
func (t *testCmd) Args() *args.ArgSet { return t.args }

func (t *testCmd) Run(c Context) error {
	if t.run != nil {
		return t.run(c)
	}
	return nil
}

type testRunner struct {
	stdin  bytes.Buffer
	stdout bytes.Buffer
	stderr bytes.Buffer
	*Runner
}

func newTestRunner() *testRunner {
	tr := &testRunner{}
	tr.Runner = &Runner{
		Stdin:  &tr.stdin,
		Stdout: &tr.stdout,
		Stderr: &tr.stderr,
	}
	return tr
}
