package cmdy

import (
	"io"
	"os"

	"github.com/ArtProcessors/cmdy/internal/istty"
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

// ProgName attempts to guess the program name from the first argument in os.Args.
func ProgName() string {
	if len(os.Args) < 1 {
		return ""
	}
	return baseName(os.Args[0])
}

// baseName is a cut-down remix of filepath.Base that saves us a dependency
// and skips use-cases that we don't need to worry about, like windows volume
// names, etc, because we are only using it to grab the program name.
func baseName(path string) string {
	// Strip trailing slashes.
	for len(path) > 0 && os.IsPathSeparator(path[len(path)-1]) {
		path = path[0 : len(path)-1]
	}
	// Find the last element
	i := len(path) - 1
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	if i >= 0 {
		path = path[i+1:]
	}
	return path
}
