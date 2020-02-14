package arg

import (
	"fmt"
	"time"

	"github.com/ArtProcessors/cmdy/usage"
)

// An ArgSet represents a set of defined command-line arguments. It is intended
// to mirror, as closely as reasonable, the structure of the flag.FlagSet
// package in the Go standard library.
type ArgSet struct {
	args      []*Arg
	optional  bool
	remaining *remaining
	hideUsage bool
}

func NewArgSet() *ArgSet {
	return &ArgSet{}
}

// HideUsage prevents the Usage string for the ArgSet being built. It is
// primarily intended to help if you have a dynamic subcommand dispatch and
// don't want spurious and redundant usage documentation for things like
// "Usage: <command> <args>...".
func (a *ArgSet) HideUsage() { a.hideUsage = true }

func (a *ArgSet) Usage() string {
	if a.hideUsage {
		return ""
	}

	usables := make([]usage.Usable, len(a.args))
	for i, a := range a.args {
		usables[i] = a
	}
	return usage.Usage(0, usables...)
}

// Invocation returns an example command invocation string intended for
// display in the "Usage:" section of a command's help message.
func (a *ArgSet) Invocation() string {
	var inv string

	for idx, arg := range a.args {
		if idx > 0 {
			inv += " "
		}
		inv += arg.Describe("", "")
	}

	return inv
}

// NArg returns the number of args that have been defined. A "remaining"
// arg counts as one arg.
func (a *ArgSet) NArg() int {
	return len(a.args)
}

func (a *ArgSet) Parse(input []string) error {
	inputLen := len(input)

	var consumed int

	for idx, arg := range a.args {
		if a.remaining != nil && arg.value == a.remaining {
			var left []string
			if idx < inputLen {
				left = input[idx:]
			}
			leftLen := len(left)
			if leftLen < a.remaining.Min {
				return fmt.Errorf("expected at least %d remaining args at position %d, found %d", a.remaining.Min, idx+1, leftLen)
			}
			if a.remaining.Max >= 0 && leftLen > a.remaining.Max {
				return fmt.Errorf("expected at most %d remaining args at position %d, found %d", a.remaining.Max, idx+1, leftLen)
			}
			for _, rem := range left {
				if err := a.remaining.Set(rem); err != nil {
					return err // FIXME
				}
				consumed++
			}
			break

		} else {
			if idx >= inputLen {
				if !arg.optional {
					return fmt.Errorf("missing arg %s at position %d", arg.Describe("", ""), idx+1)
				}
			} else {
				if err := arg.value.Set(input[idx]); err != nil {
					return fmt.Errorf("arg invalid at position %d: %v", idx+1, err)
				}
			}
		}
		consumed++
	}

	if consumed < inputLen {
		extra, s := inputLen-consumed, ""
		if extra != 1 {
			s = "s"
		}
		return fmt.Errorf("found %d additional arg%s", extra, s)
	}

	return nil
}

// Remaining collects all args after the last defined argument into the slice
// of strings pointed to by p. If more args are defined after any Remaining
// method is called, args will panic.
//
//	Use arg.AnyLen to allow an arbitrary number of remaining args.
// 	Use arg.Min(2) to require at least 2 args.
// 	Use arg.Max(2) to require at most 2 args.
// 	Use arg.MinMax(1, 3) to require at least 1 arg and at most 3 args.
//
func (a *ArgSet) Remaining(p *[]string, name string, minmax Range, usage string) {
	a.RemainingVar((*stringList)(p), name, minmax, usage)
}

// RemainingInts collects all args after the last defined argument into the
// slice of ints pointed to by p. If more args are defined after any Remaining
// method is called, args will panic.
//
// See Remaining for an explanation of minmax and usage.
func (a *ArgSet) RemainingInts(p *[]int, name string, minmax Range, usage string) {
	a.RemainingVar((*intList)(p), name, minmax, usage)
}

// RemainingInt64s collects all args after the last defined argument into the
// slice of int64s pointed to by p. If more args are defined after any Remaining
// method is called, args will panic.
//
// See Remaining for an explanation of minmax and usage.
func (a *ArgSet) RemainingInt64s(p *[]int64, name string, minmax Range, usage string) {
	a.RemainingVar((*int64List)(p), name, minmax, usage)
}

// RemainingInt64s collects all args after the last defined argument into the
// slice of int64s pointed to by p. If more args are defined after any Remaining
// method is called, args will panic.
//
// See Remaining for an explanation of minmax and usage.
func (a *ArgSet) RemainingUints(p *[]uint, name string, minmax Range, usage string) {
	a.RemainingVar((*uintList)(p), name, minmax, usage)
}

// RemainingUint64s collects all args after the last defined argument into the
// slice of uint64s pointed to by p. If more args are defined after any Remaining
// method is called, args will panic.
//
// See Remaining for an explanation of minmax and usage.
func (a *ArgSet) RemainingUint64s(p *[]uint64, name string, minmax Range, usage string) {
	a.RemainingVar((*uint64List)(p), name, minmax, usage)
}

// RemainingFloat64s collects all args after the last defined argument into the
// slice of float64s pointed to by p. If more args are defined after any Remaining
// method is called, args will panic.
//
// See Remaining for an explanation of minmax and usage.
func (a *ArgSet) RemainingFloat64s(p *[]float64, name string, minmax Range, usage string) {
	a.RemainingVar((*float64List)(p), name, minmax, usage)
}

// RemainingVar collects all args after the last defined argument into the
// slice of ArgVals pointed to by p. If more args are defined after any
// Remaining method is called, args will panic.
//
// See Remaining for an explanation of minmax and usage.
func (a *ArgSet) RemainingVar(val ArgVal, name string, minmax Range, usage string) {
	a.Var(&remaining{val, minmax}, name, usage)
}

func (a *ArgSet) String(p *string, name string, usage string) { a.Var((*stringArg)(p), name, usage) }

func (a *ArgSet) StringOptional(p *string, name string, value string, usage string) {
	a.optional = true
	*p = value
	a.String(p, name, usage)
}

func (a *ArgSet) Int(p *int, name string, usage string) { a.Var((*intArg)(p), name, usage) }

func (a *ArgSet) IntOptional(p *int, name string, value int, usage string) {
	a.optional = true
	*p = value
	a.Int(p, name, usage)
}

func (a *ArgSet) Int64(p *int64, name string, usage string) { a.Var((*int64Arg)(p), name, usage) }

func (a *ArgSet) Int64Optional(p *int64, name string, value int64, usage string) {
	a.optional = true
	*p = value
	a.Int64(p, name, usage)
}

func (a *ArgSet) Uint(p *uint, name string, usage string) { a.Var((*uintArg)(p), name, usage) }

func (a *ArgSet) UintOptional(p *uint, name string, value uint, usage string) {
	a.optional = true
	*p = value
	a.Uint(p, name, usage)
}

func (a *ArgSet) Uint64(p *uint64, name string, usage string) { a.Var((*uint64Arg)(p), name, usage) }

func (a *ArgSet) Uint64Optional(p *uint64, name string, value uint64, usage string) {
	a.optional = true
	*p = value
	a.Uint64(p, name, usage)
}

func (a *ArgSet) Float64(p *float64, name string, usage string) { a.Var((*float64Arg)(p), name, usage) }

func (a *ArgSet) Float64Optional(p *float64, name string, value float64, usage string) {
	a.optional = true
	*p = value
	a.Float64(p, name, usage)
}

func (a *ArgSet) Bool(p *bool, name string, usage string) { a.Var((*boolArg)(p), name, usage) }

func (a *ArgSet) BoolOptional(p *bool, name string, value bool, usage string) {
	a.optional = true
	*p = value
	a.Bool(p, name, usage)
}

func (a *ArgSet) Duration(p *time.Duration, name string, usage string) {
	a.Var((*durationArg)(p), name, usage)
}

func (a *ArgSet) DurationOptional(p *time.Duration, name string, value time.Duration, usage string) {
	a.optional = true
	*p = value
	a.Duration(p, name, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type ArgVal,
// which typically holds a user-defined implementation of ArgVal. For instance,
// the caller could create a flag that turns a comma-separated string into a
// slice of strings by giving the slice the methods of ArgVal; in particular,
// Set would decompose the comma-separated string into the slice.
//
// ArgVal is 100% compatible with the flag.Getter interface from the Go
// standard library; anything that can be used by flag.FlagSet.Var can be
// used by arg.ArgSet.Var with no modification.
//
func (a *ArgSet) Var(val ArgVal, name string, usage string) {
	dflt := val.String()

	// FIXME: should be possible; only one 'remaining' should be allowed but
	// we can count backwards from the end for the rest.
	if a.remaining != nil {
		panic("cannot add more arguments after accumulating remaining")
	}

	if rem, ok := val.(*remaining); ok {
		a.remaining = rem
	}

	arg := &Arg{name, usage, val, dflt, a.optional}
	a.args = append(a.args, arg)
}
