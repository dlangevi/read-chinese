package backend

import (
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {

	tempDb := path.Join(os.TempDir(), "testdb.db")
	err := NewDB(tempDb)
	if err != nil {
		os.Exit(1)
	}
	err = RunMigrateScripts()
	if err != nil {
		os.Exit(1)
	}

	code := m.Run()

	os.Remove(tempDb)
	os.Exit(code)
}

func TestBookExists(t *testing.T) {
	exists, err := bookExists("foo", "bar")
	if err != nil {
		t.Errorf("Failed to call bookExists: %v", err)
	}
	if exists {
		t.Error("The not there book foo:bar somehow exists")
	}
	addBook("foo", "bar", "cover", "location")
	exists, err = bookExists("foo", "bar")
	if err != nil {
		t.Errorf("Failed to second call bookExists: %v", err)
	}
	if !exists {
		t.Errorf("The there book foo:bar somehow not there %v", exists)
	}
}

func TestAddBook(t *testing.T) {
	_, err := addBook("foot", "bar", "cover", "location")
	if err != nil {
		t.Errorf("Failed to call addBook: %v", err)
	}
	_, err = addBook("foot", "bar", "cover", "location")
	if err == nil {
		t.Errorf("Managed to double insert foot book")
	}
}

func TestSelectBooks(t *testing.T) {
	bar1, _ := addBook("foo", "bar1", "cover", "location")
	bar2, _ := addBook("foo", "bar2", "cover", "location")
	bar3, _ := addBook("foo", "bar3", "cover", "location")

	books, err := getBooks()
	if err != nil {
		t.Errorf("Failed to get all books %v", err)
	}
	if len(books) < 3 {
		t.Errorf("Not enough all books %v", len(books))
	}
	books, err = getBooks(bar1, bar2, bar3)
	if err != nil {
		t.Errorf("Failed to get 3 books %v", err)
	}
	if len(books) != 3 {
		t.Errorf("Not enough specific books %v", len(books))
	}
	book, err := getBook(bar1)
	if err != nil {
		t.Errorf("Failed to get bar1 %v", err)
	}
	if book.Title != "bar1" {
		t.Errorf("Bad values in bar1 %v", book)
	}
}
