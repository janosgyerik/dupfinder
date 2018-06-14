#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

coverdir=$PWD/tmp

mkdir -p "$coverdir"
rm -f "$coverdir"/cover-*.out

cover() {
    local dir=$1
    local package=$2
    local name=$3
    (cd "$dir" && go test $package -coverprofile "$coverdir"/cover-$name.out)
}

cover . . dupfinder
cover . ./finder finder
cover . ./pathreader pathreader
cover . ./utils utils
cover cmd/dupfinder . cmd

{
    echo "mode: set"
    grep -hv ^mode: tmp/cover-*.out
} > "$coverdir/cover.out"

go tool cover -html="$coverdir/cover.out" -o "$coverdir/cover.html"
