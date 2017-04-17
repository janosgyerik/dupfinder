package dupfinder

import (
	"testing"
	"strings"
	"io/ioutil"
	"strconv"
	"path"
	"io"
	"reflect"
)

const testDataDir = "testdata"

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
	path    string
	size    int64
	digest  string
	content string
}

func newTestFile() testFile {
	testFileCounter += 1

	return testFile{
		path: strconv.Itoa(testFileCounter),
	}
}

func (f testFile) Path() string {
	return f.path
}

func (f testFile) Size() int64 {
	return f.size
}

func (f testFile) Digest() string {
	return f.digest
}

func (f testFile) NewReadCloser() (io.ReadCloser, error) {
	return testReadCloser{strings.NewReader(f.content)}, nil
}

type testReadCloser struct {
	reader io.Reader
}

func (t testReadCloser) Read(buf []byte) (int, error) {
	return t.reader.Read(buf)
}

func (t testReadCloser) Close() error {
	return nil
}

func fileWithSize(size int64) FileHandler {
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
	expected := []string{file1.Path(), file2.Path()}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got %#v in duplicate group, expected %#v", actual, expected)
	}
}

func Test_alldups(t*testing.T) {
	basedir := path.Join(testDataDir, "alldups")
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		t.Fatal(err)
	}

	index := newIndex()
	for _, file := range files {
		p := path.Join(basedir, file.Name())
		f := NewFileHandler(p, file)
		index.Add(f)
	}

	if len(index.Groups()) != 1 {
		t.Fatalf("got %d duplicate groups, expected 1", len(index.Groups()))
	}

	if x := len(index.Groups()[0].Paths); x != 3 {
		t.Fatalf("got %d files in duplicate group, expected 3", x)
	}
}

func Test_nodups(t*testing.T) {
	basedir := path.Join(testDataDir, "nodups")
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		t.Fatal(err)
	}

	index := newIndex()
	for _, file := range files {
		p := path.Join(basedir, file.Name())
		f := NewFileHandler(p, file)
		index.Add(f)
	}

	if len(index.Groups()) != 0 {
		t.Fatalf("got %d duplicate groups, expected none", len(index.Groups()))
	}
}

func Test_samesize(t*testing.T) {
	basedir := path.Join(testDataDir, "samesize")
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		t.Fatal(err)
	}

	index := newIndex()
	for _, file := range files {
		p := path.Join(basedir, file.Name())
		f := NewFileHandler(p, file)
		index.Add(f)
	}

	if len(index.Groups()) != 0 {
		t.Fatalf("got %d duplicate groups, expected none", len(index.Groups()))
	}
}

func TestCompareReaders_same_file(t*testing.T) {
	content := "dummy content"
	reader1 := strings.NewReader(content)
	reader2 := strings.NewReader(content)

	expected := 0

	cmp, err := CompareReaders(reader1, reader2)
	if err != nil {
		t.Errorf("Compare(f, f) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(f, f) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_size_ascending(t*testing.T) {
	smaller := strings.NewReader("dummy content")
	bigger := strings.NewReader("longer dummy content")

	expected := -1

	cmp, err := CompareReaders(smaller, bigger)
	if err != nil {
		t.Errorf("Compare(smaller, bigger) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(smaller, bigger) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_size_descending(t*testing.T) {
	smaller := strings.NewReader("dummy content")
	bigger := strings.NewReader("longer dummy content")

	expected := 1

	cmp, err := CompareReaders(bigger, smaller)
	if err != nil {
		t.Errorf("Compare(bigger, smaller) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(bigger, smaller) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_same_size_content_ascending(t*testing.T) {
	lower := strings.NewReader("dummy content a")
	higher := strings.NewReader("dummy content b")

	expected := -1

	cmp, err := CompareReaders(lower, higher)
	if err != nil {
		t.Errorf("Compare(lower, higher) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(lower, higher) == %v, want %v", cmp, expected)
	}
}

func TestCompareReaders_same_size_content_descending(t*testing.T) {
	lower := strings.NewReader("dummy content a")
	higher := strings.NewReader("dummy content b")

	expected := 1

	cmp, err := CompareReaders(higher, lower)
	if err != nil {
		t.Errorf("Compare(higher, lower) raised error: %v", err)
	}
	if cmp != expected {
		t.Errorf("Compare(higher, lower) == %v, want %v", cmp, expected)
	}
}

func TestDupTracker_merge_pools(t*testing.T) {
	tracker := newDupTracker()
	tracker.add("path1-1", "path1-2")
	tracker.add("path1-3", "path1-2")

	duplicates := tracker.getDuplicates()
	if len(duplicates) != 1 {
		t.Errorf("Found %d duplicate groups, expected %d", len(duplicates), 1)
	}
	if duplicates[0].count() != 3 {
		t.Errorf("Found %d duplicate files, expected %d", duplicates[0].count(), 3)
	}
}
