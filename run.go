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

// DefaultRunner is the global runner used by Run() and Fatal().
//
// It is intended to be used once from your main() function and is not safe for
// concurrent use. More sophisticated use cases can be supported by creating
// your own Runner directly.
func DefaultRunner() *Runner {
	if defaultRunner == nil {
		defaultRunner = NewStandardRunner()
	}
	return defaultRunner
}

// Reset is here just for testing purposes.
func Reset() {
	defaultRunner = nil
}

// Runner builds and runs your command.
//
// Runner provides access to standard input and output streams to cmdy.Command.
// Commands should access these streams via Runner rather than via os.Stdin, etc.
//
// This is not strictly required, and some situations may necessitate using the
// os streams directly, but using os streams directly without a good reason
// limits your command's testability.
//
// See NewBufferedRunner(), NewStandardRunner()
type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewStandardRunner returns a Runner configured to use os.Stdin, os.Stdout and
// os.Stderr.
func NewStandardRunner() *Runner {
	return &Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run builds and runs your command.
//
// If using Run() in your main() function, the returned error should be passed
// to Runner.Fatal(), not log.Fatal(), if you want nice errors and usage printed.
//
func (r *Runner) Run(ctx context.Context, name string, args []string, b Builder) (rerr error) {
	cmd, init := b()
	if init != nil {
		if err := init(); err != nil {
			return err
		}
	}

	var (
		flagSet = cmd.Flags()
		argSet  = cmd.Args()
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

			usageTpl, err := r.usageTpl(cmd, uerr.help, path, flagSet, argSet)
			if err != nil {
				panic(err)
			}

			var buf bytes.Buffer
			if err := usageTpl.Execute(&buf, cmd); err != nil {
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
			// suppress "flag: help requested"
			return NewHelpRequest()
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

// Fatal prints an error formatted for the end user, then calls os.Exit with
// the exit code detected in err.
//
// Calls to Fatal() will prevent any defer calls from running. See cmdy.Fatal()
// for a demonstration of the recommended usage pattern for dealing with Fatal
// errors.
//
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

func (r *Runner) usageTpl(cmd Command, fullHelp bool, path []string, flagSet *FlagSet, argSet *args.ArgSet) (tpl *template.Template, rerr error) {
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
		"ShowFullHelp": func() bool {
			return fullHelp
		},
	}

	tpl = template.New("usage").Funcs(fns)

	var usageRaw string
	if ucmd, ok := cmd.(Usage); ok {
		usageRaw = ucmd.Usage()
	}
	if usageRaw == "" {
		usageRaw = DefaultUsage
	}

	var err error
	tpl, err = tpl.Parse(usageRaw)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

// Run the command built by Builder b using the DefaultRunner, passing in the
// provided args.
//
// The args should not include the program; if using os.Args, you should
// pass 'os.Args[1:]'.
//
// The context provided should be your own master context; this allows global
// shutdown or cancellation to be propagated (provided your command blocks on
// APIs that support contexts). If no context is available, use
// context.Background().
//
func Run(ctx context.Context, args []string, b Builder) (rerr error) {
	name := ProgName()
	return DefaultRunner().Run(ctx, name, args, b)
}

// Fatal prints an error formatted for the end user using the global
// DefaultRunner, then calls os.Exit with the exit code detected in err.
//
// Calls to Fatal() will prevent any defer calls from running. This pattern is
// strongly recommended instead of a straight main() function:
//
//	func main() {
//		if err := run(); err != nil {
//			cmdy.Fatal(err)
//		}
//	}
//
//	func run() error {
//		// your command in here
//		return nil
//	}
//
func Fatal(err error) {
	DefaultRunner().Fatal(err)
}

func FormatError(err error) (msg string, code int) {
	if err == nil {
		return "", ExitSuccess
	}

	switch err := err.(type) {
	case QuietExit:
		// If we don't return here, a '0' code will be interpreted as an
		// ExitFailure. In the case of QuietExit, it's a little bit less
		// natural to assume '0' means we want a non-zero exit status even
		// though we are technically returning an error.
		return "", err.Code()

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

// NewBufferedRunner returns a Runner that wires Stdin, Stdout and Stderr up to
// bytes.Buffer instances.
func NewBufferedRunner() *BufferedRunner {
	br := &BufferedRunner{}
	br.Runner = Runner{
		Stdin:  &br.StdinBuffer,
		Stdout: &br.StdoutBuffer,
		Stderr: &br.StderrBuffer,
	}
	return br
}

type BufferedRunner struct {
	Runner
	StdinBuffer  bytes.Buffer
	StdoutBuffer bytes.Buffer
	StderrBuffer bytes.Buffer
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

// DefaultUsage is used to generate your usage string when your Command does
// not implement cmdy.Usage. You can prepend it to your own usage templates
// if you want to add to it:
//
//	const myCommandUsage = cmdy.DefaultUsage + "\n"+ `
//	Extra stuff about my command that will be stuck on the end.
//	Etc etc etc.
//	`
//
//	func (c *myCommand) Usage() string { return myCommandUsage }
//
const DefaultUsage = `
{{if Synopsis -}}
{{Synopsis}}

{{end -}}

Usage: {{Invocation}}
`
