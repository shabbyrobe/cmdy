package cmdy

import (
	"fmt"
	"strings"

	"github.com/ArtProcessors/cmdy/arg"
	"github.com/ArtProcessors/cmdy/internal/wrap"
)

// DefaultUsage exists for compatibility with earlier versions.
// It will be removed at a later date.
//
// Deprecated
const DefaultUsage = ""

// Help groups related pieces of information that can be assembled into a
// help message. Only Synopsis is required. As a shorthand, you can use
// 'cmdy.Synopsis("foo")' as an equivalent to 'Help{Synopsis: "foo"}'.
type Help struct {
	// Synopsis is the shortest possible complete description of your command,
	// ideally one sentence.
	//
	// Synopsis is required; cmdy will panic if it is empty.
	Synopsis string

	// Usage is used to specify a more complete help message that will be shown by
	// cli.Fatal() if a UsageError is returned (for example when the '-help' flag is
	// passed).
	//
	// To obtain the full help message from a usage error yourself (outside of
	// cmdy.Fatal), use cmdy.FormatError(err).
	Usage string

	Examples Examples
}

type Examples []Example

func Synopsis(synopsis string) Help {
	return Help{Synopsis: synopsis}
}

/*
Example allows you to describe a testable example for your command.

Command should contain only the flags and args required by the command
returning the example, not any parents, flags from parents, etc. The
examples are tested presuming no parent context. Examples appear in
the following basic format:

	# Desc
	$ Command
	Output

If 'Input' is specified and it is sufficiently short, Command will appear
like so:
    $ echo "Input" | Command

Otherwise it will be collapsed into an ellipsis:
    $ ... | Command

If the command has a parent and the parent is a Group, the subcommand
names used to mount your command will be prepended. For example, if the
full path for your command is this:
	$ foo bar baz cmd

...and the 'Command' portion of the help message for the command associated
with the 'cmd' subcommand is this:
	-flag1 -flag2 arg1

...then the example will appear like so:
	$ foo bar baz cmd -flag1 -flag2 arg1

If you are using nested commands dynamically, rather than via cmdy.Group,
the subcommand path may be missing from the start of the invocation.
*/
type Example struct {
	Desc    string
	Command string

	// Expected status code from command run.
	//	Code == 0   expect success
	//	Code >  0   expect specific failure code
	//	Code <  0   expect any non-zero exit
	//
	// If you are using ExampleParseOnly as your TestMode, you probably want
	// to expect cmdy.ExitUsage if you expect failure.
	Code int

	Input  string
	Output string

	// If true, the output is always hidden.
	HideOutput bool

	// If true, the example is only used for testing and hidden from any help message.
	TestOnly bool

	TestMode ExampleTestMode
}

type ExampleTestMode int

const (
	ExampleParseOnly ExampleTestMode = 0
	ExampleRun       ExampleTestMode = 1
)

func buildHelp(
	cmd Command,
	path CommandPath,
	flagSet *FlagSet,
	argSet *arg.ArgSet,
) (string, error) {
	help := cmd.Help()

	sections := []HelpSection{
		synopsisSection{&help},
		invocationSection{path, flagSet, argSet},
		usageSection{&help},
		flagSection{flagSet},
		argSection{argSet},
		exampleSection{help.Examples, path},
		commandSection{cmd},
	}

	var out strings.Builder
	var lastLen = 0
	var lastSec = len(sections)

	for idx, sec := range sections {
		if err := sec.BuildHelp(&out); err != nil {
			return "", err
		}

		total := out.Len()
		if total > lastLen && idx != lastSec {
			out.WriteByte('\n')
		}
		lastLen = out.Len()
	}

	return out.String(), nil
}

type HelpSection interface {
	BuildHelp(into *strings.Builder) error
}

type synopsisSection struct {
	help *Help
}

func (s synopsisSection) BuildHelp(into *strings.Builder) error {
	synopsis := strings.TrimSpace(s.help.Synopsis)
	if synopsis == "" {
		return nil
	}
	into.WriteString(synopsis)
	into.WriteByte('\n')
	return nil
}

type invocationSection struct {
	path    CommandPath
	flagSet *FlagSet
	argSet  *arg.ArgSet
}

func (i invocationSection) BuildHelp(into *strings.Builder) error {
	into.WriteString("Usage: ")

	for idx, p := range i.path {
		if idx > 0 {
			into.WriteByte(' ')
		}
		into.WriteString(p.Name)
	}

	if i.flagSet != nil {
		into.WriteByte(' ')
		into.WriteString(i.flagSet.Invocation())
	}

	if i.argSet != nil {
		into.WriteByte(' ')
		into.WriteString(i.argSet.Invocation())
	}

	into.WriteByte('\n')

	return nil
}

type usageSection struct {
	help *Help
}

func (us usageSection) BuildHelp(into *strings.Builder) error {
	usage := strings.TrimSpace(us.help.Usage)
	if usage == "" {
		return nil
	}
	into.WriteString(usage)
	into.WriteByte('\n')
	return nil
}

type flagSection struct {
	flagSet *FlagSet
}

func (fs flagSection) BuildHelp(into *strings.Builder) error {
	if fs.flagSet != nil {
		fu := fs.flagSet.Usage()
		if fu != "" {
			into.WriteString("Flags:\n")
			into.WriteString(fu)
		}
	}
	return nil
}

type argSection struct {
	argSet *arg.ArgSet
}

func (as argSection) BuildHelp(into *strings.Builder) error {
	if as.argSet != nil {
		au := as.argSet.Usage()
		if au != "" {
			into.WriteString("Arguments:\n")
			into.WriteString(au)
		}
	}
	return nil
}

type exampleSection struct {
	examples Examples
	path     CommandPath
}

func (es exampleSection) BuildHelp(into *strings.Builder) error {
	if len(es.examples) > 0 {
		pathStr := es.path.Invocation()
		into.WriteString("Examples:\n")
		for idx, e := range es.examples {
			if idx > 0 {
				into.WriteByte('\n')
			}
			es.renderExample(into, &e, pathStr)
		}
	}
	return nil
}

func (es exampleSection) renderExample(into *strings.Builder, e *Example, pathStr string) {
	const maxInHideSize = 20
	const maxOutWidth = wrap.DefaultWrap
	const maxOutLines = 2
	const indent = "  "

	if e.Command == "" || e.TestOnly {
		return
	}

	{ // Desc:
		if e.Desc != "" {
			descWrap := wrap.Wrapper{IndentFirst: true, Indent: indent, Prefix: "# "}
			into.WriteString(descWrap.Wrap(e.Desc))
			into.WriteByte('\n')
		}
	}

	{ // Command:
		cmd := e.Command
		if len(pathStr) > 0 {
			cmd = pathStr + " " + cmd
		}

		if e.Input != "" {
			if len(e.Input) <= maxInHideSize {
				cmd = fmt.Sprintf("echo %q | %s", e.Input, e.Command)
			} else {
				cmd = fmt.Sprintf("... | %s", e.Command)
			}
		}

		cmdWrap := wrap.Wrapper{
			IndentFirst: true,
			Indent:      indent,
			WrapWith:    " \\\n    ",
		}
		into.WriteString(cmdWrap.Wrap("$ " + cmd))
		into.WriteByte('\n')
	}

	{ // Output:
		if !e.HideOutput && e.Output != "" {
			lines := strings.SplitN(strings.TrimSpace(e.Output), "\n", maxOutLines+1)

			// replace unsplit remainder with ellipsis
			if len(lines) > maxOutLines {
				lines[maxOutLines] = "..."
			}

			for idx, line := range lines {
				// ensure individual lines are truncated if they exceed the max width:
				if len(line) >= maxOutWidth {
					line = strings.TrimSpace(line[:maxOutWidth])
					if !strings.HasPrefix(line, "...") {
						line += " ..."
					}
				}
				lines[idx] = indent + line
			}

			into.WriteString(strings.Join(lines, "\n"))
			into.WriteByte('\n')
		}
	}
}

type commandSection struct {
	cmd Command
}

func (cs commandSection) BuildHelp(into *strings.Builder) error {
	hs, ok := cs.cmd.(HelpSection)
	if !ok {
		return nil
	}
	return hs.BuildHelp(into)
}
