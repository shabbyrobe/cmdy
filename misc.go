package cmdy

import "os"

func StdinIsPipe() bool {
	return IsPipe(os.Stdin)
}

func IsPipe(f interface{ Stat() (os.FileInfo, error) }) bool {
	if f == nil {
		return false
	}
	fi, _ := f.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}
