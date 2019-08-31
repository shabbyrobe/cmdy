package cmdy

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/shabbyrobe/cmdy/arg"
	"github.com/shabbyrobe/cmdy/internal/assert"
)

func TestGroup(t *testing.T) {
	tt := assert.WrapTB(t)

	var foo, bar = errors.New("foo"), errors.New("bar")

	bldr := func() (Command, Init) {
		return NewGroup("set", Builders{
			"foo": testCmdRunBuilder(func(c Context) error {
				return foo
			}),
			"bar": testCmdRunBuilder(func(c Context) error {
				return bar
			}),
		}), nil
	}

	tt.MustEqual(foo, Run(context.Background(), []string{"foo"}, bldr))
	tt.MustEqual(bar, Run(context.Background(), []string{"bar"}, bldr))
}

func TestGroup_SubcommandArgs(t *testing.T) {
	tt := assert.WrapTB(t)

	var p string
	bldr := func() (Command, Init) {
		return NewGroup("set", Builders{
			"foo": func() (Command, Init) {
				p = ""
				as := arg.NewArgSet()
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

func TestGroup_SubcommandFlags(t *testing.T) {
	tt := assert.WrapTB(t)

	var p string
	bldr := func() (Command, Init) {
		return NewGroup("set", Builders{
			"foo": func() (Command, Init) {
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

func TestGroup_Unknown(t *testing.T) {
	tt := assert.WrapTB(t)

	var foo = errors.New("foo")

	bldr := func() (Command, Init) {
		return NewGroup("set", Builders{},
			GroupUnknown(func() (Command, Init) {
				return nil, func() error { return foo }
			}),
		), nil
	}

	tt.MustEqual(foo, Run(context.Background(), []string{"foo"}, bldr))
}

func TestGroup_Hide(t *testing.T) {
	tt := assert.WrapTB(t)
	_ = tt

	grp := NewGroup("set",
		Builders{
			"4GKwDcbp": func() (Command, Init) { return &testCmd{synopsis: "4GKwDcbp"}, nil },
			"9rdjKX3j": func() (Command, Init) { return &testCmd{synopsis: "9rdjKX3j"}, nil },
			"GM68tb0F": func() (Command, Init) { return &testCmd{synopsis: "GM68tb0F"}, nil },
			"OZJpKePU": func() (Command, Init) { return &testCmd{synopsis: "OZJpKePU"}, nil },
		},
		GroupHide("9rdjKX3j", "OZJpKePU"),
	)

	out := grp.Usage()
	tt.MustAssert(!strings.Contains(out, "OZJpKePU"))
	tt.MustAssert(!strings.Contains(out, "9rdjKX3j"))
	tt.MustAssert(strings.Contains(out, "GM68tb0F"))
	tt.MustAssert(strings.Contains(out, "4GKwDcbp"))
}
