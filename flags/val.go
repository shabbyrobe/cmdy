package flags

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var splitPattern = regexp.MustCompile(`\s*,\s*`)

type StringList []string

func (s StringList) Strings() []string { return s }

func (s *StringList) String() string {
	if s == nil {
		return ""
	}
	return strings.Join(*s, ",")
}

func (s *StringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type IntList []int

func (s IntList) Ints() []int { return s }

func (s *IntList) String() string {
	if s == nil {
		return ""
	}
	var out []string
	for _, i := range *s {
		out = append(out, fmt.Sprintf("%d", i))
	}
	return strings.Join(out, ",")
}

func (s *IntList) Set(v string) error {
	for _, part := range splitPattern.Split(v, -1) {
		if len(part) == 0 {
			continue
		}
		i, err := strconv.ParseInt(part, 10, 0)
		if err != nil {
			return err
		}
		*s = append(*s, int(i))
	}
	return nil
}

type Int64List []int64

func (s Int64List) Int64s() []int64 { return s }

func (s *Int64List) String() string {
	if s == nil {
		return ""
	}
	var out []string
	for _, i := range *s {
		out = append(out, fmt.Sprintf("%d", i))
	}
	return strings.Join(out, ",")
}

func (s *Int64List) Set(v string) error {
	for _, part := range splitPattern.Split(v, -1) {
		if len(part) == 0 {
			continue
		}
		i, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return err
		}
		*s = append(*s, i)
	}
	return nil
}
