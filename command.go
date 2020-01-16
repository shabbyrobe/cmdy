package cmdy

import (
	"github.com/shabbyrobe/cmdy/arg"
)

type Command interface {
	Help() Help

	Configure(flags *FlagSet, args *arg.ArgSet)

	Run(Context) error
}

// CommandArgs allows you to override the construction of the ArgSet in your Command. If
// your command does not implement this, it will receive a fresh instance of arg.ArgSet.
type CommandArgs interface {
	Command

	// Args defines positional arguments for your command. If you want to accept
	// all args, use github.com/shabbyrobe/cmdy/arg.All(). If no ArgSet is
	// returned, any arguments will cause an error.
	Args() *arg.ArgSet
}

// CommandFlags allows you to override the construction of the FlagSet in your Command.
// If your command does not implement this, it will receive a fresh instance of
// cmdy.FlagSet.
type CommandFlags interface {
	Command

	// Flag definitions for your command. May return nil. If no FlagSet is
	// returned, --help is still supported but all other flags will cause an
	// error.
	Flags() *FlagSet
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
