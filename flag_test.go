package cmdy

import (
	"strings"
	"testing"
	"time"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

const expectedUsage = `
  -b    Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. Integer elit ac
        efficitur sodales, nibh felis pulvinar neque, eu pellentesque odio risus sed
        risus. Nulla ac sem ex. Suspendisse in orci pellentesque, posuere massa nec
        (default: true)
  -dv=<duration> (formats: '1h2s', '-3.4ms', units: h, m, s, ms, us, ns)
        Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. (default: 1s)
  -iv=<int>
        Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. (default: 2)
  -pants
        Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. Integer
        tristique, elit ac efficitur sodales, nibh felis (default: true)
  -str=<string>
        Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. (default: "yep")
`

func TestFlagUsage(t *testing.T) {
	tt := assert.WrapTB(t)

	var pants bool
	var str string
	var iv int
	var dv time.Duration

	fs := NewFlagSet()
	fs.BoolVar(&pants, "pants", true, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. "+
		"Integer tristique, elit ac efficitur sodales, nibh felis")

	fs.StringVar(&str, "str", "yep", ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis.")

	fs.IntVar(&iv, "iv", 2, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis.")

	fs.DurationVar(&dv, "dv", 1*time.Second, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis.")

	fs.BoolVar(&pants, "b", true, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. "+
		"Integer elit ac efficitur sodales, nibh felis "+
		"pulvinar neque, eu pellentesque odio risus sed risus. Nulla "+
		"ac sem ex. Suspendisse in orci pellentesque, posuere massa nec")

	// FIXME: brittle test, but adequate for now.
	tt.MustEqual(expectedUsage, "\n"+fs.Usage())
}

func TestFlagDoubleDash(t *testing.T) {
	tt := assert.WrapTB(t)

	FlagDoubleDash = true
	defer func() {
		FlagDoubleDash = false
	}()

	var a, b, a2, b2 bool
	var s, u, s2, u2 string

	fs := NewFlagSet()
	fs.BoolVar(&a, "a", true, "")
	fs.BoolVar(&b, "b", false, "")
	fs.BoolVar(&a2, "a2", false, "")
	fs.BoolVar(&b2, "b2", true, "")

	fs.StringVar(&s, "s", "", "")
	fs.StringVar(&u, "u", "foo", "")
	fs.StringVar(&s2, "s2", "", "")
	fs.StringVar(&u2, "u2", "foo", "")

	usage := fs.Usage()
	tt.MustAssert(strings.Contains(usage, " -a "))
	tt.MustAssert(strings.Contains(usage, " --a2\n"))
	tt.MustAssert(strings.Contains(usage, " -b "))
	tt.MustAssert(strings.Contains(usage, " --b2\n"))
	tt.MustAssert(strings.Contains(usage, " -s=<string>\n"))
	tt.MustAssert(strings.Contains(usage, " --s2=<string>\n"))
	tt.MustAssert(strings.Contains(usage, " -u=<string>\n"))
	tt.MustAssert(strings.Contains(usage, " --u2=<string>\n"))
}

func TestFlagUsageCollapsing(t *testing.T) {
	// Flags without descriptions should not be followed by a blank line
	tt := assert.WrapTB(t)

	var foo, bar int
	fs := NewFlagSet()
	fs.IntVar(&foo, "foo", 0, "")
	fs.IntVar(&bar, "bar", 0, "")

	expected := "" +
		"  -bar=<int>\n" +
		"  -foo=<int>\n"

	tt.MustEqual(expected, fs.Usage())
}

type (
	hintOnlyVar string
	kindOnlyVar string
	hintBothVar string
	hintNoneVar string
)

func (h hintOnlyVar) String() string            { return string(h) }
func (h hintOnlyVar) Hint() (kind, hint string) { return "", "hint" }
func (h *hintOnlyVar) Set(s string) error       { *h = hintOnlyVar(s); return nil }

func (h kindOnlyVar) String() string            { return string(h) }
func (h kindOnlyVar) Hint() (kind, hint string) { return "kind", "" }
func (h *kindOnlyVar) Set(s string) error       { *h = kindOnlyVar(s); return nil }

func (h hintBothVar) String() string            { return string(h) }
func (h hintBothVar) Hint() (kind, hint string) { return "kind", "hint" }
func (h *hintBothVar) Set(s string) error       { *h = hintBothVar(s); return nil }

func (h hintNoneVar) String() string            { return string(h) }
func (h hintNoneVar) Hint() (kind, hint string) { return "", "" }
func (h *hintNoneVar) Set(s string) error       { *h = hintNoneVar(s); return nil }

const expectedHintableUsage = `
  -hintboth=<kind> (hint)
        hint both
  -hintnone
        hint none
  -hintonly (hint)
        hint only
  -kindonly=<kind>
        kind only
`

func TestFlagVarHintable(t *testing.T) {
	tt := assert.WrapTB(t)

	var hintOnly hintOnlyVar
	var kindOnly kindOnlyVar
	var hintBoth hintBothVar
	var hintNone hintNoneVar

	fs := NewFlagSet()
	fs.Var(&hintOnly, "hintonly", "hint only")
	fs.Var(&kindOnly, "kindonly", "kind only")
	fs.Var(&hintBoth, "hintboth", "hint both")
	fs.Var(&hintNone, "hintnone", "hint none")

	// FIXME: brittle test, but adequate for now.
	tt.MustEqual(expectedHintableUsage, "\n"+fs.Usage())
}
