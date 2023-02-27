package backend

import (
	"testing"
)

func TestSegment(t *testing.T) {
	s := testRuntime.Segmentation
	segmentedSentences, table, err := s.SegmentFullText("testdata/example_book.txt")
	if err != nil {
		t.Error("Failed to load text", err)
	}
	if len(segmentedSentences) != 9 {
		t.Error("Wrong number of sentences created", len(segmentedSentences))
	}
	if len(table) != 77 {
		t.Error("Wrong number of words found")
	}
}
