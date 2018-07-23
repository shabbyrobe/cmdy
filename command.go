package cmdy

import "flag"

type FlagSet = flag.FlagSet

type Usage interface {
	Synopsis() string
	Usage() string
}

type Command interface {
	Usage

	Flags() *FlagSet
	Args() *ArgSet
}
