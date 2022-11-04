package backend

import (
	"log"
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
  INSERT INTO words (word, interval) VALUES ($1, $2)
    ON CONFLICT(word) DO UPDATE SET 
      interval=excluded.interval
  `, word, age)
}

// TODO importCSVWords
// TODO updateInterval?

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
