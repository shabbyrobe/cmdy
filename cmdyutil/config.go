package cmdyutil

import (
	"fmt"

	"github.com/shabbyrobe/cmdy/internal/cmdstr"
)

// Config is a very crude configuration file format based on keyed sets
// of flags/args.
//
// It's an attempt to avoid introducing heavy dependencies and an attempt
// to avoid using the stdlib's insuficcient options for structured config
// (JSON, XML).
//
// The API is very unstable.
type Config struct {
	index map[string]ConfigSet
	sets  []ConfigSet
}

func NewConfig() *Config {
	return &Config{
		index: map[string]ConfigSet{},
	}
}

func (c *Config) Sets() []ConfigSet {
	return c.sets
}

func (c *Config) Args(name string) (args []string, ok bool) {
	set, ok := c.index[name]
	return set.args, ok
}

func (c *Config) AddArgs(name string, args []string) (ok bool) {
	if _, ok := c.index[name]; ok {
		return false
	}
	c.SetArgs(name, args)
	return true
}

func (c *Config) SetArgs(name string, args []string) {
	setData := ConfigSet{name: name, args: args}
	c.sets = append(c.sets, setData)
	c.index[name] = setData
}

type ConfigSet struct {
	name string
	args []string
}

type configState int

const (
	configNone configState = iota
	configSetHeader
	configArg
	configComment
)

func calculatePos(data []byte, idx int) (line, col int) {
	line = 1
	start := 0
	for i := 0; i < idx; i++ {
		if data[i] == '\n' {
			start = i
			line++
		}
	}
	return line, idx - start
}

func formatPos(data []byte, idx int) string {
	line, col := calculatePos(data, idx)
	return fmt.Sprintf("line %d, col %d", line, col)
}

func ParseConfig(data []byte) (*Config, error) {
	var (
		wsp     = [256]bool{'\r': true, '\n': true, '\t': true, ' ': true}
		setChar = [256]bool{
			'-': true, '.': true, '_': true,
			'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true, '9': true,
			'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true, 'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true, 'q': true, 'r': true, 's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true, 'y': true, 'z': true,
			'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true, 'G': true, 'H': true, 'I': true, 'J': true, 'K': true, 'L': true, 'M': true, 'N': true, 'O': true, 'P': true, 'Q': true, 'R': true, 'S': true, 'T': true, 'U': true, 'V': true, 'W': true, 'X': true, 'Y': true, 'Z': true,
		}
		stack    = [4]configState{}
		stackPos = 0
		sz       = len(data)
		i        = 0
		start    = 0

		// accumulators:
		set    = ""
		args   []string
		config = NewConfig()
	)

	for i = 0; i < sz; i++ {
	retry:
		c := data[i]

		switch stack[stackPos] {
		case configNone:
			if wsp[c] {
				continue
			}
			switch c {
			case '#':
				stackPos++
				stack[stackPos] = configComment
			case '[':
				stackPos++
				stack[stackPos], start = configSetHeader, i+1
			default:
				return nil, fmt.Errorf("config parsing failed while looking for [set]; unexpected character %q at %s", string(c), formatPos(data, i))
			}

		case configSetHeader:
			if c == ']' {
				set = string(data[start:i])
				if len(set) == 0 {
					return nil, fmt.Errorf("config parsing failed: [set] was empty at %s", formatPos(data, i))
				}

				// pop/push at the same time
				stack[stackPos], start, args = configArg, i+1, nil

			} else if !setChar[c] {
				return nil, fmt.Errorf("config parsing failed while reading [set]; unexpected character %q at %s", string(c), formatPos(data, i))
			}

		case configArg:
			nextArgs, n, err := cmdstr.Parse(data[i:], "[]#")
			if err != nil {
				return nil, fmt.Errorf("config parsing failed while reading -flags; %w", err)
			}
			args = append(args, nextArgs...)
			i += n

			if i == sz || data[i] == '[' {
				if !config.AddArgs(set, args) {
					return nil, fmt.Errorf("config parsing failed: duplicate set [%s] detected", set)
				}

				stackPos--
				if i == sz {
					goto end
				} else {
					goto retry
				}

			} else if data[i] == '#' {
				stackPos++
				stack[stackPos] = configComment

			} else {
				return nil, fmt.Errorf("config parsing failed while reading -flags; unexpected character %q at %s", string(c), formatPos(data, i))
			}

		case configComment:
			if c == '\n' {
				stackPos--
			}

		default:
			panic("unknown state")
		}
	}

end:
	if stack[stackPos] == configComment {
		stackPos--
	}
	if stack[stackPos] == configArg {
		if !config.AddArgs(set, args) {
			return nil, fmt.Errorf("config parsing failed: duplicate set [%s] detected", set)
		}
		stackPos--
	}
	if stack[stackPos] == configSetHeader {
		return nil, fmt.Errorf("config parsing failed: unclosed set")
	}

	return config, nil
}
