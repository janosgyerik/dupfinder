dupfinder
=========

Find duplicate files in specified directory trees.

TODO: not actually implemented yet

[![GoDoc](https://godoc.org/github.com/janosgyerik/dupfinder?status.svg)](https://godoc.org/github.com/janosgyerik/dupfinder)
[![Build Status](https://travis-ci.org/janosgyerik/dupfinder.svg?branch=master)](https://travis-ci.org/janosgyerik/dupfinder)

Usage
-----

To find duplicate files in the current directory and all sub-directories:

    find . -type f | dupfinder

See `dupfinder -h` for all available options.

To find duplicate files in multiple directory trees,
only considering filenames with extension `.avi`,
descending to at most 2 sub-directory levels:

    find path/to/dir path/to/other/dir -name '*.avi' -maxdepth 2 | dupfinder 

Download
--------

TODO: not actually implemented yet

Binaries for several platforms are available on SourceForge:

https://sourceforge.net/projects/dupfinder/files/

Generate test coverage report
-----------------------------

Run the commands:

    go test -coverprofile cover.out
    go tool cover -html=cover.out -o cover.html
    open cover.html
    
Or simply run the `./coverage.sh` helper script.

See more info: https://blog.golang.org/cover
