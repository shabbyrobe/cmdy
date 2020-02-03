#!/bin/bash
set -o errexit -o nounset -o pipefail

cmd-test() {
    go test -count=1 github.com/shabbyrobe/cmdy/...
}

cmd-testcvg() {
    go test -coverprofile=cover.out github.com/shabbyrobe/cmdy/...
    go tool cover -html=cover.out
}

cmd-sloc() {
    tokei --exclude '*_test.go' .
}

"cmd-$1" "${@:2}"

