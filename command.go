package cmdy

import (
	"github.com/shabbyrobe/cmdy/args"
)

type Usage interface {
	Synopsis() string
	Usage() string
}

type Command interface {
	Usage

	Flags() *FlagSet
	Args() *args.ArgSet
	Run(Context, Input) error
}

// Builder creates an instance of your Command. The instance should be a new
// instance, not a recycled instance.
type Builder func() (Command, error)

type Input interface {
	Command() Command
	RawArgs() []string
}

type input struct {
	cmd     Command
	rawArgs []string
}

func (i *input) Command() Command  { return i.cmd }
func (i *input) RawArgs() []string { return i.rawArgs }
