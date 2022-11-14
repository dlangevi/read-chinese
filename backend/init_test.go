package backend

import (
	"context"
	"log"
	"os"
	"path"
	"testing"
)

// For TestMain we load up a running full system using test data
// Indivitual componenets will test using these shared resources for now
// The goal is to eventually detangle functionality into more discrete
// code packages
func TestMain(m *testing.M) {
	tempDb := path.Join(os.TempDir(), "testdb.db")
	os.Remove(tempDb)
	ctx := context.Background()
	_, err := StartBackend(&ctx, tempDb, "./testdata/example_metadata.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// This will be the first book added, so tests can query its data with
	// BookdId = 1
	err = AddBook("张天翼", "秃秃大王", "fake.jpg", "testdata/example_book.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	code := m.Run()

	// os.Remove(tempDb)
	os.Exit(code)
}
