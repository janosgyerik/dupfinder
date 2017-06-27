package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/janosgyerik/dupfinder"
	"github.com/janosgyerik/dupfinder/finder"
	"github.com/janosgyerik/dupfinder/pathreader"
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	paths   []string
	minSize int64
	stdin   bool
	stdin0  bool
}

func parseArgs() Params {
	minSizePtr := flag.Int64("minSize", 1, "minimum file size")
	stdinPtr := flag.Bool("stdin", false, "read paths from stdin")
	zeroPtr := flag.Bool("0", false, "read paths from stdin, null-delimited")

	flag.Parse()

	var paths []string
	if *zeroPtr {
		paths = pathreader.ReadPathsFromNullDelimited(os.Stdin)
	} else if *stdinPtr {
		paths = pathreader.ReadPathsFromLines(os.Stdin)
	} else {
		paths = pathreader.FilterPaths(flag.Args())
	}

	if len(paths) == 0 {
		exit()
	}

	return Params{
		paths: paths,
		minSize: *minSizePtr,
	}
}

func main() {
	params := parseArgs()

	filefinder := finder.NewFinder(finder.Filters.MinSize(params.minSize))

	paths := []string{}
	for _, path := range params.paths {
		paths = append(paths, filefinder.FindAll(path)...)
	}

	result := dupfinder.FindDuplicates(paths)

	if len(result.Failures) > 0 {
		fmt.Println("# I/O errors in files:")
		for _, failure := range result.Failures {
			fmt.Printf("# %s\n", failure.Path)
		}
		fmt.Println()
	}

	for _, dup := range result.Groups {
		for _, path := range dup.GetPaths() {
			fmt.Println(path)
		}
		fmt.Println()
	}
}
