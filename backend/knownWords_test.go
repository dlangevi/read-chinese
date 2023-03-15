package backend

import (
	"os"
	"path"
	"testing"
)

func checkDb(myBackend *Backend) []string {
	rows := []string{}
	myBackend.DB.Select(&rows, "SELECT word FROM words")
	return rows
}

func TestWords(t *testing.T) {
	tempDb := path.Join(os.TempDir(), "generate.db")
	os.Remove(tempDb)
	defer os.Remove(tempDb)
	myBackend := createBackend(tempDb)
	myBackend.KnownWords.AddWord("真", 100)
	// Not well known enough

	rows := checkDb(myBackend)
	if len(rows) != 1 {
		t.Errorf("rows %v", rows)
	}
	myBackend.KnownWords.AddWord("你好", 100)
	myBackend.KnownWords.AddWord("真的", 100)
	myBackend.KnownWords.AddWord("哇塞", 100)
	rows = checkDb(myBackend)
	if len(rows) != 4 {
		t.Errorf("rows %v", rows)
	}
	myBackend.KnownWords.DeleteWord("真")
	rows = checkDb(myBackend)
	if len(rows) != 3 {
		t.Errorf("rows %v", rows)
	}
}
