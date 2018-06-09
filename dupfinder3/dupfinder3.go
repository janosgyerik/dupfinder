package dupfinder3

import (
	"fmt"
	"os"
	"sort"
	"io"
	"reflect"
)

const chunkSize = 4096

type Tracker interface {
	Add(path string)
	Dups() [][]string
	SetLogger(Logger)
}

type FileItem struct {
	Path string
	Size int64
}

func FileSize(path string) int64 {
	fi, e := os.Stat(path)
	if e != nil {
		panic(e)
	}
	return fi.Size()
}

func newFileItem(path string) *FileItem {
	return &FileItem{path, FileSize(path)}
}

type group struct {
	items   []*FileItem
	tracker *tracker
}

func (g *group) add(item *FileItem) {
	g.items = append(g.items, item)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (g *group) fits(item *FileItem) bool {
	p1 := g.items[0].Path
	p2 := item.Path

	f1, err := os.Open(p1)
	check(err)
	defer f1.Close()

	f2, err := os.Open(p2)
	check(err)
	defer f2.Close()

	buf1 := make([]byte, chunkSize)
	buf2 := make([]byte, chunkSize)

	for {
		n1, err1 := f1.Read(buf1)
		n2, err2 := f2.Read(buf2)

		g.tracker.logger.BytesRead(n1 + n2)

		if n1 != n2 {
			return false
		}

		if n1 == 0 {
			return err1 == io.EOF && err2 == io.EOF
		}

		if !reflect.DeepEqual(buf1, buf2) {
			return false
		}
	}
}

func (g *group) String() string {
	return fmt.Sprintf("%v", g.items)
}

func newGroup(t *tracker, item *FileItem) *group {
	g := &group{tracker: t}
	g.add(item)
	return g
}

type Logger interface {
	NewDuplicate([]*FileItem, *FileItem)
	BytesRead(count int)
}

type nullLogger struct{}

func (logger *nullLogger) NewDuplicate([]*FileItem, *FileItem) {}

func (logger *nullLogger) BytesRead(int) {}

type tracker struct {
	groups      []*group
	indexBySize map[int64][]*group
	logger      Logger
}

func (t *tracker) Add(path string) {
	item := newFileItem(path)

	for _, g := range t.indexBySize[item.Size] {
		if g.fits(item) {
			t.logger.NewDuplicate(g.items, item)
			g.add(item)
			return
		}
	}

	group := newGroup(t, item)
	t.groups = append(t.groups, group)
	t.indexBySize[item.Size] = append(t.indexBySize[item.Size], group)
}

type byPath []string

func (a byPath) Len() int           { return len(a) }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPath) Less(i, j int) bool { return a[i] < a[j] }

type bySizeAndFirstPath [][]string

func (a bySizeAndFirstPath) Len() int      { return len(a) }
func (a bySizeAndFirstPath) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bySizeAndFirstPath) Less(i, j int) bool {
	s1 := FileSize(a[i][0])
	s2 := FileSize(a[j][0])
	if s1 < s2 {
		return true
	}
	if s1 > s2 {
		return false
	}
	return a[i][0] < a[j][0]
}

func (t *tracker) Dups() [][]string {
	dups := make([][]string, 0)
	for _, g := range t.groups {
		if len(g.items) > 1 {
			paths := make([]string, 0)
			for _, item := range g.items {
				paths = append(paths, item.Path)
			}
			sort.Sort(byPath(paths))
			dups = append(dups, paths)
		}
	}
	sort.Sort(bySizeAndFirstPath(dups))
	return dups
}

func (t *tracker) SetLogger(logger Logger) {
	t.logger = logger
}

func NewTracker() Tracker {
	t := &tracker{}
	t.indexBySize = make(map[int64][]*group)
	t.logger = &nullLogger{}
	return t
}
