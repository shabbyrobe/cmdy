package cmdytest

import (
	"fmt"
	"testing"

	"github.com/ArtProcessors/cmdy"
	"github.com/ArtProcessors/cmdy/arg"
)

type testCommand struct {
	boolVal bool
}

func (cmd *testCommand) Help() cmdy.Help { return cmdy.Help{} }

func (cmd *testCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	flags.BoolVar(&cmd.boolVal, "bool", false, "test")
}

func (cmd *testCommand) Run(ctx cmdy.Context) error {
	return cmdy.ErrWithCode(100, fmt.Errorf("boom"))
}

func TestExampleTester(t *testing.T) {
	examples := cmdy.Examples{
		{Desc: "yep", Command: ""},
		{Desc: "yep", Command: "-bool"},
		{Desc: "nup", Command: "-blool", Code: cmdy.ExitUsage},
		{Desc: "nup", Command: "", Code: 100, TestMode: cmdy.ExampleRun},
	}

	for _, ex := range examples {
		t.Run(slugify(ex.Desc), func(t *testing.T) {
			tester := ExampleTester{
				TestName: "yep",
				Builder:  func() cmdy.Command { return &testCommand{} },
			}
			if err := tester.RunExample(ex); err != nil {
				t.Fatal(err)
			}
		})
	}
}
