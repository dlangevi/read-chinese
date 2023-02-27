package backend

import (
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

type Generator struct {
	userSettings    *UserConfig
	segmentation    *Segmentation
	bookLibrary     BookLibrary
	known           *KnownWords
	sentenceCache   map[string][]string
	cacheInProgress bool
	cacheComplete   bool
}

func NewGenerator(
	userSettings *UserConfig,
	s *Segmentation,
	b BookLibrary,
	known *KnownWords,
) *Generator {
	generator := Generator{
		userSettings:    userSettings,
		segmentation:    s,
		bookLibrary:     b,
		known:           known,
		sentenceCache:   map[string][]string{},
		cacheInProgress: false,
		cacheComplete:   false,
	}
	if userSettings.Meta.CacheSentences {
		go generator.GenerateSentenceTable()
	}
	return &generator
}

func (g *Generator) isT1Sentence(sentence []Token) bool {
	haventFoundAnyYet := true
	firstUnknown := ""
	for _, token := range sentence {
		if token.IsWord && !g.known.isWellKnown(token.Data) {
			if haventFoundAnyYet {
				haventFoundAnyYet = false
				firstUnknown = token.Data
				continue
			}
			if firstUnknown != token.Data {
				return false
			}
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
//
// Timing info (my computer, 21,919,827 characters across all books)
// Full index 45 seconds (running on single seperate thread)
// Seems to take up around 150mb of memory
// Searching indexed .5 seconds
// Previous method (quick culling + full search) 3 seconds,
// as the cache set is slowly loaded, this time slowly decreases eg
// generateSentences.go:74: Get Sentences: 3.150979092s
// generateSentences.go:74: Get Sentences: 2.58147285s
// generateSentences.go:74: Get Sentences: 2.398510062s
// generateSentences.go:74: Get Sentences: 1.806070047s
// generateSentences.go:74: Get Sentences: 1.790299864s
// generateSentences.go:74: Get Sentences: 604.083962ms
// generateSentences.go:74: Full Generate: 45.802128044s

func (g *Generator) GenerateSentenceTable() error {
	g.cacheInProgress = true
	defer duration(track("Full Generate"))
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
	g.cacheComplete = true
	return nil
}

// This can still be optimized a ton
func (g *Generator) GetSentencesForWord(word string, bookIds []int64) ([]string, error) {
	defer duration(track("Get Sentences"))
	if g.userSettings.Meta.CacheSentences {
		// If cache has not been completed
		if !g.cacheComplete {
			// Start it if this is the first time
			if !g.cacheInProgress {
				go g.GenerateSentenceTable()
			}
		}
	}

	books, _ := g.bookLibrary.GetSomeBooks(bookIds...)
	sentences := []string{}
	for _, book := range books {

		t1Segmented, ok := g.sentenceCache[book.Title]
		var err error = nil
		// If lookup fails, the book has not been processed yet
		if !ok {
			t1Segmented, err = GetSegmentedText(book)
			if err != nil {
				return sentences, err
			}
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
	idealLength := g.userSettings.SentenceGenerationConfig.IdealSentenceLength
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
