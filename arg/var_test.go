package arg

import (
	"testing"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

func TestIntVar(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var i int
		as := NewArgSet()
		as.Int(&i, "inty", "Usage!")
		tt.MustOK(as.Parse([]string{"1"}))
		tt.MustEqual(1, i)
		tt.MustEqual("<inty> (int) hint", as.args[0].Describe("int", "hint"))
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var i int
		as := NewArgSet()
		as.Int(&i, "inty", "Usage!")
		err := as.Parse([]string{"nope"})
		tt.MustAssert(err != nil)
	})
}

func TestInt64Var(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var i int64
		as := NewArgSet()
		as.Int64(&i, "inty", "Usage!")
		tt.MustOK(as.Parse([]string{"1"}))
		tt.MustEqual(int64(1), i)
		tt.MustEqual("<inty> (int) hint", as.args[0].Describe("int", "hint"))
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var i int64
		as := NewArgSet()
		as.Int64(&i, "inty", "Usage!")
		err := as.Parse([]string{"nope"})
		tt.MustAssert(err != nil)
	})
}

func TestFloat64Var(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var i float64
		as := NewArgSet()
		as.Float64(&i, "floaty", "Usage!")
		tt.MustOK(as.Parse([]string{"1.1"}))
		tt.MustEqual(1.1, i)
		tt.MustEqual("<floaty> (mc) floatface", as.args[0].Describe("mc", "floatface"))
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var i float64
		as := NewArgSet()
		as.Float64(&i, "floaty", "Usage!")
		err := as.Parse([]string{"nope"})
		tt.MustAssert(err != nil)
	})
}
