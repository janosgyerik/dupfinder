package main

import (
	"flag"
	"os"
	"fmt"
	"bufio"
	"github.com/janosgyerik/dupfinder"
	"github.com/janosgyerik/dupfinder/finder"
)

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	paths    []string
	minSize  int64
}

func parseArgs() Params {
	minSizePtr := flag.Int64("minSize", 0, "minimum file size")

	flag.Parse()

	paths := make([]string, 0)
	if len(flag.Args()) > 0 {
		for _, arg := range flag.Args() {
			if isFileOrDir(arg) {
				paths = append(paths, arg)
			}
		}
	} else {
		paths = readFilePathsFromStdin()
	}

	if len(paths) == 0 {
		exit()
	}

	return Params{
		paths: paths,
		minSize: *minSizePtr,
	}
}

func readFilePathsFromStdin() []string {
	paths := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if path := scanner.Text(); isFileOrDir(path) {
			paths = append(paths, path)
		}
	}

	return paths
}

func isFileOrDir(s string) bool {
	_, err := os.Stat(s)
	return err == nil
}

func main() {
	params := parseArgs()

	finder := finder.NewFinder(finder.Filters.MinSize(params.minSize))
	index := dupfinder.NewIndex()

	for _, path := range params.paths {
		for filepath := range finder.Find(path) {
			info, _ := os.Stat(filepath)
			index.Add(dupfinder.NewFileHandler(filepath, info))
		}
	}

	for _, dups := range index.Groups() {
		for _, path := range dups.Paths {
			fmt.Println(path)
		}
		fmt.Println()
	}
}
