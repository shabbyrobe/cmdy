package cmdy

import (
	"context"

	"github.com/shabbyrobe/cmdy/args"
)

type Builders map[string]Builder

type CommandSet struct {
	Builders Builders
	Default  Builder
	Unknown  Builder

	Before      func(input Context) error
	After       func(input Context) error
	FlagBuilder func() *FlagSet

	usage    string
	synopsis string

	subcommand     string
	subcommandArgs []string
}

func (cs *CommandSet) Synopsis() string { return cs.synopsis }
func (cs *CommandSet) Usage() string    { return cs.usage }

func (cs *CommandSet) Flags() *FlagSet {
	if cs.FlagBuilder != nil {
		return cs.FlagBuilder()
	}
	return nil
}

func (cs *CommandSet) Args() *args.ArgSet {
	as := args.NewArgSet()
	as.String(&cs.subcommand, "cmd", "Subcommand name")
	as.Remaining(&cs.subcommandArgs, "args", args.AnyLen, "Subcommand arguments")
	return as
}

func (cs *CommandSet) Run(ctx context.Context, in Input) error {
	panic("not implemented")
}
