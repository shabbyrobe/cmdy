package arg

import (
	"fmt"
	"log"
	"testing"

	"github.com/shabbyrobe/cmdy/internal/assert"
)

func ExampleArgSet_RemainingInts() {
	var foos []int
	as := NewArgSet()
	as.RemainingInts(&foos, "foos", AnyLen, "Usage...")

	if err := as.Parse([]string{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println(foos)

	foos = nil
	if err := as.Parse([]string{"1"}); err != nil {
		log.Fatal(err)
	}
	fmt.Println(foos)

	foos = nil
	if err := as.Parse([]string{"1", "2"}); err != nil {
		log.Fatal(err)
	}
	fmt.Println(foos)

	// Output:
	// []
	// [1]
	// [1 2]
}

func TestRemainingInts(t *testing.T) {
	for idx, tc := range []struct {
		ok  bool
		in  []string
		out []int
	}{
		{true, []string{}, nil},
		{true, []string{"1"}, []int{1}},
		{true, []string{"1", "2"}, []int{1, 2}},
		{true, []string{"0xff", "0x10"}, []int{0xff, 0x10}},
		{false, []string{"quack"}, nil},
		{false, []string{"1", "quack"}, nil},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			var foo []int
			as := NewArgSet()
			as.RemainingInts(&foo, "foo", AnyLen, "Usage...")
			err := as.Parse(tc.in)
			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(tc.out, foo)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

func TestRemainingInt64s(t *testing.T) {
	for idx, tc := range []struct {
		ok  bool
		in  []string
		out []int64
	}{
		{true, []string{}, nil},
		{true, []string{"1"}, []int64{1}},
		{true, []string{"1", "2"}, []int64{1, 2}},
		{true, []string{"0xff", "0x10"}, []int64{0xff, 0x10}},
		{false, []string{"quack"}, nil},
		{false, []string{"1", "quack"}, nil},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			var foo []int64
			as := NewArgSet()
			as.RemainingInt64s(&foo, "foo", AnyLen, "Usage...")
			err := as.Parse(tc.in)
			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(tc.out, foo)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

func TestRemainingUints(t *testing.T) {
	for idx, tc := range []struct {
		ok  bool
		in  []string
		out []uint
	}{
		{true, []string{}, nil},
		{true, []string{"1"}, []uint{1}},
		{true, []string{"1", "2"}, []uint{1, 2}},
		{true, []string{"0xff", "0x10"}, []uint{0xff, 0x10}},
		{false, []string{"quack"}, nil},
		{false, []string{"1", "quack"}, nil},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			var foo []uint
			as := NewArgSet()
			as.RemainingUints(&foo, "foo", AnyLen, "Usage...")
			err := as.Parse(tc.in)
			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(tc.out, foo)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

func TestRemainingUint64s(t *testing.T) {
	for idx, tc := range []struct {
		ok  bool
		in  []string
		out []uint64
	}{
		{true, []string{}, nil},
		{true, []string{"1"}, []uint64{1}},
		{true, []string{"1", "2"}, []uint64{1, 2}},
		{true, []string{"0xff", "0x10"}, []uint64{0xff, 0x10}},
		{false, []string{"quack"}, nil},
		{false, []string{"1", "quack"}, nil},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			var foo []uint64
			as := NewArgSet()
			as.RemainingUint64s(&foo, "foo", AnyLen, "Usage...")
			err := as.Parse(tc.in)
			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(tc.out, foo)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

func TestRemainingVar(t *testing.T) {
	for idx, tc := range []struct {
		ok  bool
		in  []string
		out []int
	}{
		{true, []string{}, nil},
		{true, []string{"1"}, []int{1}},
		{true, []string{"1", "2"}, []int{1, 2}},
		{true, []string{"0xff", "0x10"}, []int{0xff, 0x10}},
		{false, []string{"quack"}, nil},
		{false, []string{"1", "quack"}, nil},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			var foo intSliceVar
			as := NewArgSet()
			as.RemainingVar(&foo, "foo", AnyLen, "Usage...")
			err := as.Parse(tc.in)
			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(tc.out, []int(foo))
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

func TestRemainingFloat64s(t *testing.T) {
	for idx, tc := range []struct {
		ok  bool
		in  []string
		out []float64
	}{
		{true, []string{}, nil},
		{true, []string{"1"}, []float64{1}},
		{true, []string{"1.0"}, []float64{1}},
		{true, []string{"-1.0"}, []float64{-1}},
		{true, []string{"1", "2"}, []float64{1, 2}},
		{false, []string{"quack"}, nil},
		{false, []string{"1", "quack"}, nil},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			var foo []float64
			as := NewArgSet()
			as.RemainingFloat64s(&foo, "foo", AnyLen, "Usage...")
			err := as.Parse(tc.in)
			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(tc.out, foo)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

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

func TestStringOptional(t *testing.T) {
	t.Run("value-no-default", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var s string
		as := NewArgSet()
		as.StringOptional(&s, "str", "", "yep")
		tt.MustOK(as.Parse([]string{"foo"}))
		tt.MustEqual("foo", s)
		tt.MustEqual("<str> (kindo) hinto", as.args[0].Describe("kindo", "hinto"))
	})

	t.Run("no-value-default", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var s string
		as := NewArgSet()
		as.StringOptional(&s, "str", "dflt", "yep")
		tt.MustOK(as.Parse([]string{}))
		tt.MustEqual("dflt", s)
		tt.MustEqual("<str> (kindo) hinto", as.args[0].Describe("kindo", "hinto"))
	})

	t.Run("value-default", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var s string
		as := NewArgSet()
		as.StringOptional(&s, "str", "yep", "yep")
		tt.MustOK(as.Parse([]string{"foo"}))
		tt.MustEqual("foo", s)
		tt.MustEqual("<str> (kindo) hinto", as.args[0].Describe("kindo", "hinto"))
	})
}

func TestVarOptional(t *testing.T) {
	t.Run("value-no-default", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var s stringVar
		as := NewArgSet()
		as.VarOptional(&s, "str", "yep")
		tt.MustOK(as.Parse([]string{"foo"}))
		tt.MustEqual("foo", string(s))
		tt.MustEqual("<str> (kindo) hinto", as.args[0].Describe("kindo", "hinto"))
	})

	t.Run("no-value-default", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var s stringVar = "dflt"
		as := NewArgSet()
		as.VarOptional(&s, "str", "yep")
		tt.MustOK(as.Parse([]string{}))
		tt.MustEqual("dflt", string(s))
		tt.MustEqual("<str> (kindo) hinto", as.args[0].Describe("kindo", "hinto"))
	})

	t.Run("value-default", func(t *testing.T) {
		tt := assert.WrapTB(t)
		var s stringVar = "dflt"
		as := NewArgSet()
		as.VarOptional(&s, "str", "yep")
		tt.MustOK(as.Parse([]string{"foo"}))
		tt.MustEqual("foo", string(s))
		tt.MustEqual("<str> (kindo) hinto", as.args[0].Describe("kindo", "hinto"))
	})
}
