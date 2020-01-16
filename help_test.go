package cmdy

import (
	"fmt"
	"strings"
	"testing"
)

func TestExampleRender(t *testing.T) {
	ex := Example{
		Desc:    "thingo",
		Command: "floobfleebflarbflem flubflogfleedfloobfleebflarbflemflubflogfleedfloobfleeb flarb flem flub flog fleed",
		Input:   "borg",
		Output:  "it work\nyep it work\nyeppo",
	}

	var o strings.Builder
	ex.render(&o, nil, "  ")
	fmt.Println(o.String())
}
