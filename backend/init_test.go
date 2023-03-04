package backend

import (
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

	log.SetFlags(log.Ltime | log.Lshortfile)
	testRuntime = createBackend(tempDb)
	// This will be the first book added, so tests can query its data with
	// BookdId = 1
	_, err := testRuntime.BookLibrary.AddBook("张天翼", "秃秃大王", "fake.jpg", "testdata/example_book.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	code := m.Run()

	// os.Remove(tempDb)
	os.Exit(code)
}

func createBackend(dbPath string) *Backend {
	runtime := NewBackend(dbPath, "./testdata/example_metadata.json")
	// TODO lol this saves an actual dictionary in my user config
	runtime.Dictionaries.AddMigakuDictionary("example",
		"./testdata/example_dict.json", "english")
	return runtime
}
