package args

// All provides an ArgSet that collects all of the arguments into
// a slice of strings.
func All(into *[]string, name string, usage string) *ArgSet {
	as := NewArgSet()
	as.Remaining(into, name, AnyLen, usage)
	return as
}
