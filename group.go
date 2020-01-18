package cmdy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shabbyrobe/cmdy/arg"
	"github.com/shabbyrobe/cmdy/internal/wrap"
)

// Matcher is a function for choosing a command builder from a list of command builders
// when using a Group.
//
// You could use this API to implement short aliases for existing commands too, if
// you so desired (i.e. "hg co" -> "hg checkout").
//
// See GroupMatcher, GroupPrefixMatcher and PrefixMatcher.
//
// WARNING: This API may change to return a list of possible options when the
// choice is ambiguous.
type Matcher func(bldrs Builders, in string) (bld Builder, name string, rerr error)

type GroupOption func(grp *Group)

// GroupMatcher assigns a Matcher function to the group. Matcher is used for choosing a
// command builder from a list of command builders when using a Group.
func GroupMatcher(cm Matcher) GroupOption { return func(grp *Group) { grp.Matcher = cm } }

// GroupPrefixMatcher assigns the PrefixMatcher of the specified minimum length
// to the Group's Matcher.
func GroupPrefixMatcher(minLen int) GroupOption {
	return func(grp *Group) { grp.Matcher = PrefixMatcher(grp, minLen) }
}

// GroupUsage provides the usage template to the Group.
// The result of this function may be cached.
func GroupUsage(usage string) GroupOption { return func(cs *Group) { cs.usage = usage } }

// GroupFlags provides a function that creates a FlagSet to the Group.
// This function may return nil. The result of this function may be cached.
func GroupFlags(fb func() *FlagSet) GroupOption {
	return func(cs *Group) { cs.FlagBuilder = fb }
}

// GroupBefore provides a function to call Before a Group's subcommand is
// executed.
//
// Any error returned by the before function will prevent the subcommand from
// executing.
func GroupBefore(before func(Context) error) GroupOption {
	return func(cs *Group) { cs.Before = before }
}

// GroupAfter provides a function to call After a Group's subcommand is
// executed.
//
// Any error returned by the subcommand is passed to the function. If it
// is not returned, it will be swallowed.
func GroupAfter(fn func(Context, error) error) GroupOption {
	return func(cs *Group) { cs.After = fn }
}

// GroupHide hides the builders from the usage string. If the builder does not
// exist in Builders, it will panic.
func GroupHide(names ...string) GroupOption {
	return func(cs *Group) {
		for _, name := range names {
			if _, ok := cs.Builders[name]; !ok {
				panic(fmt.Errorf("cannot hide unknown builder %q", name))
			}
			if cs.hidden == nil {
				cs.hidden = make(map[string]bool, len(names))
			}
			cs.hidden[name] = true
		}
	}
}

// NOTE: This is experimental and may change.
type GroupRewriter func(grp *Group, args GroupRunState) (out *GroupRunState)

// NOTE: This is experimental and may change.
func GroupRewrite(rw GroupRewriter) GroupOption { return func(grp *Group) { grp.Rewriter = rw } }

// NOTE: This is experimental and may change.
type GroupRunState struct {
	// Builder of the subcommand to be run. May be nil if none was found for
	// the Subcommand arg. You may replace this with any builder you like.
	Builder Builder

	// The name of the builder, which may be the same as Subcommand, unless
	// it has been modified by a Matcher.
	Name string

	// The first argument passed to the Group
	Subcommand string

	// The remaining arguments passed to the Group
	SubcommandArgs []string
}

// Builders is used by Group to define a list of subcommand arguments to their
// corresponding Command Builder function. Typical usage is within a Group's
// constructor:
//
//	grp := cmdy.NewGroup("My Group", cmdy.Builders{
//		"sub1": newSub1Command,
//		"sub2": newSub2Command,
//	})
//
type Builders map[string]Builder

func (builders Builders) match(matcher Matcher, search string) (bld Builder, name string, err error) {
	if matcher != nil {
		return matcher(builders, search)
	} else {
		bld, name = builders[search], search
		return bld, name, nil
	}
}

// Group implements a command that delegates to one or more subcommand. It selects a
// single Builder from Builders based on the value of the first non-flag argument.
type Group struct {
	// Builders contains mappings between command names (received as the first
	// argument to this command) and the builder to delegate to.
	//
	// All Builders in this map will be called in order to create the Usage
	// string.
	Builders Builders

	// Allows interception of command strings so you can rewrite them to
	// other commands. Useful for aliases, or for handling the case where
	// no subcommand argument is present.
	Rewriter GroupRewriter

	Before      func(Context) error
	After       func(Context, error) error
	FlagBuilder func() *FlagSet
	Matcher     Matcher

	usage    string
	synopsis string
	hidden   map[string]bool

	state GroupRunState
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

func (cs *Group) Help() Help {
	return Help{
		Synopsis: cs.synopsis,
		Usage:    cs.Usage(),
	}
}

func (cs *Group) Usage() string {
	out := cs.usage
	if out == "" {
		out = DefaultUsage
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
		if cs.hidden == nil || !cs.hidden[name] {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	// +4 == command name, +2 == space between command name and synopsis
	indent := make([]byte, width+4+2)
	for i := 0; i < len(indent); i++ {
		indent[i] = ' '
	}

	wrp := wrap.Wrapper{Indent: string(indent)}

	for _, l := range names {
		s := cs.Builders[l]()
		syn := s.Help().Synopsis
		syn = wrp.Wrap(syn)
		out += fmt.Sprintf("    %-*s  %s\n", width, l, syn)
	}
	return out
}

func (cs *Group) Flags() *FlagSet {
	if cs.FlagBuilder != nil {
		return cs.FlagBuilder()
	}
	return nil
}

func (cs *Group) Configure(flags *FlagSet, args *arg.ArgSet) {
	args.HideUsage()
	args.StringOptional(&cs.state.Subcommand, "cmd", "", "Subcommand name")
	args.Remaining(&cs.state.SubcommandArgs, "args", arg.AnyLen, "Subcommand arguments")
}

func (grp *Group) Builder(cmd string) (bld Builder, name string, rerr error) {
	bld, name, rerr = grp.Builders.match(grp.Matcher, cmd)
	return
}

func (cs *Group) Run(ctx Context) error {
	var err error
	cs.state.Builder, cs.state.Name, err = cs.Builder(cs.state.Subcommand)
	if err != nil {
		return err
	}

	if cs.Rewriter != nil {
		newState := cs.Rewriter(cs, cs.state)
		if newState != nil {
			cs.state = *newState
		}
	}

	if cs.state.Builder == nil {
		if cs.state.Subcommand != "" {
			return UsageError(fmt.Errorf("unknown command %q", cs.state.Subcommand))
		} else {
			return UsageError(nil)
		}
	}

	if cs.Before != nil {
		if err := cs.Before(ctx); err != nil {
			return err
		}
	}

	err = ctx.Runner().Run(ctx, cs.state.Name, cs.state.SubcommandArgs, cs.state.Builder)
	if cs.After != nil {
		err = cs.After(ctx, err)
	}

	return err
}
