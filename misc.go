package cmdy

import "os"

// IsPipe probably returns true if the input is receiving piped data
// from another program, rather than from a terminal.
//
// IsPipe is not compatible with Runner.Stdin, though that may change.
//
// This may not work on Windows.
//
// Typical usage:
//
//	if IsPipe(os.Stdin) {
//		// ...
//	}
//
func IsPipe(f interface{ Stat() (os.FileInfo, error) }) bool {
	if f == nil {
		return false
	}
	fi, _ := f.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}
