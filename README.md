dupfinder
=========

Find duplicate files in specified directory trees.

[![GoDoc](https://godoc.org/github.com/janosgyerik/dupfinder?status.svg)](https://godoc.org/github.com/janosgyerik/dupfinder)
[![Build Status](https://travis-ci.org/janosgyerik/dupfinder.svg?branch=master)](https://travis-ci.org/janosgyerik/dupfinder)

Download
--------

Binaries for several platforms are available in GitHub releases:

https://github.com/janosgyerik/dupfinder/releases

Usage
-----

Find duplicate files in some directory tree:

    dupfinder path/to/dir

Some basic filtering options are available. See `dupfinder -h` for options.

For maximum control, you can use the `find` command to filter files to include,
and pass the list of files to stdin of `dupfinder -0`, for example:

    find . -type f -print0 | dupfinder -0

To find duplicate files in multiple directory trees,
only considering filenames with extension `.avi`,
descending to at most 2 sub-directory levels:

    find path/to/dir path/to/other/dir -name '*.avi' -maxdepth 2 | dupfinder -0 

Generate test coverage report
-----------------------------

Run the `./coverage.sh` helper script.
