package cmdy

import (
	"context"
	"strings"
	"testing"

	"github.com/shabbyrobe/cmdy/arg"
	"github.com/shabbyrobe/cmdy/internal/assert"
)

func TestCommand_FlagsArgs(t *testing.T) {
	tt := assert.WrapTB(t)

	var foo, bar string
	fs := NewFlagSet()
	fs.StringVar(&foo, "foo", "", "usage...")
	as := arg.NewArgSet()
	as.String(&bar, "bar", "usage...")
	c := &testCmd{
		flags: fs,
		args:  as,
	}

	tt.MustOK(Run(context.Background(), []string{"-foo", "foo", "bar"}, testBuilder(c)))
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

	usage := Run(context.Background(), []string{"-help"}, testBuilder(c))
	txt, code := FormatError(usage)
	tt.MustEqual(127, code)

	// Warning: brittle test
	tt.MustEqual("synopsis\n\nUsage: cmdy.test [options] \n\nTest", txt)
}

type testUsageVarsCmd struct {
	Stuff string
	Arg   string
	Flag  string
}

func (t *testUsageVarsCmd) Run(c Context) error { return nil }

func (t *testUsageVarsCmd) Help() Help {
	return Help{
		Synopsis: "usage vars cmd",
		Usage:    "{{.Stuff}} {{.Flag}} {{.Arg}}",
	}
}

func (t *testUsageVarsCmd) Configure(flags *FlagSet, args *arg.ArgSet) {
	flags.StringVar(&t.Flag, "flag", "", "Var")
	args.String(&t.Arg, "arg", "Var")
}

func TestUsageVars(t *testing.T) {
	assertUsage := func(tt assert.T, code int, in []string, out string) {
		c := &testUsageVarsCmd{Stuff: "stuff"}
		usage := Run(context.Background(), in, testBuilder(c))
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
