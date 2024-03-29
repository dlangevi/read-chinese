package backend

import (
	"errors"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

type Generator struct {
	sentenceCache map[string]SentenceCache
	cacheComplete bool
	mapLock       *sync.RWMutex
	cacheLock     *sync.Mutex
	bookCache     map[string][]TokenizedSentence
	backend       *Backend
}

type SentenceCache = map[string][]TokenizedSentence

func NewGenerator(
	backend *Backend,
) *Generator {
	generator := Generator{
		sentenceCache: map[string]SentenceCache{},
		cacheComplete: false,
		mapLock:       &sync.RWMutex{},
		cacheLock:     &sync.Mutex{},
		bookCache:     map[string][]TokenizedSentence{},
		backend:       backend,
	}
	go generator.GenerateSentenceTable()
	return &generator
}

func containsNew(sentence []Token, newWord string) bool {
	for _, token := range sentence {
		if token.Data == newWord {
			return true
		}
	}
	return false
}

func (g *Generator) isT1Sentence(sentence []Token) (bool, string) {
	haventFoundAnyYet := true
	firstUnknown := ""
	for _, token := range sentence {
		if token.IsWord && !g.backend.KnownWords.IsWellKnown(token.Data) {
			if haventFoundAnyYet {
				haventFoundAnyYet = false
				firstUnknown = token.Data
				continue
			}
			if firstUnknown != token.Data {
				return false, ""
			}
		}
	}
	return true, firstUnknown
}

func (g *Generator) passesKnownCheck(sentence []Token, word string) bool {
	foundWord := false
	for _, token := range sentence {
		if word == token.Data {
			foundWord = true
		}
		if token.IsWord && !g.backend.KnownWords.IsKnown(token.Data) {
			return false
		}
	}

	return foundWord
}

func (g *Generator) isT1Candidate(sentence []Token, word string) bool {
	for _, token := range sentence {
		if token.IsWord && word != token.Data && !g.backend.KnownWords.IsWellKnown(token.Data) {
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
	return false
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func (g *Generator) GenerateSentenceTableForBook(bookId int) int {
	book, err := g.backend.BookLibrary.GetBook(bookId)
	if err != nil {
		log.Println("Error loading:", bookId)
	}
	fullSegmented, err := GetSegmentedText(book)
	sentences := SentenceCache{}
	numSentences := 0
	if err != nil {
		log.Println("Error loading:", book.Title)
		return numSentences
	}
	for _, sentence := range fullSegmented {
		numSentences += 1
		if isT1, t1Word := g.isT1Sentence(sentence); isT1 {
			currentMap, ok := sentences[t1Word]
			if !ok {
				currentMap = []TokenizedSentence{}
			}
			currentMap = append(currentMap, sentence)
			sentences[t1Word] = currentMap
		}
	}
	g.mapLock.Lock()
	g.bookCache[book.Title] = fullSegmented
	g.sentenceCache[book.Title] = sentences
	g.mapLock.Unlock()

	return numSentences
}

// Timing info (my computer, 21,919,827 characters across all books)
// Full index 4 seconds (just throwing each segmentation into a seperate routine)
// Searching indexed .5 seconds
// Previous method (quick culling + full search) 3 seconds,

func (g *Generator) GenerateSentenceTable() error {
	g.cacheLock.Lock()
	defer g.cacheLock.Unlock()
	defer duration(track("Full Generate"))
	numSentences := 0
	books, _ := g.backend.BookLibrary.GetBooks()
	g.sentenceCache = map[string]SentenceCache{}
	var wg sync.WaitGroup
	for _, book := range books {
		wg.Add(1)
		go func(book Book) {
			defer wg.Done()
			numSentences += g.GenerateSentenceTableForBook(book.BookId)
		}(book)
	}
	wg.Wait()
	log.Println("Scanned", numSentences, "sentences")
	return nil
}

func (g *Generator) UpdateSentenceTable(newWord string) error {
	defer duration(track("New Word Generate"))
	books, _ := g.backend.BookLibrary.GetBooks()

	var wg sync.WaitGroup
	for _, book := range books {
		go func(book Book) {
			wg.Add(1)
			defer wg.Done()
			fullSegmented, ok := g.bookCache[book.Title]
			if !ok {
				log.Println("Error loading:", book.Title)
			}

			sentences := SentenceCache{}
			for _, sentence := range fullSegmented {
				if tokensContains(sentence, newWord) {
					if isT1, t1Word := g.isT1Sentence(sentence); isT1 {
						currentMap, ok := sentences[t1Word]
						if !ok {
							currentMap = []TokenizedSentence{}
						}
						currentMap = append(currentMap, sentence)
						sentences[t1Word] = currentMap
					}
				}
			}
			g.mapLock.Lock()
			previousCache, ok := g.sentenceCache[book.Title]
			if !ok {
				log.Println("Somehow book was missing during Update")
			}
			// The newWord is now learned, so we don't need its examples anymore
			delete(previousCache, newWord)
			for word, newSentences := range sentences {
				currentMap, ok := previousCache[word]
				if !ok {
					currentMap = []TokenizedSentence{}
				}
				currentMap = append(currentMap, newSentences...)
				previousCache[word] = currentMap
			}
			g.sentenceCache[book.Title] = previousCache
			g.mapLock.Unlock()
		}(book)
	}
	wg.Wait()
	return nil
}

type Sentence struct {
	Sentence string `json:"sentence"`
	Source   string `json:"source,omitempty"`
}

// This can still be optimized a ton
func (g *Generator) GetSentencesForWord(word string, bookIds []int) ([]Sentence, error) {
	defer duration(track("Get Sentences"))
	g.cacheLock.Lock()
	defer g.cacheLock.Unlock()

	books, _ := g.backend.BookLibrary.GetBooks(bookIds...)
	sentences := []Sentence{}

	if g.backend.KnownWords.IsKnown(word) {
		// Have to do a slower lookup in completly known sentences
		for _, book := range books {
			bookSentences, ok := g.bookCache[book.Title]
			if !ok {
				// Here we are in a weird state where a book has been added but not
				// processed
				return nil, errors.New("Book missing from sentenceCache, please restart")
			}
			for _, sentence := range bookSentences {
				if g.passesKnownCheck(sentence, word) {
					sentences = append(sentences, Sentence{
						Sentence: toString(sentence),
						Source:   book.Title,
					})
				}
			}
		}
	} else {

		for _, book := range books {
			g.mapLock.RLock()
			t1Segmented, ok := g.sentenceCache[book.Title]
			g.mapLock.RUnlock()
			if !ok {
				// Here we are in a weird state where a book has been added but not
				// processed
				return nil, errors.New("Book missing from sentenceCache, please restart")
			}
			t1Sentences, ok := t1Segmented[word]
			if ok {
				for _, sentence := range t1Sentences {
					sentences = append(sentences, Sentence{
						Sentence: toString(sentence),
						Source:   book.Title,
					})
				}
			}
		}
	}
	idealLength := g.backend.UserSettings.SentenceGenerationConfig.IdealSentenceLength
	rankSentences(sentences, idealLength)
	min := math.Min(float64(len(sentences)), 8)
	sentences = sentences[0:int(min)]

	return sentences, nil
}

func rankSentences(sentences []Sentence, idealLength int) {
	// Is a less than b?
	sort.Slice(sentences, func(a int, b int) bool {
		// Which ever score is closer to 0 is better
		aScore := math.Abs(float64(len([]rune(sentences[a].Sentence)) - idealLength))
		bScore := math.Abs(float64(len([]rune(sentences[b].Sentence)) - idealLength))
		return (aScore <= bScore)
	})
}
