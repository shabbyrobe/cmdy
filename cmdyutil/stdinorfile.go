package cmdyutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/shabbyrobe/cmdy"
)

type StdinFlag int

const (
	HyphenStdin StdinFlag = 1 << iota
)

// OpenStdinOrFile will check if the program's input is a pipe. If so, it will
// return stdin, otherwise it will return an open file.
//
// NOTE: you should probably use '-' as a filename to trigger stdin (see the
// HyphenStdin flag). If considering stdin by default, your commands might
// not work properly in bash while loops:
//
//	# Potentially bad ('find stuff' goes to 'mycmd' stdin):
//	find stuff | while read line; do mycmd; done
//
//	# Less bad:
//	find stuff | while read line; do </dev/null mycmd; done
//
// If you require '-' to trigger stdin, this won't happen by default.
//
// The returned ReadCloser must always be closed.
func OpenStdinOrFile(ctx cmdy.Context, fileName string, flag StdinFlag) (rdr io.ReadCloser, err error) {
	var input io.Reader
	var hasInput bool

	// FIXME: ReaderIsPipe is untestable when using BufferedRunner:
	if flag&HyphenStdin == 0 || fileName == "-" {
		if flag&HyphenStdin != 0 {
			fileName = ""
		}
		input = ctx.Stdin()
		if buf, ok := input.(*bytes.Buffer); ok {
			hasInput = buf.Len() > 0
		} else {
			hasInput = cmdy.ReaderIsPipe(input)
		}
	}

	if hasInput && fileName != "" {
		return nil, errStdinOrFileBoth
	} else if !hasInput && fileName == "" {
		return nil, errStdinOrFileNeither
	}

	if hasInput {
		return ioutil.NopCloser(input), nil
	} else {
		return os.Open(fileName)
	}
}

// ReadStdinOrFile will check if the program's input is a pipe. If so, it will read from
// stdin, otherwise it will read from fileName.
func ReadStdinOrFile(ctx cmdy.Context, fileName string, flag StdinFlag) (bts []byte, err error) {
	rdr, err := OpenStdinOrFile(ctx, fileName, flag)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()
	return ioutil.ReadAll(rdr)
}

func IsStdinOrFileError(err error) bool {
	return errors.Is(err, errStdinOrFileBoth) || errors.Is(err, errStdinOrFileNeither)
}

var (
	errStdinOrFileBoth    = fmt.Errorf("received file name and STDIN pipe, must be one or the other")
	errStdinOrFileNeither = fmt.Errorf("received neither file name nor STDIN pipe, must be one or the other")
)
