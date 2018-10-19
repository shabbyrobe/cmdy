package cmdy

import (
	"github.com/shabbyrobe/cmdy/args"
)

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

	{{if ShowFullHelp}}...{{end}}
		Help section contained inside the '...' should only be shown if the
		command's '--help' was requested, not if the command's usage is to
		be shown.


If your Command does not implement cmdy.Usage, cmdy.DefaultUsage is used.

Your Command instance is used as the 'data' argument to Template.Execute(),
so any exported fields from your command can be used in the template like
so: "{{.MyCommandField}}".

If a Command intends cmdy to print the usage in response to an error,
cmdy.NewUsageError or cmdy.NewUsageErrorf should be returned from Command.Run().

To obtain an actual usage string from a usage error, use cmdy.Format(err).
*/
type Usage interface {
	Usage() string
}

// Builder creates an instance of your Command. The instance should be a new
// instance, not a recycled instance.
type Builder func() (Command, error)
