package backend

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"time"

	"github.com/atselvan/ankiconnect"
	restError "github.com/privatesquare/bkst-go-utils/utils/errors"
	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

type Fields struct {
	Word        string      `json:"word"`
	Sentence    string      `json:"sentence,omitempty"`
	EnglishDefn string      `json:"englishDefn,omitempty"`
	ChineseDefn string      `json:"chineseDefn,omitempty"`
	Pinyin      string      `json:"pinyin,omitempty"`
	Images      []ImageInfo `json:"images,omitempty"`
}

type (
	AnkiInterface interface {
		CreateAnkiNote(fields Fields, tags []string) error
		GetAnkiNote(word int64) (RawAnkiNote, error)
		UpdateNoteFields(noteID int64, fields Fields) error
		UpdateSentenceAudio(noteId int64) error
		UpdateWordAudio(noteId int64) error
		UpdatePinyin(noteId int64) error
		ImportAnkiKeywords() error
		LoadProblemCards(query string) ([]ProblemCard, error)
		HealthCheck() error
		ConfigurationCheck() error
		LoadTemplate() error
		LoadModels() ([]string, error)
		LoadModelFields(model string) ([]string, error)
	}
)

type ankiInterface struct {
	backend *Backend
	anki    *ankiconnect.Client
}

func NewAnkiInterface(backend *Backend) *ankiInterface {
	return &ankiInterface{
		backend: backend,
		anki:    ankiconnect.NewClient(),
	}
}

type RawAnkiNote struct {
	NoteId int64  `json:"noteId"`
	Fields Fields `json:"fields"`
}

func (a *ankiInterface) HealthCheck() error {
	restErr := a.anki.Ping()
	if restErr != nil {
		return errors.New("Could not connect to Anki")
	}
	return nil
}

func (a *ankiInterface) ConfigurationCheck() error {
	// ActiveDeck needs to be a real deck
	decks, restErr := a.anki.Decks.GetAll()
	if restErr != nil {
		return toError(restErr)
	}
	if !slices.Contains(*decks, a.backend.UserSettings.AnkiConfig.ActiveDeck) {
		return errors.New("Active Deck does exist in Anki")
	}

	// ActiveModel needs to be a real model
	models, restErr := a.anki.Models.GetAll()
	if restErr != nil {
		return toError(restErr)
	}
	activeModel := a.backend.UserSettings.AnkiConfig.ActiveModel
	if !slices.Contains(*models, activeModel) {
		return errors.New("Chose Note type does exist in Anki")
	}

	// CardConfiguration needs to exists and have certian fields
	// and those fields all need to be real fields
	currentModelFields, restErr := a.anki.Models.GetFields(activeModel)
	if restErr != nil {
		return toError(restErr)
	}
	fieldsConfig := a.backend.UserSettings.GetMapping(activeModel)

	allEmpty := true
	for fieldName, value := range map[string]string{
		"FirstField":        fieldsConfig.FirstField,
		"Hanzi":             fieldsConfig.Hanzi,
		"ExampleSentence":   fieldsConfig.ExampleSentence,
		"EnglishDefinition": fieldsConfig.EnglishDefinition,
		"ChineseDefinition": fieldsConfig.ChineseDefinition,
		"Pinyin":            fieldsConfig.Pinyin,
		"HanziAudio":        fieldsConfig.HanziAudio,
		"SentenceAudio":     fieldsConfig.SentenceAudio,
		"Images":            fieldsConfig.Images,
		"Notes":             fieldsConfig.Notes,
	} {
		if value == "" {
			continue
		} else if !slices.Contains(*currentModelFields, value) {
			return fmt.Errorf(
				"Configured field for %v : (%v) does not exist",
				fieldName, value)
		} else {
			allEmpty = false
		}
	}

	if allEmpty {
		return errors.New("No fields have been configured for current model")
	}

	return nil
}

func (a *ankiInterface) getConfiguredMapping() (FieldsMapping, error) {

	currentModel := a.backend.UserSettings.AnkiConfig.ActiveModel
	currentSettings := a.backend.UserSettings.GetMapping(currentModel)
	// Required fields are:
	// Hanzi: yes
	// ExampleSentence: yes
	// EnglishDefinition <- one of these
	// ChineseDefinition <- one of these
	// HanziAudio: if genhanziaudio is set
	// SentenceAudio: if genAudio is set
	// Images: if imageApi is set
	// Notes: nope
	// For now just make sure hanzi is set
	if currentSettings.Hanzi == "" {
		return currentSettings, errors.New("Hanzi is not mapped")
	}

	return currentSettings, nil

}

func (a *ankiInterface) CreateAnkiNote(fields Fields, tags []string) error {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return err
	}

	restErr := a.anki.Ping()
	if restErr != nil {
		return toError(restErr)
	}

	audio := []ankiconnect.Audio{}
	pictures := []ankiconnect.Picture{}

	addAudio := func(field string, dest string) error {
		audio64, err := a.backend.TextToSpeech.Synthesize(field)
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
	if a.backend.UserSettings.AzureConfig.GenerateTermAudio {
		err := addAudio(fields.Word, currentMapping.HanziAudio)
		if err != nil {
			return err
		}
	}
	// dont generate if sentence is empty
	if a.backend.UserSettings.AzureConfig.GenerateSentenceAudio && len(fields.Sentence) > 0 {
		err := addAudio(fields.Sentence, currentMapping.SentenceAudio)
		if err != nil {
			return err
		}
	}
	for i, image := range fields.Images {
		// This should only contain valid image.Url type images
		if len(image.ImageData) > 0 {
			return errors.New("Somehow a base64 image got into anki fields")
		}
		if len(image.Url) == 0 {
			return errors.New("image.Url was empty")
		}
		milli := time.Now().UnixMilli()
		pictures = append(pictures, ankiconnect.Picture{
			URL: image.Url,
			// TODO dont guess the encoding format
			Filename: fmt.Sprintf("read-chinese-image-%v-%v.jpg", milli, i),
			Fields: []string{
				currentMapping.Images,
			},
		})
	}
	noteFields := ankiconnect.Fields{}

	if currentMapping.Hanzi != "" {
		noteFields[currentMapping.Hanzi] = fields.Word
	}

	if currentMapping.ExampleSentence != "" {
		noteFields[currentMapping.ExampleSentence] = fields.Sentence
	}

	if currentMapping.EnglishDefinition != "" {
		noteFields[currentMapping.EnglishDefinition] = fields.EnglishDefn
	}

	if currentMapping.ChineseDefinition != "" {
		noteFields[currentMapping.ChineseDefinition] = fields.ChineseDefn
	}

	if currentMapping.Pinyin != "" {
		noteFields[currentMapping.Pinyin] = fields.Pinyin
	}

	firstField := currentMapping.FirstField
	// If your first field is something else then ???
	if firstField != currentMapping.Hanzi ||
		firstField != currentMapping.ExampleSentence ||
		firstField != currentMapping.EnglishDefinition ||
		firstField != currentMapping.ChineseDefinition {
		noteFields[currentMapping.FirstField] = fmt.Sprint(time.Now().UnixMilli())
	}

	note := ankiconnect.Note{
		DeckName:  a.backend.UserSettings.AnkiConfig.ActiveDeck,
		ModelName: a.backend.UserSettings.AnkiConfig.ActiveModel,
		Fields:    noteFields,
		Tags:      tags,
		Options: &ankiconnect.Options{
			AllowDuplicate: a.backend.UserSettings.AnkiConfig.AllowDuplicates,
		},
		Audio:   audio,
		Picture: pictures,
	}
	restErr = a.anki.Notes.Add(note)
	if restErr != nil {
		return toError(restErr)
	}

	a.backend.KnownWords.AddWord(fields.Word, 0)

	return nil
}

func toError(restErr *restError.RestErr) error {
	return errors.New(fmt.Sprintf("%v, %v", restErr.Error, restErr.Message))
}

func getImageSize(base64String string) (width int, height int, err error) {
	// Decode base64 string into bytes
	imageBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return 0, 0, err
	}

	// Read image dimensions from byte array
	img, _, err := image.DecodeConfig(bytes.NewReader(imageBytes))
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

func (a *ankiInterface) GetAnkiNote(noteId int64) (RawAnkiNote, error) {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return RawAnkiNote{}, err
	}

	notes, restErr := a.anki.Notes.Get(fmt.Sprintf("nid:%v", noteId))
	if restErr != nil {
		return RawAnkiNote{}, toError(restErr)
	}
	if len(*notes) == 0 {
		return RawAnkiNote{}, errors.New("No note exists")
	}
	note := (*notes)[0]
	extract := func(field string) string {
		value, _ := note.Fields[field]
		return value.Value
	}

	images := []ImageInfo{}
	// Images will (hopefully) have html <img src=\"filename.jpg\">
	imagesString := extract(currentMapping.Images)
	if imagesString != "" {
		doc, err := html.Parse(strings.NewReader(imagesString))
		if err != nil {
			return RawAnkiNote{}, err
		}
		imageNames := []string{}
		//
		var f func(*html.Node)
		// This is definitly a little overkill but hey
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "img" {
				for _, a := range n.Attr {
					if a.Key == "src" {
						imageNames = append(imageNames, a.Val)
						break
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		for _, name := range imageNames {
			image, restErr := a.anki.Media.RetrieveMediaFile(name)
			if restErr != nil {
				return RawAnkiNote{}, toError(restErr)
			}
			width, height, err := getImageSize(*image)
			if err != nil {
				return RawAnkiNote{}, err
			}

			images = append(images, ImageInfo{
				ImageData:   *image,
				ImageWidth:  int64(width),
				ImageHeight: int64(height),
			})
		}

	}
	rawNote := RawAnkiNote{
		NoteId: note.NoteId,
		Fields: Fields{
			Word:        extract(currentMapping.Hanzi),
			Sentence:    extract(currentMapping.ExampleSentence),
			EnglishDefn: extract(currentMapping.EnglishDefinition),
			ChineseDefn: extract(currentMapping.ChineseDefinition),
			Pinyin:      extract(currentMapping.Pinyin),
			Images:      images,
		},
	}
	return rawNote, nil
}

func (a *ankiInterface) createAudio(text string, field string) (ankiconnect.Audio, error) {
	audio64, err := a.backend.TextToSpeech.Synthesize(text)
	if err != nil {
		return ankiconnect.Audio{}, err
	}
	audio := ankiconnect.Audio{
		Data: audio64,
		Filename: fmt.Sprintf("read-chinese-%v-%v.wav",
			field,
			time.Now().UnixMilli()),
		Fields: []string{
			field,
		},
	}
	return audio, nil
}

func (a *ankiInterface) createImage(url string, field string) ankiconnect.Picture {
	return ankiconnect.Picture{
		URL: url,
		Filename: fmt.Sprintf("read-chinese-%v-%v.wav",
			field,
			time.Now().UnixMilli()),
		Fields: []string{
			field,
		},
	}
}

func (a *ankiInterface) updateFieldAudio(noteId int64, sentence bool) error {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return err
	}

	note, err := a.GetAnkiNote(noteId)
	if err != nil {
		return err
	}
	var field string
	var text string
	if sentence {
		field = currentMapping.SentenceAudio
		text = note.Fields.Sentence
	} else {
		field = currentMapping.HanziAudio
		text = note.Fields.Word
	}

	audio, err := a.createAudio(text, field)
	if err != nil {
		return err
	}

	update := ankiconnect.UpdateNote{
		Id: noteId,
		Fields: ankiconnect.Fields{
			field: "",
		},
		Audio: []ankiconnect.Audio{audio},
	}

	restErr := a.anki.Notes.Update(update)
	if restErr != nil {
		return toError(restErr)
	}

	return nil
}

func (a *ankiInterface) UpdateSentenceAudio(noteId int64) error {
	return a.updateFieldAudio(noteId, true)
}

func (a *ankiInterface) UpdateWordAudio(noteId int64) error {
	return a.updateFieldAudio(noteId, false)
}

func (a *ankiInterface) UpdatePinyin(noteId int64) error {
	note, err := a.GetAnkiNote(noteId)
	if err != nil {
		return err
	}

	word := note.Fields.Word
	if word == "" {
		return errors.New("Hanzi not set for card")
	}
	pinyin := a.backend.Dictionaries.getPinyin(word)
	// TODO in this case generate it char by char?
	if pinyin == "" {
		return errors.New("Could not figure out pinyin")

	}

	return a.UpdateNoteFields(noteId, Fields{
		Word:   word,
		Pinyin: pinyin,
	})
}

func (a *ankiInterface) UpdateNoteFields(noteId int64, fields Fields) error {
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return err
	}

	update := ankiconnect.UpdateNote{
		Id:      noteId,
		Picture: []ankiconnect.Picture{},
	}

	// If there are no fields updated (ie only images added)
	// update will throw nonsense error, and so we need to
	// at least set the word (which should always be safe)
	ankiFields := ankiconnect.Fields{
		currentMapping.Hanzi: fields.Word,
	}

	if fields.Sentence != "" {
		// TODO: if there was a 'Sentence Translation' field of some sort
		// It needs to be changed. Also TODO generate our own translations?
		ankiFields[currentMapping.ExampleSentence] = fields.Sentence
		if a.backend.UserSettings.AzureConfig.GenerateSentenceAudio {
			audio, err := a.createAudio(fields.Sentence, currentMapping.SentenceAudio)
			if err != nil {
				return err
			}
			update.Audio = []ankiconnect.Audio{audio}
			// Clear the previous field data
			ankiFields[currentMapping.SentenceAudio] = ""
		}
	}
	if fields.EnglishDefn != "" {
		ankiFields[currentMapping.EnglishDefinition] = fields.EnglishDefn
	}
	if fields.ChineseDefn != "" {
		ankiFields[currentMapping.ChineseDefinition] = fields.ChineseDefn
	}
	if fields.Pinyin != "" {
		ankiFields[currentMapping.Pinyin] = fields.Pinyin
	}
	for _, image := range fields.Images {
		// Only upload images that have a url.
		// If they are base64 it means they should already be in anki?
		// TODO need to also delete old base64 images that have been removed
		if len(image.Url) > 0 {
			update.Picture = append(update.Picture,
				a.createImage(image.Url, currentMapping.Images))
		}
	}

	update.Fields = ankiFields

	restErr := a.anki.Notes.Update(update)
	if restErr != nil {
		return toError(restErr)
	}

	// TODO clear flag if needed
	return nil
}

func (a *ankiInterface) ImportAnkiKeywords() error {

	// Not able to do a good progress bar but this is
	// fine for now
	a.backend.setupProgress("Loading words from anki", 1)
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return err
	}
	cards, restErr := a.anki.Cards.Get(
		fmt.Sprintf(`"deck:%v"`, a.backend.UserSettings.AnkiConfig.ActiveDeck))
	if restErr != nil {
		return toError(restErr)
	}

	words := []WordEntry{}
	for _, card := range *cards {
		word, _ := card.Fields[currentMapping.Hanzi]
		words = append(words, WordEntry{
			Word:     word.Value,
			Interval: card.Interval,
		})
	}
	return a.backend.KnownWords.AddWords(words)
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
	Word     string   `json:"word"`
	Problems Problems `json:"problems"`
	Notes    string   `json:"notes"`
	NoteId   int64    `json:"noteId"`
}

func (a *ankiInterface) LoadProblemCards(query string) ([]ProblemCard, error) {
	problemCards := []ProblemCard{}
	currentMapping, err := a.getConfiguredMapping()
	if err != nil {
		return problemCards, err
	}
	currentDeck := a.backend.UserSettings.AnkiConfig.ActiveDeck
	currentNote := a.backend.UserSettings.AnkiConfig.ActiveModel

	trimmed := strings.TrimSpace(query)
	var prefixedQuery string
	if trimmed == "" {
		prefixedQuery = fmt.Sprintf(`"deck:%v" "note:%v"`, currentDeck, currentNote)
	} else {
		prefixedQuery = fmt.Sprintf(`"deck:%v" "note:%v" (%v)`,
			currentDeck, currentNote, trimmed)
	}

	notes, restErr := a.anki.Notes.Get(prefixedQuery)
	if restErr != nil {
		return nil, toError(restErr)
	}

	for _, note := range *notes {
		word, ok := note.Fields[currentMapping.Hanzi]
		if !ok {
			return nil, errors.New("Hanzi not found")
		}
		problemCard := ProblemCard{
			Word:     word.Value,
			Problems: Problems{},
			NoteId:   note.NoteId,
		}
		p := &problemCard.Problems

		// TODO do we need to mark
		// 	Setter: func(p *Problems) { p.Flagged = true },

		// Missing Example Sentence
		sentence, ok := note.Fields[currentMapping.ExampleSentence]
		if !ok || sentence.Value == "" {
			p.MissingSentence = true
		}

		// Missing Sentence Audio
		sentenceAudio, ok := note.Fields[currentMapping.SentenceAudio]
		if !p.MissingSentence && (!ok || sentenceAudio.Value == "") {
			p.MissingSentenceAudio = true
		}

		// Missing Image
		images, ok := note.Fields[currentMapping.Images]
		if !ok || images.Value == "" {
			p.MissingImage = true
		}

		// Missing HanziAudio
		hanziAudio, ok := note.Fields[currentMapping.HanziAudio]
		if !ok || hanziAudio.Value == "" {
			p.MissingWordAudio = true
		}

		// Missing Pinyin TODO: check for ugly pinyin? eg: ni3hao3
		pinyin, ok := note.Fields[currentMapping.Pinyin]
		if !ok || pinyin.Value == "" {
			p.MissingPinyin = true
		}

		notes, ok := note.Fields[currentMapping.Notes]
		if ok {
			problemCard.Notes = notes.Value
		}

		problemCards = append(problemCards, problemCard)
	}

	return problemCards, nil
}

func (a *ankiInterface) LoadTemplate() error {
	template, err := NewTemplateLoader().GetTemplate()
	if err != nil {
		return err
	}

	model := ankiconnect.Model{
		ModelName:     template.Name,
		InOrderFields: template.Fields,
		Css:           template.Css,
		IsCloze:       false,
		CardTemplates: []ankiconnect.CardTemplate{
			{
				Name:  template.Name,
				Front: template.Front,
				Back:  template.Back,
			},
		},
	}
	restErr := a.anki.Models.Create(model)
	if restErr != nil {
		return toError(restErr)
	}
	return nil
}

func (a *ankiInterface) LoadModels() ([]string, error) {
	models, restErr := a.anki.Models.GetAll()
	if restErr != nil {
		return nil, toError(restErr)
	}
	return *models, nil
}

func (a *ankiInterface) LoadDecks() ([]string, error) {
	decks, restErr := a.anki.Decks.GetAll()
	if restErr != nil {
		return nil, toError(restErr)
	}
	return *decks, nil
}

func (a *ankiInterface) LoadModelFields(model string) ([]string, error) {
	fields, restErr := a.anki.Models.GetFields(model)
	if restErr != nil {
		return nil, toError(restErr)
	}
	return *fields, nil

}

// TODO in client
// setSpecificValueOfCard
// and use that to unflag stuff when updated
