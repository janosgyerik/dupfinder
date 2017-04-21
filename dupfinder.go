package dupfinder

import (
	"os"
	"bytes"
	"io"
	"sort"
)

const chunkSize = 64000

func chunker(r io.Reader, ch chan <- []byte) {
	buf := make([]byte, chunkSize)
	for {
		n, err := r.Read(buf)

		if err != nil {
			if err == io.EOF {
				ch <- buf[:n]
			}
			close(ch)
			return
		}

		ch <- buf
	}
}

// TODO improve error handling: files with I/O errors should not be compared to anything, add to black list

// TODO no need to expose publicly
func CompareFiles(path1, path2 string) (int, error) {
	// TODO extract to method: compare by size
	st1, err := os.Stat(path1)
	if err != nil {
		return 0, err
	}

	st2, err := os.Stat(path2)
	if err != nil {
		return 0, err
	}

	if st1.Size() < st2.Size() {
		return -1, nil
	}
	if st1.Size() > st2.Size() {
		return 1, nil
	}

	// TODO extract to method: compare by content
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

// TODO no need to expose publicly
// TODO compare performance of serial and parallel chunkers
// TODO find on internet techniques to process files in parallel correctly
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

type Duplicates struct {
	paths map[string]bool
}

func newDuplicates() Duplicates {
	return Duplicates{make(map[string]bool)}
}

func (duplicates Duplicates) add(path string) {
	duplicates.paths[path] = true
}

func (duplicates Duplicates) count() int {
	return len(duplicates.paths)
}

func (duplicates Duplicates) GetPaths() []string {
	paths := keys(duplicates.paths)
	sort.Strings(paths)
	return paths
}

func keys(m map[string]bool) []string {
	keys := make([]string, 0)
	for key, _ := range m {
		keys = append(keys, key)
	}
	return keys
}

type dupTracker struct {
	pools map[string]Duplicates
}

func newDupTracker() dupTracker {
	return dupTracker{make(map[string]Duplicates)}
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
		pool := newDuplicates()
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
	for _, path := range tracker.getPool(path2).GetPaths() {
		tracker.addToPool(path, pool)
	}
}

func (tracker dupTracker) getPool(path string) Duplicates {
	return tracker.pools[path]
}

// methods to sort Duplicates by item count
type duplicatesList []Duplicates

func (p duplicatesList) Len() int {
	return len(p)
}
func (p duplicatesList) Less(i, j int) bool {
	return p[i].count() < p[j].count()
}
func (p duplicatesList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (tracker dupTracker) getDuplicates() []Duplicates {
	duplicates := make([]Duplicates, 0)
	for _, dups := range tracker.pools {
		duplicates = append(duplicates, dups)
		for _, path := range dups.GetPaths() {
			delete(tracker.pools, path)
		}
	}
	sort.Sort(duplicatesList(duplicates))
	return duplicates
}

func FindDuplicates(paths []string) []Duplicates {
	tracker := newDupTracker()

	mergesort(tracker, paths, 0, len(paths))

	return tracker.getDuplicates()
}

func mergesort(tracker dupTracker, paths []string, low, high int) {
	if low + 1 >= high {
		return
	}

	mid := low + (high - low) / 2
	mergesort(tracker, paths, low, mid)
	mergesort(tracker, paths, mid, high)
	merge(tracker, paths, low, mid, high)
}

func merge(tracker dupTracker, paths []string, low, mid, high int) {
	work := make([]string, 0, high - low)

	var i, j int
	for i, j = low, mid; i < mid && j < high; {
		p1 := paths[i]
		p2 := paths[j]
		cmp, err := CompareFiles(p1, p2)
		if err != nil {
			panic(err)
		}
		if cmp == 0 {
			tracker.add(p1, p2)
			work = append(work, p1, p2)
			i++
			j++
		} else if cmp < 0 {
			work = append(work, p1)
			i++
		} else {
			work = append(work, p2)
			j++
		}
	}

	work = append(work, paths[i:mid]...)
	work = append(work, paths[j:high]...)

	for i = low; i < high; i++ {
		paths[i] = work[i - low]
	}
}
