package cmdy

import (
	"context"
	"strings"
	"testing"

	"github.com/shabbyrobe/cmdy/args"
	"github.com/shabbyrobe/golib/assert"
)

func TestRun(t *testing.T) {
	t.Run("args-no-flags", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo string
		as := args.NewArgSet()
		as.String(&foo, "foo", "usage!")

		tc := &testCmd{args: as}
		bld := func() (Command, Init) { return tc, nil }
		rn := newTestRunner()

		// missing argument should produce a usage error:
		err := rn.Run(context.Background(), "test", nil, bld)
		tt.MustEqual(ExitUsage, errCode(err))

		// extra argument should produce a usage error:
		err = rn.Run(context.Background(), "test", []string{"a", "b"}, bld)
		tt.MustEqual(ExitUsage, errCode(err))

		// unexpected flag should produce a usage error:
		err = rn.Run(context.Background(), "test", []string{"--quack"}, bld)
		tt.MustEqual(ExitUsage, errCode(err), "%v", err)

		// --help should produce a usage error even with an argument:
		err = rn.Run(context.Background(), "test", []string{"--help", "arg"}, bld)
		tt.MustEqual(ExitUsage, errCode(err), "%v", err)

		// double-dash should work
		tt.MustOK(rn.Run(context.Background(), "test", []string{"--", "--quack"}, bld))

		// should work:
		tt.MustOK(rn.Run(context.Background(), "test", []string{"a"}, bld))
	})

	t.Run("flags-no-args", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo string
		fs := NewFlagSet()
		fs.StringVar(&foo, "foo", "", "usage!")

		tc := &testCmd{flags: fs}
		bld := func() (Command, Init) { return tc, nil }
		rn := newTestRunner()

		// no arguments, no flags, no worries:
		tt.MustOK(rn.Run(context.Background(), "test", nil, bld))

		// any argument should produce a usage error:
		err := rn.Run(context.Background(), "test", []string{"a"}, bld)
		tt.MustEqual(ExitUsage, errCode(err), "%v", err)

		// unexpected flag should produce a usage error:
		err = rn.Run(context.Background(), "test", []string{"--quack"}, bld)
		tt.MustEqual(ExitUsage, errCode(err), "%v", err)

		// double-dash followed by arg should produce a usage error:
		err = rn.Run(context.Background(), "test", []string{"--", "--quack"}, bld)
		tt.MustEqual(ExitUsage, errCode(err))

		// double-dash preceded by flag should work:
		tt.MustOK(rn.Run(context.Background(), "test", []string{"--foo", "bar", "--"}, bld))

		// double-dash only should work:
		tt.MustOK(rn.Run(context.Background(), "test", []string{"--"}, bld))

		// missing flag value should fail:
		err = rn.Run(context.Background(), "test", []string{"--quack"}, bld)
		tt.MustEqual(ExitUsage, errCode(err))

		// flag should work:
		tt.MustOK(rn.Run(context.Background(), "test", []string{"--foo", "bar"}, bld))
	})
}

func TestRunFullHelp(t *testing.T) {
	tt := assert.WrapTB(t)

	// need to configure one argument so we can trip a usage error by omitting it.
	var foo string
	as := args.NewArgSet()
	as.String(&foo, "foo", "usage!")

	tc := &testCmd{args: as, usage: "foo{{if ShowFullHelp}}bar{{end}}"}
	bld := func() (Command, Init) { return tc, nil }
	rn := newTestRunner()

	{ // --help should ShowFullHelp
		err := rn.Run(context.Background(), "test", []string{"--help"}, bld)
		tt.MustEqual(ExitUsage, errCode(err), "%v", err)
		msg, code := FormatError(err)
		tt.MustAssert(strings.HasPrefix(msg, "foobar\n"))
		tt.MustEqual(ExitUsage, code)
	}

	{ // UsageError should not ShowFullHelp
		err := rn.Run(context.Background(), "test", []string{}, bld)
		tt.MustEqual(ExitUsage, errCode(err), "%v", err)
		msg, code := FormatError(err)
		tt.MustAssert(strings.HasPrefix(msg, "foo\n"))
		tt.MustEqual(ExitUsage, code)
	}
}

func errCode(err error) int {
	switch err := err.(type) {
	case interface{ Code() int }:
		return err.Code()
	case nil:
		return 0
	default:
		return 1
	}
}
