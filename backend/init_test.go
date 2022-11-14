package backend

import (
	"context"
	"log"
	"os"
	"path"
	"testing"
)

// Tests can either use the global testRuntime (if they are not modifying any state)
// Or can load their own
var testRuntime *Backend

// For TestMain we load up a running full system using test data
// Indivitual componenets will test using these shared resources for now
// The goal is to eventually detangle functionality into more discrete
// code packages
func TestMain(m *testing.M) {
	tempDb := path.Join(os.TempDir(), "testdb.db")
	os.Remove(tempDb)

	testRuntime = createBackend(tempDb)
	// This will be the first book added, so tests can query its data with
	// BookdId = 1
	err := testRuntime.BookLibrary.AddBook("张天翼", "秃秃大王", "fake.jpg", "testdata/example_book.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	code := m.Run()

	// os.Remove(tempDb)
	os.Exit(code)
}

func createBackend(dbPath string) *Backend {
	ctx := context.Background()
	runtime, err := StartBackend(&ctx, dbPath, "./testdata/example_metadata.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return runtime
}
