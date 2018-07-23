package cmdy

import (
	"testing"

	"github.com/shabbyrobe/golib/assert"
)

func TestNonOptionalArgAfterOptionalArg(t *testing.T) {
	type vals struct {
		foo, bar, baz string
	}
	setup := func() (*ArgSet, *vals) {
		var v vals
		as := NewArgSet()
		as.String(&v.foo, "<foo>", "Usage...")
		as.StringOptional(&v.bar, "<bar>", "default", "Usage...")
		as.String(&v.baz, "<baz>", "Usage...")
		return as, &v
	}

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{"a", "b", "c"}))
		tt.MustEqual("a", v.foo)
		tt.MustEqual("b", v.bar)
		tt.MustEqual("c", v.baz)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual("a", v.foo)
		tt.MustEqual("default", v.bar)
		tt.MustEqual("", v.baz)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		err := as.Parse([]string{})
		tt.MustAssert(err != nil) // FIXME: check error
		tt.MustEqual("", v.foo)
		tt.MustEqual("default", v.bar)
		tt.MustEqual("", v.baz)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)
		as, v := setup()
		err := as.Parse([]string{"a", "b", "c", "d"})
		tt.MustAssert(err != nil) // FIXME: check error
		tt.MustEqual("a", v.foo)
		tt.MustEqual("b", v.bar)
		tt.MustEqual("c", v.baz)
	})
}

func TestRemainingOnly(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", AnyLen, "Usage...")

		tt.MustOK(as.Parse([]string{"a", "b", "c"}))
		tt.MustEqual([]string{"a", "b", "c"}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", AnyLen, "Usage...")

		tt.MustOK(as.Parse([]string{}))
		tt.MustEqual(0, len(foo))
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Min(1), "Usage...")

		err := as.Parse([]string{})
		tt.MustAssert(err != nil) // FIXME: check error
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Min(1), "Usage...")

		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual([]string{"a"}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Max(1), "Usage...")

		tt.MustOK(as.Parse([]string{"a"}))
		tt.MustEqual([]string{"a"}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []string
		as := NewArgSet()
		as.Remaining(&foo, "foo", Max(1), "Usage...")

		err := as.Parse([]string{"a", "b"})
		tt.MustAssert(err != nil) // FIXME: check error
	})
}

func TestRemainingInts(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []int
		as := NewArgSet()
		as.RemainingInts(&foo, "foo", AnyLen, "Usage...")

		tt.MustOK(as.Parse([]string{"1", "2", "3"}))
		tt.MustEqual([]int{1, 2, 3}, foo)
	})

	t.Run("", func(t *testing.T) {
		tt := assert.WrapTB(t)

		var foo []int
		as := NewArgSet()
		as.RemainingInts(&foo, "foo", AnyLen, "Usage...")

		err := as.Parse([]string{"1", "2", "quack"})
		tt.MustAssert(err != nil) // FIXME: check error
	})
}
