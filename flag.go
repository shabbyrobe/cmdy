package cmdy

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/shabbyrobe/cmdy/usage"
)

type FlagSet struct {
	*flag.FlagSet
	WrapWidth int
	buf       bytes.Buffer
	hideUsage bool
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{
		FlagSet: flag.NewFlagSet("", flag.ContinueOnError),
	}
	fs.FlagSet.SetOutput(&fs.buf)
	return fs
}

func (fs *FlagSet) HideUsage() { fs.hideUsage = true }

func (fs *FlagSet) Usage() string {
	if fs.hideUsage {
		return ""
	}

	var usables []usage.Usable
	fs.VisitAll(func(f *flag.Flag) {
		usables = append(usables, usableFlag{f})
	})
	return usage.Usage(fs.WrapWidth, usables...)
}

type usableFlag struct {
	flag *flag.Flag
}

func (u usableFlag) Name() string       { return u.flag.Name }
func (u usableFlag) Usage() string      { return u.flag.Usage }
func (u usableFlag) DefValue() string   { return u.flag.DefValue }
func (u usableFlag) Value() interface{} { return u.flag.Value }

func (u usableFlag) Describe(kind string, hint string) string {
	name := u.Name()
	if kind != "" && hint != "" {
		return fmt.Sprintf("-%s=<%s> (%s)", name, kind, hint)
	} else if kind != "" {
		return fmt.Sprintf("-%s=<%s>", name, kind)
	} else {
		return fmt.Sprintf("-%s", name)
	}
}
