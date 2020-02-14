package cmdy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ArtProcessors/cmdy/arg"
	"github.com/ArtProcessors/cmdy/internal/wrap"
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

// GroupUsage sets the Usage string in the Group's Help. The result of this function may
// be cached.
func GroupUsage(usage string) GroupOption {
	return func(grp *Group) { grp.help.Usage = usage }
}

// GroupExamples sets the Examples in the Group's Help. The result of this function may
// be cached.
func GroupExamples(examples ...Example) GroupOption {
	return func(grp *Group) { grp.help.Examples = examples }
}

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

	help   Help
	hidden map[string]bool

	state GroupRunState
}

var _ Command = &Group{}

func NewGroup(synopsis string, builders Builders, opts ...GroupOption) *Group {
	grp := &Group{
		help:     Help{Synopsis: synopsis},
		Builders: builders,
	}
	for _, o := range opts {
		o(grp)
	}
	return grp
}

func (grp *Group) Help() Help { return grp.help }

func (grp *Group) BuildHelp(into *strings.Builder) error {
	into.WriteString("Commands:\n")
	names := make([]string, 0, len(grp.Builders))
	width := 6
	for name := range grp.Builders {
		ln := len(name)
		if ln > width {
			width = ln
		}
		if grp.hidden == nil || !grp.hidden[name] {
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
		s := grp.Builders[l]()
		syn := s.Help().Synopsis
		syn = wrp.Wrap(syn)
		fmt.Fprintf(into, "    %-*s  %s\n", width, l, syn)
	}

	return nil
}

func (grp *Group) Flags() *FlagSet {
	if grp.FlagBuilder != nil {
		return grp.FlagBuilder()
	}
	return nil
}

func (grp *Group) Configure(flags *FlagSet, args *arg.ArgSet) {
	args.HideUsage()
	args.StringOptional(&grp.state.Subcommand, "cmd", "", "Subcommand name")
	args.Remaining(&grp.state.SubcommandArgs, "args", arg.AnyLen, "Subcommand arguments")
}

func (grp *Group) Builder(cmd string) (bld Builder, name string, rerr error) {
	bld, name, rerr = grp.Builders.match(grp.Matcher, cmd)
	return
}

func (grp *Group) Run(ctx Context) error {
	var err error
	grp.state.Builder, grp.state.Name, err = grp.Builder(grp.state.Subcommand)
	if err != nil {
		return err
	}

	if grp.Rewriter != nil {
		newState := grp.Rewriter(grp, grp.state)
		if newState != nil {
			grp.state = *newState
		}
	}

	if grp.state.Builder == nil {
		if grp.state.Subcommand != "" {
			return UsageError(fmt.Errorf("unknown command %q", grp.state.Subcommand))
		} else {
			return UsageError(nil)
		}
	}

	if grp.Before != nil {
		if err := grp.Before(ctx); err != nil {
			return err
		}
	}

	err = ctx.Runner().Run(ctx, grp.state.Name, grp.state.SubcommandArgs, grp.state.Builder)
	if grp.After != nil {
		// XXX: this intentionally replaces the error. it's intended to allow
		// implementers of After to rewrite/replace the error if desired.
		err = grp.After(ctx, err)
	}

	return err
}
