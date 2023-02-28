package backend

import (
	"os"
	"path"
	"testing"
)

func TestRank(t *testing.T) {
	sentences := []Sentence{
		{Sentence: "len = 30 012345678901234567890"},
		{Sentence: "len = 22 0123456789012"},
		{Sentence: "len = 23 01234567890123"},
		{Sentence: "len = 17 01234567"},
		{Sentence: "len = 15 012345"},
		{Sentence: "len = 28 0123456789012345678"},
		{Sentence: "len = 9  "},
		{Sentence: "len = 20 01234567890"},
	}

	rankSentences(sentences, 20)
	sentences = sentences[0:3]

	if len(sentences[0].Sentence) != 20 || len(sentences[1].Sentence) != 22 || len(sentences[2].Sentence) != 17 {

		t.Errorf("sentences had the wrong value %v", sentences)
	}
}

func TestGenerate(t *testing.T) {
	tempDb := path.Join(os.TempDir(), "generate.db")
	os.Remove(tempDb)
	defer os.Remove(tempDb)
	myBackend := createBackend(tempDb)
	bookIds := []int64{1}
	_, err := myBackend.BookLibrary.AddBook("张天翼", "秃秃大王", "fake.jpg", "testdata/example_book.txt")
	myBackend.Generator.GenerateSentenceTable()
	if err != nil {
		t.Error("Failed to insert book snip", err)
	}
	sentences, err := myBackend.Generator.GetSentencesForWord("么", bookIds)
	if err != nil {
		t.Error("Failed to get sentences", err)
	}
	if len(sentences) != 0 {
		t.Error("Too many sentences", sentences)
	}
	// Try to load this sentence
	// 真  的  么  ？
	myBackend.KnownWords.AddWord("真", 100)
	// Not well known enough
	myBackend.KnownWords.AddWord("的", 1)
	myBackend.Generator.GenerateSentenceTable()
	sentences, _ = myBackend.Generator.GetSentencesForWord("么", bookIds)
	if len(sentences) != 0 {
		t.Error("Too many sentences", sentences)
	}

	// Now it is
	myBackend.KnownWords.AddWord("的", 100)
	myBackend.Generator.GenerateSentenceTable()
	sentences, _ = myBackend.Generator.GetSentencesForWord("么", bookIds)
	if len(sentences) != 1 {
		t.Error("Not enough sentences", sentences)
	}

	// Try to load this sentence
	// 	有人  说  ，  流星  就是  这么  来  的  。
	myBackend.KnownWords.AddWords([]WordEntry{
		{"有人", 100},
		{"说", 100},
		{"流星", 100},
		{"就是", 100},
		{"来", 100},
	})
	myBackend.Generator.GenerateSentenceTable()
	sentences, _ = myBackend.Generator.GetSentencesForWord("这么", bookIds)
	if len(sentences) != 1 {
		t.Error("Not enough sentences", sentences)
	}
}
