package flags

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"testing"

	"github.com/ArtProcessors/cmdy"
	"github.com/ArtProcessors/cmdy/internal/assert"
)

func TestStringList(t *testing.T) {
	for idx, tc := range []struct {
		in  []string
		out []string
	}{
		{nil, nil},
		{[]string{"-s", ""}, []string{""}},
		{[]string{"-s", "yep"}, []string{"yep"}},
		{[]string{"-s", "yep", "-s", "yep"}, []string{"yep", "yep"}},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			var v StringList
			var fs flag.FlagSet
			fs.Var(&v, "s", "test")
			tt.MustOK(fs.Parse(tc.in))
			tt.MustEqual(StringList(tc.out), v)
		})
	}
}

type numCase struct {
	ok  bool
	in  []string
	out []string
}

func (sic numCase) AsInts() (vs []int, skip bool) {
	for _, v := range sic.out {
		v, err := strconv.ParseInt(v, 0, 0)
		if err != nil {
			return nil, true
		}
		vs = append(vs, int(v))
	}
	return vs, false
}

func (sic numCase) AsInt64s() (vs []int64, skip bool) {
	for _, v := range sic.out {
		v, err := strconv.ParseInt(v, 0, 64)
		if errors.Is(err, strconv.ErrRange) {
			return nil, true
		} else if err != nil {
			panic(v)
		}
		vs = append(vs, v)
	}
	return vs, false
}

func (sic numCase) AsFloat64s() (vs []float64, skip bool) {
	for _, fv := range sic.out {
		v, err := strconv.ParseFloat(fv, 64)
		if errors.Is(err, strconv.ErrRange) {
			return nil, true
		} else if err != nil {
			panic(v)
		}
		vs = append(vs, float64(v))
	}
	return vs, false
}

const yep, nup = true, false

// sharedNumCases contains test cases that are usable by all the numeric list
// flag var tests (floats and ints combined). Avoid underscores as well as hex, binary or
// octal literals for ints, or decimal portions, exponents, +/-inf, etc for floats.
//
// Unary negation is allowed, but will cause a skipped test for unsigned numbers.
var sharedNumCases = []numCase{
	{yep, nil, nil},
	{yep, []string{"-s", ""}, nil},
	{nup, []string{"-s", "nup"}, nil},
	{yep, []string{"-s", "1"}, []string{"1"}},
	{yep, []string{"-s", "1,2,3"}, []string{"1", "2", "3"}},
	{yep, []string{"-s", "1", "-s", "2", "-s", "3"}, []string{"1", "2", "3"}},
	{yep, []string{"-s", "1,2", "-s", "3"}, []string{"1", "2", "3"}},
	{yep, []string{"-s", " 1 , 2 , 3 "}, []string{"1", "2", "3"}},
}

var sharedIntCases = append(sharedNumCases, []numCase{
	{yep, []string{"-s", "0xf,0xf,0xf"}, []string{"0xf", "0xf", "0xf"}},
	{yep, []string{"-s", "0xf_f"}, []string{"0xff"}},
	{yep, []string{"-s", "0o777"}, []string{"0o777"}},
}...)

func TestInt64List(t *testing.T) {
	for idx, tc := range sharedIntCases {
		t.Run(fmt.Sprintf("int64/%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			out, skip := tc.AsInt64s()
			if skip {
				t.Skip()
			}

			var v Int64List
			fs := cmdy.NewFlagSet()
			fs.Var(&v, "s", "test")
			err := fs.Parse(tc.in)

			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(Int64List(out), v)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

func TestIntList(t *testing.T) {
	for idx, tc := range sharedIntCases {
		t.Run(fmt.Sprintf("int/%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			out, skip := tc.AsInts()
			if skip {
				t.Skip()
			}

			var v IntList
			fs := cmdy.NewFlagSet()
			fs.Var(&v, "s", "test")
			err := fs.Parse(tc.in)

			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(IntList(out), v)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}

func TestFloat64List(t *testing.T) {
	var cases = append(sharedNumCases, []numCase{
		{yep, []string{"-s", "0.123"}, []string{"0.123"}},
		{yep, []string{"-s", "+inf"}, []string{"+inf"}},
		{yep, []string{"-s", "0x1.b7p-1"}, []string{"0x1.b7p-1"}},
		{yep, []string{"-s", "0x1.b7p-1, 0x1.b7p-2"}, []string{"0x1.b7p-1", "0x1.b7p-2"}},
	}...)

	for idx, tc := range cases {
		t.Run(fmt.Sprintf("int/%d", idx), func(t *testing.T) {
			tt := assert.WrapTB(t)
			out, skip := tc.AsFloat64s()
			if skip {
				t.Skip()
			}

			var v Float64List
			fs := cmdy.NewFlagSet()
			fs.Var(&v, "s", "test")
			err := fs.Parse(tc.in)

			if tc.ok {
				tt.MustOK(err)
				tt.MustEqual(Float64List(out), v)
			} else {
				tt.MustAssert(err != nil)
			}
		})
	}
}
