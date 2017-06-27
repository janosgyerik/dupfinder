package pathreader

import (
	"bytes"
	"io"
	"bufio"
	"os"
)

func scanNullDelimited(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0); i >= 0 {
		// We have a full null-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

type filter func(string) bool

func readItems(reader io.Reader, splitter bufio.SplitFunc, filter filter) []string {
	items := make([]string, 0)

	scanner := bufio.NewScanner(reader)
	scanner.Split(splitter)
	for scanner.Scan() {
		if item := scanner.Text(); filter(item) {
			items = append(items, item)
		}
	}

	return items
}

func readItemsFromLines(reader io.Reader, filter filter) []string {
	return readItems(reader, bufio.ScanLines, filter)
}

func readItemsFromNullDelimited(reader io.Reader, filter filter) []string {
	return readItems(reader, scanNullDelimited, filter)
}

func readFilePaths(reader io.Reader, splitter bufio.SplitFunc) []string {
	return readItems(reader, splitter, isFileOrDir)
}

func isFileOrDir(s string) bool {
	_, err := os.Stat(s)
	return err == nil
}

func ReadPathsFromLines(reader io.Reader) []string {
	return readFilePaths(reader, bufio.ScanLines)
}

func ReadPathsFromNullDelimited(reader io.Reader) []string {
	return readFilePaths(reader, scanNullDelimited)
}

func FilterPaths(args []string) []string {
	paths := make([]string, 0)
	for _, arg := range args {
		if isFileOrDir(arg) {
			paths = append(paths, arg)
		}
	}
	return paths
}
