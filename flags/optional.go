package flags

import "strconv"

// Gah! Boilerplate everywhere!

// OptionalString is a string flag that lets you to tell if the flag has been
// set or not, unlike flag.StringVar(), which does not.
//
// OptionalString can distinguish between all three of these conditions:
//
//	$ prog -flag foo
//	$ prog -flag ""j
//	$ prog
//
type OptionalString struct {
	IsSet bool
	Value string
}

func (s *OptionalString) Set(x string) error {
	s.Value = x
	s.IsSet = true
	return nil
}

func (s *OptionalString) String() string {
	if s == nil {
		return ""
	}
	return s.Value
}

// OptionalInt is an int flag that lets you to tell if the flag has been
// set or not, unlike flag.IntVar(), which does not.
//
// OptionalInt can distinguish between all three of these conditions:
//
//	$ prog -int 1
//	$ prog -int 0
//	$ prog -int ""
//
type OptionalInt struct {
	IsSet bool
	Value int
	str   string
}

func (s *OptionalInt) Set(x string) error {
	i, err := strconv.ParseInt(x, 10, 0)
	if err != nil {
		return err
	}
	s.Value = int(i)
	s.IsSet = true
	s.str = x
	return nil
}

func (s *OptionalInt) String() string {
	if s == nil {
		return ""
	}
	return s.str
}

// OptionalInt64 is an int flag that lets you to tell if the flag has been
// set or not, unlike flag.Int64Var(), which does not.
//
// OptionalInt64 can distinguish between all three of these conditions:
//
//	$ prog -int64 1
//	$ prog -int64 0
//	$ prog -int64 ""
//
type OptionalInt64 struct {
	IsSet bool
	Value int64
	str   string
}

func (s *OptionalInt64) Set(x string) error {
	i, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return err
	}
	s.Value = i
	s.IsSet = true
	s.str = x
	return nil
}

func (s *OptionalInt64) String() string {
	if s == nil {
		return ""
	}
	return s.str
}

// OptionalUint is an int flag that lets you to tell if the flag has been
// set or not, unlike flag.UintVar(), which does not.
//
// OptionalUint can distinguish between all three of these conditions:
//
//	$ prog -uint 1
//	$ prog -uint 0
//	$ prog -uint ""
//
type OptionalUint struct {
	IsSet bool
	Value uint
	str   string
}

func (s *OptionalUint) Set(x string) error {
	i, err := strconv.ParseUint(x, 10, 0)
	if err != nil {
		return err
	}
	s.Value = uint(i)
	s.IsSet = true
	s.str = x
	return nil
}

func (s *OptionalUint) String() string {
	if s == nil {
		return ""
	}
	return s.str
}

// OptionalUint64 is an int flag that lets you to tell if the flag has been
// set or not, unlike flag.Uint64Var(), which does not.
//
// OptionalUint64 can distinguish between all three of these conditions:
//
//	$ prog -uint64 1
//	$ prog -uint64 0
//	$ prog
//
type OptionalUint64 struct {
	IsSet bool
	Value uint64
	str   string
}

func (s *OptionalUint64) Set(x string) error {
	i, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return err
	}
	s.Value = i
	s.IsSet = true
	s.str = x
	return nil
}

func (s *OptionalUint64) String() string {
	if s == nil {
		return ""
	}
	return s.str
}

// OptionalBool is an bool flag that lets you to tell if the flag has been
// set or not, unlike flag.BoolVar(), which does not.
//
// OptionalBool can distinguish between all three of these conditions:
//
//	$ prog -bool=true
//	$ prog -bool=false
//	$ prog
//
// OptionalBool supports, but can not distinguish, this additional condition:
//
//	$ prog -bool
//
// Like flag.BoolVar(), OptionalBool does not support these conditions:
//
//	$ prog -bool true
//	$ prog -bool false
//
type OptionalBool struct {
	IsSet bool
	Value bool
}

func (s *OptionalBool) IsBoolFlag() bool { return true }

func (s *OptionalBool) Set(x string) error {
	b, err := strconv.ParseBool(x)
	if err != nil {
		return err
	}
	s.Value = b
	s.IsSet = true
	return nil
}

func (s *OptionalBool) String() string {
	if s == nil {
		return ""
	}
	if !s.IsSet {
		return ""
	} else if s.Value {
		return "true"
	} else {
		return "false"
	}
}

// OptionalFloat64 is a float64 flag that lets you to tell if the flag has been
// set or not, unlike flag.Float64Var(), which does not.
//
// OptionalFloat64 can distinguish between all three of these conditions:
//
//	$ prog -float64 1.1
//	$ prog -float64 0
//	$ prog
//
type OptionalFloat64 struct {
	IsSet bool
	Value float64
	str   string
}

func (s *OptionalFloat64) Set(x string) error {
	f, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return err
	}
	s.Value = f
	s.IsSet = true
	s.str = x
	return nil
}

func (s *OptionalFloat64) String() string {
	if s == nil {
		return ""
	}
	return s.str
}
