package cmdstr

import "errors"

type stateType int

const (
	stateStart stateType = iota + 1
	stateArg
	stateEsc
)

type argType int

const (
	argNone argType = iota + 1
	argSingle
	argDouble
	argUnquot
)

var (
	ErrUnterminatedEscape = errors.New("cmdstr: unterminated escape")
	ErrUnterminatedQuote  = errors.New("cmdstr: unterminated quote")
)

// SplitString implements shell-esque splitting rules, similar to Python's shlex.
//
// The input string is stripped of leading and trailing whitespace. Whitespace is defined
// as '\r', '\n', '\t' and ' '.
//
// Strings within cmd can be single-quoted, double-quoted or unquoted. Empty strings are
// allowed.
//
// An unquoted string ends at the first occurence of whitespace or the end of cmd.
//
// Quote characters are not recognized within words; `Do"Not"Separate` is parsed as the
// single word `Do"Not"Separate`.
//
func SplitString(cmd string) ([]string, error) {
	var (
		state   = stateStart
		arg     = argNone
		current = ""
		l       = len(cmd)
		out     []string
	)

	var wsp = [256]bool{'\r': true, '\n': true, '\t': true, ' ': true}

	for i := 0; i < l; i++ {
	retry:
		c := cmd[i]
		switch state {
		case stateStart:
			if wsp[c] {
				continue
			}
			switch c {
			case '"':
				state, arg = stateArg, argDouble
			case '\'':
				state, arg = stateArg, argSingle
			default:
				state, arg = stateArg, argUnquot
				goto retry
			}

		case stateArg:
			if c == '\\' {
				state = stateEsc

			} else if len(current) > 0 &&
				(arg == argSingle && c == '\'') ||
				(arg == argDouble && c == '"') ||
				(arg == argUnquot && wsp[c]) {

				current, out = "", append(out, current)
				state, arg = stateStart, argNone

			} else {
				current += string(c)
			}

		case stateEsc:
			esc := c == '\n' ||
				(arg == argSingle && (c == '\'')) ||
				(arg == argDouble && (c == '"' || c == '\\')) ||
				(arg == argUnquot && (c == '\'' || c == '"' || c == '\\'))

			if !esc {
				current += "\\"
			}

			current += string(c)
			state = stateArg
		}
	}

	if state == stateEsc {
		return nil, ErrUnterminatedEscape
	} else if state == stateArg && arg != argUnquot {
		return nil, ErrUnterminatedQuote
	}

	if len(current) > 0 {
		out = append(out, current)
	}

	return out, nil
}
