package backend

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"read-chinese/backend/segmentation"

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

func NewBookStats() BookStats {
	return BookStats{
		ProbablyKnownWords: 0,
		KnownCharacters:    0,
		TotalCharacters:    0,
		TotalWords:         0,
		TotalKnownWords:    0,
		Targets:            []int{},
		TargetOccurances:   []int{},
		NeedToKnow:         []int{},
	}

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

func (BookLibrary) LearningTargetBook(bookId int) []WordOccuranceRow {
	words := []WordOccuranceRow{}
	err := Conn.Select(&words, `
    SELECT word, sum(count) as occurance FROM frequency 
    WHERE NOT EXISTS (
        SELECT word
        FROM words
        WHERE words.word==frequency.word
    ) 
    AND book = $1
    GROUP BY word
    ORDER BY occurance DESC
    LIMIT 200
    `, bookId)
	if err != nil {
		log.Println(err)
		return words
	}

	return words
}

func (BookLibrary) TopUnknownWords(bookId int, numWords int) []string {
	words := []string{}
	err := Conn.Select(&words, `
    SELECT word
    FROM frequency
    WHERE NOT EXISTS (
        SELECT word
        FROM words
        WHERE words.word==frequency.word
    ) 
    AND book = $1
    ORDER BY count DESC
    LIMIT $2
  `, bookId, numWords)
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

func (BookLibrary) GetBooks(bookIds ...int64) ([]Book, error) {
	books, _ := getBooks(bookIds...)
	for index := range books {
		book := &books[index]
		book.Stats = NewBookStats()
		// Compute totalKnownWords
		// Compute totalWords
		// if detailed {
		//    Compute probablyKnownWords
		//    Compute knownCharacters
		//    Compute totalCharacters
		//    Compute WordTargets
		//    Compute targets
		//    Compute targetOccurances
		//    Compute needToKnow
		// }
	}

	return books, nil
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

func AddBook(author string, title string, cover string, filepath string) error {
	bookId, err := addBook(author, title, cover, filepath)
	if err != nil {
		return err
	}
	sentences, wordTable, err := runtime.Segmentation.SegmentFullText(filepath)
	if err != nil {
		return err
	}

	// Compute WordTable and Save it
	cacheLocation := getSegmentationPath(title, author)
	err = saveCacheFile(int(bookId), sentences, cacheLocation)
	if err != nil {
		return err
	}
	_, err = saveWordTable(int(bookId), wordTable)
	// This can be nil
	return err
}

func saveCacheFile(bookId int, sentences []string, filepath string) error {
	bytes, err := json.Marshal(sentences)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, bytes, 0666)
	if err != nil {
		return err
	}
	_, err = Conn.Exec(`
  UPDATE books 
  SET segmented_file = ?1 
  WHERE bookId = ?2`, filepath, bookId)
	return err

}

func getSegmentationPath(title string, author string) string {
	fileName := fmt.Sprintf("%v-%v.json", title, author)
	cacheLocation := path.Join(
		"/home/dlangevi/.config/read-chinese/",
		"segmentationCache",
		fileName)
	return cacheLocation
}

func GetSegmentedText(book Book) ([]string, error) {
	if !book.SegmentedFile.Valid {
		return nil, errors.New("Book has not been segmented yet")
	}
	cacheLocation := getSegmentationPath(book.Title, book.Author)

	segBytes, err := os.ReadFile(cacheLocation)
	if err != nil {
		return nil, err
	}
	sentences := []string{}
	err = json.Unmarshal(segBytes, &sentences)
	return sentences, err

}

// dbSaveWordTable, // TODO once segmentation is done we can test this
func saveWordTable(bookId int, frequencyTable segmentation.FrequencyTable) (sql.Result, error) {

	wordTable := WordTable{}
	for word, count := range frequencyTable {
		wordTable = append(wordTable, WordTableRow{
			BookId: bookId,
			Word:   word,
			Count:  count,
		})
	}
	return Conn.NamedExec(`INSERT INTO frequency (book, word, count)
  VALUES (:book, :word, :count)`, wordTable)
}

// export const bookLibraryIpc = {
//   loadBooks,
//   loadBook,
// initBookStats
// computeBookData

//   deleteBook, straigt sql
//   setFavorite, straight sql
//   setRead, straigt sql
//   totalRead, straight sql

// };

//  segmentation.preloadWords ?
