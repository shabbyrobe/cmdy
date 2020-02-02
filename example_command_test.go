package cmdy_test

import (
	"context"
	"fmt"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
)

type myCommand struct {
	who string
}

func myCommandBuilder() cmdy.Command { return &myCommand{} }

func (cmd *myCommand) Help() cmdy.Help { return cmdy.Synopsis("Hello world") }

func (cmd *myCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	args.StringOptional(&cmd.who, "who", "World", "Who to hello!")
}

func (cmd *myCommand) Run(ctx cmdy.Context) error {
	fmt.Fprintf(ctx.Stdout(), "Hello %s!\n", cmd.who)
	return nil
}

func ExampleCommand() {
	runner := cmdy.NewBufferedRunner()
	prog := "example"
	ctx := context.Background()

	if err := runner.Run(ctx, prog, nil, myCommandBuilder); err != nil {
		panic(err)
	}
	if err := runner.Run(ctx, prog, []string{"Yep"}, myCommandBuilder); err != nil {
		panic(err)
	}

	fmt.Print(runner.StdoutBuffer.String())

	// Output:
	// Hello World!
	// Hello Yep!
}
