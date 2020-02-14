package cmdy

import (
	"github.com/ArtProcessors/cmdy/arg"
)

/*
Command represents a single command in a cmdy application.

A Command is created by a Builder, which is itself passed to cmdy.Run.

To dispatch subcommands to multiple Command implementations, see cmdy.Group.

Minimal 'Hello World' implementation:

	type myCommand struct {}

	func myCommandBuilder() cmdy.Command { return &myCommand{} }

	func (cmd *myCommand) Help() Help { return cmdy.Synopsis("Hello world") }

	func (cmd *myCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {}

	func (cmd *myCommand) Run(ctx Context) error {
		fmt.Fprintln(ctx.Stdout(), "Hello world!")
		return nil
	}

	func run() error {
		return cmdy.Run(context.Background(), os.Args[1:], myCommandBuilder)
	}

	func main() {
		if err := run(); err != nil {
			cmdy.Fatal(err)
		}
	}
*/
type Command interface {
	Help() Help

	// Configure allows you to add flags and args to the FlagSet and ArgSet respectively.
	// You should not call any methods on them that aren't directly related to adding
	// flags and args.
	Configure(flags *FlagSet, args *arg.ArgSet)

	Run(ctx Context) error
}

/*
Builder creates an instance of your Command. The instance returned should be a new
instance, not a recycled instance, and should only contain static dependency values that
are cheap to create:

	var goodBuilder = func() cmdy.Command {
		return &MyCommand{}
	}
	var goodBuilder = func() cmdy.Command {
		return &MyCommand{SimpleDep: "hello"}
	}
	var badBuilder = func() cmdy.Command {
		body, _ := http.Get("http://example.com")
		return &MyCommand{Stuff: body}
	}

	cmd := &MyCommand{}
	var badBuilder = func() cmdy.Command {
		return cmd
	}

*/
type Builder func() Command
