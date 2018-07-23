package flags

import "strings"

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
