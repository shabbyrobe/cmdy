package cmdy

import (
	"sort"
)

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
			}

			for i := 0; i < slen; i++ {
				if i >= ilen || str[i] != in[i] {
					break
				}
				cur++
			}

			if cur+1 == slen { // exact match
				return group.Builders[str], str, nil
			}

			if cur > 0 && cur > minLen {
				if cur < max {
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
