#!/usr/bin/env bash

set -eu

cd "$(dirname "$0")"

basename=dupfinder

usage() {
    echo "Usage: $0 [version [platform]]"
    exit 1
}

for arg; do
    case "$arg" in
    -h|--help) usage ;;
    esac
done

if test $# = 0; then
    dev=1
else
    dev=
    version=$1; shift
    test $# = 0 && set -- darwin linux windows
fi

arch=amd64
cli=cmd/$basename/main.go

mkdir -p build

build() {
    GOARCH=$arch go build -o build/$1 $cli
}

set -x
if test "$dev"; then
    build $basename
else
    for os; do
        GOOS=$os build $basename-$version-${os}_$arch
    done
fi
