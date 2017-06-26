package main

import (
	"flag"
	"os"
	"fmt"
	"bufio"
	"github.com/janosgyerik/dupfinder"
	"github.com/janosgyerik/dupfinder/finder"
	"bytes"
	"io"
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

func scanNullDelimited(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0); i >= 0 {
		// We have a full null-terminated line.
		return i + 1, dropNull(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropNull(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropNull drops a terminal \0 from the data.
func dropNull(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == 0 {
		return data[0 : len(data)-1]
	}
	return data
}

func parseArgs() Params {
	minSizePtr := flag.Int64("minSize", 1, "minimum file size")
	stdinPtr := flag.Bool("stdin", false, "read paths from stdin")
	zeroPtr := flag.Bool("0", false, "read paths from stdin, null-delimited")

	flag.Parse()

	paths := make([]string, 0)
	if *zeroPtr {
		paths = readFilePaths(os.Stdin, scanNullDelimited)
	} else if *stdinPtr {
		paths = readFilePaths(os.Stdin, bufio.ScanLines)
	} else if len(flag.Args()) > 0 {
		for _, arg := range flag.Args() {
			if isFileOrDir(arg) {
				paths = append(paths, arg)
			}
		}
	}

	if len(paths) == 0 {
		exit()
	}

	return Params{
		paths: paths,
		minSize: *minSizePtr,
	}
}

func readFilePaths(reader io.Reader, splitter bufio.SplitFunc) []string {
	paths := make([]string, 0)

	scanner := bufio.NewScanner(reader)
	scanner.Split(splitter)
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
