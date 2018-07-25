package cmdy

import (
	"github.com/shabbyrobe/cmdy/args"
)

// Usage is an optional interface you can add to a Command to specify a more
// complete help message that will be shown by cli.Fatal() if a UsageError is
// returned (for example when the '-help' flag is passed).
//
// The string returned by Usage() is a Go template
// (https://golang.org/pkg/text/template/). The template makes the following
// functions available:
//
//	{{Invocation}}
//		Full invocation string for the command, i.e.
//		'cmd sub subsub [options] <args...>'.
//		This invocation does not include parent command flags.
//  {{Synopsis}}
//		Command.Synopsis()
//	{{CommandFull}}
//		Full command name including all parent commands, i.e. 'cmd sub subsub'.
//	{{Command}}
//		Current command name, not including parent command names. i.e. for
//		command 'cmd sub subsub', only 'subsub' is returned.
//
type Usage interface {
	Usage() string
}

type Command interface {
	Synopsis() string

	Flags() *FlagSet
	Args() *args.ArgSet
	Run(Context) error
}

// Builder creates an instance of your Command. The instance should be a new
// instance, not a recycled instance.
type Builder func() (Command, error)
