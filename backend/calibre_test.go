package backend

import (
	"testing"
)

func TestCalibre(t *testing.T) {
	books, err := getCalibreBooks()
	if err != nil {
		t.Errorf("Failed to load books %v", err)
	}
	// TODO yeah dont just read from your own calibre lib
	if len(books) != 56 {
		t.Errorf("Where are the books? %v", len(books))
	}
}
