package dupfinder

import (
	"os"
	"bytes"
	"io"
)

const chunkSize = 64000

func chunker(r io.Reader, ch chan <- []byte) {
	for {
		buf := make([]byte, chunkSize)
		_, err := r.Read(buf)

		if err != nil {
			close(ch)
			return
		}

		ch <- buf
	}
}

func CompareFiles(path1, path2 string) (int, error) {
	fd1, err := os.Open(path1)
	defer fd1.Close()
	if err != nil {
		return 0, err
	}

	fd2, err := os.Open(path2)
	defer fd2.Close()
	if err != nil {
		return 0, err
	}

	return CompareReaders(fd1, fd2)
}

func CompareReaders(fd1, fd2 io.Reader) (int, error) {
	ch1 := make(chan []byte)
	ch2 := make(chan []byte)

	go chunker(fd1, ch1)
	go chunker(fd2, ch2)

	for {
		buf1 := <-ch1
		buf2 := <-ch2

		cmp := bytes.Compare(buf1, buf2)
		if cmp != 0 {
			return cmp, nil
		}

		if len(buf1) == 0 {
			return 0, nil
		}
	}
}

type Duplicates map[string]bool
//struct {
//	paths map[string]bool
//}

func (duplicates Duplicates) add(path string) {
	//duplicates.paths[path] = true
	duplicates[path] = true
}

type dupTracker struct {
	pools map[string]Duplicates
}

func (tracker dupTracker) add(path1, path2 string) {
	pool1, ok1 := tracker.pools[path1]
	pool2, ok2 := tracker.pools[path2]

	if ok1 && ok2 {
		tracker.mergePools(path1, path2)
	} else if ok1 {
		tracker.addToPool(path2, pool1)
	} else if ok2 {
		tracker.addToPool(path1, pool2)
	} else {
		//pool := make(Duplicates)
		pool := make(map[string]bool)
		//pool := Duplicates{}
		tracker.addToPool(path1, pool)
		tracker.addToPool(path2, pool)
	}
}

func (tracker dupTracker) addToPool(path string, pool Duplicates) {
	pool.add(path)
	tracker.pools[path] = pool
}

func (tracker dupTracker) mergePools(path1, path2 string) {
	pool := tracker.getPool(path1)
	for path, _ := range tracker.getPool(path2) {
		tracker.addToPool(path, pool)
	}
}

func (tracker dupTracker) getPool(path string) Duplicates {
	return tracker.pools[path]
}

func (tracker dupTracker) getDuplicates() []Duplicates {
	duplicates := make([]Duplicates, 0)
	for _, dups := range tracker.pools {
		duplicates = append(duplicates, dups)
		for path, _ := range dups {
			delete(tracker.pools, path)
		}
	}
	return duplicates
}

func FindDuplicates(paths... string) []Duplicates {
	tracker := dupTracker{make(map[string]Duplicates)}

	// naive brute-force implementation: compare all files against all
	for i := 0; i < len(paths) - 1; i++ {
		for j := i + 1; j < len(paths); j++ {
			path1 := paths[i]
			path2 := paths[j]
			if cmp, err := CompareFiles(path1, path2); err == nil && cmp == 0 {
				tracker.add(path1, path2)
			}
		}
	}

	return tracker.getDuplicates()
}
