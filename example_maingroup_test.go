package cmdy_test

import (
	"context"
	"fmt"
	"log"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
)

type mainGroupCommand struct{}

func (cmd *mainGroupCommand) Synopsis() string     { return "Example main group command" }
func (cmd *mainGroupCommand) Flags() *cmdy.FlagSet { return nil }
func (cmd *mainGroupCommand) Args() *args.ArgSet   { return nil }

func (cmd *mainGroupCommand) Run(ctx cmdy.Context) error {
	fmt.Fprintf(ctx.Stdout(), "subcommand!\n")
	return nil
}

func Example_maingroup() {
	// builders allow multiple instances of the command to be created.
	mainBuilder := func() (cmdy.Command, cmdy.Init) {
		// flag values should be scoped to the builder:
		var testFlag bool

		return cmdy.NewGroup(
			"myprog",
			cmdy.Builders{
				// Add your subcommand builders here. This has the same signature as
				// mainBuilder - you can nest cmdy.Groups arbitrarily.
				"mycmd": func() (cmdy.Command, cmdy.Init) { return &mainGroupCommand{}, nil },
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
		), nil
	}

	args := []string{"-testflag", "mycmd"}
	if err := cmdy.Run(context.Background(), args, mainBuilder); err != nil {
		log.Fatal(err)
	}

	// Output:
	// before! flag: true
	// subcommand!
	// after! flag: true
}