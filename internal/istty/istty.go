package istty

import (
	"os"
	"runtime"
)

type TTYState int

const (
	IsUnknown TTYState = 0
	IsTTY     TTYState = 1
	IsPipe    TTYState = 2
)

// isTTY tries to only return true if the input is actually a terminal.
//
// os.ModeCharDevice isn't quite enough. I tested various OS/shell combos, using the
// following code:
//
//  fi, _ := os.Stdin.Stat()
//  fmt.Println(fi.Mode())
//
// Results:
//
//                  TTY in       Pipe in      TTY out      Pipe out
//  macOS/bash      Dcrw--w----  -prw-rw----  Dcrw--w----  -prw-rw----
//  macOS/zsh       Dcrw--w----  -prw-rw----  Dcrw--w----  -prw-rw----
//  linux/bash      Dcrw--w----  -prw-------  Dcrw--w----  -prw-------
//  linux/zsh       Dcrw--w----  -prw-------  Dcrw--w----  -prw-------
//  linux/fish      Dcrw--w----  -prw-------  Dcrw--w----  -prw-------
//  win/cmd.exe     Dcrw-rw-rw-  --rw-rw-rw-  Dcrw-rw-rw-  --rw-rw-rw-
//  win/powershell  Dcrw-rw-rw-  --rw-rw-rw-  Dcrw-rw-rw-  --rw-rw-rw-
//  win/gitbash     -prw-rw-rw-  -prw-rw-rw-  -prw-rw-rw-  -prw-rw-rw-
//  win/cygwin      -prw-rw-rw-  -prw-rw-rw-  -prw-rw-rw-  -prw-rw-rw-
//
// More info:
//
//	- https://rosettacode.org/wiki/Check_output_device_is_a_terminal
//	- https://github.com/golang/crypto/blob/master/ssh/terminal/util.go#L29
//	- https://github.com/k-takata/go-iscygpty
//
func CheckTTY(v interface{}) TTYState {
	if v == nil {
		return IsUnknown
	}

	if runtime.GOOS == "windows" {
		if fdv, ok := v.(*os.File); ok {
			fd := fdv.Fd()
			IsCygwinPty(fd)
		}
	}

	type pipe interface{ Stat() (os.FileInfo, error) }
	if pv, ok := v.(pipe); ok {
		fi, _ := pv.Stat()
		if (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return IsTTY
		} else {
			return IsPipe
		}

	} else {
		return IsUnknown
	}
}
