package backend

import (
	"errors"
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

type Generator struct {
	userSettings  *UserSettings
	segmentation  *Segmentation
	bookLibrary   BookLibrary
	known         *KnownWords
	sentenceCache map[string][]string
}

func NewGenerator(
	userSettings *UserSettings,
	s *Segmentation,
	b BookLibrary,
	known *KnownWords,
) *Generator {
	return &Generator{
		userSettings: userSettings,
		segmentation: s,
		bookLibrary:  b,
		known:        known,
	}
}

func (g *Generator) isT1Sentence(sentence []Token) bool {
	unknown := 0
	for _, token := range sentence {
		if token.IsWord && !g.known.isWellKnown(token.Data) {
			unknown += 1
		}
		if unknown > 1 {
			return false
		}
	}
	return true
}

func (g *Generator) isT1Candidate(sentence []Token, word string) bool {
	for _, token := range sentence {
		if token.IsWord && word != token.Data && !g.known.isWellKnown(token.Data) {
			return false
		}
	}
	return true
}

func tokensContains(sentence []Token, word string) bool {
	for _, token := range sentence {
		if word == token.Data {
			return true
		}
	}
	return true
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

// For now just save in memory pre segmented sentences that are
// t1 candidates.
//
// For best accuracy, this set would need to be updated whenever
// a new word is 'marked known'. Theoretically, the initial build and
// the update step could be different processes (as only sentences that
// contain the 'added' word would possibly become new sentences)
//
// This will be much more instensive computing at startup, since before
// we would be able to only do the expensive segment on sentences that contain
// the word we are looking for (using plain text matching). If we want, we could
// add certian parameters (eg, sentence length) which cut down on the amount of up
// front computing if this locks the cpu for a bit
func (g *Generator) GenerateSentenceTable() error {
	defer duration(track("Generate"))
	books, _ := g.bookLibrary.GetSomeBooks()
	g.sentenceCache = map[string][]string{}
	for _, book := range books {
		fullSegmented, err := GetSegmentedText(book)
		sentences := []string{}
		if err != nil {
			return err
		}
		for _, sentence := range fullSegmented {
			segmented := g.segmentation.SegmentSentence(sentence)
			if g.isT1Sentence(segmented) {
				sentences = append(sentences, sentence)
			}
		}
		g.sentenceCache[book.Title] = sentences
	}
	return nil
}

func (g *Generator) GetSentencesForWord(word string, bookIds []int64) ([]string, error) {
	defer duration(track("Get Sentences"))
	if len(g.sentenceCache) == 0 {
		// Just do this once up front in a simple way for now
		g.GenerateSentenceTable()
	}
	books, _ := g.bookLibrary.GetSomeBooks(bookIds...)
	sentences := []string{}
	for _, book := range books {
		t1Segmented, ok := g.sentenceCache[book.Title]
		if !ok {
			return sentences, errors.New("Failed to lookup book in cache")
		}
		for _, sentence := range t1Segmented {
			if strings.Contains(sentence, word) {
				segmented := g.segmentation.SegmentSentence(sentence)
				if tokensContains(segmented, word) && g.isT1Candidate(segmented, word) {
					sentences = append(sentences, sentence)
				}
			}
		}
	}
	idealLength := g.userSettings.IdealSentenceLength
	rankSentences(&sentences, idealLength)
	min := math.Min(float64(len(sentences)), float64(idealLength))
	sentences = sentences[0:int(min)]

	return sentences, nil
}

func rankSentences(sentences *[]string, idealLength int) {
	// Is a less than b?
	sort.Slice(*sentences, func(a int, b int) bool {
		// Which ever score is closer to 0 is better
		aScore := math.Abs(float64(len([]rune((*sentences)[a])) - idealLength))
		bScore := math.Abs(float64(len([]rune((*sentences)[b])) - idealLength))
		return (aScore <= bScore)
	})
}
