package core

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type WordRow struct {
	Word      string `json:"word" db:"word"`
	Occurance int    `json:"occurance" db:"occurance"`
}

func LearningTarget(db *sqlx.DB) ([]WordRow, error) {
	words := []WordRow{}
	err := db.Select(&words, `
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
		return nil, err
	}

	return words, nil
}
