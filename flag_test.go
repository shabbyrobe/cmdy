package cmdy

import (
	"fmt"
	"testing"
	"time"

	"github.com/shabbyrobe/golib/assert"
)

func TestFlagUsage(t *testing.T) {
	tt := assert.WrapTB(t)
	_ = tt

	var pants bool
	var str string
	var iv int
	var dv time.Duration

	fs := NewFlagSet()
	fs.BoolVar(&pants, "pants", true, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. "+
		"Integer tristique, elit ac efficitur sodales, nibh felis "+
		"pulvinar neque, eu pellentesque odio risus sed risus. Nulla "+
		"ac sem ex. Suspendisse in orci pellentesque, posuere massa nec")

	fs.StringVar(&str, "str", "yep", ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. "+
		"Integer tristique, elit ac efficitur sodales, nibh felis "+
		"pulvinar neque, eu pellentesque odio risus sed risus. Nulla "+
		"ac sem ex. Suspendisse in orci pellentesque, posuere massa nec")

	fs.IntVar(&iv, "iv", 2, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. "+
		"Integer tristique, elit ac efficitur sodales, nibh felis "+
		"pulvinar neque, eu pellentesque odio risus sed risus. Nulla "+
		"ac sem ex. Suspendisse in orci pellentesque, posuere massa nec")

	fs.DurationVar(&dv, "dv", 1*time.Second, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. "+
		"Integer tristique, elit ac efficitur sodales, nibh felis "+
		"pulvinar neque, eu pellentesque odio risus sed risus. Nulla "+
		"ac sem ex. Suspendisse in orci pellentesque, posuere massa nec")

	fs.BoolVar(&pants, "b", true, ""+
		"Morbi ac elementum massa. Sed bibendum vel magna eget sagittis. "+
		"Integer tristique, elit ac efficitur sodales, nibh felis "+
		"pulvinar neque, eu pellentesque odio risus sed risus. Nulla "+
		"ac sem ex. Suspendisse in orci pellentesque, posuere massa nec")

	fmt.Println(fs.Usage())
}
