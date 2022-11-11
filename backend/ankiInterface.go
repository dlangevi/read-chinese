package backend

import (
	"errors"
	"fmt"
	"time"

	"github.com/atselvan/ankiconnect"
	restError "github.com/privatesquare/bkst-go-utils/utils/errors"
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
	anki *ankiconnect.Client
}

func NewAnkiInterface() *AnkiInterface {
	return &AnkiInterface{
		ankiconnect.NewClient(),
	}
}

type RawAnkiNote struct {
	NoteId int64  `json:"noteId"`
	Fields Fields `json:"fields"`
}

func (a *AnkiInterface) GetAnkiNoteSkeleton(word string) RawAnkiNote {
	return RawAnkiNote{
		Fields: Fields{
			Word: word,
		},
	}
}

func (a *AnkiInterface) CreateAnkiNote(fields Fields, tags []string) error {
	restErr := a.anki.Ping()
	if restErr != nil {
		return toError(restErr)
	}

	audio := []ankiconnect.Audio{}
	pictures := []ankiconnect.Picture{}

	addAudio := func(field string, dest string) error {
		audio64, err := runtime.TextToSpeech.Synthesize(field)
		if err != nil {
			return err
		}
		audio = append(audio, ankiconnect.Audio{
			Data: audio64,
			Filename: fmt.Sprintf("read-chinese-%v-%v.wav",
				dest,
				time.Now().UnixMilli()),
			Fields: []string{
				dest,
			},
		})
		return nil
	}
	if userSettings.GenerateTermAudio {
		err := addAudio(fields.Word, "HanziAudio")
		if err != nil {
			return err
		}
	}
	if userSettings.GenerateSentenceAudio {
		err := addAudio(fields.Sentence, "SentenceAudio")
		if err != nil {
			return err
		}
	}
	for i, image := range fields.ImageUrls {
		milli := time.Now().UnixMilli()
		pictures = append(pictures, ankiconnect.Picture{
			URL: image,
			// TODO dont guess the encoding format
			Filename: fmt.Sprintf("read-chinese-image-%v-%v.jpg", milli, i),
			Fields: []string{
				"Images",
			},
		})
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
		Audio:   audio,
		Picture: pictures,
	}
	restErr = a.anki.Notes.Add(note)
	if restErr != nil {
		return toError(restErr)
	}

	known.AddWord(fields.Word, 0)

	return nil
}

func toError(restErr *restError.RestErr) error {
	return errors.New(fmt.Sprintf("%v, %v", restErr.Error, restErr.Message))
}

func (a *AnkiInterface) GetAnkiNote(word string) (RawAnkiNote, error) {
	notes, restErr := a.anki.Notes.Get(fmt.Sprintf("Hanzi:%v", word))
	if restErr != nil {
		return RawAnkiNote{}, toError(restErr)
	}
	if len(*notes) == 0 {
		return RawAnkiNote{}, errors.New("No note exists")
	}
	if len(*notes) > 1 {
		return RawAnkiNote{}, errors.New("Duplicate notes exists")
	}

	note := (*notes)[0]
	extract := func(field string) string {
		value, _ := note.Fields[field]
		return value.Value
	}

	rawNote := RawAnkiNote{
		NoteId: note.NoteId,
		Fields: Fields{
			Word:        extract("Hanzi"),
			Sentence:    extract("ExampleSentence"),
			EnglishDefn: extract("EnglishDefinition"),
			ChineseDefn: extract("ChineseDefinition"),
			Pinyin:      extract("Pinyin"),
			// TODO how to load the image from a card which already exists?
			// use retrieveMediaFile to get base64 encoded image
			// ImageUrls: [],
		},
		// rawNote,
	}
	return rawNote, nil
}

func (a *AnkiInterface) UpdateNoteFields(noteID int64, fields Fields) string {
	ankiFields := ankiconnect.Fields{}
	if fields.Word != "" {
		ankiFields["Hanzi"] = fields.Word
	}
	if fields.Sentence != "" {
		ankiFields["ExampleSentence"] = fields.Sentence
	}
	if fields.EnglishDefn != "" {
		ankiFields["EnglishDefinition"] = fields.EnglishDefn
	}
	if fields.ChineseDefn != "" {
		ankiFields["ChineseDefinition"] = fields.ChineseDefn
	}
	if fields.Pinyin != "" {
		ankiFields["Pinyin"] = fields.Pinyin
	}
	// TODO include new audio and images
	update := ankiconnect.UpdateNote{
		Id:     noteID,
		Fields: ankiFields,
	}
	restErr := a.anki.Notes.Update(update)
	if restErr != nil {
		return restErr.Error
	}

	// TODO clear flag if needed

	return "success"
}

func (a *AnkiInterface) ImportAnkiKeywords() error {
	cards, restErr := a.anki.Cards.Get("deck:Reading")
	if restErr != nil {
		return toError(restErr)
	}

	words := []WordEntry{}
	for _, card := range *cards {
		word, _ := card.Fields["Hanzi"]
		words = append(words, WordEntry{
			Word:     word.Value,
			Interval: card.Interval,
		})
	}
	return known.AddWords(words)
}

type FlaggedCard struct {
	Word     string `json:"word"`
	Sentence string `json:"sentence"`
}

func (a *AnkiInterface) LoadFlaggedCards() ([]FlaggedCard, error) {
	flaggedCards := []FlaggedCard{}
	cards, restErr := a.anki.Cards.Get("flag:1")
	if restErr != nil {
		return flaggedCards, toError(restErr)
	}
	for _, card := range *cards {
		word, _ := card.Fields["Hanzi"]
		sentence, _ := card.Fields["ExampleSentence"]
		flaggedCards = append(flaggedCards, FlaggedCard{
			Word:     word.Value,
			Sentence: sentence.Value,
		})
	}
	return flaggedCards, nil
}

// TODO in client
// setSpecificValueOfCard
// and use that to unflag stuff when updated
