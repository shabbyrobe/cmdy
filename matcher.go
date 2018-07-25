package cmdy

import (
	"sort"
)

// PrefixMatcher returns a primitive Matcher for use with a command Group that
// will match a command if the input is an unambiguous prefix of one of the
// Group's Builders.
func PrefixMatcher(group *Group, minLen int) Matcher {
	if minLen <= 0 {
		panic("minLen must be > 0")
	}

	strs := make([]string, 0, len(group.Builders))
	for s := range group.Builders {
		strs = append(strs, s)
	}
	sort.Strings(strs)

	return func(bldrs Builders, in string) (bld Builder, name string, rerr error) {
		max := 0
		ilen := len(in)
		for _, str := range strs {
			var cur int
			var slen = len(str)
			if ilen > slen {
				continue
			} else if str == in {
				return group.Builders[str], str, nil
			}

			for i := 0; i < slen; i++ {
				if i >= ilen || str[i] != in[i] {
					break
				}
				cur++
			}

			if cur > 0 && cur >= minLen {
				if cur < max || cur < ilen {
					return bld, name, nil
				} else if cur == max {
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
