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
}

type fileItem struct {
	path string
	size int64
}

func newFileItem(path string) *fileItem {
	fi, e := os.Stat(path)
	if e != nil {
		panic(e)
	}
	return &fileItem{path, fi.Size()}
}

type group struct {
	items []*fileItem
}

func (g *group) add(item *fileItem) {
	g.items = append(g.items, item)
}

func (g *group) fits(item *fileItem) bool {
	p1 := g.items[0].path
	p2 := item.path

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

func newGroup(item *fileItem) *group {
	g := &group{}
	g.add(item)
	return g
}

type tracker struct {
	groups      []*group
	indexBySize map[int64][]*group
}

func (t *tracker) Add(path string) {
	item := newFileItem(path)

	for _, g := range t.indexBySize[item.size] {
		if g.fits(item) {
			g.add(item)
			return
		}
	}

	group := newGroup(item)
	t.groups = append(t.groups, group)
	t.indexBySize[item.size] = append(t.indexBySize[item.size], group)
}

type byPath []string

func (a byPath) Len() int           { return len(a) }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPath) Less(i, j int) bool { return a[i] < a[j] }

type byFirstPath [][]string

func (a byFirstPath) Len() int           { return len(a) }
func (a byFirstPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFirstPath) Less(i, j int) bool { return a[i][0] < a[j][0] }

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
	sort.Sort(byFirstPath(dups))
	return dups
}

func NewTracker() Tracker {
	t := &tracker{}
	t.indexBySize = make(map[int64][]*group)
	return t
}
