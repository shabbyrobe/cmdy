#!/bin/bash

gotest() {
    go test -coverprofile=cover.out \
        -coverpkg=github.com/shabbyrobe/cmdy,github.com/shabbyrobe/cmdy/usage,github.com/shabbyrobe/cmdy/flags,github.com/shabbyrobe/cmdy/args \
        github.com/shabbyrobe/cmdy/...
}

"$1" "${@:2}"

