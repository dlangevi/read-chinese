package backend

import (
	"read-chinese/backend/segmentation"
)

func isT1Candidate(sentence []segmentation.Token, word string) bool {
	for _, token := range sentence {
		if token.IsWord && word != token.Data && !known.isWellKnown(token.Data) {
			return false
		}
	}
	return true
}

func GetSentencesForWord(word string, bookIds ...int64) {
	// books, _ := getBooks(bookIds...)
	//
	// candidates := map[string]bool{}
	// for _, book := range books {
	//   fullSegmented := BookLibrary.GetSegmentedText(book)
	//
	// }
	//
	// tokens := runtime.Segmentation.SegmentSentence("我是你")
}
