package cmdy

import (
	"github.com/shabbyrobe/cmdy/args"
)

/*
Usage is an optional interface you can add to a Command to specify a more
complete help message that will be shown by cli.Fatal() if a UsageError is
returned (for example when the '-help' flag is passed).

The string returned by Usage() is parsed by the text/template package
(https://golang.org/pkg/text/template/). The template makes the following
functions available:

	{{Invocation}}
		Full invocation string for the command, i.e.
		'cmd sub subsub [options] <args...>'.
		This invocation does not include parent command flags.
	{{Synopsis}}
		Command.Synopsis()
	{{CommandFull}}
		Full command name including all parent commands, i.e. 'cmd sub subsub'.
	{{Command}}
		Current command name, not including parent command names. i.e. for
		command 'cmd sub subsub', only 'subsub' is returned.
*/
type Usage interface {
	Usage() string
}

type Command interface {
	// Synopsis is the shortest possible complete description of your command,
	// ideally one sentence.
	Synopsis() string

	// Flag definitions for your command. May return nil. If no FlagSet is
	// returned, --help is still supported but all other flags will cause an
	// error.
	Flags() *FlagSet

	// Args defines positional arguments for your command. If you want to accept
	// all args, use github.com/shabbyrobe/cmdy/args.All(). If no ArgSet is
	// returned, any arguments will cause an error.
	Args() *args.ArgSet

	Run(Context) error
}

// Builder creates an instance of your Command. The instance should be a new
// instance, not a recycled instance.
type Builder func() (Command, error)
