package cmdyutil

import (
	"testing"

	"github.com/shabbyrobe/cmdy/internal/assert"
)

func TestConfig(t *testing.T) {
	tt := assert.WrapTB(t)

	in := `
		[set1] -foo -bar
		[set2] -baz -qux
		[set3]
	`
	conf, err := ParseConfig([]byte(in))
	tt.MustOK(err)
	tt.MustEqual(
		[]ConfigSet{
			{"set1", []string{"-foo", "-bar"}},
			{"set2", []string{"-baz", "-qux"}},
			{"set3", nil},
		},
		conf.Sets())
}

func TestConfigComments(t *testing.T) {
	tt := assert.WrapTB(t)

	in := `
		# foo
		[set1] -foo
		# bar
		-bar
		# baz
		[set2] 
		# qux
		-baz -qux # biz
		[set3] # baz
		# wuz
	`
	conf, err := ParseConfig([]byte(in))
	tt.MustOK(err)
	tt.MustEqual(
		[]ConfigSet{
			{"set1", []string{"-foo", "-bar"}},
			{"set2", []string{"-baz", "-qux"}},
			{"set3", nil},
		},
		conf.Sets())
}
