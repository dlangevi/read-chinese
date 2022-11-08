package backend

import (
	// "encoding/json"
	// "errors"
	// "net/http"
	// "net/url"
	//
	"github.com/atselvan/ankiconnect"
)

type Fields struct {
	Word        string   `json:"word"`
	Sentence    string   `json:"sentence"`
	EnglishDefn string   `json:"englishDefn"`
	ChineseDefn string   `json:"chineseDefn"`
	Pinyin      string   `json:"pinyin"`
	ImageUrls   []string `json:"imageUrls"`
}

const (
	ankiConnectUrl = "http://localhost:8765"
)

type AnkiInterface struct {
}

type RawAnkiNote struct {
	NoteId string `json:"noteId"`
	Fields Fields `json:"fields"`
}

func (AnkiInterface) GetAnkiNoteSkeleton(word string) RawAnkiNote {
	return RawAnkiNote{
		Fields: Fields{
			Word: word,
		},
	}
}

func (AnkiInterface) CreateAnkiCard(fields Fields, tags []string) string {
	client := ankiconnect.NewClient()
	restErr := client.Ping()
	if restErr != nil {
		return restErr.Error
	}
	note := ankiconnect.Note{
		DeckName:  "Reading",
		ModelName: "Reading Card",
		Fields: ankiconnect.Fields{
			"Hanzi":             fields.Word,
			"ExampleSentence":   fields.Sentence,
			"EnglishDefinition": fields.EnglishDefn,
			"ChineseDefinition": fields.ChineseDefn,
			"Pinyin":            fields.Pinyin,
		},
		Tags: append(tags, "read-chinese"),
		Options: &ankiconnect.Options{
			AllowDuplicate: true,
		},
		// audio
		// picture
	}
	restErr = client.Notes.Add(note)
	if restErr != nil {
		return restErr.Error
	}

	return "success"
}

// TODO in client
// findCards
// cardsInfo
// updateNoteFields
// setSpecificValueOfCard

// TODO ipc calls
// getAnkiNote,
// updateAnkiCard,
// loadFlaggedCards,
// importAnkiKeywords,
