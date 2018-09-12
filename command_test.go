package cmdy

import (
	"context"
	"strings"
	"testing"

	"github.com/shabbyrobe/cmdy/args"
	"github.com/shabbyrobe/golib/assert"
)

func TestCommand_FlagsArgs(t *testing.T) {
	tt := assert.WrapTB(t)

	var foo, bar string
	fs := NewFlagSet()
	fs.StringVar(&foo, "foo", "", "usage...")
	as := args.NewArgSet()
	as.String(&bar, "bar", "usage...")
	c := &testCmd{
		flags: fs,
		args:  as,
	}

	tt.MustOK(Run(context.Background(), []string{"-foo", "foo", "bar"}, builder(c)))
	tt.MustEqual("foo", foo)
	tt.MustEqual("bar", bar)
}

func TestCommand_TemplateDefault(t *testing.T) {
	tt := assert.WrapTB(t)
	_ = tt

	c := &testCmd{
		usage: DefaultUsage + "\n" +
			"Test",
		synopsis: "synopsis",
	}

	usage := Run(context.Background(), []string{"-help"}, builder(c))
	txt, code := FormatError(usage)
	tt.MustEqual(127, code)

	// Warning: brittle test
	tt.MustEqual("synopsis\n\nUsage: cmdy.test [options]\n\nTest", txt)
}

type testUsageVarsCmd struct {
	Stuff string
	Arg   string
	Flag  string
}

func (t *testUsageVarsCmd) Run(c Context) error { return nil }
func (t *testUsageVarsCmd) Synopsis() string    { return "usage vars cmd" }
func (t *testUsageVarsCmd) Usage() string       { return "{{.Stuff}} {{.Flag}} {{.Arg}}" }

func (t *testUsageVarsCmd) Flags() *FlagSet {
	fs := NewFlagSet()
	fs.StringVar(&t.Flag, "flag", "", "Var")
	return fs
}

func (t *testUsageVarsCmd) Args() *args.ArgSet {
	as := args.NewArgSet()
	as.String(&t.Arg, "arg", "Var")
	return as
}

func TestUsageVars(t *testing.T) {
	assertUsage := func(tt assert.T, code int, in []string, out string) {
		c := &testUsageVarsCmd{Stuff: "stuff"}
		usage := Run(context.Background(), in, builder(c))
		txt, ecode := FormatError(usage)
		tt.MustEqual(code, ecode)
		tt.MustEqual(out, strings.Split(txt, "\n")[0])
	}

	t.Run("", func(t *testing.T) {
		// the args and flags will be parsed if there is an extra arg:
		in := []string{"-flag=foo", "bar", "extra"}
		out := "stuff foo bar"
		assertUsage(assert.WrapTB(t), 127, in, out)
	})

	t.Run("", func(t *testing.T) {
		// this won't get a chance to parse the flag:
		in := []string{"-wat", "-flag=foo", "bar"}
		out := "stuff"
		assertUsage(assert.WrapTB(t), 127, in, out)
	})

	t.Run("", func(t *testing.T) {
		in := []string{}
		out := "stuff"
		assertUsage(assert.WrapTB(t), 127, in, out)
	})
}
