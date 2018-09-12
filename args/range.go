package args

// Unlimited is a sentinel value used for the Max of a Range to indicate that
// it is unbounded.
//
// FIXME: This could probably be set to 0; why would you specify a Remaining
// Arg that had a Max of 0 when you could just not specify the Remaining arg
// at all?
const Unlimited = -1

// AnyLen is used for a Remaining arg that can take 0 or more arguments.
//
//	var strs []string
//	argSet.Remaining(&strs, "strs", args.AnyLen, "All remaining string arguments")
//
var AnyLen = Range{0, Unlimited}

type Range struct {
	Min int

	// Inclusive upper bound for this range; set to Unlimited for no limit.
	Max int
}

// Min is used to place a lower bound on the number of values required for
// a Remaining arg.
//
// In the following example, the args will pass validation if there are 2
// or more:
//
//	argSet.Remaining(&strs, "strs", args.Min(2), "Remaining strings")
//
//	$ myprog            // invalid
//	$ myprog a          // invalid
//	$ myprog a b        // valid
//	$ myprog a b c      // valid
//
func Min(min int) Range { return Range{min, Unlimited} }

// Max is used to place an upper bound on the number of values allowed for
// a Remaining arg.
//
// In the following example, the args will pass validation if there are 2
// or fewer:
//
//	argSet.Remaining(&strs, "strs", args.Max(2), "Remaining strings")
//
//	$ myprog            // valid
//	$ myprog a          // valid
//	$ myprog a b        // valid
//	$ myprog a b c      // invalid
//
func Max(max int) Range { return Range{0, max} }

// MinMax is used to place an upper and lower bound on the number of values
// allowed for a Remaining arg.
//
// In the following example, the args will pass validation if there are
// between 2 and 4 args inclusive:
//
//	argSet.Remaining(&strs, "strs", args.MinMax(2, 4), "Remaining strings")
//
//	$ myprog a          // invalid
//	$ myprog a b        // valid
//	$ myprog a b c      // valid
//	$ myprog a b c d    // valid
//	$ myprog a b c d e  // invalid
//
func MinMax(min, max int) Range { return Range{min, max} }
