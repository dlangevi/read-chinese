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

type AnkiInterface struct {
	anki         *ankiconnect.Client
	textToSpeech *TextToSpeech
	userSettings *UserSettings
	known        *KnownWords
}

func NewAnkiInterface(userSettings *UserSettings, known *KnownWords) *AnkiInterface {
	return &AnkiInterface{
		anki:         ankiconnect.NewClient(),
		textToSpeech: NewTextToSpeach(userSettings),
		userSettings: userSettings,
		known:        known,
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
		audio64, err := a.textToSpeech.Synthesize(field)
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
	if a.userSettings.GenerateTermAudio {
		err := addAudio(fields.Word, "HanziAudio")
		if err != nil {
			return err
		}
	}
	if a.userSettings.GenerateSentenceAudio {
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

	a.known.AddWord(fields.Word, 0)

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
	return a.known.AddWords(words)
}

type FlaggedCard struct {
	Word     string `json:"word"`
	Sentence string `json:"sentence"`
}

// Worry about using an enum or something later. I know this will work
// into Wails
type Problems struct {
	Flagged              bool `json:"Flagged"`
	MissingImage         bool `json:"MissingImage"`
	MissingSentence      bool `json:"MissingSentence"`
	MissingSentenceAudio bool `json:"MissingSentenceAudio"`
	MissingWordAudio     bool `json:"MissingWordAudio"`
	MissingPinyin        bool `json:"MissingPinyin"`
}

type ProblemCard struct {
	Word     string   `json:"Word"`
	Problems Problems `json:"Problems"`
}

func (a *AnkiInterface) LoadProblemCards() ([]ProblemCard, error) {
	problemCardsMap := map[int64]ProblemCard{}

	// For each query, get the noteIds. Add them all to the map
	// with problems

	type ProblemCase struct {
		Query  string
		Setter func(*Problems)
	}

	checks := []ProblemCase{
		{ // Cards Flagged by the user
			Query:  "-flag:0",
			Setter: func(p *Problems) { p.Flagged = true },
		},
		{ // Missing Example Sentence
			Query:  "ExampleSentence:",
			Setter: func(p *Problems) { p.MissingSentence = true },
		},
		{ // Has Sentence, but missing Sentence Audio
			Query:  "-ExampleSentence: SentenceAudio:",
			Setter: func(p *Problems) { p.MissingSentenceAudio = true },
		},
		{ // Missing Image
			Query:  "Images:",
			Setter: func(p *Problems) { p.MissingImage = true },
		},
		{ // Missing HanziAudio
			Query:  "HanziAudio:",
			Setter: func(p *Problems) { p.MissingWordAudio = true },
		},
		{ // Missing Pinyin TODO: check for ugly pinyin? eg: ni3hao3
			Query:  "Pinyin:",
			Setter: func(p *Problems) { p.MissingPinyin = true },
		},
	}

	for _, check := range checks {
		// Todo switch out Get for Search + Lookup if we want to speed it up
		ids, restErr := a.anki.Notes.Get(check.Query)
		if restErr != nil {
			return nil, toError(restErr)
		}

		for _, id := range *ids {
			word, ok := id.Fields["Hanzi"]
			if !ok {
				return nil, errors.New("Hanzi not found")
			}
			problemCard, exists := problemCardsMap[id.NoteId]
			if !exists {
				problemCard = ProblemCard{
					Word:     word.Value,
					Problems: Problems{},
				}
			}
			check.Setter(&problemCard.Problems)
			problemCardsMap[id.NoteId] = problemCard
		}

	}

	// Big lookup on all selected Note ids to map to the word in the field
	// For now just do

	// Finally map to array since wails cant do the needed ts stuff
	problemCards := []ProblemCard{}
	for _, problemCard := range problemCardsMap {
		problemCards = append(problemCards, problemCard)
	}

	return problemCards, nil
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
