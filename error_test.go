package cmdy

import (
	"errors"
	"testing"

	"github.com/shabbyrobe/cmdy/internal/assert"
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
