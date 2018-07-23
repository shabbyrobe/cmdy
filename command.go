package cmdy

import (
	"context"
	"flag"

	"github.com/shabbyrobe/cmdy/args"
)

type FlagSet = flag.FlagSet

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

type Builder func() (Command, error)

func Run(ctx context.Context, b Builder, args []string) error {
	cmd, err := b()
	if err != nil {
		return err
	}

	var (
		flagSet = cmd.Flags()
		argSet  = cmd.Args()
		remArgs = args
	)

	if flagSet != nil {
		if err := flagSet.Parse(args); err != nil {
			return err // FIXME: repair help
		}
		remArgs = flagSet.Args()
	}

	if argSet != nil {
		if err := argSet.Parse(remArgs); err != nil {
			return err // FIXME: repair help
		}
	}

	cctx := &commandContext{ctx}
	input := &input{cmd: cmd, rawArgs: args}
	return cmd.Run(cctx, input)
}

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
