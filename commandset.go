package cmdy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shabbyrobe/cmdy/args"
	"github.com/shabbyrobe/cmdy/usage"
)

type Builders map[string]Builder

type CommandMatcher func(bldrs Builders, in string) (Builder, error)

type CommandSetOption func(cs *CommandSet)

func CommandSetDefault(b Builder) CommandSetOption {
	return func(cs *CommandSet) { cs.Default = b }
}

func CommandSetUnknown(b Builder) CommandSetOption {
	return func(cs *CommandSet) { cs.Unknown = b }
}

func CommandSetUsage(usage string) CommandSetOption {
	return func(cs *CommandSet) { cs.usage = usage }
}

func CommandSetFlags(fb func() *FlagSet) CommandSetOption {
	return func(cs *CommandSet) { cs.FlagBuilder = fb }
}

func CommandSetBefore(fn func(Context, Input) error) CommandSetOption {
	return func(cs *CommandSet) { cs.Before = fn }
}

func CommandSetAfter(fn func(Context, Input, error) error) CommandSetOption {
	return func(cs *CommandSet) { cs.After = fn }
}

type CommandSet struct {
	// All Builders in this map will be called in order to create the Usage
	// string.
	Builders Builders

	Default Builder
	Unknown Builder

	Before      func(Context, Input) error
	After       func(Context, Input, error) error
	FlagBuilder func() *FlagSet

	matcher  CommandMatcher
	usage    string
	synopsis string

	// State:
	subcommand     string
	subcommandArgs []string
}

var _ Command = &CommandSet{}

func NewCommandSet(synopsis string, builders Builders, opts ...CommandSetOption) *CommandSet {
	cs := &CommandSet{
		synopsis: synopsis,
		Builders: builders,
	}
	for _, o := range opts {
		o(cs)
	}
	return cs
}

func (cs *CommandSet) Synopsis() string { return cs.synopsis }

func (cs *CommandSet) Usage() string {
	out := cs.usage
	if out == "" {
		out = cs.synopsis
	}
	out = strings.TrimSpace(out)

	if out != "" {
		out += "\n\n"
	}

	out += "Commands:\n"
	names := make([]string, 0, len(cs.Builders))
	width := 6
	for name := range cs.Builders {
		ln := len(name)
		if ln > width {
			width = ln
		}
		names = append(names, name)
	}
	sort.Strings(names)

	// +4 == command name, +2 == space between command name and synopsis
	indent := strings.Repeat(" ", width+4+2)

	for _, l := range names {
		s, err := cs.Builders[l]()
		if err == nil {
			syn := s.Synopsis()
			syn = usage.Wrap(syn, indent, 0)
			out += fmt.Sprintf("    %-*s  %s\n", width, l, syn)
		}
	}
	return out
}

func (cs *CommandSet) Flags() *FlagSet {
	if cs.FlagBuilder != nil {
		return cs.FlagBuilder()
	}
	return nil
}

func (cs *CommandSet) Args() *args.ArgSet {
	as := args.NewArgSet()
	as.HideUsage()
	as.StringOptional(&cs.subcommand, "cmd", "", "Subcommand name")
	as.Remaining(&cs.subcommandArgs, "args", args.AnyLen, "Subcommand arguments")
	return as
}

func (cs *CommandSet) Run(ctx Context, in Input) error {
	var bld Builder
	if cs.matcher != nil {
		var err error
		bld, err = cs.matcher(cs.Builders, cs.subcommand)
		if err != nil {
			return err
		}
	} else {
		bld = cs.Builders[cs.subcommand]
	}

	if bld == nil {
		if cs.Unknown != nil {
			bld = cs.Unknown
		} else {
			return NewUsageError(fmt.Errorf("unknown command %q", cs.subcommand))
		}
	}

	if cs.Before != nil {
		if err := cs.Before(ctx, in); err != nil {
			return err
		}
	}

	err := Run(ctx, cs.subcommandArgs, bld)
	if cs.After != nil {
		err = cs.After(ctx, in, err)
	}

	return err
}
