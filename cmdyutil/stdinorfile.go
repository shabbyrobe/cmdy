package cmdyutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ArtProcessors/cmdy"
)

// OpenStdinOrFile will check if the program's input is a pipe. If so, it will
// return stdin, otherwise it will return an open file.
//
// The returned ReadCloser must always be closed.
func OpenStdinOrFile(ctx cmdy.Context, fileName string) (rdr io.ReadCloser, err error) {
	// FIXME: ReaderIsPipe is untestable when using BufferedRunner:
	input := ctx.Stdin()
	hasInput := false
	if buf, ok := input.(*bytes.Buffer); ok {
		hasInput = buf.Len() > 0
	} else {
		hasInput = cmdy.ReaderIsPipe(input)
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
func ReadStdinOrFile(ctx cmdy.Context, fileName string) (bts []byte, err error) {
	rdr, err := OpenStdinOrFile(ctx, fileName)
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
