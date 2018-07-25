cmdy: Go library for implementing CLI programs
==============================================

Combines the features I like from the `flag` stdlib package with the features
I like from https://github.com/google/subcommands.

This is a very early pass of basic functionality, it is not even remotely close
to what I'm aiming for, but it's a start.

Please don't use this yet.


Features
--------

- `ArgSet`, similar to `flag.FlagSet` but for positional arguments
- Simple subcommand support
- `context.Context` support (via `cmdy.Context`, which is a `context.Context`)


Usage
-----

Subcommands are easy to create; you need a builder and a command:

```go
func myCommandBuilder() (cmdy.Command, error) {
	return &myCommand{}, nil
}

const myCommandUsage = `
{{Synopsis}}

Usage: {{Invocation}}
`

type myCommand struct {
	testFlag string
	testArg  string
	rem      []string
}

var _ cmdy.Command = &myCommand{}

func (t *myCommand) Synopsis() string { return "My command is a command that does stuff" }

// Usage is optional; it allows you to specify a Go template that will
// return a full help message.
func (t *myCommand) Usage() string { return myCommandUsage }

func (t *myCommand) Flags() *cmdy.FlagSet {
	fs := cmdy.NewFlagSet()
	fs.StringVar(&t.testFlag, "test", "", "Test flag")
	return fs
}

func (t *myCommand) Args() *args.ArgSet {
	as := args.NewArgSet()
	as.String(&t.testArg, "test", "Test arg")
	as.Remaining(&t.rem, "things", args.AnyLen, "Any number of extra string arguments.")
	return as
}

func (t *myCommand) Run(ctx cmdy.Context) error {
	fmt.Println(t.testFlag, t.testArg, t.rem)
	return nil
}

func main() {
	if err := run(); err != nil {
		cmdy.Fatal(err)
	}
}

func run() error {
	bld := func() (cmdy.Command, error) {
		return cmdy.NewGroup(
			"My command group",
			cmdy.Builders{
				"cmd": myCommandBuilder,
				"nest": func() (cmdy.Command, error) {
					return cmdy.NewGroup(
						"Nested group",
						cmdy.Builders{
							"subcmd": myCommandBuilder,
						},
					)
				},
			},
		), nil
	}
	return cmdy.Run(context.Background(), os.Args[1:], bld)
}
```
