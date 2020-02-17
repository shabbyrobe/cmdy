package arg

import (
	"strconv"
	"strings"
)

type stringVar string

func (h stringVar) String() string            { return string(h) }
func (h stringVar) Hint() (kind, hint string) { return "", "" }
func (h *stringVar) Set(s string) error       { *h = stringVar(s); return nil }

type stringSliceVar []string

func (h stringSliceVar) String() string            { return strings.Join([]string(h), ",") }
func (h stringSliceVar) Hint() (kind, hint string) { return "", "" }
func (h *stringSliceVar) Set(s string) error       { *h = append(*h, s); return nil }

type intSliceVar []int

func (h intSliceVar) String() string {
	out := ""
	for idx, v := range h {
		if idx > 0 {
			out += ","
		}
		out += strconv.FormatInt(int64(v), 10)
	}
	return out
}

func (h intSliceVar) Hint() (kind, hint string) { return "", "" }

func (h *intSliceVar) Set(s string) error {
	i, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return err
	}
	*h = append(*h, int(i))
	return nil
}
