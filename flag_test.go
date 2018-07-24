package cmdy

import (
	"testing"
	"time"

	"github.com/shabbyrobe/golib/assert"
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
