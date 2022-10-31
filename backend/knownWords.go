package backend

import (
	"log"
)

type KnownWords struct {
	// For now just map word to interval
	words      map[string]int
	characters map[rune]bool
}

func NewKnownWords() *KnownWords {
	known := &KnownWords{}
	known.words = map[string]int{}
	known.characters = map[rune]bool{}
	known.loadWords()
	return known
}

func (known *KnownWords) loadWords() {
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

func (known *KnownWords) GetWords() map[string]int {
	return known.words
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
