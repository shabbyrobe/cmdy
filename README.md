cmdy: Go library for implementing CLI programs
==============================================

[![GoDoc](https://godoc.org/github.com/shabbyrobe/cmdy?status.svg)](https://godoc.org/github.com/shabbyrobe/cmdy)

`cmdy` combines the features I like from the `flag` stdlib package with the
features I like from https://github.com/google/subcommands.

`cmdy` probably doesn't really need to exist, but I like it and use it for
my own projects. There are a lot of CLI libraries for Go but this one is mine.

`cmdy` focuses on minimalism and tries to imitate and leverage the stdlib as
much as possible. It does not attempt to replace `flag.Flag`, though it does
extend it slightly.

`cmdy` is liberally documented in [GoDoc](https://godoc.org/github.com/shabbyrobe/cmdy),
though a brief overview is provided in this README. If anything is unclear, please
submit a GitHub issue.

`cmdy` has no dependencies beyond the stdlib.


Features
--------

- `ArgSet`, similar to `flag.FlagSet` but for positional arguments. The
  `github.com/shabbyrobe/cmdy/arg` package can be used independently.
- Simple subcommand (and sub-sub command (and sub-sub-sub command)) support.
- `context.Context` support (via `cmdy.Context`, which is also a
  `context.Context`).
- Automatic (but customisable) usage and invocation strings.
- Ctrl-C propagation via `cmdy.Context` (see `cmdyutil.InterruptibleRun`).


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

func (cmd *demoCommand) Help() cmdy.Help {
    return cmdy.Synopsis("My command is a command that does stuff")
}

func (cmd *demoCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
    flags.StringVar(&cmd.testFlag, "test", "", "Test flag")
    args.String(&cmd.testArg, "test", "Test arg")
    args.Remaining(&cmd.rem, "things", arg.AnyLen, "Any number of extra string arguments.")
}

func (cmd *demoCommand) Run(ctx cmdy.Context) error {
    // Use ctx.Stdout() if you want to make it easier to test your command,
    // but it's fine to just use fmt.Println():
    fmt.Fprintln(ctx.Stdout(), cmd.testFlag, cmd.testArg, cmd.rem)
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

You can add more Usage information by returning a `cmdy.Help` structure from
the `Help()` method:

```go
const demoCommandUsage = `
Additional help for the command
`

func (cmd *demoCommand) Help() cmdy.Help {
    return cmdy.Help{
        Synopsis: "My command is a command that does stuff",
        Usage:    demoCommandUsage,
    }
}
```

You can also add examples (which can be run using 
`github.com/shabbyrobe/cmdy/cmdytest.ExampleTester` to the `cmdy.Help` structure.
Individual fields are explained in more detail in godoc:

```go
func (cmd *demoCommand) Help() cmdy.Help {
    return cmdy.Help{
        // ...
        Examples: cmdy.Examples{
            {
                Desc:    "Do a thing with this command:",
                Command: "-test foo bar baz", // <-- start at this command's flags, not parents
            },
            {
                Desc:    "This won't work",
                Command: "-flem flarg flub",
                Code:    cmdy.ExitUsage,
            },
        },
    }
}
```

