package backend

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type Book struct {
	Author        string         `db:"author" json:"author"`
	Title         string         `db:"title" json:"title"`
	Cover         string         `db:"cover" json:"cover"`
	Filepath      string         `db:"filepath" json:"filepath"`
	BookId        int64          `db:"bookId" json:"bookId"`
	Favorite      bool           `db:"favorite" json:"favorite"`
	SegmentedFile sql.NullString `db:"segmented_file"`
	HasRead       bool           `db:"has_read" json:"hasRead"`
	Stats         BookStats      `json:"stats"`
}

type BookStats struct {
	ProbablyKnownWords int   `json:"probablyKnownWords"`
	KnownCharacters    int   `json:"knownCharacters"`
	TotalCharacters    int   `json:"totalCharacters"`
	TotalWords         int   `json:"totalWords"`
	TotalKnownWords    int   `json:"totalKnownWords"`
	Targets            []int `json:"targets"`
	TargetOccurances   []int `json:"targetOccurances"`
	NeedToKnow         []int `json:"needToKnow"`
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
	return exists, err
}

// Returns summary of all book
func (BookLibrary) GetBooks() ([]Book, error) {
	books, _ := getBooks()
	for index := range books {
		book := &books[index]
		book.Stats = NewBookStats()
		// For now, only these are needed by BookCard
		book.Stats.TotalKnownWords, _ = getKnownWords(book.BookId)
		book.Stats.TotalWords, _ = getTotalWords(book.BookId)
	}

	return books, nil
}

func getKnownWords(bookId int64) (int, error) {
	var known int
	err := Conn.QueryRow(`
    SELECT SUM(count) as known 
    FROM frequency
    WHERE book = $1
    AND EXISTS (
        SELECT word
        FROM words
        WHERE words.word==frequency.word
    )`, bookId).Scan(&known)
	if err != nil {
		log.Println("Error with getKnownWords", err)
	}
	return known, err
}

func getTotalWords(bookId int64) (int, error) {
	var total int
	err := Conn.QueryRow(`
    SELECT SUM(count) as known 
    FROM frequency
    WHERE book = $1
    `, bookId).Scan(&total)
	if err != nil {
		log.Println("Error with totalWords", err)
	}
	return total, err
}

// Returns details for single book, with extra stats
func (BookLibrary) GetBook(bookId int64) (Book, error) {
	book, err := getBook(bookId)
	if err != nil {
		return book, err
	}
	book.Stats = NewBookStats()
	book.Stats.TotalKnownWords, _ = getKnownWords(bookId)
	book.Stats.TotalWords, _ = getTotalWords(bookId)
	_ = computeKnownCharacters(&book)
	_ = computeWordTargets(&book)

	return book, nil
}

func computeKnownCharacters(book *Book) error {
	words := []struct {
		Word  string
		Count int
	}{}

	err := Conn.Select(&words, `
    SELECT word, count 
    FROM frequency 
    WHERE book = $1
    `, book.BookId)
	if err != nil {
		log.Println(err)
		return err
	}
	var probablyKnownWords = 0
	var knownCharacters = 0
	var totalCharacters = 0
	for _, row := range words {
		totalCharacters += len([]rune(row.Word)) * row.Count
		allKnown := true
		for _, char := range row.Word {
			if known.isKnownChar(char) {
				knownCharacters += row.Count
			} else {
				allKnown = false
			}
		}
		if known.isKnown(row.Word) || allKnown {
			probablyKnownWords += row.Count
		}
	}
	book.Stats.ProbablyKnownWords = probablyKnownWords
	book.Stats.KnownCharacters = knownCharacters
	book.Stats.TotalCharacters = totalCharacters
	return nil
}

func computeWordTargets(book *Book) error {
	words := []struct {
		Word  string
		Count int
	}{}

	err := Conn.Select(&words, `
    SELECT word, count 
    FROM frequency 
    WHERE book = $1
    AND NOT EXISTS (
      SELECT word
      FROM words
      WHERE words.word==frequency.word
    )
    ORDER BY count DESC
    `, book.BookId)
	if err != nil {
		log.Println(err)
		return err
	}

	targets := [...]int{
		80, 84, 86, 90, 92, 94, 96, 98, 100,
	}
	targetOccurances := []int{}
	needToKnow := []int{}
	for _, target := range targets {
		targetOccurance := int(float64(target*book.Stats.TotalWords) / 100)
		targetOccurances = append(targetOccurances, int(targetOccurance))

		soFar := book.Stats.TotalKnownWords
		needToLearn := 0
		for _, row := range words {
			if soFar > targetOccurance {
				break
			}
			soFar += row.Count
			needToLearn += 1
		}
		needToKnow = append(needToKnow, needToLearn)
	}
	book.Stats.Targets = targets[:]
	book.Stats.TargetOccurances = targetOccurances
	book.Stats.NeedToKnow = needToKnow

	return nil
}

func getBook(bookId int64) (Book, error) {
	books, err := getBooks(bookId)
	if err != nil {
		return Book{}, err
	}
	if len(books) != 1 {
		return Book{}, errors.New("Book did not exist")
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

func deleteBook(bookId int64) error {
	tx, err := Conn.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(
		`DELETE FROM books WHERE bookId = $1`, bookId); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(
		`DELETE FROM frequency WHERE book = $1`, bookId); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (BookLibrary) DeleteBook(bookId int64) error {
	return deleteBook(bookId)
}

func (BookLibrary) SetFavorite(bookId int64, isFavorite bool) error {
	_, err := Conn.Exec(`
  UPDATE books 
  SET favorite = ?1 
  WHERE bookId = ?2`, isFavorite, bookId)
	return err
}

func (BookLibrary) SetRead(bookId int64, isRead bool) error {
	_, err := Conn.Exec(`
  UPDATE books 
  SET has_read = ?1 
  WHERE bookId = ?2`, isRead, bookId)
	return err
}

func (BookLibrary) TotalRead() (int, error) {
	var total sql.NullInt64
	err := Conn.QueryRow(`
    SELECT SUM(count) as total 
    FROM frequency 
    WHERE EXISTS (
      SELECT bookId
      FROM books
      WHERE has_read = true
      AND books.bookId == frequency.book
    )`).Scan(&total)
	if total.Valid {
		return int(total.Int64), err
	} else {
		return 0, err
	}

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
	cacheLocation := ConfigDir("segmentationCache", fileName)
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
func saveWordTable(bookId int, frequencyTable FrequencyTable) (sql.Result, error) {
	wordTable := WordTable{}
	for word, count := range frequencyTable {
		wordTable = append(wordTable, WordTableRow{
			BookId: bookId,
			Word:   word,
			Count:  count,
		})
	}
	tx, err := Conn.Beginx()
	if err != nil {
		return nil, err
	}
	var res sql.Result
	for rows := 0; rows < len(wordTable); rows += 5000 {
		upperLimit := rows + 5000
		if upperLimit > len(wordTable) {
			upperLimit = len(wordTable)
		}
		currentBatch := wordTable[rows:upperLimit]
		res, err = tx.NamedExec(`
    INSERT INTO frequency (book, word, count)
    VALUES (:book, :word, :count)`, currentBatch)
		if err != nil {
			return nil, err
		}

	}
	err = tx.Commit()
	return res, err
}

// TODO might want to run the segementation preloadWords on
// bookLibrary initialization
