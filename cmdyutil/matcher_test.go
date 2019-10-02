package cmdyutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	"github.com/shabbyrobe/cmdy/internal/assert"
)

type testCommand struct {
	run func(ctx cmdy.Context) error
}

func (cmd *testCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {}

func (cmd *testCommand) Synopsis() string           { return "" }
func (cmd *testCommand) Run(ctx cmdy.Context) error { return cmd.run(ctx) }

func strOutputBuilder(out string) func() cmdy.Command {
	return func() cmdy.Command {
		return &testCommand{run: func(ctx cmdy.Context) error {
			fmt.Println(out)
			return nil
		}}
	}
}

var (
	newFooCommand  = strOutputBuilder("foo")
	newFoodCommand = strOutputBuilder("food")
	newBaleCommand = strOutputBuilder("bale")
	newBarkCommand = strOutputBuilder("bark")
)

func ExampleGroup_PrefixMatcher() {
	builders := cmdy.Builders{
		"foo":  newFooCommand,
		"food": newFoodCommand,
		"bale": newBaleCommand,
		"bark": newBarkCommand,
	}

	bldr := func() cmdy.Command {
		return cmdy.NewGroup("group", builders, GroupPrefixMatcher(2))
	}

	cmdy.Run(context.Background(), []string{"foo"}, bldr)
	cmdy.Run(context.Background(), []string{"food"}, bldr)
	cmdy.Run(context.Background(), []string{"bar"}, bldr)
	cmdy.Run(context.Background(), []string{"bal"}, bldr)

	// Output:
	// foo
	// food
	// bark
	// bale
}

func TestMatcher(t *testing.T) {
	type tc struct {
		min      int
		in       string
		expected string
		options  []string
	}

	for _, c := range []tc{
		{min: 3, in: "foo", expected: "foo", options: []string{"foo", "food"}},
		{min: 3, in: "food", expected: "food", options: []string{"foo", "food"}},
		{min: 3, in: "food", expected: "foods", options: []string{"foo", "foods"}},
		{min: 4, in: "food", expected: "foods", options: []string{"foo", "foods"}},
		{min: 5, in: "food", expected: "", options: []string{"foo", "foods"}},
		{min: 3, in: "food", expected: "", options: []string{}},
	} {
		t.Run("", func(t *testing.T) {
			tt := assert.WrapTB(t)
			bldrs := cmdy.Builders{}
			for _, name := range c.options {
				bldrs[name] = func() cmdy.Command { return &testCommand{} }
			}
			grp := cmdy.NewGroup("g", bldrs, GroupPrefixMatcher(c.min))
			_, name, err := grp.Builder(c.in)
			tt.MustOK(err)
			tt.MustEqual(c.expected, name)
		})
	}
}
