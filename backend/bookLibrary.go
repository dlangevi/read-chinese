package backend

import (
	"log"
)

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
