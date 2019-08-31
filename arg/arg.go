package arg

import (
	"fmt"
)

type Arg struct {
	name     string
	usage    string
	value    ArgVal
	defValue string
	optional bool
}

func (a *Arg) Name() string     { return a.name }
func (a *Arg) Usage() string    { return a.usage }
func (a *Arg) DefValue() string { return a.defValue }

func (a *Arg) Value() interface{} {
	if rem, ok := a.value.(*remaining); ok {
		return rem.arg
	}
	return a.value
}

func (a *Arg) Describe(kind string, hint string) string {
	name := a.Name()
	if _, ok := a.value.(*remaining); ok {
		name += "..."
	}

	if kind != "" && hint != "" {
		return fmt.Sprintf("<%s> (%s) %s", name, kind, hint)
	} else if kind != "" {
		return fmt.Sprintf("<%s> (%s)", name, kind)
	} else if hint != "" {
		return fmt.Sprintf("<%s> %s", name, hint)
	} else {
		return fmt.Sprintf("<%s>", name)
	}
}

// ArgVal is exactly the same as flag.Value and should be 100% compatible.
type ArgVal interface {
	Set(string) error
	String() string
}

// ArgGetter is exactly the same as flag.Getter and should be 100% compatible.
type ArgGetter interface {
	ArgVal
	Get() interface{}
}
