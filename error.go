package cmdy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shabbyrobe/cmdy/arg"
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

// QuietExit is an error you can return to prevent cmdy.Fatal() from printing an error
// message on exit, but still cause it to call os.Exit() with the status code it
// represents.
type QuietExit int

func (e QuietExit) Code() int     { return int(e) }
func (e QuietExit) Error() string { return fmt.Sprintf("exit code %d", e) }

// ErrWithCode allows you to tag an arbitrary error with a status code which will be used
// by cmdy.Fatal() as the exit code.
func ErrWithCode(code int, err error) error {
	if ee, ok := err.(*exitError); ok {
		ee.code = code
		return ee
	}
	return &exitError{err: err, code: code}
}

// HelpRequest returns an error that will instruct cmdy.Fatal() to print the full command
// help.
func HelpRequest() error {
	return &usageError{showFullHelp: true}
}

// UsageError wraps an existing error so that cmdy.Fatal() will print the full
// command usage above the error message.
func UsageError(err error) error {
	return &usageError{err: err}
}

// UsageErrorf formats an error message so that cmdy.Fatal() will print the full
// command usage above it.
func UsageErrorf(msg string, args ...interface{}) error {
	return &usageError{err: fmt.Errorf(msg, args...)}
}

func IsUsageError(err error) bool {
	var u *usageError
	return errors.As(err, &u)
}

// ErrCode returns the error code associated with the error if it implements
// cmdy.Error, or ExitInternal if not.
func ErrCode(err error) (code int) {
	if err == nil {
		return ExitSuccess
	}
	e, ok := err.(Error)
	if !ok {
		return ExitInternal
	}
	return e.Code()
}

type exitError struct {
	code int
	err  error
}

func (e *exitError) Code() int     { return e.code }
func (e *exitError) Unwrap() error { return e.err }
func (e *exitError) Error() string { return e.err.Error() }

type usageError struct {
	err          error
	usage        string
	showFullHelp bool
	populated    bool
}

func (u *usageError) Code() int     { return ExitUsage }
func (u *usageError) Unwrap() error { return u.err }

func (u *usageError) Error() string {
	if u.err == nil {
		return "usage error"
	}
	return u.err.Error()
}

func (u *usageError) populate(usage string, path []string, flagSet *FlagSet, argSet *arg.ArgSet, examples Examples) {
	if u.populated {
		return
	}
	u.populated = true

	var out strings.Builder
	out.WriteString(strings.TrimSpace(usage))
	out.WriteByte('\n')

	// FIXME: this stuff feels like it doesn't belong here:

	if flagSet != nil {
		fu := flagSet.Usage()
		if fu != "" {
			if out.Len() > 0 {
				out.WriteByte('\n')
			}
			out.WriteString("Flags:\n")
			out.WriteString(fu)
		}
	}

	if argSet != nil {
		au := argSet.Usage()
		if au != "" {
			if out.Len() > 0 {
				out.WriteByte('\n')
			}
			out.WriteString("Arguments:\n")
			out.WriteString(au)
		}
	}

	if u.showFullHelp && len(examples) > 0 {
		if out.Len() > 0 {
			out.WriteByte('\n')
		}
		out.WriteString("Examples:\n")
		examples.render(&out, path, "  ")
	}

	u.usage = out.String()
}

type errorGroup interface {
	Errors() []error
}
