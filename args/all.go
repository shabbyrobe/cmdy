package args

func All(into *[]string, name string, usage string) *ArgSet {
	as := NewArgSet()
	as.Remaining(into, name, AnyLen, usage)
	return as
}
