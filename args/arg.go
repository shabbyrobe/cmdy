package args

import (
	"fmt"
	"time"

	"github.com/shabbyrobe/cmdy/usage"
)

const Unlimited = -1

var AnyLen = Range{0, Unlimited}

type Range struct {
	Min, Max int
}

func Min(min int) Range         { return Range{min, Unlimited} }
func Max(max int) Range         { return Range{0, max} }
func MinMax(min, max int) Range { return Range{min, max} }

type Arg struct {
	name     string
	usage    string
	value    ArgVal
	defValue string
	optional bool
}

func (a *Arg) Name() string     { return a.name }
func (a *Arg) Usage() string    { return a.usage }
func (a *Arg) DefValue() string { return a.defValue }

func (a *Arg) Value() interface{} {
	if rem, ok := a.value.(*remaining); ok {
		return rem.arg
	}
	return a.value
}

func (a *Arg) Describe(kind string, hint string) string {
	name := a.Name()
	if _, ok := a.value.(*remaining); ok {
		name += "..."
	}

	if kind != "" && hint != "" {
		return fmt.Sprintf("<%s> (%s) %s", name, kind, hint)
	} else if kind != "" {
		return fmt.Sprintf("<%s> (%s)", name, kind)
	} else {
		return fmt.Sprintf("<%s>", name)
	}
}

type ArgVal interface {
	Get() interface{}
	Set(string) error
	String() string
}

type ArgSet struct {
	args      []*Arg
	optional  bool
	remaining *remaining
	hideUsage bool
}

func NewArgSet() *ArgSet {
	return &ArgSet{}
}

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
					return fmt.Errorf("missing arg at index %d", idx+1)
				}
			} else {
				if err := arg.value.Set(input[idx]); err != nil {
					return err // FIXME
				}
			}
		}
		consumed++
	}

	if consumed < inputLen {
		return fmt.Errorf("found %d additional argument(s)", inputLen-consumed)
	}

	return nil
}

func (a *ArgSet) Remaining(p *[]string, name string, minmax Range, usage string) {
	a.RemainingVar((*StringList)(p), name, minmax, usage)
}

func (a *ArgSet) RemainingInts(p *[]int, name string, minmax Range, usage string) {
	a.RemainingVar((*IntList)(p), name, minmax, usage)
}

func (a *ArgSet) RemainingInt64s(p *[]int64, name string, minmax Range, usage string) {
	a.RemainingVar((*Int64List)(p), name, minmax, usage)
}

func (a *ArgSet) RemainingUints(p *[]uint, name string, minmax Range, usage string) {
	a.RemainingVar((*UintList)(p), name, minmax, usage)
}

func (a *ArgSet) RemainingUint64s(p *[]uint64, name string, minmax Range, usage string) {
	a.RemainingVar((*Uint64List)(p), name, minmax, usage)
}

func (a *ArgSet) RemainingFloat64s(p *[]float64, name string, minmax Range, usage string) {
	a.RemainingVar((*Float64List)(p), name, minmax, usage)
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

func (a *ArgSet) RemainingVar(val ArgVal, name string, minmax Range, usage string) {
	a.Var(&remaining{val, minmax}, name, usage)
}

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
