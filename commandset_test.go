package cmdy

import (
	"context"
	"errors"
	"testing"

	"github.com/shabbyrobe/cmdy/args"
	"github.com/shabbyrobe/golib/assert"
)

func TestCommandSet(t *testing.T) {
	tt := assert.WrapTB(t)

	var foo, bar = errors.New("foo"), errors.New("bar")

	bldr := func() (Command, error) {
		return NewCommandSet("set", Builders{
			"foo": testCmdRunBuilder(func(c Context, i Input) error {
				return foo
			}),
			"bar": testCmdRunBuilder(func(c Context, i Input) error {
				return bar
			}),
		}), nil
	}

	tt.MustEqual(foo, Run(context.Background(), []string{"foo"}, bldr))
	tt.MustEqual(bar, Run(context.Background(), []string{"bar"}, bldr))
}

func TestCommandSet_SubcommandArgs(t *testing.T) {
	tt := assert.WrapTB(t)

	var p string
	bldr := func() (Command, error) {
		return NewCommandSet("set", Builders{
			"foo": func() (Command, error) {
				p = ""
				as := args.NewArgSet()
				as.String(&p, "pants", "Usage...")
				return &testCmd{args: as}, nil
			},
		}), nil
	}

	tt.MustOK(Run(context.Background(), []string{"foo", "yep"}, bldr))
	tt.MustEqual("yep", p)

	err := Run(context.Background(), []string{"foo", "yep", "nup"}, bldr)
	tt.MustAssert(err != nil) // FIXME: check error

	err = Run(context.Background(), []string{"foo"}, bldr)
	tt.MustAssert(err != nil) // FIXME: check error
}

func TestCommandSet_SubcommandFlags(t *testing.T) {
	tt := assert.WrapTB(t)

	var p string
	bldr := func() (Command, error) {
		return NewCommandSet("set", Builders{
			"foo": func() (Command, error) {
				fs := NewFlagSet()
				fs.StringVar(&p, "pants", "", "Usage...")
				return &testCmd{flags: fs}, nil
			},
		}), nil
	}

	tt.MustOK(Run(context.Background(), []string{"foo", "-pants", "yep"}, bldr))
	tt.MustEqual("yep", p)

	err := Run(context.Background(), []string{"foo", "-pants"}, bldr)
	tt.MustAssert(err != nil) // FIXME: check error
}
