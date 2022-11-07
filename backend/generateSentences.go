package backend

import (
	"math"
	"read-chinese/backend/segmentation"
	"sort"
	"strings"
)

type Generator struct {
}

func isT1Candidate(sentence []segmentation.Token, word string) bool {
	for _, token := range sentence {
		if token.IsWord && word != token.Data && !known.isWellKnown(token.Data) {
			return false
		}
	}
	return true
}

func tokensContains(sentence []segmentation.Token, word string) bool {
	for _, token := range sentence {
		if word == token.Data {
			return true
		}
	}
	return true
}

func (Generator) GetSentencesForWord(word string, bookIds []int64) ([]string, error) {
	books, _ := getBooks(bookIds...)
	sentences := []string{}
	for _, book := range books {
		fullSegmented, err := GetSegmentedText(book)
		if err != nil {
			return sentences, err
		}
		for _, sentence := range fullSegmented {
			if strings.Contains(sentence, word) {
				segmented := runtime.Segmentation.SegmentSentence(sentence)
				if tokensContains(segmented, word) && isT1Candidate(segmented, word) {
					sentences = append(sentences, sentence)
				}
			}
		}
	}
	rankSentences(&sentences)
	min := math.Min(float64(len(sentences)), 20)
	sentences = sentences[0:int(min)]

	return sentences, nil
}

func rankSentences(sentences *[]string) {
	// Is a less than b?
	sort.Slice(*sentences, func(a int, b int) bool {
		// Which ever score is closer to 0 is better
		aScore := math.Abs(float64(len([]rune((*sentences)[a])) - 20))
		bScore := math.Abs(float64(len([]rune((*sentences)[b])) - 20))
		return (aScore <= bScore)
	})
}
