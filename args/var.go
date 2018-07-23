package cmdy

import (
	"strconv"
	"strings"
	"time"
)

type stringArg string

func (s *stringArg) Get() interface{} { return string(*s) }
func (s *stringArg) String() string   { return string(*s) }

func (s *stringArg) Set(val string) error {
	*s = stringArg(val)
	return nil
}

type intArg int

func (i *intArg) Get() interface{} { return int(*i) }
func (i *intArg) String() string   { return strconv.Itoa(int(*i)) }

func (i *intArg) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	*i = intArg(v)
	return err
}

type int64Arg int64

func (i *int64Arg) Get() interface{} { return int(*i) }
func (i *int64Arg) String() string   { return strconv.FormatInt(int64(*i), 10) }

func (i *int64Arg) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Arg(v)
	return err
}

type uintArg uint

func (i *uintArg) Get() interface{} { return uint(*i) }
func (i *uintArg) String() string   { return strconv.FormatUint(uint64(*i), 10) }

func (i *uintArg) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	*i = uintArg(v)
	return err
}

type uint64Arg uint64

func (i *uint64Arg) Get() interface{} { return int(*i) }
func (i *uint64Arg) String() string   { return strconv.FormatUint(uint64(*i), 10) }

func (i *uint64Arg) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Arg(v)
	return err
}

type float64Arg float64

func (f *float64Arg) Get() interface{} { return float64(*f) }
func (f *float64Arg) String() string   { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

func (f *float64Arg) Set(val string) error {
	v, err := strconv.ParseFloat(val, 64)
	*f = float64Arg(v)
	return err
}

type boolArg bool

func (b *boolArg) Get() interface{} { return bool(*b) }
func (b *boolArg) String() string   { return strconv.FormatBool(bool(*b)) }

func (b *boolArg) Set(val string) error {
	v, err := strconv.ParseBool(val)
	*b = boolArg(v)
	return err
}

type durationArg time.Duration

func (d *durationArg) Get() interface{} { return time.Duration(*d) }
func (d *durationArg) String() string   { return time.Duration(*d).String() }

func (d *durationArg) Set(val string) error {
	v, err := time.ParseDuration(val)
	*d = durationArg(v)
	return err
}

type remaining struct {
	arg ArgVal
	Range
}

func (r *remaining) Get() interface{}   { return r.arg.Get() }
func (r *remaining) Set(s string) error { return r.arg.Set(s) }
func (r *remaining) String() string     { return r.arg.String() }

type StringList []string

func (s *StringList) String() string {
	if s == nil {
		return ""
	}
	return strings.Join(*s, ",")
}

func (s *StringList) Get() interface{} { return []string(*s) }

func (s StringList) Strings() []string {
	out := make([]string, len(s))
	copy(out, s)
	return out
}

func (s *StringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type IntList []int

func (i *IntList) String() string {
	if i == nil {
		return ""
	}
	var out strings.Builder
	for idx, v := range *i {
		if idx != 0 {
			out.WriteByte(',')
		}
		out.WriteString(strconv.FormatInt(int64(v), 10))
	}
	return out.String()
}

func (i *IntList) Get() interface{} { return []int(*i) }

func (i IntList) Ints() []int { return []int(i) }

func (i *IntList) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	*i = append(*i, int(v))
	return err
}

type Int64List []int64

func (i *Int64List) String() string {
	if i == nil {
		return ""
	}
	var out strings.Builder
	for idx, v := range *i {
		if idx != 0 {
			out.WriteByte(',')
		}
		out.WriteString(strconv.FormatInt(v, 10))
	}
	return out.String()
}

func (i *Int64List) Get() interface{} { return []int64(*i) }

func (i Int64List) Ints() []int64 { return []int64(i) }

func (i *Int64List) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = append(*i, int64(v))
	return err
}

type UintList []uint

func (i *UintList) String() string {
	if i == nil {
		return ""
	}
	var out strings.Builder
	for idx, v := range *i {
		if idx != 0 {
			out.WriteByte(',')
		}
		out.WriteString(strconv.FormatUint(uint64(v), 10))
	}
	return out.String()
}

func (i *UintList) Get() interface{} { return []uint(*i) }

func (i UintList) Uints() []uint { return []uint(i) }

func (i *UintList) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	*i = append(*i, uint(v))
	return err
}

type Uint64List []uint64

func (i *Uint64List) String() string {
	if i == nil {
		return ""
	}
	var out strings.Builder
	for idx, v := range *i {
		if idx != 0 {
			out.WriteByte(',')
		}
		out.WriteString(strconv.FormatUint(v, 10))
	}
	return out.String()
}

func (i *Uint64List) Get() interface{} { return []uint64(*i) }

func (i Uint64List) Uints() []uint64 { return []uint64(i) }

func (i *Uint64List) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = append(*i, uint64(v))
	return err
}

type Float64List []float64

func (i *Float64List) String() string {
	if i == nil {
		return ""
	}
	var out strings.Builder
	for idx, v := range *i {
		if idx != 0 {
			out.WriteByte(',')
		}
		out.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
	}
	return out.String()
}

func (i *Float64List) Get() interface{} { return []float64(*i) }

func (i Float64List) Float64s() []float64 { return []float64(i) }

func (i *Float64List) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*i = append(*i, float64(v))
	return err
}
