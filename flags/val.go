package flags

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var splitPattern = regexp.MustCompile(`\s*,\s*`)

// StringList is a flag.Getter which allows you to accumulate multiple
// instances of the same flag into a slice of strings.
//
// If your flag is set up like so:
//	var myList StringList
//	flag.Var(&myList, "yep", "My List")
//
// ...then passing "-yep foo -yep bar" will cause myList to equal
// '[]string{"foo", "bar"}'.
//
// StringList does not make sense to use with an arg.ArgSet; use ArgSet.Remaining()
// instead.
//
type StringList []string

func (s StringList) Get() interface{}  { return []string(s) }
func (s StringList) Strings() []string { return s }

func (s StringList) String() string {
	return strings.Join(s, ",")
}

func (s *StringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

// IntList is a flag.Getter which allows you to accumulate multiple
// instances of the same flag into a slice of ints
//
// If your flag is set up like so:
//	var myList IntList
//	flag.Var(&myList, "yep", "My List")
//
// ...then passing '-yep 1 -yep 2' will cause myList to equal '[]int{1, 2}'.
//
// IntList also supports comma separated values, so passing -yep 1,2,3
// or -yep "1, 2, 3" will cause myList to equal '[]int{1, 2, 3}'.
//
// IntList uses 'strconv.ParseInt(..., 0, 0)' to parse the value; the base is implied by
// the string's prefix: base 2 for "0b", base 8 for "0" or "0o", base 16 for "0x", and
// base 10 otherwise. Underscore characters are permitted per the Go integer literal
// syntax.
//
// IntList does not make sense to use with an arg.ArgSet; use ArgSet.RemainingInts()
// instead.
//
type IntList []int

func (s IntList) Get() interface{} { return []int(s) }
func (s IntList) Ints() []int      { return s }

func (s IntList) String() string {
	if len(s) == 0 {
		return ""
	}
	var out = make([]string, 0, len(s))
	for _, i := range s {
		out = append(out, fmt.Sprintf("%d", i))
	}
	return strings.Join(out, ",")
}

func (s *IntList) Set(v string) error {
	for _, part := range splitPattern.Split(strings.TrimSpace(v), -1) {
		if len(part) == 0 {
			continue
		}
		i, err := strconv.ParseInt(part, 0, 0)
		if err != nil {
			return err
		}
		*s = append(*s, int(i))
	}
	return nil
}

// Int64List is a flag.Getter which allows you to accumulate multiple
// instances of the same flag into a slice of int64s.
//
// Int64List is identical to IntList except for the int size; see the IntList
// documentation for more details.
//
// Int64List does not make sense to use with an arg.ArgSet; use ArgSet.RemainingInt64s()
// instead.
//
type Int64List []int64

func (s Int64List) Get() []int64    { return []int64(s) }
func (s Int64List) Int64s() []int64 { return s }

func (s Int64List) String() string {
	if len(s) == 0 {
		return ""
	}
	var out = make([]string, 0, len(s))
	for _, i := range s {
		out = append(out, fmt.Sprintf("%d", i))
	}
	return strings.Join(out, ",")
}

func (s *Int64List) Set(v string) error {
	for _, part := range splitPattern.Split(strings.TrimSpace(v), -1) {
		if len(part) == 0 {
			continue
		}
		i, err := strconv.ParseInt(part, 0, 64)
		if err != nil {
			return err
		}
		*s = append(*s, i)
	}
	return nil
}

// Float64List is a flag.Getter which allows you to accumulate multiple
// instances of the same flag into a slice of float64s
//
// If your flag is set up like so:
//	var myList Float64List
//	flag.Var(&myList, "yep", "My List")
//
// ...then passing '-yep 1.0 -yep 2.0' will cause myList to equal '[]float64{1.0, 2.0}'.
//
// Float64List also supports comma separated values, so passing -yep 1.0,2.0,3.0
// or -yep "1.0, 2.0, 3.0" will cause myList to equal '[]float64{1.0, 2.0, 3.0}'.
//
// Float64List uses 'strconv.ParseFloat(..., 64)' to parse the value, and thus permits
// all syntax described by the Go documentation as valid.
//
// Float64List does not make sense to use with an arg.ArgSet; use ArgSet.RemainingFloat64s()
// instead.
//
type Float64List []float64

func (s Float64List) Get() interface{}    { return []float64(s) }
func (s Float64List) Float64s() []float64 { return s }

func (s Float64List) String() string {
	if len(s) == 0 {
		return ""
	}
	var out = make([]string, 0, len(s))
	for _, i := range s {
		out = append(out, fmt.Sprintf("%f", i))
	}
	return strings.Join(out, ",")
}

func (s *Float64List) Set(v string) error {
	for _, part := range splitPattern.Split(strings.TrimSpace(v), -1) {
		if len(part) == 0 {
			continue
		}
		i, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return err
		}
		*s = append(*s, i)
	}
	return nil
}
