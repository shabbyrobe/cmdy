package cmdyutil

import (
	"io"

	"github.com/shabbyrobe/cmdy/internal/istty"
)

// ReaderIsPipe probably returns true if the input is receiving piped data
// from another program, rather than from a terminal.
//
// This is known to work in the following environments:
//
// - Bash on macOS and Linux
// - Command Prompt on Windows
// - Windows Powershell
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
	return istty.CheckTTY(in) == istty.IsPipe
}

// WriterIsPipe probably returns true if the Writer represents a pipe to
// another program, rather than to a terminal.
//
// This is known to work in the following environments:
//
// - Bash on macOS and Linux
// - Command Prompt on Windows
// - Windows Powershell
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
	return istty.CheckTTY(out) == istty.IsPipe
}
