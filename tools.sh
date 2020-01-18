#!/bin/bash
set -o errexit -o nounset -o pipefail

cmd-testcvg() {
    go test -coverprofile=cover.out github.com/shabbyrobe/cmdy/...
    go tool cover -html=cover.out
}

cmd-sloc() {
    tokei --exclude '*_test.go' .
}

"cmd-$1" "${@:2}"

