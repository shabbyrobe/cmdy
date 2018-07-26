package cmdy

import (
	"strings"

	"github.com/shabbyrobe/cmdy/args"
)

// FIXME: these are Unix codes, but other operating systems use
// different codes.
//
// On macOS/Linux, it looks like Go uses status code 2 for a panic,
// so it's probably a good idea to avoid that. Discussion will be
// on one or both of these threads:
//
// https://groups.google.com/forum/#!msg/golang-nuts/u9NgKibJsKI/XxCdDihFDAAJ
// https://github.com/golang/go/issues/24284
//
const (
	ExitSuccess  = 0
	ExitFailure  = 1
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

// ErrWithCode allows you to wrap an error in a status code which will be used
// by cli.Fatal() as the exit code.
func ErrWithCode(code int, err error) error {
	if ee, ok := err.(*exitError); ok {
		ee.code = code
		return ee
	}
	return &exitError{err: err, code: code}
}

func (e *exitError) Code() int     { return e.code }
func (e *exitError) Error() string { return e.err.Error() }
func (e *exitError) Cause() error  { return e.err }

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
