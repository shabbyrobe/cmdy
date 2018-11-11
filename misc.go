package cmdy

import (
	"io"
	"os"
)

// ReaderIsPipe probably returns true if the input is receiving piped data
// from another program, rather than from a terminal.
//
// This may not work on Windows.
//
// Typical usage:
//
//	if ReaderIsPipe(os.Stdin) {
//		// Using stdin directly
//	}
//
//	if ReaderIsPipe(ctx.Stdin()) {
//		// Using cmdy.Context
//	}
//
func ReaderIsPipe(in io.Reader) bool {
	return isPipe(in)
}

// WriterIsPipe probably returns true if the Writer represents a pipe to
// another program, rather than to a terminal.
//
// This may not work on Windows.
//
// Typical usage:
//
//	if IsWriterPipe(os.Stdout) {
//		// Using stdout directly
//	}
//
//	if IsWriterPipe(ctx.Stdin()) {
//		// Using cmdy.Context
//	}
//
func WriterIsPipe(out io.Writer) bool {
	return isPipe(out)
}

func isPipe(v interface{}) bool {
	if v == nil {
		return false
	}
	type pipe interface{ Stat() (os.FileInfo, error) }
	if pv, ok := v.(pipe); ok {
		fi, _ := pv.Stat()
		return (fi.Mode() & os.ModeCharDevice) == 0
	} else {
		return false
	}
}
