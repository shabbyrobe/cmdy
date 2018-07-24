package cmdy

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	stderr io.Writer = os.Stderr // Substitutable for testing
)

func Run(ctx context.Context, args []string, b Builder) (rerr error) {
	cmd, err := b()
	if err != nil {
		return err
	}

	defer func() {
		if uerr, ok := rerr.(*usageError); ok {
			uerr.populate(cmd)
		}
	}()

	var (
		flagSet = cmd.Flags()
		argSet  = cmd.Args()
		remArgs = args
	)

	if flagSet != nil {
		if err := flagSet.Parse(args); err != nil {
			if err == flag.ErrHelp {
				return NewUsageError(err)
			} else {
				return err
			}
		}
		remArgs = flagSet.Args()
	}

	if argSet != nil {
		if err := argSet.Parse(remArgs); err != nil {
			return NewUsageError(err)
		}

	} else if len(remArgs) > 0 {
		return fmt.Errorf("expected 0 arguments, found %d", len(remArgs))
	}

	cctx := &commandContext{ctx}
	input := &input{cmd: cmd, rawArgs: args}
	return cmd.Run(cctx, input)
}

func Fatal(err error) {
	msg, code := FormatError(err)
	if msg != "" {
		fmt.Fprintln(stderr, msg)
	}
	os.Exit(code)
}

func FormatError(err error) (msg string, code int) {
	if err == nil {
		return "", 0
	}

	code = ExitDefault

	switch err := err.(type) {
	case *usageError:
		msg = strings.TrimSpace(err.usage)
		if err.err != nil {
			if msg != "" {
				msg += "\n\n"
			}
			msg += "error: " + err.err.Error()
		}
	case Error:
		msg, code = err.Error(), err.Code()
	case errorGroup:
		errs := err.Errors()
		last := len(errs) - 1
		for i, e := range errs {
			msg += "- " + e.Error()
			if i != last {
				msg += "\n"
			}
		}
	default:
		msg = err.Error()
	}

	if code == 0 {
		code = ExitDefault
	}

	return
}
