package cmdy_test

import (
	"context"
	"fmt"

	"github.com/ArtProcessors/cmdy"
	"github.com/ArtProcessors/cmdy/arg"
)

type mainGroupCommand struct{}

func newMainGroupCommand() cmdy.Command {
	return &mainGroupCommand{}
}

func (cmd *mainGroupCommand) Help() cmdy.Help {
	return cmdy.Synopsis("Example main group command")
}

func (cmd *mainGroupCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {}

func (cmd *mainGroupCommand) Run(ctx cmdy.Context) error {
	fmt.Fprintf(ctx.Stdout(), "subcommand!\n")
	return nil
}

func Example_maingroup() {
	// Ignore this, pretend it isn't here.
	cmdy.Reset()

	// builders allow multiple instances of the command to be created.
	mainBuilder := func() cmdy.Command {
		// flag values should be scoped to the builder:
		var testFlag bool

		return cmdy.NewGroup(
			"myprog",
			cmdy.Builders{
				// Add your subcommand builders here. This has the same signature as
				// mainBuilder - you can nest cmdy.Groups arbitrarily.
				"mycmd": newMainGroupCommand,
			},

			// Optionally override the default usage:
			cmdy.GroupUsage("Usage: yep!"),

			// Optionally provide global flags. Flags are left-associative in cmdy, so
			// any flag in here must come before the subcommand:
			cmdy.GroupFlags(func() *cmdy.FlagSet {
				set := cmdy.NewFlagSet()
				set.BoolVar(&testFlag, "testflag", false, "test flag")
				return set
			}),

			// Optionally provide global setup code to run before the Group's
			// subcommand is created and run:
			cmdy.GroupBefore(func(ctx cmdy.Context) error {
				fmt.Fprintf(ctx.Stdout(), "before! flag: %v\n", testFlag)
				return nil
			}),

			// Optionally provide global teardown code to run after the Group's
			// subcommand is run:
			cmdy.GroupAfter(func(ctx cmdy.Context, err error) error {
				fmt.Fprintf(ctx.Stdout(), "after! flag: %v\n", testFlag)

				// Careful: failing to return the passed in error here will
				// swallow the error.
				return err
			}),
		)
	}

	args := []string{"-testflag", "mycmd"}
	if err := cmdy.Run(context.Background(), args, mainBuilder); err != nil {
		cmdy.Fatal(err)
	}

	// Output:
	// before! flag: true
	// subcommand!
	// after! flag: true
}
