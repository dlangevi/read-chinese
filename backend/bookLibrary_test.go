package backend

import (
	"testing"
)

func TestBookExists(t *testing.T) {
	exists, err := bookExists(testRuntime.DB, "foo", "bar")
	if err != nil {
		t.Errorf("Failed to call bookExists: %v", err)
	}
	if exists {
		t.Error("The not there book foo:bar somehow exists")
	}
	addBook(testRuntime.DB, "foo", "bar", "cover", "location")
	exists, err = bookExists(testRuntime.DB, "foo", "bar")
	if err != nil {
		t.Errorf("Failed to second call bookExists: %v", err)
	}
	if !exists {
		t.Errorf("The there book foo:bar somehow not there %v", exists)
	}
}

func TestAddBook(t *testing.T) {
	_, err := addBook(testRuntime.DB, "foot", "bar", "cover", "location")
	if err != nil {
		t.Errorf("Failed to call addBook: %v", err)
	}
	_, err = addBook(testRuntime.DB, "foot", "bar", "cover", "location")
	if err == nil {
		t.Errorf("Managed to double insert foot book")
	}
}

func TestModifyBook(t *testing.T) {
	books := testRuntime.BookLibrary
	bookId, err := addBook(testRuntime.DB, "fish", "is", "pretty", "tasty")
	if err != nil {
		t.Errorf("Failed to call addBook: %v", err)
	}
	book, _ := getBook(testRuntime.DB, bookId)
	if book.Favorite != false {
		t.Errorf("Book started as favorited")
	}
	books.SetFavorite(book.BookId, true)
	book, _ = getBook(testRuntime.DB, bookId)
	if book.Favorite != true {
		t.Errorf("Book favorite not changed")
	}
	books.SetFavorite(book.BookId, false)
	book, _ = getBook(testRuntime.DB, bookId)
	if book.Favorite != false {
		t.Errorf("Book favorite not changed")
	}
	books.SetRead(book.BookId, true)
	book, _ = getBook(testRuntime.DB, bookId)
	if book.HasRead != true {
		t.Errorf("Book read not changed")
	}
	books.SetRead(book.BookId, false)
	book, _ = getBook(testRuntime.DB, bookId)
	if book.HasRead != false {
		t.Errorf("Book read not changed")
	}
	books.DeleteBook(bookId)
	_, err = getBook(testRuntime.DB, bookId)
	if err == nil {
		t.Errorf("Book was not deleted")
	}
	exists, err := bookExists(testRuntime.DB, "fish", "is")
	if err != nil {
		t.Errorf("Problem with exists")
	}
	if exists {
		t.Errorf("Books claims to exist")
	}
}

func TestSelectBooks(t *testing.T) {
	bar1, _ := addBook(testRuntime.DB, "foo", "bar1", "cover", "location")
	bar2, _ := addBook(testRuntime.DB, "foo", "bar2", "cover", "location")
	bar3, _ := addBook(testRuntime.DB, "foo", "bar3", "cover", "location")

	books, err := getBooks(testRuntime.DB)
	if err != nil {
		t.Errorf("Failed to get all books %v", err)
	}
	if len(books) < 3 {
		t.Errorf("Not enough all books %v", len(books))
	}
	books, err = getBooks(testRuntime.DB, bar1, bar2, bar3)
	if err != nil {
		t.Errorf("Failed to get 3 books %v", err)
	}
	if len(books) != 3 {
		t.Errorf("Not enough specific books %v", len(books))
	}
	book, err := getBook(testRuntime.DB, bar1)
	if err != nil {
		t.Errorf("Failed to get bar1 %v", err)
	}
	if book.Title != "bar1" {
		t.Errorf("Bad values in bar1 %v", book)
	}
}
