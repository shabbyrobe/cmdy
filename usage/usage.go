package usage

import (
	"fmt"
	"reflect"
	"strings"
)

const DefaultWrap = 80
const indent = "        "

var indentFlag = indent[:len(indent)-4]

type Usable interface {
	Name() string
	Usage() string
	DefValue() string
	Value() interface{}
	Describe(kind string, hint string) string
}

func Usage(width int, usables ...Usable) string {
	var out strings.Builder

	for _, usable := range usables {
		usage, kind, hint := unquoteUsage(usable)
		s := "  " + usable.Describe(kind, hint)

		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += indentFlag
		} else {
			s += "\n" + indent
		}

		defval := usable.DefValue()
		if !isZeroValue(usable, defval) {
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
	vt := reflect.TypeOf(usable.Value())
	for vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	name, hint = elemType(vt)

	return
}

func elemType(vt reflect.Type) (name, hint string) {
	if containsDuration(vt) {
		name = "duration"
		hint = "formats: '1h2s', '-3.4ms', units: h, m, s, ms, us, ns"
	} else {
		switch vt.Kind() {
		case reflect.Bool:
			name = ""
		case reflect.Float32, reflect.Float64:
			name = "float"
		case reflect.String:
			name = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			name = "int"
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			name = "uint"
		case reflect.Slice:
			name, hint = elemType(vt.Elem())
			// name += "[]"
		}
	}
	return name, hint
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

// Wrap should not be used or relied upon outside cmdy.
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

			var i, j int
			var breaking bool
			for j = range line {
				if i == width {
					breaking = true
					break
				}
				i++
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
