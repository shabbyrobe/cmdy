package cmdy

import (
	"flag"
	"fmt"

	"github.com/ArtProcessors/cmdy/usage"
)

// FlagDoubleDash allows you to globally configure whether long flag names will
// show in the help message with two dashes or one. Some people seem to dislike
// the single-dash longopts and Go supports both anyway, so if you are in this
// camp, this is the var for you.
var FlagDoubleDash = false

// FlagSet is a cmdy specific extension of flag.FlagSet; it is intended to
// behave the same way but with a few small extensions for the sake of this
// library. You should use it instead of flag.FlagSet when dealing with cmdy
// (though you can wrap an existing flag.FlagSet with FlagSetFromStd easily).
type FlagSet struct {
	*flag.FlagSet
	WrapWidth int
	hideUsage bool
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{
		FlagSet: flag.NewFlagSet("", flag.ContinueOnError),
	}
	fs.FlagSet.SetOutput(&devNull{})
	return fs
}

// FlagSetFromStd wraps an existing flag.FlagSet instance and ensures it is
// correctly configured for use with cmdy.
//
// The FlagSet will be modified. Custom output handlers are removed and the
// error handling is changed to ContinueOnError.
func FlagSetFromStd(fs *flag.FlagSet) *FlagSet {
	fs.Init(fs.Name(), flag.ContinueOnError)
	fs.SetOutput(&devNull{})
	return &FlagSet{
		FlagSet: fs,
	}
}

// HideUsage prevents the "Flags" section from appearing in the Usage string.
func (fs *FlagSet) HideUsage() { fs.hideUsage = true }

// Invocation string for the flags, for example '[-foo=<yep>] [-bar=<pants>]`.
// If there are too many flags, `[options]` is returned instead.
func (fs *FlagSet) Invocation() string {
	// FIXME: might be nice if this could be configured
	const flagInvocationSpill = 3

	var options string
	var i int

	fs.VisitAll(func(f *flag.Flag) {
		if i >= flagInvocationSpill {
			options = ""
		} else {
			if i > 0 {
				options += " "
			}
			usable := usableFlag{f}
			kind, _ := usage.Kind(usable)
			options += "[" + usable.Describe(kind, "") + "]"
		}
		i++
	})

	if options == "" {
		options = "[options]"
	}
	return options
}

// Usage returns the full usage string for the FlagSet, provided HideUsage()
// has not been set.
func (fs *FlagSet) Usage() string {
	if fs.hideUsage {
		return ""
	}

	var usables = make([]usage.Usable, 0, fs.NFlag())
	fs.VisitAll(func(f *flag.Flag) {
		usables = append(usables, usableFlag{f})
	})
	return usage.Usage(fs.WrapWidth, usables...)
}

type usableFlag struct {
	flag *flag.Flag
}

func (u usableFlag) Name() string       { return u.flag.Name }
func (u usableFlag) Usage() string      { return u.flag.Usage }
func (u usableFlag) DefValue() string   { return u.flag.DefValue }
func (u usableFlag) Value() interface{} { return u.flag.Value }

func (u usableFlag) Describe(kind string, hint string) string {
	name := u.Name()
	dashes := "-"
	if FlagDoubleDash && len(name) > 1 {
		dashes = "--"
	}

	if kind != "" && hint != "" {
		return fmt.Sprintf("%s%s=<%s> (%s)", dashes, name, kind, hint)
	} else if hint != "" {
		return fmt.Sprintf("%s%s (%s)", dashes, name, hint)
	} else if kind != "" {
		return fmt.Sprintf("%s%s=<%s>", dashes, name, kind)
	} else {
		return fmt.Sprintf("%s%s", dashes, name)
	}
}

// devNull avoids a dependency on ioutil
type devNull struct{}

func (devNull) Write(p []byte) (int, error) {
	return len(p), nil
}
