package cmdy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

var pipeSource = `
package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println(runtime.Version())
	type pipe interface{ Stat() (os.FileInfo, error) }
	for _, v := range []interface{}{os.Stdin, os.Stdout} {
		if pv, ok := v.(pipe); ok {
			fi, _ := pv.Stat()
			fmt.Println((fi.Mode() & os.ModeCharDevice) == 0)
		} else {
			fmt.Println(false)
		}
	}
}
`

func mustParseResult(out []byte) (result struct {
	version                   string
	stdinIsPipe, stdoutIsPipe bool
}) {
	outs := string(out)
	bits := strings.Split(strings.TrimSpace(outs), "\n")
	if len(bits) != 3 {
		panic(fmt.Errorf("unexpected output: %s", outs))
	}

	var err error
	result.version = bits[0]
	result.stdinIsPipe, err = strconv.ParseBool(bits[1])
	if err != nil {
		panic(err)
	}
	result.stdoutIsPipe, err = strconv.ParseBool(bits[2])
	if err != nil {
		panic(err)
	}
	return result
}

func writeTemp(data string) string {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	f.Close()
	fn := f.Name() + ".go"
	if err := os.Rename(f.Name(), fn); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(fn, []byte(data), 0600); err != nil {
		panic(err)
	}
	return fn
}

func TestIsPipeWithBash(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}

	tt := assert.WrapTB(t)
	fn := writeTemp(pipeSource)
	defer os.Remove(fn)

	cmd := func() *exec.Cmd {
		return exec.Command("bash", "-c", fmt.Sprintf("go run %q", fn))
	}

	// FIXME: got a start here, but running bash in this way doesn't seem to be able to
	// pretend the output is a terminal on macOS, even after I follow these suggestions:
	// https://stackoverflow.com/questions/32910661/pretend-to-be-a-tty-in-bash-for-any-command
	//
	// This means we can currently only test if the input is a pipe, not the output.
	{
		out, err := cmd().CombinedOutput()
		tt.MustOK(err)
		result := mustParseResult(out)
		tt.MustEqual(runtime.Version(), result.version)
		tt.MustEqual(false, result.stdinIsPipe)
	}

	{
		c := cmd()
		c.Stdin = bytes.NewReader([]byte("yep"))
		out, err := c.CombinedOutput()
		tt.MustOK(err)
		result := mustParseResult(out)
		tt.MustEqual(true, result.stdinIsPipe)
	}
}
