package cmdytest

import (
	"fmt"
	"testing"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
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
	validExamples := cmdy.Examples{
		{Desc: "yep", Command: ""},
		{Desc: "yep", Command: "-bool"},
		{Desc: "nup", Command: "-blool", Code: cmdy.ExitUsage},
		{Desc: "nup", Command: "", Code: 100, TestMode: cmdy.ExampleRun},
		{Desc: "nup", Command: "-blool", Code: -1}, // -1 allows any exit code
	}

	for _, ex := range validExamples {
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

	invalidExamples := cmdy.Examples{
		{Command: "-blool", Code: 0},
		{Command: "-bool", Code: -1},
		{Command: "", Code: -1},
		{Command: "-blool", Code: 63},
		{Command: "", Code: 63},
	}

	for _, ex := range invalidExamples {
		t.Run(slugify(ex.Desc), func(t *testing.T) {
			tester := ExampleTester{
				TestName: "yep",
				Builder:  func() cmdy.Command { return &testCommand{} },
			}
			if err := tester.RunExample(ex); err == nil {
				t.Fatalf("expected error for command %q", ex.Command)
			}
		})
	}
}
