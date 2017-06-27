package pathreader

import (
	"bytes"
	"io"
	"bufio"
	"os"
	"path/filepath"
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

func newUniqueFilter() filter {
	seen := make(map[string]bool)

	return func(s string) bool {
		normalized := filepath.Clean(s)
		if _, ok := seen[normalized]; !ok {
			seen[normalized] = true
			return true
		}
		return false
	}
}

func isFileOrDir(s string) bool {
	_, err := os.Stat(s)
	return err == nil
}

func newDefaultFilter() filter {
	isUnique := newUniqueFilter()

	return func(s string) bool {
		return isFileOrDir(s) && isUnique(s)
	}
}

func readItems(reader io.Reader, splitter bufio.SplitFunc, filter filter) []string {
	items := make([]string, 0)

	scanner := bufio.NewScanner(reader)
	scanner.Split(splitter)
	for scanner.Scan() {
		if item := scanner.Text(); filter(item) {
			items = append(items, normalize(item))
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
	return readItems(reader, splitter, newDefaultFilter())
}

func ReadPathsFromLines(reader io.Reader) []string {
	return readFilePaths(reader, bufio.ScanLines)
}

func ReadPathsFromNullDelimited(reader io.Reader) []string {
	return readFilePaths(reader, scanNullDelimited)
}

func normalize(path string) string {
	return filepath.Clean(path)
}

func FilterPaths(args []string) []string {
	filter := newDefaultFilter()

	paths := make([]string, 0)
	for _, arg := range args {
		if filter(arg) {
			paths = append(paths, normalize(arg))
		}
	}
	return paths
}
