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

func ReadPathsFromLines(reader io.Reader) []string {
	return readFilePaths(reader, scanNullDelimited)
}

func ReadPathsFromNullDelimited(reader io.Reader) []string {
	return readFilePaths(reader, bufio.ScanLines)
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
