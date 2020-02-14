package cmdy

import (
	"bytes"
	"context"
	"flag"
	"io"
	"os"

	"github.com/ArtProcessors/cmdy/arg"
)

var (
	defaultRunner = NewStandardRunner()
)

// DefaultRunner is the global runner used by Run() and Fatal().
func DefaultRunner() *Runner {
	return defaultRunner
}

// Reset is here just for testing purposes.
func Reset() {
	defaultRunner = NewStandardRunner()
}

// Runner builds and runs your command.
//
// Runner provides access to standard input and output streams to cmdy.Command.
// Commands should access these streams via Runner rather than via os.Stdin, etc.
//
// It is necessary to use Runner, and some situations may necessitate using the
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
// If calling from your main() function, the returned error should be passed to the
// Runner's Fatal() method, not log.Fatal(), if you want nice errors and usage printed.
//
// 'name' should be the top-level name of your program. You can use ProgName() to guess
// the program's name from os.Args[0].
//
func (r *Runner) Run(ctx context.Context, name string, args []string, b Builder) (rerr error) {
	var (
		cmd     = b()
		flagSet *FlagSet
		argSet  *arg.ArgSet
	)

	// FIXME: see if we can remove this; only a test depends on it at the moment:
	if acmd, ok := cmd.(interface{ Args() *arg.ArgSet }); ok {
		argSet = acmd.Args()
	}
	if fcmd, ok := cmd.(interface{ Flags() *FlagSet }); ok {
		flagSet = fcmd.Flags()
	}

	if argSet == nil {
		argSet = arg.NewArgSet()
	}
	if flagSet == nil {
		flagSet = NewFlagSet()
	}
	cmd.Configure(flagSet, argSet)

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
		// when a nested command raises a usage error, we only want the topmost
		// call to Run() to handle filling in the usage, but this defer() block
		// gets called for all commands on the stack. checking uerr.usage
		// prevents all parents of the command that raised the usageError from
		// clobbering the already-built usage.
		if uerr, ok := rerr.(*usageError); ok && uerr.usage == "" {
			path := cctx.Stack()
			help, err := buildHelp(cmd, path, flagSet, argSet)
			if err != nil {
				panic(err)
			}
			uerr.usage = help
		}
	}()

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			// suppress "flag: help requested"
			return HelpRequest()
		}

		// As at Go 1.11, the only error returned by the flag package that we
		// might not consider a usage error is the one where you define your
		// flag with a '-' in the name, but there's no way to identify it that
		// doesn't involve string matching.
		return UsageError(err)
	}

	remArgs := flagSet.Args()
	if err := argSet.Parse(remArgs); err != nil {
		return UsageError(err)
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
