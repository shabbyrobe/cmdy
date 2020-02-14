package cmdy

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ArtProcessors/cmdy/internal/assert"
)

func TestFormatErrorQuiet(t *testing.T) {
	tt := assert.WrapTB(t)
	msg, code := FormatError(QuietExit(0))
	tt.MustEqual("", msg)
	tt.MustEqual(0, code)

	msg, code = FormatError(QuietExit(1))
	tt.MustEqual("", msg)
	tt.MustEqual(1, code)
}

func TestFormatError(t *testing.T) {
	tt := assert.WrapTB(t)

	err := ErrWithCode(3, errors.New("boom"))
	msg, code := FormatError(err)
	tt.MustEqual("boom", msg)
	tt.MustEqual(3, code)
}

func TestFormatErrorWithWrappedUsageError(t *testing.T) {
	tt := assert.WrapTB(t)

	err := UsageError(ErrWithCode(3, errors.New("boom")))
	msg, code := FormatError(err)
	tt.MustEqual("error: boom", msg)
	tt.MustEqual(ExitUsage, code)
}

func TestUsageErrorCanUnwrap(t *testing.T) {
	tt := assert.WrapTB(t)
	err := fmt.Errorf("YEPPO")
	uerr := UsageError(err)
	tt.MustEqual(err, errors.Unwrap(uerr))
}
