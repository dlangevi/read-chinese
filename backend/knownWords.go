package backend

import (
	"embed"
	"fmt"
	"log"
	"path"
	"strings"
)

var known *KnownWords

type KnownWords struct {
	// For now just map word to interval
	words      map[string]int
	characters map[rune]bool
}

func NewKnownWords() *KnownWords {
	known = &KnownWords{}
	known.words = map[string]int{}
	known.characters = map[rune]bool{}
	known.syncWords()
	return known
}

func (known *KnownWords) syncWords() {
	type WordRow struct {
		Word     string
		Interval int
	}
	words := []WordRow{}
	err := Conn.Select(&words, `
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
func (known *KnownWords) AddWord(word string, age int) {
	known.words[word] = age
	Conn.MustExec(`
  INSERT OR IGNORE INTO words (word, interval) VALUES ($1, $2);
  UPDATE words 
  SET interval=$2, 
      updated_at = CURRENT_TIMESTAMP 
  WHERE word="$1"
  `, word, age)
}

type WordEntry struct {
	Word     string
	Interval int64
}

// TODO make this faster if possible
func (known *KnownWords) AddWords(words []WordEntry) error {
	tx, err := Conn.Beginx()
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
	return ok && interval >= userSettings.KnownInterval
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
		log.Println(trimmed, len(trimmed), word, len(word))

		if !known.isKnown(trimmed) && len(trimmed) > 0 {
			rows = append(rows, UnknownWordEntry{
				Word: trimmed,
			})
		}
	}
	return rows, nil
}
