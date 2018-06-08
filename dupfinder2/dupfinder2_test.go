package dupfinder2

import (
	"testing"
	"fmt"
)

func Test_find_no_groups_from_two_distinct(t *testing.T) {
	tracker := newTracker()
	tracker.Add(newItem(1))
	tracker.Add(newItem(2))

	if len(tracker.Dups()) != 0 {
		t.Fatal("expected no duplicates")
	}
}

func Test_find_a_group_from_two_equal(t *testing.T) {
	tracker := newTracker()
	tracker.Add(newItem(1))
	tracker.Add(newItem(1))

	if len(tracker.Dups()) != 1 {
		t.Fatal("expected 1 group of duplicates")
	}
}

func Test_find_two_groups(t *testing.T) {
	tracker := newTracker()
	tracker.Add(newItem(1))
	tracker.Add(newItem(1))
	tracker.Add(newItem(2))
	tracker.Add(newItem(2))
	tracker.Add(newItem(2))
	tracker.Add(newItem(3))

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

func newItem(id int) Item {
	return &testItem{id}
}
