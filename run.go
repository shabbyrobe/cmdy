package cmdy

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/shabbyrobe/cmdy/args"
)

var (
	defaultRunner *Runner
)

func DefaultRunner() *Runner {
	if defaultRunner == nil {
		defaultRunner = NewStandardRunner()
	}
	return defaultRunner
}

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewStandardRunner() *Runner {
	return &Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (r *Runner) Run(ctx context.Context, name string, args []string, b Builder) (rerr error) {
	cmd, err := b()
	if err != nil {
		return err
	}

	var (
		flagSet = cmd.Flags()
		argSet  = cmd.Args()
		remArgs = args
	)

	cctx, ok := ctx.(*commandContext)
	if !ok {
		cctx = &commandContext{
			Context: ctx,
			cmd:     cmd,
			rawArgs: args,
			runner:  r,
		}
	}

	cctx.Push(name, cmd)
	defer cctx.Pop()

	defer func() {
		if uerr, ok := rerr.(*usageError); ok {
			path := CommandPath(cctx)

			usageTpl, err := r.usageTpl(cmd, path, flagSet, argSet)
			if err != nil {
				panic(err)
			}

			var buf bytes.Buffer
			if err := usageTpl.Execute(&buf, nil); err != nil {
				panic(err)
			}

			uerr.populate(buf.String(), flagSet, argSet)
		}
	}()

	if flagSet == nil {
		flagSet = NewFlagSet() // handles --help by default, save ourselves the trouble
	}

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			err = nil // suppress "flag: help requested"
		}

		// As at Go 1.11, the only error returned by the flag package that we
		// might not consider a usage error is the one where you define your
		// flag with a '-' in the name, but there's no way to identify it that
		// doesn't involve string matching.
		return NewUsageError(err)
	}
	remArgs = flagSet.Args()

	if argSet != nil {
		if err := argSet.Parse(remArgs); err != nil {
			return NewUsageError(err)
		}

	} else if len(remArgs) > 0 {
		return NewUsageError(fmt.Errorf("expected 0 arguments, found %d", len(remArgs)))
	}

	return cmd.Run(cctx)
}

func (r *Runner) Fatal(err error) {
	msg, code := FormatError(err)
	if msg != "" {
		if _, err := io.WriteString(r.Stderr, msg); err != nil {
			panic(err)
		}
		if _, err := r.Stderr.Write([]byte{'\n'}); err != nil {
			panic(err)
		}
	}
	os.Exit(code)
}

func (r *Runner) usageTpl(cmd Command, path []string, flagSet *FlagSet, argSet *args.ArgSet) (tpl *template.Template, rerr error) {
	// Update the documentation for the Usage interface if you add new functions
	// to this map:
	fns := template.FuncMap{
		"Synopsis": func() string {
			return cmd.Synopsis()
		},
		"Invocation": func() string {
			out := strings.Join(path, " ")
			if flagSet != nil {
				out += " "
				out += flagSet.Invocation()
			}
			if argSet != nil {
				out += " "
				out += argSet.Invocation()
			}
			return out
		},
		"CommandFull": func() string {
			return strings.Join(path, " ")
		},
		"Command": func() string {
			if len(path) > 0 {
				return path[len(path)-1]
			}
			return ""
		},
	}
	tpl = template.New("usage").Funcs(fns)

	var usageRaw string
	if ucmd, ok := cmd.(Usage); ok {
		usageRaw = ucmd.Usage()
	}
	if usageRaw == "" {
		usageRaw = defaultUsage
	}

	var err error
	tpl, err = tpl.Parse(usageRaw)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func Run(ctx context.Context, args []string, b Builder) (rerr error) {
	name := ProgName()
	return DefaultRunner().Run(ctx, name, args, b)
}

func Fatal(err error) {
	DefaultRunner().Fatal(err)
}

func FormatError(err error) (msg string, code int) {
	if err == nil {
		return "", ExitSuccess
	}

	switch err := err.(type) {
	case *usageError:
		// usageError.usage is lazily populated from a Go text/template in
		// Runner.Run() before it is returned:
		msg = strings.TrimSpace(err.usage)
		code = err.Code()

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
		code = ExitFailure
	}

	return
}

// ProgName attempts to guess the program name from the first argument in os.Args.
func ProgName() string {
	if len(os.Args) < 1 {
		return ""
	}
	return baseName(os.Args[0])
}

// baseName is a cut-down remix of filepath.Base that saves us a dependency
// and skips use-cases that we don't need to worry about, like windows volume
// names, etc, because we are only using it to grab the program name.
func baseName(path string) string {
	// Strip trailing slashes.
	for len(path) > 0 && os.IsPathSeparator(path[len(path)-1]) {
		path = path[0 : len(path)-1]
	}
	// Find the last element
	i := len(path) - 1
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	if i >= 0 {
		path = path[i+1:]
	}
	return path
}

const defaultUsage = `
{{if Synopsis -}}
{{Synopsis}}

{{end -}}

Usage: {{Invocation}}
`
