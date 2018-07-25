package cmdy

import (
	"strings"

	"github.com/shabbyrobe/cmdy/args"
)

const (
	ExitDefault  = 2
	ExitUsage    = 127
	ExitInternal = 255
)

type Error interface {
	Code() int
	error
}

type exitError struct {
	code int
	err  error
}

func (e exitError) Code() int     { return e.code }
func (e exitError) Error() string { return e.err.Error() }
func (e exitError) Cause() error  { return e.err }

type usageError struct {
	err       error
	usage     string
	populated bool
}

func (u *usageError) Code() int     { return ExitUsage }
func (u *usageError) Error() string { return u.err.Error() }
func (u *usageError) Cause() error  { return u.err }

func (u *usageError) populate(usage string, flagSet *FlagSet, argSet *args.ArgSet) {
	if u.populated {
		return
	}
	u.populated = true

	out := strings.TrimSpace(usage) + "\n"

	if flagSet != nil {
		if out != "" {
			out += "\n"
		}
		fu := flagSet.Usage()
		if fu != "" {
			out += "Flags:\n" + fu
		}
	}
	if argSet != nil {
		if out != "" {
			out += "\n"
		}
		au := argSet.Usage()
		if au != "" {
			out += "Arguments:\n" + au + "\n"
		}
	}
	u.usage = out
}

func NewUsageError(err error) error {
	return &usageError{err: err}
}

type errorGroup interface {
	Errors() []error
}
