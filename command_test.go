package cmdy

import (
	"context"
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
