package backend

import (
	"testing"
)

func TestSegment(t *testing.T) {
	s := testRuntime.Segmentation
	sentences, table, err := s.SegmentFullText("testdata/example_book.txt")
	if err != nil {
		t.Error("Failed to load text", err)
	}
	if len(sentences) != 9 {
		t.Error("Wrong number of sentences created", len(sentences))
	}
	if len(table) != 77 {
		t.Error("Wrong number of words found")
	}
}
