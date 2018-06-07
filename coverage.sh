#!/bin/sh

cd "$(dirname "$0")"
mkdir -p tmp

rm -f tmp/cover-*.out
for package in . ./finder ./pathreader; do
    name=$(cd "$package"; basename "$PWD")
    go test $package -coverprofile tmp/cover-$name.out
done

{
    echo "mode: set"
    grep -hv ^mode: tmp/cover-*.out
} > tmp/cover.out
go tool cover -html=tmp/cover.out -o tmp/cover.html
