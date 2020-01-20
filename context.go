package cmdy

import (
	"context"
	"io"
	"strings"
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

	// Stack contains references to all the commands and subcommands that have been Run,
	// but have not yet completed:
	Stack() CommandPath

	// Current returns a reference to the topmost Command on the stack:
	Current() CommandRef

	// Push and Pop are used by Run to ensure that the Command being run is
	// available in the Context's stack:
	Push(name string, cmd Command)
	Pop() (name string, cmd Command)
}

// IsDone returns true if the context's Done() channel would not block.
// This provides a simpler, more literate mechanism for checking context
// completion in commands that do not block in a select{}.
func IsDone(ctx Context) bool {
	return ctx.Err() != nil
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

func (c *commandContext) Runner() *Runner    { return c.runner }
func (c *commandContext) Stack() CommandPath { return c.parents }

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

type CommandPath []CommandRef

func (cp CommandPath) Invocation() string {
	return strings.Join(cp.Names(), " ")
}

func (cp CommandPath) Names() []string {
	out := make([]string, len(cp))
	for i, r := range cp {
		out[i] = r.Name
	}
	return out
}
