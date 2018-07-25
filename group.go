package cmdy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shabbyrobe/cmdy/args"
	"github.com/shabbyrobe/cmdy/usage"
)

type Builders map[string]Builder

type CommandMatcher func(bldrs Builders, in string) (bld Builder, name string, rerr error)

type GroupOption func(cs *Group)

func GroupDefault(b Builder) GroupOption {
	return func(cs *Group) { cs.Default = b }
}

func GroupUnknown(b Builder) GroupOption {
	return func(cs *Group) { cs.Unknown = b }
}

func GroupUsage(usage string) GroupOption {
	return func(cs *Group) { cs.usage = usage }
}

func GroupFlags(fb func() *FlagSet) GroupOption {
	return func(cs *Group) { cs.FlagBuilder = fb }
}

func GroupBefore(fn func(Context) error) GroupOption {
	return func(cs *Group) { cs.Before = fn }
}

func GroupAfter(fn func(Context, error) error) GroupOption {
	return func(cs *Group) { cs.After = fn }
}

type Group struct {
	// All Builders in this map will be called in order to create the Usage
	// string.
	Builders Builders

	Default Builder
	Unknown Builder

	Before      func(Context) error
	After       func(Context, error) error
	FlagBuilder func() *FlagSet

	matcher  CommandMatcher
	usage    string
	synopsis string

	// State:
	subcommand     string
	subcommandArgs []string
}

var _ Command = &Group{}

func NewGroup(synopsis string, builders Builders, opts ...GroupOption) *Group {
	cs := &Group{
		synopsis: synopsis,
		Builders: builders,
	}
	for _, o := range opts {
		o(cs)
	}
	return cs
}

func (cs *Group) Synopsis() string { return cs.synopsis }

func (cs *Group) Usage() string {
	out := cs.usage
	if out == "" {
		out = defaultUsage
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

func (cs *Group) Flags() *FlagSet {
	if cs.FlagBuilder != nil {
		return cs.FlagBuilder()
	}
	return nil
}

func (cs *Group) Args() *args.ArgSet {
	as := args.NewArgSet()
	as.HideUsage()
	as.StringOptional(&cs.subcommand, "cmd", "", "Subcommand name")
	as.Remaining(&cs.subcommandArgs, "args", args.AnyLen, "Subcommand arguments")
	return as
}

func (cs *Group) Run(ctx Context) error {
	var (
		bld  Builder
		name string
	)
	if cs.matcher != nil {
		var err error
		bld, name, err = cs.matcher(cs.Builders, cs.subcommand)
		if err != nil {
			return err
		}
	} else {
		bld, name = cs.Builders[cs.subcommand], cs.subcommand
	}

	if bld == nil {
		if cs.Unknown != nil {
			bld = cs.Unknown
		} else {
			return NewUsageError(fmt.Errorf("unknown command %q", cs.subcommand))
		}
	}

	if cs.Before != nil {
		if err := cs.Before(ctx); err != nil {
			return err
		}
	}

	err := ctx.Runner().Run(ctx, name, cs.subcommandArgs, bld)
	if cs.After != nil {
		err = cs.After(ctx, err)
	}

	return err
}
