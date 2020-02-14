package flags

import (
	"bytes"
	"flag"
	"testing"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

func TestOptionalString(t *testing.T) {
	tt := assert.WrapTB(t)

	var v OptionalString
	var fs flag.FlagSet
	fs.Var(&v, "s", "test")

	tt.MustOK(fs.Parse([]string{}))
	tt.MustAssert(!v.IsSet)

	tt.MustOK(fs.Parse([]string{"-s", ""}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual("", v.Value)
	tt.MustEqual("", v.String())

	tt.MustOK(fs.Parse([]string{"-s", "test"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual("test", v.Value)
	tt.MustEqual("test", v.String())
}

func TestOptionalInt(t *testing.T) {
	tt := assert.WrapTB(t)

	var v OptionalInt
	var fs flag.FlagSet
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Var(&v, "i", "test")

	tt.MustOK(fs.Parse([]string{}))
	tt.MustAssert(!v.IsSet)

	tt.MustAssert(fs.Parse([]string{"-i", ""}) != nil)

	tt.MustOK(fs.Parse([]string{"-i", "0"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(0, v.Value)
	tt.MustEqual("0", v.String())

	tt.MustOK(fs.Parse([]string{"-i", "1"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(1, v.Value)
	tt.MustEqual("1", v.String())
}

func TestOptionalInt64(t *testing.T) {
	tt := assert.WrapTB(t)

	var v OptionalInt64
	var fs flag.FlagSet
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Var(&v, "i", "test")

	tt.MustOK(fs.Parse([]string{}))
	tt.MustAssert(!v.IsSet)

	tt.MustAssert(fs.Parse([]string{"-i", ""}) != nil)

	tt.MustOK(fs.Parse([]string{"-i", "0"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(int64(0), v.Value)
	tt.MustEqual("0", v.String())

	tt.MustOK(fs.Parse([]string{"-i", "1"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(int64(1), v.Value)
	tt.MustEqual("1", v.String())

	tt.MustOK(fs.Parse([]string{"-i", "-1"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(int64(-1), v.Value)
	tt.MustEqual("-1", v.String())

	tt.MustOK(fs.Parse([]string{"-i", "4611686018427387904"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(int64(4611686018427387904), v.Value)
	tt.MustEqual("4611686018427387904", v.String())
}

func TestOptionalUint(t *testing.T) {
	tt := assert.WrapTB(t)

	var v OptionalUint
	var fs flag.FlagSet
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Var(&v, "i", "test")

	tt.MustOK(fs.Parse([]string{}))
	tt.MustAssert(!v.IsSet)

	tt.MustAssert(fs.Parse([]string{"-i", ""}) != nil)

	tt.MustOK(fs.Parse([]string{"-i", "0"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(uint(0), v.Value)
	tt.MustEqual("0", v.String())

	tt.MustOK(fs.Parse([]string{"-i", "1"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(uint(1), v.Value)
	tt.MustEqual("1", v.String())
}

func TestOptionalUint64(t *testing.T) {
	tt := assert.WrapTB(t)

	var v OptionalUint64
	var fs flag.FlagSet
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Var(&v, "i", "test")

	tt.MustOK(fs.Parse([]string{}))
	tt.MustAssert(!v.IsSet)

	tt.MustAssert(fs.Parse([]string{"-i", ""}) != nil)

	tt.MustOK(fs.Parse([]string{"-i", "0"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(uint64(0), v.Value)
	tt.MustEqual("0", v.String())

	tt.MustOK(fs.Parse([]string{"-i", "1"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(uint64(1), v.Value)
	tt.MustEqual("1", v.String())

	tt.MustOK(fs.Parse([]string{"-i", "4611686018427387904"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(uint64(4611686018427387904), v.Value)
	tt.MustEqual("4611686018427387904", v.String())
}

func TestOptionalFloat64(t *testing.T) {
	tt := assert.WrapTB(t)

	var v OptionalFloat64
	var fs flag.FlagSet
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Var(&v, "f", "test")

	tt.MustOK(fs.Parse([]string{}))
	tt.MustAssert(!v.IsSet)

	tt.MustAssert(fs.Parse([]string{"-f", ""}) != nil)

	tt.MustOK(fs.Parse([]string{"-f", "0"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(float64(0), v.Value)
	tt.MustEqual("0", v.String())

	tt.MustOK(fs.Parse([]string{"-f", "1"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(float64(1), v.Value)
	tt.MustEqual("1", v.String())
}

func TestOptionalBool(t *testing.T) {
	tt := assert.WrapTB(t)

	var v OptionalBool
	var fs flag.FlagSet
	var buf bytes.Buffer
	fs.SetOutput(&buf)
	fs.Var(&v, "b", "test")

	tt.MustOK(fs.Parse([]string{}))
	tt.MustAssert(!v.IsSet)

	tt.MustOK(fs.Parse([]string{"-b"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(true, v.Value)
	tt.MustEqual("true", v.String())

	tt.MustOK(fs.Parse([]string{"-b=true"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(true, v.Value)
	tt.MustEqual("true", v.String())

	tt.MustOK(fs.Parse([]string{"-b=false"}))
	tt.MustAssert(v.IsSet)
	tt.MustEqual(false, v.Value)
	tt.MustEqual("false", v.String())
}
