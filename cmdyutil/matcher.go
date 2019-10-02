package cmdyutil

import (
	"sort"

	"github.com/shabbyrobe/cmdy"
)

func GroupPrefixMatcher(minLen int) cmdy.GroupOption {
	return func(cs *cmdy.Group) { cs.Matcher = PrefixMatcher(cs, minLen) }
}

// PrefixMatcher returns a simple Matcher for use with a command Group that will match a
// command if the input is an unambiguous prefix of one of the Group's Builders.
//
//	grp := NewGroup("grp", Builders{
//		"foo":  fooBuilder,
//		"bar":  barBuilder,
//		"bark": barkBuilder,
//		"bork": borkBuilder,
//	})
//
//	// Matches must be 2 or more characters to be considered:
//	m := PrefixMatcher(grp, 2)
//
//	$ myprog grp fo   // fooBuilder
//	$ myprog grp ba   // NOPE; bar or bark
//	$ myprog grp bar  // barBuilder
//	$ myprog grp bark // barkBuilder
//	$ myprog grp b    // NOPE; too short
//
func PrefixMatcher(group *cmdy.Group, minLen int) cmdy.Matcher {
	if minLen <= 0 {
		panic("minLen must be > 0")
	}

	strs := make([]string, 0, len(group.Builders))
	for s := range group.Builders {
		strs = append(strs, s)
	}
	sort.Strings(strs)

	return func(bldrs cmdy.Builders, in string) (bld cmdy.Builder, name string, rerr error) {
		max := 0
		inlen := len(in)
		for _, str := range strs {
			var cur int
			var curlen = len(str)
			if inlen > curlen {
				continue
			} else if str == in {
				return group.Builders[str], str, nil
			}

			for i := 0; i < curlen; i++ {
				if i >= inlen || str[i] != in[i] {
					break
				}
				cur++
			}

			if cur > 0 && cur >= minLen {
				if cur == max {
					return nil, "", nil
				} else if cur > max {
					max = cur
					bld, name = group.Builders[str], str
				}
			}
		}
		return bld, name, nil
	}
}
