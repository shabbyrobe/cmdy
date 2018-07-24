package cmdy

import (
	"bytes"
	"flag"

	"github.com/shabbyrobe/cmdy/usage"
)

type FlagSet struct {
	*flag.FlagSet
	WrapWidth int
	buf       bytes.Buffer
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{
		FlagSet: flag.NewFlagSet("", flag.ContinueOnError),
	}
	fs.FlagSet.SetOutput(&fs.buf)
	return fs
}

func (fs *FlagSet) Usage() string {
	var usables []usage.Usable
	fs.VisitAll(func(f *flag.Flag) {
		usables = append(usables, usableFlag{f})
	})
	return usage.Usage(0, usables...)
}

type usableFlag struct {
	flag *flag.Flag
}

func (u usableFlag) Name() string       { return u.flag.Name }
func (u usableFlag) Usage() string      { return u.flag.Usage }
func (u usableFlag) DefValue() string   { return u.flag.DefValue }
func (u usableFlag) Value() interface{} { return u.flag.Value }
