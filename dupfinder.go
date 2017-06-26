package dupfinder

import (
	"os"
	"bytes"
	"io"
	"sort"
	"fmt"
	"strings"
)

const chunkSize = 64000

func chunker(r io.Reader, ch chan <- []byte, ready chan bool) {
	buf := make([]byte, chunkSize)
	for <- ready {
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

type cmpResult struct {
	cmp       int
	errFirst  error
	errSecond error
	done      bool
	success   bool
}

func done(cmp int) cmpResult {
	return cmpResult{cmp: cmp, done: true, success: true}
}

func undecided() cmpResult {
	return cmpResult{}
}

func errFirst(err error) cmpResult {
	return cmpResult{errFirst: err, done: true}
}

func errSecond(err error) cmpResult {
	return cmpResult{errSecond: err, done: true}
}

func compare(a, b int64) int {
	if a < b {
		return -1
	}
	if b < a {
		return 1
	}
	return 0
}

func compareFilesBySize(path1, path2 string) cmpResult {
	st1, err := os.Stat(path1)
	if err != nil {
		return errFirst(err)
	}

	st2, err := os.Stat(path2)
	if err != nil {
		return errSecond(err)
	}

	cmp := compare(st1.Size(), st2.Size())
	if cmp != 0 {
		return done(cmp)
	}

	return undecided()
}

func compareFilesByContent(path1, path2 string) cmpResult {
	fd1, err := os.Open(path1)
	defer fd1.Close()
	if err != nil {
		return errFirst(err)
	}

	fd2, err := os.Open(path2)
	defer fd2.Close()
	if err != nil {
		return errSecond(err)
	}

	return compareReaders(fd1, fd2)
}

func compareFiles(path1, path2 string) cmpResult {
	r := compareFilesBySize(path1, path2)
	if r.done {
		return r
	}

	return compareFilesByContent(path1, path2)
}

// TODO compare performance of serial and parallel chunkers
// TODO find on internet techniques to process files in parallel correctly
func compareReaders(fd1, fd2 io.Reader) cmpResult {
	ch1 := make(chan []byte)
	ch2 := make(chan []byte)

	ready1 := make(chan bool)
	defer close(ready1)

	ready2 := make(chan bool)
	defer close(ready2)

	go chunker(fd1, ch1, ready1)
	go chunker(fd2, ch2, ready2)

	for {
		ready1 <- true
		ready2 <- true

		buf1 := <-ch1
		buf2 := <-ch2

		cmp := bytes.Compare(buf1, buf2)
		if cmp != 0 {
			return done(cmp)
		}

		if len(buf1) == 0 {
			return done(cmp)
		}
	}
}

type DupGroup struct {
	paths map[string]bool
}

func newDupGroup() DupGroup {
	return DupGroup{make(map[string]bool)}
}

func (gr DupGroup) add(path string) {
	gr.paths[path] = true
}

func (gr DupGroup) count() int {
	return len(gr.paths)
}

func (gr DupGroup) GetPaths() []string {
	paths := keys(gr.paths)
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
	pools    map[string]DupGroup
	failures []Failure
}

func newDupTracker() *dupTracker {
	return &dupTracker{make(map[string]DupGroup), make([]Failure, 0)}
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
		pool := newDupGroup()
		tracker.addToPool(path1, pool)
		tracker.addToPool(path2, pool)
	}
}

func (tracker dupTracker) addToPool(path string, pool DupGroup) {
	pool.add(path)
	tracker.pools[path] = pool
}

func (tracker dupTracker) mergePools(path1, path2 string) {
	pool := tracker.getPool(path1)
	for _, path := range tracker.getPool(path2).GetPaths() {
		tracker.addToPool(path, pool)
	}
}

func (tracker dupTracker) getPool(path string) DupGroup {
	return tracker.pools[path]
}

func (tracker *dupTracker) err(path string, err error) {
	tracker.failures = append(tracker.failures, Failure{path, err})
}

// methods to sort Duplicates by item count
type duplicatesList []DupGroup

func (p duplicatesList) Len() int {
	return len(p)
}
func (p duplicatesList) Less(i, j int) bool {
	return p[i].count() < p[j].count()
}
func (p duplicatesList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (tracker dupTracker) getDupGroups() []DupGroup {
	duplicates := make([]DupGroup, 0)
	for _, pool := range tracker.pools {
		duplicates = append(duplicates, pool)
		for _, path := range pool.GetPaths() {
			delete(tracker.pools, path)
		}
	}
	sort.Sort(duplicatesList(duplicates))
	return duplicates
}

type failuresList []Failure

func (p failuresList) Len() int {
	return len(p)
}
func (p failuresList) Less(i, j int) bool {
	return strings.Compare(p[i].Path, p[j].Path) < 0
}
func (p failuresList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (tracker dupTracker) getFailures() []Failure {
	failures := append(make([]Failure, 0), tracker.failures...)
	sort.Sort(failuresList(failures))
	return failures
}

type Failure struct {
	Path  string
	Error error
}

type Result struct {
	Groups   []DupGroup
	Failures []Failure
}

func FindDuplicates(paths []string) Result {
	tracker := newDupTracker()

	mergesort(tracker, paths, 0, len(paths))

	return Result{
		Groups: tracker.getDupGroups(),
		Failures: tracker.getFailures(),
	}
}

func mergesort(tracker *dupTracker, paths []string, low, high int) {
	if low + 1 >= high {
		return
	}

	mid := low + (high - low) / 2
	mergesort(tracker, paths, low, mid)
	mergesort(tracker, paths, mid, high)
	merge(tracker, paths, low, mid, high)
}

func merge(tracker *dupTracker, paths []string, low, mid, high int) {
	work := make([]string, 0, high - low)

	var i, j int
	for i, j = low, mid; i < mid && j < high; {
		p1 := paths[i]
		p2 := paths[j]
		r := compareFiles(p1, p2)
		if r.success {
			cmp := r.cmp
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
		} else {
			work = append(work, "")
			if r.errFirst != nil {
				tracker.err(paths[i], r.errFirst)
				i++
			} else if r.errSecond != nil {
				tracker.err(paths[j], r.errSecond)
				j++
			} else {
				panic(fmt.Sprintf("illegal state for cmpResult: %#v", r))
			}
		}
	}

	work = append(work, paths[i:mid]...)
	work = append(work, paths[j:high]...)

	for i = low; i < high; i++ {
		paths[i] = work[i - low]
	}
}
