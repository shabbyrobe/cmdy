package cmdy

import (
	"strings"
	"text/template"

	"github.com/shabbyrobe/cmdy/arg"
)

/*
DefaultUsage is used to generate your usage string when your Command does
not implement cmdy.Usage.

You can use cmdy.Usage to base your own templates on it if you don't
want to repeat the basics over and over:

	const myCommandUsage = `
	Extra stuff about my command that will be stuck on the end.
	Etc etc etc.
	`
	func (c *myCommand) Help() cmdy.Help {
		return cmdy.Help{
			Synopsis: "my command does stuff",
			Usage:    cmdy.Usage(myCommandUsage),
		}
	}
*/
const DefaultUsage = `
{{if Synopsis -}}
{{Synopsis}}

{{end -}}

Usage: {{Invocation}}
`

// Usage constructs a usage string based on DefaultUsage.
//
// It is not necessary to use this, you can construct your own usage from scratch. See
// DefaultUsage for more info.
//
func Usage(usage string) string {
	out := DefaultUsage
	if usage != "" {
		out += strings.TrimSpace(usage)
		out += "\n"
	}
	return out
}

func buildUsageTpl(help Help, showFullHelp bool, path []string, flagSet *FlagSet, argSet *arg.ArgSet) (tpl *template.Template, rerr error) {
	// Update the documentation for Help.Usage if you add new functions
	// to this map:
	fns := template.FuncMap{
		"Synopsis": func() string {
			return help.Synopsis
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
			return showFullHelp
		},
	}

	tpl = template.New("usage").Funcs(fns)

	var usageRaw = help.Usage
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
