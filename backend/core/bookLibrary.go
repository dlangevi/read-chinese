package core

import (
	"database/sql"
	"log"
)

type WordRow struct {
	Word      string `json:"word"`
	Occurance int    `json:"occurance"`
}

func LearningTarget(db *sql.DB) ([]WordRow, error) {
	rows, err := db.Query(`
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
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var words []WordRow
	for rows.Next() {
		var word WordRow
		if err := rows.Scan(&word.Word, &word.Occurance); err != nil {
			return words, err
		}
		log.Println(word)
		words = append(words, word)
	}
	if err = rows.Err(); err != nil {
		return words, err

	}
	return words, nil
}
