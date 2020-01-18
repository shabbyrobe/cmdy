package wrap

import "strings"

const DefaultWrap = 80

type Wrapper struct {
	Prefix      string
	WrapWith    string // Appears at the end of every line that breaks. Defaults to "\n" if empty.
	Indent      string
	IndentFirst bool
	Width       int
}

func (w Wrapper) Wrap(str string) string {
	str = strings.TrimSpace(str)

	width := w.Width
	if width <= 0 {
		width = DefaultWrap
	}

	var out strings.Builder
	if w.IndentFirst {
		out.WriteString(w.Indent)
	}

	prefix := w.Indent
	if w.Prefix != "" {
		out.WriteString(w.Prefix)
		prefix += w.Indent
	}

	wrapWith := w.WrapWith
	if wrapWith == "" {
		wrapWith = "\n"
	}
	wrapWith += prefix

	var ln int
	for _, line := range strings.Split(str, "\n") {
		for {
			if ln > 0 {
				out.WriteString(wrapWith)
			}

			var (
				i, j     int
				c        rune
				breaking bool
				inEsc    bool
			)

			for j, c = range line {
				if i == width {
					breaking = true
					break
				}

				// FIXME: This tries not to count ASCII escape sequences that change the
				// colour towards line width, but it breaks badly if the text contains
				// other escapes or control sequences.
				if inEsc {
					if c == 'm' {
						inEsc = false
					}
				} else {
					if c == '\033' {
						inEsc = true
					} else {
						i++
					}
				}
			}

			cur := line[:j]
			idx := strings.LastIndexAny(cur, " -")
			if idx < 0 || !breaking {
				out.WriteString(line)
				break
			} else {
				out.WriteString(line[:idx])
				line = line[idx+1:]
			}
			ln++
		}
	}

	return out.String()
}
