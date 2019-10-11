package cmdstr

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type stateType int

func (s stateType) String() string {
	switch s {
	case stateStart:
		return "start"
	case stateSingle:
		return "squot"
	case stateDouble:
		return "dquot"
	case stateBare:
		return "arg"
	case stateBareSingle:
		return "squot"
	case stateBareDouble:
		return "dquot"
	case stateEsc:
		return "esc"
	default:
		panic("unknown state type")
	}
}

const (
	stateStart stateType = iota
	stateSingle
	stateDouble
	stateBare
	stateBareSingle
	stateBareDouble
	stateEsc
)

// Pare implements shell-esque string-splitting rules, similar to Python's shlex.
//
// Parsing ends when all arguments are successfully consumed, when invalid syntax
// is encountered, or when one of the characters in 'endset' is encountered while
// scanning for the start of the next argument.
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
// Enclosing characters in single-quotes ( '' ) shall preserve the literal value of each
// character within the single-quotes. A single-quote cannot occur within single-quotes.
//
// Enclosing characters in double-quotes ( "" ) shall preserve the literal value of all
// characters within the double-quotes, with the exception of <backslash>, explained
// below.
//
// The <backslash> shall retain its special meaning as an escape character when
// present in the following two-byte sequences, producing the output on the right
// side:
//
//    0x5c 0x6e (\n)  --> 0x0a
//    0x5c 0x0a (\\n) --> <nothing>
//    0x5c 0x22 (\")  --> 0x22
//    0x5c 0x5c (\\)  --> 0x5c
//
// All other backslashes are copied literally into the output:
//
//    0x5c 0x61 (\a)  --> 0x5c 0x61
//
func Parse(cmd []byte, endset string) (out []string, n int, rerr error) {
	var (
		cur      string
		sz       = len(cmd)
		wsp      = [256]bool{'\r': true, '\n': true, '\t': true, ' ': true}
		esc      = [256]byte{'n': '\n', '"': '"', '\\': '\\'}
		end      = [256]bool{}
		stack    = [4]stateType{}
		stackPos = 0
		start    = 0
		i        = 0
	)

	for i = 0; i < len(endset); i++ {
		end[endset[i]] = true
	}

	for i = 0; i < sz; i++ {
		c := cmd[i]

	retry:
		switch stack[stackPos] {
		case stateStart:
			if wsp[c] {
				continue
			} else if end[c] {
				sz = i // break!
				goto end
			}

			switch c {
			case '"':
				stackPos++
				stack[stackPos], start = stateDouble, i+1
			case '\'':
				stackPos++
				stack[stackPos], start = stateSingle, i+1
			default:
				stackPos++
				stack[stackPos], start = stateBare, i
				goto retry
			}

		case stateDouble:
			if c == '\\' {
				stackPos++
				stack[stackPos] = stateEsc
			} else if c == '"' {
				cur += string(cmd[start:i])
				out, cur = append(out, cur), ""
				stackPos--
			}

		case stateSingle:
			if c == '\\' {
				stackPos++
				stack[stackPos] = stateEsc
			} else if c == '\'' {
				cur += string(cmd[start:i])
				out, cur = append(out, cur), ""
				stackPos--
			}

		case stateBare:
			if c == '\\' {
				stackPos++
				stack[stackPos] = stateEsc
			} else if c == '\'' {
				stackPos++
				stack[stackPos] = stateBareSingle
			} else if c == '"' {
				stackPos++
				stack[stackPos] = stateBareDouble
			} else if wsp[c] {
				cur += string(cmd[start:i])
				out, cur = append(out, cur), ""
				stackPos--
			}

		case stateBareDouble:
			if c == '\\' {
				stackPos++
				stack[stackPos] = stateEsc
			} else if c == '"' {
				cur += string(cmd[start:i])
				start = i
				stackPos--
			}

		case stateBareSingle:
			if c == '\\' {
				stackPos++
				stack[stackPos] = stateEsc
			} else if c == '\'' {
				cur += string(cmd[start:i])
				start = i
				stackPos--
			}

		case stateEsc:
			if c == '\n' {
				// Escaped newlines consume the newline
				cur += string(cmd[start : i-1])
				start = i + 1

			} else if esc[c] != 0 {
				cur += string(cmd[start:i-1]) + string(esc[c])
				start = i + 1

			} else {
				// backslash is literal, do nothing
			}
			stackPos--

		default:
			panic("unknown state")
		}
	}

end:
	if stack[stackPos] == stateBare {
		cur += string(cmd[start:sz])
		out, cur = append(out, cur), ""

	} else if stackPos > 0 {
		stackStr := ""
		for pos := 1; pos <= stackPos; pos++ {
			if pos > 1 {
				stackStr += " > "
			}
			stackStr += stack[pos].String()
		}
		return nil, i, fmt.Errorf("cmdstr: %w, stack: %v", ErrIncompleteCommand, stackStr)
	}

	return out, i, nil
}

func SplitString(cmd string) ([]string, error) {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&cmd))
	cmdBytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
	out, _, err := Parse(cmdBytes, "")
	return out, err
}

var (
	ErrIncompleteCommand = errors.New("incomplete command")
)
