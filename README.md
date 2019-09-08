cmdy: Go library for implementing CLI programs
==============================================

[![GoDoc](https://godoc.org/github.com/shabbyrobe/cmdy?status.svg)](https://godoc.org/github.com/shabbyrobe/cmdy)

`cmdy` combines the features I like from the `flag` stdlib package with the
features I like from https://github.com/google/subcommands.

`cmdy` focuses on minimalism and tries to imitate and leverage the stdlib as
much as possible. It does not attempt to replace `flag.Flag`, though it does
extend it slightly.

`cmdy` is liberally documented in [GoDoc](https://godoc.org/github.com/shabbyrobe/cmdy),
though a brief overview is provided in this README. If anything is unclear, please
submit a GitHub issue.


Features
--------

- `ArgSet`, similar to `flag.FlagSet` but for positional arguments. The
  `github.com/shabbyrobe/cmdy/arg` package can be used independently.
- Simple subcommand (and sub-sub command (and sub-sub-sub command)) support.
- `context.Context` support (via `cmdy.Context`, which is also a
  `context.Context`).
- Automatic (but customisable) usage and invocation strings.


Usage
-----

Subcommands are easy to create; you need a builder and one or more
implementations of cmdy.Command. This fairly contrived example demonstrates
the basics:

```go
type demoCommand struct {
	testFlag string
	testArg  string
	rem      []string
}

func newDemoCommand() cmdy.Command {
	return &demoCommand{}
}

func (t *demoCommand) Synopsis() string { return "My command is a command that does stuff" }

func (t *demoCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	flags.StringVar(&t.testFlag, "test", "", "Test flag")
	args.String(&t.testArg, "test", "Test arg")
	args.Remaining(&t.rem, "things", arg.AnyLen, "Any number of extra string arguments.")
}

func (t *demoCommand) Run(ctx cmdy.Context) error {
	fmt.Println(t.testFlag, t.testArg, t.rem)
	return nil
}

func main() {
	if err := run(); err != nil {
		cmdy.Fatal(err)
	}
}

func run() error {
	nestedGroupBuilder := func() cmdy.Command {
		return cmdy.NewGroup(
			"Nested group",
			cmdy.Builders{"subcmd": newDemoCommand},
		)
	}

	mainGroupBuilder := func() cmdy.Command {
		return cmdy.NewGroup(
			"My command group",
			cmdy.Builders{
				"cmd":  newDemoCommand,
				"nest": nestedGroupBuilder,
			},
		)
	}
	return cmdy.Run(context.Background(), os.Args[1:], mainGroupBuilder)
}
```

You can customise the help message by implementing the optional `cmdy.Usage`
interface. ``Usage()`` is a go ``text/template`` that has access to your
command as its vars:

```go
const myCommandUsage = `
{{Synopsis}}

Usage: {{Invocation}}

Stuff is {{.Stuff}}
`

type myCommand struct {
    Stuff string
}

var _ cmdy.Usage = &myCommand{}

func (t *myCommand) Usage() string { return myCommandUsage }

// myCommand implements the rest of cmdy.Command...
```
