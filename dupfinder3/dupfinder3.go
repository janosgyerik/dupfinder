package dupfinder3

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
)

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
	items []*FileItem
}

func (g *group) add(item *FileItem) {
	g.items = append(g.items, item)
}

func (g *group) fits(item *FileItem) bool {
	p1 := g.items[0].Path
	p2 := item.Path

	s1, err1 := ioutil.ReadFile(p1)
	if err1 != nil {
		panic("could not read file: " + p1)
	}

	s2, err2 := ioutil.ReadFile(p2)
	if err2 != nil {
		panic("could not read file: " + p2)
	}

	return reflect.DeepEqual(s1, s2)
}

func (g *group) String() string {
	return fmt.Sprintf("%v", g.items)
}

func newGroup(item *FileItem) *group {
	g := &group{}
	g.add(item)
	return g
}

type Logger interface {
	NewDuplicate([]*FileItem, *FileItem)
}

type nullLogger struct{}

func (logger *nullLogger) NewDuplicate([]*FileItem, *FileItem) {}

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

	group := newGroup(item)
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
