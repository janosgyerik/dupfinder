package dupfinder2

import "testing"

func Test_find_no_groups_from_two_distinct(t *testing.T) {
	tracker := NewTracker()
	tracker.Add(newItem(1))
	tracker.Add(newItem(2))

	if len(tracker.Dups()) != 0 {
		t.Fatal("expected no duplicates")
	}
}

func Test_find_a_group_from_two_equal(t *testing.T) {
	tracker := NewTracker()
	tracker.Add(newItem(1))
	tracker.Add(newItem(1))

	if len(tracker.Dups()) != 1 {
		t.Fatal("expected 1 group of duplicates")
	}
}

func Test_find_two_groups(t *testing.T) {
	tracker := NewTracker()
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

type testItem struct {
	id int
}

func (t *testItem) Equals(other Item) bool {
	return t.id == other.(*testItem).id
}

func newItem(id int) Item {
	return &testItem{id}
}
