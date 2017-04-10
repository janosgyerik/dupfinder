package dupfinder

import (
	"testing"
	"reflect"
	"strconv"
)

var testFileCounter = 0

func newIndex() Index {
	return NewIndex()
}

func Test_tracker_add(t*testing.T) {
	tracker := NewTracker()

	if len(tracker.Groups()) != 0 {
		t.Fatalf("got %#v groups from empty tracker", tracker.Groups())
	}

	f1 := newTestFile()
	f2 := newTestFile()
	tracker.Add(f1, f2)

	if len(tracker.Groups()) != 1 {
		t.Fatalf("got %d groups from tracker; expected 1", len(tracker.Groups()))
	}

	if x := len(tracker.Groups()[0].Paths); x != 2 {
		t.Fatalf("got %d files in group; expected 2", x)
	}

	f3 := newTestFile()
	tracker.Add(f1, f3)

	if len(tracker.Groups()) != 1 {
		t.Fatalf("got %d groups from tracker; expected 1", len(tracker.Groups()))
	}

	if x := len(tracker.Groups()[0].Paths); x != 3 {
		t.Fatalf("got %d files in group; expected 3", x)
	}

	f4 := newTestFile()
	f5 := newTestFile()
	tracker.Add(f4, f5)

	if len(tracker.Groups()) != 2 {
		t.Fatalf("got %d groups from tracker; expected 2", len(tracker.Groups()))
	}
}

type testFile struct {
	id      string
	size    int
	digest  string
	content string
}

func newTestFile() testFile {
	testFileCounter += 1

	return testFile{
		id: strconv.Itoa(testFileCounter),
	}
}

func (f testFile) Id() string {
	return f.id
}

func (f testFile) Size() int {
	return f.size
}

func (f testFile) Digest() string {
	return f.digest
}

func (f testFile) Content() string {
	return f.content
}

func fileWithSize(size int) FileHandler {
	return testFile{size: size}
}

func fileWithDigest(digest string) FileHandler {
	return testFile{digest: digest}
}

func Test_fileWithSize(t*testing.T) {
	f1_1 := fileWithSize(1)
	f1_2 := fileWithSize(1)
	f2 := fileWithSize(2)

	if f1_1.Size() != f1_2.Size() {
		t.Fatal("got different size for two fileWithSize(1)")
	}
	if f1_1.Size() == f2.Size() {
		t.Fatal("got same size for fileWithSize(1) and fileWithSize(2)")
	}
}

func Test_fileWithDigest(t*testing.T) {
	f1_1 := fileWithDigest("1")
	f1_2 := fileWithDigest("1")
	f2 := fileWithDigest("2")

	if f1_1.Digest() != f1_2.Digest() {
		t.Fatal("got different digest for two fileWithDigest(1)")
	}
	if f1_1.Digest() == f2.Digest() {
		t.Fatal("got same digest for fileWithDigest(1) and fileWithDigest(2)")
	}
}

func Test_should_not_find_duplicates_with_different_size(t*testing.T) {
	index := newIndex()

	index.Add(fileWithSize(1))
	index.Add(fileWithSize(2))

	if len(index.Groups()) > 0 {
		t.Fatal("found duplicates in two files with different size")
	}
}

func Test_should_not_find_duplicates_with_different_digest(t*testing.T) {
	index := newIndex()

	index.Add(fileWithDigest("1"))
	index.Add(fileWithDigest("2"))

	if len(index.Groups()) > 0 {
		t.Fatal("found duplicates in two files with different digest")
	}
}

func Test_should_find_duplicates_in_two_identical_files(t*testing.T) {
	index := newIndex()

	file := testFile{
		size: 1,
		digest: "1",
	}
	index.Add(file)
	index.Add(file)

	if len(index.Groups()) != 1 {
		t.Fatal("did not find duplicates in two identical files")
	}
}

func Test_should_find_duplicates_in_a_mix(t*testing.T) {
	index := newIndex()

	file1 := testFile{size: 1, digest: "1", content: "1"}
	file2 := testFile{size: 1, digest: "1", content: "1"}
	file3 := testFile{size: 1, digest: "1", content: "3"}

	index.Add(file1)
	index.Add(file2)
	index.Add(file3)

	if len(index.Groups()) != 1 {
		t.Fatal("did not find duplicates in a mix of files")
	}

	actual := index.Groups()[0].Paths
	expected := []string{file1.Id(), file2.Id()}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got %#v in duplicate group, expected %#v", actual, expected)
	}
}
