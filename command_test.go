package cmdy

import (
	"context"
	"testing"

	"github.com/ArtProcessors/cmdy/arg"
	"github.com/ArtProcessors/cmdy/internal/assert"
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

func TestCommand_Usage(t *testing.T) {
	tt := assert.WrapTB(t)
	_ = tt

	c := &testCmd{
		usage:    "Test",
		synopsis: "synopsis",
	}

	runner := NewBufferedRunner()
	usage := runner.Run(context.Background(), "cmdy", []string{"-help"}, testBuilder(c))
	txt, code := FormatError(usage)
	tt.MustEqual(0, code) // -help should return 0 exit status

	// Warning: brittle test
	expected := "synopsis\n\nUsage: cmdy [options] \n\nTest"
	tt.MustEqual(expected, txt)
}
