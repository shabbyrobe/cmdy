package cmdy

import (
	"testing"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

func TestMatcher(t *testing.T) {
	type tc struct {
		min      int
		in       string
		expected string
		options  []string
	}

	for _, c := range []tc{
		{min: 3, in: "foo", expected: "foo", options: []string{"foo", "food"}},
		{min: 3, in: "food", expected: "food", options: []string{"foo", "food"}},
		{min: 3, in: "food", expected: "foods", options: []string{"foo", "foods"}},
		{min: 4, in: "food", expected: "foods", options: []string{"foo", "foods"}},
		{min: 5, in: "food", expected: "", options: []string{"foo", "foods"}},
		{min: 3, in: "food", expected: "", options: []string{}},
	} {
		t.Run("", func(t *testing.T) {
			tt := assert.WrapTB(t)
			bldrs := Builders{}
			for _, name := range c.options {
				bldrs[name] = func() Command { return &testCmd{} }
			}
			grp := NewGroup("g", bldrs, GroupPrefixMatcher(c.min))
			_, name, err := grp.Builder(c.in)
			tt.MustOK(err)
			tt.MustEqual(c.expected, name)
		})
	}
}
