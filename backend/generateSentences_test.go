package backend

import (
	"testing"
)

func TestRank(t *testing.T) {
	sentences := []string{
		"len = 30 012345678901234567890",
		"len = 22 0123456789012",
		"len = 23 01234567890123",
		"len = 17 01234567",
		"len = 15 012345",
		"len = 28 0123456789012345678",
		"len = 9  ",
		"len = 20 01234567890",
	}

	rankSentences(&sentences)
	sentences = sentences[0:3]

	if len(sentences[0]) != 20 || len(sentences[1]) != 22 || len(sentences[2]) != 17 {

		t.Errorf("sentences had the wrong value %v", sentences)
	}

}
