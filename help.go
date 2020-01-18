package cmdy

import (
	"fmt"
	"strings"

	"github.com/shabbyrobe/cmdy/internal/wrap"
)

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
	// If Usage is an empty string, cmdy.DefaultUsage is used.
	//
	// The string in Usage is parsed by the text/template package
	// (https://golang.org/pkg/text/template/). The template makes the following functions
	// available:
	//
	//     {{Invocation}}
	//         Full invocation string for the command, i.e.
	//         'cmd sub subsub [options] <args...>'.
	//         This invocation does not include parent command flags.
	//
	//     {{Synopsis}}
	//         Command.Synopsis()
	//
	//     {{CommandFull}}
	//         Full command name including all parent commands, i.e. 'cmd sub subsub'.
	//
	//     {{Command}}
	//         Current command name, not including parent command names. i.e. for
	//         command 'cmd sub subsub', only 'subsub' is returned.
	//
	//     {{if ShowFullHelp}}...{{end}}
	//         Help section contained inside the '...' should only be shown if the
	//         command's '--help' was requested, not if the command's usage is to
	//         be shown.
	//
	//
	// Your Command instance is used as the 'data' argument to Template.Execute(),
	// so any exported fields from your command can be used in the template like
	// so: "{{.MyCommandField}}".
	//
	// If a Command intends cmdy to print the usage in response to an error,
	// cmdy.UsageError or cmdy.UsageErrorf should be returned from Command.Run().
	//
	// To obtain an actual usage string from a usage error, use cmdy.Format(err).
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

func (e Example) render(into *strings.Builder, path []string, indent string) {
	const maxInWidth = wrap.DefaultWrap
	const maxOutWidth = wrap.DefaultWrap
	const maxOutLines = 2

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
		if len(path) > 0 {
			cmd = strings.Join(path, " ") + " " + cmd
		}

		if e.Input != "" {
			const fudge = 10
			if len(e.Input)+len(e.Command)+fudge < maxInWidth { // near enough is good enough
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

func (ex Examples) render(into *strings.Builder, path []string, indent string) {
	if len(ex) == 0 {
		return
	}

	for idx, e := range ex {
		if idx > 0 {
			into.WriteByte('\n')
		}
		e.render(into, path, "  ")
	}
}
