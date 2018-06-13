package dupfinder

import (
	"fmt"
	"os"
	"sort"
	"io"
	"reflect"
)

const chunkSize = 4096

type EventListener interface {
	NewDuplicate([]string)
	BytesRead(count int)
}

type nullEventListener struct{}

func (eventListener *nullEventListener) NewDuplicate([]string) {}

func (eventListener *nullEventListener) BytesRead(int) {}

type Tracker interface {
	Add(path string)
	Dups() [][]string
	SetEventListener(EventListener)
}

type fileItem struct {
	path string
	size int64
}

func newFileItem(path string) *fileItem {
	return &fileItem{path, FileSize(path)}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type group struct {
	items   []*fileItem
	paths   []string
	tracker *tracker
}

func (g *group) add(item *fileItem) {
	g.items = append(g.items, item)
	g.paths = append(g.paths, item.path)
}

func (g *group) fits(item *fileItem) bool {
	p1 := g.items[0].path
	p2 := item.path

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

		g.tracker.eventListener.BytesRead(n1 + n2)

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

func newGroup(t *tracker, item *fileItem) *group {
	g := &group{tracker: t}
	g.add(item)
	return g
}

type tracker struct {
	groups        []*group
	indexBySize   map[int64][]*group
	eventListener EventListener
}

func (t *tracker) Add(path string) {
	item := newFileItem(path)

	for _, g := range t.indexBySize[item.size] {
		if g.fits(item) {
			g.add(item)
			t.eventListener.NewDuplicate(g.paths)
			return
		}
	}

	group := newGroup(t, item)
	t.groups = append(t.groups, group)
	t.indexBySize[item.size] = append(t.indexBySize[item.size], group)
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
				paths = append(paths, item.path)
			}
			sort.Sort(byPath(paths))
			dups = append(dups, paths)
		}
	}
	sort.Sort(bySizeAndFirstPath(dups))
	return dups
}

func (t *tracker) SetEventListener(eventListener EventListener) {
	t.eventListener = eventListener
}

func NewTracker() Tracker {
	t := &tracker{}
	t.indexBySize = make(map[int64][]*group)
	t.eventListener = &nullEventListener{}
	return t
}
