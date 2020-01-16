// Package arg implements command-line argument parsing.
//
// It is very similar to the flag handling provided by the stdlib 'flag' package (without
// the global handling), and intended to be used as a complement to it.
//
// Usage
//
// Define an ArgSet using arg.NewArgSet(), then add arguments to it:
//	var s string
//	var i int
//	var set = arg.NewArgSet()
//	set.String(&s, "yep", "Usage for the 'yep' arg")
//	set.Int(&i, "inty", "Usage for the 'inty' arg")
//
// Or you can create custom args that satisfy the flag.Value (yes, 'flag'
// is intentional here) interface (with pointer receivers) and couple them to
// arg parsing:
//	var thingo myArgType
//	arg.Var(&thingo, "thingo", "Usage for 'thingo'")
//
// For such args, the default value is just the initial value of the variable.
//
// After all flags are defined, call ArgSet.Parse() with the list of arguments,
// which will usually be the output of FlagSet.Args():
//	err = arg.Parse(myFlagSet.Parse())
//	err = arg.Parse(os.Args[1:]) // if there are no flags
//
package arg
