// +build !windows

package istty

// IsCygwinPty returns true if the file descriptor is a Cygwin/MSYS pty.
func IsCygwinPty(fd uintptr) bool {
	return false
}
