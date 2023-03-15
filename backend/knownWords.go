package backend

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WordData struct {
	Interval  int64
	LearnedOn time.Time
}

type KnownWords struct {
	// For now just map word to interval
	words      map[string]WordData
	characters map[rune]bool
	frequency  map[string]int64
	db         *sqlx.DB
	backend    *Backend
}

func NewKnownWords(db *sqlx.DB,
	backend *Backend,
) *KnownWords {
	known := &KnownWords{
		words:      map[string]WordData{},
		characters: map[rune]bool{},
		frequency:  map[string]int64{},
		db:         db,
		backend:    backend,
	}
	known.syncWords()
	known.syncFrequency()
	return known
}

func (known *KnownWords) syncWords() {
	type WordRow struct {
		Word      string
		Interval  int64
		LearnedOn time.Time `db:"created_at"`
	}
	words := []WordRow{}
	err := known.db.Select(&words, `
    SELECT word, interval, created_at
    FROM words
  `)
	if err != nil {
		log.Fatal("Failed to load words", err)
	}
	for _, word := range words {
		known.words[word.Word] = WordData{
			Interval:  word.Interval,
			LearnedOn: word.LearnedOn,
		}
		for _, char := range word.Word {
			known.characters[char] = true
		}
	}
}

func (known *KnownWords) syncFrequency() {
	type WordRow struct {
		Word  string
		Count int64
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
func (known *KnownWords) AddWord(word string, age int64) error {
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
	known.words[word] = WordData{
		Interval:  age,
		LearnedOn: time.Now(),
	}
	err = tx.Commit()
	if known.backend.ctx != nil {
		runtime.EventsEmit(known.backend.ctx, "AddedWord", word)
	}
	return err
}

type WordEntry struct {
	Word     string `db:"word"`
	Interval int64  `db:"interval"`
}

// TODO make this faster if possible
// TODO chunk in batches of 5000 words
func (known *KnownWords) AddWords(words []WordEntry) error {
	newWords := []WordEntry{}
	needsUpdate := []WordEntry{}
	alreadySeen := map[string]int64{}

	for _, word := range words {
		// Doing this ensures no duplicates are added
		_, seen := alreadySeen[word.Word]
		if seen {
			log.Println("Duplicate words: ", word.Word)
		} else {
			alreadySeen[word.Word] = word.Interval
			if !known.isKnown(word.Word) {
				newWords = append(newWords, word)
			} else {
				currentInterval := known.words[word.Word]
				if currentInterval.Interval != word.Interval {
					needsUpdate = append(needsUpdate, word)
				}
			}
		}
	}

	tx, err := known.db.Beginx()
	if err != nil {
		return err
	}
	if len(newWords) > 0 {
		_, err := tx.NamedExec(`
  INSERT INTO words (word, interval) 
  VALUES (:word, :interval)`, newWords)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	stmt, err := tx.Preparex(`
    UPDATE words 
    SET interval=$1, 
        updated_at=CURRENT_TIMESTAMP 
    WHERE word=$2`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, word := range needsUpdate {
		_, err := stmt.Exec(word.Interval, word.Word)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, word := range newWords {
		known.words[word.Word] = WordData{
			Interval:  word.Interval,
			LearnedOn: time.Now(),
		}
		for _, char := range word.Word {
			known.characters[char] = true
		}
	}

	for _, word := range needsUpdate {
		wordData := known.words[word.Word]
		wordData.Interval = word.Interval
		for _, char := range word.Word {
			known.characters[char] = true
		}
	}
	return nil
}

func (known *KnownWords) ImportCSVWords(path string) error {
	csvFile, err := os.Open(path)
	if err != nil {
		return err
	}
	r := csv.NewReader(csvFile)
	r.FieldsPerRecord = -1
	words := []WordEntry{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) == 0 {
			return fmt.Errorf("Record had 0 elements")
		}

		if len(record) == 1 {
			words = append(words, WordEntry{
				Word:     record[0],
				Interval: 10000,
			})
		} else {
			// record also has additional data
			interval, err := strconv.ParseInt(record[1], 10, 64)
			if err != nil {
				log.Println("Error: failed to parse int from", record[1])
				// fall back to interval of 10000
				words = append(words, WordEntry{
					Word:     record[0],
					Interval: 10000,
				})
			}

			words = append(words, WordEntry{
				Word:     record[0],
				Interval: interval,
			})
		}
	}
	return known.AddWords(words)
}

func (known *KnownWords) isWellKnown(word string) bool {
	interval, ok := known.words[word]
	return ok && interval.Interval >= int64(
		known.backend.UserSettings.SentenceGenerationConfig.KnownInterval)
}

func (known *KnownWords) isKnown(word string) bool {
	_, ok := known.words[word]
	return ok
}

func (known *KnownWords) isKnownChar(char rune) bool {
	_, ok := known.characters[char]
	return ok
}

func (known *KnownWords) GetOccurances(words []string) map[string]int64 {
	occurances := map[string]int64{}
	for _, word := range words {
		// Its fine if occurance is just 0
		occurance, _ := known.frequency[word]
		occurances[word] = occurance
	}
	return occurances
}

//go:embed assets/HSK
var hskWords embed.FS

func (known *KnownWords) GetUnknownHskWords(version string, level int) ([]string, error) {
	// ensure string == 2.0 or 3.0
	// ensure level == 1 - 6
	hskPath := path.Join(
		"assets",
		"HSK",
		version,
		fmt.Sprintf(`L%v.txt`, level),
	)
	rows := []string{}

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
			rows = append(rows, trimmed)
		}
	}
	return rows, nil
}

//	type StatsInfo {
//	  WordsStats
//
// }
type TimeQuery struct {
	Day             string `json:"day"`
	Known           int    `json:"known"`
	KnownCharacters int    `json:"knownCharacters"`
}

func (known *KnownWords) GetStatsInfo() ([]TimeQuery, error) {
	charMap := map[rune]time.Time{}
	dateMap := map[string]int{}
	dateCharMap := map[string]int{}

	for word, data := range known.words {
		for _, c := range word {
			prevTime, ok := charMap[c]
			if !ok {
				charMap[c] = data.LearnedOn
			} else if prevTime.After(data.LearnedOn) {
				charMap[c] = data.LearnedOn
			}
		}
		fmtTime := data.LearnedOn.Format("2006年01月02日")
		current := dateMap[fmtTime]
		dateMap[fmtTime] = current + 1
	}

	for _, learnedOn := range charMap {
		fmtTime := learnedOn.Format("2006年01月02日")
		current := dateCharMap[fmtTime]
		dateCharMap[fmtTime] = current + 1
	}

	times := []TimeQuery{}
	for date, num := range dateMap {
		knownChar := dateCharMap[date]
		times = append(times, TimeQuery{
			Day:             date,
			Known:           num,
			KnownCharacters: knownChar,
		})
	}
	sort.Slice(times, func(a, b int) bool {
		return (times[a].Day <= times[b].Day)
	})

	cumulative := 0
	cumulativeChar := 0
	for i := range times {
		cumulative += times[i].Known
		cumulativeChar += times[i].KnownCharacters
		times[i].Known = cumulative
		times[i].KnownCharacters = cumulativeChar
	}
	return times, nil
}
