package main

import (
	"flag"
	"os"
	"fmt"
	"bufio"
	"github.com/janosgyerik/dupfinder"
)

// TODO
// take list of files from stdin
// print out duplicates visually grouped

func exit() {
	flag.Usage()
	os.Exit(1)
}

type Params struct {
	paths []string
}

func parseArgs() Params {
	flag.Usage = func() {
		fmt.Printf("Usage: find . -type f | %s\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	paths := make([]string, 0)
	if len(flag.Args()) > 0 {
		for _, arg := range flag.Args() {
			if isFile(arg) {
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
	}
}

func readFilePathsFromStdin() []string {
	paths := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if path := scanner.Text(); isFile(path) {
			paths = append(paths, path)
		}
	}

	return paths
}

func isFile(s string) bool {
	if stat, err := os.Stat(s); err == nil && !stat.IsDir() {
		return true
	}
	return false
}

func main() {
	params := parseArgs()

	for _, dups := range dupfinder.FindDuplicates(params.paths...) {
		for _, path := range dups.GetPaths() {
			fmt.Println(path)
		}
		fmt.Println()
	}
}
