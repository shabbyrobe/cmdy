package cmdy

import (
	"io"
	"os"
)

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
func IsPipe(in io.Reader) bool {
	if in == nil {
		return false
	}
	type pipe interface{ Stat() (os.FileInfo, error) }
	if inPipe, ok := in.(pipe); ok {
		fi, _ := inPipe.Stat()
		return (fi.Mode() & os.ModeCharDevice) == 0
	} else {
		return false
	}
}
