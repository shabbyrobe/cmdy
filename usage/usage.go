package usage

import (
	"fmt"
	"reflect"
	"strings"
)

const DefaultWrap = 80
const indent = "        "

var indentFlag = indent[:len(indent)-4]

// Usable is a common interface that should support both flag.Flag and arg.Arg
// for the purpose of building a usage statement.
type Usable interface {
	Name() string
	Usage() string
	DefValue() string
	Value() interface{}

	// Describe is used to format the description of this Usable in the usage
	// statement according to the Usable's preferredtemplate. Flags are
	// formatted differently to args, but the inputs are the same.
	Describe(kind string, hint string) string
}

// Hinter allows flag.Var or arg.Var implementations to customise the type and
// hint portion of the usage output, for example:
//	--flag=<kind> (hint)
//
// Empty strings will remove either part from the output:
//	--flag=<kind> (hint)
//	--flag (hint)
//
type Hinter interface {
	Hint() (kind, hint string)
}

func Usage(width int, usables ...Usable) string {
	var out strings.Builder

	for _, usable := range usables {
		usage, kind, hint := unquoteUsage(usable)
		s := "  " + usable.Describe(kind, hint)

		defval := usable.DefValue()
		showDefault := !isZeroValue(usable, defval)

		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += indentFlag
		} else if usage != "" || showDefault {
			s += "\n" + indent
		}

		if showDefault {
			if containsString(reflect.TypeOf(usable.Value())) {
				usage += fmt.Sprintf(" (default: %q)", defval)
			} else {
				usage += fmt.Sprintf(" (default: %v)", defval)
			}
		}

		s += Wrap(usage, indent, width)

		out.WriteString(s)
		out.WriteByte('\n')
	}

	return out.String()
}

func containsDuration(v reflect.Type) bool {
	// Yuck. I can't find an easy way to inspect the underlying type of
	// derived types using reflection. It seems like it may not be possible
	// to get the 'time.Duration' part of 'type durationValue *time.Duration'
	// without resorting to using go/types, which is WAY too heavy for this
	// one job:
	return strings.Contains(strings.ToLower(v.Name()), "duration")
}

func containsString(v reflect.Type) bool {
	switch v.Kind() {
	case reflect.String:
		return true
	case reflect.Ptr:
		return v.Elem().Kind() == reflect.String
	default:
		return false
	}
}

func unquoteUsage(usable Usable) (usage string, name string, hint string) {
	// Look for a back-quoted name, but avoid the strings package.
	usage = usable.Usage()
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return usage, name, hint
				}
			}
			break // Only one back quote; use type name.
		}
	}

	// No explicit name, so use type if we can find one.
	// The reflection on internal type names is unfortunate, but it was the easiest
	// way to duplicate the stdlib's functionality:
	name = "value"
	name, hint = Kind(usable)

	return
}

func Kind(usable Usable) (kind, hint string) {
	vt := reflect.TypeOf(usable.Value())
	for vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	if hv, ok := usable.Value().(Hinter); ok {
		kind, hint = hv.Hint()
	} else {
		kind, hint = kindFromType(vt)
	}

	return kind, hint
}

func kindFromType(vt reflect.Type) (kind, typeHint string) {
	if containsDuration(vt) {
		kind = "duration"
		typeHint = "formats: '1h2s', '-3.4ms', units: h, m, s, ms, us, ns"

	} else {
		switch vt.Kind() {
		case reflect.Bool:
			kind = ""
		case reflect.Float32, reflect.Float64:
			kind = "float"
		case reflect.String:
			kind = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			kind = "int"
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			kind = "uint"
		case reflect.Slice:
			kind, typeHint = kindFromType(vt.Elem())
			// kind += "[]"
		}
	}

	return kind, typeHint
}

// isZeroValue guesses whether the string represents the zero
// value for a flag. It is not accurate but in practice works OK.
func isZeroValue(usable Usable, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(usable.Value())
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}

	zif := z.Interface()
	if zif != nil {
		if zs, ok := zif.(fmt.Stringer); ok && zs.String() == value {
			return true
		}
	}

	switch value {
	case "false", "", "0":
		return true
	}
	return false
}

// Wrap should not be used or relied upon outside cmdy. Changes to the API of
// Wrap will not be considered grounds for a semver bump. Use at your own risk.
func Wrap(str string, indent string, width int) string {
	str = strings.TrimSpace(str)

	if width <= 0 {
		width = DefaultWrap
	}

	out := ""
	var ln int
	for _, line := range strings.Split(str, "\n") {
		for {
			if ln > 0 {
				out += "\n" + indent
			}

			var (
				i, j     int
				c        rune
				breaking bool
				inEsc    bool
			)

			for j, c = range line {
				if i == width {
					breaking = true
					break
				}

				// Don't count ASCII escape sequences towards line width:
				if inEsc {
					if c == 'm' {
						inEsc = false
					}
				} else {
					if c == '\033' {
						inEsc = true
					} else {
						i++
					}
				}
			}

			cur := line[:j]
			idx := strings.LastIndexAny(cur, " -")
			if idx < 0 || !breaking {
				out += line
				break
			} else {
				out += line[:idx]
				line = line[idx+1:]
			}
			ln++
		}
	}

	return out
}
