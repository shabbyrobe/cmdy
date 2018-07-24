package flags

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var splitPattern = regexp.MustCompile(`\s*,\s*`)

type StringList []string

func (s *StringList) String() string {
	if s == nil {
		return ""
	}
	return strings.Join(*s, ",")
}

func (s StringList) Strings() []string {
	out := make([]string, len(s))
	copy(out, s)
	return out
}

func (s *StringList) Set(v string) error {
	*s = append(*s, v)
	return nil

}

type IntList []int

var _ flag.Value = &IntList{}

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

func (s IntList) Ints() []int {
	out := make([]int, len(s))
	copy(out, s)
	return out
}

func (s *IntList) Set(v string) error {
	for _, part := range splitPattern.Split(v, -1) {
		if len(part) == 0 {
			continue
		}
		i, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return err
		}
		*s = append(*s, int(i))
	}
	return nil
}
