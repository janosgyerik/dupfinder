package dupfinder2

import (
	"testing"
	"fmt"
)

func Test_find_no_groups_from_two_distinct(t *testing.T) {
	tracker := newTracker()
	tracker.Add(newTestItem(1))
	tracker.Add(newTestItem(2))

	if len(tracker.Dups()) != 0 {
		t.Fatal("expected no duplicates")
	}
}

func Test_find_a_group_from_two_equal(t *testing.T) {
	tracker := newTracker()
	tracker.Add(newTestItem(1))
	tracker.Add(newTestItem(1))

	if len(tracker.Dups()) != 1 {
		t.Fatal("expected 1 group of duplicates")
	}
}

func Test_find_two_groups(t *testing.T) {
	tracker := newTracker()
	tracker.Add(newTestItem(1))
	tracker.Add(newTestItem(1))
	tracker.Add(newTestItem(2))
	tracker.Add(newTestItem(2))
	tracker.Add(newTestItem(2))
	tracker.Add(newTestItem(3))

	if len(tracker.Dups()) != 2 {
		t.Fatal("expected 2 groups of duplicates")
	}
}

func Test_find_two_groups_in_files(t *testing.T) {
	tracker := NewTracker(newFakeFileFilter())
	tracker.Add(newFakeFileItem("foo", 1, "foo"))
	tracker.Add(newFakeFileItem("bar", 1, "foo"))
	tracker.Add(newFakeFileItem("a1", 2, "abc"))
	tracker.Add(newFakeFileItem("a2", 2, "abc"))
	tracker.Add(newFakeFileItem("a3", 2, "abc"))
	tracker.Add(newFakeFileItem("x1", 3, "foo"))
	tracker.Add(newFakeFileItem("x2", 4, "abc"))

	if len(tracker.Dups()) != 2 {
		t.Fatal("expected 2 groups of duplicates")
	}
}

type Key int

type KeyExtractor interface {
	Key(Item) Key
}

type keyExtractor struct {
}

func (k *keyExtractor) Key(item Item) Key {
	return Key(item.(*testItem).id)
}

type testFilter struct {
	byId         map[Key][]Group
	keyExtractor KeyExtractor
}

func (f *testFilter) CandidateGroups(item Item) []Group {
	if g, found := f.byId[f.keyExtractor.Key(item)]; found {
		return g
	}
	return nil
}

func (f *testFilter) Register(item Item, g Group) {
	f.byId[f.keyExtractor.Key(item)] = append(f.byId[f.keyExtractor.Key(item)], g)
}

func newFilter() Filter {
	return &testFilter{byId: make(map[Key][]Group), keyExtractor: &keyExtractor{}}
}

func newTracker() Tracker {
	return NewTracker(newFilter())
}

type testItem struct {
	id int
}

func (t *testItem) Equals(other Item) bool {
	return t.id == other.(*testItem).id
}

func (t *testItem) String() string {
	return fmt.Sprintf("%v", t.id)
}

func newTestItem(id int) Item {
	return &testItem{id}
}

type fakeFileItem struct {
	path    string
	size    int
	content string
}

func (f *fakeFileItem) Equals(other Item) bool {
	return f.content == other.(*fakeFileItem).content
}

func newFakeFileItem(path string, size int, content string) Item {
	return &fakeFileItem{path, size, content}
}

type fakeSizeExtractor struct {
}

func (s *fakeSizeExtractor) Key(item Item) Key {
	return Key(item.(*fakeFileItem).size)
}

func newFakeFileFilter() Filter {
	return &testFilter{byId: make(map[Key][]Group), keyExtractor: &fakeSizeExtractor{}}
}
