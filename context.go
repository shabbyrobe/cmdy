package cmdy

import (
	"context"
	"io"
)

// Context implements context.Context; Context is passed into all commands
// when they are Run().
//
// If you want your Command to be testable, you should access Stdin, Stdout
// and Stderr via Context.
//
type Context interface {
	context.Context

	RawArgs() []string
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer

	Runner() *Runner
	Stack() []CommandRef
	Current() CommandRef
	Push(name string, cmd Command)
	Pop() (name string, cmd Command)
}

type commandContext struct {
	context.Context
	cmd     Command
	rawArgs []string
	runner  *Runner
	parents []CommandRef
}

func (c *commandContext) RawArgs() []string { return c.rawArgs }
func (c *commandContext) Stdin() io.Reader  { return c.runner.Stdin }
func (c *commandContext) Stdout() io.Writer { return c.runner.Stdout }
func (c *commandContext) Stderr() io.Writer { return c.runner.Stderr }

func (c *commandContext) Runner() *Runner     { return c.runner }
func (c *commandContext) Stack() []CommandRef { return c.parents }

func (c *commandContext) Current() (ref CommandRef) {
	if len(c.parents) > 0 {
		ref = c.parents[len(c.parents)-1]
	}
	return ref
}

func (c *commandContext) Push(name string, cmd Command) {
	c.parents = append(c.parents, CommandRef{Name: name, Command: cmd})
}

func (c *commandContext) Pop() (name string, cmd Command) {
	var item CommandRef
	pl := len(c.parents)
	c.parents, item = c.parents[:pl-1], c.parents[pl-1]
	return item.Name, item.Command
}

// FIXME: this name is not good

type CommandRef struct {
	Name    string
	Command Command
}

func CommandPath(ctx Context) (out []string) {
	parents := ctx.Stack()
	out = make([]string, 0, len(parents))
	for i := 0; i < len(parents); i++ {
		out = append(out, parents[i].Name)
	}
	return out
}
