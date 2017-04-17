#!/usr/bin/env bash

cd "$(dirname "$0")"

basename=dupfinder

v=$1; shift
test "$v" || v=v0

if test $# = 0; then
    set -- darwin linux windows
fi

arch=amd64
cli=cmd/$basename/main.go

mkdir -p build

build() {
    GOARCH=$arch go build -o build/$basename-$v-${GOOS}_$arch $cli
}

set -x

for os; do
    GOOS=$os build
done
