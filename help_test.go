package cmdy

import (
	"strings"
	"testing"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

const exampleRenderResult = `
  # thingo
  $ echo "borg" | floobfleebflarbflem \
      flubflobfleedfloobfleebflarbflemflubflobfleedfloobfleeb flarb flem flub flob \
      fleed
  it works
  yep it works
  ...
`

func TestExampleRender(t *testing.T) {
	tt := assert.WrapTB(t)
	ex := Example{
		Desc:    "thingo",
		Command: "floobfleebflarbflem flubflobfleedfloobfleebflarbflemflubflobfleedfloobfleeb flarb flem flub flob fleed",
		Input:   "borg",
		Output:  "it works\nyep it works\nyeppo",
	}

	var o strings.Builder
	es := exampleSection{}
	es.renderExample(&o, &ex, "")
	tt.MustEqual(strings.TrimRight(exampleRenderResult[1:], "\n"), strings.TrimRight(o.String(), "\n"))
}
