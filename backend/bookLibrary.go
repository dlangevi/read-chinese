package backend

import (
	"errors"
	"fmt"
	"log"

	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Book struct {
	Author   string `db:"author"`
	Title    string `db:"title"`
	Cover    string `db:"cover"`
	Filepath string `db:"filepath"`
	BookId   int    `db:"bookId"`
	Favorite bool   `db:"favorite"`
	// TODO mabey get rid of this by making a default exist?
	SegmentedFile sql.NullString `db:"segmented_file"`
	HasRead       bool           `db:"has_read"`
	Stats         BookStats
}

type BookStats struct {
	ProbablyKnownWords int
	KnownCharacters    int
	TotalCharacters    int
	TotalWords         int
	TotalKnownWords    int
	Targets            []int
	TargetOccurances   []int
	NeedToKnow         []int
}

type WordTableRow struct {
	BookId int    `db:"book"`
	Word   string `db:"word"`
	Count  int    `db:"count"`
}
type WordTable []WordTableRow

type WordOccuranceRow struct {
	Word      string `json:"word" db:"word"`
	Occurance int    `json:"occurance" db:"occurance"`
}

type BookLibrary struct {
}

func (BookLibrary) LearningTarget() []WordOccuranceRow {
	words := []WordOccuranceRow{}
	err := Conn.Select(&words, `
    SELECT word, sum(count) as occurance FROM frequency 
    WHERE NOT EXISTS (
        SELECT word
        FROM words
        WHERE words.word==frequency.word
    ) 
    GROUP BY word
    ORDER BY occurance DESC
    LIMIT 200
    `)
	if err != nil {
		log.Println(err)
		return words
	}

	return words
}

// dbBookExists,
func bookExists(author string, title string) (bool, error) {
	var exists bool
	err := Conn.QueryRow(`SELECT exists (
    SELECT bookId 
    FROM books 
    WHERE author = $1
    AND title = $2
  )`, author, title).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func getBook(bookId int64) (Book, error) {
	books, err := getBooks(bookId)
	if err != nil {
		return Book{}, err
	}
	return books[0], nil
}

// dbGetBooks,
func getBooks(bookIds ...int64) ([]Book, error) {
	books := []Book{}
	var args []interface{}
	var err error
	query := `
  SELECT 
    author,
    title,
    cover,
    filepath,
    bookId,
    favorite,
    segmented_file,
    has_read
  FROM books`
	if len(bookIds) != 0 {
		query = fmt.Sprintf(`%v WHERE bookId in (?)`, query)
		query, args, err = sqlx.In(query, bookIds)
		if err != nil {
			return books, err
		}
	}
	err = Conn.Select(&books, query, args...)
	if err != nil {
		return books, err
	}
	return books, nil
}

// dbAddBook, have it return the book id
func addBook(author string, title string, cover string, filepath string) (int64, error) {
	// TODO might want to prevent duplicates using unique constaint
	exists, _ := bookExists(author, title)
	if exists {
		return 0, errors.New("Cannot not add already existing book")
	}
	res, err := Conn.Exec(`
  INSERT INTO books (author, title, cover, filepath)
  VALUES ($1, $2, $3, $4)
  `, author, title, cover, filepath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// dbSaveWordTable, // TODO once segmentation is done we can test this
func saveWordTable(wordTable WordTable) (sql.Result, error) {
	return Conn.NamedExec(`INSERT INTO frequency (book, word, count)
  VALUES (:book, :word, :count)`, wordTable)
}

// export const bookLibraryIpc = {
//   loadBooks,
//   loadBook,
// initBookStats
// computeBookData

//   topUnknownWords, straigt sql
//   deleteBook, straigt sql
//   setFavorite, straight sql
//   setRead, straigt sql
//   totalRead, straight sql

//   hskWords, // Move this somewhere else
//   learningTarget, doneish (tests?)
// };
