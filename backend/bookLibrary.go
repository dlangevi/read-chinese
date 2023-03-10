package backend

import (
	"database/sql"
	"encoding/csv"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	BookLibrary interface {
		// Add book to collection of books
		AddBook(author string, title string, cover string, filepath string) (int64, error)
		// Delete book from book collection
		DeleteBook(bookId int64) error
		// Get a book by its unique identifier
		GetBook(bookId int64) (Book, error)
		// Get a list of books by their unique identifiers
		GetSomeBooks(bookIds ...int64) ([]Book, error)
		// Get all books
		GetBooks() ([]Book, error)
		GetDetailedBooks(path string) ([]Book, error)
		// Check if book already exists in collection
		BookExists(author string, title string) (bool, error)
		HealthCheck() error

		// Resegement books
		RecalculateBooks() error

		// Get the words in the library that occure the most often
		LearningTarget() []string
		// Get the words in the library that occure the most often in favorite books
		LearningTargetFavorites() []string
		// Get the words in a specific book that occure the most often
		LearningTargetBook(bookId int) []string
		TopUnknownWords(bookId int, numWords int) []string
		computeKnownCharacters(book *Book) error
		TotalRead() (int, error)

		// Mark book as a favorite
		SetFavorite(bookId int64, isFavorite bool) error
		// Mark book as having been read
		SetRead(bookId int64, isRead bool) error
	}

	// bookLibrary implements bookLibrary
	bookLibrary struct {
		backend *Backend
		db      *sqlx.DB
	}

	Book struct {
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

	BookStats struct {
		ProbablyKnownWords int   `json:"probablyKnownWords"`
		KnownCharacters    int   `json:"knownCharacters"`
		UniqueCharacters   int   `json:"uniqueCharacters"`
		UniqueWords        int   `json:"uniqueWords"`
		TotalCharacters    int   `json:"totalCharacters"`
		TotalWords         int   `json:"totalWords"`
		TotalKnownWords    int   `json:"totalKnownWords"`
		Targets            []int `json:"targets"`
		TargetOccurances   []int `json:"targetOccurances"`
		NeedToKnow         []int `json:"needToKnow"`
	}

	WordTableRow struct {
		BookId int    `db:"book"`
		Word   string `db:"word"`
		Count  int    `db:"count"`
	}
	WordTable []WordTableRow

	WordOccuranceRow struct {
		Word      string `json:"word" db:"word"`
		Occurance int    `json:"occurance" db:"occurance"`
	}
)

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

func NewBookLibrary(
	b *Backend,
	db *sqlx.DB,
) *bookLibrary {
	return &bookLibrary{
		backend: b,
		db:      db,
	}
}

func copyCover(author string, title string, coverPath string) (string, error) {
	if coverPath == "" {
		return "", errors.New("Empty coverpath")
	}
	bytes, err := os.ReadFile(coverPath)
	if err != nil {
		return "", err
	}

	ext := path.Ext(coverPath)
	subpath := path.Join(author, fmt.Sprintf("%s%s", title, ext))
	newCoverLocation := ConfigDir("covers", subpath)
	err = os.WriteFile(newCoverLocation, bytes, 0666)
	return subpath, err
}

func (b *bookLibrary) emitBooks() error {
	books, err := b.GetSomeBooks()
	if err != nil {
		return err
	}
	runtime.EventsEmit(b.backend.ctx, "BooksUpdated", books)
	return nil
}

func (b *bookLibrary) RecalculateBooks() error {
	books, err := b.GetBooks()
	if err != nil {
		return err
	}
	b.backend.setupProgress("Resegmenting books", len(books))
	for _, book := range books {
		log.Println("Processing:", book.Title, "...")
		filepath := book.Filepath
		sentences, wordTable, err := b.backend.Segmentation.SegmentFullText(filepath)
		if err != nil {
			return err
		}

		cacheLocation := getSegmentationPath(book.Title, book.Author)
		err = saveCacheFile(b.db, int(book.BookId), sentences, cacheLocation)
		if err != nil {
			return err
		}

		err = deleteWordTable(b.db, int(book.BookId))
		if err != nil {
			return err
		}
		_, err = saveWordTable(b.db, int(book.BookId), wordTable)
		if err != nil {
			return err
		}
		b.backend.progress()
	}

	b.backend.KnownWords.syncFrequency()
	b.emitBooks()
	return nil
}

func (b *bookLibrary) AddBook(
	author string,
	title string,
	cover string,
	filepath string) (int64, error) {

	// If there is a problem copying cover maybe that is not a big deal?
	cover, err := copyCover(author, title, cover)
	if err != nil {
		log.Println("Error copying cover: ", err)
	}
	bookId, err := addBookToDb(b.db, author, title, cover, filepath)
	if err != nil {
		return 0, err
	}
	sentences, wordTable, err := b.backend.Segmentation.SegmentFullText(filepath)
	if err != nil {
		return 0, err
	}

	// Compute WordTable and Save it
	cacheLocation := getSegmentationPath(title, author)
	err = saveCacheFile(b.db, int(bookId), sentences, cacheLocation)
	if err != nil {
		return 0, err
	}
	_, err = saveWordTable(b.db, int(bookId), wordTable)
	// This can be nil
	b.emitBooks()
	return bookId, err
}

func addBookToDb(db *sqlx.DB, author string, title string, cover string, filepath string) (int64, error) {
	res, err := db.Exec(`
  INSERT INTO books (author, title, cover, filepath)
  VALUES ($1, $2, $3, $4)
  `, author, title, cover, filepath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (b *bookLibrary) DeleteBook(bookId int64) error {
	tx, err := b.db.Begin()
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
	err = tx.Commit()
	b.emitBooks()
	return err
}

// Returns details for single book, with extra stats
func (b *bookLibrary) GetBook(bookId int64) (Book, error) {
	books, err := b.GetSomeBooks(bookId)
	if err != nil {
		return Book{}, err
	}
	if len(books) != 1 {
		return Book{}, errors.New("Book did not exist")
	}
	book := books[0]
	b.computeKnownCharacters(&book)
	computeWordTargets(b.db, &book)

	return book, nil
}

func getBooks(db *sqlx.DB, bookIds ...int64) ([]Book, error) {
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
	err = db.Select(&books, query, args...)
	if err != nil {
		return books, err
	}
	return books, nil
}

func (b *bookLibrary) GetSomeBooks(bookIds ...int64) ([]Book, error) {
	books, _ := getBooks(b.db, bookIds...)
	for index := range books {
		book := &books[index]
		book.Stats = NewBookStats()
		// For now, only these are needed by BookCard
		book.Stats.TotalKnownWords, _ = getKnownWords(b.db, book.BookId)
		book.Stats.TotalWords, _ = getTotalWords(b.db, book.BookId)
	}
	return books, nil
}

// Returns summary of all books
func (b *bookLibrary) GetBooks() ([]Book, error) {
	return b.GetSomeBooks()
}

func (b *bookLibrary) GetDetailedBooks(path string) ([]Book, error) {
	books, _ := getBooks(b.db)

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.Write([]string{
		"Author",
		"Title",
		"UniqueCharacters",
		"UniqueWords",
		"TotalCharacters",
		"TotalWords",
		"KnownCharacters",
		"ProbablyKnownWords",
		"TotalKnownWords",
	})

	for index := range books {
		book := &books[index]
		book.Stats = NewBookStats()
		// For now, only these are needed by BookCard
		book.Stats.TotalKnownWords, _ = getKnownWords(b.db, book.BookId)
		book.Stats.TotalWords, _ = getTotalWords(b.db, book.BookId)
		b.computeKnownCharacters(book)
		computeWordTargets(b.db, book)
		writer.Write([]string{
			book.Author,
			book.Title,
			strconv.Itoa(book.Stats.UniqueCharacters),
			strconv.Itoa(book.Stats.UniqueWords),
			strconv.Itoa(book.Stats.TotalCharacters),
			strconv.Itoa(book.Stats.TotalWords),
			strconv.Itoa(book.Stats.KnownCharacters),
			strconv.Itoa(book.Stats.ProbablyKnownWords),
			strconv.Itoa(book.Stats.TotalKnownWords),
		})
	}

	writer.Flush()

	return books, nil
}

func (b *bookLibrary) BookExists(author string, title string) (bool, error) {
	var exists bool
	err := b.db.QueryRow(`SELECT exists (
    SELECT bookId 
    FROM books 
    WHERE author = $1
    AND title = $2
  )`, author, title).Scan(&exists)
	return exists, err
}

func (b *bookLibrary) HealthCheck() error {
	var exists bool
	err := b.db.QueryRow(`SELECT exists (
    SELECT bookId 
    FROM books 
  )`).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("User has 0 books")
	}
	return nil
}

func (b *bookLibrary) GetFavoriteFrequencies() (map[string]int, error) {
	frequencies := []WordOccuranceRow{}
	frequencyMap := map[string]int{}
	err := b.db.Select(&frequencies, `
    SELECT word, sum(count) as occurance FROM frequency 
    WHERE EXISTS (
      SELECT bookId
      FROM books
      WHERE books.bookId==frequency.book
      AND books.favorite==true
    )
    GROUP BY word
    ORDER BY occurance DESC
    `)

	if err != nil {
		return frequencyMap, err
	}

	for _, frequency := range frequencies {
		frequencyMap[frequency.Word] = frequency.Occurance
	}
	return frequencyMap, nil
}

func (b *bookLibrary) GetBookFrequencies(bookId int) (map[string]int, error) {
	frequencies := []WordOccuranceRow{}
	frequencyMap := map[string]int{}
	err := b.db.Select(&frequencies, `
    SELECT word, count as occurance FROM frequency 
    WHERE book = $1
    GROUP BY word
    ORDER BY occurance DESC
    `, bookId)

	if err != nil {
		return frequencyMap, err
	}

	for _, frequency := range frequencies {
		frequencyMap[frequency.Word] = frequency.Occurance
	}
	return frequencyMap, nil
}

func (b *bookLibrary) LearningTarget() []string {
	words := []string{}
	err := b.db.Select(&words, `
    SELECT word FROM frequency 
    WHERE NOT EXISTS (
        SELECT word
        FROM words
        WHERE words.word==frequency.word
    ) 
    GROUP BY word
    ORDER BY sum(count) DESC
    LIMIT 200
    `)
	if err != nil {
		log.Println(err)
		return words
	}

	return words
}

func (b *bookLibrary) LearningTargetFavorites() []string {
	words := []string{}
	err := b.db.Select(&words, `
    SELECT word FROM frequency 
    WHERE NOT EXISTS (
        SELECT word
        FROM words
        WHERE words.word==frequency.word
    ) AND EXISTS (
      SELECT bookId
      FROM books
      WHERE books.bookId==frequency.book
      AND books.favorite==true
    )
    GROUP BY word
    ORDER BY sum(count) DESC
    LIMIT 200
    `)
	if err != nil {
		log.Println(err)
		return words
	}

	return words
}

func (b *bookLibrary) LearningTargetBook(bookId int) []string {
	words := []string{}
	err := b.db.Select(&words, `
    SELECT word FROM frequency 
    WHERE NOT EXISTS (
        SELECT word
        FROM words
        WHERE words.word==frequency.word
    ) 
    AND book = $1
    GROUP BY word
    ORDER BY count DESC
    LIMIT 200
    `, bookId)
	if err != nil {
		log.Println(err)
		return words
	}

	return words
}

func (b *bookLibrary) TopUnknownWords(bookId int, numWords int) []string {
	words := []string{}
	err := b.db.Select(&words, `
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

func getKnownWords(db *sqlx.DB, bookId int64) (int, error) {
	var known sql.NullInt64
	err := db.QueryRow(`
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
	if known.Valid {
		return int(known.Int64), err
	} else {
		return 0, err
	}
}

func getTotalWords(db *sqlx.DB, bookId int64) (int, error) {
	var total sql.NullInt64
	err := db.QueryRow(`
    SELECT SUM(count) as known 
    FROM frequency
    WHERE book = $1
    `, bookId).Scan(&total)
	if err != nil {
		log.Println("Error with totalWords", err)
	}
	if total.Valid {
		return int(total.Int64), err
	} else {
		return 0, err
	}
}

func (b *bookLibrary) computeKnownCharacters(book *Book) error {
	words := []struct {
		Word  string
		Count int
	}{}

	err := b.db.Select(&words, `
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
	var uniqueWords = map[string]bool{}
	var uniqueCharacters = map[rune]bool{}
	for _, row := range words {
		totalCharacters += len([]rune(row.Word)) * row.Count
		allKnown := true
		uniqueWords[row.Word] = true
		for _, char := range row.Word {
			uniqueCharacters[char] = true
			if b.backend.KnownWords.isKnownChar(char) {
				knownCharacters += row.Count
			} else {
				allKnown = false
			}
		}
		if b.backend.KnownWords.isKnown(row.Word) || allKnown {
			probablyKnownWords += row.Count
		}
	}
	book.Stats.UniqueCharacters = len(uniqueCharacters)
	book.Stats.UniqueWords = len(uniqueWords)
	book.Stats.ProbablyKnownWords = probablyKnownWords
	book.Stats.KnownCharacters = knownCharacters
	book.Stats.TotalCharacters = totalCharacters
	return nil
}

func computeWordTargets(db *sqlx.DB, book *Book) error {
	words := []struct {
		Word  string
		Count int
	}{}

	err := db.Select(&words, `
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
		80, 84, 86, 90, 92, 94, 95, 96, 97, 98, 99, 100,
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

func (b *bookLibrary) TotalRead() (int, error) {
	var total sql.NullInt64
	err := b.db.QueryRow(`
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

func (b *bookLibrary) TotalReadChars() (int, error) {
	var total sql.NullInt64
	err := b.db.QueryRow(`
    SELECT SUM(LENGTH(word) * count) total 
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

func (b *bookLibrary) SetFavorite(bookId int64, isFavorite bool) error {
	_, err := b.db.Exec(`
  UPDATE books 
  SET favorite = ?1 
  WHERE bookId = ?2`, isFavorite, bookId)
	b.emitBooks()
	return err
}

func (b *bookLibrary) SetRead(bookId int64, isRead bool) error {
	_, err := b.db.Exec(`
  UPDATE books 
  SET has_read = ?1 
  WHERE bookId = ?2`, isRead, bookId)
	b.emitBooks()
	return err
}

func saveCacheFile(
	db *sqlx.DB,
	bookId int,
	sentences []TokenizedSentence,
	filepath string) error {

	file, err := os.Create(filepath)
	if err != nil {
		return nil
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(sentences); err != nil {
		return err
	}

	_, err = db.Exec(`
  UPDATE books 
  SET segmented_file = ?1 
  WHERE bookId = ?2`, filepath, bookId)
	return err

}

func getSegmentationPath(title string, author string) string {
	fileName := fmt.Sprintf("%v-%v.gob", title, author)
	cacheLocation := ConfigDir("segmentationCache", fileName)
	return cacheLocation
}

func GetSegmentedText(book Book) ([]TokenizedSentence, error) {
	if !book.SegmentedFile.Valid {
		return nil, errors.New("Book has not been segmented yet")
	}
	cacheLocation := getSegmentationPath(book.Title, book.Author)

	file, err := os.Open(cacheLocation)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	s := []TokenizedSentence{}
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}

	return s, nil
}

func deleteWordTable(db *sqlx.DB, bookId int) error {
	_, err := db.Exec(
		`DELETE FROM frequency WHERE book = $1`, bookId)
	return err
}

func saveWordTable(db *sqlx.DB, bookId int, frequencyTable FrequencyTable) (sql.Result, error) {
	wordTable := WordTable{}
	for word, count := range frequencyTable {
		wordTable = append(wordTable, WordTableRow{
			BookId: bookId,
			Word:   word,
			Count:  count,
		})
	}
	tx, err := db.Beginx()
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
