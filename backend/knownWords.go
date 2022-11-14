package backend

import (
	"embed"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
)

var known *KnownWords

type KnownWords struct {
	// For now just map word to interval
	words        map[string]int
	characters   map[rune]bool
	frequency    map[string]int
	db           *sqlx.DB
	userSettings *UserSettings
}

func NewKnownWords(db *sqlx.DB,
	userSettings *UserSettings,
) *KnownWords {
	known = &KnownWords{
		words:        map[string]int{},
		characters:   map[rune]bool{},
		frequency:    map[string]int{},
		db:           db,
		userSettings: userSettings,
	}
	known.syncWords()
	known.syncFrequency()
	return known
}

func (known *KnownWords) syncWords() {
	type WordRow struct {
		Word     string
		Interval int
	}
	words := []WordRow{}
	err := known.db.Select(&words, `
    SELECT word, interval 
    FROM words
  `)
	if err != nil {
		log.Fatal("Failed to load words", err)
	}
	for _, word := range words {
		known.words[word.Word] = word.Interval
		for _, char := range word.Word {
			known.characters[char] = true

		}
	}
}

func (known *KnownWords) syncFrequency() {
	type WordRow struct {
		Word  string
		Count int
	}
	words := []WordRow{}
	err := known.db.Select(&words, `
    SELECT word, sum(count) as count
    FROM frequency
    GROUP BY word
  `)
	if err != nil {
		log.Fatal("Failed to load frequency", err)
	}
	for _, word := range words {
		known.frequency[word.Word] = word.Count
	}
}

type WordStats struct {
	Words      int `json:"words"`
	Characters int `json:"characters"`
}

func (known *KnownWords) GetWordStats() WordStats {
	return WordStats{
		len(known.words),
		len(known.characters),
	}
}

// TODO have the updated_at automatically update
func (known *KnownWords) AddWord(word string, age int) error {
	tx, err := known.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
  INSERT OR IGNORE INTO words (word, interval) 
  VALUES ($1, $2)`, word, age)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`
  UPDATE words 
  SET interval=$2, 
      updated_at = CURRENT_TIMESTAMP 
  WHERE word="$1"
  `, word, age)
	if err != nil {
		tx.Rollback()
		return err
	}
	known.words[word] = age
	tx.Commit()
	return nil
}

type WordEntry struct {
	Word     string
	Interval int64
}

// TODO make this faster if possible
func (known *KnownWords) AddWords(words []WordEntry) error {
	tx, err := known.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(`
  INSERT OR IGNORE INTO words (word, interval) 
  VALUES (:word, :interval)`, words)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(`
  UPDATE words 
  SET interval=:interval, 
      updated_at=CURRENT_TIMESTAMP 
  WHERE word=":word"
  AND interval!=:interval`, words)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, word := range words {
		known.words[word.Word] = int(word.Interval)
	}
	return tx.Commit()
}

// TODO importCSVWords

func (known *KnownWords) isWellKnown(word string) bool {
	interval, ok := known.words[word]
	return ok && interval >= known.userSettings.KnownInterval
}

func (known *KnownWords) isKnown(word string) bool {
	_, ok := known.words[word]
	return ok
}

func (known *KnownWords) isKnownChar(char rune) bool {
	_, ok := known.characters[char]
	return ok
}

//go:embed assets/HSK
var hskWords embed.FS

func (known *KnownWords) GetUnknownHskWords(version string, level int) ([]UnknownWordEntry, error) {
	// ensure string == 2.0 or 3.0
	// ensure level == 1 - 6
	hskPath := path.Join(
		"assets",
		"HSK",
		version,
		fmt.Sprintf(`L%v.txt`, level),
	)
	rows := []UnknownWordEntry{}

	fileBytes, err := hskWords.ReadFile(hskPath)
	if err != nil {
		return rows, err
	}
	fileString := string(fileBytes)

	words := strings.Split(fileString, "\n")
	for _, word := range words {
		trimmed := strings.TrimSpace(word)
		trimmed = strings.Trim(trimmed, "\uFEFF")
		if !known.isKnown(trimmed) && len(trimmed) > 0 {
			// Its fine if occurance is just 0
			occurance, _ := known.frequency[trimmed]
			rows = append(rows, UnknownWordEntry{
				Word:      trimmed,
				Occurance: occurance,
			})
		}
	}
	return rows, nil
}
